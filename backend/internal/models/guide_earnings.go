package models

type TourEarning struct {
	TourID    int     `json:"tour_id"`
	TourTitle string  `json:"tour_title"`
	Bookings  int     `json:"bookings"`
	Revenue   float64 `json:"revenue"`
	Status    string  `json:"status"`
}

type GuideEarnings struct {
	TotalAuthorized float64       `json:"total_authorized"`
	TotalPaid       float64       `json:"total_paid"`
	ByTour          []TourEarning `json:"by_tour"`
}
