package configs

import "os"

type Config struct {
	Port             string
	Postgres         string
	ExternalAPIURL   string
	ExternalAPIToken string
}

func Load() *Config {
	return &Config{
		Port:             ":" + getEnv("PORT", "50051"),
		Postgres:         getEnv("POSTGRES", "postgres://user:pass@localhost:5432/currency?sslmode=disable"),
		ExternalAPIURL:   getEnv("EXTERNAL_API_URL", "https://v6.exchangerate-api.com/v6/"),
		ExternalAPIToken: getEnv("EXTERNAL_API_TOKEN", ""),
	}
}

func getEnv(key, def string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return def
}
