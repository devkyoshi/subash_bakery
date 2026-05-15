package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/erp-system/services/org-service/internal/service"
	"github.com/yourusername/erp-system/shared/middleware"
	"github.com/yourusername/erp-system/shared/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OrganizationHandler struct {
	orgService *service.OrganizationService
}

func NewOrganizationHandler(orgService *service.OrganizationService) *OrganizationHandler {
	return &OrganizationHandler{
		orgService: orgService,
	}
}

// CreateOrganization creates a new organization
func (h *OrganizationHandler) CreateOrganization(c *gin.Context) {
	var req service.CreateOrganizationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error(), nil)
		return
	}

	userID := middleware.GetUserID(c)
	if userID == "" {
		utils.ErrorResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated", nil)
		return
	}

	createdBy, _ := primitive.ObjectIDFromHex(userID)

	org, err := h.orgService.CreateOrganization(c.Request.Context(), req, createdBy)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "CREATE_FAILED", err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, org, "Organization created successfully")
}

// GetOrganization retrieves an organization by ID
func (h *OrganizationHandler) GetOrganization(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "Invalid organization ID", nil)
		return
	}

	org, err := h.orgService.GetOrganization(c.Request.Context(), id)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "NOT_FOUND", err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, org, "Organization retrieved successfully")
}

// ListOrganizations returns paginated organizations
func (h *OrganizationHandler) ListOrganizations(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	status := c.Query("status")

	orgs, total, err := h.orgService.ListOrganizations(c.Request.Context(), page, limit, status)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "LIST_FAILED", err.Error(), nil)
		return
	}

	utils.PaginatedResponse(c, orgs, page, limit, total)
}

// GetOrganizationOptions returns a list of organization options
func (h *OrganizationHandler) GetOrganizationOptions(c *gin.Context) {
	options, err := h.orgService.GetOrganizationOptions(c.Request.Context())
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "FETCH_FAILED", err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, options, "Organization options retrieved successfully")
}

// UpdateOrganization updates an organization
func (h *OrganizationHandler) UpdateOrganization(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "Invalid organization ID", nil)
		return
	}

	var req service.UpdateOrganizationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error(), nil)
		return
	}

	userID := middleware.GetUserID(c)
	updatedBy, _ := primitive.ObjectIDFromHex(userID)

	org, err := h.orgService.UpdateOrganization(c.Request.Context(), id, req, updatedBy)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "UPDATE_FAILED", err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, org, "Organization updated successfully")
}

// DeleteOrganization soft deletes an organization
func (h *OrganizationHandler) DeleteOrganization(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "Invalid organization ID", nil)
		return
	}

	userID := middleware.GetUserID(c)
	deletedBy, _ := primitive.ObjectIDFromHex(userID)

	if err := h.orgService.DeleteOrganization(c.Request.Context(), id, deletedBy); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "DELETE_FAILED", err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, nil, "Organization deleted successfully")
}

// CreateCompany creates a new company
func (h *OrganizationHandler) CreateCompany(c *gin.Context) {
	orgID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "Invalid organization ID", nil)
		return
	}

	var req service.CreateCompanyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error(), nil)
		return
	}

	userID := middleware.GetUserID(c)
	createdBy, _ := primitive.ObjectIDFromHex(userID)

	company, err := h.orgService.CreateCompany(c.Request.Context(), orgID, req, createdBy)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "CREATE_FAILED", err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, company, "Company created successfully")
}

// GetCompany retrieves a company by ID
func (h *OrganizationHandler) GetCompany(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "Invalid company ID", nil)
		return
	}

	company, err := h.orgService.GetCompany(c.Request.Context(), id)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "NOT_FOUND", err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, company, "Company retrieved successfully")
}

// ListCompanies returns paginated companies
func (h *OrganizationHandler) ListCompanies(c *gin.Context) {
	orgID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "Invalid organization ID", nil)
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	companies, total, err := h.orgService.ListCompanies(c.Request.Context(), orgID, page, limit)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "LIST_FAILED", err.Error(), nil)
		return
	}

	utils.PaginatedResponse(c, companies, page, limit, total)
}

// UpdateCompany updates a company
func (h *OrganizationHandler) UpdateCompany(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "Invalid company ID", nil)
		return
	}

	var req service.UpdateCompanyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error(), nil)
		return
	}

	userID := middleware.GetUserID(c)
	updatedBy, _ := primitive.ObjectIDFromHex(userID)

	company, err := h.orgService.UpdateCompany(c.Request.Context(), id, req, updatedBy)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "UPDATE_FAILED", err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, company, "Company updated successfully")
}

// DeleteCompany soft deletes a company
func (h *OrganizationHandler) DeleteCompany(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "Invalid company ID", nil)
		return
	}

	userID := middleware.GetUserID(c)
	deletedBy, _ := primitive.ObjectIDFromHex(userID)

	if err := h.orgService.DeleteCompany(c.Request.Context(), id, deletedBy); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "DELETE_FAILED", err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, nil, "Company deleted successfully")
}

