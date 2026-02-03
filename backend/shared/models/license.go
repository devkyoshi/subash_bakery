package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// LicenseType defines the type of license
type LicenseType string

const (
	LicenseTypeUserBased   LicenseType = "user_based"   // License per user
	LicenseTypeDeviceBased LicenseType = "device_based" // License per device
	LicenseTypeUsageBased  LicenseType = "usage_based"  // License based on usage metrics
	LicenseTypeConcurrent  LicenseType = "concurrent"   // Concurrent user license
	LicenseTypeUnlimited   LicenseType = "unlimited"    // Unlimited license
)

// ApplicationCategory defines application categories
type ApplicationCategory string

const (
	CategoryInventory     ApplicationCategory = "inventory"
	CategoryProcurement   ApplicationCategory = "procurement"
	CategorySales         ApplicationCategory = "sales"
	CategoryFinance       ApplicationCategory = "finance"
	CategoryCRM           ApplicationCategory = "crm"
	CategoryHR            ApplicationCategory = "hr"
	CategoryManufacturing ApplicationCategory = "manufacturing"
	CategoryWorkflow      ApplicationCategory = "workflow"
	CategoryAI            ApplicationCategory = "ai"
	CategoryAnalytics     ApplicationCategory = "analytics"
)

// LicenseStatus defines license status
type LicenseStatus string

const (
	LicenseStatusActive    LicenseStatus = "active"
	LicenseStatusSuspended LicenseStatus = "suspended"
	LicenseStatusExpired   LicenseStatus = "expired"
	LicenseStatusExceeded  LicenseStatus = "exceeded" // Usage exceeded
	LicenseStatusRevoked   LicenseStatus = "revoked"
)

// Application represents an ERP application/module
type Application struct {
	BaseModel `bson:",inline"`

	// Basic Information
	Name        string              `bson:"name" json:"name" binding:"required"`
	Code        string              `bson:"code" json:"code" binding:"required"` // Unique identifier (e.g., "INV", "CRM")
	DisplayName string              `bson:"display_name" json:"display_name"`
	Description string              `bson:"description" json:"description"`
	Category    ApplicationCategory `bson:"category" json:"category"`
	Version     string              `bson:"version" json:"version"`

	// Licensing Configuration
	SupportedLicenseTypes []LicenseType `bson:"supported_license_types" json:"supported_license_types"`
	DefaultLicenseType    LicenseType   `bson:"default_license_type" json:"default_license_type"`

	// Pricing
	BasePrice            float64 `bson:"base_price" json:"base_price"`                         // Base monthly price
	PricePerUser         float64 `bson:"price_per_user" json:"price_per_user"`                 // Additional price per user
	PricePerDevice       float64 `bson:"price_per_device" json:"price_per_device"`             // Additional price per device
	PricePerTransaction  float64 `bson:"price_per_transaction" json:"price_per_transaction"`   // For usage-based
	MinimumUsers         int     `bson:"minimum_users" json:"minimum_users"`                   // Minimum users to purchase
	IncludedTransactions int64   `bson:"included_transactions" json:"included_transactions"`   // Included in base price

	// Features & Dependencies
	Features             []ApplicationFeature   `bson:"features" json:"features"`
	RequiredApplications []primitive.ObjectID   `bson:"required_applications" json:"required_applications"` // Dependencies
	IntegratesWith       []primitive.ObjectID   `bson:"integrates_with" json:"integrates_with"`

	// API & Technical
	APIEndpoint    string            `bson:"api_endpoint" json:"api_endpoint"`
	ServiceURL     string            `bson:"service_url" json:"service_url"`
	HealthCheckURL string            `bson:"health_check_url" json:"health_check_url"`
	Configuration  map[string]string `bson:"configuration" json:"configuration"`

	// Status & Availability
	IsActive    bool                   `bson:"is_active" json:"is_active"`
	IsPublic    bool                   `bson:"is_public" json:"is_public"` // Available for purchase
	IsBeta      bool                   `bson:"is_beta" json:"is_beta"`
	ReleaseDate *time.Time             `bson:"release_date" json:"release_date"`
	EOLDate     *time.Time             `bson:"eol_date" json:"eol_date"` // End of life

	// Documentation
	DocumentationURL string   `bson:"documentation_url" json:"documentation_url"`
	SupportEmail     string   `bson:"support_email" json:"support_email"`
	Icon             string   `bson:"icon" json:"icon"`
	Screenshots      []string `bson:"screenshots" json:"screenshots"`

	// Analytics
	TotalLicenses  int `bson:"total_licenses" json:"total_licenses"`   // Total active licenses
	TotalLocations int `bson:"total_locations" json:"total_locations"` // Locations using this app

	// Metadata
	Metadata map[string]interface{} `bson:"metadata" json:"metadata"`
	Tags     []string               `bson:"tags" json:"tags"`
}

