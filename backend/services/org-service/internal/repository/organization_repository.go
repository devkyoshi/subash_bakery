package repository

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"github.com/yourusername/erp-system/shared/models"
)

type OrganizationRepository struct {
	collection *mongo.Collection
}

func NewOrganizationRepository(db *mongo.Database) *OrganizationRepository {
	return &OrganizationRepository{
		collection: db.Collection("organizations"),
	}
}

// Create creates a new organization
func (r *OrganizationRepository) Create(ctx context.Context, org *models.Organization) error {
	org.ID = primitive.NewObjectID()
	org.CreatedAt = time.Now()
	org.UpdatedAt = time.Now()

	_, err := r.collection.InsertOne(ctx, org)
	if err != nil {
		return fmt.Errorf("failed to create organization: %w", err)
	}
	return nil
}

// FindByID finds an organization by ID
func (r *OrganizationRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*models.Organization, error) {
	var org models.Organization
	filter := bson.M{
		"_id":        id,
		"deleted_at": nil, // Only return non-deleted
	}

	err := r.collection.FindOne(ctx, filter).Decode(&org)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find organization: %w", err)
	}
	return &org, nil
}

// FindByDomain finds an organization by domain
func (r *OrganizationRepository) FindByDomain(ctx context.Context, domain string) (*models.Organization, error) {
	var org models.Organization
	filter := bson.M{
		"$or": []bson.M{
			{"domain": domain},
			{"alternate_domains": domain},
		},
		"deleted_at": nil,
	}

	err := r.collection.FindOne(ctx, filter).Decode(&org)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find organization by domain: %w", err)
	}
	return &org, nil
}

// List returns paginated organizations
func (r *OrganizationRepository) List(ctx context.Context, page, limit int, status string) ([]*models.Organization, int64, error) {
	filter := bson.M{"deleted_at": nil}

	if status != "" {
		filter["status"] = status
	}

	// Get total count
	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count organizations: %w", err)
	}

	// Get paginated results
	skip := (page - 1) * limit
	opts := options.Find().
		SetSkip(int64(skip)).
		SetLimit(int64(limit)).
		SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to find organizations: %w", err)
	}
	defer cursor.Close(ctx)

	var orgs []*models.Organization
	if err := cursor.All(ctx, &orgs); err != nil {
		return nil, 0, fmt.Errorf("failed to decode organizations: %w", err)
	}

	return orgs, total, nil
}

// Update updates an organization
func (r *OrganizationRepository) Update(ctx context.Context, org *models.Organization) error {
	org.UpdatedAt = time.Now()
	org.Version++ // Increment version for optimistic locking

	filter := bson.M{
		"_id":     org.ID,
		"version": org.Version - 1, // Check version hasn't changed
	}

	update := bson.M{"$set": org}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update organization: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("organization not found or version conflict")
	}

	return nil
}

// SoftDelete soft deletes an organization
func (r *OrganizationRepository) SoftDelete(ctx context.Context, id, deletedBy primitive.ObjectID) error {
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
		return fmt.Errorf("failed to soft delete organization: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("organization not found")
	}

	return nil
}

// UpdateUsage updates organization usage counters
func (r *OrganizationRepository) UpdateUsage(ctx context.Context, id primitive.ObjectID, users, companies, locations int, storageGB float64) error {
	filter := bson.M{"_id": id}
	update := bson.M{
		"$set": bson.M{
			"current_users":     users,
			"current_companies": companies,
			"current_locations": locations,
			"storage_used_gb":   storageGB,
			"updated_at":        time.Now(),
		},
	}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update organization usage: %w", err)
	}
	return nil
}

// DomainExists checks if a domain already exists
func (r *OrganizationRepository) DomainExists(ctx context.Context, domain string, excludeID *primitive.ObjectID) (bool, error) {
	filter := bson.M{
		"$or": []bson.M{
			{"domain": domain},
			{"alternate_domains": domain},
		},
		"deleted_at": nil,
	}

	if excludeID != nil {
		filter["_id"] = bson.M{"$ne": *excludeID}
	}

	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return false, fmt.Errorf("failed to check domain existence: %w", err)
	}

	return count > 0, nil
}
