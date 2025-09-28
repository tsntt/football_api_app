package config

import (
	"os"
	"strconv"
)

type Config struct {
	Database    DatabaseConfig
	JWT         JWTConfig
	FootballAPI FootballAPIConfig
	Server      ServerConfig
	EmailAPI    EmailAPIConfig
	SMSAPI      SMSAPIConfig
}

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
	SSLMode  string
}

type JWTConfig struct {
	Secret       string
	ExpiresHours int
}

type FootballAPIConfig struct {
	Token string
	URL   string
}

type ServerConfig struct {
	Host string
	Port string
}

type EmailAPIConfig struct {
	APIKey string
	From   string
}

type SMSAPIConfig struct {
	AccountSID string
	APIKey     string
	From       string
}

func Load() *Config {
	return &Config{
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnvInt("DB_PORT", 5432),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASS", "change_this_password"),
			Name:     getEnv("DB_NAME", "football"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		JWT: JWTConfig{
			Secret:       getEnv("JWT_SECRET", "default-secret-key"),
			ExpiresHours: getEnvInt("JWT_EXPIRES_HOURS", 24),
		},
		FootballAPI: FootballAPIConfig{
			Token: getEnv("FOOTBALL_API_TOKEN", ""),
			URL:   getEnv("FOOTBALL_API_URL", "https://api.football-data.org/v4"),
		},
		Server: ServerConfig{
			Host: getEnv("SERVER_DOMAIN", "127.0.0.1"),
			Port: getEnv("SERVER_PORT", "4000"),
		},
		EmailAPI: EmailAPIConfig{
			APIKey: getEnv("MAILGUN_API_KEY", ""),
			From:   getEnv("MAILGUN_FROM", ""),
		},
		SMSAPI: SMSAPIConfig{
			AccountSID: getEnv("TWILIO_ACCOUNT_SID", ""),
			APIKey:     getEnv("TWILIO_API_KEY", ""),
			From:       getEnv("TWILIO_FROM", ""),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}
