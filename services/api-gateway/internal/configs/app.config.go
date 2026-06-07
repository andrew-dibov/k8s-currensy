package configs

type AppConfig struct {
	AppPort string
	AppKeys map[string]bool

	HistoryURL    string
	CurrencyURL   string
	ConversionURL string
}

func LoadAppConfig() *AppConfig {
	return &AppConfig{
		AppPort: ":" + loadEnv("APP_PORT", "8080"),
		AppKeys: loadKeys("APP_KEYS", map[string]bool{}),

		HistoryURL:    loadEnv("HISTORY_ADDR", "history") + ":" + loadEnv("HISTORY_PORT", "8080"),
		CurrencyURL:   loadEnv("CURRENCY_ADDR", "currency") + ":" + loadEnv("CURRENCY_PORT", "8080"),
		ConversionURL: loadEnv("CONVERSION_ADDR", "conversion") + ":" + loadEnv("CONVERSION_PORT", "8080"),
	}
}
