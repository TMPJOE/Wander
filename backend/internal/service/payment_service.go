package service

import (
	"context"
	"fmt"
	"math"
	"strconv"

	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/paymentintent"

	"wander/backend/internal/models"
	"wander/backend/internal/repository"
)

type PaymentService struct {
	bookingRepo    repository.BookingRepository
	publishableKey string
}

func NewPaymentService(bookingRepo repository.BookingRepository, secretKey string, publishableKey string) *PaymentService {
	stripe.Key = secretKey
	return &PaymentService{
		bookingRepo:    bookingRepo,
		publishableKey: publishableKey,
	}
}

// CreateIntent creates a Stripe PaymentIntent for a booking and returns the
// client secret the frontend needs to confirm the card payment.
func (s *PaymentService) CreateIntent(ctx context.Context, bookingID int, userID int) (*models.PaymentIntentResponse, error) {
	b, err := s.bookingRepo.GetByID(ctx, bookingID)
	if err != nil {
		return nil, err
	}
	if b.UserID != userID {
		return nil, models.ErrForbidden
	}
	if b.PaymentStatus == "paid" {
		return nil, fmt.Errorf("booking already paid: %w", models.ErrConflict)
	}

	// Demo currency: Stripe test mode, amounts treated as USD cents.
	amount := int64(math.Round(b.TotalPrice * 100))

	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(amount),
		Currency: stripe.String(string(stripe.CurrencyUSD)),
	}
	params.AddMetadata("booking_id", strconv.Itoa(bookingID))

	pi, err := paymentintent.New(params)
	if err != nil {
		return nil, fmt.Errorf("create stripe payment intent: %w", err)
	}

	if err := s.bookingRepo.UpdatePayment(ctx, bookingID, pi.ID, "pending"); err != nil {
		return nil, err
	}

	return &models.PaymentIntentResponse{
		ClientSecret:   pi.ClientSecret,
		PublishableKey: s.publishableKey,
		Amount:         b.TotalPrice,
		Currency:       "usd",
	}, nil
}

// ConfirmPayment verifies with Stripe that the payment succeeded and updates
// the booking's payment status accordingly.
func (s *PaymentService) ConfirmPayment(ctx context.Context, bookingID int, userID int) (*models.PaymentConfirmResponse, error) {
	b, err := s.bookingRepo.GetByID(ctx, bookingID)
	if err != nil {
		return nil, err
	}
	if b.UserID != userID {
		return nil, models.ErrForbidden
	}
	if b.StripePaymentIntentID == "" {
		return nil, fmt.Errorf("no payment intent found for booking: %w", models.ErrBadRequest)
	}

	pi, err := paymentintent.Get(b.StripePaymentIntentID, nil)
	if err != nil {
		return nil, fmt.Errorf("retrieve stripe payment intent: %w", err)
	}

	status := "pending"
	if pi.Status == stripe.PaymentIntentStatusSucceeded {
		status = "paid"
	}

	if err := s.bookingRepo.UpdatePayment(ctx, bookingID, b.StripePaymentIntentID, status); err != nil {
		return nil, err
	}

	if status != "paid" {
		return nil, fmt.Errorf("payment not completed, status: %s: %w", pi.Status, models.ErrConflict)
	}

	return &models.PaymentConfirmResponse{
		BookingID:     bookingID,
		PaymentStatus: status,
	}, nil
}
