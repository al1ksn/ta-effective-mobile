package repository

import (
	"fmt"
	"context"

	"ta-effective-mobile/internal/model"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SubscriptionRepository struct {
    db *pgxpool.Pool
}

func NewSubscriptionRepository(db *pgxpool.Pool) *SubscriptionRepository {
    return &SubscriptionRepository{db: db}
}

func (r *SubscriptionRepository) Create(ctx context.Context, sub *model.Subscription) (*model.Subscription, error) {
	query := `
		INSERT INTO subscriptions (service_name, price, user_id, start_date, end_date)
        VALUES ($1, $2, $3, $4, $5)
		RETURNING id, service_name, price, user_id, start_date, end_date, created_at, updated_at
	`

	row := r.db.QueryRow(ctx, query,
		sub.ServiceName, sub.Price, sub.UserID, sub.StartDate, sub.EndDate,
	)

	var result model.Subscription
	if err := row.Scan(
        &result.ID, &result.ServiceName, &result.Price, &result.UserID,
        &result.StartDate, &result.EndDate, &result.CreatedAt, &result.UpdatedAt,
    ); err != nil {
		return nil, fmt.Errorf("create subscription: %w", err)
	}

	return &result, nil
}

func (r *SubscriptionRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Subscription, error) {
    query := `
        SELECT id, service_name, price, user_id, start_date, end_date, created_at, updated_at
        FROM subscriptions WHERE id = $1`

    row := r.db.QueryRow(ctx, query, id)

    var sub model.Subscription
    if err := row.Scan(
        &sub.ID, &sub.ServiceName, &sub.Price, &sub.UserID,
        &sub.StartDate, &sub.EndDate, &sub.CreatedAt, &sub.UpdatedAt,
    ); err != nil {
        return nil, fmt.Errorf("get subscription by id: %w", err)
    }

    return &sub, nil
}

func (r *SubscriptionRepository) List(ctx context.Context) ([]*model.Subscription, error) {
    query := `
        SELECT id, service_name, price, user_id, start_date, end_date, created_at, updated_at
        FROM subscriptions ORDER BY created_at DESC`

    rows, err := r.db.Query(ctx, query)
    if err != nil {
        return nil, fmt.Errorf("list subscriptions: %w", err)
    }
    defer rows.Close()

    var subs []*model.Subscription
    for rows.Next() {
        var sub model.Subscription
        if err := rows.Scan(
            &sub.ID, &sub.ServiceName, &sub.Price, &sub.UserID,
            &sub.StartDate, &sub.EndDate, &sub.CreatedAt, &sub.UpdatedAt,
        ); err != nil {
            return nil, fmt.Errorf("scan subscription: %w", err)
        }
        subs = append(subs, &sub)
    }

    return subs, nil
}
