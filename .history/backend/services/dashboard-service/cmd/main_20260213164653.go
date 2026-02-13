package main

import (
	"context"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/erp-system/services/dashboard-service/config"
	"github.com/yourusername/erp-system/services/dashboard-service/internal/handler"
	"github.com/yourusername/erp-system/services/dashboard-service/internal/repository"
	"github.com/yourusername/erp-system/services/dashboard-service/internal/service"
	"github.com/yourusername/erp-system/shared/middleware"
	"github.com/yourusername/erp-system/shared/rabbitmq"
	"github.com/yourusername/erp-system/shared/utils"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	cfg := config.LoadConfig()

	// Initialize RabbitMQ
	// Retry connection
	var rabbitClient *rabbitmq.RabbitMQClient
	var err error
	for i := 0; i < 10; i++ {
		rabbitClient, err = rabbitmq.NewRabbitMQClient(cfg.RabbitMQURL)
		if err == nil {
			break
		}
		log.Printf("Failed to connect to RabbitMQ, retrying in 5s: %v", err)
		time.Sleep(5 * time.Second)
	}
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer rabbitClient.Close()

	// Initialize MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.MongoURI))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer mongoClient.Disconnect(context.Background())
	db := mongoClient.Database("erp_db")

	// Initialize Repositories
	activityRepo := repository.NewActivityRepository(db)

	// Initialize Services
	activityService := service.NewActivityService(activityRepo)
	aggregationService := service.NewAggregationService(rabbitClient, activityService)

	// Initialize Handlers
	dashboardHandler := handler.NewDashboardHandler(aggregationService)
	eventHandler := handler.NewEventHandler(activityService)

	// Start RabbitMQ Consumer
	_, err = rabbitClient.DeclareQueue("dashboard.activities")
	if err != nil {
		log.Fatalf("Failed to declare activity queue: %v", err)
	}

	err = rabbitClient.Consume("dashboard.activities", eventHandler.HandleActivityEvent)
	if err != nil {
		log.Fatalf("Failed to start activity consumer: %v", err)
	}

	// Initialize JWT Manager
	refreshExpiry, _ := time.ParseDuration(cfg.RefreshExpiry)
	jwtManager := utils.NewJWTManager(cfg.JWTSecret, 24*time.Hour, refreshExpiry)

	// Setup Router
	router := gin.Default()
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, x-organization-id")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, PATCH, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	api := router.Group("/api/v1")
	api.Use(middleware.AuthMiddleware(jwtManager))

	dashboardHandler.RegisterRoutes(api)

	log.Printf("Dashboard Service starting on port %s", cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
