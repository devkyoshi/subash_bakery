package config

import (
	"log"
	"os"
)

type Config struct {
	Port                string
	MongoURI            string
	RedisAddr           string
	RabbitMQURL         string
	JWTSecret           string
	Environment         string
	InventoryServiceURL string
}

func LoadConfig() *Config {
	cfg := &Config{
		Port:                getEnv("PORT", "8003"),
		MongoURI:            getEnv("MONGO_URI", "mongodb://admin:admin123@localhost:27017"),
		RedisAddr:           getEnv("REDIS_ADDR", "localhost:6379"),
		RabbitMQURL:         getEnv("RABBITMQ_URL", "amqp://admin:admin123@localhost:5672/"),
		JWTSecret:           getEnv("JWT_SECRET", "your-secret-key-change-in-production"),
		Environment:         getEnv("ENV", "development"),
		InventoryServiceURL: getEnv("INVENTORY_SERVICE_URL", "http://inventory-service:8004"),
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
