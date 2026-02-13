package repository

import (
	"context"
	"fmt"

	sharedModels "github.com/yourusername/erp-system/shared/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/yourusername/erp-system/services/report-service/internal/models"
)

type InventoryReportRepository struct {
	db             *mongo.Database
	stockLevelColl *mongo.Collection
	productColl    *mongo.Collection
	categoryColl   *mongo.Collection
	locationColl   *mongo.Collection
	unitColl       *mongo.Collection
}

func NewInventoryReportRepository(db *mongo.Database) *InventoryReportRepository {
	return &InventoryReportRepository{
		db:             db,
		stockLevelColl: db.Collection("stock_levels"),
		productColl:    db.Collection("products"),
		categoryColl:   db.Collection("categories"),
		locationColl:   db.Collection("locations"),
		unitColl:       db.Collection("units"),
	}
}

// GetStockLevels returns paginated stock levels for an organization
func (r *InventoryReportRepository) GetStockLevels(
	ctx context.Context,
	orgID primitive.ObjectID,
	filters models.StockLevelFilters,
	page, limit int,
) ([]*sharedModels.StockLevel, int64, error) {
	filter := r.buildStockLevelFilter(ctx, orgID, filters)

	total, err := r.stockLevelColl.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count stock levels: %w", err)
	}

	skip := int64((page - 1) * limit)
	opts := options.Find().
		SetSort(bson.M{"quantity_on_hand": -1}).
		SetSkip(skip).
		SetLimit(int64(limit))

	cursor, err := r.stockLevelColl.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to find stock levels: %w", err)
	}
	defer cursor.Close(ctx)

	var levels []*sharedModels.StockLevel
	if err := cursor.All(ctx, &levels); err != nil {
		return nil, 0, fmt.Errorf("failed to decode stock levels: %w", err)
	}

	return levels, total, nil
}

// GetAllStockLevelsForReport returns all stock levels (no pagination) for export
func (r *InventoryReportRepository) GetAllStockLevelsForReport(
	ctx context.Context,
	orgID primitive.ObjectID,
	filters models.StockLevelFilters,
) ([]*sharedModels.StockLevel, error) {
	filter := r.buildStockLevelFilter(ctx, orgID, filters)

	opts := options.Find().SetSort(bson.M{"quantity_on_hand": -1})

	cursor, err := r.stockLevelColl.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to find stock levels: %w", err)
	}
	defer cursor.Close(ctx)

	var levels []*sharedModels.StockLevel
	if err := cursor.All(ctx, &levels); err != nil {
		return nil, fmt.Errorf("failed to decode stock levels: %w", err)
	}

	return levels, nil
}

// GetProductsByIDs returns products for the given IDs
func (r *InventoryReportRepository) GetProductsByIDs(
	ctx context.Context,
	productIDs []primitive.ObjectID,
) (map[string]*sharedModels.Product, error) {
	if len(productIDs) == 0 {
		return make(map[string]*sharedModels.Product), nil
	}

	filter := bson.M{
		"_id":        bson.M{"$in": productIDs},
		"deleted_at": bson.M{"$exists": false},
	}

	cursor, err := r.productColl.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to find products: %w", err)
	}
	defer cursor.Close(ctx)

	result := make(map[string]*sharedModels.Product)
	for cursor.Next(ctx) {
		var product sharedModels.Product
		if err := cursor.Decode(&product); err != nil {
			continue
		}
		result[product.ID.Hex()] = &product
	}

	return result, nil
}

// GetCategoryNames returns a map of category ID -> category name
func (r *InventoryReportRepository) GetCategoryNames(
	ctx context.Context,
	categoryIDs []primitive.ObjectID,
) (map[string]string, error) {
	if len(categoryIDs) == 0 {
		return make(map[string]string), nil
	}

	filter := bson.M{"_id": bson.M{"$in": categoryIDs}}
	cursor, err := r.categoryColl.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to find categories: %w", err)
	}
	defer cursor.Close(ctx)

	result := make(map[string]string)
	for cursor.Next(ctx) {
		var cat struct {
			ID   primitive.ObjectID `bson:"_id"`
			Name string             `bson:"name"`
		}
		if err := cursor.Decode(&cat); err != nil {
			continue
		}
		result[cat.ID.Hex()] = cat.Name
	}

	return result, nil
}