// AssignUserToCompany assigns a user to a company
func (h *OrganizationHandler) AssignUserToCompany(c *gin.Context) {
	companyID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "Invalid company ID", nil)
		return
	}

	var req struct {
		UserID string `json:"user_id" binding:"required"`
		RoleID string `json:"role_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error(), nil)
		return
	}

	userID, err := primitive.ObjectIDFromHex(req.UserID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "Invalid user ID", nil)
		return
	}

	roleID, err := primitive.ObjectIDFromHex(req.RoleID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "Invalid role ID", nil)
		return
	}

	if err := h.orgService.AssignUserToCompany(c.Request.Context(), companyID, userID, roleID); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ASSIGNMENT_FAILED", err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, nil, "User assigned to company successfully")
}

// CreateLocation creates a new location
func (h *OrganizationHandler) CreateLocation(c *gin.Context) {
	companyID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "Invalid company ID", nil)
		return
	}

	var req service.CreateLocationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error(), nil)
		return
	}

	userID := middleware.GetUserID(c)
	createdBy, _ := primitive.ObjectIDFromHex(userID)

	location, err := h.orgService.CreateLocation(c.Request.Context(), companyID, req, createdBy)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "CREATE_FAILED", err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, location, "Location created successfully")
}

// GetLocation retrieves a location by ID
func (h *OrganizationHandler) GetLocation(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "Invalid location ID", nil)
		return
	}

	location, err := h.orgService.GetLocation(c.Request.Context(), id)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "NOT_FOUND", err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, location, "Location retrieved successfully")
}

// ListLocations returns paginated locations
func (h *OrganizationHandler) ListLocations(c *gin.Context) {
	companyID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "Invalid company ID", nil)
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	locations, total, err := h.orgService.ListLocations(c.Request.Context(), companyID, page, limit)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "LIST_FAILED", err.Error(), nil)
		return
	}

	utils.PaginatedResponse(c, locations, page, limit, total)
}

// ListAllLocations returns paginated locations for an organization
func (h *OrganizationHandler) ListAllLocations(c *gin.Context) {
	orgID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "Invalid organization ID", nil)
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	locations, total, err := h.orgService.ListLocationsByOrganization(c.Request.Context(), orgID, page, limit)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "LIST_FAILED", err.Error(), nil)
		return
	}

	utils.PaginatedResponse(c, locations, page, limit, total)
}

// GetUserAccess retrieves all organizations, companies, and locations accessible by the logged-in user
func (h *OrganizationHandler) GetUserAccess(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == "" {
		utils.ErrorResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated", nil)
		return
	}

	orgID := middleware.GetOrganizationID(c)
	if orgID == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "INVALID_REQUEST", "Organization ID not found", nil)
		return
	}

	userObjID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "Invalid user ID", nil)
		return
	}

	orgObjID, err := primitive.ObjectIDFromHex(orgID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "Invalid organization ID", nil)
		return
	}

	accessData, err := h.orgService.GetUserAccess(c.Request.Context(), userObjID, orgObjID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "FETCH_FAILED", err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, accessData, "User access data retrieved successfully")
}

// RegisterRoutes registers all organization routes
func (h *OrganizationHandler) RegisterRoutes(router *gin.RouterGroup, jwtManager *utils.JWTManager) {
	// All routes require authentication
	protected := router.Group("")
	protected.Use(middleware.AuthMiddleware(jwtManager))

	// User access endpoint
	protected.GET("/users/me/access", h.GetUserAccess)

	// Organizations
	orgs := protected.Group("/organizations")
	{
		orgs.POST("", h.CreateOrganization)
		orgs.GET("", h.ListOrganizations)
		orgs.GET("/options", h.GetOrganizationOptions) // Register BEFORE /:id
		orgs.GET("/:id", h.GetOrganization)
		orgs.PUT("/:id", h.UpdateOrganization)
		orgs.DELETE("/:id", h.DeleteOrganization)
		orgs.GET("/:id/locations", h.ListAllLocations)

		// Companies under organization
		orgs.POST("/:id/companies", h.CreateCompany)
		orgs.GET("/:id/companies", h.ListCompanies)
	}

	// Companies
	companies := protected.Group("/companies")
	{
		companies.GET("/:id", h.GetCompany)
		companies.PUT("/:id", h.UpdateCompany)
		companies.DELETE("/:id", h.DeleteCompany)
		companies.POST("/:id/users", h.AssignUserToCompany)

		// Locations under company
		companies.POST("/:id/locations", h.CreateLocation)
		companies.GET("/:id/locations", h.ListLocations)
	}

	// Locations
	locations := protected.Group("/locations")
	{
		locations.GET("/:id", h.GetLocation)
	}
}
