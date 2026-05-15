package models

// ReorderItem represents a single product's reorder status
type ReorderItem struct {
	ID            string  `json:"id"`            // Product SKU
	Name          string  `json:"name"`          // Product name
	Unit          string  `json:"unit"`          // Unit abbreviation
	Priority      string  `json:"priority"`      // CRITICAL, WARNING, NORMAL
	CurrentStock  float64 `json:"currentStock"`  // Current quantity on hand
	MinLevel      int     `json:"minLevel"`      // Minimum stock / reorder level
	RemainingDays int     `json:"remainingDays"` // Estimated days until stockout
	Pending       string  `json:"pending"`       // Pending order quantity text
	SugQty        int     `json:"sugQty"`        // Suggested reorder quantity
	LeadTime      string  `json:"leadTime"`      // Lead time as string (e.g., "5 Days")
}

// ConsumptionRow represents a single category's consumption analysis
type ConsumptionRow struct {
	Category string `json:"category"` // Category name
	AvgDaily string `json:"avgDaily"` // Average daily consumption string
	Trend    string `json:"trend"`    // Trend percentage string
	TrendDir string `json:"trendDir"` // "up", "down", "neutral"
	Forecast string `json:"forecast"` // Forecasted monthly usage string
}

// ReorderMetrics represents the summary metrics for the reorder status report
type ReorderMetrics struct {
	CriticalCount int `json:"critical_count"`
	WarningCount  int `json:"warning_count"`
	NormalCount   int `json:"normal_count"`
}

// ReorderStatusReportResponse is the full response for the reorder status report
type ReorderStatusReportResponse struct {
	Metrics         ReorderMetrics   `json:"metrics"`
	Items           []ReorderItem    `json:"items"`
	ConsumptionData []ConsumptionRow `json:"consumption_data"`
	TotalItems      int64            `json:"total_items"`
}

// ReorderStatusFilters contains filter parameters for the reorder status report
type ReorderStatusFilters struct {
	CategoryID     string `json:"category_id,omitempty"`
	LocationID     string `json:"location_id,omitempty"`
	Priority       string `json:"priority,omitempty"` // CRITICAL, WARNING, NORMAL
	Search         string `json:"search,omitempty"`
	IncludePending bool   `json:"include_pending,omitempty"`
}
