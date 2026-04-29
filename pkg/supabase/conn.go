package supabase

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/ViitoJooj/Jesterx/pkg/dotenv"
	_ "github.com/jackc/pgx/v5/stdlib"
)

var DB *sql.DB

func Conn() {
	var err error

	DB, err = sql.Open("pgx", dotenv.SupabaseURI)
	if err != nil {
		log.Panicf("Error opening database: %s", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err = DB.PingContext(ctx); err != nil {
		log.Panicf("Error pinging database: %s", err)
	}
}
