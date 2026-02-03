package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/erp-system/services/auth-service/internal/service"
	"github.com/yourusername/erp-system/shared/middleware"
	"github.com/yourusername/erp-system/shared/utils"
)

type RoleHandler struct {
	roleService *service.RoleService
}

func NewRoleHandler(roleService *service.RoleService) *RoleHandler {
	return &RoleHandler{
		roleService: roleService,
	}
}

// CreateRole creates a new role
// @Summary Create a new role
// @Tags roles
// @Accept json
// @Produce json
// @Param request body service.CreateRoleRequest true "Role details"
// @Success 201 {object} utils.Response
// @Router /roles [post]
func (h *RoleHandler) CreateRole(c *gin.Context) {
	var req service.CreateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error(), nil)
		return
	}

	role, err := h.roleService.CreateRole(c.Request.Context(), req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "CREATE_FAILED", err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, role, "Role created successfully")
}

// ListRoles retrieves all roles
// @Summary List roles
// @Tags roles
// @Produce json
// @Param organization_id query string false "Organization ID"
// @Success 200 {object} utils.Response
// @Router /roles [get]
func (h *RoleHandler) ListRoles(c *gin.Context) {
	orgID := c.Query("organization_id")
	roles, err := h.roleService.ListRoles(c.Request.Context(), orgID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "FETCH_FAILED", err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, roles, "Roles retrieved successfully")
}

// AssignRole assigns a role to a user
// @Summary Assign role to user
// @Tags roles
// @Accept json
// @Produce json
// @Param request body map[string]string true "User and Role IDs"
// @Success 200 {object} utils.Response
// @Router /roles/assign [post]
func (h *RoleHandler) AssignRole(c *gin.Context) {
	var req struct {
		UserID string `json:"user_id" binding:"required"`
		RoleID string `json:"role_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error(), nil)
		return
	}

	if err := h.roleService.AssignRole(c.Request.Context(), req.UserID, req.RoleID); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ASSIGNMENT_FAILED", err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, nil, "Role assigned successfully")
}

// ListPermissions retrieves all permissions
// @Summary List permissions
// @Tags roles
// @Produce json
// @Success 200 {object} utils.Response
// @Router /permissions [get]
func (h *RoleHandler) ListPermissions(c *gin.Context) {
	permissions, err := h.roleService.ListPermissions(c.Request.Context())
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "FETCH_FAILED", err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, permissions, "Permissions retrieved successfully")
}

// RegisterRoutes registers role routes
func (h *RoleHandler) RegisterRoutes(router *gin.RouterGroup, jwtManager *utils.JWTManager) {
	protected := router.Group("")
	protected.Use(middleware.AuthMiddleware(jwtManager))
	{
		protected.POST("/roles", h.CreateRole)
		protected.GET("/roles", h.ListRoles)
		protected.POST("/roles/assign", h.AssignRole)
		protected.GET("/permissions", h.ListPermissions)
	}
}
