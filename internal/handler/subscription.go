package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
    "time"

	"ta-effective-mobile/internal/model"
	"ta-effective-mobile/internal/repository"

	"github.com/google/uuid"
    "github.com/go-chi/chi/v5"
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

// @Summary      Создать подписку
// @Tags         subscriptions
// @Accept       json
// @Produce      json
// @Param        body body model.CreateSubscriptionRequest true "Данные подписки"
// @Success      201 {object} model.Subscription
// @Failure      400 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /subscriptions [post]
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

// @Summary      Получить подписку по ID
// @Tags         subscriptions
// @Produce      json
// @Param        id path string true "UUID подписки"
// @Success      200 {object} model.Subscription
// @Failure      400 {object} map[string]string
// @Failure      404 {object} map[string]string
// @Router       /subscriptions/{id} [get]
func (h *SubscriptionHandler) GetById(w http.ResponseWriter, r *http.Request) {
    id, err := uuid.Parse(chi.URLParam(r, "id"))
    if err != nil {
        h.respondError(w, http.StatusBadRequest, "invalid id")
        return
    }

    sub, err := h.repo.GetByID(r.Context(), id)
    if err != nil {
        h.log.Error("get subscription", "id", id, "error", err)
        h.respondError(w, http.StatusNotFound, "resource not found")
        return
    }

    h.respondJSON(w, http.StatusOK, sub)
}

// @Summary      Обновить подписку
// @Tags         subscriptions
// @Accept       json
// @Produce      json
// @Param        id   path string                          true "UUID подписки"
// @Param        body body model.UpdateSubscriptionRequest true "Поля для обновления"
// @Success      200 {object} model.Subscription
// @Failure      400 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /subscriptions/{id} [put]
func (h *SubscriptionHandler) Update(w http.ResponseWriter, r *http.Request) {
    id, err := uuid.Parse(chi.URLParam(r, "id"))
    if err != nil {
        h.respondError(w, http.StatusBadRequest, "invalid id")
        return
    }

    var req model.UpdateSubscriptionRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        h.respondError(w, http.StatusBadRequest, "invalid request body")
        return
    }

    updated, err := h.repo.Update(r.Context(), id, &req)
    if err != nil {
        h.log.Error("update subscription", "id", id, "error", err)
        h.respondError(w, http.StatusInternalServerError, "failed to update subscription")
        return
    }

    h.respondJSON(w, http.StatusOK, updated)
}

// @Summary      Удалить подписку
// @Tags         subscriptions
// @Param        id path string true "UUID подписки"
// @Success      204
// @Failure      400 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /subscriptions/{id} [delete]
func (h *SubscriptionHandler) Delete(w http.ResponseWriter, r *http.Request) {
    id, err := uuid.Parse(chi.URLParam(r, "id"))
    if err != nil {
        h.respondError(w, http.StatusBadRequest, "invalid id")
        return
    }

    if err := h.repo.Delete(r.Context(), id); err != nil {
        h.log.Error("delete subscription", "id", id, "error", err)
        h.respondError(w, http.StatusInternalServerError, "failed to delete subscription")
        return
    }

    w.WriteHeader(http.StatusNoContent)
}

// @Summary      Список всех подписок
// @Tags         subscriptions
// @Produce      json
// @Success      200 {array}  model.Subscription
// @Failure      500 {object} map[string]string
// @Router       /subscriptions [get]
func (h *SubscriptionHandler) List(w http.ResponseWriter, r *http.Request) {
    subs, err := h.repo.List(r.Context())
    if err != nil {
        h.log.Error("list subscriptions", "error", err)
        h.respondError(w, http.StatusInternalServerError, "failed to list subscriptions")
        return
    }

    if subs == nil {
        subs = []*model.Subscription{}
    }

    h.respondJSON(w, http.StatusOK, subs)
}

// @Summary      Суммарная стоимость подписок за период
// @Tags         subscriptions
// @Produce      json
// @Param        from         query string false "Начало периода (MM-YYYY)"
// @Param        to           query string false "Конец периода (MM-YYYY)"
// @Param        user_id      query string false "Фильтр по user_id"
// @Param        service_name query string false "Фильтр по названию сервиса"
// @Success      200 {object} model.TotalCostResponse
// @Failure      400 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /subscriptions/total [get]
func (h *SubscriptionHandler) TotalCost(w http.ResponseWriter, r *http.Request) {
    fromStr := r.URL.Query().Get("from")
    toStr := r.URL.Query().Get("to")

    from := time.Time{}
    to := time.Now()

    if fromStr != "" {
        t, err := time.Parse("01-2006", fromStr)
        if err != nil {
            h.respondError(w, http.StatusBadRequest, "invalid from format, expected MM-YYYY")
            return
        }
        from = t
    }

    if toStr != "" {
        t, err := time.Parse("01-2006", toStr)
        if err != nil {
            h.respondError(w, http.StatusBadRequest, "invalid to format, expected MM-YYYY")
            return
        }
        to = t
    }

    var userID *uuid.UUID
    if uid := r.URL.Query().Get("user_id"); uid != "" {
        parsed, err := uuid.Parse(uid)
        if err != nil {
            h.respondError(w, http.StatusBadRequest, "invalid user_id format")
            return
        }
        userID = &parsed
    }

    var serviceName *string
    if sn := r.URL.Query().Get("service_name"); sn != "" {
        serviceName = &sn
    }

    total, err := h.repo.TotalCost(r.Context(), from, to, userID, serviceName)
    if err != nil {
        h.log.Error("total cost", "error", err)
        h.respondError(w, http.StatusInternalServerError, "failed to calculate total cost")
        return
    }

    h.respondJSON(w, http.StatusOK, model.TotalCostResponse{Total: total})
}
