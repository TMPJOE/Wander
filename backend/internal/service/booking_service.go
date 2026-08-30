package service

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"wander/backend/internal/models"
	"wander/backend/internal/repository"
)

type BookingService struct {
	repo           repository.BookingRepository
	scheduleRepo   repository.TourScheduleRepository
	tourRepo       repository.TourRepository
	paymentService *PaymentService
	txManager      repository.TransactionManager
}

func NewBookingService(repo repository.BookingRepository, scheduleRepo repository.TourScheduleRepository, tourRepo repository.TourRepository, paymentService *PaymentService, txManager repository.TransactionManager) *BookingService {
	return &BookingService{
		repo:           repo,
		scheduleRepo:   scheduleRepo,
		tourRepo:       tourRepo,
		paymentService: paymentService,
		txManager:      txManager,
	}
}

func (s *BookingService) Create(ctx context.Context, userID int, req models.BookingCreateRequest) (*models.Booking, error) {
	var createdBooking *models.Booking

	err := s.txManager.RunInTx(ctx, func(txCtx context.Context) error {
		schedule, err := s.scheduleRepo.GetByID(txCtx, req.ScheduleID)
		if err != nil {
			return fmt.Errorf("find schedule: %w", err)
		}

		if !schedule.IsActive || schedule.AvailableSpots < req.GuestCount {
			return fmt.Errorf("insufficient spots or schedule inactive: %w", models.ErrConflict)
		}

		tour, err := s.tourRepo.GetByID(txCtx, schedule.TourID, 0)
		if err != nil {
			return fmt.Errorf("find tour: %w", err)
		}

		totalPrice := tour.PricePerPerson * float64(req.GuestCount)

		booking := models.Booking{
			UserID:     userID,
			ScheduleID: req.ScheduleID,
			TourID:     schedule.TourID,
			GuestCount: req.GuestCount,
			TotalPrice: totalPrice,
			Status:     "pending",
			Notes:      req.Notes,
		}

		// Adjust spots inside the schedule securely in transaction.
		if err := s.scheduleRepo.AdjustSpots(txCtx, req.ScheduleID, -req.GuestCount); err != nil {
			return err
		}

		b, err := s.repo.Create(txCtx, booking)
		if err != nil {
			return err
		}
		
		createdBooking = b
		return nil
	})

	if err != nil {
		return nil, err
	}

	slog.Info("audit: booking created",
		"event", "booking_created",
		"booking_id", createdBooking.ID,
		"user_id", userID,
		"tour_id", createdBooking.TourID,
		"total_price", createdBooking.TotalPrice,
	)

	return createdBooking, nil
}

func (s *BookingService) GetByID(ctx context.Context, id int, userID int, userRole string) (*models.Booking, error) {
	b, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Verify permissions: user is traveler or guide of the tour.
	if userRole != "admin" && b.UserID != userID {
		tour, err := s.tourRepo.GetByID(ctx, b.TourID, 0)
		if err != nil || tour.GuideID != userID {
			return nil, models.ErrForbidden
		}
	}

	return b, nil
}

func (s *BookingService) ListByUser(ctx context.Context, userID int, role string) ([]models.Booking, error) {
	if role == "guide" {
		return s.repo.ListByGuideID(ctx, userID)
	}
	return s.repo.ListByUserID(ctx, userID)
}

