package configs

import (
	"os"
	"strings"
)

type Config struct {
	Port         string
	RatesService string
	APIKeys      map[string]bool
}

func Load() *Config {
	envPort := os.Getenv("PORT")
	if envPort == "" {
		envPort = "8080"
	}

	envRatesService := os.Getenv("RATES_SERVICE")
	if envRatesService == "" {
		envRatesService = "localhost:50051"
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
			"test-3333": true,
			"test-4444": true,
		}
	}

	return &Config{
		Port:         ":" + envPort,
		RatesService: envRatesService,
		APIKeys:      apiKeys,
	}
}
