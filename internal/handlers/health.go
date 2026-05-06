// internal/handlers/health.go
package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

type HealthHandler struct {
	db *sql.DB
}

func NewHealthHandler(db *sql.DB) *HealthHandler {
	return &HealthHandler{db: db}
}

func (h *HealthHandler) Check(w http.ResponseWriter, r *http.Request) {
	type HealthResponse struct {
		Status   string `json:"status"`
		Database string `json:"database"`
	}

	resp := HealthResponse{
		Status:   "ok",
		Database: "connected",
	}

	if err := h.db.Ping(); err != nil {
		resp.Database = "disconnected"
		w.WriteHeader(http.StatusServiceUnavailable)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
