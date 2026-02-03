package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// BaseModel contains common fields for all models
type BaseModel struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
	DeletedAt *time.Time         `json:"deleted_at,omitempty" bson:"deleted_at,omitempty"` // Soft delete
	CreatedBy primitive.ObjectID `json:"created_by,omitempty" bson:"created_by,omitempty"`
	UpdatedBy primitive.ObjectID `json:"updated_by,omitempty" bson:"updated_by,omitempty"`
	DeletedBy primitive.ObjectID `json:"deleted_by,omitempty" bson:"deleted_by,omitempty"`
	Version   int                `json:"version" bson:"version"` // Optimistic locking
}

// User represents a system user with advanced features
type User struct {
	BaseModel          `bson:",inline"`
	OrganizationID     primitive.ObjectID `json:"organization_id" bson:"organization_id"`
	Email              string             `json:"email" bson:"email"`
	Password           string             `json:"-" bson:"password"`
	FirstName          string             `json:"first_name" bson:"first_name"`
	LastName           string             `json:"last_name" bson:"last_name"`
	FullName           string             `json:"full_name" bson:"full_name"` // Computed field
	Phone              string             `json:"phone" bson:"phone"`
	Avatar             string             `json:"avatar" bson:"avatar"`
	RoleID             primitive.ObjectID `json:"role_id" bson:"role_id"`

	// Status and verification
	Status             UserStatus         `json:"status" bson:"status"`
	IsActive           bool               `json:"is_active" bson:"is_active"`
	IsEmailVerified    bool               `json:"is_email_verified" bson:"is_email_verified"`
	IsPhoneVerified    bool               `json:"is_phone_verified" bson:"is_phone_verified"`
	EmailVerifiedAt    *time.Time         `json:"email_verified_at,omitempty" bson:"email_verified_at,omitempty"`
	PhoneVerifiedAt    *time.Time         `json:"phone_verified_at,omitempty" bson:"phone_verified_at,omitempty"`

	// OAuth integrations
	GoogleID           string             `json:"google_id,omitempty" bson:"google_id,omitempty"`
	MicrosoftID        string             `json:"microsoft_id,omitempty" bson:"microsoft_id,omitempty"`

	// Security
	MFAEnabled         bool               `json:"mfa_enabled" bson:"mfa_enabled"`
	MFASecret          string             `json:"-" bson:"mfa_secret,omitempty"`
	PasswordChangedAt  *time.Time         `json:"password_changed_at,omitempty" bson:"password_changed_at,omitempty"`
	FailedLoginAttempts int               `json:"failed_login_attempts" bson:"failed_login_attempts"`
	LastFailedLoginAt  *time.Time         `json:"last_failed_login_at,omitempty" bson:"last_failed_login_at,omitempty"`
	LockedUntil        *time.Time         `json:"locked_until,omitempty" bson:"locked_until,omitempty"`

	// Activity tracking
	LastLogin          time.Time          `json:"last_login" bson:"last_login"`
	LastLoginIP        string             `json:"last_login_ip" bson:"last_login_ip"`
	LastActivity       time.Time          `json:"last_activity" bson:"last_activity"`

	// Preferences
	Timezone           string             `json:"timezone" bson:"timezone"`
	Language           string             `json:"language" bson:"language"`
	DateFormat         string             `json:"date_format" bson:"date_format"`
	TimeFormat         string             `json:"time_format" bson:"time_format"`

	// Metadata
	Metadata           map[string]interface{} `json:"metadata,omitempty" bson:"metadata,omitempty"`
	Tags               []string           `json:"tags,omitempty" bson:"tags,omitempty"`
}

type UserStatus string

const (
	UserStatusActive    UserStatus = "active"
	UserStatusInactive  UserStatus = "inactive"
	UserStatusSuspended UserStatus = "suspended"
	UserStatusPending   UserStatus = "pending"
	UserStatusLocked    UserStatus = "locked"
)

// Role with hierarchical permissions
type Role struct {
	BaseModel      `bson:",inline"`
	OrganizationID primitive.ObjectID   `json:"organization_id" bson:"organization_id"`
	Name           string               `json:"name" bson:"name"`
	DisplayName    string               `json:"display_name" bson:"display_name"`
	Description    string               `json:"description" bson:"description"`
	Permissions    []primitive.ObjectID `json:"permissions" bson:"permissions"`
	IsSystem       bool                 `json:"is_system" bson:"is_system"`
	IsDefault      bool                 `json:"is_default" bson:"is_default"` // Default role for new users
	Priority       int                  `json:"priority" bson:"priority"`     // Higher priority = more access
	ParentRoleID   *primitive.ObjectID  `json:"parent_role_id,omitempty" bson:"parent_role_id,omitempty"` // Role inheritance
	IsActive       bool                 `json:"is_active" bson:"is_active"`
	MaxUsers       int                  `json:"max_users,omitempty" bson:"max_users,omitempty"` // Limit users with this role
	Metadata       map[string]interface{} `json:"metadata,omitempty" bson:"metadata,omitempty"`
}

// Permission with granular control
type Permission struct {
	BaseModel   `bson:",inline"`
	Name        string       `json:"name" bson:"name"`
	DisplayName string       `json:"display_name" bson:"display_name"`
	Description string       `json:"description" bson:"description"`
	Resource    string       `json:"resource" bson:"resource"`
	Action      string       `json:"action" bson:"action"`
	Scope       PermissionScope `json:"scope" bson:"scope"`
	Conditions  []PermissionCondition `json:"conditions,omitempty" bson:"conditions,omitempty"`
	IsSystem    bool         `json:"is_system" bson:"is_system"`
	Category    string       `json:"category" bson:"category"`
	Metadata    map[string]interface{} `json:"metadata,omitempty" bson:"metadata,omitempty"`
}

