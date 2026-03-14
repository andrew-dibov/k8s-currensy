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
		Port:            ":" + getEnv("PORT", "8080"),
		RedisDB:         0,
		RedisAddr:       getEnv("REDIS_ADDR", ""),
		RedisPass:       getEnv("REDIS_PASS", ""),
		CurrencyService: getEnv("CURRENCY_SERVICE", ""),
	}
}

func getEnv(key, def string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return def
}
