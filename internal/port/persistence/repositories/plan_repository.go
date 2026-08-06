package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/ViitoJooj/verkoupe/internal/domain/entities"
	"github.com/ViitoJooj/verkoupe/internal/domain/repositories/contracts"
	"github.com/ViitoJooj/verkoupe/internal/port/persistence/helpers"
)

var _ contracts.PlanContract = (*PlanRepository)(nil)

type PlanRepository struct {
	db *sql.DB
}

func NewPlanRepository(db *sql.DB) *PlanRepository {
	return &PlanRepository{
		db: db,
	}
}

func (r *PlanRepository) CreatePlan(plan *domain.VerkoupePlan) (*domain.VerkoupePlan, error) {
	if plan == nil {
		return nil, errors.New("invalid plan")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `INSERT INTO plans (name, description, max_websites, max_routers, max_products, cost_per_sale_rate, coin, price)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	RETURNING uuid, created_at, updated_at`

	err := r.db.QueryRowContext(
		ctx,
		query,
		plan.Name,
		plan.Description,
		plan.MaxWebsites,
		plan.MaxRouters,
		plan.MaxProducts,
		plan.CostPerSaleRate,
		plan.Coin,
		plan.Price,
	).Scan(
		&plan.UUID,
		&plan.CreatedAt,
		&plan.UpdatedAt,
	)

	if err != nil {
		return nil, errors.New("could not create plan")
	}

	return plan, nil
}

func (r *PlanRepository) FindPlanByUUID(uuid string) (*domain.VerkoupePlan, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `SELECT uuid, name, description, max_websites, max_routers, max_products, cost_per_sale_rate, coin, price, updated_at, created_at
	FROM plans
	WHERE uuid = $1`

	row := r.db.QueryRowContext(ctx, query, uuid)
	return helpers.ScanPlan(row)
}

func (r *PlanRepository) FindPlanByName(name string) (*domain.VerkoupePlan, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `SELECT uuid, name, description, max_websites, max_routers, max_products, cost_per_sale_rate, coin, price, updated_at, created_at
	FROM plans
	WHERE name = $1`

	row := r.db.QueryRowContext(ctx, query, name)
	return helpers.ScanPlan(row)
}

func (r *PlanRepository) GetPlans() ([]*domain.VerkoupePlan, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `SELECT uuid, name, description, max_websites, max_routers, max_products, cost_per_sale_rate, coin, price, updated_at, created_at
	FROM plans`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return helpers.ScanPlans(rows)
}

func (r *PlanRepository) UpdatePlanByUUID(uuid string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `UPDATE plans SET updated_at = NOW() WHERE uuid = $1`

	result, err := r.db.ExecContext(ctx, query, uuid)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("plan not found")
	}

	return nil
}

func (r *PlanRepository) DeletePlanByUUID(uuid string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `DELETE FROM plans WHERE uuid = $1`

	result, err := r.db.ExecContext(ctx, query, uuid)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("plan not found")
	}

	return nil
}

func (r *PlanRepository) DeletePlansByUUIDS(uuids []string) error {
	if len(uuids) == 0 {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	placeholders := make([]string, len(uuids))
	args := make([]interface{}, len(uuids))

	for i, u := range uuids {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = u
	}

	query := fmt.Sprintf(`DELETE FROM plans WHERE uuid IN (%s)`, strings.Join(placeholders, ", "))

	_, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}
