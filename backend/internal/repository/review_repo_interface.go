package repository

import (
	"context"

	"wander/backend/internal/models"
)

type ReviewRepository interface {
	Create(ctx context.Context, userID int, tourID int, bookingID *int, req models.ReviewCreateRequest) (*models.Review, error)
	Update(ctx context.Context, reviewID int, userID int, req models.ReviewCreateRequest) (*models.Review, error)
	ListByTourID(ctx context.Context, tourID int) ([]models.Review, error)
	ListByUserID(ctx context.Context, userID int) ([]models.Review, error)
}
