package configs

import (
	"os"
	"strings"
)

type Config struct {
	Port              string
	APIKeys           map[string]bool
	CurrencyService   string
	ConversionService string
	HistoryService    string
}

func Load() *Config {
	apiKeysString := os.Getenv("API_KEYS")
	apiKeys := make(map[string]bool)

	if apiKeysString != "" {
		for _, key := range strings.Split(apiKeysString, ",") {
			apiKeys[strings.TrimSpace(key)] = true
		}
	} else {
		apiKeys = map[string]bool{
			"test-1111": true,
			"test-2222": true,
		}
	}

	return &Config{
		APIKeys:           apiKeys,
		Port:              ":" + getEnv("PORT", "8080"),
		CurrencyService:   getEnv("CURRENCY_SERVICE", "localhost:50051"),
		ConversionService: getEnv("CONVERSION_SERVICE", "localhost:50051"),
		HistoryService:    getEnv("HISTORY_SERVICE", "localhost:50051"),
	}
}

func getEnv(key, def string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return def
}
