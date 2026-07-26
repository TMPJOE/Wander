package repository

import (
	"context"
	"fmt"
	"time"

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
	GetEarnings(ctx context.Context, guideID int, period string) (*models.GuideEarnings, error)
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
		       t.title as tour_title, t.location as tour_location, t.meeting_point, t.latitude as tour_latitude, t.longitude as tour_longitude,
		       (SELECT url FROM tour_images WHERE tour_id = t.id ORDER BY position ASC LIMIT 1) as tour_image,
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
	var tourImage *string
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&b.ID, &b.UserID, &b.ScheduleID, &b.TourID, &b.GuestCount, &b.TotalPrice, &b.Status, &b.Notes, &b.CreatedAt, &b.UpdatedAt,
		&b.PaymentStatus, &b.StripePaymentIntentID,
		&b.TourTitle, &b.TourLocation, &b.MeetingPoint, &b.TourLatitude, &b.TourLongitude, &tourImage, &b.GuideName, &b.GuideAvatar, &b.UserName, &b.ScheduleStart, &b.ScheduleEnd,
	)
	if err != nil {
		return nil, models.ErrNotFound
	}
	if tourImage != nil {
		b.TourImage = *tourImage
	}
	return b, nil
}

func (r *PgBookingRepository) ListByUserID(ctx context.Context, userID int) ([]models.Booking, error) {
	query := `
		SELECT b.id, b.user_id, b.schedule_id, b.tour_id, b.guest_count, b.total_price, b.status, b.notes, b.created_at, b.updated_at,
		       t.title as tour_title, t.location as tour_location,
		       (SELECT url FROM tour_images WHERE tour_id = t.id ORDER BY position ASC LIMIT 1) as tour_image,
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
		var tourImage *string
		err := rows.Scan(
			&b.ID, &b.UserID, &b.ScheduleID, &b.TourID, &b.GuestCount, &b.TotalPrice, &b.Status, &b.Notes, &b.CreatedAt, &b.UpdatedAt,
			&b.TourTitle, &b.TourLocation, &tourImage, &b.GuideName, &b.GuideAvatar, &b.ScheduleStart, &b.ScheduleEnd,
		)
		if err != nil {
			return nil, err
		}
		if tourImage != nil {
			b.TourImage = *tourImage
		}
		bookings = append(bookings, b)
	}
	return bookings, nil
}

func (r *PgBookingRepository) ListByGuideID(ctx context.Context, guideID int) ([]models.Booking, error) {
	query := `
		SELECT b.id, b.user_id, b.schedule_id, b.tour_id, b.guest_count, b.total_price, b.status, b.notes, b.created_at, b.updated_at,
		       t.title as tour_title, t.location as tour_location,
		       (SELECT url FROM tour_images WHERE tour_id = t.id ORDER BY position ASC LIMIT 1) as tour_image,
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
		var tourImage *string
		err := rows.Scan(
			&b.ID, &b.UserID, &b.ScheduleID, &b.TourID, &b.GuestCount, &b.TotalPrice, &b.Status, &b.Notes, &b.CreatedAt, &b.UpdatedAt,
			&b.TourTitle, &b.TourLocation, &tourImage, &b.UserName, &b.GuideAvatar,
			&b.ScheduleStart, &b.ScheduleEnd,
		)
		if err != nil {
			return nil, err
		}
		if tourImage != nil {
			b.TourImage = *tourImage
		}
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

func (r *PgBookingRepository) GetEarnings(ctx context.Context, guideID int, period string) (*models.GuideEarnings, error) {
	var since time.Time
	var filterTime bool

	switch period {
	case "week":
		since = time.Now().AddDate(0, 0, -7)
		filterTime = true
	case "month":
		since = time.Now().AddDate(0, -1, 0)
		filterTime = true
	case "year":
		since = time.Now().AddDate(-1, 0, 0)
		filterTime = true
	}

	// 1. Calculate totals
	totalsQuery := `
		SELECT
			COALESCE(SUM(b.total_price) FILTER (WHERE b.payment_status = 'authorized'), 0.00) as total_authorized,
			COALESCE(SUM(b.total_price) FILTER (WHERE b.payment_status = 'paid'), 0.00) as total_paid
		FROM bookings b
		JOIN tours t ON b.tour_id = t.id
		WHERE t.guide_id = $1
		  AND ($2 = false OR b.created_at >= $3)
	`
	earnings := &models.GuideEarnings{ByTour: []models.TourEarning{}}
	err := r.pool.QueryRow(ctx, totalsQuery, guideID, filterTime, since).Scan(
		&earnings.TotalAuthorized, &earnings.TotalPaid,
	)
	if err != nil {
		return nil, fmt.Errorf("calculate earnings totals: %w", err)
	}

	// 2. Calculate by-tour breakdown
	breakdownQuery := `
		SELECT 
			t.id as tour_id,
			t.title as tour_title,
			COUNT(b.id) as bookings,
			COALESCE(SUM(b.total_price), 0.00) as revenue,
			b.payment_status as status
		FROM bookings b
		JOIN tours t ON b.tour_id = t.id
		WHERE t.guide_id = $1
		  AND b.payment_status IN ('authorized', 'paid')
		  AND ($2 = false OR b.created_at >= $3)
		GROUP BY t.id, t.title, b.payment_status
		ORDER BY revenue DESC
	`
	rows, err := r.pool.Query(ctx, breakdownQuery, guideID, filterTime, since)
	if err != nil {
		return nil, fmt.Errorf("query earnings breakdown: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var te models.TourEarning
		err := rows.Scan(&te.TourID, &te.TourTitle, &te.Bookings, &te.Revenue, &te.Status)
		if err != nil {
			return nil, fmt.Errorf("scan tour earning: %w", err)
		}
		earnings.ByTour = append(earnings.ByTour, te)
	}

	return earnings, nil
}
