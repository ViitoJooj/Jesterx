package postgresql

import (
	"database/sql"
	"strings"

	"github.com/ViitoJooj/verkoupe/pkg/dotenv"
	_ "github.com/lib/pq"
)

func Conn(uri dotenv.PostgreSQL) (*sql.DB, error) {
	dsn := uri.URI
	if !strings.Contains(dsn, "sslmode") {
		if strings.Contains(dsn, "?") {
			dsn += "&sslmode=disable"
		} else {
			dsn += "?sslmode=disable"
		}
	}
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
