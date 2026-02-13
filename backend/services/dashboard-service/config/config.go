package config

import (
	"os"
)

type Config struct {
	Port          string
	RabbitMQURL   string
	MongoURI      string
	JWTSecret     string
	RefreshExpiry string
}

func LoadConfig() *Config {
	return &Config{
		Port:          getEnv("PORT", "8008"),
		RabbitMQURL:   getEnv("RABBITMQ_URL", "amqp://guest:guest@rabbitmq:5672/"),
		MongoURI:      getEnv("MONGO_URI", "mongodb://mongodb:27017"),
		JWTSecret:     getEnv("JWT_SECRET", "your_secret_key"),
		RefreshExpiry: getEnv("REFRESH_TOKEN_EXPIRY", "168h"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
