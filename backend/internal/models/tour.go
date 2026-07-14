package models

import (
	"time"
)

// Tour represents a guided tour.
type Tour struct {
	ID              int             `json:"id"`
	GuideID         int             `json:"guide_id"`
	CategoryID      int             `json:"category_id"`
	Title           string          `json:"title"`
	Description     string          `json:"description"`
	Location        string          `json:"location"`
	Latitude        *float64        `json:"latitude,omitempty"`
	Longitude       *float64        `json:"longitude,omitempty"`
	DurationMinutes int             `json:"duration_minutes"`
	PricePerPerson  float64         `json:"price_per_person"`
	MaxGuests       int             `json:"max_guests"`
	Difficulty      string          `json:"difficulty"`
	Languages       []string        `json:"languages"`
	WhatIncluded    []string        `json:"what_included"`
	MeetingPoint    string          `json:"meeting_point"`
	Images          []string        `json:"images"`
	IsPublished     bool            `json:"is_published"`
	AvgRating       float64         `json:"avg_rating"`
	ReviewCount     int             `json:"review_count"`
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at"`

	// Joined fields (not always populated).
	GuideName     string `json:"guide_name,omitempty"`
	GuideAvatar   string `json:"guide_avatar,omitempty"`
	CategoryName  string `json:"category_name,omitempty"`
	CategorySlug  string `json:"category_slug,omitempty"`
	IsFavorited   bool   `json:"is_favorited,omitempty"`
}

// TourCreateRequest is used by guides to create a tour.
type TourCreateRequest struct {
	CategoryID      int      `json:"category_id"`
	Title           string   `json:"title"`
	Description     string   `json:"description"`
	Location        string   `json:"location"`
	Latitude        *float64 `json:"latitude,omitempty"`
	Longitude       *float64 `json:"longitude,omitempty"`
	DurationMinutes int      `json:"duration_minutes"`
	PricePerPerson  float64  `json:"price_per_person"`
	MaxGuests       int      `json:"max_guests"`
	Difficulty      string   `json:"difficulty"`
	Languages       []string `json:"languages"`
	WhatIncluded    []string `json:"what_included"`
	MeetingPoint    string   `json:"meeting_point"`
	Images          []string `json:"images"`
	IsPublished     bool     `json:"is_published"`
}

// TourUpdateRequest is used by guides to update a tour.
type TourUpdateRequest struct {
	CategoryID      *int      `json:"category_id,omitempty"`
	Title           string    `json:"title,omitempty"`
	Description     string    `json:"description,omitempty"`
	Location        string    `json:"location,omitempty"`
	Latitude        *float64  `json:"latitude,omitempty"`
	Longitude       *float64  `json:"longitude,omitempty"`
	DurationMinutes *int      `json:"duration_minutes,omitempty"`
	PricePerPerson  *float64  `json:"price_per_person,omitempty"`
	MaxGuests       *int      `json:"max_guests,omitempty"`
	Difficulty      string    `json:"difficulty,omitempty"`
	Languages       []string  `json:"languages,omitempty"`
	WhatIncluded    []string  `json:"what_included,omitempty"`
	MeetingPoint    string    `json:"meeting_point,omitempty"`
	Images          []string  `json:"images,omitempty"`
	IsPublished     *bool     `json:"is_published,omitempty"`
}

// TourFilter defines filters for listing tours.
type TourFilter struct {
	CategoryID   int      `json:"category_id,omitempty"`
	CategorySlug string   `json:"category_slug,omitempty"`
	Difficulty   string   `json:"difficulty,omitempty"`
	MinPrice     *float64 `json:"min_price,omitempty"`
	MaxPrice     *float64 `json:"max_price,omitempty"`
	Query        string   `json:"q,omitempty"`
	Location     string   `json:"location,omitempty"`
	GuideID      int      `json:"guide_id,omitempty"`
	Limit        int      `json:"limit,omitempty"`
	Offset       int      `json:"offset,omitempty"`
	UserID       int      `json:"-"` // for favorite check
}
