package service

import (
	"bytes"
	"fmt"
	"math"
	"time"

	"github.com/jung-kurt/gofpdf"
	"github.com/xuri/excelize/v2"

	"github.com/yourusername/erp-system/services/report-service/internal/models"
)

type ExportService struct{}

func NewExportService() *ExportService {
	return &ExportService{}
}

// GenerateExcel creates an Excel file from PO vs GRN comparison data
func (s *ExportService) GenerateExcel(
	items []models.POvsGRNComparisonItem,
	metrics *models.POvsGRNMetrics,
	filters models.ReportFilters,
) (*bytes.Buffer, error) {
	f := excelize.NewFile()
	defer f.Close()

	// ====== Summary Sheet ======
	summarySheet := "Summary"
	f.SetSheetName("Sheet1", summarySheet)

	// Title
	f.SetCellValue(summarySheet, "A1", "PO vs GRN Comparison Report")
	f.SetCellValue(summarySheet, "A2", fmt.Sprintf("Generated: %s", time.Now().Format("02 Jan 2006, 03:04 PM")))

	// Date range
	dateRange := "All Time"
	if filters.StartDate != nil && filters.EndDate != nil {
		dateRange = fmt.Sprintf("%s to %s", filters.StartDate.Format("02 Jan 2006"), filters.EndDate.Format("02 Jan 2006"))
	} else if filters.StartDate != nil {
		dateRange = fmt.Sprintf("From %s", filters.StartDate.Format("02 Jan 2006"))
	} else if filters.EndDate != nil {
		dateRange = fmt.Sprintf("Until %s", filters.EndDate.Format("02 Jan 2006"))
	}
	f.SetCellValue(summarySheet, "A3", fmt.Sprintf("Period: %s", dateRange))

	// Metrics
	f.SetCellValue(summarySheet, "A5", "Metric")
	f.SetCellValue(summarySheet, "B5", "Value")

	metricsData := [][]interface{}{
		{"Total POs", metrics.TotalPOs},
		{"Completed POs", metrics.CompletedPOs},
		{"Partial POs", metrics.PartialPOs},
		{"Pending POs", metrics.PendingPOs},
		{"Excess Deliveries", metrics.ExcessPOs},
		{"Total PO Value", metrics.TotalPOValue},
		{"Total GRN Value", metrics.TotalGRNValue},
		{"Total Variance", metrics.TotalVariance},
		{"Variance %", fmt.Sprintf("%.2f%%", metrics.VariancePercent)},
		{"Completion Rate", fmt.Sprintf("%.1f%%", metrics.CompletedPercent)},
	}

	for i, row := range metricsData {
		f.SetCellValue(summarySheet, fmt.Sprintf("A%d", 6+i), row[0])
		f.SetCellValue(summarySheet, fmt.Sprintf("B%d", 6+i), row[1])
	}

	// Style the summary header
	titleStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true, Size: 14},
	})
	f.SetCellStyle(summarySheet, "A1", "A1", titleStyle)

	headerStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Color: "FFFFFF"},
		Fill:      excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"4472C4"}},
		Alignment: &excelize.Alignment{Horizontal: "center"},
	})
	f.SetCellStyle(summarySheet, "A5", "B5", headerStyle)

	// ====== Detail Sheet ======
	detailSheet := "PO vs GRN Detail"
	f.NewSheet(detailSheet)

	// Headers
	headers := []string{
		"PO Number", "Date", "Supplier", "SKU", "Product",
		"PO Qty", "GRN Qty", "Accepted Qty", "Rejected Qty",
		"Variance", "Variance %", "Unit Price",
		"PO Value", "GRN Value", "Value Variance", "Status",
	}
	for i, h := range headers {
		col, _ := excelize.ColumnNumberToName(i + 1)
		f.SetCellValue(detailSheet, fmt.Sprintf("%s1", col), h)
	}
	lastCol, _ := excelize.ColumnNumberToName(len(headers))
	f.SetCellStyle(detailSheet, "A1", fmt.Sprintf("%s1", lastCol), headerStyle)

	// Data rows
	for i, item := range items {
		row := i + 2
		f.SetCellValue(detailSheet, fmt.Sprintf("A%d", row), item.PONumber)
		f.SetCellValue(detailSheet, fmt.Sprintf("B%d", row), item.OrderDate)
		f.SetCellValue(detailSheet, fmt.Sprintf("C%d", row), item.SupplierName)
		f.SetCellValue(detailSheet, fmt.Sprintf("D%d", row), item.SKU)
		f.SetCellValue(detailSheet, fmt.Sprintf("E%d", row), item.ProductName)
		f.SetCellValue(detailSheet, fmt.Sprintf("F%d", row), item.POQty)
		f.SetCellValue(detailSheet, fmt.Sprintf("G%d", row), item.GRNQty)
		f.SetCellValue(detailSheet, fmt.Sprintf("H%d", row), item.AcceptedQty)
		f.SetCellValue(detailSheet, fmt.Sprintf("I%d", row), item.RejectedQty)
		f.SetCellValue(detailSheet, fmt.Sprintf("J%d", row), item.Variance)
		f.SetCellValue(detailSheet, fmt.Sprintf("K%d", row), fmt.Sprintf("%.2f%%", item.VariancePct))
		f.SetCellValue(detailSheet, fmt.Sprintf("L%d", row), item.UnitPrice)
		f.SetCellValue(detailSheet, fmt.Sprintf("M%d", row), item.POValue)
		f.SetCellValue(detailSheet, fmt.Sprintf("N%d", row), item.GRNValue)
		f.SetCellValue(detailSheet, fmt.Sprintf("O%d", row), item.ValueVariance)
		f.SetCellValue(detailSheet, fmt.Sprintf("P%d", row), item.Status)
	}

	// Auto-fit column widths
	for i := range headers {
		col, _ := excelize.ColumnNumberToName(i + 1)
		f.SetColWidth(detailSheet, col, col, 15)
	}

	// Conditional formatting for variance column
	redStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Color: "FF0000"},
	})
	greenStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Color: "008000"},
	})

	for i, item := range items {
		row := i + 2
		cell := fmt.Sprintf("J%d", row)
		if item.Variance < 0 {
			f.SetCellStyle(detailSheet, cell, cell, redStyle)
		} else if item.Variance > 0 {
			f.SetCellStyle(detailSheet, cell, cell, greenStyle)
		}
	}

	// Write to buffer
	buf := new(bytes.Buffer)
	if err := f.Write(buf); err != nil {
		return nil, fmt.Errorf("failed to write Excel file: %w", err)
	}

	return buf, nil
}

