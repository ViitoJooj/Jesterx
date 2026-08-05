package helpers

import (
	"database/sql"
	"errors"

	"github.com/ViitoJooj/verkoupe/internal/domain/entities"
)

func ScanPlans(rows *sql.Rows) ([]*domain.verkoupePlans, error) {
	var plans []*domain.verkoupePlans

	for rows.Next() {
		plan := &domain.verkoupePlans{}
		err := rows.Scan(
			&plan.UUID,
			&plan.Name,
			&plan.Description,
			&plan.MaxWebsites,
			&plan.MaxRouters,
			&plan.MaxProducts,
			&plan.CostPerSaleRate,
			&plan.Coin,
			&plan.Price,
			&plan.UpdatedAt,
			&plan.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		plans = append(plans, plan)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return plans, nil
}

func ScanPlan(row *sql.Row) (*domain.verkoupePlans, error) {
	plan := &domain.verkoupePlans{}

	err := row.Scan(
		&plan.UUID,
		&plan.Name,
		&plan.Description,
		&plan.MaxWebsites,
		&plan.MaxRouters,
		&plan.MaxProducts,
		&plan.CostPerSaleRate,
		&plan.Coin,
		&plan.Price,
		&plan.UpdatedAt,
		&plan.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("plan not found")
		}
		return nil, err
	}

	return plan, nil
}
