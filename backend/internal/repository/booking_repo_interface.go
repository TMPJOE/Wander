package repository

import (
	"context"

	"wander/backend/internal/models"
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
