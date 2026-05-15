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

type BrandRepository struct {
	collection        *mongo.Collection
	productCollection *mongo.Collection
}

func NewBrandRepository(db *mongo.Database) *BrandRepository {
	return &BrandRepository{
		collection:        db.Collection("brands"),
		productCollection: db.Collection("products"),
	}
}

// Create inserts a new brand
func (r *BrandRepository) Create(ctx context.Context, brand *models.Brand) error {
	brand.ID = primitive.NewObjectID()
	brand.CreatedAt = time.Now()
	brand.UpdatedAt = time.Now()

	_, err := r.collection.InsertOne(ctx, brand)
	return err
}

// FindByID retrieves a brand by ID
func (r *BrandRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*models.Brand, error) {
	var brand models.Brand
	filter := bson.M{
		"_id":        id,
		"deleted_at": nil,
	}

	err := r.collection.FindOne(ctx, filter).Decode(&brand)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &brand, nil
}

// FindByOrganization retrieves brands for an organization with pagination
func (r *BrandRepository) FindByOrganization(ctx context.Context, orgID primitive.ObjectID, page, limit int, isActive *bool) ([]*models.Brand, int64, error) {
	filter := bson.M{
		"organization_id": orgID,
		"deleted_at":      nil,
	}

	if isActive != nil {
		filter["is_active"] = *isActive
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

	var brands []*models.Brand
	if err = cursor.All(ctx, &brands); err != nil {
		return nil, 0, err
	}

	return brands, total, nil
}

// Update updates an existing brand
func (r *BrandRepository) Update(ctx context.Context, brand *models.Brand) error {
	brand.UpdatedAt = time.Now()

	filter := bson.M{
		"_id":        brand.ID,
		"deleted_at": nil,
	}

	update := bson.M{
		"$set": brand,
	}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}

// Delete soft deletes a brand
func (r *BrandRepository) Delete(ctx context.Context, id primitive.ObjectID, deletedBy primitive.ObjectID) error {
	now := time.Now()
	filter := bson.M{
		"_id":        id,
		"deleted_at": nil,
	}

	update := bson.M{
		"$set": bson.M{
			"deleted_at": now,
			"deleted_by": deletedBy,
			"updated_at": now,
		},
	}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}

// CodeExistsInOrg checks if a brand code exists in an organization
func (r *BrandRepository) CodeExistsInOrg(ctx context.Context, code string, orgID primitive.ObjectID, excludeID *primitive.ObjectID) (bool, error) {
	filter := bson.M{
		"organization_id": orgID,
		"code":            code,
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

// NameExistsInOrg checks if a brand name exists in an organization
func (r *BrandRepository) NameExistsInOrg(ctx context.Context, name string, orgID primitive.ObjectID, excludeID *primitive.ObjectID) (bool, error) {
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

// IsInUse checks if a brand is being used by any products
func (r *BrandRepository) IsInUse(ctx context.Context, brandID primitive.ObjectID) (bool, error) {
	filter := bson.M{
		"brand_id":   brandID,
		"deleted_at": nil,
	}

	count, err := r.productCollection.CountDocuments(ctx, filter)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// Search searches brands by name or code
func (r *BrandRepository) Search(ctx context.Context, orgID primitive.ObjectID, query string, isActive *bool, page, limit int) ([]*models.Brand, int64, error) {
	// Build the base conditions
	baseConditions := []bson.M{
		{"organization_id": orgID},
		{"deleted_at": nil},
		{"$or": []bson.M{
			{"name": bson.M{"$regex": query, "$options": "i"}},
			{"code": bson.M{"$regex": query, "$options": "i"}},
		}},
	}

	// Add is_active filter if provided
	if isActive != nil {
		baseConditions = append(baseConditions, bson.M{"is_active": *isActive})
	}

	filter := bson.M{"$and": baseConditions}

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

	var brands []*models.Brand
	if err = cursor.All(ctx, &brands); err != nil {
		return nil, 0, err
	}

	return brands, total, nil
}

// Count counts total brands for an organization
func (r *BrandRepository) Count(ctx context.Context, orgID primitive.ObjectID) (int64, error) {
	filter := bson.M{
		"organization_id": orgID,
		"deleted_at":      nil,
	}
	return r.collection.CountDocuments(ctx, filter)
}
