package repository

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/yourusername/erp-system/shared/models"
)

type StockLevelRepository struct {
	collection *mongo.Collection
}

func NewStockLevelRepository(db *mongo.Database) *StockLevelRepository {
	return &StockLevelRepository{
		collection: db.Collection("stock_levels"),
	}
}

func (r *StockLevelRepository) FindByProductAndLocation(ctx context.Context, productID, locationID primitive.ObjectID) (*models.StockLevel, error) {
	var stock models.StockLevel
	filter := bson.M{
		"product_id":  productID,
		"location_id": locationID,
	}
	err := r.collection.FindOne(ctx, filter).Decode(&stock)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &stock, nil
}

func (r *StockLevelRepository) FindByLocation(ctx context.Context, locationID primitive.ObjectID) ([]*models.StockLevel, error) {
	filter := bson.M{"location_id": locationID}
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var stocks []*models.StockLevel
	if err = cursor.All(ctx, &stocks); err != nil {
		return nil, err
	}
	return stocks, nil
}

func (r *StockLevelRepository) FindByProduct(ctx context.Context, productID primitive.ObjectID) ([]*models.StockLevel, error) {
	filter := bson.M{"product_id": productID}
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var stocks []*models.StockLevel
	if err = cursor.All(ctx, &stocks); err != nil {
		return nil, err
	}
	return stocks, nil
}

// FindByProducts retrieves stock levels for multiple products
// Returns a map of product_id+location_id -> StockLevel
func (r *StockLevelRepository) FindByProducts(ctx context.Context, productIDs []primitive.ObjectID) (map[string]*models.StockLevel, error) {
	filter := bson.M{
		"product_id": bson.M{"$in": productIDs},
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
		// Create composite key if needed, or just list them.
		// For simple aggregation, we might just return the list.
		// Use composite key: productID_locationID
		key := stock.ProductID.Hex() + "_" + stock.LocationID.Hex()
		stockMap[key] = &stock
	}

	return stockMap, nil
}

func (r *StockLevelRepository) Upsert(ctx context.Context, stock *models.StockLevel) error {
	now := time.Now()
	filter := bson.M{
		"organization_id": stock.OrganizationID,
		"product_id":      stock.ProductID,
		"location_id":     stock.LocationID,
	}

	update := bson.M{
		"$set": bson.M{
			"quantity_on_hand":    stock.QuantityOnHand,
			"quantity_available":  stock.QuantityAvailable,
			"quantity_allocated":  stock.QuantityAllocated,
			"quantity_in_transit": stock.QuantityInTransit,
			"quantity_reserved":   stock.QuantityReserved,
			"average_cost":        stock.AverageCost,
			"last_cost":           stock.LastCost,
			"total_value":         stock.TotalValue,
			"last_movement_date":  now,
			"updated_at":          now,
		},
		"$setOnInsert": bson.M{
			"_id":             primitive.NewObjectID(),
			"organization_id": stock.OrganizationID,
			"product_id":      stock.ProductID,
			"location_id":     stock.LocationID,
			"created_at":      now,
		},
	}

	opts := options.Update().SetUpsert(true)
	_, err := r.collection.UpdateOne(ctx, filter, update, opts)
	return err
}

func (r *StockLevelRepository) Find(ctx context.Context, filters map[string]interface{}, page, limit int) ([]models.StockLevel, error) {
	bsonFilters := bson.M{}
	for k, v := range filters {
		bsonFilters[k] = v
	}

	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}})
	if limit > 0 {
		opts.SetLimit(int64(limit))
		opts.SetSkip(int64((page - 1) * limit))
	}

	cursor, err := r.collection.Find(ctx, bsonFilters, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	// Initialize empty slice to ensure JSON array [] instead of null
	stocks := make([]models.StockLevel, 0)
	if err = cursor.All(ctx, &stocks); err != nil {
		return nil, err
	}
	return stocks, nil
}

func (r *StockLevelRepository) AdjustQuantity(ctx context.Context, orgID, productID, locationID primitive.ObjectID, delta float64, cost float64) error {
	filter := bson.M{
		"product_id":  productID,
		"location_id": locationID,
	}

	// Get current stock for weighted average calculation
	var currentStock models.StockLevel
	err := r.collection.FindOne(ctx, filter).Decode(&currentStock)
	if err != nil && err != mongo.ErrNoDocuments {
		return err
	}

	newQty := currentStock.QuantityOnHand + delta
	var newAvgCost float64

	// Calculate weighted average cost
	if delta > 0 && cost > 0 {
		totalValue := (currentStock.QuantityOnHand * currentStock.AverageCost) + (delta * cost)
		if newQty > 0 {
			newAvgCost = totalValue / newQty
		} else {
			newAvgCost = cost
		}
	} else {
		newAvgCost = currentStock.AverageCost
	}

	update := bson.M{
		"$inc": bson.M{
			"quantity_on_hand": delta,
		},
		"$set": bson.M{
			"quantity_available": newQty - currentStock.QuantityAllocated,
			"average_cost":       newAvgCost,
			"last_cost":          cost,
			"total_value":        newQty * newAvgCost,
			"last_movement_date": time.Now(),
			"updated_at":         time.Now(),
		},
		"$setOnInsert": bson.M{
			"_id":             primitive.NewObjectID(),
			"organization_id": orgID,
			"product_id":      productID,
			"location_id":     locationID,
			"created_at":      time.Now(),
		},
	}

	opts := options.Update().SetUpsert(true)
	_, err = r.collection.UpdateOne(ctx, filter, update, opts)
	return err
}

