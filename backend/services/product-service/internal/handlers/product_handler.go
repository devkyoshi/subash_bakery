package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/erp-system/services/product-service/internal/service"
	"github.com/yourusername/erp-system/shared/middleware"
	"github.com/yourusername/erp-system/shared/models"
	"github.com/yourusername/erp-system/shared/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProductHandler struct {
	productService *service.ProductService
}

func NewProductHandler(productService *service.ProductService) *ProductHandler {
	return &ProductHandler{
		productService: productService,
	}
}

// Route registration
func (h *ProductHandler) RegisterProductRoutes(
	router *gin.RouterGroup,
	jwtManager *utils.JWTManager,
) {
	protected := router.Group("")
	protected.Use(middleware.AuthMiddleware(jwtManager))

	products := protected.Group("/products")
	{
		products.POST("", h.CreateProduct)
		products.GET("", h.ListProducts)
		products.GET("/low-stock", h.GetLowStockProducts)
		products.GET("/sku/:sku", h.GetProductBySKU)
		products.GET("/:id", h.GetProduct)
		products.PUT("/:id", h.UpdateProduct)
		products.DELETE("/:id", h.DeleteProduct)
	}
}

// CreateProduct creates a new product
// @Summary Create a new product
// @Description Creates a new product in the inventory
// @Tags Products
// @Accept json
// @Produce json
// @Param product body models.Product true "Product data"
// @Success 201 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /products [post]
func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var req service.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error(), nil)
		return
	}

	// Get user organization from context
	userOrgIDStr, exists := c.Get("organization_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "INVALID_USER", "Organization not found in context", nil)
		return
	}

	orgID, err := primitive.ObjectIDFromHex(userOrgIDStr.(string))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "INVALID_ORG_ID", "Invalid organization ID", nil)
		return
	}

	token := c.GetHeader("Authorization")
	product, err := h.productService.CreateProduct(c.Request.Context(), req, orgID, token)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "CREATE_FAILED", err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, product, "Product created successfully")
}

// GetProduct retrieves a product by ID
// @Summary Get a product
// @Description Retrieves a product by its ID with optional location price filtering
// @Tags Products
// @Produce json
// @Param id path string true "Product ID"
// @Param location_id query string false "Location ID - filter to show only this location's price"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /products/{id} [get]
func (h *ProductHandler) GetProduct(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "Invalid product ID", nil)
		return
	}

	userOrgIDStr, exists := c.Get("organization_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "INVALID_USER", "Organization not found in context", nil)
		return
	}

	orgID, err := primitive.ObjectIDFromHex(userOrgIDStr.(string))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "INVALID_ORG_ID", "Invalid organization ID", nil)
		return
	}

	product, err := h.productService.GetProduct(c.Request.Context(), id, orgID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "NOT_FOUND", err.Error(), nil)
		return
	}

	// Filter location prices if location_id is provided
	if locationIDParam := c.Query("location_id"); locationIDParam != "" {
		locationID, err := primitive.ObjectIDFromHex(locationIDParam)
		if err != nil {
			utils.ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "Invalid location_id", nil)
			return
		}

		// Filter the location prices
		filteredPrices := []models.LocationPrice{}
		for _, lp := range product.LocationPrices {
			if lp.LocationID == locationID {
				filteredPrices = append(filteredPrices, lp)
			}
		}
		product.LocationPrices = filteredPrices
	}

	utils.SuccessResponse(c, http.StatusOK, product, "Product retrieved successfully")
}

// GetProductBySKU retrieves a product by SKU
// @Summary Get product by SKU
// @Description Retrieves a product by its SKU
// @Tags Products
// @Produce json
// @Param organization_id query string true "Organization ID"
// @Param sku path string true "Product SKU"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /products/sku/{sku} [get]
func (h *ProductHandler) GetProductBySKU(c *gin.Context) {
	sku := c.Param("sku")
	if sku == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", "SKU is required", nil)
		return
	}

	orgIDParam := c.Query("organization_id")
	if orgIDParam == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", "organization_id is required", nil)
		return
	}

	orgID, err := primitive.ObjectIDFromHex(orgIDParam)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "Invalid organization ID", nil)
		return
	}

	userOrgIDStr, exists := c.Get("organization_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "INVALID_USER", "Organization not found in context", nil)
		return
	}

	userOrg, err := primitive.ObjectIDFromHex(userOrgIDStr.(string))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "INVALID_ORG_ID", "Invalid organization ID", nil)
		return
	}

	product, err := h.productService.GetProductBySKU(c.Request.Context(), orgID, sku, userOrg)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "NOT_FOUND", err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, product, "Product retrieved successfully")
}