// ApplicationFeature represents a feature of an application
type ApplicationFeature struct {
	Name        string `bson:"name" json:"name"`
	Description string `bson:"description" json:"description"`
	IsAdvanced  bool   `bson:"is_advanced" json:"is_advanced"` // Requires higher tier
	IsAIPowered bool   `bson:"is_ai_powered" json:"is_ai_powered"`
}

// LocationLicense represents a license assigned to a location
type LocationLicense struct {
	BaseModel `bson:",inline"`

	// Ownership
	OrganizationID primitive.ObjectID `bson:"organization_id" json:"organization_id" binding:"required"`
	LocationID     primitive.ObjectID `bson:"location_id" json:"location_id" binding:"required"`
	ApplicationID  primitive.ObjectID `bson:"application_id" json:"application_id" binding:"required"`

	// License Configuration
	LicenseType   LicenseType   `bson:"license_type" json:"license_type" binding:"required"`
	Status        LicenseStatus `bson:"status" json:"status"`
	LicenseKey    string        `bson:"license_key" json:"license_key"`         // Unique license key
	ActivationKey string        `bson:"activation_key" json:"activation_key"`   // For device activation

	// Limits
	MaxUsers         int   `bson:"max_users" json:"max_users"`                   // For user-based
	MaxDevices       int   `bson:"max_devices" json:"max_devices"`               // For device-based
	MaxTransactions  int64 `bson:"max_transactions" json:"max_transactions"`     // For usage-based
	MaxConcurrent    int   `bson:"max_concurrent" json:"max_concurrent"`         // For concurrent
	MaxStorageGB     float64 `bson:"max_storage_gb" json:"max_storage_gb"`
	MaxAPICallsDaily int64 `bson:"max_api_calls_daily" json:"max_api_calls_daily"`

	// Current Usage
	CurrentUsers        int     `bson:"current_users" json:"current_users"`
	CurrentDevices      int     `bson:"current_devices" json:"current_devices"`
	CurrentTransactions int64   `bson:"current_transactions" json:"current_transactions"`
	CurrentStorageGB    float64 `bson:"current_storage_gb" json:"current_storage_gb"`
	CurrentAPICallsToday int64  `bson:"current_api_calls_today" json:"current_api_calls_today"`

	// Period
	StartDate     time.Time  `bson:"start_date" json:"start_date"`
	EndDate       *time.Time `bson:"end_date" json:"end_date"`           // nil for perpetual
	TrialEndDate  *time.Time `bson:"trial_end_date" json:"trial_end_date"`
	LastRenewedAt *time.Time `bson:"last_renewed_at" json:"last_renewed_at"`
	IsTrial       bool       `bson:"is_trial" json:"is_trial"`
	IsPerpetual   bool       `bson:"is_perpetual" json:"is_perpetual"` // Never expires

	// Billing
	BillingCycle  BillingCycle `bson:"billing_cycle" json:"billing_cycle"`
	PricePerCycle float64      `bson:"price_per_cycle" json:"price_per_cycle"`
	Currency      string       `bson:"currency" json:"currency"`
	AutoRenew     bool         `bson:"auto_renew" json:"auto_renew"`

	// Features
	EnabledFeatures  []string `bson:"enabled_features" json:"enabled_features"`
	DisabledFeatures []string `bson:"disabled_features" json:"disabled_features"`

	// Restrictions
	AllowedIPs       []string  `bson:"allowed_ips" json:"allowed_ips"`           // IP whitelist
	AllowedDomains   []string  `bson:"allowed_domains" json:"allowed_domains"`   // Email domain restrictions
	GeoRestrictions  []string  `bson:"geo_restrictions" json:"geo_restrictions"` // Country codes
	RestrictionsNote string    `bson:"restrictions_note" json:"restrictions_note"`

	// Compliance & Audit
	ComplianceLevel     string    `bson:"compliance_level" json:"compliance_level"` // e.g., "HIPAA", "SOC2"
	LastAuditDate       *time.Time `bson:"last_audit_date" json:"last_audit_date"`
	CertificationExpiry *time.Time `bson:"certification_expiry" json:"certification_expiry"`

	// Status tracking
	ActivatedAt   *time.Time `bson:"activated_at" json:"activated_at"`
	SuspendedAt   *time.Time `bson:"suspended_at" json:"suspended_at"`
	SuspendReason string     `bson:"suspend_reason" json:"suspend_reason"`
	RevokedAt     *time.Time `bson:"revoked_at" json:"revoked_at"`
	RevokeReason  string     `bson:"revoke_reason" json:"revoke_reason"`
	ExpiresAt     *time.Time `bson:"expires_at" json:"expires_at"`

	// Notifications
	NotifyBeforeExpiry int  `bson:"notify_before_expiry" json:"notify_before_expiry"` // Days
	NotifyOnOveruse    bool `bson:"notify_on_overuse" json:"notify_on_overuse"`
	NotifyAdmins       []primitive.ObjectID `bson:"notify_admins" json:"notify_admins"`

	// Metadata
	PurchaseOrderID string                 `bson:"purchase_order_id" json:"purchase_order_id"`
	Notes           string                 `bson:"notes" json:"notes"`
	Metadata        map[string]interface{} `bson:"metadata" json:"metadata"`
	Tags            []string               `bson:"tags" json:"tags"`
}

