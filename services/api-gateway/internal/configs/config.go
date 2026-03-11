package configs

import (
	"os"
	"strings"
)

type Config struct {
	Port              string
	HistoryService    string
	CurrencyService   string
	ConversionService string
	APIKeys           map[string]bool
}

func Load() *Config {
	return &Config{
		Port:              ":" + getEnv("PORT", "8080"),
		HistoryService:    getEnv("HISTORY_SERVICE", "localhost:50051"),
		CurrencyService:   getEnv("CURRENCY_SERVICE", "localhost:50051"),
		ConversionService: getEnv("CONVERSION_SERVICE", "localhost:50051"),
		APIKeys:           parseKeys("API_KEYS", map[string]bool{"test-1111": true, "test-2222": true}),
	}
}

func getEnv(key, def string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return def
}

func parseKeys(key string, def map[string]bool) map[string]bool {
	str := os.Getenv(key)
	keys := make(map[string]bool)

	if str != "" {
		for _, key := range strings.Split(str, ",") {
			keys[strings.TrimSpace(key)] = true
		}
	} else {
		keys = def
	}

	return keys
}
