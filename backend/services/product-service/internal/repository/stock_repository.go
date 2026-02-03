package repository

import (
	"context"

	"github.com/yourusername/erp-system/shared/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type StockLevelRepository struct {
	collection *mongo.Collection
}

func NewStockLevelRepository(db *mongo.Database) *StockLevelRepository {
	return &StockLevelRepository{
		collection: db.Collection("stock_levels"),
	}
}

// FindByProductAndLocation retrieves stock level for a specific product at a location
func (r *StockLevelRepository) FindByProductAndLocation(ctx context.Context, productID, locationID primitive.ObjectID) (*models.StockLevel, error) {
	var stockLevel models.StockLevel
	filter := bson.M{
		"product_id":  productID,
		"location_id": locationID,
		"deleted_at":  nil,
	}

	err := r.collection.FindOne(ctx, filter).Decode(&stockLevel)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &stockLevel, nil
}

// FindByProducts retrieves stock levels for multiple products
// Returns a map of product_id+location_id -> StockLevel
func (r *StockLevelRepository) FindByProducts(ctx context.Context, productIDs []primitive.ObjectID) (map[string]*models.StockLevel, error) {
	filter := bson.M{
		"product_id": bson.M{"$in": productIDs},
		"deleted_at": nil,
	}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	stockMap := make(map[string]*models.StockLevel)
	for cursor.Next(ctx) {
		var stock models.StockLevel
		if err := cursor.Decode(&stock); err != nil {
			continue
		}
		// Create composite key: productID_locationID
		key := stock.ProductID.Hex() + "_" + stock.LocationID.Hex()
		stockMap[key] = &stock
	}

	return stockMap, nil
}
