package postgresql

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

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

	var db *sql.DB
	var err error

	for attempt := 1; attempt <= 5; attempt++ {
		db, err = sql.Open("postgres", dsn)
		if err != nil {
			return nil, err
		}

		if err = db.Ping(); err == nil {
			return db, nil
		}

		db.Close()
		wait := time.Duration(attempt) * time.Second
		fmt.Printf("[postgresql] attempt %d/5 failed: %v — retrying in %v\n", attempt, err, wait)
		time.Sleep(wait)
	}

	return nil, fmt.Errorf("postgresql: all 5 connection attempts failed: %w", err)
}
