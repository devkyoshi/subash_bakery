package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/yourusername/erp-system/services/inventory-service/internal/service"
	"github.com/yourusername/erp-system/shared/middleware"
	"github.com/yourusername/erp-system/shared/utils"
)

type BatchHandler struct {
	batchService *service.BatchService
}

func NewBatchHandler(batchService *service.BatchService) *BatchHandler {
	return &BatchHandler{
		batchService: batchService,
	}
}

func (h *BatchHandler) CreateBatch(c *gin.Context) {
	orgIDStr := c.Param("org_id")
	orgID, err := primitive.ObjectIDFromHex(orgIDStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "Invalid organization ID", nil)
		return
	}

	var req service.CreateBatchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error(), nil)
		return
	}

	userID := middleware.GetUserID(c)
	createdBy, _ := primitive.ObjectIDFromHex(userID)

	batch, err := h.batchService.CreateBatch(c.Request.Context(), orgID, req, createdBy)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "CREATE_FAILED", err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, batch, "Batch created successfully")
}

func (h *BatchHandler) GetBatch(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "Invalid batch ID", nil)
		return
	}

	batch, err := h.batchService.GetBatch(c.Request.Context(), id)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "NOT_FOUND", err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, batch, "Batch retrieved successfully")
}

func (h *BatchHandler) GetBatches(c *gin.Context) {
	productIDStr := c.Param("product_id")
	productID, err := primitive.ObjectIDFromHex(productIDStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "Invalid product ID", nil)
		return
	}

	locationIDStr := c.Query("location_id")
	if locationIDStr == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "MISSING_PARAMS", "location_id required", nil)
		return
	}

	locationID, _ := primitive.ObjectIDFromHex(locationIDStr)
	activeOnly := c.Query("active") == "true"

	batches, err := h.batchService.GetBatchesByProduct(c.Request.Context(), productID, locationID, activeOnly)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "FETCH_FAILED", err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, batches, "Batches retrieved successfully")
}

func (h *BatchHandler) UpdateBatchQuantity(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "Invalid batch ID", nil)
		return
	}

	var req struct {
		Delta float64 `json:"delta" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error(), nil)
		return
	}

	if err := h.batchService.UpdateBatchQuantity(c.Request.Context(), id, req.Delta); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "UPDATE_FAILED", err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, nil, "Batch quantity updated successfully")
}
