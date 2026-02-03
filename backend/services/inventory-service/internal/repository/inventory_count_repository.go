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

type InventoryCountRepository struct {
	collection *mongo.Collection
}

func NewInventoryCountRepository(db *mongo.Database) *InventoryCountRepository {
	return &InventoryCountRepository{
		collection: db.Collection("inventory_counts"),
	}
}

func (r *InventoryCountRepository) Create(ctx context.Context, count *models.InventoryCount) error {
	count.ID = primitive.NewObjectID()
	if count.CreatedAt.IsZero() {
		count.CreatedAt = time.Now()
	}
	count.UpdatedAt = time.Now()
	_, err := r.collection.InsertOne(ctx, count)
	return err
}

func (r *InventoryCountRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*models.InventoryCount, error) {
	var count models.InventoryCount
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&count)
	if err != nil {
		return nil, err
	}
	return &count, nil
}

func (r *InventoryCountRepository) FindByOrganization(ctx context.Context, orgID primitive.ObjectID, filters map[string]interface{}, page, limit int) ([]*models.InventoryCount, error) {
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

	var counts []*models.InventoryCount
	if err = cursor.All(ctx, &counts); err != nil {
		return nil, err
	}
	return counts, nil
}

func (r *InventoryCountRepository) UpdateItem(ctx context.Context, countID primitive.ObjectID, item models.InventoryCountItem) error {
	// Try to update using a pull and push strategy to ensure clean update of the item in array
	// Simpler heuristic: Pull the item with same ProductID (and BatchID if applicable) first, then push.

	pullFilter := bson.M{
		"_id": countID,
	}
	pullUpdate := bson.M{
		"$pull": bson.M{
			"items": bson.M{"product_id": item.ProductID},
		},
	}
	if item.BatchID != nil {
		pullUpdate = bson.M{
			"$pull": bson.M{
				"items": bson.M{"product_id": item.ProductID, "batch_id": item.BatchID},
			},
		}
	} else {
		pullUpdate = bson.M{
			"$pull": bson.M{
				"items": bson.M{"product_id": item.ProductID, "batch_id": nil},
			},
		}
	}

	_, _ = r.collection.UpdateOne(ctx, pullFilter, pullUpdate)

	// Now push
	pushUpdate := bson.M{
		"$push": bson.M{"items": item},
		"$set":  bson.M{"updated_at": time.Now()},
	}
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": countID}, pushUpdate)
	return err
}

func (r *InventoryCountRepository) Complete(ctx context.Context, countID, completedBy primitive.ObjectID, summary map[string]interface{}) error {
	update := bson.M{
		"$set": bson.M{
			"status":       "completed",
			"completed_by": completedBy,
			"completed_at": time.Now(),
			"updated_at":   time.Now(),
			"summary":      summary,
		},
	}
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": countID}, update)
	return err
}

func (r *InventoryCountRepository) Update(ctx context.Context, countID primitive.ObjectID, updates map[string]interface{}) error {
	updates["updated_at"] = time.Now()
	update := bson.M{"$set": updates}
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": countID}, update)
	return err
}

func (r *InventoryCountRepository) Delete(ctx context.Context, id, deletedBy primitive.ObjectID) error {
	// Soft delete
	update := bson.M{
		"$set": bson.M{
			"deleted_at": time.Now(),
			"deleted_by": deletedBy,
		},
	}
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	return err
	// Or hard delete if intended? Service check says "DeleteCount". Usually soft delete is preferred.
	// But let's stick to soft delete based on "deletedBy" param.
}
