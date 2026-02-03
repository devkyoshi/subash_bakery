package repository

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/yourusername/erp-system/shared/models"
)

type SupplierRepository struct {
	collection *mongo.Collection
}

func NewSupplierRepository(db *mongo.Database) *SupplierRepository {
	return &SupplierRepository{
		collection: db.Collection("suppliers"),
	}
}

func (r *SupplierRepository) Create(ctx context.Context, supplier *models.Supplier) error {
	supplier.ID = primitive.NewObjectID()
	supplier.CreatedAt = time.Now()
	supplier.UpdatedAt = time.Now()
	supplier.Version = 1

	_, err := r.collection.InsertOne(ctx, supplier)
	if err != nil {
		return fmt.Errorf("failed to create supplier: %w", err)
	}
	return nil
}

func (r *SupplierRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*models.Supplier, error) {
	var supplier models.Supplier
	filter := bson.M{
		"_id":        id,
		"deleted_at": nil,
	}

	err := r.collection.FindOne(ctx, filter).Decode(&supplier)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("supplier not found")
		}
		return nil, fmt.Errorf("failed to find supplier: %w", err)
	}
	return &supplier, nil
}

func (r *SupplierRepository) FindByCode(ctx context.Context, orgID primitive.ObjectID, code string) (*models.Supplier, error) {
	var supplier models.Supplier
	filter := bson.M{
		"organization_id": orgID,
		"supplier_code":   code,
		"deleted_at":      nil,
	}

	err := r.collection.FindOne(ctx, filter).Decode(&supplier)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("supplier not found")
		}
		return nil, fmt.Errorf("failed to find supplier: %w", err)
	}
	return &supplier, nil
}

func (r *SupplierRepository) FindAll(ctx context.Context, orgID primitive.ObjectID, filters map[string]interface{}, page, limit int) ([]*models.Supplier, int64, error) {
	filter := bson.M{
		"organization_id": orgID,
		"deleted_at":      nil,
	}

	// Apply additional filters
	if status, ok := filters["status"].(string); ok && status != "" {
		filter["status"] = status
	}
	if search, ok := filters["search"].(string); ok && search != "" {
		filter["$or"] = []bson.M{
			{"supplier_code": bson.M{"$regex": search, "$options": "i"}},
			{"company_name": bson.M{"$regex": search, "$options": "i"}},
			{"contact_person": bson.M{"$regex": search, "$options": "i"}},
			{"email": bson.M{"$regex": search, "$options": "i"}},
		}
	}

	// Count total documents
	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count suppliers: %w", err)
	}

	// Find with pagination
	opts := options.Find().
		SetSort(bson.D{{Key: "created_at", Value: -1}}).
		SetSkip(int64((page - 1) * limit)).
		SetLimit(int64(limit))

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to find suppliers: %w", err)
	}
	defer cursor.Close(ctx)

	var suppliers []*models.Supplier
	if err := cursor.All(ctx, &suppliers); err != nil {
		return nil, 0, fmt.Errorf("failed to decode suppliers: %w", err)
	}

	return suppliers, total, nil
}

func (r *SupplierRepository) Update(ctx context.Context, supplier *models.Supplier) error {
	supplier.UpdatedAt = time.Now()
	supplier.Version++

	filter := bson.M{
		"_id":     supplier.ID,
		"version": supplier.Version - 1,
	}

	update := bson.M{
		"$set": supplier,
	}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update supplier: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("supplier not found or version mismatch")
	}

	return nil
}

func (r *SupplierRepository) SoftDelete(ctx context.Context, id, deletedBy primitive.ObjectID) error {
	now := time.Now()
	filter := bson.M{
		"_id":        id,
		"deleted_at": nil,
	}

	update := bson.M{
		"$set": bson.M{
			"deleted_at": &now,
			"deleted_by": &deletedBy,
			"updated_at": now,
		},
		"$inc": bson.M{"version": 1},
	}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to delete supplier: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("supplier not found or already deleted")
	}

	return nil
}

func (r *SupplierRepository) CodeExists(ctx context.Context, orgID primitive.ObjectID, code string, excludeID *primitive.ObjectID) (bool, error) {
	filter := bson.M{
		"organization_id": orgID,
		"supplier_code":   code,
		"deleted_at":      nil,
	}

	if excludeID != nil {
		filter["_id"] = bson.M{"$ne": *excludeID}
	}

	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return false, fmt.Errorf("failed to check code existence: %w", err)
	}

	return count > 0, nil
}

func (r *SupplierRepository) UpdateAnalytics(ctx context.Context, supplierID primitive.ObjectID, totalOrders int, totalPurchaseValue, outstandingBalance float64) error {
	filter := bson.M{"_id": supplierID}
	update := bson.M{
		"$set": bson.M{
			"total_orders":         totalOrders,
			"total_purchase_value": totalPurchaseValue,
			"outstanding_balance":  outstandingBalance,
			"updated_at":           time.Now(),
		},
		"$inc": bson.M{"version": 1},
	}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update supplier analytics: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("supplier not found")
	}

	return nil
}
