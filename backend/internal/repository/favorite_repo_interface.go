package repository

import (
	"context"

	"wander/backend/internal/models"
)

type FavoriteRepository interface {
	Add(ctx context.Context, userID int, tourID int) error
	Remove(ctx context.Context, userID int, tourID int) error
	List(ctx context.Context, userID int) ([]models.Tour, error)
}
