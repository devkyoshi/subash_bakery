package service

import (
	"context"
	"fmt"
	"time"

	"github.com/yourusername/erp-system/services/inventory-service/internal/repository"
	"github.com/yourusername/erp-system/shared/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type StockAdjustmentService struct {
	adjustmentRepo *repository.StockAdjustmentRepository
	stockRepo      *repository.StockLevelRepository
	movementRepo   *repository.StockMovementRepository
}

func NewStockAdjustmentService(
	adjustmentRepo *repository.StockAdjustmentRepository,
	stockRepo *repository.StockLevelRepository,
	movementRepo *repository.StockMovementRepository,
) *StockAdjustmentService {
	return &StockAdjustmentService{
		adjustmentRepo: adjustmentRepo,
		stockRepo:      stockRepo,
		movementRepo:   movementRepo,
	}
}

type CreateStockAdjustmentRequest struct {
	LocationID     string                       `json:"location_id" binding:"required"`
	AdjustmentNo   string                       `json:"adjustment_no"`
	AdjustmentDate time.Time                    `json:"adjustment_date"`
	Reason         string                       `json:"reason" binding:"required"`
	ReasonDetails  string                       `json:"reason_details"`
	Items          []StockAdjustmentItemRequest `json:"items" binding:"required,min=1"`
	Notes          string                       `json:"notes"`
}

type StockAdjustmentItemRequest struct {
	ProductID   string  `json:"product_id" binding:"required"`
	ExpectedQty float64 `json:"expected_qty" binding:"required,gte=0"`
	ActualQty   float64 `json:"actual_qty" binding:"required,gte=0"`
	UOM         string  `json:"uom"`
	UnitCost    float64 `json:"unit_cost" binding:"gte=0"`
	BatchID     string  `json:"batch_id"`
	Reason      string  `json:"reason"`
}

// CreateAdjustment creates a new stock adjustment
func (s *StockAdjustmentService) CreateAdjustment(ctx context.Context, orgID primitive.ObjectID, req CreateStockAdjustmentRequest, createdBy primitive.ObjectID) (*models.StockAdjustment, error) {
	locationID, err := primitive.ObjectIDFromHex(req.LocationID)
	if err != nil {
		return nil, fmt.Errorf("invalid location ID: %w", err)
	}

	// Process items
	items := make([]models.StockAdjustmentItem, 0, len(req.Items))
	for _, itemReq := range req.Items {
		productID, err := primitive.ObjectIDFromHex(itemReq.ProductID)
		if err != nil {
			return nil, fmt.Errorf("invalid product ID: %w", err)
		}

		item := models.StockAdjustmentItem{
			ProductID:     productID,
			ExpectedQty:   itemReq.ExpectedQty,
			ActualQty:     itemReq.ActualQty,
			DifferenceQty: itemReq.ActualQty - itemReq.ExpectedQty,
			UOM:           itemReq.UOM,
			UnitCost:      itemReq.UnitCost,
			TotalCost:     (itemReq.ActualQty - itemReq.ExpectedQty) * itemReq.UnitCost,
			Reason:        itemReq.Reason,
		}

		if itemReq.BatchID != "" {
			batchID, _ := primitive.ObjectIDFromHex(itemReq.BatchID)
			item.BatchID = &batchID
		}

		items = append(items, item)
	}

	adjustment := &models.StockAdjustment{
		OrganizationID: orgID,
		LocationID:     locationID,
		AdjustmentNo:   req.AdjustmentNo,
		AdjustmentDate: req.AdjustmentDate,
		Reason:         req.Reason,
		ReasonDetails:  req.ReasonDetails,
		Items:          items,
		Status:         "draft",
		Notes:          req.Notes,
	}

	if adjustment.AdjustmentDate.IsZero() {
		adjustment.AdjustmentDate = time.Now()
	}

	adjustment.BaseModel.CreatedBy = createdBy

	if err := s.adjustmentRepo.Create(ctx, adjustment); err != nil {
		return nil, fmt.Errorf("failed to create adjustment: %w", err)
	}

	return adjustment, nil
}

// GetAdjustment retrieves an adjustment by ID
func (s *StockAdjustmentService) GetAdjustment(ctx context.Context, id primitive.ObjectID) (*models.StockAdjustment, error) {
	adjustment, err := s.adjustmentRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("adjustment not found: %w", err)
	}
	return adjustment, nil
}

// ListAdjustments retrieves adjustments for an organization
func (s *StockAdjustmentService) ListAdjustments(ctx context.Context, orgID primitive.ObjectID, filters map[string]interface{}, page, limit int) ([]*models.StockAdjustment, error) {
	adjustments, err := s.adjustmentRepo.FindByOrganization(ctx, orgID, filters, page, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to list adjustments: %w", err)
	}
	return adjustments, nil
}