type StockMovementRepository struct {
	collection *mongo.Collection
}

func NewStockMovementRepository(db *mongo.Database) *StockMovementRepository {
	return &StockMovementRepository{
		collection: db.Collection("stock_movements"),
	}
}

func (r *StockMovementRepository) Create(ctx context.Context, movement *models.StockMovement) error {
	movement.ID = primitive.NewObjectID()
	movement.CreatedAt = time.Now()
	movement.UpdatedAt = time.Now()
	if movement.MovementDate.IsZero() {
		movement.MovementDate = time.Now()
	}

	_, err := r.collection.InsertOne(ctx, movement)
	return err
}

func (r *StockMovementRepository) FindByProduct(ctx context.Context, productID primitive.ObjectID, page, limit int) ([]*models.StockMovement, error) {
	filter := bson.M{"product_id": productID}
	opts := options.Find().SetSort(bson.D{{Key: "movement_date", Value: -1}})

	if limit > 0 {
		opts.SetLimit(int64(limit))
		opts.SetSkip(int64((page - 1) * limit))
	}

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var movements []*models.StockMovement
	if err = cursor.All(ctx, &movements); err != nil {
		return nil, err
	}
	return movements, nil
}

func (r *StockMovementRepository) FindByLocation(ctx context.Context, locationID primitive.ObjectID, startDate, endDate time.Time) ([]*models.StockMovement, error) {
	filter := bson.M{
		"$or": []bson.M{
			{"from_location_id": locationID},
			{"to_location_id": locationID},
		},
		"movement_date": bson.M{
			"$gte": startDate,
			"$lte": endDate,
		},
	}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var movements []*models.StockMovement
	if err = cursor.All(ctx, &movements); err != nil {
		return nil, err
	}
	return movements, nil
}

type BatchRepository struct {
	collection *mongo.Collection
}

func NewBatchRepository(db *mongo.Database) *BatchRepository {
	return &BatchRepository{
		collection: db.Collection("batches"),
	}
}

func (r *BatchRepository) Create(ctx context.Context, batch *models.Batch) error {
	batch.ID = primitive.NewObjectID()
	batch.CreatedAt = time.Now()
	batch.UpdatedAt = time.Now()
	batch.IsActive = true
	batch.CurrentQuantity = batch.InitialQuantity

	_, err := r.collection.InsertOne(ctx, batch)
	return err
}

func (r *BatchRepository) FindByProduct(ctx context.Context, productID, locationID primitive.ObjectID, activeOnly bool) ([]*models.Batch, error) {
	filter := bson.M{
		"product_id":  productID,
		"location_id": locationID,
	}
	if activeOnly {
		filter["is_active"] = true
		filter["current_quantity"] = bson.M{"$gt": 0}
	}

	opts := options.Find().SetSort(bson.D{{Key: "expiry_date", Value: 1}})
	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var batches []*models.Batch
	if err = cursor.All(ctx, &batches); err != nil {
		return nil, err
	}
	return batches, nil
}

func (r *BatchRepository) UpdateQuantity(ctx context.Context, batchID primitive.ObjectID, delta float64) error {
	filter := bson.M{"_id": batchID}
	update := bson.M{
		"$inc": bson.M{"current_quantity": delta},
		"$set": bson.M{"updated_at": time.Now()},
	}
	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

func (r *BatchRepository) FindBatchByID(ctx context.Context, id primitive.ObjectID) (*models.Batch, error) {
	var batch models.Batch
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&batch)
	if err != nil {
		return nil, err
	}
	return &batch, nil
}

// ==================== StockMovementRepository Additional Methods ====================

func (r *StockMovementRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*models.StockMovement, error) {
	var movement models.StockMovement
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&movement)
	if err != nil {
		return nil, err
	}
	return &movement, nil
}

func (r *StockMovementRepository) Find(ctx context.Context, filters map[string]interface{}, page, limit int) ([]*models.StockMovement, error) {
	skip := (page - 1) * limit
	opts := options.Find().
		SetSort(bson.D{{Key: "created_at", Value: -1}}).
		SetSkip(int64(skip)).
		SetLimit(int64(limit))

	// Convert filters map to bson.M for proper MongoDB querying
	bsonFilters := bson.M{}
	for k, v := range filters {
		bsonFilters[k] = v
	}

	cursor, err := r.collection.Find(ctx, bsonFilters, opts)
	if err != nil {
		return []*models.StockMovement{}, err
	}
	defer cursor.Close(ctx)

	// Initialize with empty slice instead of nil to avoid JSON null
	movements := make([]*models.StockMovement, 0)
	if err = cursor.All(ctx, &movements); err != nil {
		return []*models.StockMovement{}, err
	}

	return movements, nil
}
