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

type DeviceHandler struct {
	deviceService *service.DeviceService
}

func NewDeviceHandler(deviceService *service.DeviceService) *DeviceHandler {
	return &DeviceHandler{
		deviceService: deviceService,
	}
}

// CreateDevice registers a new device for an organization
// @Summary Register a new device
// @Tags devices
// @Accept json
// @Produce json
// @Param request body service.CreateDeviceRequest true "Device details"
// @Success 201 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Router /org-devices [post]
func (h *DeviceHandler) CreateDevice(c *gin.Context) {
	var req service.CreateDeviceRequest
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

	device, err := h.deviceService.CreateDevice(c.Request.Context(), req, createdBy)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "CREATE_FAILED", err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, device, "Device registered successfully")
}

// GetDevice retrieves a device by ID
// @Summary Get device by ID
// @Tags devices
// @Produce json
// @Param id path string true "Device ID"
// @Success 200 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /org-devices/{id} [get]
func (h *DeviceHandler) GetDevice(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "Invalid device ID", nil)
		return
	}

	device, err := h.deviceService.GetDevice(c.Request.Context(), id)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "NOT_FOUND", err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, device, "Device retrieved successfully")
}

// ListDevices returns paginated devices for an organization
// @Summary List devices
// @Tags devices
// @Produce json
// @Param organization_id query string true "Organization ID"
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Param search query string false "Search term"
// @Param is_active query bool false "Active status"
// @Param device_type query string false "Device type"
// @Success 200 {object} utils.Response
// @Router /org-devices [get]
func (h *DeviceHandler) ListDevices(c *gin.Context) {
	orgID := c.Query("organization_id")
	if orgID == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", "Organization ID is required", nil)
		return
	}

	orgObjID, err := primitive.ObjectIDFromHex(orgID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "Invalid organization ID", nil)
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	filters := make(map[string]interface{})
	if search := c.Query("search"); search != "" {
		filters["search"] = search
	}
	if isActive := c.Query("is_active"); isActive != "" {
		filters["is_active"] = isActive == "true"
	}
	if deviceType := c.Query("device_type"); deviceType != "" {
		filters["device_type"] = deviceType
	}

	devices, total, err := h.deviceService.ListDevices(c.Request.Context(), orgObjID, filters, page, limit)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "FETCH_FAILED", err.Error(), nil)
		return
	}

	utils.PaginatedResponse(c, devices, page, limit, total)
}

// UpdateDevice updates a device
// @Summary Update device
// @Tags devices
// @Accept json
// @Produce json
// @Param id path string true "Device ID"
// @Param request body service.UpdateDeviceRequest true "Update data"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Router /org-devices/{id} [put]
func (h *DeviceHandler) UpdateDevice(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "Invalid device ID", nil)
		return
	}

	var req service.UpdateDeviceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error(), nil)
		return
	}

	userID := middleware.GetUserID(c)
	updatedBy, _ := primitive.ObjectIDFromHex(userID)

	device, err := h.deviceService.UpdateDevice(c.Request.Context(), id, req, updatedBy)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "UPDATE_FAILED", err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, device, "Device updated successfully")
}

// DeleteDevice soft deletes a device
// @Summary Delete device
// @Tags devices
// @Param id path string true "Device ID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Router /org-devices/{id} [delete]
func (h *DeviceHandler) DeleteDevice(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "Invalid device ID", nil)
		return
	}

	userID := middleware.GetUserID(c)
	deletedBy, _ := primitive.ObjectIDFromHex(userID)

	if err := h.deviceService.DeleteDevice(c.Request.Context(), id, deletedBy); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "DELETE_FAILED", err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, nil, "Device deleted successfully")
}

// LookupDeviceByMAC looks up a device by MAC address (public endpoint for registration)
// @Summary Lookup device by MAC address
// @Tags devices
// @Produce json
// @Param mac_address query string true "MAC address"
// @Success 200 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /org-devices/lookup [get]
func (h *DeviceHandler) LookupDeviceByMAC(c *gin.Context) {
	macAddress := c.Query("mac_address")
	if macAddress == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", "MAC address is required", nil)
		return
	}

	result, err := h.deviceService.LookupDeviceByMAC(c.Request.Context(), macAddress)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "DEVICE_NOT_FOUND", err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, result, "Device found")
}

// RegisterRoutes registers all device routes
func (h *DeviceHandler) RegisterRoutes(router *gin.RouterGroup, jwtManager *utils.JWTManager) {
	devices := router.Group("/org-devices")

	// Public endpoint - for mobile app to look up device org during registration
	devices.GET("/lookup", h.LookupDeviceByMAC)

	// Protected routes - admin only
	protected := devices.Group("")
	protected.Use(middleware.AuthMiddleware(jwtManager))
	{
		protected.POST("", h.CreateDevice)
		protected.GET("", h.ListDevices)
		protected.GET("/:id", h.GetDevice)
		protected.PUT("/:id", h.UpdateDevice)
		protected.DELETE("/:id", h.DeleteDevice)
	}
}
