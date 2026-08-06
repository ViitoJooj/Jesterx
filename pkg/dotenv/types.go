package dotenv

type Config struct {
	Application Application
	PostgreSQL  PostgreSQL
	Security    Security
}

type Application struct {
	Port       string
	Enviroment string
	ViewUrl    string
	DaemonUrl  string
}

type PostgreSQL struct {
	URI      string
	User     string
	Password string
	Host     string
	Port     string
	SSLMode  string
}

type Security struct {
	PasetoSecretKey string
}
