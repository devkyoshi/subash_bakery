package models

// StockLevelComparisonItem represents a single product's stock level comparison
type StockLevelComparisonItem struct {
	ProductID    string  `json:"product_id"`
	SKU          string  `json:"sku"`
	ProductName  string  `json:"product_name"`
	CategoryID   string  `json:"category_id"`
	CategoryName string  `json:"category_name"`
	LocationID   string  `json:"location_id"`
	LocationName string  `json:"location_name"`
	Unit         string  `json:"unit"`
	SystemQty    float64 `json:"system_qty"`    // Current quantity_on_hand from stock_levels
	AvailableQty float64 `json:"available_qty"` // quantity_available (on_hand - allocated)
	AllocatedQty float64 `json:"allocated_qty"` // quantity_allocated
	InTransitQty float64 `json:"in_transit_qty"`
	ReorderLevel int     `json:"reorder_level"`
	MinStock     int     `json:"min_stock"`
	MaxStock     int     `json:"max_stock"`
	AverageCost  float64 `json:"average_cost"`
	TotalValue   float64 `json:"total_value"`
	StockStatus  string  `json:"stock_status"` // OPTIMAL, LOW, CRITICAL, OVERSTOCK, OUT_OF_STOCK
}

// StockLevelMetrics represents summary metrics for the stock level report
type StockLevelMetrics struct {
	TotalProducts   int     `json:"total_products"`
	OptimalCount    int     `json:"optimal_count"`
	LowStockCount   int     `json:"low_stock_count"`
	CriticalCount   int     `json:"critical_count"`
	OverstockCount  int     `json:"overstock_count"`
	OutOfStockCount int     `json:"out_of_stock_count"`
	TotalStockValue float64 `json:"total_stock_value"`
	TotalOnHand     float64 `json:"total_on_hand"`
	TotalAllocated  float64 `json:"total_allocated"`
	TotalAvailable  float64 `json:"total_available"`
}

// StockStatusDistribution represents the distribution of stock statuses
type StockStatusDistribution struct {
	Name  string  `json:"name"`
	Value float64 `json:"value"`
	Color string  `json:"color"`
}

// StockLevelReportResponse is the full response for the stock level report
type StockLevelReportResponse struct {
	Metrics            StockLevelMetrics         `json:"metrics"`
	StatusDistribution []StockStatusDistribution `json:"status_distribution"`
	Items              []StockLevelComparisonItem `json:"items"`
	TotalItems         int64                     `json:"total_items"`
}

// StockLevelFilters contains filter parameters for the stock level report
type StockLevelFilters struct {
	CategoryID  string `json:"category_id,omitempty"`
	LocationID  string `json:"location_id,omitempty"`
	StockStatus string `json:"stock_status,omitempty"` // OPTIMAL, LOW, CRITICAL, OVERSTOCK, OUT_OF_STOCK
	Search      string `json:"search,omitempty"`       // search by product name or SKU
}
