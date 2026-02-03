package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/yourusername/erp-system/services/inventory-service/internal/service"
	"github.com/yourusername/erp-system/shared/utils"
)

type StockLevelHandler struct {
	stockLevelService *service.StockLevelService
}

func NewStockLevelHandler(stockLevelService *service.StockLevelService) *StockLevelHandler {
	return &StockLevelHandler{
		stockLevelService: stockLevelService,
	}
}

func (h *StockLevelHandler) GetStockLevel(c *gin.Context) {
	productIDStr := c.Query("product_id")
	locationIDStr := c.Query("location_id")

	if productIDStr == "" || locationIDStr == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "MISSING_PARAMS", "product_id and location_id required", nil)
		return
	}

	productID, _ := primitive.ObjectIDFromHex(productIDStr)
	locationID, _ := primitive.ObjectIDFromHex(locationIDStr)

	stock, err := h.stockLevelService.GetStockLevel(c.Request.Context(), productID, locationID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "NOT_FOUND", err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, stock, "Stock level retrieved successfully")
}

func (h *StockLevelHandler) GetStockLevelsBatch(c *gin.Context) {
	var req struct {
		ProductIDs []string `json:"product_ids" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error(), nil)
		return
	}

	var productIDs []primitive.ObjectID
	for _, idStr := range req.ProductIDs {
		id, err := primitive.ObjectIDFromHex(idStr)
		if err == nil {
			productIDs = append(productIDs, id)
		}
	}

	if len(productIDs) == 0 {
		utils.SuccessResponse(c, http.StatusOK, map[string]interface{}{}, "No valid product IDs provided")
		return
	}

	stockMap, err := h.stockLevelService.GetStockByProductIDs(c.Request.Context(), productIDs)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "FETCH_FAILED", err.Error(), nil)
		return
	}

	// Transform map keys to just match what we need (maybe return as is, or list)
	// Returning the map: "productID_locationID": StockLevel

	utils.SuccessResponse(c, http.StatusOK, stockMap, "Stock levels retrieved successfully")
}

func (h *StockLevelHandler) GetStockByProduct(c *gin.Context) {
	productIDStr := c.Param("product_id")
	productID, err := primitive.ObjectIDFromHex(productIDStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "Invalid product ID", nil)
		return
	}

	stocks, err := h.stockLevelService.GetStockByProduct(c.Request.Context(), productID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "FETCH_FAILED", err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, stocks, "Stock retrieved successfully")
}

func (h *StockLevelHandler) GetStockByLocation(c *gin.Context) {
	locationIDStr := c.Param("location_id")
	locationID, err := primitive.ObjectIDFromHex(locationIDStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "Invalid location ID", nil)
		return
	}

	stocks, err := h.stockLevelService.GetStockByLocation(c.Request.Context(), locationID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "FETCH_FAILED", err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, stocks, "Stock retrieved successfully")
}

func (h *StockLevelHandler) AllocateStock(c *gin.Context) {
	var req struct {
		ProductID  string  `json:"product_id" binding:"required"`
		LocationID string  `json:"location_id" binding:"required"`
		Quantity   float64 `json:"quantity" binding:"required,gt=0"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error(), nil)
		return
	}

	productID, _ := primitive.ObjectIDFromHex(req.ProductID)
	locationID, _ := primitive.ObjectIDFromHex(req.LocationID)

	if err := h.stockLevelService.AllocateStock(c.Request.Context(), productID, locationID, req.Quantity); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ALLOCATE_FAILED", err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, nil, "Stock allocated successfully")
}

func (h *StockLevelHandler) ReleaseStock(c *gin.Context) {
	var req struct {
		ProductID  string  `json:"product_id" binding:"required"`
		LocationID string  `json:"location_id" binding:"required"`
		Quantity   float64 `json:"quantity" binding:"required,gt=0"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error(), nil)
		return
	}

	productID, _ := primitive.ObjectIDFromHex(req.ProductID)
	locationID, _ := primitive.ObjectIDFromHex(req.LocationID)

	if err := h.stockLevelService.ReleaseStock(c.Request.Context(), productID, locationID, req.Quantity); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "RELEASE_FAILED", err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, nil, "Stock released successfully")
}
