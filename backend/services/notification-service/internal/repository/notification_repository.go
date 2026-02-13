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

type NotificationRepository struct {
	collection *mongo.Collection
}

func NewNotificationRepository(db *mongo.Database) *NotificationRepository {
	return &NotificationRepository{
		collection: db.Collection("notifications"),
	}
}

func (r *NotificationRepository) Create(ctx context.Context, notification *models.Notification) error {
	notification.CreatedAt = time.Now()
	notification.IsRead = false
	result, err := r.collection.InsertOne(ctx, notification)
	if err != nil {
		return err
	}
	notification.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *NotificationRepository) FindByUserID(ctx context.Context, userID primitive.ObjectID, limit int64) ([]*models.Notification, error) {
	filter := bson.M{"user_id": userID}
	opts := options.Find().SetSort(bson.M{"created_at": -1}).SetLimit(limit)

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var notifications []*models.Notification
	if err = cursor.All(ctx, &notifications); err != nil {
		return nil, err
	}
	return notifications, nil
}

func (r *NotificationRepository) MarkAsRead(ctx context.Context, id primitive.ObjectID, userID primitive.ObjectID) error {
	filter := bson.M{"_id": id, "user_id": userID}
	update := bson.M{"$set": bson.M{"is_read": true}}
	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

func (r *NotificationRepository) MarkAllAsRead(ctx context.Context, userID primitive.ObjectID) error {
	filter := bson.M{"user_id": userID, "is_read": false}
	update := bson.M{"$set": bson.M{"is_read": true}}
	_, err := r.collection.UpdateMany(ctx, filter, update)
	return err
}
