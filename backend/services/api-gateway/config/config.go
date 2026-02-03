package config

import (
	"os"
)

type Config struct {
	Port             string
	AuthServiceURL   string
	OrgServiceURL    string
	ProductServiceURL string
	InventoryServiceURL string
	ProcurementServiceURL string
}

func LoadConfig() *Config {
	return &Config{
		Port:             getEnv("PORT", "8080"),
		AuthServiceURL:   getEnv("AUTH_SERVICE_URL", "http://auth-service:8001"),
		OrgServiceURL:    getEnv("ORG_SERVICE_URL", "http://org-service:8002"),
		ProductServiceURL: getEnv("PRODUCT_SERVICE_URL", "http://product-service:8003"),
		InventoryServiceURL: getEnv("INVENTORY_SERVICE_URL", "http://inventory-service:8004"),
		ProcurementServiceURL: getEnv("PROCUREMENT_SERVICE_URL", "http://procurement-service:8005"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
