package handler

import (
	"net/http"

	"wander/backend/internal/service"
)

// Handler holds all HTTP handlers.
type Handler struct {
	AuthHandler     *AuthHandler
	UserHandler     *UserHandler
	CategoryHandler *CategoryHandler
	TourHandler     *TourHandler
	ScheduleHandler *TourScheduleHandler
	BookingHandler  *BookingHandler
	ReviewHandler   *ReviewHandler
	FavoriteHandler *FavoriteHandler
	MessageHandler  *MessageHandler
	UploadHandler   *UploadHandler
}

// NewHandler creates a new Handler with all sub-handlers.
func NewHandler(
	authService *service.AuthService,
	userService *service.UserService,
	categoryService *service.CategoryService,
	tourService *service.TourService,
	scheduleService *service.TourScheduleService,
	bookingService *service.BookingService,
	reviewService *service.ReviewService,
	favoriteService *service.FavoriteService,
	messageService *service.MessageService,
) *Handler {
	return &Handler{
		AuthHandler:     NewAuthHandler(authService),
		UserHandler:     NewUserHandler(userService),
		CategoryHandler: NewCategoryHandler(categoryService),
		TourHandler:     NewTourHandler(tourService),
		ScheduleHandler: NewTourScheduleHandler(scheduleService),
		BookingHandler:  NewBookingHandler(bookingService),
		ReviewHandler:   NewReviewHandler(reviewService),
		FavoriteHandler: NewFavoriteHandler(favoriteService),
		MessageHandler:  NewMessageHandler(messageService),
	}
}

// HealthCheck returns a simple health check response.
func (h *Handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok","service":"wander-backend"}`))
}
