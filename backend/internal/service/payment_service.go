package service

import (
	"context"
	"fmt"
	"math"
	"strconv"

	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/paymentintent"
	"github.com/stripe/stripe-go/v76/refund"

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

// CreateIntent returns a Stripe PaymentIntent client_secret for the checkout page.
// If an existing PaymentIntent is still usable (requires_payment_method, requires_action),
// it is reused. If already authorized (requires_capture), AlreadyAuthorized is set so the
// frontend can skip confirmCardPayment and go straight to the server confirm step.
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

	// If there is already a PaymentIntent for this booking, try to reuse it.
	if b.StripePaymentIntentID != "" {
		pi, err := paymentintent.Get(b.StripePaymentIntentID, nil)
		if err == nil {
			switch pi.Status {
			case stripe.PaymentIntentStatusRequiresPaymentMethod,
				stripe.PaymentIntentStatusRequiresAction,
				stripe.PaymentIntentStatusRequiresConfirmation:
				// Still needs the card — return existing client_secret so the frontend can confirm.
				return &models.PaymentIntentResponse{
					ClientSecret:   pi.ClientSecret,
					PublishableKey: s.publishableKey,
					Amount:         b.TotalPrice,
					Currency:       "usd",
				}, nil
			case stripe.PaymentIntentStatusRequiresCapture:
				// Card already authorized — tell the frontend to skip confirmCardPayment.
				if dbErr := s.bookingRepo.UpdatePayment(ctx, bookingID, pi.ID, "authorized"); dbErr != nil {
					fmt.Printf("warning: sync authorized status for booking %d: %v\n", bookingID, dbErr)
				}
				return &models.PaymentIntentResponse{
					ClientSecret:      pi.ClientSecret,
					PublishableKey:    s.publishableKey,
					Amount:            b.TotalPrice,
					Currency:          "usd",
					AlreadyAuthorized: true,
				}, nil
			}
			// For any other terminal state (succeeded, canceled), fall through and create a new PI.
		}
	}

	amount := int64(math.Round(b.TotalPrice * 100))

	params := &stripe.PaymentIntentParams{
		CaptureMethod: stripe.String("manual"),
		Amount:        stripe.Int64(amount),
		Currency:      stripe.String(string(stripe.CurrencyUSD)),
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

// AuthorizePayment verifies with Stripe that card authorization succeeded.
func (s *PaymentService) AuthorizePayment(ctx context.Context, bookingID int, userID int) (*models.PaymentConfirmResponse, error) {
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

	var status string
	switch pi.Status {
	case stripe.PaymentIntentStatusRequiresCapture:
		status = "authorized"
	case stripe.PaymentIntentStatusSucceeded:
		status = "paid"
	default:
		return nil, fmt.Errorf("payment not authorized, status: %s: %w", pi.Status, models.ErrConflict)
	}

	if err := s.bookingRepo.UpdatePayment(ctx, bookingID, b.StripePaymentIntentID, status); err != nil {
		return nil, err
	}

	return &models.PaymentConfirmResponse{
		BookingID:     bookingID,
		PaymentStatus: status,
	}, nil
}

// CapturePayment captures an authorized PaymentIntent when a guide confirms the booking.
func (s *PaymentService) CapturePayment(ctx context.Context, bookingID int) error {
	b, err := s.bookingRepo.GetByID(ctx, bookingID)
	if err != nil {
		return err
	}
	if b.StripePaymentIntentID == "" {
		// If no payment intent (e.g. cash or test mode without stripe intent), mark as paid directly
		return s.bookingRepo.UpdatePayment(ctx, bookingID, "", "paid")
	}

	_, err = paymentintent.Capture(b.StripePaymentIntentID, nil)
	if err != nil {
		return fmt.Errorf("capture stripe payment intent: %w", err)
	}

	return s.bookingRepo.UpdatePayment(ctx, bookingID, b.StripePaymentIntentID, "paid")
}

// CancelPaymentIntent cancels an uncaptured PaymentIntent (releases authorization hold).
func (s *PaymentService) CancelPaymentIntent(ctx context.Context, bookingID int) error {
	b, err := s.bookingRepo.GetByID(ctx, bookingID)
	if err != nil {
		return err
	}
	if b.StripePaymentIntentID == "" {
		return nil
	}

	_, err = paymentintent.Cancel(b.StripePaymentIntentID, nil)
	if err != nil {
		// If already canceled or unable to cancel, log and proceed
		fmt.Printf("warning: cancel payment intent %s failed: %v\n", b.StripePaymentIntentID, err)
	}

	return s.bookingRepo.UpdatePayment(ctx, bookingID, b.StripePaymentIntentID, "failed")
}

// RefundPayment issues a full or partial refund for a paid booking via Stripe.
func (s *PaymentService) RefundPayment(ctx context.Context, bookingID int, amount float64) error {
	b, err := s.bookingRepo.GetByID(ctx, bookingID)
	if err != nil {
		return err
	}
	if b.StripePaymentIntentID == "" {
		return nil
	}

	cents := int64(math.Round(amount * 100))
	params := &stripe.RefundParams{
		PaymentIntent: stripe.String(b.StripePaymentIntentID),
		Amount:        stripe.Int64(cents),
	}

	_, err = refund.New(params)
	if err != nil {
		return fmt.Errorf("refund stripe payment: %w", err)
	}

	return nil
}