// UserLicenseAssignment tracks which users are assigned to which licenses
type UserLicenseAssignment struct {
	BaseModel `bson:",inline"`

	LicenseID      primitive.ObjectID `bson:"license_id" json:"license_id" binding:"required"`
	UserID         primitive.ObjectID `bson:"user_id" json:"user_id" binding:"required"`
	LocationID     primitive.ObjectID `bson:"location_id" json:"location_id" binding:"required"`
	ApplicationID  primitive.ObjectID `bson:"application_id" json:"application_id" binding:"required"`
	OrganizationID primitive.ObjectID `bson:"organization_id" json:"organization_id"`

	// Assignment Details
	Role           string     `bson:"role" json:"role"`                     // User's role in the app
	Permissions    []string   `bson:"permissions" json:"permissions"`       // Specific permissions
	AssignedAt     time.Time  `bson:"assigned_at" json:"assigned_at"`
	LastAccessedAt *time.Time `bson:"last_accessed_at" json:"last_accessed_at"`
	AccessCount    int64      `bson:"access_count" json:"access_count"`

	// Status
	IsActive   bool       `bson:"is_active" json:"is_active"`
	RevokedAt  *time.Time `bson:"revoked_at" json:"revoked_at"`
	RevokeBy   primitive.ObjectID `bson:"revoke_by" json:"revoke_by"`
	RevokeReason string   `bson:"revoke_reason" json:"revoke_reason"`

	Metadata map[string]interface{} `bson:"metadata" json:"metadata"`
}

// DeviceLicenseAssignment tracks which devices are activated with licenses
type DeviceLicenseAssignment struct {
	BaseModel `bson:",inline"`

	LicenseID      primitive.ObjectID `bson:"license_id" json:"license_id" binding:"required"`
	LocationID     primitive.ObjectID `bson:"location_id" json:"location_id" binding:"required"`
	ApplicationID  primitive.ObjectID `bson:"application_id" json:"application_id" binding:"required"`
	OrganizationID primitive.ObjectID `bson:"organization_id" json:"organization_id"`

	// Device Information
	DeviceID       string     `bson:"device_id" json:"device_id" binding:"required"` // Unique device identifier
	DeviceName     string     `bson:"device_name" json:"device_name"`
	DeviceType     string     `bson:"device_type" json:"device_type"`     // desktop, mobile, tablet, pos
	DeviceModel    string     `bson:"device_model" json:"device_model"`
	OS             string     `bson:"os" json:"os"`
	OSVersion      string     `bson:"os_version" json:"os_version"`
	MACAddress     string     `bson:"mac_address" json:"mac_address"`
	IPAddress      string     `bson:"ip_address" json:"ip_address"`
	Hostname       string     `bson:"hostname" json:"hostname"`

	// Activation
	ActivationCode string     `bson:"activation_code" json:"activation_code"`
	ActivatedAt    time.Time  `bson:"activated_at" json:"activated_at"`
	LastSeenAt     *time.Time `bson:"last_seen_at" json:"last_seen_at"`
	LastSyncAt     *time.Time `bson:"last_sync_at" json:"last_sync_at"`

	// Status
	IsActive       bool       `bson:"is_active" json:"is_active"`
	IsOnline       bool       `bson:"is_online" json:"is_online"`
	DeactivatedAt  *time.Time `bson:"deactivated_at" json:"deactivated_at"`
	DeactivateBy   primitive.ObjectID `bson:"deactivate_by" json:"deactivate_by"`
	DeactivateReason string   `bson:"deactivate_reason" json:"deactivate_reason"`

	// User Assignment (optional - device can be shared)
	AssignedUserID *primitive.ObjectID `bson:"assigned_user_id" json:"assigned_user_id"`

	Metadata map[string]interface{} `bson:"metadata" json:"metadata"`
}

