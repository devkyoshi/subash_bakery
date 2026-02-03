package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/yourusername/erp-system/services/license-service/internal/repository"
	"github.com/yourusername/erp-system/shared/models"
)

type LicenseService struct {
	licenseRepo      *repository.LicenseRepository
	appRepo          *repository.ApplicationRepository
	userAssignRepo   *repository.UserAssignmentRepository
	deviceAssignRepo *repository.DeviceAssignmentRepository
}

func NewLicenseService(
	licenseRepo *repository.LicenseRepository,
	appRepo *repository.ApplicationRepository,
	userAssignRepo *repository.UserAssignmentRepository,
	deviceAssignRepo *repository.DeviceAssignmentRepository,
) *LicenseService {
	return &LicenseService{
		licenseRepo:      licenseRepo,
		appRepo:          appRepo,
		userAssignRepo:   userAssignRepo,
		deviceAssignRepo: deviceAssignRepo,
	}
}

// CreateLicenseRequest represents a request to create a license
type CreateLicenseRequest struct {
	ApplicationID     string                `json:"application_id" binding:"required"`
	LocationID        string                `json:"location_id" binding:"required"`
	LicenseType       models.LicenseType    `json:"license_type" binding:"required"`
	MaxUsers          int                   `json:"max_users"`
	MaxDevices        int                   `json:"max_devices"`
	MaxTransactions   int64                 `json:"max_transactions"`
	MaxConcurrent     int                   `json:"max_concurrent"`
	MaxStorageGB      float64               `json:"max_storage_gb"`
	MaxAPICallsDaily  int64                 `json:"max_api_calls_daily"`
	BillingCycle      models.BillingCycle   `json:"billing_cycle"`
	PricePerCycle     float64               `json:"price_per_cycle"`
	IsTrial           bool                  `json:"is_trial"`
	TrialDays         int                   `json:"trial_days"`
	IsPerpetual       bool                  `json:"is_perpetual"`
	DurationMonths    int                   `json:"duration_months"`
	EnabledFeatures   []string              `json:"enabled_features"`
	DisabledFeatures  []string              `json:"disabled_features"`
}

// CreateLicense creates a new license for a location
func (s *LicenseService) CreateLicense(ctx context.Context, orgID primitive.ObjectID, req CreateLicenseRequest, createdBy primitive.ObjectID) (*models.LocationLicense, error) {
	// Parse IDs
	appID, err := primitive.ObjectIDFromHex(req.ApplicationID)
	if err != nil {
		return nil, fmt.Errorf("invalid application ID: %w", err)
	}

	locationID, err := primitive.ObjectIDFromHex(req.LocationID)
	if err != nil {
		return nil, fmt.Errorf("invalid location ID: %w", err)
	}

	// Verify application exists
	app, err := s.appRepo.FindByID(ctx, appID)
	if err != nil {
		return nil, fmt.Errorf("application not found: %w", err)
	}

	if !app.IsActive {
		return nil, fmt.Errorf("application is not active")
	}

	// Check if license already exists for this location and application
	existing, _ := s.licenseRepo.FindByLocationAndApplication(ctx, locationID, appID)
	if existing != nil && existing.Status == models.LicenseStatusActive {
		return nil, fmt.Errorf("active license already exists for this application at this location")
	}

	// Validate license type is supported
	supported := false
	for _, lt := range app.SupportedLicenseTypes {
		if lt == req.LicenseType {
			supported = true
			break
		}
	}
	if !supported {
		return nil, fmt.Errorf("license type '%s' is not supported for this application", req.LicenseType)
	}

	// Create license
	license := &models.LocationLicense{
		OrganizationID:   orgID,
		LocationID:       locationID,
		ApplicationID:    appID,
		LicenseType:      req.LicenseType,
		Status:           models.LicenseStatusActive,
		LicenseKey:       uuid.New().String(),
		ActivationKey:    uuid.New().String(),
		MaxUsers:         req.MaxUsers,
		MaxDevices:       req.MaxDevices,
		MaxTransactions:  req.MaxTransactions,
		MaxConcurrent:    req.MaxConcurrent,
		MaxStorageGB:     req.MaxStorageGB,
		MaxAPICallsDaily: req.MaxAPICallsDaily,
		StartDate:        time.Now(),
		IsTrial:          req.IsTrial,
		IsPerpetual:      req.IsPerpetual,
		BillingCycle:     req.BillingCycle,
		PricePerCycle:    req.PricePerCycle,
		Currency:         "USD",
		AutoRenew:        true,
		EnabledFeatures:  req.EnabledFeatures,
		DisabledFeatures: req.DisabledFeatures,
	}

	// Handle trial period
	if req.IsTrial && req.TrialDays > 0 {
		trialEnd := time.Now().AddDate(0, 0, req.TrialDays)
		license.TrialEndDate = &trialEnd
		license.PricePerCycle = 0 // No charge during trial
	}

	// Handle expiry
	if !req.IsPerpetual {
		if req.DurationMonths > 0 {
			endDate := time.Now().AddDate(0, req.DurationMonths, 0)
			license.EndDate = &endDate
			expiresAt := endDate
			license.ExpiresAt = &expiresAt
		}
	}

	license.BaseModel.CreatedBy = createdBy

	if err := s.licenseRepo.Create(ctx, license); err != nil {
		return nil, fmt.Errorf("failed to create license: %w", err)
	}

	// Update application analytics
	totalLicenses, _ := s.licenseRepo.CountByApplication(ctx, appID)
	s.appRepo.UpdateAnalytics(ctx, appID, int(totalLicenses), app.TotalLocations+1)

	return license, nil
}

