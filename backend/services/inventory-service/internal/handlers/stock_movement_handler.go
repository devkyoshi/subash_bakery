package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/yourusername/erp-system/services/inventory-service/internal/service"
	"github.com/yourusername/erp-system/shared/middleware"
	"github.com/yourusername/erp-system/shared/utils"
)

type StockMovementHandler struct {
	stockService *service.StockService
}

func NewStockMovementHandler(stockService *service.StockService) *StockMovementHandler {
	return &StockMovementHandler{
		stockService: stockService,
	}
}

func (h *StockMovementHandler) CreateStockMovement(c *gin.Context) {
	orgIDStr := c.Param("org_id")
	orgID, err := primitive.ObjectIDFromHex(orgIDStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "Invalid organization ID", nil)
		return
	}

	var req service.StockMovementRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error(), nil)
		return
	}

	userID := middleware.GetUserID(c)
	createdBy, _ := primitive.ObjectIDFromHex(userID)

	movement, err := h.stockService.CreateStockMovement(c.Request.Context(), orgID, req, createdBy)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "CREATE_FAILED", err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, movement, "Stock movement created successfully")
}

func (h *StockMovementHandler) GetStockMovement(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "Invalid movement ID", nil)
		return
	}

	movement, err := h.stockService.GetStockMovement(c.Request.Context(), id)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "NOT_FOUND", err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, movement, "Stock movement retrieved successfully")
}

func (h *StockMovementHandler) GetStockMovements(c *gin.Context) {
	productIDStr := c.Param("product_id")
	productID, err := primitive.ObjectIDFromHex(productIDStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "Invalid product ID", nil)
		return
	}

	page := utils.GetPageParam(c)
	limit := utils.GetLimitParam(c)

	movements, err := h.stockService.GetStockMovements(c.Request.Context(), productID, page, limit)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "FETCH_FAILED", err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, movements, "Stock movements retrieved successfully")
}

func (h *StockMovementHandler) ListStockMovements(c *gin.Context) {
	orgIDStr := c.Param("org_id")
	orgID, err := primitive.ObjectIDFromHex(orgIDStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "Invalid organization ID", nil)
		return
	}

	filters := make(map[string]interface{})
	if movementType := c.Query("movement_type"); movementType != "" {
		filters["movement_type"] = movementType
	}
	if locationID := c.Query("location_id"); locationID != "" {
		filters["location_id"] = locationID
	}
	if productID := c.Query("product_id"); productID != "" {
		filters["product_id"] = productID
	}

	page := utils.GetPageParam(c)
	limit := utils.GetLimitParam(c)

	movements, err := h.stockService.ListStockMovements(c.Request.Context(), orgID, filters, page, limit)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "LIST_FAILED", err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, movements, "Stock movements retrieved successfully")
}

func (h *StockMovementHandler) GetStockMovementsByLocation(c *gin.Context) {
	locationIDStr := c.Param("location_id")
	locationID, err := primitive.ObjectIDFromHex(locationIDStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "Invalid location ID", nil)
		return
	}

	movements, err := h.stockService.GetStockMovementsByLocation(c.Request.Context(), locationID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "FETCH_FAILED", err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, movements, "Stock movements retrieved successfully")
}
