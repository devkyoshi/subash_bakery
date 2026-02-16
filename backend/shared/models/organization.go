package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Organization with advanced multi-tenancy features
type Organization struct {
	BaseModel        `bson:",inline"`
	Name             string   `json:"name" bson:"name"`
	LegalName        string   `json:"legal_name" bson:"legal_name"`
	Domain           string   `json:"domain" bson:"domain"` // Primary domain
	AlternateDomains []string `json:"alternate_domains,omitempty" bson:"alternate_domains,omitempty"`
	Logo             string   `json:"logo" bson:"logo"`
	Favicon          string   `json:"favicon,omitempty" bson:"favicon,omitempty"`

	// Contact Information
	Email   string `json:"email" bson:"email"`
	Phone   string `json:"phone" bson:"phone"`
	Website string `json:"website" bson:"website"`

	// Business Information
	TaxID              string `json:"tax_id" bson:"tax_id"`
	RegistrationNumber string `json:"registration_number" bson:"registration_number"`
	Industry           string `json:"industry" bson:"industry"`
	CompanySize        string `json:"company_size" bson:"company_size"` // small, medium, large, enterprise

	// Status
	Status           OrganizationStatus `json:"status" bson:"status"`
	IsActive         bool               `json:"is_active" bson:"is_active"`
	ActivatedAt      *time.Time         `json:"activated_at,omitempty" bson:"activated_at,omitempty"`
	SuspendedAt      *time.Time         `json:"suspended_at,omitempty" bson:"suspended_at,omitempty"`
	SuspensionReason string             `json:"suspension_reason,omitempty" bson:"suspension_reason,omitempty"`

	// Subscription & Billing
	SubscriptionID primitive.ObjectID `json:"subscription_id,omitempty" bson:"subscription_id,omitempty"`
	BillingEmail   string             `json:"billing_email" bson:"billing_email"`
	BillingAddress *Address           `json:"billing_address,omitempty" bson:"billing_address,omitempty"`

	// Limits & Usage
	MaxUsers         int     `json:"max_users" bson:"max_users"`
	MaxCompanies     int     `json:"max_companies" bson:"max_companies"`
	MaxLocations     int     `json:"max_locations" bson:"max_locations"`
	CurrentUsers     int     `json:"current_users" bson:"current_users"`
	CurrentCompanies int     `json:"current_companies" bson:"current_companies"`
	CurrentLocations int     `json:"current_locations" bson:"current_locations"`
	StorageUsedGB    float64 `json:"storage_used_gb" bson:"storage_used_gb"`
	StorageLimitGB   float64 `json:"storage_limit_gb" bson:"storage_limit_gb"`

	// Settings
	Settings OrganizationSettings `json:"settings" bson:"settings"`

	// Branding
	BrandColors *BrandColors `json:"brand_colors,omitempty" bson:"brand_colors,omitempty"`

	// Metadata
	Metadata map[string]interface{} `json:"metadata,omitempty" bson:"metadata,omitempty"`
	Tags     []string               `json:"tags,omitempty" bson:"tags,omitempty"`
}

type OrganizationStatus string

const (
	OrganizationStatusActive    OrganizationStatus = "active"
	OrganizationStatusInactive  OrganizationStatus = "inactive"
	OrganizationStatusSuspended OrganizationStatus = "suspended"
	OrganizationStatusTrial     OrganizationStatus = "trial"
	OrganizationStatusCancelled OrganizationStatus = "cancelled"
)

type OrganizationSettings struct {
	Timezone                 string         `json:"timezone" bson:"timezone"`
	DateFormat               string         `json:"date_format" bson:"date_format"`
	TimeFormat               string         `json:"time_format" bson:"time_format"`
	Currency                 string         `json:"currency" bson:"currency"`
	Language                 string         `json:"language" bson:"language"`
	EnabledModules           []string       `json:"enabled_modules" bson:"enabled_modules"`
	AllowUserRegistration    bool           `json:"allow_user_registration" bson:"allow_user_registration"`
	RequireEmailVerification bool           `json:"require_email_verification" bson:"require_email_verification"`
	EnableMFA                bool           `json:"enable_mfa" bson:"enable_mfa"`
	SessionTimeout           int            `json:"session_timeout" bson:"session_timeout"` // minutes
	PasswordPolicy           PasswordPolicy `json:"password_policy" bson:"password_policy"`
}

