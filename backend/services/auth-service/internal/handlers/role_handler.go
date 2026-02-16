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

// GetRole retrieves a role by ID
// @Summary Get role by ID
// @Tags roles
// @Produce json
// @Param id path string true "Role ID"
// @Success 200 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /roles/{id} [get]
func (h *RoleHandler) GetRole(c *gin.Context) {
	roleID := c.Param("id")
	if roleID == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "Role ID is required", nil)
		return
	}

	role, err := h.roleService.GetRole(c.Request.Context(), roleID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "NOT_FOUND", err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, role, "Role retrieved successfully")
}

// UpdateRole updates a role
// @Summary Update role
// @Tags roles
// @Accept json
// @Produce json
// @Param id path string true "Role ID"
// @Success 200 {object} utils.Response
// @Router /roles/{id} [put]
func (h *RoleHandler) UpdateRole(c *gin.Context) {
	roleID := c.Param("id")
	var req service.UpdateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error(), nil)
		return
	}

	role, err := h.roleService.UpdateRole(c.Request.Context(), roleID, req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "UPDATE_FAILED", err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, role, "Role updated successfully")
}

// DeleteRole deletes a role
// @Summary Delete role
// @Tags roles
// @Param id path string true "Role ID"
// @Success 200 {object} utils.Response
// @Router /roles/{id} [delete]
func (h *RoleHandler) DeleteRole(c *gin.Context) {
	roleID := c.Param("id")
	if err := h.roleService.DeleteRole(c.Request.Context(), roleID); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "DELETE_FAILED", err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, nil, "Role deleted successfully")
}

// RegisterRoutes registers role routes
func (h *RoleHandler) RegisterRoutes(router *gin.RouterGroup, jwtManager *utils.JWTManager) {
	protected := router.Group("")
	protected.Use(middleware.AuthMiddleware(jwtManager))
	{
		protected.POST("/roles", h.CreateRole)
		protected.GET("/roles", h.ListRoles)
		protected.GET("/roles/:id", h.GetRole)
		protected.PUT("/roles/:id", h.UpdateRole)
		protected.DELETE("/roles/:id", h.DeleteRole)
		protected.POST("/roles/assign", h.AssignRole)
		protected.GET("/permissions", h.ListPermissions)
	}
}