// UpdateAdjustment updates a pending adjustment
func (s *StockAdjustmentService) UpdateAdjustment(ctx context.Context, id primitive.ObjectID, req CreateStockAdjustmentRequest, updatedBy primitive.ObjectID) (*models.StockAdjustment, error) {
	adjustment, err := s.adjustmentRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("adjustment not found: %w", err)
	}

	if adjustment.Status != "draft" {
		return nil, fmt.Errorf("can only update draft adjustments")
	}

	// Build updated items
	items := make([]models.StockAdjustmentItem, 0, len(req.Items))
	for _, itemReq := range req.Items {
		productID, err := primitive.ObjectIDFromHex(itemReq.ProductID)
		if err != nil {
			return nil, fmt.Errorf("invalid product ID: %w", err)
		}

		item := models.StockAdjustmentItem{
			ProductID:     productID,
			ExpectedQty:   itemReq.ExpectedQty,
			ActualQty:     itemReq.ActualQty,
			DifferenceQty: itemReq.ActualQty - itemReq.ExpectedQty,
			UOM:           itemReq.UOM,
			UnitCost:      itemReq.UnitCost,
			TotalCost:     (itemReq.ActualQty - itemReq.ExpectedQty) * itemReq.UnitCost,
			Reason:        itemReq.Reason,
		}

		if itemReq.BatchID != "" {
			batchID, _ := primitive.ObjectIDFromHex(itemReq.BatchID)
			item.BatchID = &batchID
		}

		items = append(items, item)
	}

	adjustment.Items = items
	adjustment.Reason = req.Reason
	adjustment.ReasonDetails = req.ReasonDetails
	adjustment.Notes = req.Notes
	adjustment.UpdatedBy = updatedBy
	adjustment.UpdatedAt = time.Now()

	updates := map[string]interface{}{
		"items":          items,
		"reason":         req.Reason,
		"reason_details": req.ReasonDetails,
		"notes":          req.Notes,
		"updated_by":     updatedBy,
		"updated_at":     time.Now(),
	}

	if err := s.adjustmentRepo.Update(ctx, id, updates); err != nil {
		return nil, fmt.Errorf("failed to update adjustment: %w", err)
	}

	return s.adjustmentRepo.FindByID(ctx, id)
}

// ApproveAdjustment approves an adjustment and applies stock changes
func (s *StockAdjustmentService) ApproveAdjustment(ctx context.Context, id, approvedBy primitive.ObjectID) error {
	adjustment, err := s.adjustmentRepo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("adjustment not found: %w", err)
	}

	if adjustment.Status != "draft" && adjustment.Status != "pending" {
		return fmt.Errorf("adjustment cannot be approved in current status: %s", adjustment.Status)
	}

	// Apply stock adjustments
	for _, item := range adjustment.Items {
		if item.DifferenceQty != 0 {
			// Adjust stock level
			if err := s.stockRepo.AdjustQuantity(ctx, adjustment.OrganizationID, item.ProductID, adjustment.LocationID, item.DifferenceQty, item.UnitCost); err != nil {
				return fmt.Errorf("failed to adjust stock for product %s: %w", item.ProductID.Hex(), err)
			}

			// Create stock movement record
			movement := &models.StockMovement{
				OrganizationID: adjustment.OrganizationID,
				ProductID:      item.ProductID,
				MovementType:   models.MovementAdjustment,
				Quantity:       item.DifferenceQty,
				UOM:            item.UOM,
				UnitCost:       item.UnitCost,
				TotalCost:      item.TotalCost,
				ReferenceType:  "stock_adjustment",
				ReferenceID:    adjustment.ID,
				ReferenceNo:    adjustment.AdjustmentNo,
				Reason:         adjustment.Reason,
				Notes:          item.Reason,
				MovementDate:   time.Now(),
			}

			if item.DifferenceQty > 0 {
				movement.ToLocationID = &adjustment.LocationID
			} else {
				movement.FromLocationID = &adjustment.LocationID
			}

			movement.BaseModel.CreatedBy = approvedBy
			if err := s.movementRepo.Create(ctx, movement); err != nil {
				return fmt.Errorf("failed to create movement record: %w", err)
			}
		}
	}

	// Update adjustment status
	updates := map[string]interface{}{
		"status":      "approved",
		"approved_by": approvedBy,
		"approved_at": time.Now(),
		"updated_by":  approvedBy,
	}

	if err := s.adjustmentRepo.Update(ctx, id, updates); err != nil {
		return fmt.Errorf("failed to approve adjustment: %w", err)
	}

	return nil
}

// RejectAdjustment rejects an adjustment
func (s *StockAdjustmentService) RejectAdjustment(ctx context.Context, id, rejectedBy primitive.ObjectID, reason string) error {
	updates := map[string]interface{}{
		"status":          "rejected",
		"rejected_reason": reason,
		"approved_by":     rejectedBy,
		"approved_at":     time.Now(),
	}

	if err := s.adjustmentRepo.Update(ctx, id, updates); err != nil {
		return fmt.Errorf("failed to reject adjustment: %w", err)
	}

	return nil
}

// DeleteAdjustment soft deletes an adjustment
func (s *StockAdjustmentService) DeleteAdjustment(ctx context.Context, id, deletedBy primitive.ObjectID) error {
	adjustment, err := s.adjustmentRepo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("adjustment not found: %w", err)
	}

	if adjustment.Status == "approved" {
		return fmt.Errorf("cannot delete approved adjustment")
	}

	if err := s.adjustmentRepo.Delete(ctx, id, deletedBy); err != nil {
		return fmt.Errorf("failed to delete adjustment: %w", err)
	}

	return nil
}
