package configs

import (
	"os"
)

func loadEnv(k string, d string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return d
}
