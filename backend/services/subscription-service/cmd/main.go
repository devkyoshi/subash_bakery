package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/yourusername/erp-system/services/subscription-service/config"
	"github.com/yourusername/erp-system/services/subscription-service/internal/handlers"
	"github.com/yourusername/erp-system/services/subscription-service/internal/repository"
	"github.com/yourusername/erp-system/services/subscription-service/internal/service"
	"github.com/yourusername/erp-system/shared/database"
	"github.com/yourusername/erp-system/shared/middleware"
	"github.com/yourusername/erp-system/shared/utils"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Connect to MongoDB
	mongoDB, err := database.NewMongoDB(cfg.MongoURI, "erp_db")
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer mongoDB.Close()

	// Connect to Redis
	redisClient, err := database.NewRedisClient(cfg.RedisAddr)
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	defer redisClient.Close()

	// Create indexes
	if err := createIndexes(mongoDB.Database); err != nil {
		log.Fatalf("Failed to create indexes: %v", err)
	}

	// Initialize JWT manager
	jwtManager := utils.NewJWTManager(cfg.JWTSecret, 15*time.Minute, 168*time.Hour)

	// Initialize repositories
	planRepo := repository.NewPlanRepository(mongoDB.Database)
	subscriptionRepo := repository.NewSubscriptionRepository(mongoDB.Database)

	// Initialize services
	subscriptionService := service.NewSubscriptionService(planRepo, subscriptionRepo)

	// Initialize handlers
	subscriptionHandler := handlers.NewSubscriptionHandler(subscriptionService)

	// Set Gin mode
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Create Gin router
	router := gin.Default()

	// Middleware
	router.Use(middleware.CORSMiddleware())

	// Rate limiting
	rateLimiter := middleware.NewRateLimiter(redisClient.Client, 100, time.Minute)
	router.Use(rateLimiter.Middleware())

	// Health check
	router.GET("/health", func(c *gin.Context) {
		utils.SuccessResponse(c, 200, gin.H{"status": "healthy"}, "Service is healthy")
	})

	// API routes
	v1 := router.Group("/api/v1")
	subscriptionHandler.RegisterRoutes(v1, jwtManager)

	// Start server
	srv := &http.Server{
		Addr:    "0.0.0.0:" + cfg.Port,
		Handler: router,
	}

	// Graceful shutdown
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	log.Printf("Subscription Service started on port %s", cfg.Port)

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exited")
}

func createIndexes(db *mongo.Database) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Subscription plans collection indexes
	plansCollection := db.Collection("subscription_plans")
	plansIndexes := []mongo.IndexModel{
		{
			Keys: bson.D{{Key: "tier", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "is_public", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "is_active", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "display_order", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "deleted_at", Value: 1}},
		},
	}

	_, err := plansCollection.Indexes().CreateMany(ctx, plansIndexes)
	if err != nil {
		return fmt.Errorf("failed to create plans indexes: %w", err)
	}

	// Organization subscriptions collection indexes
	subscriptionsCollection := db.Collection("organization_subscriptions")
	subscriptionsIndexes := []mongo.IndexModel{
		{
			Keys: bson.D{{Key: "organization_id", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "plan_id", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "status", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "next_billing_date", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "deleted_at", Value: 1}},
		},
	}

	_, err = subscriptionsCollection.Indexes().CreateMany(ctx, subscriptionsIndexes)
	if err != nil {
		return fmt.Errorf("failed to create subscriptions indexes: %w", err)
	}

	log.Println("Database indexes created successfully")
	return nil
}
