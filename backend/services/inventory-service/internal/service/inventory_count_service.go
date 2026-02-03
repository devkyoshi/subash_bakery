package service

import (
	"context"
	"fmt"
	"time"

	"github.com/yourusername/erp-system/services/inventory-service/internal/repository"
	"github.com/yourusername/erp-system/shared/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type InventoryCountService struct {
	countRepo      *repository.InventoryCountRepository
	stockRepo      *repository.StockLevelRepository
	adjustmentRepo *repository.StockAdjustmentRepository
}

func NewInventoryCountService(
	countRepo *repository.InventoryCountRepository,
	stockRepo *repository.StockLevelRepository,
	adjustmentRepo *repository.StockAdjustmentRepository,
) *InventoryCountService {
	return &InventoryCountService{
		countRepo:      countRepo,
		stockRepo:      stockRepo,
		adjustmentRepo: adjustmentRepo,
	}
}

type CreateInventoryCountRequest struct {
	LocationID string    `json:"location_id" binding:"required"`
	CountNo    string    `json:"count_no"`
	CountDate  time.Time `json:"count_date"`
	CountType  string    `json:"count_type" binding:"required"`
	Notes      string    `json:"notes"`
}

type UpdateCountItemRequest struct {
	ProductID  string  `json:"product_id" binding:"required"`
	CountedQty float64 `json:"counted_qty" binding:"required,gte=0"`
	BatchID    string  `json:"batch_id"`
	Notes      string  `json:"notes"`
}

func (s *InventoryCountService) CreateCount(ctx context.Context, orgID primitive.ObjectID, req CreateInventoryCountRequest, createdBy primitive.ObjectID) (*models.InventoryCount, error) {
	locationID, err := primitive.ObjectIDFromHex(req.LocationID)
	if err != nil {
		return nil, fmt.Errorf("invalid location ID: %w", err)
	}

	count := &models.InventoryCount{
		OrganizationID: orgID,
		LocationID:     locationID,
		CountNo:        req.CountNo,
		CountDate:      req.CountDate,
		CountType:      req.CountType,
		Status:         "in_progress",
		StartedBy:      createdBy,
		StartedAt:      time.Now(),
		Notes:          req.Notes,
		Items:          []models.InventoryCountItem{},
	}

	if count.CountDate.IsZero() {
		count.CountDate = time.Now()
	}

	count.BaseModel.CreatedBy = createdBy

	if err := s.countRepo.Create(ctx, count); err != nil {
		return nil, fmt.Errorf("failed to create count: %w", err)
	}

	return count, nil
}

func (s *InventoryCountService) GetCount(ctx context.Context, id primitive.ObjectID) (*models.InventoryCount, error) {
	count, err := s.countRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("count not found: %w", err)
	}
	return count, nil
}

func (s *InventoryCountService) ListCounts(ctx context.Context, orgID primitive.ObjectID, filters map[string]interface{}, page, limit int) ([]*models.InventoryCount, error) {
	counts, err := s.countRepo.FindByOrganization(ctx, orgID, filters, page, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to list counts: %w", err)
	}
	return counts, nil
}

func (s *InventoryCountService) UpdateCountItem(ctx context.Context, countID primitive.ObjectID, req UpdateCountItemRequest, countedBy primitive.ObjectID) error {
	count, err := s.countRepo.FindByID(ctx, countID)
	if err != nil {
		return fmt.Errorf("count not found: %w", err)
	}

	if count.Status != "in_progress" {
		return fmt.Errorf("count is not in progress")
	}

	productID, err := primitive.ObjectIDFromHex(req.ProductID)
	if err != nil {
		return fmt.Errorf("invalid product ID: %w", err)
	}

	stock, err := s.stockRepo.FindByProductAndLocation(ctx, productID, count.LocationID)
	if err != nil {
		return fmt.Errorf("stock not found: %w", err)
	}

	item := models.InventoryCountItem{
		ProductID:     productID,
		SystemQty:     stock.QuantityOnHand,
		CountedQty:    req.CountedQty,
		VarianceQty:   req.CountedQty - stock.QuantityOnHand,
		UOM:           "unit",
		UnitCost:      stock.AverageCost,
		VarianceValue: (req.CountedQty - stock.QuantityOnHand) * stock.AverageCost,
		CountedBy:     countedBy,
		CountedAt:     time.Now(),
		Notes:         req.Notes,
	}

	if req.BatchID != "" {
		batchID, _ := primitive.ObjectIDFromHex(req.BatchID)
		item.BatchID = &batchID
	}

	if err := s.countRepo.UpdateItem(ctx, countID, item); err != nil {
		return fmt.Errorf("failed to update count item: %w", err)
	}

	return nil
}

