package config
package config

import (
	"os"
	"time"





































}	return defaultDuration	}		return duration	if duration, err := time.ParseDuration(value); err == nil {func parseDuration(value string, defaultDuration time.Duration) time.Duration {}	return defaultValue	}		return value	if value := os.Getenv(key); value != "" {func getEnv(key, defaultValue string) string {}	}		ProcurementServiceURL: getEnv("PROCUREMENT_SERVICE_URL", "http://procurement-service:8005"),		Environment:           getEnv("ENVIRONMENT", "development"),		RefreshExpiry:         getEnv("REFRESH_TOKEN_EXPIRY", "168h"),		JWTSecret:             getEnv("JWT_SECRET", "your-super-secret-jwt-key-change-in-production"),		RabbitMQURL:           getEnv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/"),		MongoURI:              getEnv("MONGO_URI", "mongodb://admin:admin123@localhost:27017"),		Port:                  getEnv("PORT", "8009"),	return &Config{func LoadConfig() *Config {}	ProcurementServiceURL string	Environment           string	RefreshExpiry         string	JWTSecret             string	RabbitMQURL           string	MongoURI              string	Port                  stringtype Config struct {)