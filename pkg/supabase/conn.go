package supabase

import (
	"database/sql"
	"log"

	"github.com/ViitoJooj/Jesterx/pkg/dotenv"
	_ "github.com/jackc/pgx/v5/stdlib"
)

var DB *sql.DB
var openDB = sql.Open
var pingDB = func(db *sql.DB) error {
	return db.Ping()
}

func Conn() {
	var err error

	DB, err = openDB("pgx", dotenv.Supabase_uri)
	if err != nil {
		log.Panicf("Error opening database: %s", err)
	}

	if err = pingDB(DB); err != nil {
		log.Panicf("Error pinging database: %s", err)
	}
}
