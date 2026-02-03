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

	"github.com/yourusername/erp-system/services/inventory-service/config"
	"github.com/yourusername/erp-system/services/inventory-service/internal/client"
	"github.com/yourusername/erp-system/services/inventory-service/internal/handlers"
	"github.com/yourusername/erp-system/services/inventory-service/internal/repository"
	"github.com/yourusername/erp-system/services/inventory-service/internal/service"
	"github.com/yourusername/erp-system/shared/middleware"
	"github.com/yourusername/erp-system/shared/utils"
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
	db := mongoClient.Database("erp_inventory")

	// Create indexes
	if err := createIndexes(db); err != nil {
		log.Fatalf("Failed to create indexes: %v", err)
	}

	// Initialize repositories
	stockLevelRepo := repository.NewStockLevelRepository(db)
	stockMovementRepo := repository.NewStockMovementRepository(db)
	batchRepo := repository.NewBatchRepository(db)
	adjustmentRepo := repository.NewStockAdjustmentRepository(db)
	countRepo := repository.NewInventoryCountRepository(db)
	serialRepo := repository.NewSerialNumberRepository(db)
	unitRepo := repository.NewUnitRepository(db)
	unitChartRepo := repository.NewUnitChartRepository(db)

	// Initialize services
	stockLevelService := service.NewStockLevelService(stockLevelRepo)
	stockService := service.NewStockService(stockLevelRepo, stockMovementRepo, batchRepo)
	batchService := service.NewBatchService(batchRepo, stockLevelRepo)
	adjustmentService := service.NewStockAdjustmentService(adjustmentRepo, stockLevelRepo, stockMovementRepo)
	countService := service.NewInventoryCountService(countRepo, stockLevelRepo, adjustmentRepo)
	serialNumberService := service.NewSerialNumberService(serialRepo, stockLevelRepo)
	unitService := service.NewUnitService(unitRepo, unitChartRepo)
	unitChartService := service.NewUnitChartService(unitChartRepo, unitRepo)

	// Initialize clients
	productClient := client.NewProductClient(cfg)
	orgClient := client.NewOrgClient(cfg)

	// Initialize handlers
	inventoryHandler := handlers.NewInventoryHandler(
		stockLevelService,
		stockService,
		batchService,
		adjustmentService,
		countService,
		serialNumberService,
		productClient,
		orgClient,
	)
	unitHandler := handlers.NewUnitHandler(unitService)
	unitChartHandler := handlers.NewUnitChartHandler(unitChartService)

	// Initialize JWT manager
	jwtManager := utils.NewJWTManager(cfg.JWTSecret, cfg.JWTExpiry, cfg.RefreshExpiry)

	// Setup Gin router
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()
	router.Use(middleware.CORSMiddleware())

	// TODO: Implement rate limiting
	// router.Use(middleware.RateLimitMiddleware())

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "healthy",
			"service": "inventory-service",
			"time":    time.Now(),
		})
	})

	// API routes
	api := router.Group("/api/v1")
	inventoryHandler.RegisterRoutes(api, jwtManager)
	unitHandler.RegisterRoutes(api, jwtManager)
	unitChartHandler.RegisterRoutes(api, jwtManager)

	// Start server
	port := fmt.Sprintf("0.0.0.0:%s", cfg.Port)
	log.Printf("Inventory Service starting on port %s", cfg.Port)
	if err := router.Run(port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func createIndexes(db *mongo.Database) error {
	ctx := context.Background()

	// Stock levels indexes
	stockIndexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "organization_id", Value: 1}, {Key: "product_id", Value: 1}, {Key: "location_id", Value: 1}},
			Options: options.Index().SetUnique(true).SetName("idx_stock_org_product_location"),
		},
		{
			Keys:    bson.D{{Key: "product_id", Value: 1}, {Key: "location_id", Value: 1}},
			Options: options.Index().SetName("idx_stock_product_location"),
		},
		{
			Keys:    bson.D{{Key: "location_id", Value: 1}},
			Options: options.Index().SetName("idx_stock_location"),
		},
	}
	if _, err := db.Collection("stock_levels").Indexes().CreateMany(ctx, stockIndexes); err != nil {
		return fmt.Errorf("failed to create stock indexes: %w", err)
	}

	// Stock movements indexes
	movementIndexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "organization_id", Value: 1}, {Key: "movement_date", Value: -1}},
			Options: options.Index().SetName("idx_movement_org_date"),
		},
		{
			Keys:    bson.D{{Key: "product_id", Value: 1}, {Key: "movement_date", Value: -1}},
			Options: options.Index().SetName("idx_movement_product_date"),
		},
		{
			Keys:    bson.D{{Key: "reference_type", Value: 1}, {Key: "reference_id", Value: 1}},
			Options: options.Index().SetName("idx_movement_reference"),
		},
	}
	if _, err := db.Collection("stock_movements").Indexes().CreateMany(ctx, movementIndexes); err != nil {
		return fmt.Errorf("failed to create movement indexes: %w", err)
	}

	// Batches indexes
	batchIndexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "organization_id", Value: 1}, {Key: "batch_number", Value: 1}},
			Options: options.Index().SetUnique(true).SetName("idx_batch_org_number"),
		},
		{
			Keys:    bson.D{{Key: "product_id", Value: 1}, {Key: "is_active", Value: 1}},
			Options: options.Index().SetName("idx_batch_product_active"),
		},
		{
			Keys:    bson.D{{Key: "expiry_date", Value: 1}},
			Options: options.Index().SetName("idx_batch_expiry"),
		},
	}
	if _, err := db.Collection("batches").Indexes().CreateMany(ctx, batchIndexes); err != nil {
		return fmt.Errorf("failed to create batch indexes: %w", err)
	}

	// Stock adjustments indexes
	adjustmentIndexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "organization_id", Value: 1}, {Key: "adjustment_date", Value: -1}},
			Options: options.Index().SetName("idx_adjustment_org_date"),
		},
		{
			Keys:    bson.D{{Key: "location_id", Value: 1}, {Key: "status", Value: 1}},
			Options: options.Index().SetName("idx_adjustment_location_status"),
		},
		{
			Keys:    bson.D{{Key: "adjustment_no", Value: 1}},
			Options: options.Index().SetName("idx_adjustment_no"),
		},
	}
	if _, err := db.Collection("stock_adjustments").Indexes().CreateMany(ctx, adjustmentIndexes); err != nil {
		return fmt.Errorf("failed to create adjustment indexes: %w", err)
	}

	// Inventory counts indexes
	countIndexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "organization_id", Value: 1}, {Key: "count_date", Value: -1}},
			Options: options.Index().SetName("idx_count_org_date"),
		},
		{
			Keys:    bson.D{{Key: "location_id", Value: 1}, {Key: "status", Value: 1}},
			Options: options.Index().SetName("idx_count_location_status"),
		},
		{
			Keys:    bson.D{{Key: "count_no", Value: 1}},
			Options: options.Index().SetName("idx_count_no"),
		},
	}
	if _, err := db.Collection("inventory_counts").Indexes().CreateMany(ctx, countIndexes); err != nil {
		return fmt.Errorf("failed to create count indexes: %w", err)
	}

	// Serial numbers indexes
	serialIndexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "organization_id", Value: 1}, {Key: "serial_no", Value: 1}},
			Options: options.Index().SetUnique(true).SetName("idx_serial_org_no"),
		},
		{
			Keys:    bson.D{{Key: "product_id", Value: 1}, {Key: "is_available", Value: 1}},
			Options: options.Index().SetName("idx_serial_product_available"),
		},
		{
			Keys:    bson.D{{Key: "status", Value: 1}},
			Options: options.Index().SetName("idx_serial_status"),
		},
	}
	if _, err := db.Collection("serial_numbers").Indexes().CreateMany(ctx, serialIndexes); err != nil {
		return fmt.Errorf("failed to create serial number indexes: %w", err)
	}

	// Units indexes
	unitIndexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "code", Value: 1}},
			Options: options.Index().SetUnique(true).SetName("idx_unit_code"),
		},
		{
			Keys:    bson.D{{Key: "unit_type", Value: 1}, {Key: "is_active", Value: 1}},
			Options: options.Index().SetName("idx_unit_type_active"),
		},
		{
			Keys:    bson.D{{Key: "is_base_unit", Value: 1}, {Key: "unit_type", Value: 1}},
			Options: options.Index().SetName("idx_unit_base_type"),
		},
	}
	if _, err := db.Collection("units").Indexes().CreateMany(ctx, unitIndexes); err != nil {
		return fmt.Errorf("failed to create unit indexes: %w", err)
	}

	// Unit charts indexes
	unitChartIndexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "from_unit_id", Value: 1}, {Key: "to_unit_id", Value: 1}},
			Options: options.Index().SetUnique(true).SetName("idx_unit_chart_from_to"),
		},
		{
			Keys:    bson.D{{Key: "from_unit_id", Value: 1}, {Key: "is_active", Value: 1}},
			Options: options.Index().SetName("idx_unit_chart_from_active"),
		},
		{
			Keys:    bson.D{{Key: "to_unit_id", Value: 1}, {Key: "is_active", Value: 1}},
			Options: options.Index().SetName("idx_unit_chart_to_active"),
		},
	}
	if _, err := db.Collection("unit_charts").Indexes().CreateMany(ctx, unitChartIndexes); err != nil {
		return fmt.Errorf("failed to create unit chart indexes: %w", err)
	}

	log.Println("Database indexes created successfully")
	return nil
}
