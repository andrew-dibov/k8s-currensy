package configs

import "os"

type Config struct {
	Port        string
	APIToken    string
	ExternalAPI string
	PostgresDB  string
}

func Load() *Config {
	return &Config{
		Port:        ":" + getEnv("PORT", ""),
		PostgresDB:  getEnv("POSTGRES_DB", ""),
		APIToken:    getEnv("EXTERNAL_API_TOKEN", ""),
		ExternalAPI: getEnv("EXTERNAL_API", "https://v6.exchangerate-api.com/v6/"),
	}
}

func getEnv(key, def string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return def
}
