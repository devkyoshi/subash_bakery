package service

import (
	"context"
	"fmt"
	"math"

	"github.com/yourusername/erp-system/services/report-service/internal/models"
	"github.com/yourusername/erp-system/services/report-service/internal/repository"
	sharedModels "github.com/yourusername/erp-system/shared/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type StockLevelService struct {
	inventoryRepo *repository.InventoryReportRepository
}

func NewStockLevelService(inventoryRepo *repository.InventoryReportRepository) *StockLevelService {
	return &StockLevelService{
		inventoryRepo: inventoryRepo,
	}
}

// GetStockLevelReport returns the paginated stock level comparison report
func (s *StockLevelService) GetStockLevelReport(
	ctx context.Context,
	orgID primitive.ObjectID,
	filters models.StockLevelFilters,
	page, limit int,
) (*models.StockLevelReportResponse, error) {
	// Get paginated stock levels
	stockLevels, totalItems, err := s.inventoryRepo.GetStockLevels(ctx, orgID, filters, page, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get stock levels: %w", err)
	}

	// Get all stock levels for metrics (without pagination)
	allStockLevels, err := s.inventoryRepo.GetAllStockLevelsForReport(ctx, orgID, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to get all stock levels for metrics: %w", err)
	}

	// Resolve product details
	productIDs := s.collectProductIDs(stockLevels)
	products, err := s.inventoryRepo.GetProductsByIDs(ctx, productIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to get products: %w", err)
	}

	// Resolve category names
	categoryIDs := s.collectCategoryIDs(products)
	categoryNames, _ := s.inventoryRepo.GetCategoryNames(ctx, categoryIDs)

	// Resolve location names
	locationIDs := s.collectLocationIDs(stockLevels)
	locationNames, _ := s.inventoryRepo.GetLocationNames(ctx, locationIDs)

	// Resolve unit names for products
	unitNames := s.resolveUnitNames(ctx, products)

	// Build comparison items for the current page
	items := s.buildStockLevelItems(stockLevels, products, categoryNames, locationNames, unitNames)

	// Apply stock_status post-filter if needed
	if filters.StockStatus != "" {
		items = s.filterByStatus(items, filters.StockStatus)
	}

	// Build metrics from ALL stock levels
	allItems := s.buildStockLevelItems(allStockLevels, products, categoryNames, locationNames, unitNames)
	metrics := s.calculateMetrics(allItems)
	statusDist := s.buildStatusDistribution(metrics)

	return &models.StockLevelReportResponse{
		Metrics:            metrics,
		StatusDistribution: statusDist,
		Items:              items,
		TotalItems:         totalItems,
	}, nil
}

// GetAllStockLevelsForExport returns all stock level items for Excel/PDF export
func (s *StockLevelService) GetAllStockLevelsForExport(
	ctx context.Context,
	orgID primitive.ObjectID,
	filters models.StockLevelFilters,
) ([]models.StockLevelComparisonItem, *models.StockLevelMetrics, error) {
	stockLevels, err := s.inventoryRepo.GetAllStockLevelsForReport(ctx, orgID, filters)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get stock levels: %w", err)
	}

	productIDs := s.collectProductIDs(stockLevels)
	products, err := s.inventoryRepo.GetProductsByIDs(ctx, productIDs)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get products: %w", err)
	}

	categoryIDs := s.collectCategoryIDs(products)
	categoryNames, _ := s.inventoryRepo.GetCategoryNames(ctx, categoryIDs)

	locationIDs := s.collectLocationIDs(stockLevels)
	locationNames, _ := s.inventoryRepo.GetLocationNames(ctx, locationIDs)

	unitNames := s.resolveUnitNames(ctx, products)

	items := s.buildStockLevelItems(stockLevels, products, categoryNames, locationNames, unitNames)

	if filters.StockStatus != "" {
		items = s.filterByStatus(items, filters.StockStatus)
	}

	metrics := s.calculateMetrics(items)
	return items, &metrics, nil
}

func (s *StockLevelService) collectProductIDs(levels []*sharedModels.StockLevel) []primitive.ObjectID {
	seen := make(map[primitive.ObjectID]bool)
	var ids []primitive.ObjectID
	for _, sl := range levels {
		if !seen[sl.ProductID] {
			seen[sl.ProductID] = true
			ids = append(ids, sl.ProductID)
		}
	}
	return ids
}

func (s *StockLevelService) collectCategoryIDs(products map[string]*sharedModels.Product) []primitive.ObjectID {
	seen := make(map[primitive.ObjectID]bool)
	var ids []primitive.ObjectID
	for _, p := range products {
		if !p.CategoryID.IsZero() && !seen[p.CategoryID] {
			seen[p.CategoryID] = true
			ids = append(ids, p.CategoryID)
		}
	}
	return ids
}

func (s *StockLevelService) collectLocationIDs(levels []*sharedModels.StockLevel) []primitive.ObjectID {
	seen := make(map[primitive.ObjectID]bool)
	var ids []primitive.ObjectID
	for _, sl := range levels {
		if !seen[sl.LocationID] {
			seen[sl.LocationID] = true
			ids = append(ids, sl.LocationID)
		}
	}
	return ids
}

