package configs

import "os"

type Config struct {
	Port            string
	CurrencyService string
	RedisAddr       string
	RedisPass       string
	RedisDB         int
}

func Load() *Config {
	return &Config{
		Port:            ":" + getEnv("PORT", "50052"),
		RedisDB:         0,
		RedisPass:       getEnv("REDIS_PASSWORD", ""),
		RedisAddr:       getEnv("REDIS_ADDR", "redis:6379"),
		CurrencyService: getEnv("CURRENCY_SERVICE", "currency:50051"),
	}
}

func getEnv(key, def string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return def
}
