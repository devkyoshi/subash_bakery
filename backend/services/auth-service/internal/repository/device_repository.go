package repository

import (
	"context"
	"fmt"

	"github.com/yourusername/erp-system/shared/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// DeviceRepository is a read-only repository for looking up registered devices
// during user registration. The devices collection is managed by the org-service.
type DeviceRepository struct {
	collection *mongo.Collection
}

func NewDeviceRepository(db *mongo.Database) *DeviceRepository {
	return &DeviceRepository{
		collection: db.Collection("devices"),
	}
}

// FindByMACAddress finds an active device by its MAC address
func (r *DeviceRepository) FindByMACAddress(ctx context.Context, macAddress string) (*models.Device, error) {
	var device models.Device
	filter := bson.M{
		"mac_address": macAddress,
		"is_active":   true,
		"deleted_at":  nil,
	}

	err := r.collection.FindOne(ctx, filter).Decode(&device)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find device by MAC address: %w", err)
	}
	return &device, nil
}
