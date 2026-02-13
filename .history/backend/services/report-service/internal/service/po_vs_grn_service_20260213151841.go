package service

import (
	"context"
	"fmt"
	"log"
	"math"
	"sort"

	"github.com/yourusername/erp-system/services/report-service/internal/models"
	"github.com/yourusername/erp-system/services/report-service/internal/repository"
	sharedModels "github.com/yourusername/erp-system/shared/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type POvsGRNService struct {
	procurementRepo *repository.ProcurementReportRepository
	supplierColl    *mongo.Collection
}

func NewPOvsGRNService(procurementRepo *repository.ProcurementReportRepository, db *mongo.Database) *POvsGRNService {
	return &POvsGRNService{
		procurementRepo: procurementRepo,
		supplierColl:    db.Collection("suppliers"),
	}
}

// GetPOvsGRNComparison returns the full comparison data for the PO vs GRN report page
func (s *POvsGRNService) GetPOvsGRNComparison(
	ctx context.Context,
	orgID primitive.ObjectID,
	filters models.ReportFilters,
	page, limit int,
) (*models.POvsGRNReportResponse, error) {
	// 1. Get POs with pagination for the table
	pos, totalItems, err := s.procurementRepo.GetPurchaseOrders(ctx, orgID, filters, page, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get purchase orders: %w", err)
	}

	// 2. Get PO metric counts (across ALL matching POs, not just the page)
	statusCounts, err := s.procurementRepo.GetPOMetricCounts(ctx, orgID, filters)
	if err != nil {
		log.Printf("Warning: failed to get PO metric counts: %v", err)
		statusCounts = map[string]int{}
	}

	// 3. Collect PO IDs for the current page to fetch their GRNs
	poIDs := make([]primitive.ObjectID, 0, len(pos))
	for _, po := range pos {
		poIDs = append(poIDs, po.ID)
	}

	// 4. Get GRNs for the page POs
	grns, err := s.procurementRepo.GetGRNsByPurchaseOrderIDs(ctx, orgID, poIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to get GRNs: %w", err)
	}

	// 5. Build a map of PO ID → [GRNs]
	grnMap := make(map[string][]*sharedModels.GoodsReceiptNote)
	for _, grn := range grns {
		key := grn.PurchaseOrderID.Hex()
		grnMap[key] = append(grnMap[key], grn)
	}

	// 6. Collect unique supplier IDs and resolve names
	supplierIDs := make(map[string]bool)
	for _, po := range pos {
		supplierIDs[po.SupplierID.Hex()] = true
	}
	supplierNames := s.resolveSupplierNames(ctx, supplierIDs)

	// 7. Build comparison items (one per PO item, aggregated across GRNs)
	items := s.buildComparisonItems(pos, grnMap, supplierNames)

	// 8. Calculate metrics from status counts
	metrics := s.calculateMetrics(statusCounts, items)

	// 9. Build variance distribution
	varianceDist := s.buildVarianceDistribution(statusCounts, metrics)

	// 10. Build action items from the comparison
	actionItems := s.buildActionItems(items)

	return &models.POvsGRNReportResponse{
		Metrics:              metrics,
		VarianceDistribution: varianceDist,
		Items:                items,
		ActionItems:          actionItems,
		TotalItems:           totalItems,
	}, nil
}

// GetAllPOvsGRNForExport returns ALL comparison data (no pagination) for export
func (s *POvsGRNService) GetAllPOvsGRNForExport(
	ctx context.Context,
	orgID primitive.ObjectID,
	filters models.ReportFilters,
) ([]models.POvsGRNComparisonItem, *models.POvsGRNMetrics, error) {
	pos, err := s.procurementRepo.GetAllPurchaseOrdersForReport(ctx, orgID, filters)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get purchase orders: %w", err)
	}

	poIDs := make([]primitive.ObjectID, 0, len(pos))
	for _, po := range pos {
		poIDs = append(poIDs, po.ID)
	}

	grns, err := s.procurementRepo.GetGRNsByPurchaseOrderIDs(ctx, orgID, poIDs)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get GRNs: %w", err)
	}

	grnMap := make(map[string][]*sharedModels.GoodsReceiptNote)
	for _, grn := range grns {
		key := grn.PurchaseOrderID.Hex()
		grnMap[key] = append(grnMap[key], grn)
	}

	supplierIDs := make(map[string]bool)
	for _, po := range pos {
		supplierIDs[po.SupplierID.Hex()] = true
	}
	supplierNames := s.resolveSupplierNames(ctx, supplierIDs)

	items := s.buildComparisonItems(pos, grnMap, supplierNames)

	statusCounts, _ := s.procurementRepo.GetPOMetricCounts(ctx, orgID, filters)
	metrics := s.calculateMetrics(statusCounts, items)

	return items, &metrics, nil
}

