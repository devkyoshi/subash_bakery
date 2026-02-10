package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/erp-system/services/notification-service/internal/models"
	"github.com/yourusername/erp-system/services/notification-service/internal/repository"
	"github.com/yourusername/erp-system/shared/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DeviceHandler struct {
	deviceRepo *repository.DeviceRepository
}

func NewDeviceHandler(deviceRepo *repository.DeviceRepository) *DeviceHandler {
	return &DeviceHandler{
		deviceRepo: deviceRepo,
	}
}

func (h *DeviceHandler) RegisterRoutes(router *gin.RouterGroup) {
	devices := router.Group("/devices")
	{
		devices.POST("", h.RegisterDevice)
	}
}

type RegisterDeviceRequest struct {
	Token    string `json:"token" binding:"required"`
	Platform string `json:"platform" binding:"required"`
	Name     string `json:"name"`
}

func (h *DeviceHandler) RegisterDevice(c *gin.Context) {
	var req RegisterDeviceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "INVALID_PAYLOAD", "Invalid request payload", nil)
		return
	}

	// Get user ID and Org ID from context (set by Auth middleware)
	userIDStr := c.GetString("user_id")
	orgIDStr := c.GetString("organization_id")

	// If context is empty (e.g. testing without auth), handle gracefully
	if userIDStr == "" || orgIDStr == "" {
		utils.ErrorResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated", nil)
		return
	}

	userID, _ := primitive.ObjectIDFromHex(userIDStr)
	orgID, _ := primitive.ObjectIDFromHex(orgIDStr)

	device := &models.DeviceToken{
		UserID:         userID,
		OrganizationID: orgID,
		Token:          req.Token,
		Platform:       req.Platform,
		Name:           req.Name,
	}

	if err := h.deviceRepo.Register(c.Request.Context(), device); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to register device", nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, nil, "Device registered successfully")
}
