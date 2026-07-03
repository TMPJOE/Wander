package repository

import (
	"context"
	"sync"
	"time"

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

// InMemoryUserRepository is a simple in-memory implementation for development.
type InMemoryUserRepository struct {
	mu      sync.RWMutex
	users   map[int]*models.User
	nextID  int
}

// NewInMemoryUserRepository creates a new in-memory user repository.
func NewInMemoryUserRepository() UserRepository {
	return &InMemoryUserRepository{
		users:  make(map[int]*models.User),
		nextID: 1,
	}
}

func (r *InMemoryUserRepository) Create(ctx context.Context, req models.UserCreateRequest, hashedPassword string) (*models.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	user := &models.User{
		ID:        r.nextID,
		Email:     req.Email,
		Username:  req.Username,
		Password:  hashedPassword,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	r.users[user.ID] = user
	r.nextID++
	return user, nil
}

func (r *InMemoryUserRepository) GetByID(ctx context.Context, id int) (*models.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, ok := r.users[id]
	if !ok {
		return nil, models.ErrNotFound
	}
	return user, nil
}

func (r *InMemoryUserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, user := range r.users {
		if user.Email == email {
			return user, nil
		}
	}
	return nil, models.ErrNotFound
}

func (r *InMemoryUserRepository) Update(ctx context.Context, id int, req models.UserUpdateRequest) (*models.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	user, ok := r.users[id]
	if !ok {
		return nil, models.ErrNotFound
	}

	if req.Email != "" {
		user.Email = req.Email
	}
	if req.Username != "" {
		user.Username = req.Username
	}
	if req.FirstName != "" {
		user.FirstName = req.FirstName
	}
	if req.LastName != "" {
		user.LastName = req.LastName
	}
	user.UpdatedAt = time.Now()
	return user, nil
}

func (r *InMemoryUserRepository) Delete(ctx context.Context, id int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.users[id]; !ok {
		return models.ErrNotFound
	}
	delete(r.users, id)
	return nil
}

func (r *InMemoryUserRepository) List(ctx context.Context, limit, offset int) ([]models.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	users := make([]models.User, 0, len(r.users))
	for _, user := range r.users {
		users = append(users, *user)
	}
	return users, nil
}
