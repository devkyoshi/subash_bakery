package service

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"github.com/yourusername/erp-system/services/subscription-service/internal/repository"
	"github.com/yourusername/erp-system/shared/models"
)

type SubscriptionService struct {
	planRepo         *repository.PlanRepository
	subscriptionRepo *repository.SubscriptionRepository
}

func NewSubscriptionService(
	planRepo *repository.PlanRepository,
	subscriptionRepo *repository.SubscriptionRepository,
) *SubscriptionService {
	return &SubscriptionService{
		planRepo:         planRepo,
		subscriptionRepo: subscriptionRepo,
	}
}

// CreatePlan creates a new subscription plan
func (s *SubscriptionService) CreatePlan(ctx context.Context, req CreatePlanRequest, createdBy primitive.ObjectID) (*models.SubscriptionPlan, error) {
	plan := &models.SubscriptionPlan{
		Name:                     req.Name,
		DisplayName:              req.DisplayName,
		Description:              req.Description,
		Tier:                     req.Tier,
		PriceMonthly:             req.PriceMonthly,
		PriceQuarterly:           req.PriceQuarterly,
		PriceYearly:              req.PriceYearly,
		Currency:                 req.Currency,
		TrialDays:                req.TrialDays,
		Applications:             req.Applications,
		Features:                 req.Features,
		MaxUsers:                 req.MaxUsers,
		MaxCompanies:             req.MaxCompanies,
		MaxLocations:             req.MaxLocations,
		StorageGB:                req.StorageGB,
		APICallsPerMonth:         req.APICallsPerMonth,
		MaxWorkflows:             req.MaxWorkflows,
		MaxCustomForms:           req.MaxCustomForms,
		AICreditsPerMonth:        req.AICreditsPerMonth,
		EnableAIAgent:            req.EnableAIAgent,
		EnableAdvancedAnalytics:  req.EnableAdvancedAnalytics,
		EnableWorkflowAutomation: req.EnableWorkflowAutomation,
		IsPublic:                 req.IsPublic,
		IsActive:                 true,
		IsFeatured:               req.IsFeatured,
		DisplayOrder:             req.DisplayOrder,
	}

	plan.CreatedBy = createdBy

	if err := s.planRepo.Create(ctx, plan); err != nil {
		return nil, err
	}

	return plan, nil
}

// GetPlan retrieves a plan by ID
func (s *SubscriptionService) GetPlan(ctx context.Context, id primitive.ObjectID) (*models.SubscriptionPlan, error) {
	plan, err := s.planRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if plan == nil {
		return nil, fmt.Errorf("plan not found")
	}
	return plan, nil
}

// ListPlans returns all plans
func (s *SubscriptionService) ListPlans(ctx context.Context, tier string, publicOnly bool) ([]*models.SubscriptionPlan, error) {
	var isPublic *bool
	if publicOnly {
		t := true
		isPublic = &t
	}

	return s.planRepo.List(ctx, tier, isPublic)
}

// Subscribe creates a subscription for an organization
func (s *SubscriptionService) Subscribe(ctx context.Context, orgID, planID primitive.ObjectID, billingCycle models.BillingCycle, startTrial bool) (*models.OrganizationSubscription, error) {
	// Check if plan exists
	plan, err := s.planRepo.FindByID(ctx, planID)
	if err != nil {
		return nil, err
	}
	if plan == nil {
		return nil, fmt.Errorf("plan not found")
	}

	// Check if organization already has an active subscription
	existing, err := s.subscriptionRepo.FindByOrganization(ctx, orgID)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, fmt.Errorf("organization already has an active subscription")
	}

	// Calculate price based on billing cycle
	var price float64
	switch billingCycle {
	case models.BillingCycleMonthly:
		price = plan.PriceMonthly
	case models.BillingCycleQuarterly:
		price = plan.PriceQuarterly
	case models.BillingCycleYearly:
		price = plan.PriceYearly
	default:
		return nil, fmt.Errorf("invalid billing cycle")
	}

	now := time.Now()
	subscription := &models.OrganizationSubscription{
		OrganizationID:     orgID,
		PlanID:             planID,
		Status:             models.SubscriptionStatusActive,
		BillingCycle:       billingCycle,
		CurrentPrice:       price,
		Currency:           plan.Currency,
		StartDate:          now,
		CurrentPeriodStart: now,
		AutoRenew:          true,
		CurrentUsage: models.SubscriptionUsage{
			LastUpdated: now,
		},
	}

	// Set end date based on billing cycle
	switch billingCycle {
	case models.BillingCycleMonthly:
		subscription.EndDate = now.AddDate(0, 1, 0)
		subscription.CurrentPeriodEnd = now.AddDate(0, 1, 0)
		subscription.NextBillingDate = now.AddDate(0, 1, 0)
	case models.BillingCycleQuarterly:
		subscription.EndDate = now.AddDate(0, 3, 0)
		subscription.CurrentPeriodEnd = now.AddDate(0, 3, 0)
		subscription.NextBillingDate = now.AddDate(0, 3, 0)
	case models.BillingCycleYearly:
		subscription.EndDate = now.AddDate(1, 0, 0)
		subscription.CurrentPeriodEnd = now.AddDate(1, 0, 0)
		subscription.NextBillingDate = now.AddDate(1, 0, 0)
	}

	// Handle trial
	if startTrial && plan.TrialDays > 0 {
		subscription.Status = models.SubscriptionStatusTrial
		trialStart := now
		trialEnd := now.AddDate(0, 0, plan.TrialDays)
		subscription.TrialStartDate = &trialStart
		subscription.TrialEndDate = &trialEnd
		subscription.CurrentPrice = 0 // No charge during trial
	}

	if err := s.subscriptionRepo.Create(ctx, subscription); err != nil {
		return nil, err
	}

	return subscription, nil
}

