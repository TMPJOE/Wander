package service

import (
	"context"
	"fmt"

	"wander/backend/internal/models"
	"wander/backend/internal/repository"
)

type TourService struct {
	repo repository.TourRepository
}

func NewTourService(repo repository.TourRepository) *TourService {
	return &TourService{repo: repo}
}

func (s *TourService) Create(ctx context.Context, guideID int, req models.TourCreateRequest) (*models.Tour, error) {
	if req.Title == "" {
		return nil, fmt.Errorf("title is required: %w", models.ErrBadRequest)
	}
	if req.PricePerPerson < 0 {
		return nil, fmt.Errorf("price cannot be negative: %w", models.ErrBadRequest)
	}
	return s.repo.Create(ctx, guideID, req)
}

func (s *TourService) GetByID(ctx context.Context, id int, userID int) (*models.Tour, error) {
	return s.repo.GetByID(ctx, id, userID)
}

func (s *TourService) Update(ctx context.Context, id int, guideID int, req models.TourUpdateRequest) (*models.Tour, error) {
	t, err := s.repo.GetByID(ctx, id, 0)
	if err != nil {
		return nil, err
	}
	if t.GuideID != guideID {
		return nil, models.ErrForbidden
	}

	return s.repo.Update(ctx, id, req)
}

func (s *TourService) Delete(ctx context.Context, id int, guideID int) error {
	t, err := s.repo.GetByID(ctx, id, 0)
	if err != nil {
		return err
	}
	if t.GuideID != guideID {
		return models.ErrForbidden
	}
	return s.repo.Delete(ctx, id)
}

func (s *TourService) List(ctx context.Context, filter models.TourFilter) ([]models.Tour, error) {
	return s.repo.List(ctx, filter)
}

func (s *TourService) GetStats(ctx context.Context, guideID int) (*models.GuideStats, error) {
	return s.repo.GetStats(ctx, guideID)
}
