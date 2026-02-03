package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/erp-system/services/product-service/config"
	"github.com/yourusername/erp-system/services/product-service/internal/client"
	"github.com/yourusername/erp-system/services/product-service/internal/handlers"
	"github.com/yourusername/erp-system/services/product-service/internal/repository"
	"github.com/yourusername/erp-system/services/product-service/internal/service"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/yourusername/erp-system/shared/database"
	"github.com/yourusername/erp-system/shared/middleware"
	"github.com/yourusername/erp-system/shared/utils"
	"go.mongodb.org/mongo-driver/bson"
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

	if os.Getenv("ENVIRONMENT") == "development" {
		if err := createIndexes(mongoDB.Database); err != nil {
			log.Fatalf("Failed to create indexes: %v", err)
		}
	}

	// Initialize JWT manager
	jwtManager := utils.NewJWTManager(cfg.JWTSecret, 15*time.Minute, 168*time.Hour)

	// Initialize repositories
	productRepo := repository.NewProductRepository(mongoDB.Database)
	brandRepo := repository.NewBrandRepository(mongoDB.Database)
	categoryRepo := repository.NewCategoryRepository(mongoDB.Database)
	orgRepo := repository.NewOrganizationRepository(mongoDB.Database)
	unitRepo := repository.NewUnitRepository(mongoDB.Database)
	stockRepo := repository.NewStockLevelRepository(mongoDB.Database)

	// Initialize Clients
	inventoryClient := client.NewInventoryClient(cfg)

	// Initialize services
	productService := service.NewProductService(productRepo, categoryRepo, brandRepo, orgRepo, unitRepo, stockRepo, inventoryClient)
	brandService := service.NewBrandService(brandRepo, orgRepo)
	categoryService := service.NewCategoryService(categoryRepo, orgRepo)
	unitService := service.NewUnitService(unitRepo)

	// Initialize handlers
	productHandler := handlers.NewProductHandler(productService)
	brandHandler := handlers.NewBrandHandler(brandService)
	categoryHandler := handlers.NewCategoryHandler(categoryService)
	unitHandler := handlers.NewUnitHandler(unitService)

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
	productHandler.RegisterProductRoutes(v1, jwtManager)
	brandHandler.RegisterBrandRoutes(v1, jwtManager)
	categoryHandler.RegisterCategoryRoutes(v1, jwtManager)
	unitHandler.RegisterUnitRoutes(v1, jwtManager)

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

	log.Printf("Product Service started on port %s", cfg.Port)

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
	ctx := context.Background()

	// Brand collection indexes
	brandCollection := db.Collection("brands")

	// Index on organization_id and name (unique within organization)
	_, err := brandCollection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{
			{Key: "organization_id", Value: 1},
			{Key: "name", Value: 1},
		},
		Options: options.Index().
			SetUnique(true).
			SetPartialFilterExpression(bson.M{
				"deleted_at": nil,
			}),
	})
	if err != nil {
		return err
	}

	// Index on organization_id and code (unique within organization)
	_, err = brandCollection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{
			{Key: "organization_id", Value: 1},
			{Key: "code", Value: 1},
		},
		Options: options.Index().
			SetUnique(true).
			SetPartialFilterExpression(bson.M{
				"deleted_at": nil,
				"code":       bson.M{"$exists": true},
			}),
	})
	if err != nil {
		return err
	}

	// Index on organization_id for listing
	_, err = brandCollection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{
			{Key: "organization_id", Value: 1},
			{Key: "deleted_at", Value: 1},
		},
	})
	if err != nil {
		return err
	}

	// Product collection indexes
	productCollection := db.Collection("products")

	// Index on organization_id and SKU (unique within organization)
	_, err = productCollection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{
			{Key: "organization_id", Value: 1},
			{Key: "sku", Value: 1},
		},
		Options: options.Index().
			SetUnique(true).
			SetPartialFilterExpression(bson.M{
				"deleted_at": nil,
			}),
	})
	if err != nil {
		return err
	}

	// Index on brand_id for checking brand usage
	_, err = productCollection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{
			{Key: "brand_id", Value: 1},
			{Key: "deleted_at", Value: 1},
		},
	})
	if err != nil {
		return err
	}

	// Category collection indexes
	categoryCollection := db.Collection("categories")

	// Index on organization_id and name (unique within same parent)
	_, err = categoryCollection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{
			{Key: "organization_id", Value: 1},
			{Key: "parent_id", Value: 1},
			{Key: "name", Value: 1},
		},
		Options: options.Index().
			SetUnique(true).
			SetPartialFilterExpression(bson.M{
				"deleted_at": nil,
			}),
	})
	if err != nil {
		return err
	}

	// Index on organization_id and path (unique paths)
	_, err = categoryCollection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{
			{Key: "organization_id", Value: 1},
			{Key: "path", Value: 1},
		},
		Options: options.Index().
			SetUnique(true).
			SetPartialFilterExpression(bson.M{
				"deleted_at": nil,
			}),
	})
	if err != nil {
		return err
	}

	// Index on parent_id for finding children
	_, err = categoryCollection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{
			{Key: "parent_id", Value: 1},
			{Key: "deleted_at", Value: 1},
		},
	})
	if err != nil {
		return err
	}

	// Index on organization_id and level for filtering
	_, err = categoryCollection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{
			{Key: "organization_id", Value: 1},
			{Key: "level", Value: 1},
			{Key: "deleted_at", Value: 1},
		},
	})
	if err != nil {
		return err
	}

	// Index on category_id in products for checking category usage
	_, err = productCollection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{
			{Key: "category_id", Value: 1},
			{Key: "deleted_at", Value: 1},
		},
	})
	if err != nil {
		return err
	}

	return nil
}
