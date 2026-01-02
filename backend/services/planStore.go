package services

import (
	"context"
	"database/sql"
	"errors"
	"jesterx-core/config"
	"strings"

	"github.com/lib/pq"
)

type PlanConfig struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	PriceCents  int64    `json:"price_cents"`
	Description string   `json:"description"`
	Features    []string `json:"features"`
	SiteLimit   int      `json:"site_limit"`
}

var defaultPlanConfigs = []PlanConfig{
	{
		ID:          "free",
		Name:        "Free",
		PriceCents:  0,
		Description: "Plano gratuito para começar a testar páginas e lojas.",
		Features:    []string{"1 site", "Editor básico", "Hospedagem incluída"},
		SiteLimit:   0,
	},
	{
		ID:          "business",
		Name:        "Business",
		PriceCents:  4900,
		Description: "Para quem precisa sair do zero com rapidez.",
		Features:    []string{"1 site", "Páginas ilimitadas", "Templates básicos", "Suporte por email"},
		SiteLimit:   1,
	},
	{
		ID:          "pro",
		Name:        "Pro",
		PriceCents:  9900,
		Description: "Mais recursos e liberdade para escalar.",
		Features:    []string{"Até 10 sites", "Páginas ilimitadas", "Todos os templates", "Suporte prioritário", "Analytics avançado"},
		SiteLimit:   10,
	},
	{
		ID:          "enterprise",
		Name:        "Enterprise",
		PriceCents:  19900,
		Description: "Para equipes e operações que precisam de alta capacidade.",
		Features:    []string{"Até 50 sites", "Páginas ilimitadas", "Templates customizados", "Integrações e API dedicadas", "Suporte 24/7"},
		SiteLimit:   50,
	},
}

func defaultPlanByID(id string) (PlanConfig, bool) {
	for _, p := range defaultPlanConfigs {
		if p.ID == id {
			return p, true
		}
	}
	return PlanConfig{}, false
}

func EnsurePlanStore(ctx context.Context) error {
	if err := createPlanTable(ctx); err != nil {
		return err
	}
	return seedDefaultPlans(ctx)
}

func createPlanTable(ctx context.Context) error {
	_, err := config.DB.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS plan_configs (
			id VARCHAR(50) PRIMARY KEY,
			name VARCHAR(120) NOT NULL,
			price_cents BIGINT NOT NULL DEFAULT 0,
			description TEXT,
			features TEXT[] NOT NULL DEFAULT '{}',
			site_limit INT NOT NULL DEFAULT 0,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		);
	`)
	if err != nil {
		return err
	}

	return nil
}

func seedDefaultPlans(ctx context.Context) error {
	for _, plan := range defaultPlanConfigs {
		_, err := config.DB.ExecContext(ctx, `
			INSERT INTO plan_configs (id, name, price_cents, description, features, site_limit)
			VALUES ($1, $2, $3, $4, $5, $6)
			ON CONFLICT (id) DO UPDATE
			SET name = EXCLUDED.name,
				price_cents = EXCLUDED.price_cents,
				description = EXCLUDED.description,
				features = EXCLUDED.features,
				site_limit = EXCLUDED.site_limit,
				updated_at = NOW()
		`, plan.ID, plan.Name, plan.PriceCents, plan.Description, pq.Array(plan.Features), plan.SiteLimit)
		if err != nil {
			return err
		}
	}
	return nil
}

func ListPlanConfigs(ctx context.Context) ([]PlanConfig, error) {
	if err := EnsurePlanStore(ctx); err != nil {
		return nil, err
	}

	rows, err := config.DB.QueryContext(ctx, `
		SELECT id, name, price_cents, COALESCE(description, ''), features, site_limit
		FROM plan_configs
		ORDER BY price_cents ASC, id ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var plans []PlanConfig
	for rows.Next() {
		var plan PlanConfig
		if err := rows.Scan(&plan.ID, &plan.Name, &plan.PriceCents, &plan.Description, pq.Array(&plan.Features), &plan.SiteLimit); err != nil {
			return nil, err
		}
		plans = append(plans, plan)
	}

	if len(plans) == 0 {
		return defaultPlanConfigs, nil
	}

	return plans, nil
}

func GetPlanConfig(ctx context.Context, id string) (PlanConfig, error) {
	if err := EnsurePlanStore(ctx); err != nil {
		return PlanConfig{}, err
	}

	var plan PlanConfig
	err := config.DB.QueryRowContext(ctx, `
		SELECT id, name, price_cents, COALESCE(description, ''), features, site_limit
		FROM plan_configs
		WHERE id = $1
	`, id).Scan(&plan.ID, &plan.Name, &plan.PriceCents, &plan.Description, pq.Array(&plan.Features), &plan.SiteLimit)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			if def, ok := defaultPlanByID(id); ok {
				return def, nil
			}
		}
		return PlanConfig{}, err
	}

	return plan, nil
}

func UpdatePlanConfig(ctx context.Context, planID string, input PlanConfig) (PlanConfig, error) {
	if err := EnsurePlanStore(ctx); err != nil {
		return PlanConfig{}, err
	}

	if input.PriceCents < 0 || input.SiteLimit < 0 {
		return PlanConfig{}, errors.New("invalid plan values")
	}

	if strings.TrimSpace(input.Name) == "" {
		return PlanConfig{}, errors.New("plan name is required")
	}

	if strings.TrimSpace(planID) == "" {
		return PlanConfig{}, errors.New("plan id is required")
	}

	_, err := config.DB.ExecContext(ctx, `
		INSERT INTO plan_configs (id, name, price_cents, description, features, site_limit)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (id) DO UPDATE
		SET name = EXCLUDED.name,
			price_cents = EXCLUDED.price_cents,
			description = EXCLUDED.description,
			features = EXCLUDED.features,
			site_limit = EXCLUDED.site_limit,
			updated_at = NOW()
	`, planID, input.Name, input.PriceCents, input.Description, pq.Array(input.Features), input.SiteLimit)
	if err != nil {
		return PlanConfig{}, err
	}

	return GetPlanConfig(ctx, planID)
}
