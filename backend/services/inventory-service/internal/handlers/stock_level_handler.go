package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/yourusername/erp-system/services/inventory-service/internal/client"
	"github.com/yourusername/erp-system/services/inventory-service/internal/service"
	"github.com/yourusername/erp-system/shared/models"
	"github.com/yourusername/erp-system/shared/utils"
)

type StockLevelHandler struct {
	stockLevelService *service.StockLevelService
	productClient     *client.ProductClient
	orgClient         *client.OrgClient
}

func NewStockLevelHandler(stockLevelService *service.StockLevelService, productClient *client.ProductClient, orgClient *client.OrgClient) *StockLevelHandler {
	return &StockLevelHandler{
		stockLevelService: stockLevelService,
		productClient:     productClient,
		orgClient:         orgClient,
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

func (h *StockLevelHandler) ListStockLevels(c *gin.Context) {
	filters := make(map[string]interface{})

	if orgIDStr := c.Query("organization_id"); orgIDStr != "" {
		if orgID, err := primitive.ObjectIDFromHex(orgIDStr); err == nil {
			filters["organization_id"] = orgID
		}
	}
	if productIDStr := c.Query("product_id"); productIDStr != "" {
		if productID, err := primitive.ObjectIDFromHex(productIDStr); err == nil {
			filters["product_id"] = productID
		}
	}
	if locationIDStr := c.Query("location_id"); locationIDStr != "" {
		if locationID, err := primitive.ObjectIDFromHex(locationIDStr); err == nil {
			filters["location_id"] = locationID
		}
	}

	page := utils.GetPageParam(c)
	limit := utils.GetLimitParam(c)

	stocks, err := h.stockLevelService.ListStockLevels(c.Request.Context(), filters, page, limit)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "FETCH_FAILED", err.Error(), nil)
		return
	}

	// Enrich data with Product and Location details
	if len(stocks) > 0 {
		token := c.GetHeader("Authorization")

		// Collect unique IDs
		productIDs := make([]primitive.ObjectID, 0)
		seenProducts := make(map[string]bool)

		var orgID primitive.ObjectID
		if len(stocks) > 0 {
			orgID = stocks[0].OrganizationID
		}

		for _, s := range stocks {
			if !seenProducts[s.ProductID.Hex()] {
				productIDs = append(productIDs, s.ProductID)
				seenProducts[s.ProductID.Hex()] = true
			}
		}

		// Fetch details concurrently
		// Only fetch products if we have IDs
		if len(productIDs) > 0 {
			productsMap, err := h.productClient.GetProductsBatch(c.Request.Context(), productIDs, token)
			if err == nil {
				// Populate Product details
				for i := range stocks {
					if p, ok := productsMap[stocks[i].ProductID.Hex()]; ok {
						stocks[i].ProductName = p.Name
						stocks[i].SKU = p.SKU
					}
				}
			} else {
				fmt.Printf("Error fetching products: %v\n", err)
			}
		}

		// Fetch locations for the organization
		locationsMap, err := h.orgClient.GetLocationsByOrganization(c.Request.Context(), orgID, token)
		if err == nil {
			// Populate Location details
			for i := range stocks {
				if l, ok := locationsMap[stocks[i].LocationID.Hex()]; ok {
					stocks[i].LocationName = l.Name
				}
			}
		} else {
			fmt.Printf("Error fetching locations: %v\n", err)
		}
	}

	//TODO: Get total count of stock levels
	response := map[string]interface{}{
		"data":  stocks,
		"page":  page,
		"limit": limit,
		"total": len(stocks),
	}

	utils.SuccessResponse(c, http.StatusOK, response, "Stock levels retrieved successfully")
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

func (h *StockLevelHandler) GetDashboardStats(c *gin.Context) {
	orgIDStr := c.Query("organization_id")
	if orgIDStr == "" {
		// Try header
		orgIDStr = c.GetHeader("x-organization-id")
	}

	if orgIDStr == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "MISSING_ORG_ID", "Organization ID required", nil)
		return
	}

	orgID, err := primitive.ObjectIDFromHex(orgIDStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "INVALID_ORG_ID", "Invalid Organization ID", nil)
		return
	}

	stats, err := h.stockLevelService.GetDashboardStats(c.Request.Context(), orgID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "FETCH_FAILED", err.Error(), nil)
		return
	}

	// Populate product names for low stock items
	lowStockItems := stats["low_stock_items"].([]*models.StockLevel)
	if len(lowStockItems) > 0 {
		token := c.GetHeader("Authorization")
		productIDs := make([]primitive.ObjectID, len(lowStockItems))
		for i, item := range lowStockItems {
			productIDs[i] = item.ProductID
		}

		productsMap, err := h.productClient.GetProductsBatch(c.Request.Context(), productIDs, token)
		if err == nil {
			for _, item := range lowStockItems {
				if p, ok := productsMap[item.ProductID.Hex()]; ok {
					item.ProductName = p.Name
					item.SKU = p.SKU
				}
			}
		}
	}

	utils.SuccessResponse(c, http.StatusOK, stats, "Dashboard stats retrieved successfully")
}
