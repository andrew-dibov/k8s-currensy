package configs

import "os"

type Config struct {
	Port         string
	KafkaBrokers string
	KafkaTopic   string
	KafkaGroup   string
	PostgresDB   string
}

func Load() *Config {
	return &Config{
		Port:         ":" + getEnv("PORT", "8080"),
		KafkaGroup:   getEnv("KAFKA_GROUP", ""),
		KafkaTopic:   getEnv("KAFKA_TOPIC", ""),
		KafkaBrokers: getEnv("KAFKA_BROKERS", ""),
	}
}

func getEnv(key, def string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return def
}
