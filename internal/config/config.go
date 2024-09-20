package config

import "os"

type Config struct {
	DatabaseURI string
}

func LoadConfig() Config {
	return Config{
		DatabaseURI: getEnvOrDefault("DATABASE_URI", "receipts.db"),
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