// GeneratePDF creates a PDF file from PO vs GRN comparison data
func (s *ExportService) GeneratePDF(
	items []models.POvsGRNComparisonItem,
	metrics *models.POvsGRNMetrics,
	filters models.ReportFilters,
) (*bytes.Buffer, error) {
	pdf := gofpdf.New("L", "mm", "A4", "") // Landscape for table width
	pdf.SetAutoPageBreak(true, 15)

	// ====== Cover / Summary Page ======
	pdf.AddPage()

	// Title
	pdf.SetFont("Arial", "B", 20)
	pdf.CellFormat(0, 12, "PO vs GRN Comparison Report", "", 1, "C", false, 0, "")
	pdf.Ln(4)

	// Subtitle
	pdf.SetFont("Arial", "", 10)
	pdf.SetTextColor(100, 100, 100)
	pdf.CellFormat(0, 6, fmt.Sprintf("Generated: %s", time.Now().Format("02 Jan 2006, 03:04 PM")), "", 1, "C", false, 0, "")

	// Date range
	dateRange := "All Time"
	if filters.StartDate != nil && filters.EndDate != nil {
		dateRange = fmt.Sprintf("%s to %s", filters.StartDate.Format("02 Jan 2006"), filters.EndDate.Format("02 Jan 2006"))
	} else if filters.StartDate != nil {
		dateRange = fmt.Sprintf("From %s", filters.StartDate.Format("02 Jan 2006"))
	} else if filters.EndDate != nil {
		dateRange = fmt.Sprintf("Until %s", filters.EndDate.Format("02 Jan 2006"))
	}
	pdf.CellFormat(0, 6, fmt.Sprintf("Period: %s", dateRange), "", 1, "C", false, 0, "")
	pdf.Ln(10)

	// Metrics summary boxes
	pdf.SetTextColor(0, 0, 0)
	pdf.SetFont("Arial", "B", 12)
	pdf.CellFormat(0, 8, "Summary", "", 1, "L", false, 0, "")
	pdf.Ln(2)

	// Metrics table
	pdf.SetFont("Arial", "B", 9)
	pdf.SetFillColor(68, 114, 196)
	pdf.SetTextColor(255, 255, 255)
	pdf.CellFormat(60, 7, "Metric", "1", 0, "C", true, 0, "")
	pdf.CellFormat(40, 7, "Value", "1", 1, "C", true, 0, "")

	pdf.SetFont("Arial", "", 9)
	pdf.SetTextColor(0, 0, 0)
	pdf.SetFillColor(240, 240, 240)

	metricsRows := [][]string{
		{"Total POs", fmt.Sprintf("%d", metrics.TotalPOs)},
		{"Completed POs", fmt.Sprintf("%d", metrics.CompletedPOs)},
		{"Partial POs", fmt.Sprintf("%d", metrics.PartialPOs)},
		{"Pending POs", fmt.Sprintf("%d", metrics.PendingPOs)},
		{"Total PO Value", fmt.Sprintf("%.2f", metrics.TotalPOValue)},
		{"Total GRN Value", fmt.Sprintf("%.2f", metrics.TotalGRNValue)},
		{"Total Variance", fmt.Sprintf("%.2f", metrics.TotalVariance)},
		{"Variance %", fmt.Sprintf("%.2f%%", metrics.VariancePercent)},
		{"Completion Rate", fmt.Sprintf("%.1f%%", metrics.CompletedPercent)},
	}

	for i, row := range metricsRows {
		fill := i%2 == 0
		pdf.CellFormat(60, 6, row[0], "1", 0, "L", fill, 0, "")
		pdf.CellFormat(40, 6, row[1], "1", 1, "R", fill, 0, "")
	}

	// ====== Detail Table Page ======
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 12)
	pdf.CellFormat(0, 8, "PO vs GRN Detail", "", 1, "L", false, 0, "")
	pdf.Ln(2)

	// Table headers
	colWidths := []float64{25, 22, 40, 25, 25, 25, 20, 22, 22, 22, 30}
	colHeaders := []string{"PO No", "Date", "Supplier", "PO Qty", "GRN Qty", "Variance", "Var%", "PO Value", "GRN Value", "Val Var", "Status"}

	pdf.SetFont("Arial", "B", 7)
	pdf.SetFillColor(68, 114, 196)
	pdf.SetTextColor(255, 255, 255)

	for i, header := range colHeaders {
		pdf.CellFormat(colWidths[i], 6, header, "1", 0, "C", true, 0, "")
	}
	pdf.Ln(-1)

	// Table rows
	pdf.SetFont("Arial", "", 7)
	pdf.SetTextColor(0, 0, 0)

	for i, item := range items {
		fill := i%2 == 0
		if fill {
			pdf.SetFillColor(245, 245, 245)
		} else {
			pdf.SetFillColor(255, 255, 255)
		}

		pdf.CellFormat(colWidths[0], 5, item.PONumber, "1", 0, "L", fill, 0, "")
		pdf.CellFormat(colWidths[1], 5, item.OrderDate, "1", 0, "C", fill, 0, "")

		// Truncate supplier name if too long
		supplierName := item.SupplierName
		if len(supplierName) > 22 {
			supplierName = supplierName[:20] + ".."
		}
		pdf.CellFormat(colWidths[2], 5, supplierName, "1", 0, "L", fill, 0, "")

		pdf.CellFormat(colWidths[3], 5, fmt.Sprintf("%.2f", item.POQty), "1", 0, "R", fill, 0, "")
		pdf.CellFormat(colWidths[4], 5, fmt.Sprintf("%.2f", item.GRNQty), "1", 0, "R", fill, 0, "")

		// Variance with color
		varianceStr := fmt.Sprintf("%.2f", item.Variance)
		if item.Variance > 0 {
			varianceStr = "+" + varianceStr
		}
		if item.Variance < 0 {
			pdf.SetTextColor(220, 38, 38)
		} else if item.Variance > 0 {
			pdf.SetTextColor(22, 163, 74)
		}
		pdf.CellFormat(colWidths[5], 5, varianceStr, "1", 0, "R", fill, 0, "")
		pdf.SetTextColor(0, 0, 0)

		pdf.CellFormat(colWidths[6], 5, fmt.Sprintf("%.1f%%", item.VariancePct), "1", 0, "R", fill, 0, "")
		pdf.CellFormat(colWidths[7], 5, fmt.Sprintf("%.2f", item.POValue), "1", 0, "R", fill, 0, "")
		pdf.CellFormat(colWidths[8], 5, fmt.Sprintf("%.2f", item.GRNValue), "1", 0, "R", fill, 0, "")
		pdf.CellFormat(colWidths[9], 5, fmt.Sprintf("%.2f", math.Abs(item.ValueVariance)), "1", 0, "R", fill, 0, "")
		pdf.CellFormat(colWidths[10], 5, item.Status, "1", 0, "C", fill, 0, "")
		pdf.Ln(-1)

		// Check if we need a new page
		if pdf.GetY() > 190 {
			pdf.AddPage()
			pdf.SetFont("Arial", "B", 7)
			pdf.SetFillColor(68, 114, 196)
			pdf.SetTextColor(255, 255, 255)
			for i, header := range colHeaders {
				pdf.CellFormat(colWidths[i], 6, header, "1", 0, "C", true, 0, "")
			}
			pdf.Ln(-1)
			pdf.SetFont("Arial", "", 7)
			pdf.SetTextColor(0, 0, 0)
		}
	}

	// Footer
	pdf.Ln(5)
	pdf.SetFont("Arial", "I", 8)
	pdf.SetTextColor(150, 150, 150)
	pdf.CellFormat(0, 5, fmt.Sprintf("Total items: %d | Report generated by ERP System", len(items)), "", 1, "L", false, 0, "")

	// Write to buffer
	buf := new(bytes.Buffer)
	if err := pdf.Output(buf); err != nil {
		return nil, fmt.Errorf("failed to write PDF: %w", err)
	}

	return buf, nil
}
