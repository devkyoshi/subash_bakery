package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/erp-system/services/notification-service/config"
	"github.com/yourusername/erp-system/services/notification-service/internal/handler"
	"github.com/yourusername/erp-system/services/notification-service/internal/repository"
	"github.com/yourusername/erp-system/services/notification-service/internal/service"
	"github.com/yourusername/erp-system/shared/middleware"
	"github.com/yourusername/erp-system/shared/rabbitmq"
	"github.com/yourusername/erp-system/shared/utils"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	cfg := config.LoadConfig()

	// Initialize MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.MongoURI))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer mongoClient.Disconnect(context.Background())

	if err := mongoClient.Ping(ctx, nil); err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
	}
	log.Println("Connected to MongoDB successfully")

	db := mongoClient.Database(cfg.DBName)

	// Initialize RabbitMQ
	var rabbitClient *rabbitmq.RabbitMQClient
	var rabbitErr error
	for i := 0; i < 30; i++ {
		rabbitClient, rabbitErr = rabbitmq.NewRabbitMQClient(cfg.RabbitMQURL)
		if rabbitErr == nil {
			log.Println("Connected to RabbitMQ successfully")
			break
		}
		log.Printf("Failed to connect to RabbitMQ (attempt %d/30): %v", i+1, rabbitErr)
		time.Sleep(2 * time.Second)
	}

	if rabbitErr != nil {
		log.Fatalf("Failed to connect to RabbitMQ after retries: %v", rabbitErr)
	}
	defer rabbitClient.Close()

	// Initialize Repositories
	deviceRepo := repository.NewDeviceRepository(db)
	notifRepo := repository.NewNotificationRepository(db)

	// Initialize Services
	notifService, err := service.NewNotificationService(cfg, deviceRepo, notifRepo)
	if err != nil {
		log.Printf("Warning: Failed to initialize Notification Service (Firebase): %v", err)
		// Proceeding without Firebase for testing or partial functionality
		// In production, might want to fail hard if notifications are critical
	}

	// Initialize Handlers
	deviceHandler := handler.NewDeviceHandler(deviceRepo, notifService)
	eventHandler := handler.NewEventHandler(rabbitClient, notifService)
	notifHandler := handler.NewNotificationHandler(notifService)

	// Start RabbitMQ Consumer
	go func() {
		if err := eventHandler.Start(context.Background()); err != nil {
			log.Printf("Failed to start event handler: %v", err)
		}
	}()

	// Initialize JWT Manager
	jwtManager := utils.NewJWTManager(cfg.JWTSecret, 15*time.Minute, 168*time.Hour)

	// Setup Gin Router
	router := gin.Default()

	api := router.Group("/api/v1")
	api.Use(middleware.AuthMiddleware(jwtManager))
	deviceHandler.RegisterRoutes(api)
	notifHandler.RegisterRoutes(api)

	// Health Check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "healthy",
			"service": "notification-service",
		})
	})

	// Start Server
	port := fmt.Sprintf("0.0.0.0:%s", cfg.Port)
	log.Printf("Notification Service starting on port %s", cfg.Port)
	if err := router.Run(port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