// GetLicense retrieves a license by ID
func (s *LicenseService) GetLicense(ctx context.Context, id primitive.ObjectID) (*models.LocationLicense, error) {
	license, err := s.licenseRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("license not found: %w", err)
	}
	return license, nil
}

// ListLicensesByLocation returns all licenses for a location
func (s *LicenseService) ListLicensesByLocation(ctx context.Context, locationID primitive.ObjectID, activeOnly bool) ([]*models.LocationLicense, error) {
	licenses, err := s.licenseRepo.FindByLocation(ctx, locationID, activeOnly)
	if err != nil {
		return nil, fmt.Errorf("failed to list licenses: %w", err)
	}
	return licenses, nil
}

// ListLicensesByOrganization returns all licenses for an organization
func (s *LicenseService) ListLicensesByOrganization(ctx context.Context, orgID primitive.ObjectID, page, limit int) ([]*models.LocationLicense, error) {
	licenses, err := s.licenseRepo.FindByOrganization(ctx, orgID, page, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to list licenses: %w", err)
	}
	return licenses, nil
}

// UpdateUsageRequest represents a usage update request
type UpdateUsageRequest struct {
	Users        int     `json:"users"`
	Devices      int     `json:"devices"`
	Transactions int64   `json:"transactions"`
	StorageGB    float64 `json:"storage_gb"`
}

// UpdateLicenseUsage updates license usage counters
func (s *LicenseService) UpdateLicenseUsage(ctx context.Context, licenseID primitive.ObjectID, req UpdateUsageRequest) error {
	// Get license
	license, err := s.licenseRepo.FindByID(ctx, licenseID)
	if err != nil {
		return fmt.Errorf("license not found: %w", err)
	}

	// Check limits and update status if exceeded
	exceeded := false
	if license.MaxUsers > 0 && req.Users > license.MaxUsers {
		exceeded = true
	}
	if license.MaxDevices > 0 && req.Devices > license.MaxDevices {
		exceeded = true
	}
	if license.MaxTransactions > 0 && req.Transactions > license.MaxTransactions {
		exceeded = true
	}
	if license.MaxStorageGB > 0 && req.StorageGB > license.MaxStorageGB {
		exceeded = true
	}

	if exceeded && license.Status == models.LicenseStatusActive {
		s.licenseRepo.UpdateStatus(ctx, licenseID, models.LicenseStatusExceeded, "Usage limits exceeded")
	}

	// Update usage
	if err := s.licenseRepo.UpdateUsage(ctx, licenseID, req.Users, req.Devices, req.Transactions, req.StorageGB); err != nil {
		return fmt.Errorf("failed to update usage: %w", err)
	}

	return nil
}

// SuspendLicense suspends a license
func (s *LicenseService) SuspendLicense(ctx context.Context, licenseID primitive.ObjectID, reason string) error {
	if err := s.licenseRepo.UpdateStatus(ctx, licenseID, models.LicenseStatusSuspended, reason); err != nil {
		return fmt.Errorf("failed to suspend license: %w", err)
	}
	return nil
}

// ActivateLicense activates a suspended license
func (s *LicenseService) ActivateLicense(ctx context.Context, licenseID primitive.ObjectID) error {
	if err := s.licenseRepo.UpdateStatus(ctx, licenseID, models.LicenseStatusActive, ""); err != nil {
		return fmt.Errorf("failed to activate license: %w", err)
	}
	return nil
}

// RevokeLicense revokes a license
func (s *LicenseService) RevokeLicense(ctx context.Context, licenseID, revokedBy primitive.ObjectID, reason string) error {
	if err := s.licenseRepo.SoftDelete(ctx, licenseID, revokedBy); err != nil {
		return fmt.Errorf("failed to revoke license: %w", err)
	}
	return nil
}

// AssignUserRequest represents a request to assign a user to a license
type AssignUserRequest struct {
	UserID      string   `json:"user_id" binding:"required"`
	Role        string   `json:"role"`
	Permissions []string `json:"permissions"`
}

