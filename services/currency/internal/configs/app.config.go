package configs

type AppConfig struct {
	AppPort string

	ExternalURL   string
	ExternalToken string

	PostgresURL string
}

func LoadAppConfig() *AppConfig {
	return &AppConfig{
		AppPort: ":" + loadEnv("APP_PORT", "8080"),

		ExternalURL:   loadEnv("EXTERNAL_URL", "https://v6.exchangerate-api.com/v6/"),
		ExternalToken: loadEnv("EXTERNAL_TOKEN", ""),

		PostgresURL: loadPostgresURL(),
	}
}
