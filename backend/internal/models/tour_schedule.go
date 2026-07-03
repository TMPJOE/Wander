package models

import "time"

// TourSchedule represents an available time slot for a tour.
type TourSchedule struct {
	ID             int       `json:"id"`
	TourID         int       `json:"tour_id"`
	StartTime      time.Time `json:"start_time"`
	EndTime        time.Time `json:"end_time"`
	AvailableSpots int       `json:"available_spots"`
	IsActive       bool      `json:"is_active"`
	CreatedAt      time.Time `json:"created_at"`
}

// ScheduleCreateRequest is used to add a new schedule to a tour.
type ScheduleCreateRequest struct {
	StartTime      time.Time `json:"start_time"`
	EndTime        time.Time `json:"end_time"`
	AvailableSpots int       `json:"available_spots"`
}

// ScheduleUpdateRequest is used to update a schedule.
type ScheduleUpdateRequest struct {
	StartTime      *time.Time `json:"start_time,omitempty"`
	EndTime        *time.Time `json:"end_time,omitempty"`
	AvailableSpots *int       `json:"available_spots,omitempty"`
	IsActive       *bool      `json:"is_active,omitempty"`
}
