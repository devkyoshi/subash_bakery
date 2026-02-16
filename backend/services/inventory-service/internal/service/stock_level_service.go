package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/yourusername/erp-system/services/inventory-service/internal/client"
	"github.com/yourusername/erp-system/services/inventory-service/internal/repository"
	"github.com/yourusername/erp-system/shared/models"
	"github.com/yourusername/erp-system/shared/rabbitmq"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type StockLevelService struct {
	stockRepo     *repository.StockLevelRepository
	rabbitClient  *rabbitmq.RabbitMQClient
	productClient *client.ProductClient
}

func NewStockLevelService(stockRepo *repository.StockLevelRepository, rabbitClient *rabbitmq.RabbitMQClient, productClient *client.ProductClient) *StockLevelService {
	return &StockLevelService{
		stockRepo:     stockRepo,
		rabbitClient:  rabbitClient,
		productClient: productClient,
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

	// Check for low stock after adjustment
	// Fetch the updated stock to get current levels
	stock, err := s.stockRepo.FindByProductAndLocation(ctx, productID, locationID)
	if err == nil && stock != nil {
		s.checkLowStock(ctx, stock)
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

	s.checkLowStock(ctx, stock)

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

// GetDashboardStats retrieves critical stock count and low stock items
func (s *StockLevelService) GetDashboardStats(ctx context.Context, orgID primitive.ObjectID, token string) (map[string]interface{}, error) {
	criticalCount, inStockCount, outOfStockCount, lowStockItems, err := s.stockRepo.GetDashboardStats(ctx, orgID)
	if err != nil {
		return nil, err
	}

	// Enrich with product details if token is provided
	if token != "" && len(lowStockItems) > 0 {
		productIDs := make([]primitive.ObjectID, len(lowStockItems))
		for i, item := range lowStockItems {
			productIDs[i] = item.ProductID
		}

		products, err := s.productClient.GetProductsBatch(ctx, productIDs, token)
		if err == nil {
			for _, item := range lowStockItems {
				if prod, ok := products[item.ProductID.Hex()]; ok {
					item.ProductName = prod.Name
					item.SKU = prod.SKU
				}
			}
		} else {
			log.Printf("Failed to enrich stock items: %v", err)
		}
	}

	return map[string]interface{}{
		"critical_stock_count": criticalCount,
		"in_stock_count":       inStockCount,
		"out_of_stock_count":   outOfStockCount,
		"low_stock_items":      lowStockItems,
	}, nil
}

// checkLowStock checks if stock is below threshold and publishes event
func (s *StockLevelService) checkLowStock(ctx context.Context, stock *models.StockLevel) {
	// TODO: dynamically fetch threshold from Product service/DB.
	// For now, hardcoded threshold of 10
	threshold := 10.0

	if stock.QuantityAvailable <= threshold {
		event := map[string]interface{}{
			"type":               "low_stock",
			"product_id":         stock.ProductID.Hex(),
			"location_id":        stock.LocationID.Hex(),
			"organization_id":    stock.OrganizationID.Hex(),
			"current_stock":      stock.QuantityAvailable,
			"threshold":          threshold,
			"warehouse_zone":     stock.WarehouseZone,
			"quantity_on_hand":   stock.QuantityOnHand,
			"quantity_allocated": stock.QuantityAllocated,
		}

		err := s.rabbitClient.Publish(ctx, "notification_events", "inventory.low_stock", event)
		if err != nil {
			log.Printf("Failed to publish low stock event: %v", err)
		}
	}
}
