package repository

import (
	"context"

	"wander/backend/internal/models"
)

type TourScheduleRepository interface {
	Create(ctx context.Context, s models.TourSchedule) (*models.TourSchedule, error)
	GetByID(ctx context.Context, id int) (*models.TourSchedule, error)
	ListByTourID(ctx context.Context, tourID int, onlyActive bool) ([]models.TourSchedule, error)
	Update(ctx context.Context, id int, req models.ScheduleUpdateRequest) (*models.TourSchedule, error)
	Delete(ctx context.Context, id int) error
	AdjustSpots(ctx context.Context, id int, delta int) error
}
