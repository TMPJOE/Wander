package models

// PaymentIntentResponse is returned after creating a Stripe PaymentIntent.
type PaymentIntentResponse struct {
	ClientSecret      string  `json:"client_secret"`
	PublishableKey    string  `json:"publishable_key"`
	Amount            float64 `json:"amount"`
	Currency          string  `json:"currency"`
	// AlreadyAuthorized is true when the card hold was already placed in a previous session.
	// The frontend should skip stripe.confirmCardPayment() and call the server confirm endpoint directly.
	AlreadyAuthorized bool    `json:"already_authorized"`
}

// PaymentConfirmResponse is returned after confirming a payment succeeded.
type PaymentConfirmResponse struct {
	BookingID     int    `json:"booking_id"`
	PaymentStatus string `json:"payment_status"`
}
