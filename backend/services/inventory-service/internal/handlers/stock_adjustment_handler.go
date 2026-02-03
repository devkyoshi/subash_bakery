package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/yourusername/erp-system/services/inventory-service/internal/service"
	"github.com/yourusername/erp-system/shared/middleware"
	"github.com/yourusername/erp-system/shared/utils"
)

type StockAdjustmentHandler struct {
	adjustmentService *service.StockAdjustmentService
}

func NewStockAdjustmentHandler(adjustmentService *service.StockAdjustmentService) *StockAdjustmentHandler {
	return &StockAdjustmentHandler{
		adjustmentService: adjustmentService,
	}
}

func (h *StockAdjustmentHandler) CreateStockAdjustment(c *gin.Context) {
	orgIDStr := c.Param("org_id")
	orgID, err := primitive.ObjectIDFromHex(orgIDStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "Invalid organization ID", nil)
		return
	}

	var req service.CreateStockAdjustmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error(), nil)
		return
	}

	userID := middleware.GetUserID(c)
	createdBy, _ := primitive.ObjectIDFromHex(userID)

	adjustment, err := h.adjustmentService.CreateAdjustment(c.Request.Context(), orgID, req, createdBy)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "CREATE_FAILED", err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, adjustment, "Stock adjustment created successfully")
}

func (h *StockAdjustmentHandler) GetStockAdjustment(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "Invalid adjustment ID", nil)
		return
	}

	adjustment, err := h.adjustmentService.GetAdjustment(c.Request.Context(), id)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "NOT_FOUND", err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, adjustment, "Adjustment retrieved successfully")
}

func (h *StockAdjustmentHandler) ListStockAdjustments(c *gin.Context) {
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

	page := utils.GetPageParam(c)
	limit := utils.GetLimitParam(c)

	adjustments, err := h.adjustmentService.ListAdjustments(c.Request.Context(), orgID, filters, page, limit)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "LIST_FAILED", err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, adjustments, "Adjustments retrieved successfully")
}

func (h *StockAdjustmentHandler) UpdateStockAdjustment(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "Invalid adjustment ID", nil)
		return
	}

	var req service.CreateStockAdjustmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error(), nil)
		return
	}

	userID := middleware.GetUserID(c)
	updatedBy, _ := primitive.ObjectIDFromHex(userID)

	adjustment, err := h.adjustmentService.UpdateAdjustment(c.Request.Context(), id, req, updatedBy)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "UPDATE_FAILED", err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, adjustment, "Adjustment updated successfully")
}

func (h *StockAdjustmentHandler) ApproveStockAdjustment(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "Invalid adjustment ID", nil)
		return
	}

	userID := middleware.GetUserID(c)
	approvedBy, _ := primitive.ObjectIDFromHex(userID)

	if err := h.adjustmentService.ApproveAdjustment(c.Request.Context(), id, approvedBy); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "APPROVE_FAILED", err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, nil, "Adjustment approved successfully")
}

func (h *StockAdjustmentHandler) RejectStockAdjustment(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "Invalid adjustment ID", nil)
		return
	}

	var req struct {
		Reason string `json:"reason" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error(), nil)
		return
	}

	userID := middleware.GetUserID(c)
	rejectedBy, _ := primitive.ObjectIDFromHex(userID)

	if err := h.adjustmentService.RejectAdjustment(c.Request.Context(), id, rejectedBy, req.Reason); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "REJECT_FAILED", err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, nil, "Adjustment rejected successfully")
}

func (h *StockAdjustmentHandler) DeleteStockAdjustment(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "Invalid adjustment ID", nil)
		return
	}

	userID := middleware.GetUserID(c)
	deletedBy, _ := primitive.ObjectIDFromHex(userID)

	if err := h.adjustmentService.DeleteAdjustment(c.Request.Context(), id, deletedBy); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "DELETE_FAILED", err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, nil, "Adjustment deleted successfully")
}
