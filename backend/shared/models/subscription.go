package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BillingCycle string

const (
	BillingCycleMonthly    BillingCycle = "monthly"
	BillingCycleQuarterly  BillingCycle = "quarterly"
	BillingCycleYearly     BillingCycle = "yearly"
	BillingCycleBiAnnually BillingCycle = "bi_annually"
)

type SubscriptionStatus string

const (
	SubscriptionStatusTrial      SubscriptionStatus = "trial"
	SubscriptionStatusActive     SubscriptionStatus = "active"
	SubscriptionStatusPastDue    SubscriptionStatus = "past_due"
	SubscriptionStatusSuspended  SubscriptionStatus = "suspended"
	SubscriptionStatusCancelled  SubscriptionStatus = "cancelled"
	SubscriptionStatusExpired    SubscriptionStatus = "expired"
)

type PlanTier string

const (
	PlanTierFree       PlanTier = "free"
	PlanTierBasic      PlanTier = "basic"
	PlanTierPro        PlanTier = "pro"
	PlanTierEnterprise PlanTier = "enterprise"
	PlanTierCustom     PlanTier = "custom"
)

// SubscriptionPlan with advanced features
type SubscriptionPlan struct {
	BaseModel          `bson:",inline"`
	Name               string               `json:"name" bson:"name"`
	DisplayName        string               `json:"display_name" bson:"display_name"`
	Description        string               `json:"description" bson:"description"`
	Tier               PlanTier             `json:"tier" bson:"tier"`

	// Pricing
	PriceMonthly       float64              `json:"price_monthly" bson:"price_monthly"`
	PriceQuarterly     float64              `json:"price_quarterly" bson:"price_quarterly"`
	PriceYearly        float64              `json:"price_yearly" bson:"price_yearly"`
	PriceBiAnnually    float64              `json:"price_bi_annually" bson:"price_bi_annually"`
	Currency           string               `json:"currency" bson:"currency"`
	TrialDays          int                  `json:"trial_days" bson:"trial_days"`
	SetupFee           float64              `json:"setup_fee" bson:"setup_fee"`

	// Applications & Features
	Applications       []primitive.ObjectID `json:"applications" bson:"applications"`
	Features           []PlanFeature        `json:"features" bson:"features"`
	AddOns             []primitive.ObjectID `json:"add_ons,omitempty" bson:"add_ons,omitempty"` // Available add-ons

	// Limits
	MaxUsers           int                  `json:"max_users" bson:"max_users"` // 0 = unlimited
	MaxCompanies       int                  `json:"max_companies" bson:"max_companies"`
	MaxLocations       int                  `json:"max_locations" bson:"max_locations"`
	StorageGB          float64              `json:"storage_gb" bson:"storage_gb"`
	APICallsPerMonth   int64                `json:"api_calls_per_month" bson:"api_calls_per_month"`
	FileUploadSizeMB   int                  `json:"file_upload_size_mb" bson:"file_upload_size_mb"`
	MaxWorkflows       int                  `json:"max_workflows" bson:"max_workflows"`
	MaxCustomForms     int                  `json:"max_custom_forms" bson:"max_custom_forms"`

	// AI & Advanced Features
	AICreditsPerMonth  int                  `json:"ai_credits_per_month" bson:"ai_credits_per_month"`
	EnableAIAgent      bool                 `json:"enable_ai_agent" bson:"enable_ai_agent"`
	EnableAdvancedAnalytics bool            `json:"enable_advanced_analytics" bson:"enable_advanced_analytics"`
	EnableWorkflowAutomation bool           `json:"enable_workflow_automation" bson:"enable_workflow_automation"`

	// Visibility & Status
	IsPublic           bool                 `json:"is_public" bson:"is_public"`
	IsActive           bool                 `json:"is_active" bson:"is_active"`
	IsFeatured         bool                 `json:"is_featured" bson:"is_featured"`
	DisplayOrder       int                  `json:"display_order" bson:"display_order"`

	// Marketing
	Badge              string               `json:"badge,omitempty" bson:"badge,omitempty"` // "Most Popular", "Best Value"
	RecommendedFor     []string             `json:"recommended_for,omitempty" bson:"recommended_for,omitempty"`

	// Metadata
	Metadata           map[string]interface{} `json:"metadata,omitempty" bson:"metadata,omitempty"`
}

type PlanFeature struct {
	Name        string `json:"name" bson:"name"`
	DisplayName string `json:"display_name" bson:"display_name"`
	Description string `json:"description,omitempty" bson:"description,omitempty"`
	IsIncluded  bool   `json:"is_included" bson:"is_included"`
	Limit       int    `json:"limit,omitempty" bson:"limit,omitempty"` // 0 = unlimited, -1 = not included
	Category    string `json:"category,omitempty" bson:"category,omitempty"`
}

