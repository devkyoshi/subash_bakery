package repository

import (
	"context"

	"github.com/yourusername/erp-system/shared/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UnitRepository struct {
	collection      *mongo.Collection
	chartCollection *mongo.Collection
}

func NewUnitRepository(db *mongo.Database) *UnitRepository {
	return &UnitRepository{
		collection:      db.Collection("units"),
		chartCollection: db.Collection("unit_charts"),
	}
}

// FindByIDs retrieves units by multiple IDs
func (r *UnitRepository) FindByIDs(ctx context.Context, ids []primitive.ObjectID) ([]*models.Unit, error) {
	if len(ids) == 0 {
		return []*models.Unit{}, nil
	}

	filter := bson.M{
		"_id": bson.M{"$in": ids},
	}

	cursor, err := r.collection.Find(ctx, filter)
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

// ListUnits lists all active units
func (r *UnitRepository) ListUnits(ctx context.Context) ([]*models.Unit, error) {
	filter := bson.M{
		"is_active": true,
	}

	cursor, err := r.collection.Find(ctx, filter)
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

// ListUnitCharts lists all active unit conversion charts
func (r *UnitRepository) ListUnitCharts(ctx context.Context) ([]*models.UnitChart, error) {
	filter := bson.M{
		"is_active": true,
	}

	cursor, err := r.chartCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var charts []*models.UnitChart
	if err = cursor.All(ctx, &charts); err != nil {
		return nil, err
	}

	return charts, nil
}
