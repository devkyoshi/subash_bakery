package service

import (
	"context"
	"fmt"
	"math"

	"github.com/yourusername/erp-system/services/report-service/internal/models"
	"github.com/yourusername/erp-system/services/report-service/internal/repository"
	sharedModels "github.com/yourusername/erp-system/shared/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ReorderStatusService struct {
	inventoryRepo *repository.InventoryReportRepository
	db            *mongo.Database
}

func NewReorderStatusService(inventoryRepo *repository.InventoryReportRepository, db *mongo.Database) *ReorderStatusService {
	return &ReorderStatusService{
		inventoryRepo: inventoryRepo,
		db:            db,
	}
}

// GetReorderStatusReport returns the reorder status report
func (s *ReorderStatusService) GetReorderStatusReport(
	ctx context.Context,
	orgID primitive.ObjectID,
	filters models.ReorderStatusFilters,
	page, limit int,
) (*models.ReorderStatusReportResponse, error) {
	// Convert reorder filters to stock level filters for data retrieval
	slFilters := models.StockLevelFilters{
		CategoryID: filters.CategoryID,
		LocationID: filters.LocationID,
		Search:     filters.Search,
	}

	// Get all stock levels to compute metrics
	allStockLevels, err := s.inventoryRepo.GetAllStockLevelsForReport(ctx, orgID, slFilters)
	if err != nil {
		return nil, fmt.Errorf("failed to get stock levels: %w", err)
	}

	// Resolve product details
	productIDs := s.collectProductIDs(allStockLevels)
	products, err := s.inventoryRepo.GetProductsByIDs(ctx, productIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to get products: %w", err)
	}

	// Resolve category names
	categoryIDs := s.collectCategoryIDs(products)
	categoryNames, _ := s.inventoryRepo.GetCategoryNames(ctx, categoryIDs)

	// Resolve location names
	locationIDs := s.collectLocationIDs(allStockLevels)
	locationNames, _ := s.inventoryRepo.GetLocationNames(ctx, locationIDs)

	// Resolve unit names
	unitNames := s.resolveUnitNames(ctx, products)

	// Get pending PO quantities per product
	pendingQty := s.getPendingPOQuantities(ctx, orgID)

	// Build reorder items from all stock levels
	allItems := s.buildReorderItems(allStockLevels, products, categoryNames, locationNames, unitNames, pendingQty)

	// Apply priority filter
	if filters.Priority != "" {
		allItems = s.filterByPriority(allItems, filters.Priority)
	}

	// Compute metrics from all items (before pagination)
	metrics := s.calculateMetrics(allItems)

	// Compute consumption data by category
	consumptionData := s.buildConsumptionData(allStockLevels, products, categoryNames)

	// Get total count
	totalItems := int64(len(allItems))

	// Apply pagination
	start := (page - 1) * limit
	end := start + limit
	if start > len(allItems) {
		start = len(allItems)
	}
	if end > len(allItems) {
		end = len(allItems)
	}
	pagedItems := allItems[start:end]

	return &models.ReorderStatusReportResponse{
		Metrics:         metrics,
		Items:           pagedItems,
		ConsumptionData: consumptionData,
		TotalItems:      totalItems,
	}, nil
}

func (s *ReorderStatusService) buildReorderItems(
	levels []*sharedModels.StockLevel,
	products map[string]*sharedModels.Product,
	categoryNames map[string]string,
	locationNames map[string]string,
	unitNames map[string]string,
	pendingQty map[string]float64,
) []models.ReorderItem {
	var items []models.ReorderItem

	for _, sl := range levels {
		product := products[sl.ProductID.Hex()]
		if product == nil {
			continue
		}

		unit := unitNames[product.BaseUnitID.Hex()]
		if unit == "" {
			unit = "pcs"
		}

		priority := s.determinePriority(sl, product)
		remainingDays := s.estimateRemainingDays(sl, product)
		sugQty := s.calculateSuggestedQty(sl, product)

		// Pending orders text
		pending := "—"
		if qty, ok := pendingQty[sl.ProductID.Hex()]; ok && qty > 0 {
			pending = fmt.Sprintf("%.0f %s", qty, unit)
		}

		// Lead time text
		leadTime := "—"
		if product.LeadTimeDays > 0 {
			leadTime = fmt.Sprintf("%d Days", product.LeadTimeDays)
		}

		item := models.ReorderItem{
			ID:            product.SKU,
			Name:          product.Name,
			Unit:          unit,
			Priority:      priority,
			CurrentStock:  math.Round(sl.QuantityOnHand*100) / 100,
			MinLevel:      product.ReorderLevel,
			RemainingDays: remainingDays,
			Pending:       pending,
			SugQty:        sugQty,
			LeadTime:      leadTime,
		}

		items = append(items, item)
	}

	// Sort: CRITICAL first, then WARNING, then NORMAL
	s.sortByPriority(items)

	return items
}

