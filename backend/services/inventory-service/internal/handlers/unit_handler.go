package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/erp-system/services/inventory-service/internal/service"
	"github.com/yourusername/erp-system/shared/middleware"
	"github.com/yourusername/erp-system/shared/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UnitHandler struct {
	unitService *service.UnitService
}

func NewUnitHandler(unitService *service.UnitService) *UnitHandler {
	return &UnitHandler{
		unitService: unitService,
	}
}

func (h *UnitHandler) RegisterRoutes(router *gin.RouterGroup, jwtManager *utils.JWTManager) {
	units := router.Group("/units")
	units.Use(middleware.AuthMiddleware(jwtManager))
	{
		units.POST("", h.CreateUnit)
		units.GET("", h.GetUnits)
		units.GET("/:id", h.GetUnit)
		units.PUT("/:id", h.UpdateUnit)
		units.DELETE("/:id", h.DeleteUnit)
	}
}

// CreateUnit godoc
// @Summary Create a new unit of measure
// @Tags Units
// @Accept json
// @Produce json
// @Param request body service.CreateUnitRequest true "Unit Data"
// @Success 201 {object} models.Unit
// @Failure 400 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /units [post]
func (h *UnitHandler) CreateUnit(c *gin.Context) {
	var req service.CreateUnitRequest
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

	unit, err := h.unitService.CreateUnit(c.Request.Context(), req, createdBy)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "CREATE_FAILED", err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, unit, "Unit created successfully")
}

// GetUnit godoc
// @Summary Get a unit by ID
// @Tags Units
// @Produce json
// @Param id path string true "Unit ID"
// @Success 200 {object} models.Unit
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Router /units/{id} [get]
func (h *UnitHandler) GetUnit(c *gin.Context) {
	idParam := c.Param("id")
	unitID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "Invalid unit ID", nil)
		return
	}

	unit, err := h.unitService.GetUnit(c.Request.Context(), unitID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "NOT_FOUND", err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, unit, "Unit retrieved successfully")
}

// GetUnits godoc
// @Summary Get all units
// @Tags Units
// @Produce json
// @Param unit_type query string false "Filter by unit type"
// @Param active_only query bool false "Filter active only"
// @Success 200 {array} models.Unit
// @Failure 500 {object} utils.ErrorResponse
// @Router /units [get]
func (h *UnitHandler) GetUnits(c *gin.Context) {
	unitType := c.Query("unit_type")
	activeOnly := c.Query("active_only") == "true"
	idsStr := c.Query("ids")

	var unitTypePtr *string
	if unitType != "" {
		unitTypePtr = &unitType
	}

	var ids []primitive.ObjectID
	if idsStr != "" {
		idList := strings.Split(idsStr, ",")
		for _, id := range idList {
			if oid, err := primitive.ObjectIDFromHex(strings.TrimSpace(id)); err == nil {
				ids = append(ids, oid)
			}
		}
	}

	units, err := h.unitService.GetUnits(c.Request.Context(), unitTypePtr, activeOnly, ids)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "FETCH_FAILED", err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, units, "Units retrieved successfully")
}

// UpdateUnit godoc
// @Summary Update a unit
// @Tags Units
// @Accept json
// @Produce json
// @Param id path string true "Unit ID"
// @Param request body service.UpdateUnitRequest true "Update Data"
// @Success 200 {object} models.Unit
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /units/{id} [put]
func (h *UnitHandler) UpdateUnit(c *gin.Context) {
	idParam := c.Param("id")
	unitID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "Invalid unit ID", nil)
		return
	}

	var req service.UpdateUnitRequest
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

	unit, err := h.unitService.UpdateUnit(c.Request.Context(), unitID, req, updatedBy)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "UPDATE_FAILED", err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, unit, "Unit updated successfully")
}

// DeleteUnit godoc
// @Summary Delete a unit
// @Tags Units
// @Produce json
// @Param id path string true "Unit ID"
// @Success 200 {object} utils.SuccessResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /units/{id} [delete]
func (h *UnitHandler) DeleteUnit(c *gin.Context) {
	idParam := c.Param("id")
	unitID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "Invalid unit ID", nil)
		return
	}

	if err := h.unitService.DeleteUnit(c.Request.Context(), unitID); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "DELETE_FAILED", err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, nil, "Unit deleted successfully")
}
