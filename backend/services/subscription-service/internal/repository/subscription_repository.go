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

type SubscriptionRepository struct {
	collection *mongo.Collection
}

func NewSubscriptionRepository(db *mongo.Database) *SubscriptionRepository {
	return &SubscriptionRepository{
		collection: db.Collection("organization_subscriptions"),
	}
}

// Create creates a new subscription
func (r *SubscriptionRepository) Create(ctx context.Context, subscription *models.OrganizationSubscription) error {
	subscription.ID = primitive.NewObjectID()
	subscription.CreatedAt = time.Now()
	subscription.UpdatedAt = time.Now()

	_, err := r.collection.InsertOne(ctx, subscription)
	if err != nil {
		return fmt.Errorf("failed to create subscription: %w", err)
	}
	return nil
}

// FindByID finds a subscription by ID
func (r *SubscriptionRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*models.OrganizationSubscription, error) {
	var subscription models.OrganizationSubscription
	filter := bson.M{
		"_id":        id,
		"deleted_at": nil,
	}

	err := r.collection.FindOne(ctx, filter).Decode(&subscription)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find subscription: %w", err)
	}
	return &subscription, nil
}

// FindByOrganization finds the active subscription for an organization
func (r *SubscriptionRepository) FindByOrganization(ctx context.Context, orgID primitive.ObjectID) (*models.OrganizationSubscription, error) {
	var subscription models.OrganizationSubscription
	filter := bson.M{
		"organization_id": orgID,
		"deleted_at":      nil,
		"status": bson.M{
			"$in": []models.SubscriptionStatus{
				models.SubscriptionStatusActive,
				models.SubscriptionStatusTrial,
			},
		},
	}

	err := r.collection.FindOne(ctx, filter).Decode(&subscription)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find subscription: %w", err)
	}
	return &subscription, nil
}

// Update updates a subscription
func (r *SubscriptionRepository) Update(ctx context.Context, subscription *models.OrganizationSubscription) error {
	subscription.UpdatedAt = time.Now()
	subscription.Version++

	filter := bson.M{
		"_id":     subscription.ID,
		"version": subscription.Version - 1,
	}

	update := bson.M{"$set": subscription}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update subscription: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("subscription not found or version conflict")
	}

	return nil
}

// UpdateUsage updates subscription usage
func (r *SubscriptionRepository) UpdateUsage(ctx context.Context, id primitive.ObjectID, usage models.SubscriptionUsage) error {
	usage.LastUpdated = time.Now()

	filter := bson.M{"_id": id}
	update := bson.M{
		"$set": bson.M{
			"current_usage": usage,
			"updated_at":    time.Now(),
		},
	}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update subscription usage: %w", err)
	}
	return nil
}

// Cancel cancels a subscription
func (r *SubscriptionRepository) Cancel(ctx context.Context, id primitive.ObjectID, reason string) error {
	now := time.Now()
	filter := bson.M{"_id": id}
	update := bson.M{
		"$set": bson.M{
			"status":              models.SubscriptionStatusCancelled,
			"cancelled_at":        now,
			"cancellation_reason": reason,
			"auto_renew":          false,
			"updated_at":          now,
		},
	}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to cancel subscription: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("subscription not found")
	}

	return nil
}
