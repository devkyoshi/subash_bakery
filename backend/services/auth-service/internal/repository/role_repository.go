package repository

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"github.com/yourusername/erp-system/shared/models"
)

type RoleRepository struct {
	collection *mongo.Collection
}

func NewRoleRepository(db *mongo.Database) *RoleRepository {
	return &RoleRepository{
		collection: db.Collection("roles"),
	}
}

// Create creates a new role
func (r *RoleRepository) Create(ctx context.Context, role *models.Role) error {
	role.ID = primitive.NewObjectID()
	role.CreatedAt = time.Now()
	role.UpdatedAt = time.Now()

	_, err := r.collection.InsertOne(ctx, role)
	if err != nil {
		return fmt.Errorf("failed to create role: %w", err)
	}
	return nil
}

// FindByID finds a role by ID
func (r *RoleRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*models.Role, error) {
	var role models.Role
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&role)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find role by ID: %w", err)
	}
	return &role, nil
}

// FindAll finds all roles, optionally filtered by organization
func (r *RoleRepository) FindAll(ctx context.Context, orgID *primitive.ObjectID) ([]*models.Role, error) {
	filter := bson.M{}
	if orgID != nil {
		filter["organization_id"] = orgID
	}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to find roles: %w", err)
	}
	defer cursor.Close(ctx)

	var roles []*models.Role
	if err := cursor.All(ctx, &roles); err != nil {
		return nil, fmt.Errorf("failed to decode roles: %w", err)
	}
	return roles, nil
}

// Update updates a role
func (r *RoleRepository) Update(ctx context.Context, role *models.Role) error {
	role.UpdatedAt = time.Now()
	filter := bson.M{"_id": role.ID}
	update := bson.M{"$set": role}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update role: %w", err)
	}
	return nil
}

// Delete permanently deletes a role (or soft delete if preferred, but base model has DeletedAt)
func (r *RoleRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	// Soft delete implementation
	update := bson.M{
		"$set": bson.M{
			"deleted_at": time.Now(),
		},
	}
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	if err != nil {
		return fmt.Errorf("failed to delete role: %w", err)
	}
	return nil
}
