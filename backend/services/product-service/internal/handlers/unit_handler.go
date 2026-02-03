package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/erp-system/services/product-service/internal/service"
	"github.com/yourusername/erp-system/shared/middleware"
	"github.com/yourusername/erp-system/shared/utils"
)

type UnitHandler struct {
	unitService *service.UnitService
}

func NewUnitHandler(unitService *service.UnitService) *UnitHandler {
	return &UnitHandler{
		unitService: unitService,
	}
}

// RegisterUnitRoutes registers unit routes
func (h *UnitHandler) RegisterUnitRoutes(router *gin.RouterGroup, jwtManager *utils.JWTManager) {
	protected := router.Group("")
	protected.Use(middleware.AuthMiddleware(jwtManager))

	units := protected.Group("/units")
	{
		units.GET("", h.ListUnits)
	}
}

// ListUnits retrieves all units basic info
// @Summary List units
// @Description Retrieves all available units
// @Tags Units
// @Produce json
// @Success 200 {object} utils.Response
// @Router /units [get]
func (h *UnitHandler) ListUnits(c *gin.Context) {
	units, err := h.unitService.ListUnits(c.Request.Context())
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "FETCH_FAILED", err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, units, "Units retrieved successfully")
}
