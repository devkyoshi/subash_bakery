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

type SerialNumberService struct {
	serialRepo *repository.SerialNumberRepository
	stockRepo  *repository.StockLevelRepository
}

func NewSerialNumberService(serialRepo *repository.SerialNumberRepository, stockRepo *repository.StockLevelRepository) *SerialNumberService {
	return &SerialNumberService{
		serialRepo: serialRepo,
		stockRepo:  stockRepo,
	}
}

type CreateSerialNumberRequest struct {
	ProductID       string     `json:"product_id" binding:"required"`
	LocationID      string     `json:"location_id"`
	SerialNo        string     `json:"serial_no" binding:"required"`
	ManufactureDate *time.Time `json:"manufacture_date"`
	WarrantyExpiry  *time.Time `json:"warranty_expiry"`
	BatchID         string     `json:"batch_id"`
	UnitCost        float64    `json:"unit_cost" binding:"gte=0"`
}

// CreateSerialNumber creates a new serial number
func (s *SerialNumberService) CreateSerialNumber(ctx context.Context, orgID primitive.ObjectID, req CreateSerialNumberRequest, createdBy primitive.ObjectID) (*models.SerialNumber, error) {
	productID, err := primitive.ObjectIDFromHex(req.ProductID)
	if err != nil {
		return nil, fmt.Errorf("invalid product ID: %w", err)
	}

	// Check if serial number exists
	exists, err := s.serialRepo.SerialNoExists(ctx, orgID, req.SerialNo, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to check serial number: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("serial number already exists")
	}

	serialNumber := &models.SerialNumber{
		OrganizationID:  orgID,
		ProductID:       productID,
		SerialNo:        req.SerialNo,
		ManufactureDate: req.ManufactureDate,
		WarrantyExpiry:  req.WarrantyExpiry,
		UnitCost:        req.UnitCost,
		Status:          "available",
		IsAvailable:     true,
	}

	if req.LocationID != "" {
		locationID, _ := primitive.ObjectIDFromHex(req.LocationID)
		serialNumber.LocationID = &locationID
	}

	if req.BatchID != "" {
		batchID, _ := primitive.ObjectIDFromHex(req.BatchID)
		serialNumber.BatchID = &batchID
	}

	serialNumber.BaseModel.CreatedBy = createdBy

	if err := s.serialRepo.Create(ctx, serialNumber); err != nil {
		return nil, fmt.Errorf("failed to create serial number: %w", err)
	}

	return serialNumber, nil
}

// GetSerialNumber retrieves a serial number by ID
func (s *SerialNumberService) GetSerialNumber(ctx context.Context, id primitive.ObjectID) (*models.SerialNumber, error) {
	serialNumber, err := s.serialRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("serial number not found: %w", err)
	}
	return serialNumber, nil
}

// GetSerialNumberBySerial retrieves a serial number by serial string
func (s *SerialNumberService) GetSerialNumberBySerial(ctx context.Context, orgID primitive.ObjectID, serialNo string) (*models.SerialNumber, error) {
	serialNumber, err := s.serialRepo.FindBySerialNo(ctx, orgID, serialNo)
	if err != nil {
		return nil, fmt.Errorf("serial number not found: %w", err)
	}
	return serialNumber, nil
}

// ListSerialNumbersByProduct retrieves all serial numbers for a product
func (s *SerialNumberService) ListSerialNumbersByProduct(ctx context.Context, productID primitive.ObjectID, filters map[string]interface{}) ([]*models.SerialNumber, error) {
	serialNumbers, err := s.serialRepo.FindByProduct(ctx, productID, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to list serial numbers: %w", err)
	}
	return serialNumbers, nil
}

// UpdateSerialNumber updates serial number fields
func (s *SerialNumberService) UpdateSerialNumber(ctx context.Context, id primitive.ObjectID, updates map[string]interface{}) error {
	serialNumber, err := s.serialRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	// Don't allow updates if sold
	if serialNumber.Status == "sold" {
		return fmt.Errorf("cannot update sold serial number")
	}

	updates["updated_at"] = time.Now()
	bsonUpdates := bson.M{"$set": updates}
	return s.serialRepo.Update(ctx, id, bsonUpdates)
}

// AllocateSerialNumber allocates a serial number to a customer/order
func (s *SerialNumberService) AllocateSerialNumber(ctx context.Context, serialID, customerID, salesOrderID primitive.ObjectID) error {
	if err := s.serialRepo.Allocate(ctx, serialID, customerID, salesOrderID); err != nil {
		return fmt.Errorf("failed to allocate serial number: %w", err)
	}
	return nil
}

// MarkAsSold marks a serial number as sold
func (s *SerialNumberService) MarkAsSold(ctx context.Context, serialID primitive.ObjectID) error {
	if err := s.serialRepo.MarkAsSold(ctx, serialID); err != nil {
		return fmt.Errorf("failed to mark as sold: %w", err)
	}
	return nil
}

// DeleteSerialNumber soft-deletes a serial number
func (s *SerialNumberService) DeleteSerialNumber(ctx context.Context, id, deletedBy primitive.ObjectID) error {
	serialNumber, err := s.serialRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	// Don't allow deletion if sold
	if serialNumber.Status == "sold" {
		return fmt.Errorf("cannot delete sold serial number")
	}

	return s.serialRepo.Delete(ctx, id, deletedBy)
}
