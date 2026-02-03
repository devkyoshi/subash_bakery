package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/erp-system/services/product-service/internal/service"
	"github.com/yourusername/erp-system/shared/middleware"
	"github.com/yourusername/erp-system/shared/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BrandHandler struct {
	brandService *service.BrandService
}

func NewBrandHandler(brandService *service.BrandService) *BrandHandler {
	return &BrandHandler{
		brandService: brandService,
	}
}

// Route registration
func (h *BrandHandler) RegisterBrandRoutes(router *gin.RouterGroup, jwtManager *utils.JWTManager) {
	protected := router.Group("")
	protected.Use(middleware.AuthMiddleware(jwtManager))

	brands := protected.Group("/brands")
	{
		brands.POST("", h.CreateBrand)
		brands.GET("", h.ListBrands)
		brands.GET("/search", h.SearchBrands)
		brands.GET("/:id", h.GetBrand)
		brands.PUT("/:id", h.UpdateBrand)
		brands.DELETE("/:id", h.DeleteBrand)
	}
}

// CreateBrand creates a new brand
// @Summary Create a new brand
// @Tags brands
// @Accept json
// @Produce json
// @Param request body service.CreateBrandRequest true "Brand details"
// @Success 201 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Router /brands [post]
func (h *BrandHandler) CreateBrand(c *gin.Context) {
	var req service.CreateBrandRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error(), nil)
		return
	}

	userID := middleware.GetUserID(c)
	createdBy, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, "INVALID_USER", "Invalid user ID", nil)
		return
	}

	brand, err := h.brandService.CreateBrand(c.Request.Context(), req, createdBy)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "CREATE_FAILED", err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, brand, "Brand created successfully")
}

// GetBrand retrieves a brand by ID
// @Summary Get a brand by ID
// @Tags brands
// @Produce json
// @Param id path string true "Brand ID"
// @Success 200 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /brands/{id} [get]
func (h *BrandHandler) GetBrand(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "Invalid brand ID", nil)
		return
	}

	brand, err := h.brandService.GetBrand(c.Request.Context(), id)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "NOT_FOUND", err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, brand, "Brand retrieved successfully")
}

// ListBrands retrieves all brands for an organization
// @Summary List brands
// @Tags brands
// @Produce json
// @Param org_id query string true "Organization ID"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Param is_active query bool false "Filter by active status"
// @Param q query string false "Search query (optional)"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Router /brands [get]
func (h *BrandHandler) ListBrands(c *gin.Context) {
	orgID, err := primitive.ObjectIDFromHex(c.Query("org_id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "Invalid organization ID", nil)
		return
	}

	page := utils.GetPageParam(c)
	limit := utils.GetLimitParam(c)

	filter := service.BrandFilter{
		Page:  page,
		Limit: limit,
		Query: c.Query("q"), // Optional search query
	}

	// Handle is_active filter - optional, only filter when provided
	if val, exists := c.GetQuery("is_active"); exists {
		parsed, err := strconv.ParseBool(val)
		if err != nil {
			utils.ErrorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid is_active value", nil)
			return
		}
		filter.IsActive = &parsed
	}

	brands, total, err := h.brandService.GetBrandsByOrganization(c.Request.Context(), orgID, filter)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "LIST_FAILED", err.Error(), nil)
		return
	}

	response := gin.H{
		"brands": brands,
		"pagination": gin.H{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	}

	utils.SuccessResponse(c, http.StatusOK, response, "Brands retrieved successfully")
}

// UpdateBrand updates an existing brand
// @Summary Update a brand
// @Tags brands
// @Accept json
// @Produce json
// @Param id path string true "Brand ID"
// @Param request body service.UpdateBrandRequest true "Brand update details"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Router /brands/{id} [put]
func (h *BrandHandler) UpdateBrand(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "Invalid brand ID", nil)
		return
	}

	var req service.UpdateBrandRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error(), nil)
		return
	}

	userID := middleware.GetUserID(c)
	updatedBy, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, "INVALID_USER", "Invalid user ID", nil)
		return
	}

	brand, err := h.brandService.UpdateBrand(c.Request.Context(), id, req, updatedBy)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "UPDATE_FAILED", err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, brand, "Brand updated successfully")
}

// DeleteBrand deletes a brand
// @Summary Delete a brand
// @Tags brands
// @Produce json
// @Param id path string true "Brand ID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Router /brands/{id} [delete]
func (h *BrandHandler) DeleteBrand(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "Invalid brand ID", nil)
		return
	}

	userID := middleware.GetUserID(c)
	deletedBy, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, "INVALID_USER", "Invalid user ID", nil)
		return
	}

	if err := h.brandService.DeleteBrand(c.Request.Context(), id, deletedBy); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "DELETE_FAILED", err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, nil, "Brand deleted successfully")
}

// SearchBrands searches brands by name or code
// @Summary Search brands
// @Tags brands
// @Produce json
// @Param org_id query string true "Organization ID"
// @Param q query string true "Search query"
// @Param is_active query bool false "Filter by active status"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Router /brands/search [get]
func (h *BrandHandler) SearchBrands(c *gin.Context) {
	orgID, err := primitive.ObjectIDFromHex(c.Query("org_id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "Invalid organization ID", nil)
		return
	}

	query := c.Query("q")
	if query == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", "Search query is required", nil)
		return
	}

	page := utils.GetPageParam(c)
	limit := utils.GetLimitParam(c)
	var isActive *bool

	if val, exists := c.GetQuery("is_active"); exists {
		parsed, err := strconv.ParseBool(val)
		if err != nil {
			utils.ErrorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid is_active value", nil)
			return
		}
		isActive = &parsed
	}

	brands, total, err := h.brandService.SearchBrands(c.Request.Context(), orgID, query, isActive, page, limit)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "SEARCH_FAILED", err.Error(), nil)
		return
	}

	response := gin.H{
		"brands": brands,
		"pagination": gin.H{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	}

	utils.SuccessResponse(c, http.StatusOK, response, "Search completed successfully")
}
