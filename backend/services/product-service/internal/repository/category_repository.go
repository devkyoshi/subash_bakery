package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/yourusername/erp-system/shared/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CategoryRepository struct {
	collection        *mongo.Collection
	productCollection *mongo.Collection
}

func NewCategoryRepository(db *mongo.Database) *CategoryRepository {
	return &CategoryRepository{
		collection:        db.Collection("categories"),
		productCollection: db.Collection("products"),
	}
}

// Create inserts a new category
func (r *CategoryRepository) Create(ctx context.Context, category *models.ProductCategory) error {
	category.ID = primitive.NewObjectID()
	category.CreatedAt = time.Now()
	category.UpdatedAt = time.Now()
	category.ProductCount = 0

	// Initialize subcategories with IDs and timestamps
	if len(category.Subcategories) > 0 {
		initialized := make([]models.ProductSubcategory, 0, len(category.Subcategories))
		for _, sub := range category.Subcategories {
			sub.ID = primitive.NewObjectID()
			sub.CreatedAt = time.Now()
			sub.UpdatedAt = time.Now()
			sub.DeletedAt = nil
			sub.ProductCount = 0
			initialized = append(initialized, sub)
		}
		category.Subcategories = initialized
	}

	_, err := r.collection.InsertOne(ctx, category)
	return err
}

// FindByID retrieves a category by ID
func (r *CategoryRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*models.ProductCategory, error) {
	var category models.ProductCategory
	filter := bson.M{
		"_id":        id,
		"deleted_at": nil,
	}

	err := r.collection.FindOne(ctx, filter).Decode(&category)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &category, nil
}

// FindParentBySubcategoryID retrieves the parent category for a given subcategory ID
func (r *CategoryRepository) FindParentBySubcategoryID(ctx context.Context, subcategoryID primitive.ObjectID) (*models.ProductCategory, error) {
	var category models.ProductCategory
	filter := bson.M{
		"subcategories._id": subcategoryID,
		"deleted_at":        nil,
	}

	err := r.collection.FindOne(ctx, filter).Decode(&category)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &category, nil
}

