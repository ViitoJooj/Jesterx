package postgresql

import (
	"database/sql"

	"github.com/ViitoJooj/Jesterx/pkg/dotenv"
)

func Conn(uri dotenv.PostgreSQL) (*sql.DB, error) {
	db, err := sql.Open("postgres", uri.URI)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
