package services

import (
	"database/sql"
	"jesterx-core/config"
)

// resolveUserTenant tries to find a tenant for the user when no header was provided.
// It returns the first tenant_id where the user has any role.
func resolveUserTenant(userID string) (string, error) {
	var tenantID string
	err := config.DB.QueryRow(`SELECT tenant_id FROM tenant_users WHERE user_id = $1 ORDER BY tenant_id LIMIT 1`, userID).Scan(&tenantID)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}
		return "", err
	}
	return tenantID, nil
}