func (s *InventoryCountService) CompleteCount(ctx context.Context, countID, completedBy primitive.ObjectID, createAdjustment bool) error {
	count, err := s.countRepo.FindByID(ctx, countID)
	if err != nil {
		return fmt.Errorf("count not found: %w", err)
	}

	if count.Status != "in_progress" {
		return fmt.Errorf("count is not in progress")
	}

	totalItems := len(count.Items)
	var totalVariance, varianceValue float64
	for _, item := range count.Items {
		totalVariance += item.VarianceQty
		varianceValue += item.VarianceValue
	}

	summary := map[string]interface{}{
		"total_items_counted": totalItems,
		"total_variance":      totalVariance,
		"variance_value":      varianceValue,
	}

	if err := s.countRepo.Complete(ctx, countID, completedBy, summary); err != nil {
		return fmt.Errorf("failed to complete count: %w", err)
	}

	if createAdjustment && totalVariance != 0 {
		adjustmentItems := make([]models.StockAdjustmentItem, 0, len(count.Items))
		for _, item := range count.Items {
			if item.VarianceQty != 0 {
				adjustmentItems = append(adjustmentItems, models.StockAdjustmentItem{
					ProductID:     item.ProductID,
					ExpectedQty:   item.SystemQty,
					ActualQty:     item.CountedQty,
					DifferenceQty: item.VarianceQty,
					UOM:           item.UOM,
					UnitCost:      item.UnitCost,
					TotalCost:     item.VarianceValue,
					BatchID:       item.BatchID,
					Reason:        fmt.Sprintf("Inventory count variance - %s", count.CountNo),
				})
			}
		}

		adjustment := &models.StockAdjustment{
			OrganizationID: count.OrganizationID,
			LocationID:     count.LocationID,
			AdjustmentNo:   fmt.Sprintf("ADJ-%s", count.CountNo),
			AdjustmentDate: time.Now(),
			Reason:         "inventory_count",
			ReasonDetails:  fmt.Sprintf("Auto-generated from inventory count %s", count.CountNo),
			Items:          adjustmentItems,
			Status:         "pending",
			Notes:          fmt.Sprintf("Created from inventory count ID: %s", count.ID.Hex()),
		}

		adjustment.BaseModel.CreatedBy = completedBy

		if err := s.adjustmentRepo.Create(ctx, adjustment); err != nil {
			return fmt.Errorf("failed to create adjustment: %w", err)
		}
	}

	return nil
}

func (s *InventoryCountService) CancelCount(ctx context.Context, countID primitive.ObjectID) error {
	count, err := s.countRepo.FindByID(ctx, countID)
	if err != nil {
		return fmt.Errorf("count not found: %w", err)
	}

	if count.Status == "completed" {
		return fmt.Errorf("cannot cancel completed count")
	}

	updates := map[string]interface{}{
		"status": "cancelled",
	}

	if err := s.countRepo.Update(ctx, countID, updates); err != nil {
		return fmt.Errorf("failed to cancel count: %w", err)
	}

	return nil
}

func (s *InventoryCountService) DeleteCount(ctx context.Context, id, deletedBy primitive.ObjectID) error {
	count, err := s.countRepo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("count not found: %w", err)
	}

	if count.Status == "in_progress" {
		return fmt.Errorf("cannot delete in-progress count")
	}

	if err := s.countRepo.Delete(ctx, id, deletedBy); err != nil {
		return fmt.Errorf("failed to delete count: %w", err)
	}

	return nil
}