// FindByOrganization retrieves categories for an organization with optional filters
func (r *CategoryRepository) FindByOrganization(ctx context.Context, orgID primitive.ObjectID, isActive *bool, query string, page, limit int) ([]*models.ProductCategory, int64, error) {
	filter := bson.M{
		"organization_id": orgID,
		"deleted_at":      nil,
	}

	if isActive != nil {
		filter["is_active"] = *isActive
	}

	// Add search query filter
	if query != "" {
		filter["$or"] = []bson.M{
			{"name": bson.M{"$regex": query, "$options": "i"}},
			{"code": bson.M{"$regex": query, "$options": "i"}},
			{"description": bson.M{"$regex": query, "$options": "i"}},
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
		SetSort(bson.D{{Key: "level", Value: 1}, {Key: "name", Value: 1}})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var categories []*models.ProductCategory
	if err = cursor.All(ctx, &categories); err != nil {
		return nil, 0, err
	}

	return categories, total, nil
}

// FindRootCategories retrieves all root-level categories (level 0)
// FindByPath retrieves a category by its path (kept for compatibility if needed)
func (r *CategoryRepository) FindByPath(ctx context.Context, orgID primitive.ObjectID, path string) (*models.ProductCategory, error) {
	var category models.ProductCategory
	filter := bson.M{
		"organization_id": orgID,
		"path":            path,
		"deleted_at":      nil,
	}

	err := r.collection.FindOne(ctx, filter).Decode(&category)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &category, nil
}

// Update updates an existing category
func (r *CategoryRepository) Update(ctx context.Context, category *models.ProductCategory) error {
	category.UpdatedAt = time.Now()

	// Ensure subcategories have updated timestamps when modified
	if len(category.Subcategories) > 0 {
		for i := range category.Subcategories {
			category.Subcategories[i].UpdatedAt = time.Now()
		}
	}

	filter := bson.M{
		"_id":        category.ID,
		"deleted_at": nil,
	}

	update := bson.M{
		"$set": category,
	}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

// Delete soft deletes a category
func (r *CategoryRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
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

// HasChildren checks if a category has any subcategories (embedded)
func (r *CategoryRepository) HasChildren(ctx context.Context, id primitive.ObjectID) (bool, error) {
	category, err := r.FindByID(ctx, id)
	if err != nil {
		return false, err
	}
	if category == nil {
		return false, nil
	}

	// Check if subcategories array is non-empty
	for _, sub := range category.Subcategories {
		if sub.DeletedAt == nil {
			return true, nil
		}
	}

	return false, nil
}

// HasProducts checks if a category has any products
func (r *CategoryRepository) HasProducts(ctx context.Context, categoryID primitive.ObjectID) (bool, error) {
	filter := bson.M{
		"category_id": categoryID,
		"deleted_at":  nil,
	}

	count, err := r.productCollection.CountDocuments(ctx, filter)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// HasActiveProducts checks if a category has any active products
func (r *CategoryRepository) HasActiveProducts(ctx context.Context, categoryID primitive.ObjectID) (bool, error) {
	filter := bson.M{
		"category_id": categoryID,
		"is_active":   true,
		"deleted_at":  nil,
	}

	count, err := r.productCollection.CountDocuments(ctx, filter)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// UpdateProductCount updates the product count for a category
func (r *CategoryRepository) UpdateProductCount(ctx context.Context, categoryID primitive.ObjectID) error {
	filter := bson.M{
		"category_id": categoryID,
		"deleted_at":  nil,
	}

	count, err := r.productCollection.CountDocuments(ctx, filter)
	if err != nil {
		return err
	}

	update := bson.M{
		"$set": bson.M{
			"product_count": count,
			"updated_at":    time.Now(),
		},
	}

	_, err = r.collection.UpdateOne(ctx, bson.M{"_id": categoryID}, update)
	return err
}

// GetCategoryTree retrieves all categories for an organization (subcategories are embedded)
func (r *CategoryRepository) GetCategoryTree(ctx context.Context, orgID primitive.ObjectID) ([]*models.ProductCategory, error) {
	filter := bson.M{
		"organization_id": orgID,
		"deleted_at":      nil,
	}

	opts := options.Find().SetSort(bson.D{{Key: "name", Value: 1}})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var categories []*models.ProductCategory
	if err = cursor.All(ctx, &categories); err != nil {
		return nil, err
	}

	return categories, nil
}

// CheckNameExists checks if a category name exists in the organization at the top level
func (r *CategoryRepository) CheckNameExists(ctx context.Context, orgID primitive.ObjectID, name string, excludeID *primitive.ObjectID) (bool, error) {
	filter := bson.M{
		"organization_id": orgID,
		"name":            name,
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

// CheckSubcategoryNameExists checks if a subcategory name exists within a specific category
func (r *CategoryRepository) CheckSubcategoryNameExists(ctx context.Context, categoryID primitive.ObjectID, name string, excludeSubcategoryID *primitive.ObjectID) (bool, error) {
	category, err := r.FindByID(ctx, categoryID)
	if err != nil {
		return false, err
	}
	if category == nil {
		return false, fmt.Errorf("category not found")
	}

	for _, sub := range category.Subcategories {
		if sub.DeletedAt != nil {
			continue
		}
		if sub.Name == name {
			if excludeSubcategoryID == nil || sub.ID != *excludeSubcategoryID {
				return true, nil
			}
		}
	}

	return false, nil
}

// FindByIDs retrieves multiple categories by their IDs
func (r *CategoryRepository) FindByIDs(ctx context.Context, ids []primitive.ObjectID) ([]*models.ProductCategory, error) {
	filter := bson.M{
		"_id":        bson.M{"$in": ids},
		"deleted_at": nil,
	}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var categories []*models.ProductCategory
	if err = cursor.All(ctx, &categories); err != nil {
		return nil, err
	}

	return categories, nil
}

// HasProductsInCategories checks if any of the given categories have products
func (r *CategoryRepository) HasProductsInCategories(ctx context.Context, categoryIDs []primitive.ObjectID) (map[string]bool, error) {
	result := make(map[string]bool)

	for _, categoryID := range categoryIDs {
		hasProducts, err := r.HasProducts(ctx, categoryID)
		if err != nil {
			return nil, err
		}
		result[categoryID.Hex()] = hasProducts
	}

	return result, nil
}

// DeleteMultiple soft deletes multiple categories
func (r *CategoryRepository) DeleteMultiple(ctx context.Context, ids []primitive.ObjectID) error {
	now := time.Now()

	filter := bson.M{
		"_id":        bson.M{"$in": ids},
		"deleted_at": nil,
	}

	update := bson.M{
		"$set": bson.M{
			"deleted_at": now,
			"updated_at": now,
		},
	}

	_, err := r.collection.UpdateMany(ctx, filter, update)
	return err
}

// Count counts total categories for an organization
func (r *CategoryRepository) Count(ctx context.Context, orgID primitive.ObjectID) (int64, error) {
	filter := bson.M{
		"organization_id": orgID,
		"deleted_at":      nil,
	}
	return r.collection.CountDocuments(ctx, filter)
}
