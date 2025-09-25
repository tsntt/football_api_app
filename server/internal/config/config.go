package config

import (
	"log"
	"os"
)

// Config armazena as configurações da aplicação
type Config struct {
	DatabaseURL    string
	JWTSecret      string
	FootballAPIKey string
	Port           string
}

// LoadConfig carrega as configurações a partir das variáveis de ambiente
func Load() *Config {
	cfg := &Config{
		DatabaseURL:    getEnv("DATABASE_URL", "user=postgres password=password dbname=futebol_api sslmode=disable"),
		JWTSecret:      getEnv("JWT_SECRET", "default_secret_key"),
		FootballAPIKey: getEnv("FOOTBALL_API_KEY", ""),
		Port:           getEnv("PORT", "4000"),
	}

	if cfg.FootballAPIKey == "" {
		log.Println("AVISO: FOOTBALL_API_KEY não está definida.")
	}

	return cfg
}

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}
