package service

import (
	"context"

	"wander/backend/internal/models"
	"wander/backend/internal/repository"
)

type TourScheduleService struct {
	repo     repository.TourScheduleRepository
	tourRepo repository.TourRepository
}

func NewTourScheduleService(repo repository.TourScheduleRepository, tourRepo repository.TourRepository) *TourScheduleService {
	return &TourScheduleService{repo: repo, tourRepo: tourRepo}
}

func (s *TourScheduleService) Create(ctx context.Context, guideID int, schedule models.TourSchedule) (*models.TourSchedule, error) {
	tour, err := s.tourRepo.GetByID(ctx, schedule.TourID, 0)
	if err != nil {
		return nil, err
	}
	if tour.GuideID != guideID {
		return nil, models.ErrForbidden
	}
	return s.repo.Create(ctx, schedule)
}

func (s *TourScheduleService) ListByTourID(ctx context.Context, tourID int, onlyActive bool) ([]models.TourSchedule, error) {
	return s.repo.ListByTourID(ctx, tourID, onlyActive)
}

func (s *TourScheduleService) Update(ctx context.Context, guideID int, id int, req models.ScheduleUpdateRequest) (*models.TourSchedule, error) {
	sched, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	tour, err := s.tourRepo.GetByID(ctx, sched.TourID, 0)
	if err != nil {
		return nil, err
	}
	if tour.GuideID != guideID {
		return nil, models.ErrForbidden
	}
	return s.repo.Update(ctx, id, req)
}

func (s *TourScheduleService) Delete(ctx context.Context, guideID int, id int) error {
	sched, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	tour, err := s.tourRepo.GetByID(ctx, sched.TourID, 0)
	if err != nil {
		return err
	}
	if tour.GuideID != guideID {
		return models.ErrForbidden
	}
	return s.repo.Delete(ctx, id)
}
