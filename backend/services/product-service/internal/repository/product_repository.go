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

type ProductRepository struct {
	collection *mongo.Collection
}

func NewProductRepository(db *mongo.Database) *ProductRepository {
	return &ProductRepository{
		collection: db.Collection("products"),
	}
}

// Create inserts a new product
func (r *ProductRepository) Create(ctx context.Context, product *models.Product) error {
	product.ID = primitive.NewObjectID()
	product.CreatedAt = time.Now()
	product.UpdatedAt = time.Now()

	_, err := r.collection.InsertOne(ctx, product)
	return err
}

// FindByID retrieves a product by ID
func (r *ProductRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*models.Product, error) {
	var product models.Product
	filter := bson.M{
		"_id":        id,
		"deleted_at": nil,
	}

	err := r.collection.FindOne(ctx, filter).Decode(&product)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &product, nil
}

// FindBySKU retrieves a product by SKU
func (r *ProductRepository) FindBySKU(ctx context.Context, orgID primitive.ObjectID, sku string) (*models.Product, error) {
	var product models.Product
	filter := bson.M{
		"organization_id": orgID,
		"sku":             sku,
		"deleted_at":      nil,
	}

	err := r.collection.FindOne(ctx, filter).Decode(&product)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &product, nil
}

// FindByOrganization retrieves products for an organization with filters
func (r *ProductRepository) FindByOrganization(ctx context.Context, orgID primitive.ObjectID, filters map[string]interface{}, page, limit int) ([]*models.Product, int64, error) {
	filter := bson.M{
		"organization_id": orgID,
		"deleted_at":      nil,
	}

	// Apply additional filters
	if categoryID, ok := filters["category_id"].(primitive.ObjectID); ok {
		filter["category_id"] = categoryID
	}
	if subcategoryID, ok := filters["subcategory_id"].(primitive.ObjectID); ok {
		filter["subcategory_id"] = subcategoryID
	}
	if brandID, ok := filters["brand_id"].(primitive.ObjectID); ok {
		filter["brand_id"] = brandID
	}
	if status, ok := filters["status"].(models.ProductStatus); ok {
		filter["status"] = status
	}
	if productType, ok := filters["type"].(models.ProductType); ok {
		filter["type"] = productType
	}
	if trackInventory, ok := filters["track_inventory"].(bool); ok {
		filter["track_inventory"] = trackInventory
	}

	// Filter by location ID - products that have prices for this location
	if locationID, ok := filters["location_id"].(primitive.ObjectID); ok {
		filter["location_prices.location_id"] = locationID
	}

	// Search by name or SKU
	if search, ok := filters["search"].(string); ok && search != "" {
		filter["$or"] = []bson.M{
			{"name": bson.M{"$regex": search, "$options": "i"}},
			{"sku": bson.M{"$regex": search, "$options": "i"}},
			{"description": bson.M{"$regex": search, "$options": "i"}},
		}
	}

	// Count total
	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	// Setup pagination
	skip := (page - 1) * limit
	opts := options.Find().
		SetSkip(int64(skip)).
		SetLimit(int64(limit)).
		SetSort(bson.D{{Key: "name", Value: 1}})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var products []*models.Product
	if err = cursor.All(ctx, &products); err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

// Update updates an existing product
func (r *ProductRepository) Update(ctx context.Context, product *models.Product) error {
	product.UpdatedAt = time.Now()

	filter := bson.M{
		"_id":        product.ID,
		"deleted_at": nil,
	}

	update := bson.M{
		"$set": product,
	}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

// Delete soft deletes a product
func (r *ProductRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	now := time.Now()

	filter := bson.M{
		"_id":        id,
		"deleted_at": nil,
	}

	update := bson.M{
		"$set": bson.M{
			"deleted_at": now,
			"updated_at": now,
		},
	}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

// CheckSKUExists checks if a SKU already exists
func (r *ProductRepository) CheckSKUExists(ctx context.Context, orgID primitive.ObjectID, sku string, excludeID *primitive.ObjectID) (bool, error) {
	filter := bson.M{
		"organization_id": orgID,
		"sku":             sku,
		"deleted_at":      nil,
	}

	if excludeID != nil {
		filter["_id"] = bson.M{"$ne": *excludeID}
	}

	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// UpdateStockLevels updates stock levels for a product
func (r *ProductRepository) UpdateStockLevels(ctx context.Context, productID primitive.ObjectID, totalStock, availableStock, allocatedStock, inTransitStock float64) error {
	filter := bson.M{
		"_id":        productID,
		"deleted_at": nil,
	}

	update := bson.M{
		"$set": bson.M{
			"total_stock":      totalStock,
			"available_stock":  availableStock,
			"allocated_stock":  allocatedStock,
			"in_transit_stock": inTransitStock,
			"updated_at":       time.Now(),
		},
	}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

// GetLowStockProducts retrieves products below reorder level
func (r *ProductRepository) GetLowStockProducts(ctx context.Context, orgID primitive.ObjectID) ([]*models.Product, error) {
	filter := bson.M{
		"organization_id": orgID,
		"deleted_at":      nil,
		"track_inventory": true,
		"$expr": bson.M{
			"$lte": []interface{}{"$available_stock", "$reorder_level"},
		},
	}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var products []*models.Product
	if err = cursor.All(ctx, &products); err != nil {
		return nil, err
	}

	return products, nil
}

// CountByCategory counts products in a category
func (r *ProductRepository) CountByCategory(ctx context.Context, categoryID primitive.ObjectID) (int64, error) {
	filter := bson.M{
		"category_id": categoryID,
		"deleted_at":  nil,
	}

	return r.collection.CountDocuments(ctx, filter)
}

// CountBySubcategory counts products in a subcategory
func (r *ProductRepository) CountBySubcategory(ctx context.Context, subcategoryID primitive.ObjectID) (int64, error) {
	filter := bson.M{
		"subcategory_id": subcategoryID,
		"deleted_at":     nil,
	}

	return r.collection.CountDocuments(ctx, filter)
}

// CountByBrand counts products for a brand
func (r *ProductRepository) CountByBrand(ctx context.Context, brandID primitive.ObjectID) (int64, error) {
	filter := bson.M{
		"brand_id":   brandID,
		"deleted_at": nil,
	}

	return r.collection.CountDocuments(ctx, filter)
}

// Count counts total products for an organization
func (r *ProductRepository) Count(ctx context.Context, orgID primitive.ObjectID) (int64, error) {
	filter := bson.M{
		"organization_id": orgID,
		"deleted_at":      nil,
	}
	return r.collection.CountDocuments(ctx, filter)
}
