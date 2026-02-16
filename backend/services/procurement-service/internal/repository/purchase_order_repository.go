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

type PurchaseOrderRepository struct {
	collection *mongo.Collection
}

func NewPurchaseOrderRepository(db *mongo.Database) *PurchaseOrderRepository {
	return &PurchaseOrderRepository{
		collection: db.Collection("purchase_orders"),
	}
}

func (r *PurchaseOrderRepository) Create(ctx context.Context, po *models.PurchaseOrder) error {
	po.ID = primitive.NewObjectID()
	po.CreatedAt = time.Now()
	po.UpdatedAt = time.Now()
	po.Version = 1

	_, err := r.collection.InsertOne(ctx, po)
	if err != nil {
		return fmt.Errorf("failed to create purchase order: %w", err)
	}
	return nil
}

func (r *PurchaseOrderRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*models.PurchaseOrder, error) {
	var po models.PurchaseOrder
	filter := bson.M{
		"_id":        id,
		"deleted_at": nil,
	}

	err := r.collection.FindOne(ctx, filter).Decode(&po)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("purchase order not found")
		}
		return nil, fmt.Errorf("failed to find purchase order: %w", err)
	}
	return &po, nil
}

func (r *PurchaseOrderRepository) FindByPONumber(ctx context.Context, orgID primitive.ObjectID, poNumber string) (*models.PurchaseOrder, error) {
	var po models.PurchaseOrder
	filter := bson.M{
		"organization_id": orgID,
		"po_number":       poNumber,
		"deleted_at":      nil,
	}

	err := r.collection.FindOne(ctx, filter).Decode(&po)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("purchase order not found")
		}
		return nil, fmt.Errorf("failed to find purchase order: %w", err)
	}
	return &po, nil
}

func (r *PurchaseOrderRepository) FindAll(ctx context.Context, orgID primitive.ObjectID, filters map[string]interface{}, page, limit int) ([]*models.PurchaseOrder, int64, error) {
	filter := bson.M{
		"organization_id": orgID,
		"deleted_at":      nil,
	}

	// Apply additional filters
	if status, ok := filters["status"].(string); ok && status != "" {
		filter["status"] = status
	}
	if supplierID, ok := filters["supplier_id"].(primitive.ObjectID); ok {
		filter["supplier_id"] = supplierID
	}
	if search, ok := filters["search"].(string); ok && search != "" {
		filter["$or"] = []bson.M{
			{"po_number": bson.M{"$regex": search, "$options": "i"}},
			{"reference_number": bson.M{"$regex": search, "$options": "i"}},
		}
	}

	// Count total documents
	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count purchase orders: %w", err)
	}

	// Find with pagination
	opts := options.Find().
		SetSort(bson.D{{Key: "order_date", Value: -1}}).
		SetSkip(int64((page - 1) * limit)).
		SetLimit(int64(limit))

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to find purchase orders: %w", err)
	}
	defer cursor.Close(ctx)

	var pos []*models.PurchaseOrder
	if err := cursor.All(ctx, &pos); err != nil {
		return nil, 0, fmt.Errorf("failed to decode purchase orders: %w", err)
	}

	return pos, total, nil
}

func (r *PurchaseOrderRepository) Update(ctx context.Context, po *models.PurchaseOrder) error {
	po.UpdatedAt = time.Now()
	po.Version++

	filter := bson.M{
		"_id":     po.ID,
		"version": po.Version - 1,
	}

	update := bson.M{
		"$set": po,
	}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update purchase order: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("purchase order not found or version mismatch")
	}

	return nil
}

func (r *PurchaseOrderRepository) UpdateStatus(ctx context.Context, id primitive.ObjectID, status models.POStatus) error {
	filter := bson.M{
		"_id":        id,
		"deleted_at": nil,
	}

	update := bson.M{
		"$set": bson.M{
			"status":     status,
			"updated_at": time.Now(),
		},
		"$inc": bson.M{"version": 1},
	}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update status: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("purchase order not found")
	}

	return nil
}

func (r *PurchaseOrderRepository) SoftDelete(ctx context.Context, id, deletedBy primitive.ObjectID) error {
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
		return fmt.Errorf("failed to delete purchase order: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("purchase order not found or already deleted")
	}

	return nil
}

func (r *PurchaseOrderRepository) PONumberExists(ctx context.Context, orgID primitive.ObjectID, poNumber string, excludeID *primitive.ObjectID) (bool, error) {
	filter := bson.M{
		"organization_id": orgID,
		"po_number":       poNumber,
		"deleted_at":      nil,
	}

	if excludeID != nil {
		filter["_id"] = bson.M{"$ne": *excludeID}
	}

	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return false, fmt.Errorf("failed to check PO number existence: %w", err)
	}

	return count > 0, nil
}

func (r *PurchaseOrderRepository) Approve(ctx context.Context, id, approvedBy primitive.ObjectID) error {
	now := time.Now()
	filter := bson.M{
		"_id":         id,
		"deleted_at":  nil,
		"approved_by": nil,
	}

	update := bson.M{
		"$set": bson.M{
			"approved_by":   &approvedBy,
			"approved_date": &now,
			"status":        models.POStatusSent,
			"updated_at":    now,
		},
		"$inc": bson.M{"version": 1},
	}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to approve purchase order: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("purchase order not found or already approved")
	}

	return nil
}

// GetDashboardStats retrieves dashboard statistics
func (r *PurchaseOrderRepository) GetDashboardStats(ctx context.Context, orgID primitive.ObjectID) (int64, int64, []*models.PurchaseOrder, error) {
	// Pending PO Count (Sent, Confirmed, Partial) = Active Orders
	pendingCount, err := r.collection.CountDocuments(ctx, bson.M{
		"organization_id": orgID,
		"deleted_at":      nil,
		"status": bson.M{
			"$in": []models.POStatus{
				models.POStatusSent,
				models.POStatusConfirmed,
				models.POStatusPartiallyReceived,
			},
		},
	})
	if err != nil {
		return 0, 0, nil, fmt.Errorf("failed to count pending POs: %w", err)
	}

	// Pending Approvals (Draft status) - Limit to top 5
	opts := options.Find().
		SetSort(bson.D{{Key: "created_at", Value: -1}}).
		SetLimit(5)

	cursor, err := r.collection.Find(ctx, bson.M{
		"organization_id": orgID,
		"deleted_at":      nil,
		"status":          models.POStatusDraft,
	}, opts)
	if err != nil {
		return 0, 0, nil, fmt.Errorf("failed to find pending approvals: %w", err)
	}
	defer cursor.Close(ctx)

	var pendingApprovals []*models.PurchaseOrder
	if err = cursor.All(ctx, &pendingApprovals); err != nil {
		return 0, 0, nil, fmt.Errorf("failed to decode pending approvals: %w", err)
	}

	// Initialize slice
	if pendingApprovals == nil {
		pendingApprovals = make([]*models.PurchaseOrder, 0)
	}

	// Total POs
	totalPOs, err := r.collection.CountDocuments(ctx, bson.M{
		"organization_id": orgID,
		"deleted_at":      nil,
	})
	if err != nil {
		return 0, 0, nil, fmt.Errorf("failed to count total POs: %w", err)
	}

	return pendingCount, totalPOs, pendingApprovals, nil
}
