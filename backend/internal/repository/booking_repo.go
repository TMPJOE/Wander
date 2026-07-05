package repository

import (
	"context"
	"fmt"
	"strings"

	"wander/backend/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type BookingRepository interface {
	Create(ctx context.Context, b models.Booking) (*models.Booking, error)
	GetByID(ctx context.Context, id int) (*models.Booking, error)
	ListByUserID(ctx context.Context, userID int) ([]models.Booking, error)
	ListByGuideID(ctx context.Context, guideID int) ([]models.Booking, error)
	UpdateStatus(ctx context.Context, id int, status string) error
	UpdatePayment(ctx context.Context, id int, intentID string, status string) error
}

type PgBookingRepository struct {
	pool *pgxpool.Pool
}

func NewPgBookingRepository(pool *pgxpool.Pool) BookingRepository {
	return &PgBookingRepository{pool: pool}
}

func (r *PgBookingRepository) Create(ctx context.Context, b models.Booking) (*models.Booking, error) {
	query := `
		INSERT INTO bookings (user_id, schedule_id, tour_id, guest_count, total_price, status, notes)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, user_id, schedule_id, tour_id, guest_count, total_price, status, notes, created_at, updated_at
	`
	res := &models.Booking{}
	err := r.pool.QueryRow(ctx, query, b.UserID, b.ScheduleID, b.TourID, b.GuestCount, b.TotalPrice, b.Status, b.Notes).Scan(
		&res.ID, &res.UserID, &res.ScheduleID, &res.TourID, &res.GuestCount, &res.TotalPrice, &res.Status, &res.Notes, &res.CreatedAt, &res.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("create booking: %w", err)
	}
	return res, nil
}

func (r *PgBookingRepository) GetByID(ctx context.Context, id int) (*models.Booking, error) {
	query := `
		SELECT b.id, b.user_id, b.schedule_id, b.tour_id, b.guest_count, b.total_price, b.status, b.notes, b.created_at, b.updated_at,
		       b.payment_status, b.stripe_payment_intent_id,
		       t.title as tour_title, t.location as tour_location, t.images as tour_images,
		       u.first_name || ' ' || u.last_name as guide_name, u.avatar_url as guide_avatar,
		       tu.first_name || ' ' || tu.last_name as user_name,
		       s.start_time as schedule_start, s.end_time as schedule_end
		FROM bookings b
		JOIN tours t ON b.tour_id = t.id
		JOIN users u ON t.guide_id = u.id
		JOIN users tu ON b.user_id = tu.id
		JOIN tour_schedules s ON b.schedule_id = s.id
		WHERE b.id = $1
	`
	b := &models.Booking{}
	var tourImages []byte
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&b.ID, &b.UserID, &b.ScheduleID, &b.TourID, &b.GuestCount, &b.TotalPrice, &b.Status, &b.Notes, &b.CreatedAt, &b.UpdatedAt,
		&b.PaymentStatus, &b.StripePaymentIntentID,
		&b.TourTitle, &b.TourLocation, &tourImages, &b.GuideName, &b.GuideAvatar, &b.UserName, &b.ScheduleStart, &b.ScheduleEnd,
	)
	if err != nil {
		return nil, models.ErrNotFound
	}
	b.TourImage = parseFirstImage(tourImages)
	return b, nil
}

func (r *PgBookingRepository) ListByUserID(ctx context.Context, userID int) ([]models.Booking, error) {
	query := `
		SELECT b.id, b.user_id, b.schedule_id, b.tour_id, b.guest_count, b.total_price, b.status, b.notes, b.created_at, b.updated_at,
		       t.title as tour_title, t.location as tour_location, t.images as tour_images,
		       u.first_name || ' ' || u.last_name as guide_name, u.avatar_url as guide_avatar,
		       s.start_time as schedule_start, s.end_time as schedule_end
		FROM bookings b
		JOIN tours t ON b.tour_id = t.id
		JOIN users u ON t.guide_id = u.id
		JOIN tour_schedules s ON b.schedule_id = s.id
		WHERE b.user_id = $1
		ORDER BY s.start_time DESC
	`
	rows, err := r.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("list bookings user: %w", err)
	}
	defer rows.Close()

	var bookings []models.Booking
	for rows.Next() {
		var b models.Booking
		var tourImages []byte
		err := rows.Scan(
			&b.ID, &b.UserID, &b.ScheduleID, &b.TourID, &b.GuestCount, &b.TotalPrice, &b.Status, &b.Notes, &b.CreatedAt, &b.UpdatedAt,
			&b.TourTitle, &b.TourLocation, &tourImages, &b.GuideName, &b.GuideAvatar, &b.ScheduleStart, &b.ScheduleEnd,
		)
		if err != nil {
			return nil, err
		}
		b.TourImage = parseFirstImage(tourImages)
		bookings = append(bookings, b)
	}
	return bookings, nil
}

func (r *PgBookingRepository) ListByGuideID(ctx context.Context, guideID int) ([]models.Booking, error) {
	query := `
		SELECT b.id, b.user_id, b.schedule_id, b.tour_id, b.guest_count, b.total_price, b.status, b.notes, b.created_at, b.updated_at,
		       t.title as tour_title, t.location as tour_location, t.images as tour_images,
		       tu.first_name || ' ' || tu.last_name as user_name, tu.avatar_url as user_avatar,
		       s.start_time as schedule_start, s.end_time as schedule_end
		FROM bookings b
		JOIN tours t ON b.tour_id = t.id
		JOIN users tu ON b.user_id = tu.id
		JOIN tour_schedules s ON b.schedule_id = s.id
		WHERE t.guide_id = $1
		ORDER BY s.start_time DESC
	`
	rows, err := r.pool.Query(ctx, query, guideID)
	if err != nil {
		return nil, fmt.Errorf("list bookings guide: %w", err)
	}
	defer rows.Close()

	var bookings []models.Booking
	for rows.Next() {
		var b models.Booking
		var tourImages []byte
		err := rows.Scan(
			&b.ID, &b.UserID, &b.ScheduleID, &b.TourID, &b.GuestCount, &b.TotalPrice, &b.Status, &b.Notes, &b.CreatedAt, &b.UpdatedAt,
			&b.TourTitle, &b.TourLocation, &tourImages, &b.UserName, &b.GuideAvatar,
			&b.ScheduleStart, &b.ScheduleEnd,
		)
		if err != nil {
			return nil, err
		}
		b.TourImage = parseFirstImage(tourImages)
		bookings = append(bookings, b)
	}
	return bookings, nil
}

func (r *PgBookingRepository) UpdateStatus(ctx context.Context, id int, status string) error {
	_, err := r.pool.Exec(ctx, "UPDATE bookings SET status = $1 WHERE id = $2", status, id)
	if err != nil {
		return fmt.Errorf("update booking status: %w", err)
	}
	return nil
}

func (r *PgBookingRepository) UpdatePayment(ctx context.Context, id int, intentID string, status string) error {
	_, err := r.pool.Exec(ctx, "UPDATE bookings SET payment_status = $1, stripe_payment_intent_id = $2 WHERE id = $3", status, intentID, id)
	if err != nil {
		return fmt.Errorf("update booking payment: %w", err)
	}
	return nil
}

func parseFirstImage(imagesRaw []byte) string {
	if len(imagesRaw) == 0 {
		return ""
	}
	s := string(imagesRaw)
	s = strings.Trim(s, "[]\" ")
	parts := strings.Split(s, "\",\"")
	if len(parts) > 0 && parts[0] != "" {
		return parts[0]
	}
	return ""
}
