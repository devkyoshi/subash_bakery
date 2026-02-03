package repository

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/yourusername/erp-system/shared/models"
)

type LicenseRepository struct {
	collection *mongo.Collection
}

func NewLicenseRepository(db *mongo.Database) *LicenseRepository {
	return &LicenseRepository{
		collection: db.Collection("location_licenses"),
	}
}

// Create creates a new license
func (r *LicenseRepository) Create(ctx context.Context, license *models.LocationLicense) error {
	license.ID = primitive.NewObjectID()
	license.CreatedAt = time.Now()
	license.UpdatedAt = time.Now()
	license.Status = models.LicenseStatusActive

	_, err := r.collection.InsertOne(ctx, license)
	return err
}

// FindByID finds a license by ID
func (r *LicenseRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*models.LocationLicense, error) {
	var license models.LocationLicense
	filter := bson.M{
		"_id":        id,
		"deleted_at": nil,
	}
	err := r.collection.FindOne(ctx, filter).Decode(&license)
	if err != nil {
		return nil, err
	}
	return &license, nil
}

// FindByLocation returns all licenses for a location
func (r *LicenseRepository) FindByLocation(ctx context.Context, locationID primitive.ObjectID, activeOnly bool) ([]*models.LocationLicense, error) {
	filter := bson.M{
		"location_id": locationID,
		"deleted_at":  nil,
	}
	
	if activeOnly {
		filter["status"] = models.LicenseStatusActive
	}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var licenses []*models.LocationLicense
	if err = cursor.All(ctx, &licenses); err != nil {
		return nil, err
	}

	return licenses, nil
}

// FindByOrganization returns all licenses for an organization
func (r *LicenseRepository) FindByOrganization(ctx context.Context, orgID primitive.ObjectID, page, limit int) ([]*models.LocationLicense, error) {
	filter := bson.M{
		"organization_id": orgID,
		"deleted_at":      nil,
	}

	opts := options.Find()
	if limit > 0 {
		opts.SetLimit(int64(limit))
		opts.SetSkip(int64((page - 1) * limit))
	}
	opts.SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var licenses []*models.LocationLicense
	if err = cursor.All(ctx, &licenses); err != nil {
		return nil, err
	}

	return licenses, nil
}

// FindByLocationAndApplication finds a license by location and application
func (r *LicenseRepository) FindByLocationAndApplication(ctx context.Context, locationID, appID primitive.ObjectID) (*models.LocationLicense, error) {
	var license models.LocationLicense
	filter := bson.M{
		"location_id":    locationID,
		"application_id": appID,
		"deleted_at":     nil,
	}
	err := r.collection.FindOne(ctx, filter).Decode(&license)
	if err != nil {
		return nil, err
	}
	return &license, nil
}

// Update updates a license
func (r *LicenseRepository) Update(ctx context.Context, license *models.LocationLicense) error {
	license.UpdatedAt = time.Now()
	license.Version++

	filter := bson.M{
		"_id":     license.ID,
		"version": license.Version - 1,
	}

	update := bson.M{
		"$set": license,
	}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}

// UpdateUsage updates license usage counters
func (r *LicenseRepository) UpdateUsage(ctx context.Context, licenseID primitive.ObjectID, users, devices int, transactions int64, storageGB float64) error {
	filter := bson.M{"_id": licenseID}
	update := bson.M{
		"$set": bson.M{
			"current_users":        users,
			"current_devices":      devices,
			"current_transactions": transactions,
			"current_storage_gb":   storageGB,
			"updated_at":           time.Now(),
		},
	}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

// UpdateStatus updates license status
func (r *LicenseRepository) UpdateStatus(ctx context.Context, licenseID primitive.ObjectID, status models.LicenseStatus, reason string) error {
	now := time.Now()
	update := bson.M{
		"$set": bson.M{
			"status":     status,
			"updated_at": now,
		},
	}

	if status == models.LicenseStatusSuspended {
		update["$set"].(bson.M)["suspended_at"] = now
		update["$set"].(bson.M)["suspend_reason"] = reason
	} else if status == models.LicenseStatusRevoked {
		update["$set"].(bson.M)["revoked_at"] = now
		update["$set"].(bson.M)["revoke_reason"] = reason
	}

	filter := bson.M{"_id": licenseID}
	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

// SoftDelete soft deletes a license
func (r *LicenseRepository) SoftDelete(ctx context.Context, id, deletedBy primitive.ObjectID) error {
	now := time.Now()
	filter := bson.M{
		"_id":        id,
		"deleted_at": nil,
	}

	update := bson.M{
		"$set": bson.M{
			"deleted_at": now,
			"deleted_by": deletedBy,
			"updated_at": now,
			"status":     models.LicenseStatusRevoked,
		},
	}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}

// CountByApplication counts licenses for an application
func (r *LicenseRepository) CountByApplication(ctx context.Context, appID primitive.ObjectID) (int64, error) {
	filter := bson.M{
		"application_id": appID,
		"status":         models.LicenseStatusActive,
		"deleted_at":     nil,
	}
	return r.collection.CountDocuments(ctx, filter)
}

// CountActiveByLocation counts active licenses for a location
func (r *LicenseRepository) CountActiveByLocation(ctx context.Context, locationID primitive.ObjectID) (int64, error) {
	filter := bson.M{
		"location_id": locationID,
		"status":      models.LicenseStatusActive,
		"deleted_at":  nil,
	}
	return r.collection.CountDocuments(ctx, filter)
}

// FindExpiring finds licenses expiring within days
func (r *LicenseRepository) FindExpiring(ctx context.Context, withinDays int) ([]*models.LocationLicense, error) {
	now := time.Now()
	expiryDate := now.AddDate(0, 0, withinDays)

	filter := bson.M{
		"status":      models.LicenseStatusActive,
		"deleted_at":  nil,
		"is_perpetual": false,
		"expires_at": bson.M{
			"$gte": now,
			"$lte": expiryDate,
		},
	}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var licenses []*models.LocationLicense
	if err = cursor.All(ctx, &licenses); err != nil {
		return nil, err
	}

	return licenses, nil
}