// AssignUserToLicense assigns a user to a license
func (s *LicenseService) AssignUserToLicense(ctx context.Context, licenseID primitive.ObjectID, req AssignUserRequest, assignedBy primitive.ObjectID) error {
	// Get license
	license, err := s.licenseRepo.FindByID(ctx, licenseID)
	if err != nil {
		return fmt.Errorf("license not found: %w", err)
	}

	if license.Status != models.LicenseStatusActive {
		return fmt.Errorf("license is not active")
	}

	// Parse user ID
	userID, err := primitive.ObjectIDFromHex(req.UserID)
	if err != nil {
		return fmt.Errorf("invalid user ID: %w", err)
	}

	// Check if user limit is reached (for user-based licenses)
	if license.LicenseType == models.LicenseTypeUserBased && license.MaxUsers > 0 {
		currentUsers, _ := s.userAssignRepo.CountActiveByLicense(ctx, licenseID)
		if int(currentUsers) >= license.MaxUsers {
			return fmt.Errorf("user limit reached for this license")
		}
	}

	// Create assignment
	assignment := &models.UserLicenseAssignment{
		LicenseID:      licenseID,
		UserID:         userID,
		LocationID:     license.LocationID,
		ApplicationID:  license.ApplicationID,
		OrganizationID: license.OrganizationID,
		Role:           req.Role,
		Permissions:    req.Permissions,
	}

	assignment.BaseModel.CreatedBy = assignedBy

	if err := s.userAssignRepo.AssignUser(ctx, assignment); err != nil {
		return fmt.Errorf("failed to assign user: %w", err)
	}

	// Update license usage
	activeUsers, _ := s.userAssignRepo.CountActiveByLicense(ctx, licenseID)
	s.licenseRepo.UpdateUsage(ctx, licenseID, int(activeUsers), license.CurrentDevices, license.CurrentTransactions, license.CurrentStorageGB)

	return nil
}

// RevokeUserFromLicense revokes a user's access to a license
func (s *LicenseService) RevokeUserFromLicense(ctx context.Context, assignmentID, revokedBy primitive.ObjectID, reason string) error {
	if err := s.userAssignRepo.RevokeUser(ctx, assignmentID, revokedBy, reason); err != nil {
		return fmt.Errorf("failed to revoke user: %w", err)
	}
	return nil
}

// ActivateDeviceRequest represents a request to activate a device
type ActivateDeviceRequest struct {
	DeviceID    string `json:"device_id" binding:"required"`
	DeviceName  string `json:"device_name"`
	DeviceType  string `json:"device_type"`
	DeviceModel string `json:"device_model"`
	OS          string `json:"os"`
	OSVersion   string `json:"os_version"`
	MACAddress  string `json:"mac_address"`
	IPAddress   string `json:"ip_address"`
}

// ActivateDevice activates a device for a license
func (s *LicenseService) ActivateDevice(ctx context.Context, licenseID primitive.ObjectID, req ActivateDeviceRequest, activatedBy primitive.ObjectID) error {
	// Get license
	license, err := s.licenseRepo.FindByID(ctx, licenseID)
	if err != nil {
		return fmt.Errorf("license not found: %w", err)
	}

	if license.Status != models.LicenseStatusActive {
		return fmt.Errorf("license is not active")
	}

	// Check if device is already activated
	exists, _ := s.deviceAssignRepo.DeviceExists(ctx, req.DeviceID, licenseID)
	if exists {
		return fmt.Errorf("device already activated for this license")
	}

	// Check if device limit is reached (for device-based licenses)
	if license.LicenseType == models.LicenseTypeDeviceBased && license.MaxDevices > 0 {
		currentDevices, _ := s.deviceAssignRepo.CountActiveByLicense(ctx, licenseID)
		if int(currentDevices) >= license.MaxDevices {
			return fmt.Errorf("device limit reached for this license")
		}
	}

	// Create device assignment
	assignment := &models.DeviceLicenseAssignment{
		LicenseID:      licenseID,
		LocationID:     license.LocationID,
		ApplicationID:  license.ApplicationID,
		OrganizationID: license.OrganizationID,
		DeviceID:       req.DeviceID,
		DeviceName:     req.DeviceName,
		DeviceType:     req.DeviceType,
		DeviceModel:    req.DeviceModel,
		OS:             req.OS,
		OSVersion:      req.OSVersion,
		MACAddress:     req.MACAddress,
		IPAddress:      req.IPAddress,
		ActivationCode: uuid.New().String(),
	}

	assignment.BaseModel.CreatedBy = activatedBy

	if err := s.deviceAssignRepo.ActivateDevice(ctx, assignment); err != nil {
		return fmt.Errorf("failed to activate device: %w", err)
	}

	// Update license usage
	activeDevices, _ := s.deviceAssignRepo.CountActiveByLicense(ctx, licenseID)
	s.licenseRepo.UpdateUsage(ctx, licenseID, license.CurrentUsers, int(activeDevices), license.CurrentTransactions, license.CurrentStorageGB)

	return nil
}

// DeactivateDevice deactivates a device
func (s *LicenseService) DeactivateDevice(ctx context.Context, assignmentID, deactivatedBy primitive.ObjectID, reason string) error {
	if err := s.deviceAssignRepo.DeactivateDevice(ctx, assignmentID, deactivatedBy, reason); err != nil {
		return fmt.Errorf("failed to deactivate device: %w", err)
	}
	return nil
}
