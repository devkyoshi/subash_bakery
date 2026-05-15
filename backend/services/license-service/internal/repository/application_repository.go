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

type ApplicationRepository struct {
	collection *mongo.Collection
}

func NewApplicationRepository(db *mongo.Database) *ApplicationRepository {
	return &ApplicationRepository{
		collection: db.Collection("applications"),
	}
}

// Create creates a new application
func (r *ApplicationRepository) Create(ctx context.Context, app *models.Application) error {
	app.ID = primitive.NewObjectID()
	app.CreatedAt = time.Now()
	app.UpdatedAt = time.Now()
	app.IsActive = true

	_, err := r.collection.InsertOne(ctx, app)
	return err
}

// FindByID finds an application by ID
func (r *ApplicationRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*models.Application, error) {
	var app models.Application
	filter := bson.M{
		"_id":        id,
		"deleted_at": nil,
	}
	err := r.collection.FindOne(ctx, filter).Decode(&app)
	if err != nil {
		return nil, err
	}
	return &app, nil
}

// FindByCode finds an application by code
func (r *ApplicationRepository) FindByCode(ctx context.Context, code string) (*models.Application, error) {
	var app models.Application
	filter := bson.M{
		"code":       code,
		"deleted_at": nil,
	}
	err := r.collection.FindOne(ctx, filter).Decode(&app)
	if err != nil {
		return nil, err
	}
	return &app, nil
}

// FindAll returns all active applications
func (r *ApplicationRepository) FindAll(ctx context.Context, category models.ApplicationCategory, publicOnly bool, page, limit int) ([]*models.Application, error) {
	filter := bson.M{"deleted_at": nil}
	
	if category != "" {
		filter["category"] = category
	}
	
	if publicOnly {
		filter["is_public"] = true
		filter["is_active"] = true
	}

	opts := options.Find()
	if limit > 0 {
		opts.SetLimit(int64(limit))
		opts.SetSkip(int64((page - 1) * limit))
	}
	opts.SetSort(bson.D{{Key: "display_name", Value: 1}})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var apps []*models.Application
	if err = cursor.All(ctx, &apps); err != nil {
		return nil, err
	}

	return apps, nil
}

// Update updates an application
func (r *ApplicationRepository) Update(ctx context.Context, app *models.Application) error {
	app.UpdatedAt = time.Now()
	// app.Version++

	// filter := bson.M{
	// 	"_id":     app.ID,
	// 	"version": app.Version - 1,
	// }

	// update := bson.M{
	// 	"$set": app,
	// }

	// // result, err := r.collection.UpdateOne(ctx, filter, update)
	// // if err != nil {
	// // 	return err
	// // }

	// if result.MatchedCount == 0 {
	// 	return mongo.ErrNoDocuments
	// }

	return nil
}

// SoftDelete soft deletes an application
func (r *ApplicationRepository) SoftDelete(ctx context.Context, id, deletedBy primitive.ObjectID) error {
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
			"is_active":  false,
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

// CodeExists checks if an application code already exists
func (r *ApplicationRepository) CodeExists(ctx context.Context, code string, excludeID *primitive.ObjectID) (bool, error) {
	filter := bson.M{
		"code":       code,
		"deleted_at": nil,
	}

	if excludeID != nil {
		filter["_id"] = bson.M{"$ne": *excludeID}
	}

	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// UpdateAnalytics updates application analytics
func (r *ApplicationRepository) UpdateAnalytics(ctx context.Context, appID primitive.ObjectID, totalLicenses, totalLocations int) error {
	filter := bson.M{"_id": appID}
	update := bson.M{
		"$set": bson.M{
			"total_licenses":  totalLicenses,
			"total_locations": totalLocations,
			"updated_at":      time.Now(),
		},
	}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}