func (s *ReorderStatusService) determinePriority(sl *sharedModels.StockLevel, product *sharedModels.Product) string {
	onHand := sl.QuantityOnHand

	if onHand <= 0 {
		return "CRITICAL"
	}
	if product.MinStockLevel > 0 && onHand <= float64(product.MinStockLevel) {
		return "CRITICAL"
	}
	if product.ReorderLevel > 0 && onHand <= float64(product.ReorderLevel) {
		return "WARNING"
	}
	return "NORMAL"
}

func (s *ReorderStatusService) estimateRemainingDays(sl *sharedModels.StockLevel, product *sharedModels.Product) int {
	onHand := sl.QuantityOnHand
	if onHand <= 0 {
		return 0
	}

	// Estimate daily consumption based on reorder level and lead time
	// If we have lead time and reorder level, approximate daily usage
	if product.LeadTimeDays > 0 && product.ReorderLevel > 0 {
		dailyUsage := float64(product.ReorderLevel) / float64(product.LeadTimeDays)
		if dailyUsage > 0 {
			days := int(onHand / dailyUsage)
			if days > 90 {
				days = 90
			}
			return days
		}
	}

	// Fallback: use ratio of current stock to min level
	if product.MinStockLevel > 0 {
		ratio := onHand / float64(product.MinStockLevel)
		days := int(ratio * 7) // approximate 7 days per ratio unit
		if days > 90 {
			days = 90
		}
		if days < 0 {
			days = 0
		}
		return days
	}

	return 30 // default if no info
}

func (s *ReorderStatusService) calculateSuggestedQty(sl *sharedModels.StockLevel, product *sharedModels.Product) int {
	// If reorder quantity is defined, use it
	if product.ReorderQuantity > 0 {
		return product.ReorderQuantity
	}

	// Otherwise, suggest enough to reach max stock level
	if product.MaxStockLevel > 0 {
		needed := float64(product.MaxStockLevel) - sl.QuantityOnHand
		if needed > 0 {
			return int(math.Ceil(needed))
		}
	}

	// Fallback: suggest double the reorder level
	if product.ReorderLevel > 0 {
		return product.ReorderLevel * 2
	}

	return 0
}

func (s *ReorderStatusService) calculateMetrics(items []models.ReorderItem) models.ReorderMetrics {
	metrics := models.ReorderMetrics{}
	for _, item := range items {
		switch item.Priority {
		case "CRITICAL":
			metrics.CriticalCount++
		case "WARNING":
			metrics.WarningCount++
		case "NORMAL":
			metrics.NormalCount++
		}
	}
	return metrics
}

func (s *ReorderStatusService) filterByPriority(items []models.ReorderItem, priority string) []models.ReorderItem {
	var filtered []models.ReorderItem
	for _, item := range items {
		if item.Priority == priority {
			filtered = append(filtered, item)
		}
	}
	return filtered
}

func (s *ReorderStatusService) sortByPriority(items []models.ReorderItem) {
	priorityOrder := map[string]int{
		"CRITICAL": 0,
		"WARNING":  1,
		"NORMAL":   2,
	}
	for i := 0; i < len(items); i++ {
		for j := i + 1; j < len(items); j++ {
			pi := priorityOrder[items[i].Priority]
			pj := priorityOrder[items[j].Priority]
			if pi > pj || (pi == pj && items[i].RemainingDays > items[j].RemainingDays) {
				items[i], items[j] = items[j], items[i]
			}
		}
	}
}

