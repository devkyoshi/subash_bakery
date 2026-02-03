package repository

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/yourusername/erp-system/shared/models"
)

type UserAssignmentRepository struct {
	collection *mongo.Collection
}

func NewUserAssignmentRepository(db *mongo.Database) *UserAssignmentRepository {
	return &UserAssignmentRepository{
		collection: db.Collection("user_license_assignments"),
	}
}

// AssignUser assigns a user to a license
func (r *UserAssignmentRepository) AssignUser(ctx context.Context, assignment *models.UserLicenseAssignment) error {
	assignment.ID = primitive.NewObjectID()
	assignment.CreatedAt = time.Now()
	assignment.UpdatedAt = time.Now()
	assignment.AssignedAt = time.Now()
	assignment.IsActive = true

	_, err := r.collection.InsertOne(ctx, assignment)
	return err
}

// RevokeUser revokes a user assignment
func (r *UserAssignmentRepository) RevokeUser(ctx context.Context, assignmentID, revokedBy primitive.ObjectID, reason string) error {
	now := time.Now()
	filter := bson.M{"_id": assignmentID}
	update := bson.M{
		"$set": bson.M{
			"is_active":     false,
			"revoked_at":    now,
			"revoke_by":     revokedBy,
			"revoke_reason": reason,
			"updated_at":    now,
		},
	}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

// FindByLicense finds all user assignments for a license
func (r *UserAssignmentRepository) FindByLicense(ctx context.Context, licenseID primitive.ObjectID, activeOnly bool) ([]*models.UserLicenseAssignment, error) {
	filter := bson.M{"license_id": licenseID}
	if activeOnly {
		filter["is_active"] = true
	}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var assignments []*models.UserLicenseAssignment
	if err = cursor.All(ctx, &assignments); err != nil {
		return nil, err
	}

	return assignments, nil
}

// FindByUser finds all license assignments for a user
func (r *UserAssignmentRepository) FindByUser(ctx context.Context, userID primitive.ObjectID, activeOnly bool) ([]*models.UserLicenseAssignment, error) {
	filter := bson.M{"user_id": userID}
	if activeOnly {
		filter["is_active"] = true
	}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var assignments []*models.UserLicenseAssignment
	if err = cursor.All(ctx, &assignments); err != nil {
		return nil, err
	}

	return assignments, nil
}

// CountActiveByLicense counts active user assignments for a license
func (r *UserAssignmentRepository) CountActiveByLicense(ctx context.Context, licenseID primitive.ObjectID) (int64, error) {
	filter := bson.M{
		"license_id": licenseID,
		"is_active":  true,
	}
	return r.collection.CountDocuments(ctx, filter)
}

// UserHasLicense checks if a user is assigned to a specific application
func (r *UserAssignmentRepository) UserHasLicense(ctx context.Context, userID, applicationID primitive.ObjectID) (bool, error) {
	filter := bson.M{
		"user_id":        userID,
		"application_id": applicationID,
		"is_active":      true,
	}
	count, err := r.collection.CountDocuments(ctx, filter)
	return count > 0, err
}

type DeviceAssignmentRepository struct {
	collection *mongo.Collection
}

func NewDeviceAssignmentRepository(db *mongo.Database) *DeviceAssignmentRepository {
	return &DeviceAssignmentRepository{
		collection: db.Collection("device_license_assignments"),
	}
}

// ActivateDevice activates a device for a license
func (r *DeviceAssignmentRepository) ActivateDevice(ctx context.Context, assignment *models.DeviceLicenseAssignment) error {
	assignment.ID = primitive.NewObjectID()
	assignment.CreatedAt = time.Now()
	assignment.UpdatedAt = time.Now()
	assignment.ActivatedAt = time.Now()
	assignment.IsActive = true
	assignment.IsOnline = true

	_, err := r.collection.InsertOne(ctx, assignment)
	return err
}

// DeactivateDevice deactivates a device
func (r *DeviceAssignmentRepository) DeactivateDevice(ctx context.Context, assignmentID, deactivatedBy primitive.ObjectID, reason string) error {
	now := time.Now()
	filter := bson.M{"_id": assignmentID}
	update := bson.M{
		"$set": bson.M{
			"is_active":          false,
			"is_online":          false,
			"deactivated_at":     now,
			"deactivate_by":      deactivatedBy,
			"deactivate_reason":  reason,
			"updated_at":         now,
		},
	}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

// UpdateDeviceStatus updates device online status
func (r *DeviceAssignmentRepository) UpdateDeviceStatus(ctx context.Context, deviceID primitive.ObjectID, isOnline bool) error {
	now := time.Now()
	filter := bson.M{"_id": deviceID}
	update := bson.M{
		"$set": bson.M{
			"is_online":   isOnline,
			"last_seen_at": now,
			"updated_at":   now,
		},
	}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

// FindByLicense finds all device assignments for a license
func (r *DeviceAssignmentRepository) FindByLicense(ctx context.Context, licenseID primitive.ObjectID, activeOnly bool) ([]*models.DeviceLicenseAssignment, error) {
	filter := bson.M{"license_id": licenseID}
	if activeOnly {
		filter["is_active"] = true
	}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var assignments []*models.DeviceLicenseAssignment
	if err = cursor.All(ctx, &assignments); err != nil {
		return nil, err
	}

	return assignments, nil
}

// FindByDeviceID finds a device assignment by device ID
func (r *DeviceAssignmentRepository) FindByDeviceID(ctx context.Context, deviceID string, licenseID primitive.ObjectID) (*models.DeviceLicenseAssignment, error) {
	var assignment models.DeviceLicenseAssignment
	filter := bson.M{
		"device_id":  deviceID,
		"license_id": licenseID,
		"is_active":  true,
	}
	err := r.collection.FindOne(ctx, filter).Decode(&assignment)
	if err != nil {
		return nil, err
	}
	return &assignment, nil
}

// CountActiveByLicense counts active device assignments for a license
func (r *DeviceAssignmentRepository) CountActiveByLicense(ctx context.Context, licenseID primitive.ObjectID) (int64, error) {
	filter := bson.M{
		"license_id": licenseID,
		"is_active":  true,
	}
	return r.collection.CountDocuments(ctx, filter)
}

// DeviceExists checks if a device is already activated
func (r *DeviceAssignmentRepository) DeviceExists(ctx context.Context, deviceID string, licenseID primitive.ObjectID) (bool, error) {
	filter := bson.M{
		"device_id":  deviceID,
		"license_id": licenseID,
		"is_active":  true,
	}
	count, err := r.collection.CountDocuments(ctx, filter)
	return count > 0, err
}