// resolveSupplierNames fetches supplier names from the procurement DB
func (s *POvsGRNService) resolveSupplierNames(ctx context.Context, supplierIDs map[string]bool) map[string]string {
	names := make(map[string]string)

	for idStr := range supplierIDs {
		id, err := primitive.ObjectIDFromHex(idStr)
		if err != nil {
			continue
		}

		var supplier struct {
			CompanyName string `bson:"company_name"`
		}
		err = s.supplierColl.FindOne(ctx, bson.M{"_id": id}).Decode(&supplier)
		if err != nil {
			names[idStr] = "Unknown Supplier"
			continue
		}
		names[idStr] = supplier.CompanyName
	}

	return names
}

// buildComparisonItems creates line-level PO vs GRN comparison items
func (s *POvsGRNService) buildComparisonItems(
	pos []*sharedModels.PurchaseOrder,
	grnMap map[string][]*sharedModels.GoodsReceiptNote,
	supplierNames map[string]string,
) []models.POvsGRNComparisonItem {
	var items []models.POvsGRNComparisonItem

	for _, po := range pos {
		poGRNs := grnMap[po.ID.Hex()]

		for _, poItem := range po.Items {
			// Aggregate received quantities across all GRNs for this PO
			var totalReceived, totalAccepted, totalRejected float64
			for _, grn := range poGRNs {
				for _, grnItem := range grn.Items {
					if grnItem.ProductID == poItem.ProductID {
						totalReceived += grnItem.ReceivedQuantity
						totalAccepted += grnItem.AcceptedQuantity
						totalRejected += grnItem.RejectedQuantity
					}
				}
			}

			variance := totalReceived - poItem.Quantity
			variancePct := float64(0)
			if poItem.Quantity > 0 {
				variancePct = math.Round((variance/poItem.Quantity)*10000) / 100
			}

			poValue := poItem.Quantity * poItem.UnitPrice
			grnValue := totalReceived * poItem.UnitPrice
			valueVariance := grnValue - poValue

			// Determine status
			status := "PENDING"
			if totalReceived == 0 && po.Status != string(sharedModels.POStatusReceived) {
				status = "PENDING"
			} else if math.Abs(variance) < 0.01 {
				status = "MATCHED"
			} else if variance > 0 {
				status = "EXCESS"
			} else {
				status = "PARTIAL"
			}

			supplierName := supplierNames[po.SupplierID.Hex()]
			if supplierName == "" {
				supplierName = "Unknown"
			}

			item := models.POvsGRNComparisonItem{
				POID:          po.ID.Hex(),
				PONumber:      po.PONumber,
				OrderDate:     po.OrderDate.Format("2006-01-02"),
				SupplierID:    po.SupplierID.Hex(),
				SupplierName:  supplierName,
				ProductID:     poItem.ProductID.Hex(),
				SKU:           poItem.SKU,
				ProductName:   poItem.Description,
				POQty:         poItem.Quantity,
				GRNQty:        totalReceived,
				AcceptedQty:   totalAccepted,
				RejectedQty:   totalRejected,
				Variance:      math.Round(variance*100) / 100,
				VariancePct:   variancePct,
				UnitPrice:     poItem.UnitPrice,
				POValue:       math.Round(poValue*100) / 100,
				GRNValue:      math.Round(grnValue*100) / 100,
				ValueVariance: math.Round(valueVariance*100) / 100,
				Status:        status,
			}

			items = append(items, item)
		}
	}

	return items
}

