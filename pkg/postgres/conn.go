package supabase

import (
	"database/sql"
	"fmt"

	"github.com/ViitoJooj/Jesterx/pkg/dotenv"
	_ "github.com/lib/pq"
)

var DB *sql.DB

func Conn(cfg dotenv.PostgresConfig) (*sql.DB, error) {
	strConn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.PostgresHost, cfg.PostgresPort, cfg.PostgresUser, cfg.PostgresPassword, cfg.PostgresDBName, cfg.PostgresSSLMode)

	DB, err := sql.Open("postgres", strConn)
	if err != nil {
		return nil, fmt.Errorf("error connecting to the database: %w", err)
	}

	if err = DB.Ping(); err != nil {
		return nil, fmt.Errorf("error establishing a connection to the database: %w", err)
	}

	return DB, nil

}
