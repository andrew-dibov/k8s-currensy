package configs

type AppConfig struct {
	AppPort string

	PostgresURL string

	KafkaBrokers string
	KafkaGroup   string
	KafkaTopic   string
}

func LoadAppConfig() *AppConfig {
	return &AppConfig{
		AppPort: ":" + loadEnv("APP_PORT", "8080"),

		PostgresURL: loadPostgresURL(),

		KafkaBrokers: loadEnv("KAFKA_BROKERS", ""),
		KafkaGroup:   loadEnv("KAFKA_GROUP", ""),
		KafkaTopic:   loadEnv("KAFKA_TOPIC", ""),
	}
}