// ListProducts retrieves products with filters
// @Summary List products
// @Description Retrieves products with optional filters
// @Tags Products
// @Produce json
// @Param organization_id query string true "Organization ID"
// @Param category_id query string false "Category ID"
// @Param brand_id query string false "Brand ID"
// @Param status query string false "Product status"
// @Param type query string false "Product type"
// @Param track_inventory query bool false "Track inventory"
// @Param location_id query string false "Location ID - filter products by location price"
// @Param search query string false "Search by name, SKU, or description"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Router /products [get]
func (h *ProductHandler) ListProducts(c *gin.Context) {
	orgIDParam := c.Query("organization_id")
	if orgIDParam == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", "organization_id is required", nil)
		return
	}

	orgID, err := primitive.ObjectIDFromHex(orgIDParam)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "Invalid organization ID", nil)
		return
	}

	// Build filters
	filters := make(map[string]interface{})

	if categoryIDParam := c.Query("category_id"); categoryIDParam != "" {
		categoryID, err := primitive.ObjectIDFromHex(categoryIDParam)
		if err != nil {
			utils.ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "Invalid category_id", nil)
			return
		}
		filters["category_id"] = categoryID
	}

	if subcategoryIDParam := c.Query("subcategory_id"); subcategoryIDParam != "" {
		subcategoryID, err := primitive.ObjectIDFromHex(subcategoryIDParam)
		if err != nil {
			utils.ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "Invalid subcategory_id", nil)
			return
		}
		filters["subcategory_id"] = subcategoryID
	}

	if brandIDParam := c.Query("brand_id"); brandIDParam != "" {
		brandID, err := primitive.ObjectIDFromHex(brandIDParam)
		if err != nil {
			utils.ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "Invalid brand_id", nil)
			return
		}
		filters["brand_id"] = brandID
	}

	if status := c.Query("status"); status != "" {
		filters["status"] = models.ProductStatus(status)
	}

	if productType := c.Query("type"); productType != "" {
		filters["type"] = models.ProductType(productType)
	}

	if trackInventoryParam := c.Query("track_inventory"); trackInventoryParam != "" {
		trackInventory, err := strconv.ParseBool(trackInventoryParam)
		if err != nil {
			utils.ErrorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid track_inventory", nil)
			return
		}
		filters["track_inventory"] = trackInventory
	}

	if search := c.Query("search"); search != "" {
		filters["search"] = search
	}

	if locationIDParam := c.Query("location_id"); locationIDParam != "" {
		locationID, err := primitive.ObjectIDFromHex(locationIDParam)
		if err != nil {
			utils.ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "Invalid location_id", nil)
			return
		}
		filters["location_id"] = locationID
	}

	// Pagination
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	products, total, err := h.productService.ListProducts(c.Request.Context(), orgID, filters, page, limit)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "LIST_FAILED", err.Error(), nil)
		return
	}

	utils.PaginatedResponse(c, products, page, limit, total)
}

// GetLowStockProducts retrieves products below reorder level
// @Summary Get low stock products
// @Description Retrieves products that are below their reorder level
// @Tags Products
// @Produce json
// @Param organization_id query string true "Organization ID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Router /products/low-stock [get]
func (h *ProductHandler) GetLowStockProducts(c *gin.Context) {
	orgIDParam := c.Query("organization_id")
	if orgIDParam == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", "organization_id is required", nil)
		return
	}

	orgID, err := primitive.ObjectIDFromHex(orgIDParam)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "Invalid organization ID", nil)
		return
	}

	products, err := h.productService.GetLowStockProducts(c.Request.Context(), orgID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "FETCH_FAILED", err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, products, "Low stock products retrieved successfully")
}

// UpdateProduct updates an existing product
// @Summary Update a product
// @Description Updates an existing product's information
// @Tags Products
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Param product body models.Product true "Product data"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /products/{id} [put]
func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "Invalid product ID", nil)
		return
	}

	var req service.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error(), nil)
		return
	}

	userOrgIDStr, exists := c.Get("organization_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "INVALID_USER", "Organization not found in context", nil)
		return
	}

	orgID, err := primitive.ObjectIDFromHex(userOrgIDStr.(string))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "INVALID_ORG_ID", "Invalid organization ID", nil)
		return
	}

	product, err := h.productService.UpdateProduct(c.Request.Context(), id, req, orgID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "UPDATE_FAILED", err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, product, "Product updated successfully")
}

// DeleteProduct deletes a product
// @Summary Delete a product
// @Description Soft deletes a product
// @Tags Products
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /products/{id} [delete]
func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "Invalid product ID", nil)
		return
	}

	userOrgIDStr, exists := c.Get("organization_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "INVALID_USER", "Organization not found in context", nil)
		return
	}

	orgID, err := primitive.ObjectIDFromHex(userOrgIDStr.(string))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "INVALID_ORG_ID", "Invalid organization ID", nil)
		return
	}

	if err := h.productService.DeleteProduct(c.Request.Context(), id, orgID); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "DELETE_FAILED", err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, nil, "Product deleted successfully")
}
