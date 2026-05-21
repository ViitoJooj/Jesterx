package postgres

import "database/sql"

type connection struct {
	db *sql.DB
}
