package models

import "time"

// Favorite represents a user's favorited tour.
type Favorite struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	TourID    int       `json:"tour_id"`
	CreatedAt time.Time `json:"created_at"`
}
