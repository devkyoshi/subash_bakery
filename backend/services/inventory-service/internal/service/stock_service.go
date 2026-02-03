package service

import (
	"context"
	"fmt"
	"time"

	"github.com/yourusername/erp-system/services/inventory-service/internal/repository"
	"github.com/yourusername/erp-system/shared/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type StockService struct {
	stockRepo    *repository.StockLevelRepository
	movementRepo *repository.StockMovementRepository
	batchRepo    *repository.BatchRepository
}

func NewStockService(
	stockRepo *repository.StockLevelRepository,
	movementRepo *repository.StockMovementRepository,
	batchRepo *repository.BatchRepository,
) *StockService {
	return &StockService{
		stockRepo:    stockRepo,
		movementRepo: movementRepo,
		batchRepo:    batchRepo,
	}
}

type StockMovementRequest struct {
	ProductID      string              `json:"product_id" binding:"required"`
	MovementType   models.MovementType `json:"movement_type" binding:"required"`
	FromLocationID string              `json:"from_location_id"`
	ToLocationID   string              `json:"to_location_id"`
	Quantity       float64             `json:"quantity" binding:"required,gt=0"`
	UnitCost       float64             `json:"unit_cost"`
	ReferenceType  string              `json:"reference_type"`
	ReferenceNo    string              `json:"reference_no"`
	Reason         string              `json:"reason"`
	Notes          string              `json:"notes"`
	BatchNumber    string              `json:"batch_number"`
}

func (s *StockService) CreateStockMovement(ctx context.Context, orgID primitive.ObjectID, req StockMovementRequest, createdBy primitive.ObjectID) (*models.StockMovement, error) {
	// Parse IDs
	productID, err := primitive.ObjectIDFromHex(req.ProductID)
	if err != nil {
		return nil, fmt.Errorf("invalid product ID: %w", err)
	}

	uom := "unit" // Default UOM
	if req.UnitCost == 0 {
		uom = ""
	}

	movement := &models.StockMovement{
		OrganizationID: orgID,
		ProductID:      productID,
		MovementType:   req.MovementType,
		Quantity:       req.Quantity,
		UOM:            uom,
		UnitCost:       req.UnitCost,
		TotalCost:      req.Quantity * req.UnitCost,
		ReferenceType:  req.ReferenceType,
		ReferenceNo:    req.ReferenceNo,
		Reason:         req.Reason,
		Notes:          req.Notes,
		MovementDate:   time.Now(),
	}

	movement.BaseModel.CreatedBy = createdBy

	// Parse location IDs
	if req.FromLocationID != "" {
		fromLocID, _ := primitive.ObjectIDFromHex(req.FromLocationID)
		movement.FromLocationID = &fromLocID
	}
	if req.ToLocationID != "" {
		toLocID, _ := primitive.ObjectIDFromHex(req.ToLocationID)
		movement.ToLocationID = &toLocID
	}

	// Validate movement type has required locations
	switch req.MovementType {
	case models.MovementIn:
		if movement.ToLocationID == nil {
			return nil, fmt.Errorf("to_location_id required for IN movement")
		}
	case models.MovementOut:
		if movement.FromLocationID == nil {
			return nil, fmt.Errorf("from_location_id required for OUT movement")
		}
	case models.MovementTransfer:
		if movement.FromLocationID == nil || movement.ToLocationID == nil {
			return nil, fmt.Errorf("both from_location_id and to_location_id required for TRANSFER")
		}
	}

	// Create movement record
	if err := s.movementRepo.Create(ctx, movement); err != nil {
		return nil, fmt.Errorf("failed to create movement: %w", err)
	}

	// Update stock levels based on movement type
	switch req.MovementType {
	case models.MovementIn:
		if err := s.stockRepo.AdjustQuantity(ctx, orgID, productID, *movement.ToLocationID, req.Quantity, req.UnitCost); err != nil {
			return nil, fmt.Errorf("failed to adjust stock: %w", err)
		}

		// Create batch if batch number is provided
		if req.BatchNumber != "" {
			batch := &models.Batch{
				OrganizationID:  orgID,
				ProductID:       productID,
				LocationID:      *movement.ToLocationID,
				BatchNumber:     req.BatchNumber,
				InitialQuantity: req.Quantity,
				CurrentQuantity: req.Quantity,
				UnitCost:        req.UnitCost,
				TotalCost:       req.Quantity * req.UnitCost,
				ReceiveDate:     time.Now(),
			}
			batch.BaseModel.CreatedBy = createdBy
			s.batchRepo.Create(ctx, batch)
		}

	case models.MovementOut:
		if err := s.stockRepo.AdjustQuantity(ctx, orgID, productID, *movement.FromLocationID, -req.Quantity, 0); err != nil {
			return nil, fmt.Errorf("failed to adjust stock: %w", err)
		}

	case models.MovementTransfer:
		// Decrease from source
		if err := s.stockRepo.AdjustQuantity(ctx, orgID, productID, *movement.FromLocationID, -req.Quantity, 0); err != nil {
			return nil, fmt.Errorf("failed to adjust source stock: %w", err)
		}
		// Increase at destination
		if err := s.stockRepo.AdjustQuantity(ctx, orgID, productID, *movement.ToLocationID, req.Quantity, req.UnitCost); err != nil {
			return nil, fmt.Errorf("failed to adjust destination stock: %w", err)
		}

	case models.MovementAdjustment:
		locationID := movement.FromLocationID
		if locationID == nil {
			locationID = movement.ToLocationID
		}
		if err := s.stockRepo.AdjustQuantity(ctx, orgID, productID, *locationID, req.Quantity, req.UnitCost); err != nil {
			return nil, fmt.Errorf("failed to adjust stock: %w", err)
		}
	}

	return movement, nil
}

