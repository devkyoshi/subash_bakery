package config

import (
	"os"
)

type Config struct {
	Port                  string
	MongoURI              string
	RabbitMQURL           string
	JWTSecret             string
	RefreshExpiry         string
	Environment           string
	ProcurementServiceURL string
}

func LoadConfig() *Config {
	return &Config{
		Port:                  getEnv("PORT", "8009"),
		MongoURI:              getEnv("MONGO_URI", "mongodb://admin:admin123@localhost:27017"),
		RabbitMQURL:           getEnv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/"),
		JWTSecret:             getEnv("JWT_SECRET", "your-super-secret-jwt-key-change-in-production"),
		RefreshExpiry:         getEnv("REFRESH_TOKEN_EXPIRY", "168h"),
		Environment:           getEnv("ENVIRONMENT", "development"),
		ProcurementServiceURL: getEnv("PROCUREMENT_SERVICE_URL", "http://procurement-service:8005"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
