package repository

import (
	"context"
	"time"

	"github.com/yourusername/erp-system/shared/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type StockAdjustmentRepository struct {
	collection *mongo.Collection
}

func NewStockAdjustmentRepository(db *mongo.Database) *StockAdjustmentRepository {
	return &StockAdjustmentRepository{
		collection: db.Collection("stock_adjustments"),
	}
}

func (r *StockAdjustmentRepository) Create(ctx context.Context, adjustment *models.StockAdjustment) error {
	adjustment.ID = primitive.NewObjectID()
	if adjustment.CreatedAt.IsZero() {
		adjustment.CreatedAt = time.Now()
	}
	adjustment.UpdatedAt = time.Now()
	_, err := r.collection.InsertOne(ctx, adjustment)
	return err
}

func (r *StockAdjustmentRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*models.StockAdjustment, error) {
	var adjustment models.StockAdjustment
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&adjustment)
	if err != nil {
		return nil, err
	}
	return &adjustment, nil
}

func (r *StockAdjustmentRepository) FindByOrganization(ctx context.Context, orgID primitive.ObjectID, filters map[string]interface{}, page, limit int) ([]*models.StockAdjustment, error) {
	bsonFilters := bson.M{"organization_id": orgID}
	for k, v := range filters {
		bsonFilters[k] = v
	}

	opts := options.Find().
		SetSort(bson.D{{Key: "created_at", Value: -1}}).
		SetSkip(int64((page - 1) * limit)).
		SetLimit(int64(limit))

	cursor, err := r.collection.Find(ctx, bsonFilters, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var adjustments []*models.StockAdjustment
	if err = cursor.All(ctx, &adjustments); err != nil {
		return nil, err
	}
	return adjustments, nil
}

func (r *StockAdjustmentRepository) Update(ctx context.Context, id primitive.ObjectID, updates map[string]interface{}) error {
	updates["updated_at"] = time.Now()
	update := bson.M{
		"$set": updates,
	}
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	return err
}

func (r *StockAdjustmentRepository) UpdateStatus(ctx context.Context, id primitive.ObjectID, status string, updatedBy primitive.ObjectID) error {
	update := bson.M{
		"$set": bson.M{
			"status":     status,
			"updated_by": updatedBy,
			"updated_at": time.Now(),
		},
	}
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	return err
}

func (r *StockAdjustmentRepository) Delete(ctx context.Context, id, deletedBy primitive.ObjectID) error {
	update := bson.M{
		"$set": bson.M{
			"deleted_at": time.Now(),
			"deleted_by": deletedBy,
		},
	}
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	return err
}
