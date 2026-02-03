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

type PlanRepository struct {
	collection *mongo.Collection
}

func NewPlanRepository(db *mongo.Database) *PlanRepository {
	return &PlanRepository{
		collection: db.Collection("subscription_plans"),
	}
}

// Create creates a new subscription plan
func (r *PlanRepository) Create(ctx context.Context, plan *models.SubscriptionPlan) error {
	plan.ID = primitive.NewObjectID()
	plan.CreatedAt = time.Now()
	plan.UpdatedAt = time.Now()

	_, err := r.collection.InsertOne(ctx, plan)
	if err != nil {
		return fmt.Errorf("failed to create plan: %w", err)
	}
	return nil
}

// FindByID finds a plan by ID
func (r *PlanRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*models.SubscriptionPlan, error) {
	var plan models.SubscriptionPlan
	filter := bson.M{
		"_id":        id,
		"deleted_at": nil,
	}

	err := r.collection.FindOne(ctx, filter).Decode(&plan)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find plan: %w", err)
	}
	return &plan, nil
}

// List returns all active plans
func (r *PlanRepository) List(ctx context.Context, tier string, isPublic *bool) ([]*models.SubscriptionPlan, error) {
	filter := bson.M{"deleted_at": nil}

	if tier != "" {
		filter["tier"] = tier
	}
	if isPublic != nil {
		filter["is_public"] = *isPublic
	}

	opts := options.Find().SetSort(bson.D{{Key: "display_order", Value: 1}})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to find plans: %w", err)
	}
	defer cursor.Close(ctx)

	var plans []*models.SubscriptionPlan
	if err := cursor.All(ctx, &plans); err != nil {
		return nil, fmt.Errorf("failed to decode plans: %w", err)
	}

	return plans, nil
}

// Update updates a plan
func (r *PlanRepository) Update(ctx context.Context, plan *models.SubscriptionPlan) error {
	plan.UpdatedAt = time.Now()
	plan.Version++

	filter := bson.M{
		"_id":     plan.ID,
		"version": plan.Version - 1,
	}

	update := bson.M{"$set": plan}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update plan: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("plan not found or version conflict")
	}

	return nil
}

// SoftDelete soft deletes a plan
func (r *PlanRepository) SoftDelete(ctx context.Context, id, deletedBy primitive.ObjectID) error {
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
		return fmt.Errorf("failed to soft delete plan: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("plan not found")
	}

	return nil
}
