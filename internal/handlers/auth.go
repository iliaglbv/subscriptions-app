// internal/handlers/auth.go
package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"subscriptions-app/internal/middleware"
	"subscriptions-app/internal/models"
	"subscriptions-app/internal/repository"
	"subscriptions-app/internal/utils"
	"subscriptions-app/internal/validator"
)

type AuthHandler struct {
	userRepo  *repository.UserRepository
	validator *validator.Validator
	jwtSecret string
	jwtExp    time.Duration
}

func NewAuthHandler(userRepo *repository.UserRepository, v *validator.Validator, secret string, exp time.Duration) *AuthHandler {
	return &AuthHandler{
		userRepo:  userRepo,
		validator: v,
		jwtSecret: secret,
		jwtExp:    exp,
	}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req models.CreateUserRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Decode error: %v", err) // ← Добавили лог
		//log.Printf("Raw body: %s", string(body)) // ← Опционально: логировать тело
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}

	if err := h.validator.Validate(&req); err != nil {
		http.Error(w, `{"error":"validation failed","details":"`+err.Error()+`"}`, http.StatusBadRequest)
		return
	}

	hash, err := utils.HashPassword(req.Password)
	if err != nil {
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	id, err := h.userRepo.Create(req.Username, req.Email, hash)
	if err != nil {
		http.Error(w, `{"error":"username or email already exists"}`, http.StatusConflict)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":       id,
		"username": req.Username,
		"email":    req.Email,
	})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}

	user, err := h.userRepo.GetByUsername(req.Username)
	if err != nil || user == nil {
		http.Error(w, `{"error":"invalid credentials"}`, http.StatusUnauthorized)
		return
	}

	if err := utils.CheckPassword(req.Password, user.PasswordHash); err != nil {
		http.Error(w, `{"error":"invalid credentials"}`, http.StatusUnauthorized)
		return
	}

	token, err := utils.GenerateJWTToken(user.ID, user.Username, h.jwtSecret, h.jwtExp)
	if err != nil {
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"access_token": token,
		"token_type":   "Bearer",
		"expires_in":   int(h.jwtExp.Seconds()),
		"user": map[string]interface{}{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
		},
	})
}

func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)

	user, err := h.userRepo.GetByID(userID)
	if err != nil || user == nil {
		http.Error(w, `{"error":"user not found"}`, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":         user.ID,
		"username":   user.Username,
		"email":      user.Email,
		"created_at": user.CreatedAt,
	})
}
