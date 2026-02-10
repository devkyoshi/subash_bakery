package config

import (
	"os"
)

type Config struct {
	Port                    string
	MongoURI                string
	DBName                  string
	RabbitMQURL             string
	FirebaseCredentialsPath string
}

func LoadConfig() *Config {
	return &Config{
		Port:                    getEnv("PORT", "8007"),
		MongoURI:                getEnv("MONGO_URI", "mongodb://admin:admin123@localhost:27017"),
		DBName:                  getEnv("DB_NAME", "erp_db"),
		RabbitMQURL:             getEnv("RABBITMQ_URL", "amqp://admin:admin123@localhost:5672/"),
		FirebaseCredentialsPath: getEnv("FIREBASE_CREDENTIALS_PATH", "firebase-credentials.json"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
