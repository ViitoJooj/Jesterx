package mongo

import "database/sql"

var DB *sql.DB

func Conn() {
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
12	if err != nil {
13		log.Fatalf("Failed to connect to the database: %v", err)
14	}
}