// GetLocationNames returns a map of location ID -> location name
func (r *InventoryReportRepository) GetLocationNames(
	ctx context.Context,
	locationIDs []primitive.ObjectID,
) (map[string]string, error) {
	if len(locationIDs) == 0 {
		return make(map[string]string), nil
	}

	filter := bson.M{"_id": bson.M{"$in": locationIDs}}
	cursor, err := r.locationColl.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to find locations: %w", err)
	}
	defer cursor.Close(ctx)

	result := make(map[string]string)
	for cursor.Next(ctx) {
		var loc struct {
			ID   primitive.ObjectID `bson:"_id"`
			Name string             `bson:"name"`
		}
		if err := cursor.Decode(&loc); err != nil {
			continue
		}
		result[loc.ID.Hex()] = loc.Name
	}

	return result, nil
}

// GetUnitName returns a unit's abbreviation by ID
func (r *InventoryReportRepository) GetUnitName(
	ctx context.Context,
	unitID primitive.ObjectID,
) string {
	if unitID.IsZero() {
		return ""
	}

	var unit struct {
		Abbreviation string `bson:"abbreviation"`
		Name         string `bson:"name"`
	}
	err := r.unitColl.FindOne(ctx, bson.M{"_id": unitID}).Decode(&unit)
	if err != nil {
		return ""
	}
	if unit.Abbreviation != "" {
		return unit.Abbreviation
	}
	return unit.Name
}

// GetProductIDsByFilters returns product IDs matching search/category filters
func (r *InventoryReportRepository) GetProductIDsByFilters(
	ctx context.Context,
	orgID primitive.ObjectID,
	filters models.StockLevelFilters,
) ([]primitive.ObjectID, error) {
	filter := bson.M{
		"organization_id": orgID,
		"deleted_at":      bson.M{"$exists": false},
	}

	if filters.CategoryID != "" {
		catID, err := primitive.ObjectIDFromHex(filters.CategoryID)
		if err == nil {
			filter["category_id"] = catID
		}
	}

	if filters.Search != "" {
		filter["$or"] = []bson.M{
			{"name": bson.M{"$regex": filters.Search, "$options": "i"}},
			{"sku": bson.M{"$regex": filters.Search, "$options": "i"}},
		}
	}

	cursor, err := r.productColl.Find(ctx, filter, options.Find().SetProjection(bson.M{"_id": 1}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var ids []primitive.ObjectID
	for cursor.Next(ctx) {
		var doc struct {
			ID primitive.ObjectID `bson:"_id"`
		}
		if err := cursor.Decode(&doc); err != nil {
			continue
		}
		ids = append(ids, doc.ID)
	}

	return ids, nil
}

func (r *InventoryReportRepository) buildStockLevelFilter(
	ctx context.Context,
	orgID primitive.ObjectID,
	filters models.StockLevelFilters,
) bson.M {
	filter := bson.M{
		"organization_id": orgID,
		"deleted_at":      bson.M{"$exists": false},
	}

	if filters.LocationID != "" {
		locID, err := primitive.ObjectIDFromHex(filters.LocationID)
		if err == nil {
			filter["location_id"] = locID
		}
	}

	// If there's a search or category filter, first resolve matching product IDs
	if filters.Search != "" || filters.CategoryID != "" {
		productIDs, err := r.GetProductIDsByFilters(ctx, orgID, filters)
		if err == nil && len(productIDs) > 0 {
			filter["product_id"] = bson.M{"$in": productIDs}
		} else if filters.Search != "" || filters.CategoryID != "" {
			// No matching products found — return empty result
			filter["product_id"] = primitive.NilObjectID
		}
	}

	return filter
}
