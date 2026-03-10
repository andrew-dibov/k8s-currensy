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
	envPort := os.Getenv("PORT")
	if envPort == "" {
		envPort = "8080"
	}

	envApiKeys := os.Getenv("API_KEYS")
	apiKeys := make(map[string]bool)

	if envApiKeys != "" {
		for _, key := range strings.Split(envApiKeys, ",") {
			apiKeys[strings.TrimSpace(key)] = true
		}
	} else {
		apiKeys = map[string]bool{
			"test-1111": true,
			"test-2222": true,
		}
	}

	envCurrencyService := os.Getenv("CURRENCY_SERVICE")
	if envCurrencyService == "" {
		envCurrencyService = "localhost:50051"
	}

	envConversionService := os.Getenv("CONVERSION_SERVICE")
	if envConversionService == "" {
		envConversionService = "localhost:50051"
	}

	envHistoryService := os.Getenv("HISTORY_SERVICE")
	if envHistoryService == "" {
		envHistoryService = "localhost:50051"
	}

	return &Config{
		Port:              ":" + envPort,
		APIKeys:           apiKeys,
		CurrencyService:   envCurrencyService,
		ConversionService: envConversionService,
		HistoryService:    envHistoryService,
	}
}