// GetSubscription retrieves an organization's subscription
func (s *SubscriptionService) GetSubscription(ctx context.Context, orgID primitive.ObjectID) (*models.OrganizationSubscription, error) {
	subscription, err := s.subscriptionRepo.FindByOrganization(ctx, orgID)
	if err != nil {
		return nil, err
	}
	if subscription == nil {
		return nil, fmt.Errorf("no active subscription found")
	}
	return subscription, nil
}

// UpdateUsage updates subscription usage
func (s *SubscriptionService) UpdateUsage(ctx context.Context, orgID primitive.ObjectID, usage models.SubscriptionUsage) error {
	subscription, err := s.subscriptionRepo.FindByOrganization(ctx, orgID)
	if err != nil {
		return err
	}
	if subscription == nil {
		return fmt.Errorf("no active subscription found")
	}

	return s.subscriptionRepo.UpdateUsage(ctx, subscription.ID, usage)
}

// CancelSubscription cancels a subscription
func (s *SubscriptionService) CancelSubscription(ctx context.Context, orgID primitive.ObjectID, reason string) error {
	subscription, err := s.subscriptionRepo.FindByOrganization(ctx, orgID)
	if err != nil {
		return err
	}
	if subscription == nil {
		return fmt.Errorf("no active subscription found")
	}

	return s.subscriptionRepo.Cancel(ctx, subscription.ID, reason)
}

// ChangePlan changes the subscription plan
func (s *SubscriptionService) ChangePlan(ctx context.Context, orgID, newPlanID primitive.ObjectID, immediate bool) (*models.OrganizationSubscription, error) {
	// Get current subscription
	subscription, err := s.subscriptionRepo.FindByOrganization(ctx, orgID)
	if err != nil {
		return nil, err
	}
	if subscription == nil {
		return nil, fmt.Errorf("no active subscription found")
	}

	// Get new plan
	newPlan, err := s.planRepo.FindByID(ctx, newPlanID)
	if err != nil {
		return nil, err
	}
	if newPlan == nil {
		return nil, fmt.Errorf("new plan not found")
	}

	if immediate {
		// Immediate upgrade/downgrade
		subscription.PlanID = newPlanID

		// Update price based on billing cycle
		switch subscription.BillingCycle {
		case models.BillingCycleMonthly:
			subscription.CurrentPrice = newPlan.PriceMonthly
		case models.BillingCycleQuarterly:
			subscription.CurrentPrice = newPlan.PriceQuarterly
		case models.BillingCycleYearly:
			subscription.CurrentPrice = newPlan.PriceYearly
		}
	} else {
		// Schedule downgrade for next billing period
		subscription.WillDowngrade = true
		subscription.DowngradePlanID = &newPlanID
	}

	if err := s.subscriptionRepo.Update(ctx, subscription); err != nil {
		return nil, err
	}

	return subscription, nil
}

// Request/Response DTOs
type CreatePlanRequest struct {
	Name                     string                 `json:"name" binding:"required"`
	DisplayName              string                 `json:"display_name" binding:"required"`
	Description              string                 `json:"description"`
	Tier                     models.PlanTier        `json:"tier" binding:"required"`
	PriceMonthly             float64                `json:"price_monthly"`
	PriceQuarterly           float64                `json:"price_quarterly"`
	PriceYearly              float64                `json:"price_yearly"`
	Currency                 string                 `json:"currency"`
	TrialDays                int                    `json:"trial_days"`
	Applications             []primitive.ObjectID   `json:"applications"`
	Features                 []models.PlanFeature   `json:"features"`
	MaxUsers                 int                    `json:"max_users"`
	MaxCompanies             int                    `json:"max_companies"`
	MaxLocations             int                    `json:"max_locations"`
	StorageGB                float64                `json:"storage_gb"`
	APICallsPerMonth         int64                  `json:"api_calls_per_month"`
	MaxWorkflows             int                    `json:"max_workflows"`
	MaxCustomForms           int                    `json:"max_custom_forms"`
	AICreditsPerMonth        int                    `json:"ai_credits_per_month"`
	EnableAIAgent            bool                   `json:"enable_ai_agent"`
	EnableAdvancedAnalytics  bool                   `json:"enable_advanced_analytics"`
	EnableWorkflowAutomation bool                   `json:"enable_workflow_automation"`
	IsPublic                 bool                   `json:"is_public"`
	IsFeatured               bool                   `json:"is_featured"`
	DisplayOrder             int                    `json:"display_order"`
}

type SubscribeRequest struct {
	PlanID       string                `json:"plan_id" binding:"required"`
	BillingCycle models.BillingCycle   `json:"billing_cycle" binding:"required"`
	StartTrial   bool                  `json:"start_trial"`
}

type UpdateUsageRequest struct {
	Users             int     `json:"users"`
	Companies         int     `json:"companies"`
	Locations         int     `json:"locations"`
	StorageUsedGB     float64 `json:"storage_used_gb"`
	APICallsUsed      int64   `json:"api_calls_used"`
	AICreditsUsed     int     `json:"ai_credits_used"`
	WorkflowsActive   int     `json:"workflows_active"`
	CustomFormsActive int     `json:"custom_forms_active"`
}

type ChangePlanRequest struct {
	NewPlanID string `json:"new_plan_id" binding:"required"`
	Immediate bool   `json:"immediate"`
}
