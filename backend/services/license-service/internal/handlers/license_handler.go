package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/yourusername/erp-system/services/license-service/internal/service"
	"github.com/yourusername/erp-system/shared/middleware"
	"github.com/yourusername/erp-system/shared/utils"
)

type LicenseHandler struct {
	licenseService *service.LicenseService
}

func NewLicenseHandler(licenseService *service.LicenseService) *LicenseHandler {
	return &LicenseHandler{
		licenseService: licenseService,
	}
}

// CreateLicense creates a new license
func (h *LicenseHandler) CreateLicense(c *gin.Context) {
	orgID, err := primitive.ObjectIDFromHex(c.Param("org_id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "Invalid organization ID", nil)
		return
	}

	var req service.CreateLicenseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error(), nil)
		return
	}

	userID := middleware.GetUserID(c)
	createdBy, _ := primitive.ObjectIDFromHex(userID)

	license, err := h.licenseService.CreateLicense(c.Request.Context(), orgID, req, createdBy)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "CREATE_FAILED", err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, license, "License created successfully")
}

// GetLicense retrieves a license by ID
func (h *LicenseHandler) GetLicense(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "Invalid license ID", nil)
		return
	}

	license, err := h.licenseService.GetLicense(c.Request.Context(), id)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "NOT_FOUND", err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, license, "License retrieved successfully")
}

// ListLicensesByLocation returns all licenses for a location
func (h *LicenseHandler) ListLicensesByLocation(c *gin.Context) {
	locationID, err := primitive.ObjectIDFromHex(c.Param("location_id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "Invalid location ID", nil)
		return
	}

	activeOnly := c.Query("active") == "true"

	licenses, err := h.licenseService.ListLicensesByLocation(c.Request.Context(), locationID, activeOnly)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "LIST_FAILED", err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, licenses, "Licenses retrieved successfully")
}

// ListLicensesByOrganization returns all licenses for an organization
func (h *LicenseHandler) ListLicensesByOrganization(c *gin.Context) {
	orgID, err := primitive.ObjectIDFromHex(c.Param("org_id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "Invalid organization ID", nil)
		return
	}

	page := utils.GetPageParam(c)
	limit := utils.GetLimitParam(c)

	licenses, err := h.licenseService.ListLicensesByOrganization(c.Request.Context(), orgID, page, limit)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "LIST_FAILED", err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, licenses, "Licenses retrieved successfully")
}

// UpdateLicenseUsage updates license usage
func (h *LicenseHandler) UpdateLicenseUsage(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "Invalid license ID", nil)
		return
	}

	var req service.UpdateUsageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error(), nil)
		return
	}

	if err := h.licenseService.UpdateLicenseUsage(c.Request.Context(), id, req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "UPDATE_FAILED", err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, nil, "License usage updated successfully")
}

// SuspendLicense suspends a license
func (h *LicenseHandler) SuspendLicense(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "Invalid license ID", nil)
		return
	}

	var req struct {
		Reason string `json:"reason"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error(), nil)
		return
	}

	if err := h.licenseService.SuspendLicense(c.Request.Context(), id, req.Reason); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "SUSPEND_FAILED", err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, nil, "License suspended successfully")
}

// ActivateLicense activates a license
func (h *LicenseHandler) ActivateLicense(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "Invalid license ID", nil)
		return
	}

	if err := h.licenseService.ActivateLicense(c.Request.Context(), id); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ACTIVATE_FAILED", err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, nil, "License activated successfully")
}

// RevokeLicense revokes a license
func (h *LicenseHandler) RevokeLicense(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "Invalid license ID", nil)
		return
	}

	var req struct {
		Reason string `json:"reason"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error(), nil)
		return
	}

	userID := middleware.GetUserID(c)
	revokedBy, _ := primitive.ObjectIDFromHex(userID)

	if err := h.licenseService.RevokeLicense(c.Request.Context(), id, revokedBy, req.Reason); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "REVOKE_FAILED", err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, nil, "License revoked successfully")
}

