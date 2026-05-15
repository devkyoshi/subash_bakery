package repository

import (
	"context"
	"fmt"

	"github.com/yourusername/erp-system/services/dashboard-service/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ActivityRepository struct {
	collection *mongo.Collection
}

func NewActivityRepository(db *mongo.Database) *ActivityRepository {
	return &ActivityRepository{
		collection: db.Collection("activities"),
	}
}

func (r *ActivityRepository) Create(ctx context.Context, activity *models.Activity) error {
	_, err := r.collection.InsertOne(ctx, activity)
	return err
}

func (r *ActivityRepository) GetRecent(ctx context.Context, orgID primitive.ObjectID, limit int64) ([]*models.Activity, error) {
	filter := bson.M{"organization_id": orgID}
	opts := options.Find().SetSort(bson.M{"created_at": -1}).SetLimit(limit)

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to find activities: %w", err)
	}
	defer cursor.Close(ctx)

	var activities []*models.Activity
	if err := cursor.All(ctx, &activities); err != nil {
		return nil, fmt.Errorf("failed to decode activities: %w", err)
	}

	return activities, nil
}
