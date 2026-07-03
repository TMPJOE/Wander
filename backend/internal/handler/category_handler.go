package handler

import (
	"net/http"

	"wander/backend/internal/service"
	"wander/backend/internal/utils"
)

type CategoryHandler struct {
	service *service.CategoryService
}

func NewCategoryHandler(service *service.CategoryService) *CategoryHandler {
	return &CategoryHandler{service: service}
}

func (h *CategoryHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /categories", h.List)
}

func (h *CategoryHandler) List(w http.ResponseWriter, r *http.Request) {
	categories, err := h.service.List(r.Context())
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, "Error al listar categorías", err.Error())
		return
	}
	utils.SendSuccess(w, http.StatusOK, "Categorías recuperadas", categories)
}
