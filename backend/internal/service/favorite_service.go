package service

import (
	"context"

	"wander/backend/internal/models"
	"wander/backend/internal/repository"
)

type FavoriteService struct {
	repo repository.FavoriteRepository
}

func NewFavoriteService(repo repository.FavoriteRepository) *FavoriteService {
	return &FavoriteService{repo: repo}
}

func (s *FavoriteService) Add(ctx context.Context, userID int, tourID int) error {
	return s.repo.Add(ctx, userID, tourID)
}

func (s *FavoriteService) Remove(ctx context.Context, userID int, tourID int) error {
	return s.repo.Remove(ctx, userID, tourID)
}

func (s *FavoriteService) List(ctx context.Context, userID int) ([]models.Tour, error) {
	return s.repo.List(ctx, userID)
}
