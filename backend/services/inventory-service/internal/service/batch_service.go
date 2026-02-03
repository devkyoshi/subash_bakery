package service

import (
	"context"
	"fmt"
	"time"

	"github.com/yourusername/erp-system/services/inventory-service/internal/repository"
	"github.com/yourusername/erp-system/shared/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BatchService struct {
	batchRepo *repository.BatchRepository
	stockRepo *repository.StockLevelRepository
}

func NewBatchService(batchRepo *repository.BatchRepository, stockRepo *repository.StockLevelRepository) *BatchService {
	return &BatchService{
		batchRepo: batchRepo,
		stockRepo: stockRepo,
	}
}

type CreateBatchRequest struct {
	ProductID       string     `json:"product_id" binding:"required"`
	LocationID      string     `json:"location_id" binding:"required"`
	BatchNumber     string     `json:"batch_number" binding:"required"`
	InitialQuantity float64    `json:"initial_quantity" binding:"required,gt=0"`
	UnitCost        float64    `json:"unit_cost" binding:"required,gte=0"`
	ManufactureDate time.Time  `json:"manufacture_date"`
	ExpiryDate      *time.Time `json:"expiry_date"`
	SupplierID      string     `json:"supplier_id"`
	PurchaseOrderID string     `json:"purchase_order_id"`
	QCStatus        string     `json:"qc_status"`
	QCNotes         string     `json:"qc_notes"`
}

func (s *BatchService) CreateBatch(ctx context.Context, orgID primitive.ObjectID, req CreateBatchRequest, createdBy primitive.ObjectID) (*models.Batch, error) {
	productID, err := primitive.ObjectIDFromHex(req.ProductID)
	if err != nil {
		return nil, fmt.Errorf("invalid product ID: %w", err)
	}

	locationID, err := primitive.ObjectIDFromHex(req.LocationID)
	if err != nil {
		return nil, fmt.Errorf("invalid location ID: %w", err)
	}

	batch := &models.Batch{
		OrganizationID:  orgID,
		ProductID:       productID,
		LocationID:      locationID,
		BatchNumber:     req.BatchNumber,
		InitialQuantity: req.InitialQuantity,
		CurrentQuantity: req.InitialQuantity,
		UnitCost:        req.UnitCost,
		TotalCost:       req.InitialQuantity * req.UnitCost,
		ManufactureDate: req.ManufactureDate,
		ExpiryDate:      req.ExpiryDate,
		ReceiveDate:     time.Now(),
		QCStatus:        req.QCStatus,
		QCNotes:         req.QCNotes,
		IsActive:        true,
	}

	if req.SupplierID != "" {
		supplierID, _ := primitive.ObjectIDFromHex(req.SupplierID)
		batch.SupplierID = supplierID
	}

	if req.PurchaseOrderID != "" {
		poID, _ := primitive.ObjectIDFromHex(req.PurchaseOrderID)
		batch.PurchaseOrderID = poID
	}

	batch.BaseModel.CreatedBy = createdBy

	if err := s.batchRepo.Create(ctx, batch); err != nil {
		return nil, fmt.Errorf("failed to create batch: %w", err)
	}

	if err := s.stockRepo.AdjustQuantity(ctx, orgID, productID, locationID, req.InitialQuantity, req.UnitCost); err != nil {
		return nil, fmt.Errorf("failed to update stock level: %w", err)
	}

	return batch, nil
}

func (s *BatchService) GetBatchesByProduct(ctx context.Context, productID, locationID primitive.ObjectID, activeOnly bool) ([]*models.Batch, error) {
	batches, err := s.batchRepo.FindByProduct(ctx, productID, locationID, activeOnly)
	if err != nil {
		return nil, fmt.Errorf("failed to get batches: %w", err)
	}
	return batches, nil
}

func (s *BatchService) UpdateBatchQuantity(ctx context.Context, batchID primitive.ObjectID, delta float64) error {
	if err := s.batchRepo.UpdateQuantity(ctx, batchID, delta); err != nil {
		return fmt.Errorf("failed to update batch quantity: %w", err)
	}
	return nil
}

func (s *BatchService) GetBatch(ctx context.Context, id primitive.ObjectID) (*models.Batch, error) {
	batch, err := s.batchRepo.FindBatchByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("batch not found: %w", err)
	}
	return batch, nil
}
