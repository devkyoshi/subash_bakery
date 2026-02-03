package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/yourusername/erp-system/services/procurement-service/internal/service"
	"github.com/yourusername/erp-system/shared/middleware"
	"github.com/yourusername/erp-system/shared/models"
	"github.com/yourusername/erp-system/shared/utils"
)

type ProcurementHandler struct {
	procurementService *service.ProcurementService
}

func NewProcurementHandler(procurementService *service.ProcurementService) *ProcurementHandler {
	return &ProcurementHandler{
		procurementService: procurementService,
	}
}

// ============== Supplier Handlers ==============

func (h *ProcurementHandler) CreateSupplier(c *gin.Context) {
	var req service.CreateSupplierRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error(), nil)
		return
	}

	orgIDStr := c.Param("org_id")
	orgID, err := primitive.ObjectIDFromHex(orgIDStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "INVALID_ORG_ID", "Invalid organization ID", nil)
		return
	}

	userIDStr := middleware.GetUserID(c)
	userID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, "INVALID_USER", "Invalid user ID", nil)
		return
	}

	supplier, err := h.procurementService.CreateSupplier(c.Request.Context(), orgID, req, userID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "CREATE_SUPPLIER_ERROR", err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, supplier, "Supplier created successfully")
}

func (h *ProcurementHandler) GetSupplier(c *gin.Context) {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "Invalid supplier ID", nil)
		return
	}

	supplier, err := h.procurementService.GetSupplier(c.Request.Context(), id)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "SUPPLIER_NOT_FOUND", err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, supplier, "Supplier retrieved successfully")
}

func (h *ProcurementHandler) ListSuppliers(c *gin.Context) {
	orgIDStr := c.Param("org_id")
	orgID, err := primitive.ObjectIDFromHex(orgIDStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "INVALID_ORG_ID", "Invalid organization ID", nil)
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	filters := make(map[string]interface{})
	if status := c.Query("status"); status != "" {
		filters["status"] = status
	}
	if search := c.Query("search"); search != "" {
		filters["search"] = search
	}

	suppliers, total, err := h.procurementService.ListSuppliers(c.Request.Context(), orgID, filters, page, limit)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "LIST_SUPPLIERS_ERROR", err.Error(), nil)
		return
	}

	utils.PaginatedResponse(c, suppliers, page, limit, total)
}

func (h *ProcurementHandler) UpdateSupplier(c *gin.Context) {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "Invalid supplier ID", nil)
		return
	}

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error(), nil)
		return
	}

	supplier, err := h.procurementService.UpdateSupplier(c.Request.Context(), id, updates)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "UPDATE_SUPPLIER_ERROR", err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, supplier, "Supplier updated successfully")
}

func (h *ProcurementHandler) DeleteSupplier(c *gin.Context) {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "Invalid supplier ID", nil)
		return
	}

	userIDStr := middleware.GetUserID(c)
	userID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, "INVALID_USER", "Invalid user ID", nil)
		return
	}

	if err := h.procurementService.DeleteSupplier(c.Request.Context(), id, userID); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "DELETE_SUPPLIER_ERROR", err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, nil, "Supplier deleted successfully")
}

// ============== Purchase Order Handlers ==============

func (h *ProcurementHandler) CreatePurchaseOrder(c *gin.Context) {
	var req service.CreatePurchaseOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error(), nil)
		return
	}

	orgIDStr := c.Param("org_id")
	orgID, err := primitive.ObjectIDFromHex(orgIDStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "INVALID_ORG_ID", "Invalid organization ID", nil)
		return
	}

	userIDStr := middleware.GetUserID(c)
	userID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, "INVALID_USER", "Invalid user ID", nil)
		return
	}

	po, err := h.procurementService.CreatePurchaseOrder(c.Request.Context(), orgID, req, userID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "CREATE_PO_ERROR", err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, po, "Purchase order created successfully")
}

func (h *ProcurementHandler) GetPurchaseOrder(c *gin.Context) {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "Invalid purchase order ID", nil)
		return
	}

	po, err := h.procurementService.GetPurchaseOrder(c.Request.Context(), id)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "PO_NOT_FOUND", err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, po, "Purchase order retrieved successfully")
}

func (h *ProcurementHandler) ListPurchaseOrders(c *gin.Context) {
	orgIDStr := c.Param("org_id")
	orgID, err := primitive.ObjectIDFromHex(orgIDStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "INVALID_ORG_ID", "Invalid organization ID", nil)
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	filters := make(map[string]interface{})
	if status := c.Query("status"); status != "" {
		filters["status"] = status
	}
	if supplierIDStr := c.Query("supplier_id"); supplierIDStr != "" {
		if supplierID, err := primitive.ObjectIDFromHex(supplierIDStr); err == nil {
			filters["supplier_id"] = supplierID
		}
	}
	if search := c.Query("search"); search != "" {
		filters["search"] = search
	}

	pos, total, err := h.procurementService.ListPurchaseOrders(c.Request.Context(), orgID, filters, page, limit)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "LIST_PO_ERROR", err.Error(), nil)
		return
	}

	utils.PaginatedResponse(c, pos, page, limit, total)
}

