package configs

import (
	"fmt"
	"net/url"
	"os"
)

func loadEnv(k string, d string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return d
}

func loadPostgresURL() string {
	user := loadEnv("POSTGRES_USER", "currency")
	pass := url.QueryEscape(loadEnv("POSTGRES_PASS", ""))
	host := loadEnv("POSTGRES_HOST", "postgres-currency")
	port := loadEnv("POSTGRES_PORT", "5432")
	db := loadEnv("POSTGRES_DB", "currency")

	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, pass, host, port, db)
}