func (s *ReorderStatusService) buildConsumptionData(
	levels []*sharedModels.StockLevel,
	products map[string]*sharedModels.Product,
	categoryNames map[string]string,
) []models.ConsumptionRow {
	// Aggregate by category
	type categoryAgg struct {
		totalOnHand   float64
		productCount  int
		totalReorder  float64
		totalLeadTime int
		leadTimeCount int
		categoryName  string
	}

	catMap := make(map[string]*categoryAgg)

	for _, sl := range levels {
		product := products[sl.ProductID.Hex()]
		if product == nil {
			continue
		}
		catID := product.CategoryID.Hex()
		catName := categoryNames[catID]
		if catName == "" {
			catName = "Uncategorized"
		}

		agg, ok := catMap[catID]
		if !ok {
			agg = &categoryAgg{categoryName: catName}
			catMap[catID] = agg
		}
		agg.totalOnHand += sl.QuantityOnHand
		agg.productCount++
		agg.totalReorder += float64(product.ReorderLevel)
		if product.LeadTimeDays > 0 {
			agg.totalLeadTime += product.LeadTimeDays
			agg.leadTimeCount++
		}
	}

	var rows []models.ConsumptionRow
	for _, agg := range catMap {
		if agg.productCount == 0 {
			continue
		}

		// Estimate average daily consumption from reorder level and lead time
		avgLeadTime := 7.0 // default
		if agg.leadTimeCount > 0 {
			avgLeadTime = float64(agg.totalLeadTime) / float64(agg.leadTimeCount)
		}

		dailyConsumption := 0.0
		if avgLeadTime > 0 && agg.totalReorder > 0 {
			dailyConsumption = agg.totalReorder / avgLeadTime
		}

		// Monthly forecast (30 days)
		monthlyForecast := dailyConsumption * 30

		// Trend: compare current stock to reorder levels
		trendDir := "neutral"
		trendStr := "0%"
		if agg.totalReorder > 0 {
			ratio := agg.totalOnHand / agg.totalReorder
			if ratio > 1.2 {
				trendDir = "up"
				pct := (ratio - 1.0) * 100
				trendStr = fmt.Sprintf("+ %.0f%%", pct)
			} else if ratio < 0.8 {
				trendDir = "down"
				pct := (1.0 - ratio) * 100
				trendStr = fmt.Sprintf("↓ %.0f%%", pct)
			} else {
				trendStr = fmt.Sprintf("-%.1f%%", math.Abs(ratio-1.0)*100)
			}
		}

		row := models.ConsumptionRow{
			Category: agg.categoryName,
			AvgDaily: fmt.Sprintf("%.1f", dailyConsumption),
			Trend:    trendStr,
			TrendDir: trendDir,
			Forecast: fmt.Sprintf("%.0f", monthlyForecast),
		}
		rows = append(rows, row)
	}

	return rows
}

// getPendingPOQuantities returns map of product_id -> total pending quantity from open POs
func (s *ReorderStatusService) getPendingPOQuantities(ctx context.Context, orgID primitive.ObjectID) map[string]float64 {
	result := make(map[string]float64)

	poColl := s.db.Collection("purchase_orders")

	// Find POs with pending statuses
	filter := bson.M{
		"organization_id": orgID,
		"status":          bson.M{"$in": []string{"draft", "sent", "confirmed", "partial"}},
		"deleted_at":      bson.M{"$exists": false},
	}

	cursor, err := poColl.Find(ctx, filter)
	if err != nil {
		return result
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var po sharedModels.PurchaseOrder
		if err := cursor.Decode(&po); err != nil {
			continue
		}
		for _, item := range po.Items {
			remaining := item.Quantity - item.QuantityReceived
			if remaining > 0 {
				result[item.ProductID.Hex()] += remaining
			}
		}
	}

	return result
}

func (s *ReorderStatusService) collectProductIDs(levels []*sharedModels.StockLevel) []primitive.ObjectID {
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

func (s *ReorderStatusService) collectCategoryIDs(products map[string]*sharedModels.Product) []primitive.ObjectID {
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

func (s *ReorderStatusService) collectLocationIDs(levels []*sharedModels.StockLevel) []primitive.ObjectID {
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

func (s *ReorderStatusService) resolveUnitNames(ctx context.Context, products map[string]*sharedModels.Product) map[string]string {
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
