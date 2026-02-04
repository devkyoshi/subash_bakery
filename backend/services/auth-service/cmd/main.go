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

	"github.com/yourusername/erp-system/services/auth-service/config"
	"github.com/yourusername/erp-system/services/auth-service/internal/handlers"
	"github.com/yourusername/erp-system/services/auth-service/internal/repository"
	"github.com/yourusername/erp-system/services/auth-service/internal/seed"
	"github.com/yourusername/erp-system/services/auth-service/internal/service"
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

	// Seed database with system permissions and roles
	seeder := seed.NewSeeder(
		repository.NewPermissionRepository(mongoDB.Database),
		repository.NewRoleRepository(mongoDB.Database),
	)
	seedCtx, seedCancel := context.WithTimeout(context.Background(), 30*time.Second)
	if err := seeder.SeedAll(seedCtx); err != nil {
		seedCancel()
		log.Fatalf("Failed to seed database: %v", err)
	}
	seedCancel()

	// Initialize JWT manager
	jwtManager := utils.NewJWTManager(cfg.JWTSecret, cfg.JWTExpiry, cfg.RefreshTokenExpiry)

	// Initialize repositories
	userRepo := repository.NewUserRepository(mongoDB.Database)
	sessionRepo := repository.NewSessionRepository(mongoDB.Database)
	roleRepo := repository.NewRoleRepository(mongoDB.Database)
	permissionRepo := repository.NewPermissionRepository(mongoDB.Database)

	// Initialize services
	authService := service.NewAuthService(userRepo, sessionRepo, roleRepo, permissionRepo, jwtManager, cfg)
	roleService := service.NewRoleService(roleRepo, permissionRepo, userRepo)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)
	roleHandler := handlers.NewRoleHandler(roleService)

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
	authHandler.RegisterRoutes(v1, jwtManager)
	roleHandler.RegisterRoutes(v1, jwtManager)

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

	log.Printf("Auth Service started on port %s", cfg.Port)

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

	// Users collection indexes
	usersCollection := db.Collection("users")
	usersIndexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "email", Value: 1}},
			Options: &options.IndexOptions{Unique: &[]bool{true}[0]},
		},
		{
			Keys: bson.D{{Key: "organization_id", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "google_id", Value: 1}},
		},
	}

	_, err := usersCollection.Indexes().CreateMany(ctx, usersIndexes)
	if err != nil {
		return fmt.Errorf("failed to create users indexes: %w", err)
	}

	// Sessions collection indexes
	sessionsCollection := db.Collection("sessions")
	sessionsIndexes := []mongo.IndexModel{
		{
			Keys: bson.D{{Key: "refresh_token", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "user_id", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "expires_at", Value: 1}},
		},
	}

	_, err = sessionsCollection.Indexes().CreateMany(ctx, sessionsIndexes)
	if err != nil {
		return fmt.Errorf("failed to create sessions indexes: %w", err)
	}

	log.Println("Database indexes created successfully")
	return nil
}
