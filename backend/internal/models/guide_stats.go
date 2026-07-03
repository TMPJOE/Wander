package models

// GuideStats holds aggregate stats for a guide's dashboard.
type GuideStats struct {
	TotalTours      int     `json:"total_tours"`
	PublishedTours  int     `json:"published_tours"`
	TotalBookings   int     `json:"total_bookings"`
	PendingBookings int     `json:"pending_bookings"`
	TotalRevenue    float64 `json:"total_revenue"`
	AvgRating       float64 `json:"avg_rating"`
	TotalReviews    int     `json:"total_reviews"`
}
