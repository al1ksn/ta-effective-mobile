package repository

import (
	"fmt"
	"context"
	"time"

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

func (r *SubscriptionRepository) Update(ctx context.Context, id uuid.UUID, req *model.UpdateSubscriptionRequest) (*model.Subscription, error) {
    query := `
        UPDATE subscriptions SET
            service_name = COALESCE($1, service_name),
            price        = COALESCE($2, price),
            end_date     = COALESCE($3, end_date),
            updated_at   = NOW()
        WHERE id = $4
        RETURNING id, service_name, price, user_id, start_date, end_date, created_at, updated_at`

    var endDate *time.Time
    if req.EndDate != nil {
        t, err := time.Parse("01-2006", *req.EndDate)
        if err != nil {
            return nil, fmt.Errorf("invalid end_date format: %w", err)
        }
        endDate = &t
    }

    row := r.db.QueryRow(ctx, query, req.ServiceName, req.Price, endDate, id)

    var sub model.Subscription
    if err := row.Scan(
        &sub.ID, &sub.ServiceName, &sub.Price, &sub.UserID,
        &sub.StartDate, &sub.EndDate, &sub.CreatedAt, &sub.UpdatedAt,
    ); err != nil {
        return nil, fmt.Errorf("update subscription: %w", err)
    }

    return &sub, nil
}

func (r *SubscriptionRepository) Delete(ctx context.Context, id uuid.UUID) error {
    query := `DELETE FROM subscriptions WHERE id = $1`
    
    res, err := r.db.Exec(ctx, query, id)
    if err != nil {
        return fmt.Errorf("delete subscription: %w", err)
    }
    if res.RowsAffected() == 0 {
        return fmt.Errorf("subscription not found")
    }
    return nil
}
