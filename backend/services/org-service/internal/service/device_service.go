package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/yourusername/erp-system/services/org-service/internal/repository"
	"github.com/yourusername/erp-system/shared/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DeviceService struct {
	deviceRepo *repository.DeviceRepository
	orgRepo    *repository.OrganizationRepository
}

func NewDeviceService(
	deviceRepo *repository.DeviceRepository,
	orgRepo *repository.OrganizationRepository,
) *DeviceService {
	return &DeviceService{
		deviceRepo: deviceRepo,
		orgRepo:    orgRepo,
	}
}

// --- DTOs ---

type CreateDeviceRequest struct {
	OrganizationID string `json:"organization_id" binding:"required"`
	Name           string `json:"name" binding:"required"`
	MACAddress     string `json:"mac_address" binding:"required"`
	DeviceType     string `json:"device_type" binding:"required"`
	Description    string `json:"description"`
	Location       string `json:"location"`
}

type UpdateDeviceRequest struct {
	Name        *string `json:"name"`
	MACAddress  *string `json:"mac_address"`
	DeviceType  *string `json:"device_type"`
	Description *string `json:"description"`
	Location    *string `json:"location"`
	IsActive    *bool   `json:"is_active"`
}

type DeviceLookupResponse struct {
	OrganizationID string `json:"organization_id"`
	DeviceID       string `json:"device_id"`
	DeviceName     string `json:"device_name"`
}

// --- Service Methods ---

// CreateDevice registers a new device for an organization
func (s *DeviceService) CreateDevice(ctx context.Context, req CreateDeviceRequest, createdBy primitive.ObjectID) (*models.Device, error) {
	// Normalize MAC address (uppercase, colon-separated)
	macAddress := normalizeMACAddress(req.MACAddress)
	if macAddress == "" {
		return nil, fmt.Errorf("invalid MAC address format")
	}

	// Verify organization exists
	orgID, err := primitive.ObjectIDFromHex(req.OrganizationID)
	if err != nil {
		return nil, fmt.Errorf("invalid organization ID: %w", err)
	}

	org, err := s.orgRepo.FindByID(ctx, orgID)
	if err != nil {
		return nil, err
	}
	if org == nil {
		return nil, fmt.Errorf("organization not found")
	}

	// Check MAC address uniqueness
	exists, err := s.deviceRepo.MACAddressExists(ctx, macAddress, nil)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, fmt.Errorf("a device with this MAC address is already registered")
	}

	// Validate device type
	deviceType := models.DeviceType(req.DeviceType)
	if !isValidDeviceType(deviceType) {
		return nil, fmt.Errorf("invalid device type: %s", req.DeviceType)
	}

	device := &models.Device{
		OrganizationID: orgID,
		Name:           req.Name,
		MACAddress:     macAddress,
		DeviceType:     deviceType,
		Description:    req.Description,
		Location:       req.Location,
		IsActive:       true,
	}

	device.CreatedBy = createdBy

	if err := s.deviceRepo.Create(ctx, device); err != nil {
		return nil, err
	}

	return device, nil
}

// GetDevice retrieves a device by ID
func (s *DeviceService) GetDevice(ctx context.Context, id primitive.ObjectID) (*models.Device, error) {
	device, err := s.deviceRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if device == nil {
		return nil, fmt.Errorf("device not found")
	}
	return device, nil
}

// ListDevices returns paginated devices for an organization
func (s *DeviceService) ListDevices(ctx context.Context, orgID primitive.ObjectID, filters map[string]interface{}, page, limit int) ([]*models.Device, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	return s.deviceRepo.List(ctx, orgID, filters, page, limit)
}

// UpdateDevice updates an existing device
func (s *DeviceService) UpdateDevice(ctx context.Context, id primitive.ObjectID, req UpdateDeviceRequest, updatedBy primitive.ObjectID) (*models.Device, error) {
	device, err := s.deviceRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if device == nil {
		return nil, fmt.Errorf("device not found")
	}

	// Update MAC address if changed
	if req.MACAddress != nil && *req.MACAddress != device.MACAddress {
		macAddress := normalizeMACAddress(*req.MACAddress)
		if macAddress == "" {
			return nil, fmt.Errorf("invalid MAC address format")
		}

		exists, err := s.deviceRepo.MACAddressExists(ctx, macAddress, &id)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, fmt.Errorf("a device with this MAC address is already registered")
		}

		device.MACAddress = macAddress
	}

	if req.Name != nil {
		device.Name = *req.Name
	}
	if req.DeviceType != nil {
		deviceType := models.DeviceType(*req.DeviceType)
		if !isValidDeviceType(deviceType) {
			return nil, fmt.Errorf("invalid device type: %s", *req.DeviceType)
		}
		device.DeviceType = deviceType
	}
	if req.Description != nil {
		device.Description = *req.Description
	}
	if req.Location != nil {
		device.Location = *req.Location
	}
	if req.IsActive != nil {
		device.IsActive = *req.IsActive
	}

	device.UpdatedBy = updatedBy

	if err := s.deviceRepo.Update(ctx, device); err != nil {
		return nil, err
	}

	return device, nil
}

// DeleteDevice soft deletes a device
func (s *DeviceService) DeleteDevice(ctx context.Context, id, deletedBy primitive.ObjectID) error {
	device, err := s.deviceRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if device == nil {
		return fmt.Errorf("device not found")
	}

	return s.deviceRepo.SoftDelete(ctx, id, deletedBy)
}

// LookupDeviceByMAC finds a device by MAC address and returns organization info
// This is used by the auth-service registration flow (public endpoint)
func (s *DeviceService) LookupDeviceByMAC(ctx context.Context, macAddress string) (*DeviceLookupResponse, error) {
	macAddress = normalizeMACAddress(macAddress)
	if macAddress == "" {
		return nil, fmt.Errorf("invalid MAC address format")
	}

	device, err := s.deviceRepo.FindByMACAddress(ctx, macAddress)
	if err != nil {
		return nil, err
	}
	if device == nil {
		return nil, fmt.Errorf("device not registered")
	}

	return &DeviceLookupResponse{
		OrganizationID: device.OrganizationID.Hex(),
		DeviceID:       device.ID.Hex(),
		DeviceName:     device.Name,
	}, nil
}

// --- Helpers ---

// normalizeMACAddress normalizes a MAC address to uppercase colon-separated format
func normalizeMACAddress(mac string) string {
	// Remove common separators and spaces
	mac = strings.TrimSpace(mac)
	mac = strings.ReplaceAll(mac, "-", "")
	mac = strings.ReplaceAll(mac, ":", "")
	mac = strings.ReplaceAll(mac, ".", "")
	mac = strings.ToUpper(mac)

	// Validate length (6 bytes = 12 hex chars)
	if len(mac) != 12 {
		return ""
	}

	// Validate hex characters
	for _, c := range mac {
		if !((c >= '0' && c <= '9') || (c >= 'A' && c <= 'F')) {
			return ""
		}
	}

	// Format as XX:XX:XX:XX:XX:XX
	parts := make([]string, 6)
	for i := 0; i < 6; i++ {
		parts[i] = mac[i*2 : i*2+2]
	}
	return strings.Join(parts, ":")
}

func isValidDeviceType(dt models.DeviceType) bool {
	switch dt {
	case models.DeviceTypePOS, models.DeviceTypeTablet, models.DeviceTypeMobile,
		models.DeviceTypeDesktop, models.DeviceTypeKiosk, models.DeviceTypeOther:
		return true
	}
	return false
}
