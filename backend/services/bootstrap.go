package services

import (
	"context"
	"jesterx-core/config"
	"strings"

	"github.com/lib/pq"
)

func SetupPlatformData(ctx context.Context) error {
	if err := ensureUserColumns(ctx); err != nil {
		return err
	}

	if err := EnsurePlanStore(ctx); err != nil {
		return err
	}

	if err := syncAdminRoles(ctx); err != nil {
		return err
	}

	return nil
}

func ensureUserColumns(ctx context.Context) error {
	statements := []string{
		`ALTER TABLE users ADD COLUMN IF NOT EXISTS role VARCHAR(50) NOT NULL DEFAULT 'platform_user';`,
		`ALTER TABLE users ADD COLUMN IF NOT EXISTS banned BOOLEAN NOT NULL DEFAULT FALSE;`,
	}

	for _, stmt := range statements {
		if _, err := config.DB.ExecContext(ctx, stmt); err != nil {
			return err
		}
	}

	return nil
}

func syncAdminRoles(ctx context.Context) error {
	if len(config.AdminEmails) == 0 {
		return nil
	}

	var normalized []string
	for _, email := range config.AdminEmails {
		email = strings.ToLower(strings.TrimSpace(email))
		if email != "" {
			normalized = append(normalized, email)
		}
	}

	if len(normalized) == 0 {
		return nil
	}

	_, err := config.DB.ExecContext(ctx, `
		UPDATE users
		SET role = 'platform_admin'
		WHERE lower(email) = ANY($1)
	`, pq.Array(normalized))

	return err
}