type PermissionScope string

const (
	PermissionScopeGlobal       PermissionScope = "global"        // All organizations
	PermissionScopeOrganization PermissionScope = "organization"  // Within organization
	PermissionScopeCompany      PermissionScope = "company"       // Within company
	PermissionScopeLocation     PermissionScope = "location"      // Within location
	PermissionScopeOwn          PermissionScope = "own"           // Only own resources
)

type PermissionCondition struct {
	Field    string      `json:"field" bson:"field"`
	Operator string      `json:"operator" bson:"operator"` // eq, ne, gt, lt, in, contains
	Value    interface{} `json:"value" bson:"value"`
}

// Session with enhanced security and token rotation
type Session struct {
	BaseModel    `bson:",inline"`
	UserID       primitive.ObjectID `json:"user_id" bson:"user_id"`
	RefreshToken string             `json:"refresh_token" bson:"refresh_token"`
	TokenFamily  string             `json:"token_family" bson:"token_family"` // For token rotation detection
	UserAgent    string             `json:"user_agent" bson:"user_agent"`
	IPAddress    string             `json:"ip_address" bson:"ip_address"`
	Device       DeviceInfo         `json:"device" bson:"device"`
	Location     *GeoLocation       `json:"location,omitempty" bson:"location,omitempty"`
	ExpiresAt    time.Time          `json:"expires_at" bson:"expires_at"`
	LastUsedAt   time.Time          `json:"last_used_at" bson:"last_used_at"`
	IsRevoked    bool               `json:"is_revoked" bson:"is_revoked"`
	RevokedAt    *time.Time         `json:"revoked_at,omitempty" bson:"revoked_at,omitempty"`
	RevokedBy    *primitive.ObjectID `json:"revoked_by,omitempty" bson:"revoked_by,omitempty"`
	RevokeReason string             `json:"revoke_reason,omitempty" bson:"revoke_reason,omitempty"`
}

type DeviceInfo struct {
	Type        string `json:"type" bson:"type"`                 // mobile, desktop, tablet
	OS          string `json:"os" bson:"os"`
	Browser     string `json:"browser" bson:"browser"`
	DeviceID    string `json:"device_id,omitempty" bson:"device_id,omitempty"`
	IsTrusted   bool   `json:"is_trusted" bson:"is_trusted"`
}

type GeoLocation struct {
	Country     string  `json:"country" bson:"country"`
	City        string  `json:"city" bson:"city"`
	Latitude    float64 `json:"latitude" bson:"latitude"`
	Longitude   float64 `json:"longitude" bson:"longitude"`
}

// APIKey with rate limiting and scopes
type APIKey struct {
	BaseModel      `bson:",inline"`
	OrganizationID primitive.ObjectID `json:"organization_id" bson:"organization_id"`
	UserID         primitive.ObjectID `json:"user_id" bson:"user_id"`
	Name           string             `json:"name" bson:"name"`
	Description    string             `json:"description" bson:"description"`
	Key            string             `json:"key" bson:"key"`
	KeyPrefix      string             `json:"key_prefix" bson:"key_prefix"` // First 8 chars for display
	Scopes         []string           `json:"scopes" bson:"scopes"`         // Permissions for this key
	RateLimit      int                `json:"rate_limit" bson:"rate_limit"` // Requests per hour
	IsActive       bool               `json:"is_active" bson:"is_active"`
	ExpiresAt      *time.Time         `json:"expires_at,omitempty" bson:"expires_at,omitempty"`
	LastUsedAt     *time.Time         `json:"last_used_at,omitempty" bson:"last_used_at,omitempty"`
	UsageCount     int64              `json:"usage_count" bson:"usage_count"`
	AllowedIPs     []string           `json:"allowed_ips,omitempty" bson:"allowed_ips,omitempty"`
	Metadata       map[string]interface{} `json:"metadata,omitempty" bson:"metadata,omitempty"`
}

// UserClaims for JWT
type UserClaims struct {
	UserID         string   `json:"user_id"`
	OrganizationID string   `json:"organization_id"`
	Email          string   `json:"email"`
	RoleID         string   `json:"role_id"`
	Permissions    []string `json:"permissions,omitempty"`
	Scopes         []string `json:"scopes,omitempty"`
}

// AuditLog for tracking all changes
type AuditLog struct {
	ID             primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	OrganizationID primitive.ObjectID `json:"organization_id" bson:"organization_id"`
	UserID         primitive.ObjectID `json:"user_id" bson:"user_id"`
	Action         string             `json:"action" bson:"action"`         // create, update, delete, login, etc.
	Resource       string             `json:"resource" bson:"resource"`     // users, products, orders, etc.
	ResourceID     primitive.ObjectID `json:"resource_id" bson:"resource_id"`
	Changes        map[string]interface{} `json:"changes,omitempty" bson:"changes,omitempty"`
	IPAddress      string             `json:"ip_address" bson:"ip_address"`
	UserAgent      string             `json:"user_agent" bson:"user_agent"`
	Status         string             `json:"status" bson:"status"`         // success, failure
	ErrorMessage   string             `json:"error_message,omitempty" bson:"error_message,omitempty"`
	Metadata       map[string]interface{} `json:"metadata,omitempty" bson:"metadata,omitempty"`
	Timestamp      time.Time          `json:"timestamp" bson:"timestamp"`
}
