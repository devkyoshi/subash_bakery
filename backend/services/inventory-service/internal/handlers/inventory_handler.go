package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/yourusername/erp-system/services/inventory-service/internal/client"
	"github.com/yourusername/erp-system/services/inventory-service/internal/service"
	"github.com/yourusername/erp-system/shared/middleware"
	"github.com/yourusername/erp-system/shared/utils"
)

type InventoryHandler struct {
	stockLevelHandler      *StockLevelHandler
	stockMovementHandler   *StockMovementHandler
	batchHandler           *BatchHandler
	stockAdjustmentHandler *StockAdjustmentHandler
	inventoryCountHandler  *InventoryCountHandler
	serialNumberHandler    *SerialNumberHandler
}

func NewInventoryHandler(
	stockLevelService *service.StockLevelService,
	stockService *service.StockService,
	batchService *service.BatchService,
	adjustmentService *service.StockAdjustmentService,
	countService *service.InventoryCountService,
	serialNumberService *service.SerialNumberService,
	productClient *client.ProductClient,
	orgClient *client.OrgClient,
	userClient *client.UserClient,
) *InventoryHandler {
	return &InventoryHandler{
		stockLevelHandler:      NewStockLevelHandler(stockLevelService, productClient, orgClient),
		stockMovementHandler:   NewStockMovementHandler(stockService, productClient, orgClient),
		batchHandler:           NewBatchHandler(batchService),
		stockAdjustmentHandler: NewStockAdjustmentHandler(adjustmentService, productClient, orgClient, userClient),
		inventoryCountHandler:  NewInventoryCountHandler(countService),
		serialNumberHandler:    NewSerialNumberHandler(serialNumberService),
	}
}

// RegisterRoutes registers all inventory routes
func (h *InventoryHandler) RegisterRoutes(api *gin.RouterGroup, jwtManager *utils.JWTManager) {
	// Public routes
	api.GET("/health", func(c *gin.Context) {
		utils.SuccessResponse(c, http.StatusOK, map[string]string{"status": "healthy"}, "Inventory Service is running")
	})
	api.POST("/inventory/stock/bulk", h.stockLevelHandler.GetStockLevelsBatch)
	api.GET("/inventory/dashboard/stats", h.stockLevelHandler.GetDashboardStats)

	// Protected routes
	protected := api.Group("/inventory")
	protected.Use(middleware.AuthMiddleware(jwtManager))

	// Stock level routes
	protected.GET("/stock-levels", h.stockLevelHandler.ListStockLevels)
	protected.GET("/stock-levels/item", h.stockLevelHandler.GetStockLevel) // Re-map single item fetch if needed, or rely on List with filters
	protected.GET("/products/:product_id/stock", h.stockLevelHandler.GetStockByProduct)
	protected.GET("/locations/:location_id/stock", h.stockLevelHandler.GetStockByLocation)
	protected.POST("/stock/allocate", h.stockLevelHandler.AllocateStock)
	protected.POST("/stock/release", h.stockLevelHandler.ReleaseStock)

	// Stock movement routes
	protected.POST("/organizations/:org_id/stock-movements", h.stockMovementHandler.CreateStockMovement)
	protected.GET("/organizations/:org_id/stock-movements", h.stockMovementHandler.ListStockMovements)
	protected.GET("/stock-movements/:id", h.stockMovementHandler.GetStockMovement)
	protected.GET("/products/:product_id/movements", h.stockMovementHandler.GetStockMovements)
	protected.GET("/locations/:location_id/movements", h.stockMovementHandler.GetStockMovementsByLocation)

	// Batch routes
	protected.POST("/organizations/:org_id/batches", h.batchHandler.CreateBatch)
	protected.GET("/batches/:id", h.batchHandler.GetBatch)
	protected.PUT("/batches/:id/quantity", h.batchHandler.UpdateBatchQuantity)
	protected.GET("/products/:product_id/batches", h.batchHandler.GetBatches)

	// Stock adjustment routes
	protected.POST("/organizations/:org_id/stock-adjustments", h.stockAdjustmentHandler.CreateStockAdjustment)
	protected.GET("/organizations/:org_id/stock-adjustments", h.stockAdjustmentHandler.ListStockAdjustments)
	protected.GET("/stock-adjustments/:id", h.stockAdjustmentHandler.GetStockAdjustment)
	protected.PUT("/stock-adjustments/:id", h.stockAdjustmentHandler.UpdateStockAdjustment)
	protected.POST("/stock-adjustments/:id/approve", h.stockAdjustmentHandler.ApproveStockAdjustment)
	protected.POST("/stock-adjustments/:id/reject", h.stockAdjustmentHandler.RejectStockAdjustment)
	protected.DELETE("/stock-adjustments/:id", h.stockAdjustmentHandler.DeleteStockAdjustment)

	// Inventory count routes
	protected.POST("/organizations/:org_id/inventory-counts", h.inventoryCountHandler.CreateInventoryCount)
	protected.GET("/organizations/:org_id/inventory-counts", h.inventoryCountHandler.ListInventoryCounts)
	protected.GET("/inventory-counts/:id", h.inventoryCountHandler.GetInventoryCount)
	protected.POST("/inventory-counts/:id/items", h.inventoryCountHandler.UpdateCountItem)
	protected.POST("/inventory-counts/:id/complete", h.inventoryCountHandler.CompleteInventoryCount)
	protected.POST("/inventory-counts/:id/cancel", h.inventoryCountHandler.CancelInventoryCount)
	protected.DELETE("/inventory-counts/:id", h.inventoryCountHandler.DeleteInventoryCount)

	// Serial number routes
	protected.POST("/organizations/:org_id/serial-numbers", h.serialNumberHandler.CreateSerialNumber)
	protected.GET("/serial-numbers/:id", h.serialNumberHandler.GetSerialNumber)
	protected.GET("/organizations/:org_id/serial-numbers/:serial_no", h.serialNumberHandler.GetSerialNumberBySerial)
	protected.GET("/products/:product_id/serial-numbers", h.serialNumberHandler.ListSerialNumbers)
	protected.PUT("/serial-numbers/:id", h.serialNumberHandler.UpdateSerialNumber)
	protected.POST("/serial-numbers/:id/allocate", h.serialNumberHandler.AllocateSerialNumber)
	protected.POST("/serial-numbers/:id/sold", h.serialNumberHandler.MarkSerialAsSold)
	protected.DELETE("/serial-numbers/:id", h.serialNumberHandler.DeleteSerialNumber)
}
