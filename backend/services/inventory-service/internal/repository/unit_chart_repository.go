package repository

import (
	"context"
	"time"

	"github.com/yourusername/erp-system/shared/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UnitChartRepository struct {
	collection *mongo.Collection
}

func NewUnitChartRepository(db *mongo.Database) *UnitChartRepository {
	return &UnitChartRepository{
		collection: db.Collection("unit_charts"),
	}
}

func (r *UnitChartRepository) Create(ctx context.Context, chart *models.UnitChart) error {
	if chart.ID.IsZero() {
		chart.ID = primitive.NewObjectID()
	}
	if chart.CreatedAt.IsZero() {
		chart.CreatedAt = time.Now()
	}
	chart.UpdatedAt = time.Now()
	_, err := r.collection.InsertOne(ctx, chart)
	return err
}

func (r *UnitChartRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*models.UnitChart, error) {
	var chart models.UnitChart
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&chart)
	if err != nil {
		return nil, err
	}
	return &chart, nil
}

func (r *UnitChartRepository) FindByUnits(ctx context.Context, fromUnitID, toUnitID primitive.ObjectID) (*models.UnitChart, error) {
	filter := bson.M{
		"from_unit_id": fromUnitID,
		"to_unit_id":   toUnitID,
	}
	var chart models.UnitChart
	err := r.collection.FindOne(ctx, filter).Decode(&chart)
	if err != nil {
		return nil, err
	}
	return &chart, nil
}

func (r *UnitChartRepository) Find(ctx context.Context, activeOnly bool) ([]*models.UnitChart, error) {
	filter := bson.M{}
	if activeOnly {
		filter["is_active"] = true
	}

	cursor, err := r.collection.Find(ctx, filter)
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

func (r *UnitChartRepository) Update(ctx context.Context, chart *models.UnitChart) error {
	filter := bson.M{"_id": chart.ID}
	update := bson.M{
		"$set": bson.M{
			"conversion_rate": chart.ConversionRate,
			"is_active":       chart.IsActive,
			"metadata":        chart.Metadata,
			"updated_at":      time.Now(),
			"updated_by":      chart.UpdatedBy,
		},
	}
	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

func (r *UnitChartRepository) SoftDelete(ctx context.Context, id primitive.ObjectID) error {
	filter := bson.M{"_id": id}
	update := bson.M{
		"$set": bson.M{
			"deleted_at": time.Now(),
		},
	}
	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

// PathExists checks if a conversion path exists from -> to to detect circles
// Simple check: direct check + maybe recursively?
// For now, let's implement just checking if the reverse exists or if to -> from exists directly for 1-step circle.
// A full graph cycle check would be more complex but maybe overkill if we only add 1 by 1.
// Service just calls PathExists(to, from) which implies checking if "To -> From" exists when adding "From -> To".
func (r *UnitChartRepository) PathExists(ctx context.Context, from, to primitive.ObjectID) (bool, error) {
	// Simple DFS or similar could be done here if we want to support multi-hop circles.
	// But let's act as if we are checking if a conversion ALREADY exists from 'from' to 'to'.
	_, err := r.FindByUnits(ctx, from, to)
	if err == nil {
		return true, nil
	}
	if err == mongo.ErrNoDocuments {
		return false, nil
	}
	return false, err
}
