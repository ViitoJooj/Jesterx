package helpers

import (
	"database/sql"
	"gen-you-ecommerce/config"
)

func UserHasTenantAccess(userID, tenantID string) (bool, string, error) {
	db := config.DB
	var role string

	err := db.QueryRow(`SELECT role FROM tenant_users WHERE user_id = $1 AND tenant_id = $2`, userID, tenantID).Scan(&role)

	if err == sql.ErrNoRows {
		return false, "", nil
	}
	if err != nil {
		return false, "", err
	}

	return true, role, nil
}
