package dotenv

type Config struct {
	Application Application
	PostgreSQL  PostgreSQL
}

type Application struct {
	Port       string
	Enviroment string
}

type PostgreSQL struct {
	URI      string
	User     string
	Password string
	Host     string
	Port     string
	SSLMode  string
}
