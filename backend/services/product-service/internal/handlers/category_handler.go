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

type CategoryHandler struct {
	categoryService *service.CategoryService
}

func NewCategoryHandler(categoryService *service.CategoryService) *CategoryHandler {
	return &CategoryHandler{
		categoryService: categoryService,
	}
}

// Route registration
func (h *CategoryHandler) RegisterCategoryRoutes(router *gin.RouterGroup, jwtManager *utils.JWTManager) {
	protected := router.Group("")
	protected.Use(middleware.AuthMiddleware(jwtManager))

	categories := protected.Group("/categories")
	{
		categories.POST("", h.CreateCategory)
		categories.GET("", h.ListCategories)
		categories.GET("/tree", h.GetCategoryTree)
		categories.GET("/root", h.GetRootCategories)
		categories.GET("/:id", h.GetCategory)
		categories.GET("/:id/children", h.GetChildren)
		categories.PUT("/:id", h.UpdateCategory)
		categories.DELETE("/:id", h.DeleteCategory)
	}
}

// CreateCategory creates a new category
// @Summary Create a new category
// @Description Creates a new product category with optional parent (for subcategories). Supports creating multiple subcategories in a single request.
// @Tags Categories
// @Accept json
// @Produce json
// @Param category body models.ProductCategory true "Category data"
// @Success 201 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /categories [post]
func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	var req service.CreateCategoryRequest
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

	category, err := h.categoryService.CreateCategory(c.Request.Context(), req, orgID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "CREATE_FAILED", err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, category, "Category created successfully")
}

// GetCategory retrieves a category by ID
// @Summary Get a category
// @Description Retrieves a category by its ID
// @Tags Categories
// @Produce json
// @Param id path string true "Category ID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /categories/{id} [get]
func (h *CategoryHandler) GetCategory(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "Invalid category ID", nil)
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

	category, err := h.categoryService.GetCategory(c.Request.Context(), id, orgID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "NOT_FOUND", err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, category, "Category retrieved successfully")
}

// ListCategories retrieves categories with optional filters
// @Summary List categories
// @Description Retrieves categories with optional filters for parent, level, and active status
// @Tags Categories
// @Produce json
// @Param organization_id query string true "Organization ID"
// @Param parent_id query string false "Parent category ID"
// @Param level query int false "Category level"
// @Param is_active query bool false "Active status"
// @Param q query string false "Search query for category name"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Router /categories [get]
func (h *CategoryHandler) ListCategories(c *gin.Context) {
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

	// Parse optional filters
	var isActive *bool
	if isActiveParam := c.Query("is_active"); isActiveParam != "" {
		active, err := strconv.ParseBool(isActiveParam)
		if err != nil {
			utils.ErrorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid is_active", nil)
			return
		}
		isActive = &active
	}

	query := c.Query("q")

	// Pagination
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	categories, total, err := h.categoryService.ListCategories(c.Request.Context(), orgID, isActive, query, page, limit)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "LIST_FAILED", err.Error(), nil)
		return
	}

	utils.PaginatedResponse(c, categories, page, limit, total)
}

// GetCategoryTree retrieves the full category tree
// @Summary Get category tree
// @Description Retrieves the full hierarchical category tree for an organization
// @Tags Categories
// @Produce json
// @Param organization_id query string true "Organization ID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Router /categories/tree [get]
func (h *CategoryHandler) GetCategoryTree(c *gin.Context) {
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

	tree, err := h.categoryService.GetCategoryTree(c.Request.Context(), orgID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "FETCH_FAILED", err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, tree, "Category tree retrieved successfully")
}

// GetRootCategories retrieves all root-level categories
// @Summary Get root categories
// @Description Retrieves all root-level (level 0) categories for an organization
// @Tags Categories
// @Produce json
// @Param organization_id query string true "Organization ID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Router /categories/root [get]
func (h *CategoryHandler) GetRootCategories(c *gin.Context) {
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

	categories, err := h.categoryService.GetRootCategories(c.Request.Context(), orgID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "FETCH_FAILED", err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, categories, "Root categories retrieved successfully")
}

// GetChildren retrieves direct children of a category
// @Summary Get category children
// @Description Retrieves direct children (subcategories) of a category
// @Tags Categories
// @Produce json
// @Param id path string true "Parent category ID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Router /categories/{id}/children [get]
func (h *CategoryHandler) GetChildren(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "Invalid category ID", nil)
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

	children, err := h.categoryService.GetChildren(c.Request.Context(), id, orgID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "FETCH_FAILED", err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, children, "Children categories retrieved successfully")
}

// UpdateCategory updates an existing category
// @Summary Update a category
// @Description Updates an existing category's information
// @Tags Categories
// @Accept json
// @Produce json
// @Param id path string true "Category ID"
// @Param category body models.ProductCategory true "Category data"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /categories/{id} [put]
func (h *CategoryHandler) UpdateCategory(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "Invalid category ID", nil)
		return
	}

	var req service.UpdateCategoryRequest
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

	category, err := h.categoryService.UpdateCategory(c.Request.Context(), id, req, orgID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "UPDATE_FAILED", err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, category, "Category updated successfully")
}

// DeleteCategory deletes a category
// @Summary Delete a category
// @Description Soft deletes a category (only if it has no subcategories or products)
// @Tags Categories
// @Produce json
// @Param id path string true "Category ID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /categories/{id} [delete]
func (h *CategoryHandler) DeleteCategory(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "Invalid category ID", nil)
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

	if err := h.categoryService.DeleteCategory(c.Request.Context(), id, orgID); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "DELETE_FAILED", err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, nil, "Category deleted successfully")
}
