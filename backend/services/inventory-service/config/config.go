package config

import (
	"os"
	"time"
)

type Config struct {
	Port              string
	MongoURI          string
	RedisAddr         string
	RabbitMQURL       string
	JWTSecret         string
	JWTExpiry         time.Duration
	RefreshExpiry     time.Duration
	Environment       string
	ProductServiceURL string
	OrgServiceURL     string
	AuthServiceURL    string
}

func LoadConfig() *Config {
	return &Config{
		Port:              getEnv("PORT", "8006"),
		MongoURI:          getEnv("MONGO_URI", "mongodb://admin:admin123@localhost:27017"),
		RedisAddr:         getEnv("REDIS_ADDR", "localhost:6379"),
		RabbitMQURL:       getEnv("RABBITMQ_URL", "amqp://admin:admin123@localhost:5672/"),
		JWTSecret:         getEnv("JWT_SECRET", "your-super-secret-jwt-key-change-in-production"),
		JWTExpiry:         parseDuration(getEnv("JWT_EXPIRY", "15m"), 15*time.Minute),
		RefreshExpiry:     parseDuration(getEnv("REFRESH_TOKEN_EXPIRY", "168h"), 168*time.Hour),
		Environment:       getEnv("ENVIRONMENT", "development"),
		ProductServiceURL: getEnv("PRODUCT_SERVICE_URL", "http://product-service:8003"),
		OrgServiceURL:     getEnv("ORG_SERVICE_URL", "http://org-service:8002"),
		AuthServiceURL:    getEnv("AUTH_SERVICE_URL", "http://auth-service:8001"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func parseDuration(value string, defaultDuration time.Duration) time.Duration {
	if duration, err := time.ParseDuration(value); err == nil {
		return duration
	}
	return defaultDuration
}
