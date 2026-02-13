package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ReportType defines the type of report
type ReportType string

const (
	ReportTypePOvsGRN      ReportType = "po_vs_grn"
	ReportTypeReorderStatus ReportType = "reorder_status"
)

// ReportFormat defines the export format
type ReportFormat string

const (
	ReportFormatPDF   ReportFormat = "pdf"
	ReportFormatExcel ReportFormat = "excel"
)

// GeneratedReport tracks generated reports
type GeneratedReport struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	OrganizationID primitive.ObjectID `bson:"organization_id" json:"organization_id"`
	Type           ReportType         `bson:"type" json:"type"`
	Format         ReportFormat       `bson:"format" json:"format"`
	Title          string             `bson:"title" json:"title"`
	Filters        ReportFilters      `bson:"filters" json:"filters"`
	GeneratedBy    primitive.ObjectID `bson:"generated_by" json:"generated_by"`
	GeneratedAt    time.Time          `bson:"generated_at" json:"generated_at"`
	FileSize       int64              `bson:"file_size" json:"file_size"`
	FilePath       string             `bson:"file_path" json:"file_path"`
	ExpiresAt      time.Time          `bson:"expires_at" json:"expires_at"`
}

// ReportFilters holds the filter criteria for a report
type ReportFilters struct {
	StartDate  *time.Time `bson:"start_date,omitempty" json:"start_date,omitempty"`
	EndDate    *time.Time `bson:"end_date,omitempty" json:"end_date,omitempty"`
	SupplierID string     `bson:"supplier_id,omitempty" json:"supplier_id,omitempty"`
	Status     string     `bson:"status,omitempty" json:"status,omitempty"`
	LocationID string     `bson:"location_id,omitempty" json:"location_id,omitempty"`
}

// POvsGRNComparisonItem represents a single line-item comparison between a PO item and its GRN receipt
type POvsGRNComparisonItem struct {
	POID          string  `json:"po_id"`
	PONumber      string  `json:"po_number"`
	OrderDate     string  `json:"order_date"`
	SupplierID    string  `json:"supplier_id"`
	SupplierName  string  `json:"supplier_name"`
	ProductID     string  `json:"product_id"`
	SKU           string  `json:"sku"`
	ProductName   string  `json:"product_name"`
	POQty         float64 `json:"po_qty"`
	GRNQty        float64 `json:"grn_qty"`
	AcceptedQty   float64 `json:"accepted_qty"`
	RejectedQty   float64 `json:"rejected_qty"`
	Variance      float64 `json:"variance"`
	VariancePct   float64 `json:"variance_pct"`
	UnitPrice     float64 `json:"unit_price"`
	POValue       float64 `json:"po_value"`
	GRNValue      float64 `json:"grn_value"`
	ValueVariance float64 `json:"value_variance"`
	Status        string  `json:"status"` // MATCHED, PARTIAL, EXCESS, PENDING
}

// POvsGRNMetrics holds the summary metrics for the PO vs GRN comparison
type POvsGRNMetrics struct {
	TotalPOs         int     `json:"total_pos"`
	CompletedPOs     int     `json:"completed_pos"`
	PartialPOs       int     `json:"partial_pos"`
	PendingPOs       int     `json:"pending_pos"`
	ExcessPOs        int     `json:"excess_pos"`
	TotalVariance    float64 `json:"total_variance"`
	TotalPOValue     float64 `json:"total_po_value"`
	TotalGRNValue    float64 `json:"total_grn_value"`
	VariancePercent  float64 `json:"variance_percent"`
	CompletedPercent float64 `json:"completed_percent"`
}

// VarianceDistribution holds the variance distribution for the chart
type VarianceDistribution struct {
	Name  string  `json:"name"`
	Value float64 `json:"value"`
	Color string  `json:"color"`
}

// ActionItem represents an action item derived from PO vs GRN analysis
type ActionItem struct {
	ID          string `json:"id"`
	Type        string `json:"type"` // critical, warning, info
	Title       string `json:"title"`
	Description string `json:"description"`
	POID        string `json:"po_id"`
	PONumber    string `json:"po_number"`
}

// POvsGRNReportResponse is the complete response for the PO vs GRN comparison page
type POvsGRNReportResponse struct {
	Metrics              POvsGRNMetrics          `json:"metrics"`
	VarianceDistribution []VarianceDistribution  `json:"variance_distribution"`
	Items                []POvsGRNComparisonItem `json:"items"`
	ActionItems          []ActionItem            `json:"action_items"`
	TotalItems           int64                   `json:"total_items"`
}
