package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/erp-system/services/inventory-service/internal/service"
	"github.com/yourusername/erp-system/shared/middleware"
	"github.com/yourusername/erp-system/shared/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UnitChartHandler struct {
	unitChartService *service.UnitChartService
}

func NewUnitChartHandler(unitChartService *service.UnitChartService) *UnitChartHandler {
	return &UnitChartHandler{
		unitChartService: unitChartService,
	}
}

func (h *UnitChartHandler) RegisterRoutes(router *gin.RouterGroup, jwtManager *utils.JWTManager) {
	unitCharts := router.Group("/unit-charts")
	unitCharts.Use(middleware.AuthMiddleware(jwtManager))
	{
		unitCharts.POST("", h.CreateUnitChart)
		unitCharts.GET("", h.GetUnitCharts)
		unitCharts.GET("/:id", h.GetUnitChart)
		unitCharts.PUT("/:id", h.UpdateUnitChart)
		unitCharts.DELETE("/:id", h.DeleteUnitChart)
		unitCharts.GET("/conversion-rate", h.GetConversionRate)
	}
}

// CreateUnitChart godoc
// @Summary Create a new unit chart (conversion rule)
// @Tags UnitCharts
// @Accept json
// @Produce json
// @Param request body service.CreateUnitChartRequest true "Unit Chart Data"
// @Success 201 {object} models.UnitChart
// @Failure 400 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /unit-charts [post]
func (h *UnitChartHandler) CreateUnitChart(c *gin.Context) {
	var req service.CreateUnitChartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error(), nil)
		return
	}

	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated", nil)
		return
	}

	createdBy, err := primitive.ObjectIDFromHex(userID.(string))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "INVALID_USER_ID", err.Error(), nil)
		return
	}

	chart, err := h.unitChartService.CreateUnitChart(c.Request.Context(), req, createdBy)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "CREATE_FAILED", err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, chart, "Unit chart created successfully")
}

// GetUnitChart godoc
// @Summary Get a unit chart by ID
// @Tags UnitCharts
// @Produce json
// @Param id path string true "Unit Chart ID"
// @Success 200 {object} models.UnitChart
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Router /unit-charts/{id} [get]
func (h *UnitChartHandler) GetUnitChart(c *gin.Context) {
	idParam := c.Param("id")
	chartID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "Invalid unit chart ID", nil)
		return
	}

	chart, err := h.unitChartService.GetUnitChart(c.Request.Context(), chartID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "NOT_FOUND", err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, chart, "Unit chart retrieved successfully")
}

// GetUnitCharts godoc
// @Summary Get all unit charts with conversions
// @Tags UnitCharts
// @Produce json
// @Param active_only query bool false "Filter active only"
// @Success 200 {array} service.UnitResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /unit-charts [get]
func (h *UnitChartHandler) GetUnitCharts(c *gin.Context) {
	activeOnly := c.Query("active_only") == "true"

	charts, err := h.unitChartService.GetUnitCharts(c.Request.Context(), activeOnly)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "FETCH_FAILED", err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, charts, "Unit charts retrieved successfully")
}


// UpdateUnitChart godoc
// @Summary Update a unit chart
// @Tags UnitCharts
// @Accept json
// @Produce json
// @Param id path string true "Unit Chart ID"
// @Param request body service.UpdateUnitChartRequest true "Update Data"
// @Success 200 {object} models.UnitChart
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /unit-charts/{id} [put]
func (h *UnitChartHandler) UpdateUnitChart(c *gin.Context) {
	idParam := c.Param("id")
	chartID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "Invalid unit chart ID", nil)
		return
	}

	var req service.UpdateUnitChartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error(), nil)
		return
	}

	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated", nil)
		return
	}

	updatedBy, err := primitive.ObjectIDFromHex(userID.(string))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "INVALID_USER_ID", err.Error(), nil)
		return
	}

	chart, err := h.unitChartService.UpdateUnitChart(c.Request.Context(), chartID, req, updatedBy)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "UPDATE_FAILED", err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, chart, "Unit chart updated successfully")
}

// DeleteUnitChart godoc
// @Summary Delete a unit chart
// @Tags UnitCharts
// @Produce json
// @Param id path string true "Unit Chart ID"
// @Success 200 {object} utils.SuccessResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /unit-charts/{id} [delete]
func (h *UnitChartHandler) DeleteUnitChart(c *gin.Context) {
	idParam := c.Param("id")
	chartID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "Invalid unit chart ID", nil)
		return
	}

	if err := h.unitChartService.DeleteUnitChart(c.Request.Context(), chartID); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "DELETE_FAILED", err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, nil, "Unit chart deleted successfully")
}

// GetConversionRate godoc
// @Summary Get conversion rate between two units
// @Tags UnitCharts
// @Produce json
// @Param from_unit_id query string true "From Unit ID"
// @Param to_unit_id query string true "To Unit ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Router /unit-charts/conversion-rate [get]
func (h *UnitChartHandler) GetConversionRate(c *gin.Context) {
	fromUnitIDStr := c.Query("from_unit_id")
	toUnitIDStr := c.Query("to_unit_id")

	if fromUnitIDStr == "" || toUnitIDStr == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "MISSING_PARAMS", "from_unit_id and to_unit_id are required", nil)
		return
	}

	fromUnitID, err := primitive.ObjectIDFromHex(fromUnitIDStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "Invalid from_unit_id", nil)
		return
	}

	toUnitID, err := primitive.ObjectIDFromHex(toUnitIDStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "Invalid to_unit_id", nil)
		return
	}

	rate, err := h.unitChartService.GetConversionRate(c.Request.Context(), fromUnitID, toUnitID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "NOT_FOUND", err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, gin.H{
		"from_unit_id":    fromUnitID,
		"to_unit_id":      toUnitID,
		"conversion_rate": rate,
	}, "Conversion rate retrieved successfully")
}