func (s *StockLevelService) resolveUnitNames(ctx context.Context, products map[string]*sharedModels.Product) map[string]string {
	unitNames := make(map[string]string)
	for _, p := range products {
		if !p.BaseUnitID.IsZero() {
			idStr := p.BaseUnitID.Hex()
			if _, exists := unitNames[idStr]; !exists {
				unitNames[idStr] = s.inventoryRepo.GetUnitName(ctx, p.BaseUnitID)
			}
		}
	}
	return unitNames
}

func (s *StockLevelService) buildStockLevelItems(
	levels []*sharedModels.StockLevel,
	products map[string]*sharedModels.Product,
	categoryNames map[string]string,
	locationNames map[string]string,
	unitNames map[string]string,
) []models.StockLevelComparisonItem {
	var items []models.StockLevelComparisonItem

	for _, sl := range levels {
		product := products[sl.ProductID.Hex()]
		if product == nil {
			continue
		}

		categoryName := categoryNames[product.CategoryID.Hex()]
		if categoryName == "" {
			categoryName = "Uncategorized"
		}

		locationName := locationNames[sl.LocationID.Hex()]
		if locationName == "" {
			locationName = "Unknown"
		}

		unit := unitNames[product.BaseUnitID.Hex()]
		if unit == "" {
			unit = "pcs"
		}

		stockStatus := s.determineStockStatus(sl, product)

		item := models.StockLevelComparisonItem{
			ProductID:    sl.ProductID.Hex(),
			SKU:          product.SKU,
			ProductName:  product.Name,
			CategoryID:   product.CategoryID.Hex(),
			CategoryName: categoryName,
			LocationID:   sl.LocationID.Hex(),
			LocationName: locationName,
			Unit:         unit,
			SystemQty:    math.Round(sl.QuantityOnHand*100) / 100,
			AvailableQty: math.Round(sl.QuantityAvailable*100) / 100,
			AllocatedQty: math.Round(sl.QuantityAllocated*100) / 100,
			InTransitQty: math.Round(sl.QuantityInTransit*100) / 100,
			ReorderLevel: product.ReorderLevel,
			MinStock:     product.MinStockLevel,
			MaxStock:     product.MaxStockLevel,
			AverageCost:  math.Round(sl.AverageCost*100) / 100,
			TotalValue:   math.Round(sl.TotalValue*100) / 100,
			StockStatus:  stockStatus,
		}

		items = append(items, item)
	}

	return items
}

func (s *StockLevelService) determineStockStatus(sl *sharedModels.StockLevel, product *sharedModels.Product) string {
	onHand := sl.QuantityOnHand

	if onHand <= 0 {
		return "OUT_OF_STOCK"
	}

	if product.MinStockLevel > 0 && onHand <= float64(product.MinStockLevel) {
		return "CRITICAL"
	}

	if product.ReorderLevel > 0 && onHand <= float64(product.ReorderLevel) {
		return "LOW"
	}

	if product.MaxStockLevel > 0 && onHand >= float64(product.MaxStockLevel) {
		return "OVERSTOCK"
	}

	return "OPTIMAL"
}

func (s *StockLevelService) filterByStatus(items []models.StockLevelComparisonItem, status string) []models.StockLevelComparisonItem {
	var filtered []models.StockLevelComparisonItem
	for _, item := range items {
		if item.StockStatus == status {
			filtered = append(filtered, item)
		}
	}
	return filtered
}

func (s *StockLevelService) calculateMetrics(items []models.StockLevelComparisonItem) models.StockLevelMetrics {
	metrics := models.StockLevelMetrics{
		TotalProducts: len(items),
	}

	for _, item := range items {
		metrics.TotalStockValue += item.TotalValue
		metrics.TotalOnHand += item.SystemQty
		metrics.TotalAllocated += item.AllocatedQty
		metrics.TotalAvailable += item.AvailableQty

		switch item.StockStatus {
		case "OPTIMAL":
			metrics.OptimalCount++
		case "LOW":
			metrics.LowStockCount++
		case "CRITICAL":
			metrics.CriticalCount++
		case "OVERSTOCK":
			metrics.OverstockCount++
		case "OUT_OF_STOCK":
			metrics.OutOfStockCount++
		}
	}

	metrics.TotalStockValue = math.Round(metrics.TotalStockValue*100) / 100
	metrics.TotalOnHand = math.Round(metrics.TotalOnHand*100) / 100
	metrics.TotalAllocated = math.Round(metrics.TotalAllocated*100) / 100
	metrics.TotalAvailable = math.Round(metrics.TotalAvailable*100) / 100

	return metrics
}

func (s *StockLevelService) buildStatusDistribution(metrics models.StockLevelMetrics) []models.StockStatusDistribution {
	total := metrics.TotalProducts
	if total == 0 {
		total = 1
	}

	return []models.StockStatusDistribution{
		{Name: "Optimal", Value: math.Round(float64(metrics.OptimalCount) / float64(total) * 100), Color: "#22c55e"},
		{Name: "Low Stock", Value: math.Round(float64(metrics.LowStockCount) / float64(total) * 100), Color: "#eab308"},
		{Name: "Critical", Value: math.Round(float64(metrics.CriticalCount) / float64(total) * 100), Color: "#ef4444"},
		{Name: "Overstock", Value: math.Round(float64(metrics.OverstockCount) / float64(total) * 100), Color: "#3b82f6"},
		{Name: "Out of Stock", Value: math.Round(float64(metrics.OutOfStockCount) / float64(total) * 100), Color: "#6b7280"},
	}
}
