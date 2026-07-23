package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv string
	Port   string

	AccountServiceAddr string

	Auth AuthConfig
}

type AuthConfig struct {
	JwtSecret string
}

func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("Aviso: arquivo .env não encontrado, lendo variáveis do ambiente de sistema")
	}

	return &Config{
		AppEnv:             getEnv("APP_ENV", "development"),
		Port:               getEnv("PORT", "8080"),
		AccountServiceAddr: getEnv("ACCOUNT_SERVICE_ADDR", "0.0.0.0"),
		Auth: AuthConfig{
			JwtSecret: getEnv("JWT_SECRET", "super-secret-key-change-me"),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}

func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	valueStr := getEnv(key, "")
	if value, err := time.ParseDuration(valueStr); err == nil {
		return value
	}
	return defaultValue
}