// AssignUser assigns a user to a license
func (h *LicenseHandler) AssignUser(c *gin.Context) {
	licenseID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "Invalid license ID", nil)
		return
	}

	var req service.AssignUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error(), nil)
		return
	}

	userID := middleware.GetUserID(c)
	assignedBy, _ := primitive.ObjectIDFromHex(userID)

	if err := h.licenseService.AssignUserToLicense(c.Request.Context(), licenseID, req, assignedBy); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ASSIGN_FAILED", err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, nil, "User assigned to license successfully")
}

// RevokeUser revokes a user from a license
func (h *LicenseHandler) RevokeUser(c *gin.Context) {
	assignmentID, err := primitive.ObjectIDFromHex(c.Param("assignment_id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "Invalid assignment ID", nil)
		return
	}

	var req struct {
		Reason string `json:"reason"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error(), nil)
		return
	}

	userID := middleware.GetUserID(c)
	revokedBy, _ := primitive.ObjectIDFromHex(userID)

	if err := h.licenseService.RevokeUserFromLicense(c.Request.Context(), assignmentID, revokedBy, req.Reason); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "REVOKE_FAILED", err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, nil, "User revoked from license successfully")
}

// ActivateDevice activates a device for a license
func (h *LicenseHandler) ActivateDevice(c *gin.Context) {
	licenseID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "Invalid license ID", nil)
		return
	}

	var req service.ActivateDeviceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error(), nil)
		return
	}

	userID := middleware.GetUserID(c)
	activatedBy, _ := primitive.ObjectIDFromHex(userID)

	if err := h.licenseService.ActivateDevice(c.Request.Context(), licenseID, req, activatedBy); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ACTIVATE_FAILED", err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, nil, "Device activated successfully")
}

// DeactivateDevice deactivates a device
func (h *LicenseHandler) DeactivateDevice(c *gin.Context) {
	assignmentID, err := primitive.ObjectIDFromHex(c.Param("assignment_id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "Invalid assignment ID", nil)
		return
	}

	var req struct {
		Reason string `json:"reason"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error(), nil)
		return
	}

	userID := middleware.GetUserID(c)
	deactivatedBy, _ := primitive.ObjectIDFromHex(userID)

	if err := h.licenseService.DeactivateDevice(c.Request.Context(), assignmentID, deactivatedBy, req.Reason); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "DEACTIVATE_FAILED", err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, nil, "Device deactivated successfully")
}

// RegisterRoutes registers all license routes
func (h *LicenseHandler) RegisterRoutes(router *gin.RouterGroup, appHandler *ApplicationHandler, jwtManager *utils.JWTManager) {
	// Public routes (no auth required) - list available applications
	public := router.Group("/applications")
	{
		public.GET("", appHandler.ListApplications)
		public.GET("/:id", appHandler.GetApplication)
		public.GET("/code/:code", appHandler.GetApplicationByCode)
	}

	// Protected routes
	protected := router.Group("")
	protected.Use(middleware.AuthMiddleware(jwtManager))

	// Applications (admin only in production)
	apps := protected.Group("/applications")
	{
		apps.POST("", appHandler.CreateApplication)
		apps.PUT("/:id", appHandler.UpdateApplication)
		apps.DELETE("/:id", appHandler.DeleteApplication)
	}

	// Licenses
	licenses := protected.Group("/organizations/:org_id/licenses")
	{
		licenses.POST("", h.CreateLicense)
		licenses.GET("", h.ListLicensesByOrganization)
	}

	locations := protected.Group("/locations/:location_id/licenses")
	{
		locations.GET("", h.ListLicensesByLocation)
	}

	license := protected.Group("/licenses/:id")
	{
		license.GET("", h.GetLicense)
		license.PUT("/usage", h.UpdateLicenseUsage)
		license.POST("/suspend", h.SuspendLicense)
		license.POST("/activate", h.ActivateLicense)
		license.POST("/revoke", h.RevokeLicense)
		
		// User assignments
		license.POST("/users", h.AssignUser)
		license.DELETE("/users/:assignment_id", h.RevokeUser)
		
		// Device assignments
		license.POST("/devices", h.ActivateDevice)
		license.DELETE("/devices/:assignment_id", h.DeactivateDevice)
	}
}