// LicenseUsageLog tracks usage for usage-based licenses
type LicenseUsageLog struct {
	BaseModel `bson:",inline"`

	LicenseID      primitive.ObjectID `bson:"license_id" json:"license_id" binding:"required"`
	ApplicationID  primitive.ObjectID `bson:"application_id" json:"application_id" binding:"required"`
	LocationID     primitive.ObjectID `bson:"location_id" json:"location_id" binding:"required"`
	OrganizationID primitive.ObjectID `bson:"organization_id" json:"organization_id"`

	// Usage Period
	PeriodStart time.Time `bson:"period_start" json:"period_start"`
	PeriodEnd   time.Time `bson:"period_end" json:"period_end"`

	// Metrics
	TransactionCount int64   `bson:"transaction_count" json:"transaction_count"`
	APICallCount     int64   `bson:"api_call_count" json:"api_call_count"`
	StorageUsedGB    float64 `bson:"storage_used_gb" json:"storage_used_gb"`
	DataProcessedGB  float64 `bson:"data_processed_gb" json:"data_processed_gb"`
	ActiveUsers      int     `bson:"active_users" json:"active_users"`
	ActiveDevices    int     `bson:"active_devices" json:"active_devices"`
	PeakConcurrent   int     `bson:"peak_concurrent" json:"peak_concurrent"`

	// Costs (for usage-based billing)
	BaseCost        float64 `bson:"base_cost" json:"base_cost"`
	UsageCost       float64 `bson:"usage_cost" json:"usage_cost"`
	OverageCost     float64 `bson:"overage_cost" json:"overage_cost"`
	TotalCost       float64 `bson:"total_cost" json:"total_cost"`

	// Custom Metrics
	CustomMetrics map[string]interface{} `bson:"custom_metrics" json:"custom_metrics"`

	Metadata map[string]interface{} `bson:"metadata" json:"metadata"`
}

// LicenseAlert represents alerts for license usage or expiry
type LicenseAlert struct {
	BaseModel `bson:",inline"`

	LicenseID      primitive.ObjectID `bson:"license_id" json:"license_id" binding:"required"`
	OrganizationID primitive.ObjectID `bson:"organization_id" json:"organization_id"`
	LocationID     primitive.ObjectID `bson:"location_id" json:"location_id"`
	ApplicationID  primitive.ObjectID `bson:"application_id" json:"application_id"`

	// Alert Details
	AlertType    string    `bson:"alert_type" json:"alert_type"`       // expiry_warning, usage_exceeded, etc.
	Severity     string    `bson:"severity" json:"severity"`           // info, warning, critical
	Title        string    `bson:"title" json:"title"`
	Message      string    `bson:"message" json:"message"`
	TriggeredAt  time.Time `bson:"triggered_at" json:"triggered_at"`

	// Status
	IsRead         bool       `bson:"is_read" json:"is_read"`
	ReadAt         *time.Time `bson:"read_at" json:"read_at"`
	ReadBy         primitive.ObjectID `bson:"read_by" json:"read_by"`
	IsResolved     bool       `bson:"is_resolved" json:"is_resolved"`
	ResolvedAt     *time.Time `bson:"resolved_at" json:"resolved_at"`
	ResolvedBy     primitive.ObjectID `bson:"resolved_by" json:"resolved_by"`
	ResolutionNote string     `bson:"resolution_note" json:"resolution_note"`

	// Notification
	NotifiedUsers []primitive.ObjectID `bson:"notified_users" json:"notified_users"`
	NotifiedAt    *time.Time           `bson:"notified_at" json:"notified_at"`

	Metadata map[string]interface{} `bson:"metadata" json:"metadata"`
}
