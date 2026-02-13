package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/yourusername/erp-system/services/procurement-service/config"
	"github.com/yourusername/erp-system/services/procurement-service/internal/client"
	"github.com/yourusername/erp-system/services/procurement-service/internal/handlers"
	"github.com/yourusername/erp-system/services/procurement-service/internal/repository"
	"github.com/yourusername/erp-system/services/procurement-service/internal/service"
	"github.com/yourusername/erp-system/shared/middleware"
	"github.com/yourusername/erp-system/shared/rabbitmq"
	"github.com/yourusername/erp-system/shared/utils"
)

func main() {
	cfg := config.LoadConfig()
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
	db := mongoClient.Database("erp_procurement")

	// Create indexes
	if err := createIndexes(db); err != nil {
		log.Fatalf("Failed to create indexes: %v", err)
	}

	// Initialize repositories
	supplierRepo := repository.NewSupplierRepository(db)
	poRepo := repository.NewPurchaseOrderRepository(db)
	grnRepo := repository.NewGRNRepository(db)

	// Initialize clients
	productClient := client.NewProductClient(cfg)
	userClient := client.NewUserClient(cfg)
	inventoryClient := client.NewInventoryClient(cfg)

	// Initialize RabbitMQ
	rabbitClient, err := rabbitmq.NewRabbitMQClient(cfg.RabbitMQURL)
	if err != nil {
		log.Printf("Failed to connect to RabbitMQ: %v", err)
	} else {
		defer rabbitClient.Close()
	}

	// Initialize services
	procurementService := service.NewProcurementService(supplierRepo, poRepo, grnRepo, productClient, userClient, inventoryClient, rabbitClient)

	// Initialize handlers
	procurementHandler := handlers.NewProcurementHandler(procurementService)

	// Initialize JWT manager
	jwtManager := utils.NewJWTManager(cfg.JWTSecret, cfg.JWTExpiry, cfg.RefreshExpiry)

	// Setup Gin router
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()
	router.Use(middleware.CORSMiddleware())

	// Rate limiting
	redisClient := redis.NewClient(&redis.Options{
		Addr: cfg.RedisAddr,
	})
	rateLimiter := middleware.NewRateLimiter(redisClient, 100, time.Minute)
	router.Use(rateLimiter.Middleware())

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "healthy",
			"service": "procurement-service",
			"time":    time.Now(),
		})
	})

	// API routes
	api := router.Group("/api/v1")
	procurementHandler.RegisterRoutes(api, jwtManager)

	// Initialize RPC Handler
	if rabbitClient != nil {
		rpcHandler := handlers.NewRPCHandler(procurementService)
		_, err := rabbitClient.DeclareQueue("procurement.dashboard.stats")
		if err != nil {
			log.Fatalf("Failed to declare RPC queue: %v", err)
		}
		if err := rabbitClient.RPCServe("procurement.dashboard.stats", rpcHandler.HandleDashboardStats); err != nil {
			log.Fatalf("Failed to start RPC server: %v", err)
		}
	}

	// Start server
	port := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("Procurement Service starting on port %s", cfg.Port)
	if err := router.Run(port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func createIndexes(db *mongo.Database) error {
	ctx := context.Background()

	// Supplier indexes
	supplierIndexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "organization_id", Value: 1}, {Key: "supplier_code", Value: 1}},
			Options: options.Index().SetUnique(true).SetName("idx_supplier_org_code"),
		},
		{
			Keys:    bson.D{{Key: "organization_id", Value: 1}, {Key: "status", Value: 1}},
			Options: options.Index().SetName("idx_supplier_org_status"),
		},
		{
			Keys:    bson.D{{Key: "email", Value: 1}},
			Options: options.Index().SetName("idx_supplier_email"),
		},
		{
			Keys:    bson.D{{Key: "deleted_at", Value: 1}},
			Options: options.Index().SetName("idx_supplier_deleted"),
		},
	}
	if _, err := db.Collection("suppliers").Indexes().CreateMany(ctx, supplierIndexes); err != nil {
		return fmt.Errorf("failed to create supplier indexes: %w", err)
	}

	// Purchase Order indexes
	poIndexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "organization_id", Value: 1}, {Key: "po_number", Value: 1}},
			Options: options.Index().SetUnique(true).SetName("idx_po_org_number"),
		},
		{
			Keys:    bson.D{{Key: "organization_id", Value: 1}, {Key: "status", Value: 1}},
			Options: options.Index().SetName("idx_po_org_status"),
		},
		{
			Keys:    bson.D{{Key: "supplier_id", Value: 1}, {Key: "order_date", Value: -1}},
			Options: options.Index().SetName("idx_po_supplier_date"),
		},
		{
			Keys:    bson.D{{Key: "deleted_at", Value: 1}},
			Options: options.Index().SetName("idx_po_deleted"),
		},
	}
	if _, err := db.Collection("purchase_orders").Indexes().CreateMany(ctx, poIndexes); err != nil {
		return fmt.Errorf("failed to create purchase order indexes: %w", err)
	}

	// GRN indexes
	grnIndexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "organization_id", Value: 1}, {Key: "grn_number", Value: 1}},
			Options: options.Index().SetUnique(true).SetName("idx_grn_org_number"),
		},
		{
			Keys:    bson.D{{Key: "purchase_order_id", Value: 1}},
			Options: options.Index().SetName("idx_grn_po"),
		},
		{
			Keys:    bson.D{{Key: "supplier_id", Value: 1}},
			Options: options.Index().SetName("idx_grn_supplier"),
		},
		{
			Keys:    bson.D{{Key: "status", Value: 1}},
			Options: options.Index().SetName("idx_grn_status"),
		},
		{
			Keys:    bson.D{{Key: "deleted_at", Value: 1}},
			Options: options.Index().SetName("idx_grn_deleted"),
		},
	}
	if _, err := db.Collection("grns").Indexes().CreateMany(ctx, grnIndexes); err != nil {
		return fmt.Errorf("failed to create GRN indexes: %w", err)
	}

	log.Println("Database indexes created successfully")
	return nil
}
