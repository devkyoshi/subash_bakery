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

type GRNRepository struct {
	collection *mongo.Collection
}

func NewGRNRepository(db *mongo.Database) *GRNRepository {
	return &GRNRepository{
		collection: db.Collection("grns"),
	}
}

func (r *GRNRepository) Create(ctx context.Context, grn *models.GoodsReceiptNote) error {
	grn.ID = primitive.NewObjectID()
	grn.CreatedAt = time.Now()
	grn.UpdatedAt = time.Now()
	grn.Version = 1

	_, err := r.collection.InsertOne(ctx, grn)
	if err != nil {
		return fmt.Errorf("failed to create GRN: %w", err)
	}
	return nil
}

func (r *GRNRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*models.GoodsReceiptNote, error) {
	var grn models.GoodsReceiptNote
	filter := bson.M{
		"_id":        id,
		"deleted_at": nil,
	}

	err := r.collection.FindOne(ctx, filter).Decode(&grn)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("GRN not found")
		}
		return nil, fmt.Errorf("failed to find GRN: %w", err)
	}
	return &grn, nil
}

func (r *GRNRepository) FindByGRNNumber(ctx context.Context, orgID primitive.ObjectID, grnNumber string) (*models.GoodsReceiptNote, error) {
	var grn models.GoodsReceiptNote
	filter := bson.M{
		"organization_id": orgID,
		"grn_number":      grnNumber,
		"deleted_at":      nil,
	}

	err := r.collection.FindOne(ctx, filter).Decode(&grn)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("GRN not found")
		}
		return nil, fmt.Errorf("failed to find GRN: %w", err)
	}
	return &grn, nil
}

func (r *GRNRepository) FindAll(ctx context.Context, orgID primitive.ObjectID, filters map[string]interface{}, page, limit int) ([]*models.GoodsReceiptNote, int64, error) {
	filter := bson.M{
		"organization_id": orgID,
		"deleted_at":      nil,
	}

	// Apply additional filters
	if status, ok := filters["status"].(string); ok && status != "" {
		filter["status"] = status
	}
	if poID, ok := filters["purchase_order_id"].(primitive.ObjectID); ok {
		filter["purchase_order_id"] = poID
	}
	if supplierID, ok := filters["supplier_id"].(primitive.ObjectID); ok {
		filter["supplier_id"] = supplierID
	}
	if search, ok := filters["search"].(string); ok && search != "" {
		filter["$or"] = []bson.M{
			{"grn_number": bson.M{"$regex": search, "$options": "i"}},
		}
	}

	// Count total documents
	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count GRNs: %w", err)
	}

	// Find with pagination
	opts := options.Find().
		SetSort(bson.D{{Key: "receipt_date", Value: -1}}).
		SetSkip(int64((page - 1) * limit)).
		SetLimit(int64(limit))

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to find GRNs: %w", err)
	}
	defer cursor.Close(ctx)

	var grns []*models.GoodsReceiptNote
	if err := cursor.All(ctx, &grns); err != nil {
		return nil, 0, fmt.Errorf("failed to decode GRNs: %w", err)
	}

	return grns, total, nil
}

func (r *GRNRepository) Update(ctx context.Context, grn *models.GoodsReceiptNote) error {
	grn.UpdatedAt = time.Now()
	grn.Version++

	filter := bson.M{
		"_id":     grn.ID,
		"version": grn.Version - 1,
	}

	update := bson.M{
		"$set": grn,
	}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update GRN: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("GRN not found or version mismatch")
	}

	return nil
}

func (r *GRNRepository) SoftDelete(ctx context.Context, id, deletedBy primitive.ObjectID) error {
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
		return fmt.Errorf("failed to delete GRN: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("GRN not found or already deleted")
	}

	return nil
}

func (r *GRNRepository) GRNNumberExists(ctx context.Context, orgID primitive.ObjectID, grnNumber string, excludeID *primitive.ObjectID) (bool, error) {
	filter := bson.M{
		"organization_id": orgID,
		"grn_number":      grnNumber,
		"deleted_at":      nil,
	}

	if excludeID != nil {
		filter["_id"] = bson.M{"$ne": *excludeID}
	}

	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return false, fmt.Errorf("failed to check GRN number existence: %w", err)
	}

	return count > 0, nil
}

func (r *GRNRepository) FindByPurchaseOrder(ctx context.Context, poID primitive.ObjectID) ([]*models.GoodsReceiptNote, error) {
	filter := bson.M{
		"purchase_order_id": poID,
		"deleted_at":        nil,
	}

	opts := options.Find().SetSort(bson.D{{Key: "receipt_date", Value: -1}})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to find GRNs by purchase order: %w", err)
	}
	defer cursor.Close(ctx)

	var grns []*models.GoodsReceiptNote
	if err := cursor.All(ctx, &grns); err != nil {
		return nil, fmt.Errorf("failed to decode GRNs: %w", err)
	}

	return grns, nil
}

func (r *GRNRepository) CompleteInspection(ctx context.Context, id, inspectedBy primitive.ObjectID, qcStatus, qcNotes string) error {
	now := time.Now()
	filter := bson.M{
		"_id":          id,
		"deleted_at":   nil,
		"inspected_by": nil,
	}

	update := bson.M{
		"$set": bson.M{
			"status":       models.GRNStatusInspected,
			"inspected_by": &inspectedBy,
			"qc_status":    qcStatus,
			"qc_notes":     qcNotes,
			"updated_at":   now,
		},
		"$inc": bson.M{"version": 1},
	}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to complete inspection: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("GRN not found or already inspected")
	}

	return nil
}

// GetPendingCount retrieves count of pending GRNs
func (r *GRNRepository) GetPendingCount(ctx context.Context, orgID primitive.ObjectID) (int64, error) {
	// Pending GRNs (Draft, Received, Inspected) - Not Accepted or Rejected
	count, err := r.collection.CountDocuments(ctx, bson.M{
		"organization_id": orgID,
		"deleted_at":      nil,
		"status": bson.M{
			"$in": []models.GRNStatus{
				models.GRNStatusDraft,
				models.GRNStatusReceived,
				models.GRNStatusInspected,
			},
		},
	})
	if err != nil {
		return 0, fmt.Errorf("failed to count pending GRNs: %w", err)
	}
	return count, nil
}
