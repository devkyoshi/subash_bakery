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

type LocationRepository struct {
	collection *mongo.Collection
}

func NewLocationRepository(db *mongo.Database) *LocationRepository {
	return &LocationRepository{
		collection: db.Collection("locations"),
	}
}

// Create creates a new location
func (r *LocationRepository) Create(ctx context.Context, location *models.Location) error {
	location.ID = primitive.NewObjectID()
	location.CreatedAt = time.Now()
	location.UpdatedAt = time.Now()

	_, err := r.collection.InsertOne(ctx, location)
	if err != nil {
		return fmt.Errorf("failed to create location: %w", err)
	}
	return nil
}

// FindByID finds a location by ID
func (r *LocationRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*models.Location, error) {
	var location models.Location
	filter := bson.M{
		"_id":        id,
		"deleted_at": nil,
	}

	err := r.collection.FindOne(ctx, filter).Decode(&location)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find location: %w", err)
	}
	return &location, nil
}

// FindByCompany finds all locations for a company
func (r *LocationRepository) FindByCompany(ctx context.Context, companyID primitive.ObjectID, page, limit int) ([]*models.Location, int64, error) {
	filter := bson.M{
		"company_id": companyID,
		"deleted_at": nil,
	}

	// Get total count
	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count locations: %w", err)
	}

	// Get paginated results
	skip := (page - 1) * limit
	opts := options.Find().
		SetSkip(int64(skip)).
		SetLimit(int64(limit)).
		SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to find locations: %w", err)
	}
	defer cursor.Close(ctx)

	var locations []*models.Location
	if err := cursor.All(ctx, &locations); err != nil {
		return nil, 0, fmt.Errorf("failed to decode locations: %w", err)
	}

	return locations, total, nil
}

// FindByOrganization finds all locations for an organization
func (r *LocationRepository) FindByOrganization(ctx context.Context, orgID primitive.ObjectID, locationType string, page, limit int) ([]*models.Location, int64, error) {
	filter := bson.M{
		"organization_id": orgID,
		"deleted_at":      nil,
	}

	if locationType != "" {
		filter["type"] = locationType
	}

	// Get total count
	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count locations: %w", err)
	}

	// Get paginated results
	skip := (page - 1) * limit
	opts := options.Find().
		SetSkip(int64(skip)).
		SetLimit(int64(limit)).
		SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to find locations: %w", err)
	}
	defer cursor.Close(ctx)

	var locations []*models.Location
	if err := cursor.All(ctx, &locations); err != nil {
		return nil, 0, fmt.Errorf("failed to decode locations: %w", err)
	}

	return locations, total, nil
}

// Update updates a location
func (r *LocationRepository) Update(ctx context.Context, location *models.Location) error {
	location.UpdatedAt = time.Now()
	location.Version++

	filter := bson.M{
		"_id":     location.ID,
		"version": location.Version - 1,
	}

	update := bson.M{"$set": location}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update location: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("location not found or version conflict")
	}

	return nil
}

// SoftDelete soft deletes a location
func (r *LocationRepository) SoftDelete(ctx context.Context, id, deletedBy primitive.ObjectID) error {
	now := time.Now()
	filter := bson.M{"_id": id}
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
		return fmt.Errorf("failed to soft delete location: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("location not found")
	}

	return nil
}

// CodeExists checks if a location code exists within a company
func (r *LocationRepository) CodeExists(ctx context.Context, companyID primitive.ObjectID, code string, excludeID *primitive.ObjectID) (bool, error) {
	filter := bson.M{
		"company_id": companyID,
		"code":       code,
		"deleted_at": nil,
	}

	if excludeID != nil {
		filter["_id"] = bson.M{"$ne": *excludeID}
	}

	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return false, fmt.Errorf("failed to check code existence: %w", err)
	}

	return count > 0, nil
}

// CountByOrganization counts locations for an organization
func (r *LocationRepository) CountByOrganization(ctx context.Context, orgID primitive.ObjectID) (int, error) {
	filter := bson.M{
		"organization_id": orgID,
		"deleted_at":      nil,
	}

	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, fmt.Errorf("failed to count locations: %w", err)
	}

	return int(count), nil
}

// FindByIDs finds multiple locations by their IDs
func (r *LocationRepository) FindByIDs(ctx context.Context, locationIDs []primitive.ObjectID) ([]*models.Location, error) {
	if len(locationIDs) == 0 {
		return []*models.Location{}, nil
	}

	filter := bson.M{
		"_id":        bson.M{"$in": locationIDs},
		"deleted_at": nil,
	}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to find locations: %w", err)
	}
	defer cursor.Close(ctx)

	var locations []*models.Location
	if err := cursor.All(ctx, &locations); err != nil {
		return nil, fmt.Errorf("failed to decode locations: %w", err)
	}

	return locations, nil
}
