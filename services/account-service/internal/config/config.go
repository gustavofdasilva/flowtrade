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

	DB   DBConfig
	Auth AuthConfig
}

type DBConfig struct {
	Host                   string
	Port                   string
	User                   string
	Pass                   string
	Name                   string
	SSLMode                string
	MaxOpenConn            int
	MaxIdleConn            int
	ConnMaxLifetimeMinutes int
}

type AuthConfig struct {
	JWTSecret            string
	JWTTokenDuration     time.Duration
	RefreshTokenDuration time.Duration
}

func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("Aviso: arquivo .env não encontrado, lendo variáveis do ambiente de sistema")
	}

	return &Config{
		AppEnv: getEnv("APP_ENV", "development"),
		Port:   getEnv("PORT", "8080"),

		DB: DBConfig{
			Host:                   getEnv("DB_HOST", "localhost"),
			Port:                   getEnv("DB_PORT", "5432"),
			User:                   getEnv("DB_USER", "postgres"),
			Pass:                   getEnv("DB_PASS", "postgres"),
			Name:                   getEnv("DB_NAME", "account_db"),
			SSLMode:                getEnv("DB_SSLMODE", "disable"),
			MaxOpenConn:            getEnvAsInt("DB_MAX_OPEN_CONNS", 25),
			MaxIdleConn:            getEnvAsInt("DB_MAX_IDLE_CONNS", 10),
			ConnMaxLifetimeMinutes: getEnvAsInt("DB_CONN_MAX_LIFETIME_MIN", 15),
		},

		Auth: AuthConfig{
			JWTSecret:            getEnv("JWT_SECRET", "super-secret-key-change-me"),
			JWTTokenDuration:     getEnvAsDuration("JWT_TOKEN_DURATION", 15*time.Minute),
			RefreshTokenDuration: getEnvAsDuration("REFRESH_TOKEN_DURATION", 7*24*time.Hour),
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
