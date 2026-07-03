package service

import (
	"context"
	"fmt"

	"wander/backend/internal/models"
	"wander/backend/internal/repository"
)

type ReviewService struct {
	repo        repository.ReviewRepository
	bookingRepo repository.BookingRepository
}

func NewReviewService(repo repository.ReviewRepository, bookingRepo repository.BookingRepository) *ReviewService {
	return &ReviewService{repo: repo, bookingRepo: bookingRepo}
}

func (s *ReviewService) Create(ctx context.Context, userID int, tourID int, req models.ReviewCreateRequest) (*models.Review, error) {
	if req.Rating < 1 || req.Rating > 5 {
		return nil, fmt.Errorf("rating must be 1-5: %w", models.ErrBadRequest)
	}

	// Verify the user has a confirmed/completed booking for this tour.
	bookings, err := s.bookingRepo.ListByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	var validBooking *models.Booking
	for _, b := range bookings {
		if b.TourID == tourID && (b.Status == "confirmed" || b.Status == "completed") {
			validBooking = &b
			break
		}
	}

	if validBooking == nil {
		return nil, fmt.Errorf("must have a confirmed booking to review: %w", models.ErrForbidden)
	}

	return s.repo.Create(ctx, userID, tourID, &validBooking.ID, req)
}

func (s *ReviewService) ListByTourID(ctx context.Context, tourID int) ([]models.Review, error) {
	return s.repo.ListByTourID(ctx, tourID)
}
