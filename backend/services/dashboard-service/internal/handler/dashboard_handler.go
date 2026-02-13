package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/erp-system/services/dashboard-service/internal/service"
	"github.com/yourusername/erp-system/shared/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DashboardHandler struct {
	aggregationService *service.AggregationService
}

func NewDashboardHandler(aggregationService *service.AggregationService) *DashboardHandler {
	return &DashboardHandler{
		aggregationService: aggregationService,
	}
}

func (h *DashboardHandler) RegisterRoutes(router *gin.RouterGroup) {
	router.GET("/dashboard/overview", h.GetOverview)
}

func (h *DashboardHandler) GetOverview(c *gin.Context) {
	orgIDStr := c.Query("organization_id")
	if orgIDStr == "" {
		orgIDStr = c.GetHeader("x-organization-id")
	}

	if orgIDStr == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "MISSING_ORG_ID", "Organization ID required", nil)
		return
	}

	orgID, err := primitive.ObjectIDFromHex(orgIDStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "INVALID_ORG_ID", "Invalid Organization ID", nil)
		return
	}

	overview, err := h.aggregationService.GetDashboardOverview(c.Request.Context(), orgID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, overview, "Dashboard overview retrieved successfully")
}