func (s *BookingService) Cancel(ctx context.Context, id int, userID int) (*models.CancellationResult, error) {
	b, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if b.UserID != userID {
		return nil, models.ErrForbidden
	}

	if b.Status == "cancelled" || b.Status == "completed" {
		return nil, fmt.Errorf("booking already %s: %w", b.Status, models.ErrConflict)
	}

	schedule, err := s.scheduleRepo.GetByID(ctx, b.ScheduleID)
	if err != nil {
		return nil, fmt.Errorf("find schedule: %w", err)
	}

	hoursUntilTour := time.Until(schedule.StartTime).Hours()
	var refundPct int
	var refundAmount float64
	var msg string

	if hoursUntilTour > 48 {
		refundPct = 100
		refundAmount = b.TotalPrice
		msg = "Full refund issued (cancelled > 48 hours in advance)."
	} else {
		refundPct = 50
		refundAmount = b.TotalPrice * 0.5
		msg = "50% partial refund issued (cancelled within 48 hours)."
	}

	// Handle Stripe payment status
	if s.paymentService != nil {
		if b.PaymentStatus == "authorized" {
			_ = s.paymentService.CancelPaymentIntent(ctx, id)
		} else if b.PaymentStatus == "paid" && refundAmount > 0 {
			_ = s.paymentService.RefundPayment(ctx, id, refundAmount)
		}
	}

	err = s.txManager.RunInTx(ctx, func(txCtx context.Context) error {
		if err := s.repo.UpdateStatus(txCtx, id, "cancelled"); err != nil {
			return err
		}
		if err := s.scheduleRepo.AdjustSpots(txCtx, b.ScheduleID, b.GuestCount); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	slog.Info("audit: booking cancelled",
		"event", "booking_cancelled",
		"booking_id", id,
		"user_id", userID,
		"refund_amount", refundAmount,
		"message", msg,
	)

	return &models.CancellationResult{
		BookingID:    id,
		Status:       "cancelled",
		RefundAmount: refundAmount,
		RefundPct:    refundPct,
		Message:      msg,
	}, nil
}

func (s *BookingService) Confirm(ctx context.Context, id int, guideID int) error {
	b, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	tour, err := s.tourRepo.GetByID(ctx, b.TourID, 0)
	if err != nil {
		return err
	}

	if tour.GuideID != guideID {
		return models.ErrForbidden
	}

	if b.Status != "pending" {
		return fmt.Errorf("booking is not pending: %w", models.ErrConflict)
	}

	if err := s.repo.UpdateStatus(ctx, id, "confirmed"); err != nil {
		return err
	}

	slog.Info("audit: booking confirmed", "event", "booking_confirmed", "booking_id", id, "guide_id", guideID)

	if s.paymentService != nil {
		if err := s.paymentService.CapturePayment(ctx, id); err != nil {
			fmt.Printf("warning: capture payment failed for booking %d: %v\n", id, err)
		}
	}

	return nil
}

func (s *BookingService) Complete(ctx context.Context, id int, guideID int) error {
	b, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	tour, err := s.tourRepo.GetByID(ctx, b.TourID, 0)
	if err != nil {
		return err
	}

	if tour.GuideID != guideID {
		return models.ErrForbidden
	}

	if b.Status != "confirmed" {
		return fmt.Errorf("booking is not confirmed: %w", models.ErrConflict)
	}

	if err := s.repo.UpdateStatus(ctx, id, "completed"); err != nil {
		return err
	}

	slog.Info("audit: booking completed", "event", "booking_completed", "booking_id", id, "guide_id", guideID)
	return nil
}

func (s *BookingService) Reject(ctx context.Context, id int, guideID int) error {
	b, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	tour, err := s.tourRepo.GetByID(ctx, b.TourID, 0)
	if err != nil {
		return err
	}

	if tour.GuideID != guideID {
		return models.ErrForbidden
	}

	if b.Status == "cancelled" || b.Status == "completed" {
		return fmt.Errorf("booking already %s: %w", b.Status, models.ErrConflict)
	}

	err = s.txManager.RunInTx(ctx, func(txCtx context.Context) error {
		if err := s.repo.UpdateStatus(txCtx, id, "cancelled"); err != nil {
			return err
		}
		return s.scheduleRepo.AdjustSpots(txCtx, b.ScheduleID, b.GuestCount)
	})
	
	if err != nil {
		return err
	}

	slog.Info("audit: booking rejected", "event", "booking_rejected", "booking_id", id, "guide_id", guideID)

	if s.paymentService != nil {
		_ = s.paymentService.CancelPaymentIntent(ctx, id)
	}

	return nil
}

func (s *BookingService) GetEarnings(ctx context.Context, guideID int, period string) (*models.GuideEarnings, error) {
	return s.repo.GetEarnings(ctx, guideID, period)
}

