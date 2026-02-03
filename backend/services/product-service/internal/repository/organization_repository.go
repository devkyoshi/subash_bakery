package repository

import (
	"context"
	"fmt"

	"github.com/yourusername/erp-system/shared/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type OrganizationRepository struct {
	collection *mongo.Collection
}

func NewOrganizationRepository(db *mongo.Database) *OrganizationRepository {
	return &OrganizationRepository{
		collection: db.Collection("organizations"),
	}
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

// Exists checks if an organization exists
func (r *OrganizationRepository) Exists(ctx context.Context, id primitive.ObjectID) (bool, error) {
	filter := bson.M{
		"_id":        id,
		"deleted_at": nil,
	}

	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return false, fmt.Errorf("failed to check organization existence: %w", err)
	}

	return count > 0, nil
}
