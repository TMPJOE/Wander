package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"wander/backend/internal/models"
)

type ReviewRepository interface {
	Create(ctx context.Context, userID int, tourID int, bookingID *int, req models.ReviewCreateRequest) (*models.Review, error)
	ListByTourID(ctx context.Context, tourID int) ([]models.Review, error)
}

type PgReviewRepository struct {
	pool *pgxpool.Pool
}

func NewPgReviewRepository(pool *pgxpool.Pool) ReviewRepository {
	return &PgReviewRepository{pool: pool}
}

func (r *PgReviewRepository) Create(ctx context.Context, userID int, tourID int, bookingID *int, req models.ReviewCreateRequest) (*models.Review, error) {
	query := `
		INSERT INTO reviews (user_id, tour_id, booking_id, rating, comment)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, user_id, tour_id, booking_id, rating, comment, created_at
	`
	res := &models.Review{}
	err := r.pool.QueryRow(ctx, query, userID, tourID, bookingID, req.Rating, req.Comment).Scan(
		&res.ID, &res.UserID, &res.TourID, &res.BookingID, &res.Rating, &res.Comment, &res.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("create review: %w", err)
	}
	return res, nil
}

func (r *PgReviewRepository) ListByTourID(ctx context.Context, tourID int) ([]models.Review, error) {
	query := `
		SELECT r.id, r.user_id, r.tour_id, r.booking_id, r.rating, r.comment, r.created_at,
		       u.first_name || ' ' || u.last_name as user_name, u.avatar_url as user_avatar
		FROM reviews r
		JOIN users u ON r.user_id = u.id
		WHERE r.tour_id = $1
		ORDER BY r.created_at DESC
	`
	rows, err := r.pool.Query(ctx, query, tourID)
	if err != nil {
		return nil, fmt.Errorf("list reviews: %w", err)
	}
	defer rows.Close()

	var reviews []models.Review
	for rows.Next() {
		var rev models.Review
		err := rows.Scan(
			&rev.ID, &rev.UserID, &rev.TourID, &rev.BookingID, &rev.Rating, &rev.Comment, &rev.CreatedAt,
			&rev.UserName, &rev.UserAvatar,
		)
		if err != nil {
			return nil, err
		}
		reviews = append(reviews, rev)
	}
	return reviews, nil
}
