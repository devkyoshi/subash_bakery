package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/yourusername/erp-system/shared/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DeviceRepository struct {
	collection *mongo.Collection
}

func NewDeviceRepository(db *mongo.Database) *DeviceRepository {
	return &DeviceRepository{
		collection: db.Collection("devices"),
	}
}

// Create inserts a new device
func (r *DeviceRepository) Create(ctx context.Context, device *models.Device) error {
	device.ID = primitive.NewObjectID()
	device.CreatedAt = time.Now()
	device.UpdatedAt = time.Now()

	_, err := r.collection.InsertOne(ctx, device)
	if err != nil {
		return fmt.Errorf("failed to create device: %w", err)
	}
	return nil
}

// FindByID finds a device by its ID (excluding soft-deleted)
func (r *DeviceRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*models.Device, error) {
	var device models.Device
	filter := bson.M{
		"_id":        id,
		"deleted_at": nil,
	}

	err := r.collection.FindOne(ctx, filter).Decode(&device)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find device: %w", err)
	}
	return &device, nil
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

// MACAddressExists checks if a MAC address is already registered
func (r *DeviceRepository) MACAddressExists(ctx context.Context, macAddress string, excludeID *primitive.ObjectID) (bool, error) {
	filter := bson.M{
		"mac_address": macAddress,
		"deleted_at":  nil,
	}

	if excludeID != nil {
		filter["_id"] = bson.M{"$ne": *excludeID}
	}

	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return false, fmt.Errorf("failed to check MAC address existence: %w", err)
	}
	return count > 0, nil
}

// List returns paginated devices for an organization
func (r *DeviceRepository) List(ctx context.Context, orgID primitive.ObjectID, filters map[string]interface{}, page, limit int) ([]*models.Device, int64, error) {
	filter := bson.M{
		"organization_id": orgID,
		"deleted_at":      nil,
	}

	// Apply search filter
	if search, ok := filters["search"].(string); ok && search != "" {
		filter["$or"] = []bson.M{
			{"name": bson.M{"$regex": search, "$options": "i"}},
			{"mac_address": bson.M{"$regex": search, "$options": "i"}},
			{"description": bson.M{"$regex": search, "$options": "i"}},
		}
	}

	// Apply status filter
	if isActive, ok := filters["is_active"].(bool); ok {
		filter["is_active"] = isActive
	}

	// Apply device type filter
	if deviceType, ok := filters["device_type"].(string); ok && deviceType != "" {
		filter["device_type"] = deviceType
	}

	// Get total count
	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count devices: %w", err)
	}

	// Get paginated results
	skip := (page - 1) * limit
	opts := options.Find().
		SetSkip(int64(skip)).
		SetLimit(int64(limit)).
		SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to find devices: %w", err)
	}
	defer cursor.Close(ctx)

	var devices []*models.Device
	if err := cursor.All(ctx, &devices); err != nil {
		return nil, 0, fmt.Errorf("failed to decode devices: %w", err)
	}

	return devices, total, nil
}

// Update updates a device
func (r *DeviceRepository) Update(ctx context.Context, device *models.Device) error {
	device.UpdatedAt = time.Now()
	device.Version++

	filter := bson.M{
		"_id":     device.ID,
		"version": device.Version - 1,
	}

	update := bson.M{"$set": device}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update device: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("device not found or version conflict")
	}

	return nil
}

// SoftDelete soft deletes a device
func (r *DeviceRepository) SoftDelete(ctx context.Context, id, deletedBy primitive.ObjectID) error {
	now := time.Now()
	filter := bson.M{"_id": id, "deleted_at": nil}
	update := bson.M{
		"$set": bson.M{
			"deleted_at": now,
			"deleted_by": deletedBy,
			"updated_at": now,
			"is_active":  false,
		},
	}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to soft delete device: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("device not found")
	}

	return nil
}

// CountByOrganization counts devices for an organization
func (r *DeviceRepository) CountByOrganization(ctx context.Context, orgID primitive.ObjectID) (int, error) {
	filter := bson.M{
		"organization_id": orgID,
		"deleted_at":      nil,
	}

	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, fmt.Errorf("failed to count devices: %w", err)
	}
	return int(count), nil
}
