package configs

type AppConfig struct {
	AppPort     string
	CurrencyURL string

	RedisDB   int
	RedisURL  string
	RedisPass string

	KafkaTopic   string
	KafkaBrokers string
}

func LoadAppConfig() *AppConfig {
	return &AppConfig{
		AppPort:     ":" + loadEnv("APP_PORT", "8080"),
		CurrencyURL: loadEnv("CURRENCY_ADDR", "currency") + ":" + loadEnv("CURRENCY_PORT", "8080"),

		RedisDB:   0,
		RedisURL:  loadEnv("REDIS_ADDR", "conversion-redis") + ":" + loadEnv("REDIS_PORT", "6379"),
		RedisPass: loadEnv("REDIS_PASS", ""),

		KafkaTopic:   loadEnv("KAFKA_TOPIC", "conversion"),
		KafkaBrokers: loadEnv("KAFKA_BROKERS", ""),
	}
}
