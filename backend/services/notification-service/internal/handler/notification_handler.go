package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/erp-system/services/notification-service/internal/service"
	"github.com/yourusername/erp-system/shared/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type NotificationHandler struct {
	notifService *service.NotificationService
}

func NewNotificationHandler(notifService *service.NotificationService) *NotificationHandler {
	return &NotificationHandler{
		notifService: notifService,
	}
}

func (h *NotificationHandler) RegisterRoutes(router *gin.RouterGroup) {
	notifications := router.Group("/notifications")
	{
		notifications.GET("", h.GetNotifications)
		notifications.PATCH("/:id/read", h.MarkAsRead)
		notifications.PATCH("/read-all", h.MarkAllAsRead)
	}
}

func (h *NotificationHandler) GetNotifications(c *gin.Context) {
	userIDStr := c.GetString("user_id")
	if userIDStr == "" {
		utils.ErrorResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated", nil)
		return
	}
	userID, _ := primitive.ObjectIDFromHex(userIDStr)

	notifications, err := h.notifService.GetNotifications(c.Request.Context(), userID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to fetch notifications", nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, notifications, "Notifications fetched successfully")
}

func (h *NotificationHandler) MarkAsRead(c *gin.Context) {
	userIDStr := c.GetString("user_id")
	if userIDStr == "" {
		utils.ErrorResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated", nil)
		return
	}
	userID, _ := primitive.ObjectIDFromHex(userIDStr)

	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "Invalid notification ID", nil)
		return
	}

	if err := h.notifService.MarkAsRead(c.Request.Context(), id, userID); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to mark notification as read", nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, nil, "Notification marked as read")
}

func (h *NotificationHandler) MarkAllAsRead(c *gin.Context) {
	userIDStr := c.GetString("user_id")
	if userIDStr == "" {
		utils.ErrorResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated", nil)
		return
	}
	userID, _ := primitive.ObjectIDFromHex(userIDStr)

	if err := h.notifService.MarkAllAsRead(c.Request.Context(), userID); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to mark all notifications as read", nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, nil, "All notifications marked as read")
}
