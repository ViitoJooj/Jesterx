package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/ViitoJooj/Jesterx/internal/domain"
)

func NewPaymentRepository(db *sql.DB) *connection {
	return &connection{db: db}
}

func (r *connection) ListActivePlans() ([]domain.Plan, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		SELECT id, name, COALESCE(description, ''), COALESCE(description_md, ''), price,
		       COALESCE(billing_cycle, 'monthly'), active,
		       COALESCE(max_sites, 1), COALESCE(max_routes, 5),
		       created_at, updated_at
		FROM plans
		WHERE active = true
		ORDER BY price ASC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	plans := make([]domain.Plan, 0)
	for rows.Next() {
		var plan domain.Plan
		if err := rows.Scan(
			&plan.ID, &plan.Name, &plan.Description, &plan.DescriptionM,
			&plan.Price, &plan.BillingCycle, &plan.Active,
			&plan.MaxSites, &plan.MaxRoutes,
			&plan.CreatedAt, &plan.UpdatedAt,
		); err != nil {
			return nil, err
		}
		plans = append(plans, plan)
	}
	return plans, rows.Err()
}

func (r *connection) FindPlanByID(id int64) (*domain.Plan, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		SELECT id, name, COALESCE(description, ''), COALESCE(description_md, ''), price,
		       COALESCE(billing_cycle, 'monthly'), active,
		       COALESCE(max_sites, 1), COALESCE(max_routes, 5),
		       created_at, updated_at
		FROM plans WHERE id = $1
	`

	var plan domain.Plan
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&plan.ID, &plan.Name, &plan.Description, &plan.DescriptionM,
		&plan.Price, &plan.BillingCycle, &plan.Active,
		&plan.MaxSites, &plan.MaxRoutes,
		&plan.CreatedAt, &plan.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &plan, nil
}

func (r *connection) FindPlanByName(name string) (*domain.Plan, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		SELECT id, name, COALESCE(description, ''), COALESCE(description_md, ''), price,
		       COALESCE(billing_cycle, 'monthly'), active,
		       COALESCE(max_sites, 1), COALESCE(max_routes, 5),
		       created_at, updated_at
		FROM plans WHERE LOWER(name) = LOWER($1) AND active = true LIMIT 1
	`

	var plan domain.Plan
	err := r.db.QueryRowContext(ctx, query, name).Scan(
		&plan.ID, &plan.Name, &plan.Description, &plan.DescriptionM,
		&plan.Price, &plan.BillingCycle, &plan.Active,
		&plan.MaxSites, &plan.MaxRoutes,
		&plan.CreatedAt, &plan.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &plan, nil
}

func (r *connection) FindLatestCompletedPaymentByUserID(userID string) (*domain.Payment, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var p domain.Payment
	err := r.db.QueryRowContext(ctx, `
		SELECT id, user_id, website_id, reference_id, type, quantity, amount, currency, status, purchased_in
		FROM payments
		WHERE user_id = $1 AND status = 'completed'
		ORDER BY purchased_in DESC LIMIT 1
	`, userID).Scan(
		&p.ID, &p.UserID, &p.WebsiteID, &p.ReferenceID, &p.Type,
		&p.Quantity, &p.Amount, &p.Currency, &p.Status, &p.PurchasedIn,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *connection) CreatePayment(payment domain.Payment) (*domain.Payment, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		INSERT INTO payments (user_id, website_id, reference_id, type, quantity, amount, currency, status, purchased_in)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW())
		RETURNING id, purchased_in
	`

	created := payment
	err := r.db.QueryRowContext(
		ctx,
		query,
		payment.UserID,
		payment.WebsiteID,
		payment.ReferenceID,
		payment.Type,
		payment.Quantity,
		payment.Amount,
		payment.Currency,
		payment.Status,
	).Scan(&created.ID, &created.PurchasedIn)
	if err != nil {
		return nil, err
	}

	return &created, nil
}

func (r *connection) UpdatePaymentStatusByReference(referenceID, status string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `UPDATE payments SET status = $1 WHERE reference_id = $2`
	_, err := r.db.ExecContext(ctx, query, status, referenceID)
	return err
}
