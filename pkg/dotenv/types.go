package dotenv

type Config struct {
	PostgreSQL PostgreSQL
}

type PostgreSQL struct {
	URI      string
	User     string
	Password string
	Host     string
	Port     string
	SSLMode  string
}
