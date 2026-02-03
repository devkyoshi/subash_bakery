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
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/yourusername/erp-system/services/org-service/config"
	"github.com/yourusername/erp-system/services/org-service/internal/handlers"
	"github.com/yourusername/erp-system/services/org-service/internal/repository"
	"github.com/yourusername/erp-system/services/org-service/internal/service"
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
	orgRepo := repository.NewOrganizationRepository(mongoDB.Database)
	companyRepo := repository.NewCompanyRepository(mongoDB.Database)
	locationRepo := repository.NewLocationRepository(mongoDB.Database)
	locationUserRepo := repository.NewLocationUserRepository(mongoDB.Database)

	// Initialize services
	orgService := service.NewOrganizationService(orgRepo, companyRepo, locationRepo, locationUserRepo)

	// Initialize handlers
	orgHandler := handlers.NewOrganizationHandler(orgService)

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
	orgHandler.RegisterRoutes(v1, jwtManager)

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

	log.Printf("Organization Service started on port %s", cfg.Port)

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

	// Organizations collection indexes
	orgsCollection := db.Collection("organizations")
	orgsIndexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "domain", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{{Key: "status", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "deleted_at", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "created_at", Value: -1}},
		},
	}

	_, err := orgsCollection.Indexes().CreateMany(ctx, orgsIndexes)
	if err != nil {
		return fmt.Errorf("failed to create organizations indexes: %w", err)
	}

	// Companies collection indexes
	companiesCollection := db.Collection("companies")
	companiesIndexes := []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "organization_id", Value: 1},
				{Key: "code", Value: 1},
			},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{{Key: "organization_id", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "deleted_at", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "created_at", Value: -1}},
		},
	}

	_, err = companiesCollection.Indexes().CreateMany(ctx, companiesIndexes)
	if err != nil {
		return fmt.Errorf("failed to create companies indexes: %w", err)
	}

	// Locations collection indexes
	locationsCollection := db.Collection("locations")
	locationsIndexes := []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "company_id", Value: 1},
				{Key: "code", Value: 1},
			},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{{Key: "company_id", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "organization_id", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "type", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "deleted_at", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "created_at", Value: -1}},
		},
	}

	_, err = locationsCollection.Indexes().CreateMany(ctx, locationsIndexes)
	if err != nil {
		return fmt.Errorf("failed to create locations indexes: %w", err)
	}

	// Location Users collection indexes
	locationUsersCollection := db.Collection("location_users")
	locationUsersIndexes := []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "user_id", Value: 1},
				{Key: "location_id", Value: 1},
			},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{{Key: "user_id", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "location_id", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "is_active", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "deleted_at", Value: 1}},
		},
	}

	_, err = locationUsersCollection.Indexes().CreateMany(ctx, locationUsersIndexes)
	if err != nil {
		return fmt.Errorf("failed to create location_users indexes: %w", err)
	}

	log.Println("Database indexes created successfully")
	return nil
}