// OrganizationSubscription with billing history
type OrganizationSubscription struct {
	BaseModel          `bson:",inline"`
	OrganizationID     primitive.ObjectID `json:"organization_id" bson:"organization_id"`
	PlanID             primitive.ObjectID `json:"plan_id" bson:"plan_id"`
	Status             SubscriptionStatus `json:"status" bson:"status"`
	BillingCycle       BillingCycle       `json:"billing_cycle" bson:"billing_cycle"`

	// Pricing
	CurrentPrice       float64            `json:"current_price" bson:"current_price"`
	Currency           string             `json:"currency" bson:"currency"`
	Discount           *SubscriptionDiscount `json:"discount,omitempty" bson:"discount,omitempty"`

	// Dates
	StartDate          time.Time          `json:"start_date" bson:"start_date"`
	EndDate            time.Time          `json:"end_date" bson:"end_date"`
	CurrentPeriodStart time.Time          `json:"current_period_start" bson:"current_period_start"`
	CurrentPeriodEnd   time.Time          `json:"current_period_end" bson:"current_period_end"`
	TrialStartDate     *time.Time         `json:"trial_start_date,omitempty" bson:"trial_start_date,omitempty"`
	TrialEndDate       *time.Time         `json:"trial_end_date,omitempty" bson:"trial_end_date,omitempty"`

	// Renewal & Cancellation
	AutoRenew          bool               `json:"auto_renew" bson:"auto_renew"`
	WillDowngrade      bool               `json:"will_downgrade" bson:"will_downgrade"`
	DowngradePlanID    *primitive.ObjectID `json:"downgrade_plan_id,omitempty" bson:"downgrade_plan_id,omitempty"`
	CancelledAt        *time.Time         `json:"cancelled_at,omitempty" bson:"cancelled_at,omitempty"`
	CancellationReason string             `json:"cancellation_reason,omitempty" bson:"cancellation_reason,omitempty"`

	// Add-ons
	ActiveAddOns       []SubscriptionAddOn `json:"active_add_ons,omitempty" bson:"active_add_ons,omitempty"`

	// Payment
	PaymentMethodID    string             `json:"payment_method_id,omitempty" bson:"payment_method_id,omitempty"`
	LastPaymentDate    *time.Time         `json:"last_payment_date,omitempty" bson:"last_payment_date,omitempty"`
	NextBillingDate    time.Time          `json:"next_billing_date" bson:"next_billing_date"`

	// Usage tracking
	CurrentUsage       SubscriptionUsage  `json:"current_usage" bson:"current_usage"`

	// Metadata
	Metadata           map[string]interface{} `json:"metadata,omitempty" bson:"metadata,omitempty"`
}

type SubscriptionDiscount struct {
	Code        string    `json:"code" bson:"code"`
	Type        string    `json:"type" bson:"type"` // percentage, fixed
	Value       float64   `json:"value" bson:"value"`
	Duration    string    `json:"duration" bson:"duration"` // forever, once, repeating
	ValidUntil  *time.Time `json:"valid_until,omitempty" bson:"valid_until,omitempty"`
}

type SubscriptionAddOn struct {
	AddOnID     primitive.ObjectID `json:"add_on_id" bson:"add_on_id"`
	Name        string             `json:"name" bson:"name"`
	Price       float64            `json:"price" bson:"price"`
	Quantity    int                `json:"quantity" bson:"quantity"`
	AddedAt     time.Time          `json:"added_at" bson:"added_at"`
	IsActive    bool               `json:"is_active" bson:"is_active"`
}

type SubscriptionUsage struct {
	Users          int     `json:"users" bson:"users"`
	Companies      int     `json:"companies" bson:"companies"`
	Locations      int     `json:"locations" bson:"locations"`
	StorageUsedGB  float64 `json:"storage_used_gb" bson:"storage_used_gb"`
	APICallsUsed   int64   `json:"api_calls_used" bson:"api_calls_used"`
	AICreditsUsed  int     `json:"ai_credits_used" bson:"ai_credits_used"`
	WorkflowsActive int    `json:"workflows_active" bson:"workflows_active"`
	CustomFormsActive int  `json:"custom_forms_active" bson:"custom_forms_active"`
	LastUpdated    time.Time `json:"last_updated" bson:"last_updated"`
}

