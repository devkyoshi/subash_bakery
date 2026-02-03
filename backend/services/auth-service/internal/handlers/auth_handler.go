package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/erp-system/services/auth-service/internal/service"
	"github.com/yourusername/erp-system/shared/middleware"
	"github.com/yourusername/erp-system/shared/utils"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Register handles user registration
// @Summary Register a new user
// @Tags auth
// @Accept json
// @Produce json
// @Param request body service.RegisterRequest true "Registration details"
// @Success 201 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req service.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error(), nil)
		return
	}

	// Validate email
	if !utils.IsValidEmail(req.Email) {
		utils.ErrorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid email format", nil)
		return
	}

	// Validate password
	if !utils.IsValidPassword(req.Password) {
		utils.ErrorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", "Password must be at least 8 characters", nil)
		return
	}

	authResponse, err := h.authService.Register(c.Request.Context(), req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "REGISTRATION_FAILED", err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, authResponse, "User registered successfully")
}

// Login handles user login
// @Summary Login
// @Tags auth
// @Accept json
// @Produce json
// @Param request body service.LoginRequest true "Login credentials"
// @Success 200 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req service.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error(), nil)
		return
	}

	userAgent := c.GetHeader("User-Agent")
	ipAddress := c.ClientIP()

	authResponse, err := h.authService.Login(c.Request.Context(), req, userAgent, ipAddress)
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, "LOGIN_FAILED", err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, authResponse, "Login successful")
}

// RefreshToken handles token refresh
// @Summary Refresh access token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body map[string]string true "Refresh token"
// @Success 200 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Router /auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error(), nil)
		return
	}

	authResponse, err := h.authService.RefreshToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, "REFRESH_FAILED", err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, authResponse, "Token refreshed successfully")
}

// Logout handles user logout
// @Summary Logout
// @Tags auth
// @Accept json
// @Produce json
// @Param request body map[string]string true "Refresh token"
// @Success 200 {object} utils.Response
// @Router /auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error(), nil)
		return
	}

	if err := h.authService.Logout(c.Request.Context(), req.RefreshToken); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "LOGOUT_FAILED", err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, nil, "Logged out successfully")
}

// GetGoogleLoginURL returns the Google OAuth login URL
// @Summary Get Google login URL
// @Tags auth
// @Produce json
// @Success 200 {object} utils.Response
// @Router /auth/google [get]
func (h *AuthHandler) GetGoogleLoginURL(c *gin.Context) {
	url := h.authService.GetGoogleLoginURL()
	utils.SuccessResponse(c, http.StatusOK, gin.H{"url": url}, "Google login URL generated")
}

// GoogleCallback handles the Google OAuth callback
// @Summary Google OAuth callback
// @Tags auth
// @Produce json
// @Param code query string true "Authorization code"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Router /auth/google/callback [get]
func (h *AuthHandler) GoogleCallback(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", "Authorization code is required", nil)
		return
	}

	userAgent := c.GetHeader("User-Agent")
	ipAddress := c.ClientIP()

	authResponse, err := h.authService.GoogleLogin(c.Request.Context(), code, userAgent, ipAddress)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "GOOGLE_LOGIN_FAILED", err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, authResponse, "Google login successful")
}

// GetMe returns the current user's profile
// @Summary Get current user profile
// @Tags auth
// @Security BearerAuth
// @Produce json
// @Success 200 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Router /auth/me [get]
func (h *AuthHandler) GetMe(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == "" {
		utils.ErrorResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated", nil)
		return
	}

	user, err := h.authService.GetMe(c.Request.Context(), userID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "USER_NOT_FOUND", err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, user, "User profile retrieved successfully")
}

// UpdateUserOrganization updates a user's organization
// @Summary Update user organization
// @Tags auth
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param request body map[string]string true "Organization ID"
// @Success 200 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Router /auth/users/{id}/organization [put]
func (h *AuthHandler) UpdateUserOrganization(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "User ID is required", nil)
		return
	}

	var req struct {
		OrganizationID string `json:"organization_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error(), nil)
		return
	}

	user, err := h.authService.UpdateUserOrganization(c.Request.Context(), userID, req.OrganizationID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "UPDATE_FAILED", err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, user, "User organization updated successfully")
}

// RegisterRoutes registers all auth routes
func (h *AuthHandler) RegisterRoutes(router *gin.RouterGroup, jwtManager *utils.JWTManager) {
	auth := router.Group("/auth")
	{
		auth.POST("/register", h.Register)
		auth.POST("/login", h.Login)
		auth.POST("/refresh", h.RefreshToken)
		auth.POST("/logout", h.Logout)
		auth.GET("/google", h.GetGoogleLoginURL)
		auth.GET("/google/callback", h.GoogleCallback)

		// Protected routes
		protected := auth.Group("")
		protected.Use(middleware.AuthMiddleware(jwtManager))
		{
			protected.GET("/me", h.GetMe)
			protected.PUT("/users/:id/organization", h.UpdateUserOrganization)
		}
	}
}
