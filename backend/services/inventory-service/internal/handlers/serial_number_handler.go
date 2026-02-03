package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/yourusername/erp-system/services/inventory-service/internal/service"
	"github.com/yourusername/erp-system/shared/middleware"
	"github.com/yourusername/erp-system/shared/utils"
)

type SerialNumberHandler struct {
	serialNumberService *service.SerialNumberService
}

func NewSerialNumberHandler(serialNumberService *service.SerialNumberService) *SerialNumberHandler {
	return &SerialNumberHandler{
		serialNumberService: serialNumberService,
	}
}

func (h *SerialNumberHandler) CreateSerialNumber(c *gin.Context) {
	orgIDStr := c.Param("org_id")
	orgID, err := primitive.ObjectIDFromHex(orgIDStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "Invalid organization ID", nil)
		return
	}

	var req service.CreateSerialNumberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error(), nil)
		return
	}

	userID := middleware.GetUserID(c)
	createdBy, _ := primitive.ObjectIDFromHex(userID)

	serialNumber, err := h.serialNumberService.CreateSerialNumber(c.Request.Context(), orgID, req, createdBy)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "CREATE_FAILED", err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, serialNumber, "Serial number created successfully")
}

func (h *SerialNumberHandler) GetSerialNumber(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "Invalid serial number ID", nil)
		return
	}

	serialNumber, err := h.serialNumberService.GetSerialNumber(c.Request.Context(), id)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "NOT_FOUND", err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, serialNumber, "Serial number retrieved successfully")
}

func (h *SerialNumberHandler) GetSerialNumberBySerial(c *gin.Context) {
	orgIDStr := c.Param("org_id")
	orgID, err := primitive.ObjectIDFromHex(orgIDStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "Invalid organization ID", nil)
		return
	}

	serialNo := c.Param("serial_no")
	if serialNo == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "MISSING_PARAMS", "serial_no required", nil)
		return
	}

	serialNumber, err := h.serialNumberService.GetSerialNumberBySerial(c.Request.Context(), orgID, serialNo)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "NOT_FOUND", err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, serialNumber, "Serial number retrieved successfully")
}

func (h *SerialNumberHandler) ListSerialNumbers(c *gin.Context) {
	productIDStr := c.Param("product_id")
	productID, err := primitive.ObjectIDFromHex(productIDStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "Invalid product ID", nil)
		return
	}

	filters := make(map[string]interface{})
	if status := c.Query("status"); status != "" {
		filters["status"] = status
	}
	if locationID := c.Query("location_id"); locationID != "" {
		filters["location_id"] = locationID
	}
	if available := c.Query("available"); available == "true" {
		filters["is_available"] = true
	} else if available == "false" {
		filters["is_available"] = false
	}

	serialNumbers, err := h.serialNumberService.ListSerialNumbersByProduct(c.Request.Context(), productID, filters)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "LIST_FAILED", err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, serialNumbers, "Serial numbers retrieved successfully")
}

func (h *SerialNumberHandler) UpdateSerialNumber(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "Invalid serial number ID", nil)
		return
	}

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error(), nil)
		return
	}

	if err := h.serialNumberService.UpdateSerialNumber(c.Request.Context(), id, updates); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "UPDATE_FAILED", err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, nil, "Serial number updated successfully")
}

func (h *SerialNumberHandler) AllocateSerialNumber(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "Invalid serial number ID", nil)
		return
	}

	var req struct {
		CustomerID   string `json:"customer_id" binding:"required"`
		SalesOrderID string `json:"sales_order_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error(), nil)
		return
	}

	customerID, _ := primitive.ObjectIDFromHex(req.CustomerID)
	salesOrderID, _ := primitive.ObjectIDFromHex(req.SalesOrderID)

	if err := h.serialNumberService.AllocateSerialNumber(c.Request.Context(), id, customerID, salesOrderID); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ALLOCATE_FAILED", err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, nil, "Serial number allocated successfully")
}

func (h *SerialNumberHandler) MarkSerialAsSold(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "Invalid serial number ID", nil)
		return
	}

	if err := h.serialNumberService.MarkAsSold(c.Request.Context(), id); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "UPDATE_FAILED", err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, nil, "Serial number marked as sold successfully")
}

func (h *SerialNumberHandler) DeleteSerialNumber(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "Invalid serial number ID", nil)
		return
	}

	userID := middleware.GetUserID(c)
	deletedBy, _ := primitive.ObjectIDFromHex(userID)

	if err := h.serialNumberService.DeleteSerialNumber(c.Request.Context(), id, deletedBy); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "DELETE_FAILED", err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, nil, "Serial number deleted successfully")
}
