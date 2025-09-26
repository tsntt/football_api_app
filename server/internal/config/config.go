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
	Port string
}

func Load() *Config {
	return &Config{
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnvInt("DB_PORT", 5432),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASS", "postgres"),
			Name:     getEnv("DB_NAME", "football_api"),
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
			Port: getEnv("SERVER_PORT", "8080"),
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
