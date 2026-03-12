package configs

import "os"

type Config struct {
	Port        string
	APIToken    string
	API         string
	PostgresURL string
}

func Load() *Config {
	return &Config{
		Port:        ":" + getEnv("PORT", "50051"),
		APIToken:    getEnv("API_TOKEN", ""),
		API:         getEnv("API", "https://v6.exchangerate-api.com/v6/"),
		PostgresURL: getEnv("POSTGRES_URL", "postgres://currency:pass1234@postgres:5432/currency?sslmode=disable"),
	}
}

func getEnv(key, def string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return def
}
