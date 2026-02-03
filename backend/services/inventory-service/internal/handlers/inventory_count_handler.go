package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/yourusername/erp-system/services/inventory-service/internal/service"
	"github.com/yourusername/erp-system/shared/middleware"
	"github.com/yourusername/erp-system/shared/utils"
)

type InventoryCountHandler struct {
	countService *service.InventoryCountService
}

func NewInventoryCountHandler(countService *service.InventoryCountService) *InventoryCountHandler {
	return &InventoryCountHandler{
		countService: countService,
	}
}

func (h *InventoryCountHandler) CreateInventoryCount(c *gin.Context) {
	orgIDStr := c.Param("org_id")
	orgID, err := primitive.ObjectIDFromHex(orgIDStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "Invalid organization ID", nil)
		return
	}

	var req service.CreateInventoryCountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error(), nil)
		return
	}

	userID := middleware.GetUserID(c)
	createdBy, _ := primitive.ObjectIDFromHex(userID)

	count, err := h.countService.CreateCount(c.Request.Context(), orgID, req, createdBy)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "CREATE_FAILED", err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, count, "Inventory count created successfully")
}

func (h *InventoryCountHandler) GetInventoryCount(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "Invalid count ID", nil)
		return
	}

	count, err := h.countService.GetCount(c.Request.Context(), id)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "NOT_FOUND", err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, count, "Count retrieved successfully")
}

func (h *InventoryCountHandler) ListInventoryCounts(c *gin.Context) {
	orgIDStr := c.Param("org_id")
	orgID, err := primitive.ObjectIDFromHex(orgIDStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "Invalid organization ID", nil)
		return
	}

	filters := make(map[string]interface{})
	if status := c.Query("status"); status != "" {
		filters["status"] = status
	}
	if locationID := c.Query("location_id"); locationID != "" {
		filters["location_id"] = locationID
	}
	if countType := c.Query("count_type"); countType != "" {
		filters["count_type"] = countType
	}

	page := utils.GetPageParam(c)
	limit := utils.GetLimitParam(c)

	counts, err := h.countService.ListCounts(c.Request.Context(), orgID, filters, page, limit)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "LIST_FAILED", err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, counts, "Counts retrieved successfully")
}

func (h *InventoryCountHandler) UpdateCountItem(c *gin.Context) {
	countID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "Invalid count ID", nil)
		return
	}

	var req service.UpdateCountItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error(), nil)
		return
	}

	userID := middleware.GetUserID(c)
	countedBy, _ := primitive.ObjectIDFromHex(userID)

	if err := h.countService.UpdateCountItem(c.Request.Context(), countID, req, countedBy); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "UPDATE_FAILED", err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, nil, "Count item updated successfully")
}

func (h *InventoryCountHandler) CompleteInventoryCount(c *gin.Context) {
	countID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "Invalid count ID", nil)
		return
	}

	var req struct {
		CreateAdjustment bool `json:"create_adjustment"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		req.CreateAdjustment = false
	}

	userID := middleware.GetUserID(c)
	completedBy, _ := primitive.ObjectIDFromHex(userID)

	if err := h.countService.CompleteCount(c.Request.Context(), countID, completedBy, req.CreateAdjustment); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "COMPLETE_FAILED", err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, nil, "Count completed successfully")
}

func (h *InventoryCountHandler) CancelInventoryCount(c *gin.Context) {
	countID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "Invalid count ID", nil)
		return
	}

	if err := h.countService.CancelCount(c.Request.Context(), countID); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "CANCEL_FAILED", err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, nil, "Count cancelled successfully")
}

func (h *InventoryCountHandler) DeleteInventoryCount(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "Invalid count ID", nil)
		return
	}

	userID := middleware.GetUserID(c)
	deletedBy, _ := primitive.ObjectIDFromHex(userID)

	if err := h.countService.DeleteCount(c.Request.Context(), id, deletedBy); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "DELETE_FAILED", err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, nil, "Count deleted successfully")
}
