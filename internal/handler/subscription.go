package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
    "time"

	"ta-effective-mobile/internal/model"
	"ta-effective-mobile/internal/repository"

	"github.com/google/uuid"
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

func (h *SubscriptionHandler) Create(w http.ResponseWriter, r *http.Request) {
    var req model.CreateSubscriptionRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        h.respondError(w, http.StatusBadRequest, "invalid request body")
        return
    }

    if req.ServiceName == "" || req.Price <= 0 || req.UserID == "" || req.StartDate == "" {
        h.respondError(w, http.StatusBadRequest, "service_name, price, user_id and start_date are required")
        return
    }

    userID, err := uuid.Parse(req.UserID)
    if err != nil {
        h.respondError(w, http.StatusBadRequest, "invalid user_id format")
        return
    }

    startDate, err := time.Parse("01-2006", req.StartDate)
    if err != nil {
        h.respondError(w, http.StatusBadRequest, "invalid start_date format, expected MM-YYYY")
        return
    }

    sub := &model.Subscription{
        ServiceName: req.ServiceName,
        Price:       req.Price,
        UserID:      userID,
        StartDate:   startDate,
    }

    if req.EndDate != nil {
        endDate, err := time.Parse("01-2006", *req.EndDate)
        if err != nil {
            h.respondError(w, http.StatusBadRequest, "invalid end_date format, expected MM-YYYY")
            return
        }
        sub.EndDate = &endDate
    }

    created, err := h.repo.Create(r.Context(), sub)
    if err != nil {
        h.log.Error("create subscription", "error", err)
        h.respondError(w, http.StatusInternalServerError, "failed to create subscription")
        return
    }

    h.respondJSON(w, http.StatusCreated, created)
}