type PasswordPolicy struct {
	MinLength           int  `json:"min_length" bson:"min_length"`
	RequireUppercase    bool `json:"require_uppercase" bson:"require_uppercase"`
	RequireLowercase    bool `json:"require_lowercase" bson:"require_lowercase"`
	RequireNumbers      bool `json:"require_numbers" bson:"require_numbers"`
	RequireSpecialChars bool `json:"require_special_chars" bson:"require_special_chars"`
	ExpiryDays          int  `json:"expiry_days" bson:"expiry_days"` // 0 = never expires
	PreventReuseCount   int  `json:"prevent_reuse_count" bson:"prevent_reuse_count"`
}

type BrandColors struct {
	Primary   string `json:"primary" bson:"primary"`
	Secondary string `json:"secondary" bson:"secondary"`
	Accent    string `json:"accent" bson:"accent"`
	Success   string `json:"success" bson:"success"`
	Warning   string `json:"warning" bson:"warning"`
	Error     string `json:"error" bson:"error"`
}

// OrganizationOption for lightweight dropdowns
type OrganizationOption struct {
	ID   primitive.ObjectID `json:"id" bson:"_id"`
	Name string             `json:"name" bson:"name"`
}

// Company with enhanced features
type Company struct {
	BaseModel      `bson:",inline"`
	OrganizationID primitive.ObjectID `json:"organization_id" bson:"organization_id"`
	Name           string             `json:"name" bson:"name"`
	LegalName      string             `json:"legal_name" bson:"legal_name"`
	Code           string             `json:"code" bson:"code"` // Unique within organization

	// Business Information
	TaxID              string `json:"tax_id" bson:"tax_id"`
	RegistrationNumber string `json:"registration_number" bson:"registration_number"`
	VATNumber          string `json:"vat_number,omitempty" bson:"vat_number,omitempty"`

	// Contact
	Email   string `json:"email" bson:"email"`
	Phone   string `json:"phone" bson:"phone"`
	Fax     string `json:"fax,omitempty" bson:"fax,omitempty"`
	Website string `json:"website,omitempty" bson:"website,omitempty"`

	// Address
	Address Address `json:"address" bson:"address"`

	// Banking
	BankAccounts []BankAccount `json:"bank_accounts,omitempty" bson:"bank_accounts,omitempty"`

	// Settings
	Settings CompanySettings `json:"settings" bson:"settings"`

	// Status
	IsActive  bool `json:"is_active" bson:"is_active"`
	IsDefault bool `json:"is_default" bson:"is_default"` // Default company

	// Parent company (for company groups)
	ParentCompanyID *primitive.ObjectID `json:"parent_company_id,omitempty" bson:"parent_company_id,omitempty"`

	// Metadata
	Metadata map[string]interface{} `json:"metadata,omitempty" bson:"metadata,omitempty"`
	Tags     []string               `json:"tags,omitempty" bson:"tags,omitempty"`
}

type CompanySettings struct {
	FiscalYearStart     string `json:"fiscal_year_start" bson:"fiscal_year_start"` // MM-DD
	Currency            string `json:"currency" bson:"currency"`
	Timezone            string `json:"timezone" bson:"timezone"`
	EnableMultiCurrency bool   `json:"enable_multi_currency" bson:"enable_multi_currency"`
}

type BankAccount struct {
	BankName      string `json:"bank_name" bson:"bank_name"`
	AccountNumber string `json:"account_number" bson:"account_number"`
	AccountName   string `json:"account_name" bson:"account_name"`
	IBAN          string `json:"iban,omitempty" bson:"iban,omitempty"`
	SwiftCode     string `json:"swift_code,omitempty" bson:"swift_code,omitempty"`
	Branch        string `json:"branch,omitempty" bson:"branch,omitempty"`
	IsDefault     bool   `json:"is_default" bson:"is_default"`
}

// Location with advanced features
type Location struct {
	BaseModel      `bson:",inline"`
	CompanyID      primitive.ObjectID `json:"company_id" bson:"company_id"`
	OrganizationID primitive.ObjectID `json:"organization_id" bson:"organization_id"` // Denormalized for queries
	Name           string             `json:"name" bson:"name"`
	Code           string             `json:"code" bson:"code"` // Unique within company

	// Type & Category
	Type     LocationType `json:"type" bson:"type"`
	Category string       `json:"category,omitempty" bson:"category,omitempty"`

	// Contact
	Email string `json:"email" bson:"email"`
	Phone string `json:"phone" bson:"phone"`
	Fax   string `json:"fax,omitempty" bson:"fax,omitempty"`

	// Address
	Address Address `json:"address" bson:"address"`

	// Manager
	ManagerID *primitive.ObjectID `json:"manager_id,omitempty" bson:"manager_id,omitempty"`

	// Applications & Licenses
	Applications []ApplicationAssignment `json:"applications" bson:"applications"`

	// Warehouse specific (if type is warehouse)
	WarehouseInfo *WarehouseInfo `json:"warehouse_info,omitempty" bson:"warehouse_info,omitempty"`

	// Store specific (if type is store)
	StoreInfo *StoreInfo `json:"store_info,omitempty" bson:"store_info,omitempty"`

	// Operating hours
	OperatingHours []OperatingHours `json:"operating_hours,omitempty" bson:"operating_hours,omitempty"`

	// Settings
	Settings LocationSettings `json:"settings" bson:"settings"`

	// Status
	IsActive  bool `json:"is_active" bson:"is_active"`
	IsDefault bool `json:"is_default" bson:"is_default"`

	// Parent location (for location hierarchy)
	ParentLocationID *primitive.ObjectID `json:"parent_location_id,omitempty" bson:"parent_location_id,omitempty"`

	// Metadata
	Metadata map[string]interface{} `json:"metadata,omitempty" bson:"metadata,omitempty"`
	Tags     []string               `json:"tags,omitempty" bson:"tags,omitempty"`
}

