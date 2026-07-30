package service

import (
	"context"
	"fmt"
	"log/slog"

	"wander/backend/internal/models"
	"wander/backend/internal/repository"
	"wander/backend/internal/storage"
)

type TourService struct {
	repo    repository.TourRepository
	storage storage.Provider
}

// NewTourService builds a TourService. storage may be nil (deletion of
// orphaned uploaded files is skipped in that case), which keeps tests and
// other callers that don't care about file cleanup working unchanged.
func NewTourService(repo repository.TourRepository, storage storage.Provider) *TourService {
	return &TourService{repo: repo, storage: storage}
}

// deleteImages deletes the given image URLs from the storage provider
// best-effort. It never returns an error: a failure to remove a file must not
// block a tour deletion/update, and orphaned files can be reclaimed later by
// a janitor. Failures are logged.
func (s *TourService) deleteImages(urls []string) {
	if s.storage == nil || len(urls) == 0 {
		return
	}
	if err := s.storage.DeleteMany(context.Background(), urls); err != nil {
		slog.Warn("failed to delete tour images from storage", "count", len(urls), "error", err)
	}
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

	// If the caller is replacing the image set, clean up any URLs that are no
	// longer referenced. This is the case where a guide uploaded an image by
	// mistake and then removed it from the tour in the editor: the file would
	// otherwise persist in storage forever.
	if req.Images != nil {
		newSet := make(map[string]struct{}, len(req.Images))
		for _, u := range req.Images {
			newSet[u] = struct{}{}
		}
		var orphans []string
		for _, u := range t.Images {
			if _, ok := newSet[u]; !ok {
				orphans = append(orphans, u)
			}
		}
		s.deleteImages(orphans)
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

	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}

	// Best-effort cleanup of the tour's uploaded images now that the DB rows
	// are gone (tour_images cascade on tour delete, but the files in storage
	// do not).
	s.deleteImages(t.Images)
	return nil
}

func (s *TourService) List(ctx context.Context, filter models.TourFilter) ([]models.Tour, error) {
	return s.repo.List(ctx, filter)
}

func (s *TourService) GetStats(ctx context.Context, guideID int, period string) (*models.GuideStats, error) {
	return s.repo.GetStats(ctx, guideID, period)
}
