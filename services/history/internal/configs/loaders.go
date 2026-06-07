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
	user := loadEnv("POSTGRES_USER", "history")
	pass := url.QueryEscape(loadEnv("POSTGRES_PASS", ""))
	host := loadEnv("POSTGRES_HOST", "postgres-history")
	port := loadEnv("POSTGRES_PORT", "5432")
	db := loadEnv("POSTGRES_DB", "history")

	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, pass, host, port, db)
}
