package supabase

import (
	"database/sql"
	"log"

	"github.com/ViitoJooj/Jesterx/pkg/dotenv"
	_ "github.com/jackc/pgx/v5/stdlib"
)

var DB *sql.DB

func Conn() {
	var err error

	DB, err = sql.Open("pgx", dotenv.Supabase_uri)
	if err != nil {
		log.Panicf("Error opening database: %s", err)
	}

	if err = DB.Ping(); err != nil {
		log.Panicf("Error pinging database: %s", err)
	}
}