type LocationType string

const (
	LocationTypeHeadOffice         LocationType = "head_office"
	LocationTypeBranch             LocationType = "branch"
	LocationTypeWarehouse          LocationType = "warehouse"
	LocationTypeStore              LocationType = "store"
	LocationTypeFactory            LocationType = "factory"
	LocationTypeDistributionCenter LocationType = "distribution_center"
)

type ApplicationAssignment struct {
	ApplicationID primitive.ObjectID `json:"application_id" bson:"application_id"`
	LicenseID     primitive.ObjectID `json:"license_id" bson:"license_id"`
	AssignedAt    time.Time          `json:"assigned_at" bson:"assigned_at"`
	IsActive      bool               `json:"is_active" bson:"is_active"`
}

type WarehouseInfo struct {
	TotalArea        float64 `json:"total_area" bson:"total_area"`             // in sq meters
	StorageCapacity  int     `json:"storage_capacity" bson:"storage_capacity"` // number of pallets
	DockingBays      int     `json:"docking_bays" bson:"docking_bays"`
	RefrigeratedArea float64 `json:"refrigerated_area,omitempty" bson:"refrigerated_area,omitempty"`
	HasColdStorage   bool    `json:"has_cold_storage" bson:"has_cold_storage"`
}

type StoreInfo struct {
	FloorArea         float64 `json:"floor_area" bson:"floor_area"` // in sq meters
	POSCount          int     `json:"pos_count" bson:"pos_count"`
	ParkingSpaces     int     `json:"parking_spaces,omitempty" bson:"parking_spaces,omitempty"`
	HasOnlineOrdering bool    `json:"has_online_ordering" bson:"has_online_ordering"`
}

type OperatingHours struct {
	DayOfWeek int    `json:"day_of_week" bson:"day_of_week"` // 0 = Sunday, 1 = Monday, etc.
	OpenTime  string `json:"open_time" bson:"open_time"`     // HH:MM format
	CloseTime string `json:"close_time" bson:"close_time"`   // HH:MM format
	IsClosed  bool   `json:"is_closed" bson:"is_closed"`
}

type LocationSettings struct {
	Timezone                   string `json:"timezone" bson:"timezone"`
	AllowBackdatedTransactions bool   `json:"allow_backdated_transactions" bson:"allow_backdated_transactions"`
	RequireApproval            bool   `json:"require_approval" bson:"require_approval"`
}

type Address struct {
	Street      string  `json:"street" bson:"street"`
	Street2     string  `json:"street2,omitempty" bson:"street2,omitempty"`
	City        string  `json:"city" bson:"city"`
	State       string  `json:"state" bson:"state"`
	PostalCode  string  `json:"postal_code" bson:"postal_code"`
	Country     string  `json:"country" bson:"country"`
	CountryCode string  `json:"country_code,omitempty" bson:"country_code,omitempty"` // ISO 3166-1 alpha-2
	Latitude    float64 `json:"latitude,omitempty" bson:"latitude,omitempty"`
	Longitude   float64 `json:"longitude,omitempty" bson:"longitude,omitempty"`
}

// LocationUser with role and permissions
type LocationUser struct {
	BaseModel   `bson:",inline"`
	LocationID  primitive.ObjectID `json:"location_id" bson:"location_id"`
	UserID      primitive.ObjectID `json:"user_id" bson:"user_id"`
	RoleID      primitive.ObjectID `json:"role_id" bson:"role_id"` // Location-specific role
	IsPrimary   bool               `json:"is_primary" bson:"is_primary"`
	IsActive    bool               `json:"is_active" bson:"is_active"`
	AssignedAt  time.Time          `json:"assigned_at" bson:"assigned_at"`
	AccessLevel string             `json:"access_level" bson:"access_level"` // full, read_only, restricted
}
