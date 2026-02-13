package main

import (
	"context"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/yourusername/erp-system/services/report-service/config"
	"github.com/yourusername/erp-system/services/report-service/internal/handler"
	"github.com/yourusername/erp-system/services/report-service/internal/repository"
	"github.com/yourusername/erp-system/services/report-service/internal/service"
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

	// Ping MongoDB to verify connection
	if err := mongoClient.Ping(ctx, nil); err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
	}
	log.Println("Connected to MongoDB successfully")

	// We read from the procurement database directly (read-only access)
	procurementDB := mongoClient.Database("erp_procurement")

	// Report service's own database for storing generated report metadata
	reportDB := mongoClient.Database("erp_reports")

	// Initialize repositories
	procurementReportRepo := repository.NewProcurementReportRepository(procurementDB)
	_ = repository.NewGeneratedReportRepository(reportDB)

	// Initialize services
	povsGRNService := service.NewPOvsGRNService(procurementReportRepo, procurementDB)
	exportService := service.NewExportService()

	// Initialize handlers
	reportHandler := handler.NewReportHandler(povsGRNService, exportService)

	// JWT Setup
	refreshExpiry, _ := time.ParseDuration(cfg.RefreshExpiry)
	jwtManager := utils.NewJWTManager(cfg.JWTSecret, 24*time.Hour, refreshExpiry)

	// Gin Router
	router := gin.Default()

	// CORS Middleware
	router.Use(middleware.CORSMiddleware())

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "healthy",
			"service": "report-service",
		})
	})

	// API Routes
	api := router.Group("/api/v1")
	reportHandler.RegisterRoutes(api, jwtManager)

	log.Printf("Report Service starting on port %s", cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}