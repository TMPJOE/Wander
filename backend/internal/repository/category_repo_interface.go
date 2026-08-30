package repository

import (
	"context"

	"wander/backend/internal/models"
)

type CategoryRepository interface {
	List(ctx context.Context) ([]models.Category, error)
	GetByID(ctx context.Context, id int) (*models.Category, error)
}
