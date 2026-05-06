// internal/handlers/subscriptions.go
package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"subscriptions-app/internal/middleware"
	"subscriptions-app/internal/models"
	"subscriptions-app/internal/repository"
	"subscriptions-app/internal/validator"
)

type SubscriptionHandler struct {
	subRepo   *repository.SubscriptionRepository
	validator *validator.Validator
}

func NewSubscriptionHandler(subRepo *repository.SubscriptionRepository, v *validator.Validator) *SubscriptionHandler {
	return &SubscriptionHandler{
		subRepo:   subRepo,
		validator: v,
	}
}

func (h *SubscriptionHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)

	var req models.CreateSubscriptionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}

	if err := h.validator.Validate(&req); err != nil {
		http.Error(w, `{"error":"validation failed","details":"`+err.Error()+`"}`, http.StatusBadRequest)
		return
	}

	id, err := h.subRepo.Create(&req, userID)
	if err != nil {
		http.Error(w, `{"error":"failed to create subscription"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{"id": id})
}

func (h *SubscriptionHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)

	subs, err := h.subRepo.GetAllByUserID(userID)
	if err != nil {
		http.Error(w, `{"error":"failed to fetch subscriptions"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(subs)
}

func (h *SubscriptionHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	idStr := r.PathValue("id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, `{"error":"invalid subscription id"}`, http.StatusBadRequest)
		return
	}

	sub, err := h.subRepo.GetByID(id, userID)
	if err != nil {
		http.Error(w, `{"error":"subscription not found"}`, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sub)
}

func (h *SubscriptionHandler) Update(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	idStr := r.PathValue("id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, `{"error":"invalid subscription id"}`, http.StatusBadRequest)
		return
	}

	var req models.CreateSubscriptionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}

	if err := h.validator.Validate(&req); err != nil {
		http.Error(w, `{"error":"validation failed"}`, http.StatusBadRequest)
		return
	}

	if err := h.subRepo.Update(id, userID, &req); err != nil {
		http.Error(w, `{"error":"failed to update subscription"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "updated"})
}

func (h *SubscriptionHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	idStr := r.PathValue("id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, `{"error":"invalid subscription id"}`, http.StatusBadRequest)
		return
	}

	if err := h.subRepo.Delete(id, userID); err != nil {
		http.Error(w, `{"error":"failed to delete subscription"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *SubscriptionHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)

	total, err := h.subRepo.GetTotalMonthlyCost(userID)
	if err != nil {
		http.Error(w, `{"error":"failed to fetch stats"}`, http.StatusInternalServerError)
		return
	}

	byCategory, err := h.subRepo.GetStatsByCategory(userID)
	if err != nil {
		http.Error(w, `{"error":"failed to fetch stats"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"total_monthly_cost": total,
		"by_category":        byCategory,
	})
}