// calculateMetrics computes the summary metrics from status counts and items
func (s *POvsGRNService) calculateMetrics(statusCounts map[string]int, items []models.POvsGRNComparisonItem) models.POvsGRNMetrics {
	totalPOs := 0
	for _, count := range statusCounts {
		totalPOs += count
	}

	completedPOs := statusCounts[string(sharedModels.POStatusReceived)]
	partialPOs := statusCounts[string(sharedModels.POStatusPartiallyReceived)]
	pendingPOs := statusCounts[string(sharedModels.POStatusDraft)] +
		statusCounts[string(sharedModels.POStatusSent)] +
		statusCounts[string(sharedModels.POStatusConfirmed)]

	// Count excess from items
	excessCount := 0
	for _, item := range items {
		if item.Status == "EXCESS" {
			excessCount++
		}
	}

	// Calculate totals from items
	var totalPOValue, totalGRNValue, totalVariance float64
	for _, item := range items {
		totalPOValue += item.POValue
		totalGRNValue += item.GRNValue
		totalVariance += math.Abs(item.ValueVariance)
	}

	variancePercent := float64(0)
	if totalPOValue > 0 {
		variancePercent = math.Round((totalVariance/totalPOValue)*10000) / 100
	}

	completedPercent := float64(0)
	if totalPOs > 0 {
		completedPercent = math.Round((float64(completedPOs)/float64(totalPOs))*10000) / 100
	}

	return models.POvsGRNMetrics{
		TotalPOs:         totalPOs,
		CompletedPOs:     completedPOs,
		PartialPOs:       partialPOs,
		PendingPOs:       pendingPOs,
		ExcessPOs:        excessCount,
		TotalVariance:    math.Round(totalVariance*100) / 100,
		TotalPOValue:     math.Round(totalPOValue*100) / 100,
		TotalGRNValue:    math.Round(totalGRNValue*100) / 100,
		VariancePercent:  variancePercent,
		CompletedPercent: completedPercent,
	}
}

// buildVarianceDistribution creates the pie chart data
func (s *POvsGRNService) buildVarianceDistribution(statusCounts map[string]int, metrics models.POvsGRNMetrics) []models.VarianceDistribution {
	total := metrics.TotalPOs
	if total == 0 {
		total = 1 // prevent division by zero
	}

	matched := metrics.CompletedPOs
	partial := metrics.PartialPOs
	excess := metrics.ExcessPOs
	pending := metrics.PendingPOs

	return []models.VarianceDistribution{
		{
			Name:  "Matched",
			Value: math.Round(float64(matched) / float64(total) * 100),
			Color: "#22c55e",
		},
		{
			Name:  "Shortage",
			Value: math.Round(float64(partial) / float64(total) * 100),
			Color: "#eab308",
		},
		{
			Name:  "Excess",
			Value: math.Round(float64(excess) / float64(total) * 100),
			Color: "#3b82f6",
		},
		{
			Name:  "Pending",
			Value: math.Round(float64(pending) / float64(total) * 100),
			Color: "#ef4444",
		},
	}
}

// buildActionItems generates action items from items with significant variances
func (s *POvsGRNService) buildActionItems(items []models.POvsGRNComparisonItem) []models.ActionItem {
	var actionItems []models.ActionItem

	// Sort items by absolute variance descending
	type itemWithAbsVariance struct {
		item        models.POvsGRNComparisonItem
		absVariance float64
	}

	var sortable []itemWithAbsVariance
	for _, item := range items {
		if item.Status != "MATCHED" && item.Status != "PENDING" {
			sortable = append(sortable, itemWithAbsVariance{
				item:        item,
				absVariance: math.Abs(item.Variance),
			})
		}
	}

	sort.Slice(sortable, func(i, j int) bool {
		return sortable[i].absVariance > sortable[j].absVariance
	})

	// Generate action items for the top variances (max 5)
	maxActions := 5
	if len(sortable) < maxActions {
		maxActions = len(sortable)
	}

	for i := 0; i < maxActions; i++ {
		item := sortable[i].item
		actionID := fmt.Sprintf("action-%d", i+1)

		var actionItem models.ActionItem
		if item.Variance < 0 {
			// Shortage
			severity := "warning"
			title := "Review Required"
			if math.Abs(item.VariancePct) > 10 {
				severity = "critical"
				title = "High Discrepancy"
			}
			actionItem = models.ActionItem{
				ID:          actionID,
				Type:        severity,
				Title:       title,
				Description: fmt.Sprintf("%s: %.2f variance (%.1f%%). Follow up regarding shortage on #%s for %s.", item.SupplierName, item.Variance, item.VariancePct, item.PONumber, item.ProductName),
				POID:        item.POID,
				PONumber:    item.PONumber,
			}
		} else {
			// Excess
			actionItem = models.ActionItem{
				ID:          actionID,
				Type:        "info",
				Title:       "Excess Delivery",
				Description: fmt.Sprintf("%s: +%.2f variance (%.1f%%). Review excess delivery for #%s for %s.", item.SupplierName, item.Variance, item.VariancePct, item.PONumber, item.ProductName),
				POID:        item.POID,
				PONumber:    item.PONumber,
			}
		}

		actionItems = append(actionItems, actionItem)
	}

	return actionItems
}
