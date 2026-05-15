package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Device represents a registered organization device (e.g. a POS terminal)
// identified by its MAC address. Users registering from this device will be
// automatically associated with the device's organization.
type Device struct {
	BaseModel      `bson:",inline"`
	OrganizationID primitive.ObjectID     `json:"organization_id" bson:"organization_id"`
	Name           string                 `json:"name" bson:"name"`
	MACAddress     string                 `json:"mac_address" bson:"mac_address"`
	DeviceType     DeviceType             `json:"device_type" bson:"device_type"`
	Description    string                 `json:"description,omitempty" bson:"description,omitempty"`
	Location       string                 `json:"location,omitempty" bson:"location,omitempty"` // Physical location label
	IsActive       bool                   `json:"is_active" bson:"is_active"`
	LastSeenAt     *int64                 `json:"last_seen_at,omitempty" bson:"last_seen_at,omitempty"`
	Metadata       map[string]interface{} `json:"metadata,omitempty" bson:"metadata,omitempty"`
}

type DeviceType string

const (
	DeviceTypePOS     DeviceType = "pos"
	DeviceTypeTablet  DeviceType = "tablet"
	DeviceTypeMobile  DeviceType = "mobile"
	DeviceTypeDesktop DeviceType = "desktop"
	DeviceTypeKiosk   DeviceType = "kiosk"
	DeviceTypeOther   DeviceType = "other"
)

type DeviceStatus string

const (
	DeviceStatusActive   DeviceStatus = "active"
	DeviceStatusInactive DeviceStatus = "inactive"
)
