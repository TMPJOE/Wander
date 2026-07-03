package repository

import (
	"context"

	"wander/backend/internal/models"
)

// UserRepository defines the interface for user data access.
type UserRepository interface {
	Create(ctx context.Context, req models.UserCreateRequest, hashedPassword string) (*models.User, error)
	GetByID(ctx context.Context, id int) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	Update(ctx context.Context, id int, req models.UserUpdateRequest) (*models.User, error)
	Delete(ctx context.Context, id int) error
	List(ctx context.Context, limit, offset int) ([]models.User, error)
}
