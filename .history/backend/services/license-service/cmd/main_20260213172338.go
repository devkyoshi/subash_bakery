package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/yourusername/erp-system/services/license-service/config"
	"github.com/yourusername/erp-system/services/license-service/internal/handlers"
	"github.com/yourusername/erp-system/services/license-service/internal/repository"
	"github.com/yourusername/erp-system/services/license-service/internal/service"
	"github.com/yourusername/erp-system/shared/middleware"
	"github.com/yourusername/erp-system/shared/utils"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.MongoURI))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer client.Disconnect(context.Background())

	// Ping MongoDB
	if err := client.Ping(ctx, nil); err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
	}

	log.Println("Connected to MongoDB successfully")

	// Get database
	db := client.Database("erp_db")

	// Create indexes
	if err := createIndexes(db); err != nil {
		log.Fatalf("Failed to create indexes: %v", err)
	}

	// Initialize repositories
	appRepo := repository.NewApplicationRepository(db)
	licenseRepo := repository.NewLicenseRepository(db)
	userAssignRepo := repository.NewUserAssignmentRepository(db)
	deviceAssignRepo := repository.NewDeviceAssignmentRepository(db)

	// Initialize services
	appService := service.NewApplicationService(appRepo)
	licenseService := service.NewLicenseService(licenseRepo, appRepo, userAssignRepo, deviceAssignRepo)

	// Initialize handlers
	appHandler := handlers.NewApplicationHandler(appService)
	licenseHandler := handlers.NewLicenseHandler(licenseService)

	// Initialize JWT manager
	jwtManager := utils.NewJWTManager(cfg.JWTSecret, cfg.JWTExpiry, cfg.RefreshExpiry)

	// Setup Gin router
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	// Middleware
	router.Use(middleware.CORSMiddleware())
	// router.Use(middleware.RateLimitMiddleware())

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "healthy",
			"service": "license-service",
			"time":    time.Now(),
		})
	})

	// API routes
	api := router.Group("/api/v1")
	licenseHandler.RegisterRoutes(api, appHandler, jwtManager)

	// Start server
	port := fmt.Sprintf("0.0.0.0:%s", cfg.Port)
	log.Printf("License Service starting on port %s", cfg.Port)
	if err := router.Run(port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func createIndexes(db *mongo.Database) error {
	ctx := context.Background()

	// Applications indexes
	appIndexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "code", Value: 1}},
			Options: options.Index().SetUnique(true).SetName("idx_app_code"),
		},
		{
			Keys:    bson.D{{Key: "category", Value: 1}, {Key: "is_active", Value: 1}},
			Options: options.Index().SetName("idx_app_category_active"),
		},
		{
			Keys:    bson.D{{Key: "is_public", Value: 1}, {Key: "is_active", Value: 1}},
			Options: options.Index().SetName("idx_app_public_active"),
		},
		{
			Keys:    bson.D{{Key: "deleted_at", Value: 1}},
			Options: options.Index().SetName("idx_app_deleted"),
		},
	}

	if _, err := db.Collection("applications").Indexes().CreateMany(ctx, appIndexes); err != nil {
		return fmt.Errorf("failed to create application indexes: %w", err)
	}

	// Location licenses indexes
	licenseIndexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "organization_id", Value: 1}, {Key: "deleted_at", Value: 1}},
			Options: options.Index().SetName("idx_license_org"),
		},
		{
			Keys:    bson.D{{Key: "location_id", Value: 1}, {Key: "status", Value: 1}},
			Options: options.Index().SetName("idx_license_location_status"),
		},
		{
			Keys:    bson.D{{Key: "location_id", Value: 1}, {Key: "application_id", Value: 1}},
			Options: options.Index().SetName("idx_license_location_app"),
		},
		{
			Keys:    bson.D{{Key: "application_id", Value: 1}, {Key: "status", Value: 1}},
			Options: options.Index().SetName("idx_license_app_status"),
		},
		{
			Keys:    bson.D{{Key: "status", Value: 1}, {Key: "expires_at", Value: 1}},
			Options: options.Index().SetName("idx_license_status_expiry"),
		},
		{
			Keys:    bson.D{{Key: "license_key", Value: 1}},
			Options: options.Index().SetUnique(true).SetName("idx_license_key"),
		},
		{
			Keys:    bson.D{{Key: "deleted_at", Value: 1}},
			Options: options.Index().SetName("idx_license_deleted"),
		},
	}

	if _, err := db.Collection("location_licenses").Indexes().CreateMany(ctx, licenseIndexes); err != nil {
		return fmt.Errorf("failed to create license indexes: %w", err)
	}

	// User license assignments indexes
	userAssignIndexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "license_id", Value: 1}, {Key: "is_active", Value: 1}},
			Options: options.Index().SetName("idx_user_assign_license_active"),
		},
		{
			Keys:    bson.D{{Key: "user_id", Value: 1}, {Key: "is_active", Value: 1}},
			Options: options.Index().SetName("idx_user_assign_user_active"),
		},
		{
			Keys:    bson.D{{Key: "user_id", Value: 1}, {Key: "application_id", Value: 1}, {Key: "is_active", Value: 1}},
			Options: options.Index().SetName("idx_user_assign_user_app_active"),
		},
		{
			Keys:    bson.D{{Key: "organization_id", Value: 1}},
			Options: options.Index().SetName("idx_user_assign_org"),
		},
	}

	if _, err := db.Collection("user_license_assignments").Indexes().CreateMany(ctx, userAssignIndexes); err != nil {
		return fmt.Errorf("failed to create user assignment indexes: %w", err)
	}

	// Device license assignments indexes
	deviceAssignIndexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "license_id", Value: 1}, {Key: "is_active", Value: 1}},
			Options: options.Index().SetName("idx_device_assign_license_active"),
		},
		{
			Keys:    bson.D{{Key: "device_id", Value: 1}, {Key: "license_id", Value: 1}},
			Options: options.Index().SetName("idx_device_assign_device_license"),
		},
		{
			Keys:    bson.D{{Key: "organization_id", Value: 1}},
			Options: options.Index().SetName("idx_device_assign_org"),
		},
		{
			Keys:    bson.D{{Key: "is_online", Value: 1}, {Key: "last_seen_at", Value: 1}},
			Options: options.Index().SetName("idx_device_assign_online"),
		},
	}

	if _, err := db.Collection("device_license_assignments").Indexes().CreateMany(ctx, deviceAssignIndexes); err != nil {
		return fmt.Errorf("failed to create device assignment indexes: %w", err)
	}

	log.Println("Database indexes created successfully")
	return nil
}
