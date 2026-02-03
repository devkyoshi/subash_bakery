package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"github.com/yourusername/erp-system/services/subscription-service/internal/service"
	"github.com/yourusername/erp-system/shared/middleware"
	"github.com/yourusername/erp-system/shared/models"
	"github.com/yourusername/erp-system/shared/utils"
)

type SubscriptionHandler struct {
	subscriptionService *service.SubscriptionService
}

func NewSubscriptionHandler(subscriptionService *service.SubscriptionService) *SubscriptionHandler {
	return &SubscriptionHandler{
		subscriptionService: subscriptionService,
	}
}

// CreatePlan creates a new subscription plan
func (h *SubscriptionHandler) CreatePlan(c *gin.Context) {
	var req service.CreatePlanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error(), nil)
		return
	}

	userID := middleware.GetUserID(c)
	createdBy, _ := primitive.ObjectIDFromHex(userID)

	plan, err := h.subscriptionService.CreatePlan(c.Request.Context(), req, createdBy)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "CREATE_FAILED", err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, plan, "Plan created successfully")
}

// GetPlan retrieves a plan by ID
func (h *SubscriptionHandler) GetPlan(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "Invalid plan ID", nil)
		return
	}

	plan, err := h.subscriptionService.GetPlan(c.Request.Context(), id)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "NOT_FOUND", err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, plan, "Plan retrieved successfully")
}

// ListPlans returns all plans
func (h *SubscriptionHandler) ListPlans(c *gin.Context) {
	tier := c.Query("tier")
	publicOnly := c.Query("public") == "true"

	plans, err := h.subscriptionService.ListPlans(c.Request.Context(), tier, publicOnly)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "LIST_FAILED", err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, plans, "Plans retrieved successfully")
}

// Subscribe creates a subscription
func (h *SubscriptionHandler) Subscribe(c *gin.Context) {
	orgID, err := primitive.ObjectIDFromHex(c.Param("org_id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "Invalid organization ID", nil)
		return
	}

	var req service.SubscribeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error(), nil)
		return
	}

	planID, err := primitive.ObjectIDFromHex(req.PlanID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "Invalid plan ID", nil)
		return
	}

	subscription, err := h.subscriptionService.Subscribe(c.Request.Context(), orgID, planID, req.BillingCycle, req.StartTrial)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "SUBSCRIBE_FAILED", err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, subscription, "Subscription created successfully")
}

// GetSubscription retrieves an organization's subscription
func (h *SubscriptionHandler) GetSubscription(c *gin.Context) {
	orgID, err := primitive.ObjectIDFromHex(c.Param("org_id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "Invalid organization ID", nil)
		return
	}

	subscription, err := h.subscriptionService.GetSubscription(c.Request.Context(), orgID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "NOT_FOUND", err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, subscription, "Subscription retrieved successfully")
}

// UpdateUsage updates subscription usage
func (h *SubscriptionHandler) UpdateUsage(c *gin.Context) {
	orgID, err := primitive.ObjectIDFromHex(c.Param("org_id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "Invalid organization ID", nil)
		return
	}

	var req service.UpdateUsageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error(), nil)
		return
	}

	usage := models.SubscriptionUsage{
		Users:             req.Users,
		Companies:         req.Companies,
		Locations:         req.Locations,
		StorageUsedGB:     req.StorageUsedGB,
		APICallsUsed:      req.APICallsUsed,
		AICreditsUsed:     req.AICreditsUsed,
		WorkflowsActive:   req.WorkflowsActive,
		CustomFormsActive: req.CustomFormsActive,
	}

	if err := h.subscriptionService.UpdateUsage(c.Request.Context(), orgID, usage); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "UPDATE_FAILED", err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, nil, "Usage updated successfully")
}

// CancelSubscription cancels a subscription
func (h *SubscriptionHandler) CancelSubscription(c *gin.Context) {
	orgID, err := primitive.ObjectIDFromHex(c.Param("org_id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "Invalid organization ID", nil)
		return
	}

	var req struct {
		Reason string `json:"reason"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error(), nil)
		return
	}

	if err := h.subscriptionService.CancelSubscription(c.Request.Context(), orgID, req.Reason); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "CANCEL_FAILED", err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, nil, "Subscription cancelled successfully")
}

// ChangePlan changes the subscription plan
func (h *SubscriptionHandler) ChangePlan(c *gin.Context) {
	orgID, err := primitive.ObjectIDFromHex(c.Param("org_id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "Invalid organization ID", nil)
		return
	}

	var req service.ChangePlanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error(), nil)
		return
	}

	newPlanID, err := primitive.ObjectIDFromHex(req.NewPlanID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "Invalid plan ID", nil)
		return
	}

	subscription, err := h.subscriptionService.ChangePlan(c.Request.Context(), orgID, newPlanID, req.Immediate)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "CHANGE_PLAN_FAILED", err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, subscription, "Plan changed successfully")
}

// RegisterRoutes registers all subscription routes
func (h *SubscriptionHandler) RegisterRoutes(router *gin.RouterGroup, jwtManager *utils.JWTManager) {
	// Public routes (no auth required)
	public := router.Group("/plans")
	{
		public.GET("", h.ListPlans)
		public.GET("/:id", h.GetPlan)
	}

	// Protected routes
	protected := router.Group("")
	protected.Use(middleware.AuthMiddleware(jwtManager))

	// Plans (admin only in production)
	plans := protected.Group("/plans")
	{
		plans.POST("", h.CreatePlan)
	}

	// Subscriptions
	subscriptions := protected.Group("/organizations/:org_id/subscription")
	{
		subscriptions.POST("", h.Subscribe)
		subscriptions.GET("", h.GetSubscription)
		subscriptions.PUT("/usage", h.UpdateUsage)
		subscriptions.POST("/cancel", h.CancelSubscription)
		subscriptions.POST("/change-plan", h.ChangePlan)
	}
}
