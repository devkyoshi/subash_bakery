package repository

import (
	"context"
	"time"

	"github.com/yourusername/erp-system/shared/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UnitRepository struct {
	collection *mongo.Collection
}

func NewUnitRepository(db *mongo.Database) *UnitRepository {
	return &UnitRepository{
		collection: db.Collection("units"),
	}
}

func (r *UnitRepository) Create(ctx context.Context, unit *models.Unit) error {
	if unit.ID.IsZero() {
		unit.ID = primitive.NewObjectID()
	}
	if unit.CreatedAt.IsZero() {
		unit.CreatedAt = time.Now()
	}
	unit.UpdatedAt = time.Now()
	_, err := r.collection.InsertOne(ctx, unit)
	return err
}

func (r *UnitRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*models.Unit, error) {
	var unit models.Unit
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&unit)
	if err != nil {
		return nil, err
	}
	return &unit, nil
}

func (r *UnitRepository) Find(ctx context.Context, filters map[string]interface{}, activeOnly bool) ([]*models.Unit, error) {
	bsonFilters := bson.M{}
	for k, v := range filters {
		// specific handling for "ids" to map to "_id"
		if k == "ids" {
			if ids, ok := v.([]primitive.ObjectID); ok && len(ids) > 0 {
				bsonFilters["_id"] = bson.M{"$in": ids}
			}
			continue
		}
		bsonFilters[k] = v
	}
	if activeOnly {
		bsonFilters["is_active"] = true
	}

	cursor, err := r.collection.Find(ctx, bsonFilters)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var units []*models.Unit
	if err = cursor.All(ctx, &units); err != nil {
		return nil, err
	}
	return units, nil
}

func (r *UnitRepository) Update(ctx context.Context, id primitive.ObjectID, updates map[string]interface{}) error {
	updates["updated_at"] = time.Now()
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": updates})
	return err
}

func (r *UnitRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	// Soft delete
	update := bson.M{
		"$set": bson.M{
			"deleted_at": time.Now(),
		},
	}
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	return err
}
