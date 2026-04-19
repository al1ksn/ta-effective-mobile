package handler

import (
    "log/slog"
    "net/http"
	"encoding/json"

	"ta-effective-mobile/internal/repository"
)

type SubscriptionHandler struct {
    repo *repository.SubscriptionRepository
    log  *slog.Logger
}

func NewSubscriptionHandler(repo *repository.SubscriptionRepository, log *slog.Logger) *SubscriptionHandler {
    return &SubscriptionHandler{repo: repo, log: log}
}

func (h *SubscriptionHandler) respondJSON(w http.ResponseWriter, status int, data any) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    json.NewEncoder(w).Encode(data)
}

func (h *SubscriptionHandler) respondError(w http.ResponseWriter, status int, msg string) {
    h.respondJSON(w, status, map[string]string{"error": msg})
}
