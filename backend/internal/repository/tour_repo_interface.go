package repository

import (
	"context"

	"wander/backend/internal/models"
)

type TourRepository interface {
	Create(ctx context.Context, guideID int, req models.TourCreateRequest) (*models.Tour, error)
	GetByID(ctx context.Context, id int, userID int) (*models.Tour, error)
	Update(ctx context.Context, id int, req models.TourUpdateRequest) (*models.Tour, error)
	Delete(ctx context.Context, id int) error
	List(ctx context.Context, filter models.TourFilter) ([]models.Tour, error)
	GetStats(ctx context.Context, guideID int) (*models.GuideStats, error)
}