// PlanAddOn for additional features
type PlanAddOn struct {
	BaseModel       `bson:",inline"`
	Name            string  `json:"name" bson:"name"`
	DisplayName     string  `json:"display_name" bson:"display_name"`
	Description     string  `json:"description" bson:"description"`
	PriceMonthly    float64 `json:"price_monthly" bson:"price_monthly"`
	PriceYearly     float64 `json:"price_yearly" bson:"price_yearly"`
	Currency        string  `json:"currency" bson:"currency"`
	Type            string  `json:"type" bson:"type"` // storage, users, ai_credits, etc.
	Quantity        int     `json:"quantity" bson:"quantity"` // Units provided
	IsActive        bool    `json:"is_active" bson:"is_active"`
	CompatiblePlans []primitive.ObjectID `json:"compatible_plans" bson:"compatible_plans"`
	Metadata        map[string]interface{} `json:"metadata,omitempty" bson:"metadata,omitempty"`
}

// FeatureFlag for A/B testing and gradual rollouts
type FeatureFlag struct {
	BaseModel      `bson:",inline"`
	OrganizationID *primitive.ObjectID `json:"organization_id,omitempty" bson:"organization_id,omitempty"` // Nil = global
	Feature        string             `json:"feature" bson:"feature"`
	DisplayName    string             `json:"display_name" bson:"display_name"`
	Description    string             `json:"description" bson:"description"`
	IsEnabled      bool               `json:"is_enabled" bson:"is_enabled"`
	RolloutPercent int                `json:"rollout_percent" bson:"rollout_percent"` // 0-100
	Config         map[string]interface{} `json:"config,omitempty" bson:"config,omitempty"`
	TargetTiers    []PlanTier         `json:"target_tiers,omitempty" bson:"target_tiers,omitempty"` // Which tiers get this
	ExpiresAt      *time.Time         `json:"expires_at,omitempty" bson:"expires_at,omitempty"`
	Metadata       map[string]interface{} `json:"metadata,omitempty" bson:"metadata,omitempty"`
}

// Invoice for billing
type Invoice struct {
	BaseModel          `bson:",inline"`
	OrganizationID     primitive.ObjectID `json:"organization_id" bson:"organization_id"`
	SubscriptionID     primitive.ObjectID `json:"subscription_id" bson:"subscription_id"`
	InvoiceNumber      string             `json:"invoice_number" bson:"invoice_number"`
	Status             InvoiceStatus      `json:"status" bson:"status"`

	// Amounts
	Subtotal           float64            `json:"subtotal" bson:"subtotal"`
	Tax                float64            `json:"tax" bson:"tax"`
	Discount           float64            `json:"discount" bson:"discount"`
	Total              float64            `json:"total" bson:"total"`
	AmountPaid         float64            `json:"amount_paid" bson:"amount_paid"`
	AmountDue          float64            `json:"amount_due" bson:"amount_due"`
	Currency           string             `json:"currency" bson:"currency"`

	// Line items
	LineItems          []InvoiceLineItem  `json:"line_items" bson:"line_items"`

	// Dates
	IssueDate          time.Time          `json:"issue_date" bson:"issue_date"`
	DueDate            time.Time          `json:"due_date" bson:"due_date"`
	PaidAt             *time.Time         `json:"paid_at,omitempty" bson:"paid_at,omitempty"`

	// Payment
	PaymentMethodID    string             `json:"payment_method_id,omitempty" bson:"payment_method_id,omitempty"`
	TransactionID      string             `json:"transaction_id,omitempty" bson:"transaction_id,omitempty"`

	// PDF
	PDFUrl             string             `json:"pdf_url,omitempty" bson:"pdf_url,omitempty"`

	// Metadata
	Metadata           map[string]interface{} `json:"metadata,omitempty" bson:"metadata,omitempty"`
}

type InvoiceStatus string

const (
	InvoiceStatusDraft   InvoiceStatus = "draft"
	InvoiceStatusOpen    InvoiceStatus = "open"
	InvoiceStatusPaid    InvoiceStatus = "paid"
	InvoiceStatusVoid    InvoiceStatus = "void"
	InvoiceStatusUncollectible InvoiceStatus = "uncollectible"
)

type InvoiceLineItem struct {
	Description string  `json:"description" bson:"description"`
	Quantity    int     `json:"quantity" bson:"quantity"`
	UnitPrice   float64 `json:"unit_price" bson:"unit_price"`
	Amount      float64 `json:"amount" bson:"amount"`
	Type        string  `json:"type" bson:"type"` // subscription, add_on, usage
}

// UsageRecord for usage-based billing
type UsageRecord struct {
	ID             primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	OrganizationID primitive.ObjectID `json:"organization_id" bson:"organization_id"`
	SubscriptionID primitive.ObjectID `json:"subscription_id" bson:"subscription_id"`
	MetricName     string             `json:"metric_name" bson:"metric_name"` // api_calls, storage, ai_credits
	Quantity       int64              `json:"quantity" bson:"quantity"`
	Timestamp      time.Time          `json:"timestamp" bson:"timestamp"`
	Metadata       map[string]interface{} `json:"metadata,omitempty" bson:"metadata,omitempty"`
}
