package repository

import (
	"context"
	"time"

	"github.com/yourusername/erp-system/services/notification-service/internal/models"
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
		collection: db.Collection("device_tokens"),
	}
}

// Register adds or updates a device token.
// Upsert based on token string.
func (r *DeviceRepository) Register(ctx context.Context, device *models.DeviceToken) error {
	filter := bson.M{"token": device.Token}
	update := bson.M{
		"$set": bson.M{
			"user_id":         device.UserID,
			"organization_id": device.OrganizationID,
			"platform":        device.Platform,
			"name":            device.Name,
			"last_used_at":    time.Now(),
			"updated_at":      time.Now(),
		},
		"$setOnInsert": bson.M{
			"_id":        primitive.NewObjectID(),
			"created_at": time.Now(),
		},
	}
	opts := options.Update().SetUpsert(true)
	_, err := r.collection.UpdateOne(ctx, filter, update, opts)
	return err
}

// FindByOrganizationID finds all device tokens for an organization
func (r *DeviceRepository) FindByOrganizationID(ctx context.Context, orgID primitive.ObjectID) ([]string, error) {
	filter := bson.M{"organization_id": orgID}
	// Projection: only return token string
	opts := options.Find().SetProjection(bson.M{"token": 1})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []struct {
		Token string `bson:"token"`
	}
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	tokens := make([]string, len(results))
	for i, r := range results {
		tokens[i] = r.Token
	}

	return tokens, nil
}

// Delete removes a token (e.g. on logout)
func (r *DeviceRepository) Delete(ctx context.Context, token string) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"token": token})
	return err
}

// PruneOldTokens removes tokens older than X days
func (r *DeviceRepository) PruneOldTokens(ctx context.Context, olderThan time.Duration) error {
	threshold := time.Now().Add(-olderThan)
	_, err := r.collection.DeleteMany(ctx, bson.M{"last_used_at": bson.M{"$lt": threshold}})
	return err
}
