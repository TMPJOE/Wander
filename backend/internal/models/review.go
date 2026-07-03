package models

import "time"

// Review represents a user's review of a tour.
type Review struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	TourID    int       `json:"tour_id"`
	BookingID *int      `json:"booking_id,omitempty"`
	Rating    int       `json:"rating"`
	Comment   string    `json:"comment"`
	CreatedAt time.Time `json:"created_at"`

	// Joined fields.
	UserName   string `json:"user_name,omitempty"`
	UserAvatar string `json:"user_avatar,omitempty"`
}

// ReviewCreateRequest is used to submit a review.
type ReviewCreateRequest struct {
	Rating  int    `json:"rating"`
	Comment string `json:"comment"`
}
