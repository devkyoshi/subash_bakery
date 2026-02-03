package config

import (
	"log"
	"os"
	"time"
)

type Config struct {
	Port                string
	MongoURI            string
	RedisAddr           string
	JWTSecret           string
	JWTExpiry           time.Duration
	RefreshTokenExpiry  time.Duration
	GoogleClientID      string
	GoogleClientSecret  string
	GoogleRedirectURL   string
	Environment         string
}

func LoadConfig() *Config {
	jwtExpiry, _ := time.ParseDuration(getEnv("JWT_EXPIRY", "15m"))
	refreshExpiry, _ := time.ParseDuration(getEnv("REFRESH_TOKEN_EXPIRY", "168h"))

	cfg := &Config{
		Port:                getEnv("PORT", "8001"),
		MongoURI:            getEnv("MONGO_URI", "mongodb://admin:admin123@localhost:27017"),
		RedisAddr:           getEnv("REDIS_ADDR", "localhost:6379"),
		JWTSecret:           getEnv("JWT_SECRET", "your-secret-key-change-in-production"),
		JWTExpiry:           jwtExpiry,
		RefreshTokenExpiry:  refreshExpiry,
		GoogleClientID:      getEnv("GOOGLE_CLIENT_ID", ""),
		GoogleClientSecret:  getEnv("GOOGLE_CLIENT_SECRET", ""),
		GoogleRedirectURL:   getEnv("GOOGLE_REDIRECT_URL", "http://localhost:8080/auth/google/callback"),
		Environment:         getEnv("ENV", "development"),
	}

	log.Printf("Config loaded: Port=%s, Environment=%s", cfg.Port, cfg.Environment)
	return cfg
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
