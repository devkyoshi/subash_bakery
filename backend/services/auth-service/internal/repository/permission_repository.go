package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/yourusername/erp-system/shared/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type PermissionRepository struct {
	collection *mongo.Collection
}

func NewPermissionRepository(db *mongo.Database) *PermissionRepository {
	return &PermissionRepository{
		collection: db.Collection("permissions"),
	}
}

// Create creates a new permission
func (r *PermissionRepository) Create(ctx context.Context, permission *models.Permission) error {
	permission.ID = primitive.NewObjectID()
	permission.CreatedAt = time.Now()
	permission.UpdatedAt = time.Now()

	_, err := r.collection.InsertOne(ctx, permission)
	if err != nil {
		return fmt.Errorf("failed to create permission: %w", err)
	}
	return nil
}

// FindAll finds all permissions
func (r *PermissionRepository) FindAll(ctx context.Context) ([]*models.Permission, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("failed to find permissions: %w", err)
	}
	defer cursor.Close(ctx)

	var permissions []*models.Permission
	if err := cursor.All(ctx, &permissions); err != nil {
		return nil, fmt.Errorf("failed to decode permissions: %w", err)
	}
	return permissions, nil
}

// FindByIDs finds permissions by a list of IDs
func (r *PermissionRepository) FindByIDs(ctx context.Context, ids []primitive.ObjectID) ([]*models.Permission, error) {
	if len(ids) == 0 {
		return []*models.Permission{}, nil
	}

	filter := bson.M{"_id": bson.M{"$in": ids}}
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to find permissions: %w", err)
	}
	defer cursor.Close(ctx)

	var permissions []*models.Permission
	if err := cursor.All(ctx, &permissions); err != nil {
		return nil, fmt.Errorf("failed to decode permissions: %w", err)
	}
	return permissions, nil
}

// FindByName finds a permission by name
func (r *PermissionRepository) FindByName(ctx context.Context, name string) (*models.Permission, error) {
	var permission models.Permission
	err := r.collection.FindOne(ctx, bson.M{"name": name}).Decode(&permission)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find permission by name: %w", err)
	}
	return &permission, nil
}

// BulkCreate creates multiple permissions at once
func (r *PermissionRepository) BulkCreate(ctx context.Context, permissions []*models.Permission) error {
	if len(permissions) == 0 {
		return nil
	}

	now := time.Now()
	docs := make([]interface{}, len(permissions))
	for i, permission := range permissions {
		permission.ID = primitive.NewObjectID()
		permission.CreatedAt = now
		permission.UpdatedAt = now
		docs[i] = permission
	}

	_, err := r.collection.InsertMany(ctx, docs)
	if err != nil {
		return fmt.Errorf("failed to bulk create permissions: %w", err)
	}
	return nil
}
