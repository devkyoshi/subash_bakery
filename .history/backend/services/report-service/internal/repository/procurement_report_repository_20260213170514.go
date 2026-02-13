package repository

import (
	"context"
	"fmt"
	"time"

	sharedModels "github.com/yourusername/erp-system/shared/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/yourusername/erp-system/services/report-service/internal/models"
)

type ProcurementReportRepository struct {
	poDB    *mongo.Database
	poColl  *mongo.Collection
	grnColl *mongo.Collection
}

func NewProcurementReportRepository(db *mongo.Database) *ProcurementReportRepository {
	return &ProcurementReportRepository{
		poDB:    db,
		poColl:  db.Collection("purchase_orders"),
		grnColl: db.Collection("grns"),
	}
}

// GetPurchaseOrders retrieves paginated purchase orders with filters
func (r *ProcurementReportRepository) GetPurchaseOrders(
	ctx context.Context,
	orgID primitive.ObjectID,
	filters models.ReportFilters,
	page, limit int,
) ([]*sharedModels.PurchaseOrder, int64, error) {
	filter := r.buildPOFilter(orgID, filters)

	total, err := r.poColl.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count purchase orders: %w", err)
	}

	skip := int64((page - 1) * limit)
	opts := options.Find().
		SetSort(bson.M{"order_date": -1}).
		SetSkip(skip).
		SetLimit(int64(limit))

	cursor, err := r.poColl.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to find purchase orders: %w", err)
	}
	defer cursor.Close(ctx)

	var pos []*sharedModels.PurchaseOrder
	if err := cursor.All(ctx, &pos); err != nil {
		return nil, 0, fmt.Errorf("failed to decode purchase orders: %w", err)
	}

	return pos, total, nil
}

// GetAllPurchaseOrdersForReport retrieves all purchase orders matching filters (no pagination, for export)
func (r *ProcurementReportRepository) GetAllPurchaseOrdersForReport(
	ctx context.Context,
	orgID primitive.ObjectID,
	filters models.ReportFilters,
) ([]*sharedModels.PurchaseOrder, error) {
	filter := r.buildPOFilter(orgID, filters)

	opts := options.Find().SetSort(bson.M{"order_date": -1})

	cursor, err := r.poColl.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to find purchase orders: %w", err)
	}
	defer cursor.Close(ctx)

	var pos []*sharedModels.PurchaseOrder
	if err := cursor.All(ctx, &pos); err != nil {
		return nil, fmt.Errorf("failed to decode purchase orders: %w", err)
	}

	return pos, nil
}

// GetGRNsByPurchaseOrderIDs retrieves GRNs for the given PO IDs
func (r *ProcurementReportRepository) GetGRNsByPurchaseOrderIDs(
	ctx context.Context,
	orgID primitive.ObjectID,
	poIDs []primitive.ObjectID,
) ([]*sharedModels.GoodsReceiptNote, error) {
	if len(poIDs) == 0 {
		return nil, nil
	}

	filter := bson.M{
		"organization_id":   orgID,
		"purchase_order_id": bson.M{"$in": poIDs},
		"deleted_at":        bson.M{"$exists": false},
	}

	cursor, err := r.grnColl.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to find GRNs: %w", err)
	}
	defer cursor.Close(ctx)

	var grns []*sharedModels.GoodsReceiptNote
	if err := cursor.All(ctx, &grns); err != nil {
		return nil, fmt.Errorf("failed to decode GRNs: %w", err)
	}

	return grns, nil
}

// GetPOMetricCounts returns aggregated counts of POs by status
func (r *ProcurementReportRepository) GetPOMetricCounts(
	ctx context.Context,
	orgID primitive.ObjectID,
	filters models.ReportFilters,
) (map[string]int, error) {
	matchStage := bson.M{
		"organization_id": orgID,
		"deleted_at":      bson.M{"$exists": false},
	}

	// Apply date filters
	if filters.StartDate != nil || filters.EndDate != nil {
		dateFilter := bson.M{}
		if filters.StartDate != nil {
			dateFilter["$gte"] = *filters.StartDate
		}
		if filters.EndDate != nil {
			endOfDay := filters.EndDate.Add(24*time.Hour - time.Nanosecond)
			dateFilter["$lte"] = endOfDay
		}
		matchStage["order_date"] = dateFilter
	}

	if filters.SupplierID != "" {
		supplierOID, err := primitive.ObjectIDFromHex(filters.SupplierID)
		if err == nil {
			matchStage["supplier_id"] = supplierOID
		}
	}

	if filters.Status != "" {
		matchStage["status"] = filters.Status
	}

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: matchStage}},
		{{Key: "$group", Value: bson.M{
			"_id":   "$status",
			"count": bson.M{"$sum": 1},
		}}},
	}

	cursor, err := r.poColl.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, fmt.Errorf("failed to aggregate PO metrics: %w", err)
	}
	defer cursor.Close(ctx)

	result := make(map[string]int)
	for cursor.Next(ctx) {
		var entry struct {
			Status string `bson:"_id"`
			Count  int    `bson:"count"`
		}
		if err := cursor.Decode(&entry); err != nil {
			continue
		}
		result[entry.Status] = entry.Count
	}

	return result, nil
}

// buildPOFilter constructs the MongoDB filter for purchase orders
func (r *ProcurementReportRepository) buildPOFilter(orgID primitive.ObjectID, filters models.ReportFilters) bson.M {
	filter := bson.M{
		"organization_id": orgID,
		"deleted_at":      bson.M{"$exists": false},
	}

	if filters.StartDate != nil || filters.EndDate != nil {
		dateFilter := bson.M{}
		if filters.StartDate != nil {
			dateFilter["$gte"] = *filters.StartDate
		}
		if filters.EndDate != nil {
			endOfDay := filters.EndDate.Add(24*time.Hour - time.Nanosecond)
			dateFilter["$lte"] = endOfDay
		}
		filter["order_date"] = dateFilter
	}

	if filters.SupplierID != "" {
		supplierOID, err := primitive.ObjectIDFromHex(filters.SupplierID)
		if err == nil {
			filter["supplier_id"] = supplierOID
		}
	}

	if filters.Status != "" {
		filter["status"] = filters.Status
	}

	return filter
}