func (h *ProcurementHandler) ApprovePurchaseOrder(c *gin.Context) {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "Invalid purchase order ID", nil)
		return
	}

	userIDStr := middleware.GetUserID(c)
	userID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, "INVALID_USER", "Invalid user ID", nil)
		return
	}

	if err := h.procurementService.ApprovePurchaseOrder(c.Request.Context(), id, userID); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "APPROVE_PO_ERROR", err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, nil, "Purchase order approved successfully")
}

func (h *ProcurementHandler) UpdatePOStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "Invalid purchase order ID", nil)
		return
	}

	var req struct {
		Status models.POStatus `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error(), nil)
		return
	}

	if err := h.procurementService.UpdatePOStatus(c.Request.Context(), id, req.Status); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "UPDATE_STATUS_ERROR", err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, nil, "Purchase order status updated successfully")
}

func (h *ProcurementHandler) DeletePurchaseOrder(c *gin.Context) {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "Invalid purchase order ID", nil)
		return
	}

	userIDStr := middleware.GetUserID(c)
	userID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, "INVALID_USER", "Invalid user ID", nil)
		return
	}

	if err := h.procurementService.DeletePurchaseOrder(c.Request.Context(), id, userID); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "DELETE_PO_ERROR", err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, nil, "Purchase order deleted successfully")
}

// ============== GRN Handlers ==============

func (h *ProcurementHandler) CreateGRN(c *gin.Context) {
	var req service.CreateGRNRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error(), nil)
		return
	}

	orgIDStr := c.Param("org_id")
	orgID, err := primitive.ObjectIDFromHex(orgIDStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "INVALID_ORG_ID", "Invalid organization ID", nil)
		return
	}

	userIDStr := middleware.GetUserID(c)
	userID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, "INVALID_USER", "Invalid user ID", nil)
		return
	}

	grn, err := h.procurementService.CreateGRN(c.Request.Context(), orgID, req, userID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "CREATE_GRN_ERROR", err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, grn, "GRN created successfully")
}

func (h *ProcurementHandler) GetGRN(c *gin.Context) {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "Invalid GRN ID", nil)
		return
	}

	grn, err := h.procurementService.GetGRN(c.Request.Context(), id)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "GRN_NOT_FOUND", err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, grn, "GRN retrieved successfully")
}

func (h *ProcurementHandler) ListGRNs(c *gin.Context) {
	orgIDStr := c.Param("org_id")
	orgID, err := primitive.ObjectIDFromHex(orgIDStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "INVALID_ORG_ID", "Invalid organization ID", nil)
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	filters := make(map[string]interface{})
	if status := c.Query("status"); status != "" {
		filters["status"] = status
	}
	if poIDStr := c.Query("purchase_order_id"); poIDStr != "" {
		if poID, err := primitive.ObjectIDFromHex(poIDStr); err == nil {
			filters["purchase_order_id"] = poID
		}
	}
	if search := c.Query("search"); search != "" {
		filters["search"] = search
	}

	grns, total, err := h.procurementService.ListGRNs(c.Request.Context(), orgID, filters, page, limit)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "LIST_GRN_ERROR", err.Error(), nil)
		return
	}

	utils.PaginatedResponse(c, grns, page, limit, total)
}

func (h *ProcurementHandler) CompleteInspection(c *gin.Context) {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "Invalid GRN ID", nil)
		return
	}

	var req struct {
		QCStatus string `json:"qc_status" binding:"required"`
		QCNotes  string `json:"qc_notes,omitempty"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error(), nil)
		return
	}

	userIDStr := middleware.GetUserID(c)
	userID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, "INVALID_USER", "Invalid user ID", nil)
		return
	}

	if err := h.procurementService.CompleteInspection(c.Request.Context(), id, userID, req.QCStatus, req.QCNotes); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "INSPECTION_ERROR", err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, nil, "Inspection completed successfully")
}

// ============== Route Registration ==============

func (h *ProcurementHandler) RegisterRoutes(router *gin.RouterGroup, jwtManager *utils.JWTManager) {
	protected := router.Group("")
	protected.Use(middleware.AuthMiddleware(jwtManager))

	// Supplier routes
	protected.POST("/organizations/:org_id/suppliers", h.CreateSupplier)
	protected.GET("/organizations/:org_id/suppliers", h.ListSuppliers)
	protected.GET("/suppliers/:id", h.GetSupplier)
	protected.PUT("/suppliers/:id", h.UpdateSupplier)
	protected.DELETE("/suppliers/:id", h.DeleteSupplier)

	// Purchase Order routes
	protected.POST("/organizations/:org_id/purchase-orders", h.CreatePurchaseOrder)
	protected.GET("/organizations/:org_id/purchase-orders", h.ListPurchaseOrders)
	protected.GET("/purchase-orders/:id", h.GetPurchaseOrder)
	protected.PUT("/purchase-orders/:id/status", h.UpdatePOStatus)
	protected.POST("/purchase-orders/:id/approve", h.ApprovePurchaseOrder)
	protected.DELETE("/purchase-orders/:id", h.DeletePurchaseOrder)

	// GRN routes
	protected.POST("/organizations/:org_id/grns", h.CreateGRN)
	protected.GET("/organizations/:org_id/grns", h.ListGRNs)
	protected.GET("/grns/:id", h.GetGRN)
	protected.POST("/grns/:id/inspect", h.CompleteInspection)
}
