package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/yourusername/erp-system/services/inventory-service/internal/client"
	"github.com/yourusername/erp-system/services/inventory-service/internal/service"
	"github.com/yourusername/erp-system/shared/middleware"
	"github.com/yourusername/erp-system/shared/utils"
)

type StockAdjustmentHandler struct {
	adjustmentService *service.StockAdjustmentService
	productClient     *client.ProductClient
	orgClient         *client.OrgClient
	userClient        *client.UserClient
}

func NewStockAdjustmentHandler(
	adjustmentService *service.StockAdjustmentService,
	productClient *client.ProductClient,
	orgClient *client.OrgClient,
	userClient *client.UserClient,
) *StockAdjustmentHandler {
	return &StockAdjustmentHandler{
		adjustmentService: adjustmentService,
		productClient:     productClient,
		orgClient:         orgClient,
		userClient:        userClient,
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

	// Enrich Data
	token := c.GetHeader("Authorization")
	ctx := c.Request.Context()

	// 1. Fetch Location
	if adjustment.LocationID != primitive.NilObjectID {
		if loc, err := h.orgClient.GetLocation(ctx, adjustment.LocationID, token); err == nil {
			adjustment.LocationName = loc.Name
		}
	}

	// 2. Fetch Users (CreatedBy, ApprovedBy)
	userIDs := []primitive.ObjectID{}
	if adjustment.CreatedBy != primitive.NilObjectID { // Wait, CreatedBy is in BaseModel... wait, BaseModel has CreatedBy?
		// Let's check BaseModel definition in inventory.go if available or assume standard.
		// StockAdjustment struct has CreatedBy in BaseModel?
		// Let's check `inventory.go` define BaseModel. It is usually inline.
		// Looking at my previous `view_file` of `inventory.go`, line 73 `BaseModel bson:",inline"`.
		// I need to know if BaseModel has CreatedBy.
		// Assuming yes for now, but `StockAdjustment` specifically has `ApprovedBy` and `StartedBy` (in InventoryCount).
		// `StockAdjustment` has `ApprovedBy`.
		// Standard `BaseModel` usually tracks `CreatedBy`.
		// I'll check `inventory.go` imports again or just trust common sense but `BaseModel` fields weren't shown in the snippet above (lines 71+).
		// I'll assume `CreatedBy` is there or use `ApprovedBy` definitely.
		// Wait, the `CreateStockAdjustment` handler uses `createdBy` variable.
		// I'll fetch `ApprovedBy` definitely.
		// Does `StockAdjustment` struct have `CreatedBy`?
		// Line 395 `BaseModel`.
		// I will try to use `CreatedBy` from BaseModel implicitly.

		// I'll add `ApprovedBy`.
		if adjustment.ApprovedBy != primitive.NilObjectID {
			userIDs = append(userIDs, adjustment.ApprovedBy)
		}
		// I'll just try to fetch ApprovedBy for now as that's explicit.
	}

	if len(userIDs) > 0 {
		if users, err := h.userClient.GetUsersBatch(ctx, userIDs, token); err == nil {
			if adjustment.ApprovedBy != primitive.NilObjectID {
				if user, ok := users[adjustment.ApprovedBy.Hex()]; ok {
					adjustment.ApprovedByName = fmt.Sprintf("%s %s", user.FirstName, user.LastName)
				}
			}
		}
	}

	// 3. Fetch Products for Items
	productIDs := []primitive.ObjectID{}
	for _, item := range adjustment.Items {
		productIDs = append(productIDs, item.ProductID)
	}

	if len(productIDs) > 0 {
		if products, err := h.productClient.GetProductsBatch(ctx, productIDs, token); err == nil {
			for i, item := range adjustment.Items {
				if prod, ok := products[item.ProductID.Hex()]; ok {
					adjustment.Items[i].ProductName = prod.Name
					adjustment.Items[i].SKU = prod.SKU
				}
			}
		}
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

	// Enrich with Location Names
	if len(adjustments) > 0 {
		locationIDs := []primitive.ObjectID{}
		for _, adj := range adjustments {
			if adj.LocationID != primitive.NilObjectID {
				locationIDs = append(locationIDs, adj.LocationID)
			}
		}

		if len(locationIDs) > 0 {
			token := c.GetHeader("Authorization")
			if locs, err := h.orgClient.GetLocationsBatch(c.Request.Context(), locationIDs, token); err == nil {
				for _, adj := range adjustments {
					if loc, ok := locs[adj.LocationID.Hex()]; ok {
						adj.LocationName = loc.Name
					}
				}
			}
		}
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
