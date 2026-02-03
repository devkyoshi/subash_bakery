package service

import (
	"context"
	"fmt"
	"time"

	"github.com/yourusername/erp-system/services/inventory-service/internal/repository"
	"github.com/yourusername/erp-system/shared/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type StockLevelService struct {
	stockRepo *repository.StockLevelRepository
}

func NewStockLevelService(stockRepo *repository.StockLevelRepository) *StockLevelService {
	return &StockLevelService{
		stockRepo: stockRepo,
	}
}

// GetStockLevel retrieves stock level for a product at a location
func (s *StockLevelService) GetStockLevel(ctx context.Context, productID, locationID primitive.ObjectID) (*models.StockLevel, error) {
	stock, err := s.stockRepo.FindByProductAndLocation(ctx, productID, locationID)
	if err != nil {
		return nil, fmt.Errorf("stock not found: %w", err)
	}
	return stock, nil
}

// GetStockByProduct retrieves all stock levels for a product across all locations
func (s *StockLevelService) GetStockByProduct(ctx context.Context, productID primitive.ObjectID) ([]*models.StockLevel, error) {
	stocks, err := s.stockRepo.FindByProduct(ctx, productID)
	if err != nil {
		return nil, fmt.Errorf("failed to get stock: %w", err)
	}
	return stocks, nil
}

// GetStockByProductIDs retrieves stock levels for multiple products
func (s *StockLevelService) GetStockByProductIDs(ctx context.Context, productIDs []primitive.ObjectID) (map[string]*models.StockLevel, error) {
	return s.stockRepo.FindByProducts(ctx, productIDs)
}

// GetStockByLocation retrieves all stock at a specific location
func (s *StockLevelService) GetStockByLocation(ctx context.Context, locationID primitive.ObjectID) ([]*models.StockLevel, error) {
	stocks, err := s.stockRepo.FindByLocation(ctx, locationID)
	if err != nil {
		return nil, fmt.Errorf("failed to get stock: %w", err)
	}
	return stocks, nil
}

// ListStockLevels retrieves stock levels with filters and pagination
func (s *StockLevelService) ListStockLevels(ctx context.Context, filters map[string]interface{}, page, limit int) ([]models.StockLevel, error) {
	return s.stockRepo.Find(ctx, filters, page, limit)
}

// UpsertStockLevel creates or updates stock level
func (s *StockLevelService) UpsertStockLevel(ctx context.Context, stock *models.StockLevel) error {
	if err := s.stockRepo.Upsert(ctx, stock); err != nil {
		return fmt.Errorf("failed to upsert stock level: %w", err)
	}
	return nil
}

// AdjustQuantity adjusts stock quantity at a location
func (s *StockLevelService) AdjustQuantity(ctx context.Context, orgID, productID, locationID primitive.ObjectID, delta float64, cost float64) error {
	if err := s.stockRepo.AdjustQuantity(ctx, orgID, productID, locationID, delta, cost); err != nil {
		return fmt.Errorf("failed to adjust quantity: %w", err)
	}
	return nil
}

// AllocateStock reserves stock for an order
func (s *StockLevelService) AllocateStock(ctx context.Context, productID, locationID primitive.ObjectID, quantity float64) error {
	stock, err := s.stockRepo.FindByProductAndLocation(ctx, productID, locationID)
	if err != nil {
		return fmt.Errorf("failed to find stock: %w", err)
	}

	// If stock level doesn't exist, create it with 0 initial stock
	if stock == nil {
		now := time.Now()
		stock = &models.StockLevel{
			ProductID:         productID,
			LocationID:        locationID,
			QuantityOnHand:    0,
			QuantityAvailable: 0,
			QuantityAllocated: 0,
			QuantityInTransit: 0,
			QuantityReserved:  0,
			AverageCost:       0,
			LastCost:          0,
			TotalValue:        0,
		}
		stock.ID = primitive.NewObjectID()
		stock.CreatedAt = now
		stock.UpdatedAt = now
	}

	if stock.QuantityAvailable < quantity {
		return fmt.Errorf("insufficient stock: available %.2f, requested %.2f", stock.QuantityAvailable, quantity)
	}

	stock.QuantityAllocated += quantity
	stock.QuantityAvailable -= quantity
	stock.UpdatedAt = time.Now()

	if err := s.stockRepo.Upsert(ctx, stock); err != nil {
		return fmt.Errorf("failed to allocate stock: %w", err)
	}

	return nil
}

// ReleaseStock releases allocated stock
func (s *StockLevelService) ReleaseStock(ctx context.Context, productID, locationID primitive.ObjectID, quantity float64) error {
	stock, err := s.stockRepo.FindByProductAndLocation(ctx, productID, locationID)
	if err != nil {
		return fmt.Errorf("stock not found: %w", err)
	}

	stock.QuantityAllocated -= quantity
	stock.QuantityAvailable += quantity
	stock.UpdatedAt = time.Now()

	if err := s.stockRepo.Upsert(ctx, stock); err != nil {
		return fmt.Errorf("failed to release stock: %w", err)
	}

	return nil
}
