package service

import (
	"context"
	"fmt"

	"wander/backend/internal/models"
	"wander/backend/internal/repository"
)

// UserService handles business logic for users.
type UserService struct {
	repo repository.UserRepository
}

// NewUserService creates a new UserService.
func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

// Create creates a new user.
func (s *UserService) Create(ctx context.Context, req models.UserCreateRequest) (*models.User, error) {
	// Validate that email is not already in use.
	_, err := s.repo.GetByEmail(ctx, req.Email)
	if err == nil {
		return nil, fmt.Errorf("email already in use: %w", models.ErrConflict)
	}

	// Hash password (simplified for development; use bcrypt in production).
	hashedPassword := hashPassword(req.Password)

	return s.repo.Create(ctx, req, hashedPassword)
}

// GetByID retrieves a user by ID.
func (s *UserService) GetByID(ctx context.Context, id int) (*models.User, error) {
	return s.repo.GetByID(ctx, id)
}

// GetByEmail retrieves a user by email.
func (s *UserService) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	return s.repo.GetByEmail(ctx, email)
}

// Update updates a user.
func (s *UserService) Update(ctx context.Context, id int, req models.UserUpdateRequest) (*models.User, error) {
	return s.repo.Update(ctx, id, req)
}

// Delete deletes a user.
func (s *UserService) Delete(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}

// List lists all users.
func (s *UserService) List(ctx context.Context, limit, offset int) ([]models.User, error) {
	return s.repo.List(ctx, limit, offset)
}

// hashPassword is a placeholder for actual password hashing.
func hashPassword(password string) string {
	// TODO: Use bcrypt.GenerateFromPassword in production.
	return password
}