func (s *StockService) GetStockLevels(ctx context.Context, productID, locationID primitive.ObjectID) (*models.StockLevel, error) {
	stock, err := s.stockRepo.FindByProductAndLocation(ctx, productID, locationID)
	if err != nil {
		return nil, fmt.Errorf("stock not found: %w", err)
	}
	return stock, nil
}

func (s *StockService) GetStockByLocation(ctx context.Context, locationID primitive.ObjectID) ([]*models.StockLevel, error) {
	stocks, err := s.stockRepo.FindByLocation(ctx, locationID)
	if err != nil {
		return nil, fmt.Errorf("failed to get stock: %w", err)
	}
	return stocks, nil
}

func (s *StockService) GetStockMovements(ctx context.Context, productID primitive.ObjectID, page, limit int) ([]*models.StockMovement, error) {
	movements, err := s.movementRepo.FindByProduct(ctx, productID, page, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get movements: %w", err)
	}
	return movements, nil
}

func (s *StockService) GetStockMovement(ctx context.Context, id primitive.ObjectID) (*models.StockMovement, error) {
	return s.movementRepo.FindByID(ctx, id)
}

func (s *StockService) ListStockMovements(ctx context.Context, organizationID primitive.ObjectID, filters map[string]interface{}, page, limit int) ([]models.StockMovement, error) {
	filters["organization_id"] = organizationID
	// Only get non-deleted records (deleted_at is nil or doesn't exist)
	filters["deleted_at"] = bson.M{"$exists": false}
	return s.movementRepo.Find(ctx, filters, page, limit)
}

func (s *StockService) GetStockMovementsByLocation(ctx context.Context, locationID primitive.ObjectID) ([]models.StockMovement, error) {
	filters := map[string]interface{}{
		"$or": []bson.M{
			bson.M{"from_location_id": locationID},
			bson.M{"to_location_id": locationID},
		},
		"deleted_at": bson.M{"$exists": false},
	}
	return s.movementRepo.Find(ctx, filters, 1, 100)
}

func (s *StockService) GetBatches(ctx context.Context, productID, locationID primitive.ObjectID, activeOnly bool) ([]*models.Batch, error) {
	batches, err := s.batchRepo.FindByProduct(ctx, productID, locationID, activeOnly)
	if err != nil {
		return nil, fmt.Errorf("failed to get batches: %w", err)
	}
	return batches, nil
}
