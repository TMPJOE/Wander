package models

import "time"

// Booking represents a user's booking for a tour schedule.
type Booking struct {
	ID         int       `json:"id"`
	UserID     int       `json:"user_id"`
	ScheduleID int       `json:"schedule_id"`
	TourID     int       `json:"tour_id"`
	GuestCount int       `json:"guest_count"`
	TotalPrice float64   `json:"total_price"`
	Status     string    `json:"status"`
	Notes      string    `json:"notes"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`

	// Payment tracking (Stripe).
	PaymentStatus         string `json:"payment_status,omitempty"`
	StripePaymentIntentID string `json:"stripe_payment_intent_id,omitempty"`

	// Joined fields.
	TourTitle     string    `json:"tour_title,omitempty"`
	TourImage     string    `json:"tour_image,omitempty"`
	TourLocation  string    `json:"tour_location,omitempty"`
	GuideName     string    `json:"guide_name,omitempty"`
	GuideAvatar   string    `json:"guide_avatar,omitempty"`
	UserName      string    `json:"user_name,omitempty"`
	ScheduleStart time.Time `json:"schedule_start,omitempty"`
	ScheduleEnd   time.Time `json:"schedule_end,omitempty"`
}

// BookingCreateRequest is used to create a booking.
type BookingCreateRequest struct {
	ScheduleID int    `json:"schedule_id"`
	GuestCount int    `json:"guest_count"`
	Notes      string `json:"notes,omitempty"`
}
