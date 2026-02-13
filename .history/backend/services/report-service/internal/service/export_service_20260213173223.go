package service

import (
	"bytes"
	"fmt"
	"math"
	"os"
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

// GeneratePDF creates a branded PDF file from PO vs GRN comparison data
func (s *ExportService) GeneratePDF(
	items []models.POvsGRNComparisonItem,
	metrics *models.POvsGRNMetrics,
	filters models.ReportFilters,
) (*bytes.Buffer, error) {
	pdf := gofpdf.New("L", "mm", "A4", "") // Landscape
	pdf.SetAutoPageBreak(false, 0)
	pageW, pageH := pdf.GetPageSize()

	// Brand colors
	brandR, brandG, brandB := 180, 60, 40       // Deep warm red/brown
	accentR, accentG, accentB := 240, 130, 60   // Warm orange accent
	darkR, darkG, darkB := 45, 45, 50           // Near-black for text
	mutedR, mutedG, mutedB := 120, 120, 130     // Gray for secondary text
	headerR, headerG, headerB := 55, 65, 80     // Dark slate for table headers
	bgLightR, bgLightG, bgLightB := 250, 248, 245 // Warm off-white

	// Date range string
	dateRange := "All Time"
	if filters.StartDate != nil && filters.EndDate != nil {
		dateRange = fmt.Sprintf("%s to %s", filters.StartDate.Format("02 Jan 2006"), filters.EndDate.Format("02 Jan 2006"))
	} else if filters.StartDate != nil {
		dateRange = fmt.Sprintf("From %s", filters.StartDate.Format("02 Jan 2006"))
	} else if filters.EndDate != nil {
		dateRange = fmt.Sprintf("Until %s", filters.EndDate.Format("02 Jan 2006"))
	}

	// ============================================================
	// Helper: Draw page header (logo + brand bar) on every page
	// ============================================================
	drawPageHeader := func() {
		// Top brand bar
		pdf.SetFillColor(brandR, brandG, brandB)
		pdf.Rect(0, 0, pageW, 4, "F")

		// Accent stripe below
		pdf.SetFillColor(accentR, accentG, accentB)
		pdf.Rect(0, 4, pageW, 1.5, "F")

		// Logo
		logoPath := "assets/report/logo.png"
		if fileExists(logoPath) {
			pdf.ImageOptions(logoPath, 12, 9, 14, 0, false, gofpdf.ImageOptions{ImageType: "PNG"}, 0, "")
		}

		// Company name
		pdf.SetXY(28, 10)
		pdf.SetFont("Arial", "B", 14)
		pdf.SetTextColor(brandR, brandG, brandB)
		pdf.CellFormat(60, 6, "Subash Bakery", "", 0, "L", false, 0, "")

		// Tagline
		pdf.SetXY(28, 16)
		pdf.SetFont("Arial", "", 7)
		pdf.SetTextColor(mutedR, mutedG, mutedB)
		pdf.CellFormat(60, 4, "Enterprise Resource Planning", "", 0, "L", false, 0, "")

		// Right-aligned report info
		pdf.SetXY(pageW-90, 10)
		pdf.SetFont("Arial", "", 7)
		pdf.SetTextColor(mutedR, mutedG, mutedB)
		pdf.CellFormat(78, 4, fmt.Sprintf("Generated: %s", time.Now().Format("02 Jan 2006, 03:04 PM")), "", 0, "R", false, 0, "")
		pdf.SetXY(pageW-90, 14)
		pdf.CellFormat(78, 4, fmt.Sprintf("Period: %s", dateRange), "", 0, "R", false, 0, "")
		pdf.SetXY(pageW-90, 18)
		pdf.CellFormat(78, 4, fmt.Sprintf("Total Records: %d", len(items)), "", 0, "R", false, 0, "")

		// Separator line
		pdf.SetDrawColor(220, 220, 220)
		pdf.SetLineWidth(0.3)
		pdf.Line(12, 25, pageW-12, 25)

		pdf.SetY(28)
	}

	// ============================================================
	// Helper: Draw page footer
	// ============================================================
	drawPageFooter := func(pageNum int) {
		pdf.SetY(pageH - 15)
		pdf.SetDrawColor(220, 220, 220)
		pdf.SetLineWidth(0.2)
		pdf.Line(12, pageH-16, pageW-12, pageH-16)

		pdf.SetFont("Arial", "", 6.5)
		pdf.SetTextColor(mutedR, mutedG, mutedB)
		pdf.CellFormat(pageW/2-12, 4, "Subash Bakery - Confidential", "", 0, "L", false, 0, "")
		pdf.CellFormat(pageW/2-12, 4, fmt.Sprintf("Page %d", pageNum), "", 0, "R", false, 0, "")
	}

	// ============================================================
	// PAGE 1 — Cover / Summary
	// ============================================================
	pdf.AddPage()
	drawPageHeader()

	// Report title
	pdf.SetFont("Arial", "B", 22)
	pdf.SetTextColor(darkR, darkG, darkB)
	pdf.CellFormat(0, 12, "PO vs GRN Comparison Report", "", 1, "L", false, 0, "")
	pdf.Ln(1)

	// Subtitle bar
	pdf.SetFillColor(bgLightR, bgLightG, bgLightB)
	pdf.SetFont("Arial", "", 9)
	pdf.SetTextColor(mutedR, mutedG, mutedB)
	pdf.CellFormat(0, 7, fmt.Sprintf("    Procurement analysis for %s", dateRange), "", 1, "L", true, 0, "")
	pdf.Ln(6)

	// ---- KPI Cards Row ----
	cardW := (pageW - 24 - 16) / 5 // 5 cards with gaps
	cardH := float64(28)
	startX := float64(12)

	type kpiCard struct {
		label string
		value string
		sub   string
		r, g, b int
	}
	kpis := []kpiCard{
		{"Total POs", fmt.Sprintf("%d", metrics.TotalPOs), "Purchase Orders", 55, 120, 200},
		{"Completed", fmt.Sprintf("%d", metrics.CompletedPOs), fmt.Sprintf("%.0f%% of total", metrics.CompletedPercent), 34, 160, 90},
		{"Partial", fmt.Sprintf("%d", metrics.PartialPOs), "Awaiting delivery", 220, 160, 40},
		{"Pending", fmt.Sprintf("%d", metrics.PendingPOs), "Not yet received", 220, 70, 55},
		{"Variance", fmt.Sprintf("%.1f%%", metrics.VariancePercent), fmt.Sprintf("%.2f total", metrics.TotalVariance), brandR, brandG, brandB},
	}

	cardY := pdf.GetY()
	for i, kpi := range kpis {
		x := startX + float64(i)*(cardW+4)

		// Card background
		pdf.SetFillColor(250, 250, 252)
		pdf.SetDrawColor(230, 230, 235)
		pdf.RoundedRect(x, cardY, cardW, cardH, 2, "1234", "FD")

		// Color accent bar at top of card
		pdf.SetFillColor(kpi.r, kpi.g, kpi.b)
		pdf.Rect(x+1, cardY+1, cardW-2, 2.5, "F")

		// Value
		pdf.SetXY(x+4, cardY+7)
		pdf.SetFont("Arial", "B", 16)
		pdf.SetTextColor(darkR, darkG, darkB)
		pdf.CellFormat(cardW-8, 8, kpi.value, "", 0, "L", false, 0, "")

		// Label
		pdf.SetXY(x+4, cardY+16)
		pdf.SetFont("Arial", "B", 7.5)
		pdf.SetTextColor(kpi.r, kpi.g, kpi.b)
		pdf.CellFormat(cardW-8, 4, kpi.label, "", 0, "L", false, 0, "")

		// Sub text
		pdf.SetXY(x+4, cardY+21)
		pdf.SetFont("Arial", "", 6.5)
		pdf.SetTextColor(mutedR, mutedG, mutedB)
		pdf.CellFormat(cardW-8, 4, kpi.sub, "", 0, "L", false, 0, "")
	}

	pdf.SetY(cardY + cardH + 8)

	// ---- Financial Summary Table ----
	pdf.SetFont("Arial", "B", 11)
	pdf.SetTextColor(darkR, darkG, darkB)
	pdf.CellFormat(0, 8, "Financial Summary", "", 1, "L", false, 0, "")
	pdf.Ln(1)

	// Two-column financial metrics
	finColW := float64(135)
	finRows := [][]string{
		{"Total PO Value", fmt.Sprintf("%.2f", metrics.TotalPOValue)},
		{"Total GRN Value", fmt.Sprintf("%.2f", metrics.TotalGRNValue)},
		{"Total Variance (Absolute)", fmt.Sprintf("%.2f", metrics.TotalVariance)},
		{"Variance Percentage", fmt.Sprintf("%.2f%%", metrics.VariancePercent)},
		{"Completion Rate", fmt.Sprintf("%.1f%%", metrics.CompletedPercent)},
	}

	// Table header
	pdf.SetFont("Arial", "B", 8)
	pdf.SetFillColor(headerR, headerG, headerB)
	pdf.SetTextColor(255, 255, 255)
	pdf.CellFormat(finColW*0.6, 7, "  Metric", "1", 0, "L", true, 0, "")
	pdf.CellFormat(finColW*0.4, 7, "Value  ", "1", 1, "R", true, 0, "")

	pdf.SetFont("Arial", "", 8.5)
	for i, row := range finRows {
		if i%2 == 0 {
			pdf.SetFillColor(bgLightR, bgLightG, bgLightB)
		} else {
			pdf.SetFillColor(255, 255, 255)
		}
		pdf.SetTextColor(darkR, darkG, darkB)
		pdf.CellFormat(finColW*0.6, 6.5, "  "+row[0], "LB", 0, "L", true, 0, "")

		// Highlight variance row
		if i == 2 && metrics.TotalVariance > 0 {
			pdf.SetTextColor(220, 70, 55)
		}
		pdf.CellFormat(finColW*0.4, 6.5, row[1]+"  ", "RB", 1, "R", true, 0, "")
	}

	drawPageFooter(1)

	// ============================================================
	// PAGE 2+ — Detail Table
	// ============================================================
	pdf.AddPage()
	drawPageHeader()

	pdf.SetFont("Arial", "B", 14)
	pdf.SetTextColor(darkR, darkG, darkB)
	pdf.CellFormat(0, 10, "Detailed PO vs GRN Comparison", "", 1, "L", false, 0, "")
	pdf.Ln(2)

	// Table layout
	colWidths := []float64{24, 22, 38, 36, 22, 22, 24, 20, 22, 22, 22, 24}
	colHeaders := []string{"PO No", "Date", "Supplier", "Product", "PO Qty", "GRN Qty", "Variance", "Var%", "PO Value", "GRN Value", "Val Var", "Status"}

	drawTableHeader := func() {
		pdf.SetFont("Arial", "B", 7)
		pdf.SetFillColor(headerR, headerG, headerB)
		pdf.SetTextColor(255, 255, 255)
		for i, header := range colHeaders {
			pdf.CellFormat(colWidths[i], 7, header, "1", 0, "C", true, 0, "")
		}
		pdf.Ln(-1)
	}

	drawTableHeader()
	pageNum := 2

	// Table rows
	for i, item := range items {
		// Check for page break
		if pdf.GetY() > pageH-25 {
			drawPageFooter(pageNum)
			pageNum++
			pdf.AddPage()
			drawPageHeader()
			pdf.SetFont("Arial", "B", 11)
			pdf.SetTextColor(darkR, darkG, darkB)
			pdf.CellFormat(0, 8, "Detailed PO vs GRN Comparison (continued)", "", 1, "L", false, 0, "")
			pdf.Ln(1)
			drawTableHeader()
		}

		pdf.SetFont("Arial", "", 7)

		// Alternating row colors
		if i%2 == 0 {
			pdf.SetFillColor(bgLightR, bgLightG, bgLightB)
		} else {
			pdf.SetFillColor(255, 255, 255)
		}
		fill := true

		pdf.SetTextColor(darkR, darkG, darkB)
		pdf.CellFormat(colWidths[0], 5.5, item.PONumber, "1", 0, "L", fill, 0, "")
		pdf.CellFormat(colWidths[1], 5.5, item.OrderDate, "1", 0, "C", fill, 0, "")

		supplierName := item.SupplierName
		if len(supplierName) > 20 {
			supplierName = supplierName[:18] + ".."
		}
		pdf.CellFormat(colWidths[2], 5.5, supplierName, "1", 0, "L", fill, 0, "")

		productName := item.ProductName
		if productName == "" {
			productName = item.SKU
		}
		if len(productName) > 20 {
			productName = productName[:18] + ".."
		}
		pdf.CellFormat(colWidths[3], 5.5, productName, "1", 0, "L", fill, 0, "")

		pdf.CellFormat(colWidths[4], 5.5, fmt.Sprintf("%.2f", item.POQty), "1", 0, "R", fill, 0, "")
		pdf.CellFormat(colWidths[5], 5.5, fmt.Sprintf("%.2f", item.GRNQty), "1", 0, "R", fill, 0, "")

		// Variance with color coding
		varianceStr := fmt.Sprintf("%.2f", item.Variance)
		if item.Variance > 0 {
			varianceStr = "+" + varianceStr
		}
		if item.Variance < 0 {
			pdf.SetTextColor(220, 50, 50)
		} else if item.Variance > 0 {
			pdf.SetTextColor(34, 160, 90)
		} else {
			pdf.SetTextColor(mutedR, mutedG, mutedB)
		}
		pdf.SetFont("Arial", "B", 7)
		pdf.CellFormat(colWidths[6], 5.5, varianceStr, "1", 0, "R", fill, 0, "")
		pdf.SetFont("Arial", "", 7)
		pdf.SetTextColor(darkR, darkG, darkB)

		pdf.CellFormat(colWidths[7], 5.5, fmt.Sprintf("%.1f%%", item.VariancePct), "1", 0, "R", fill, 0, "")
		pdf.CellFormat(colWidths[8], 5.5, fmt.Sprintf("%.2f", item.POValue), "1", 0, "R", fill, 0, "")
		pdf.CellFormat(colWidths[9], 5.5, fmt.Sprintf("%.2f", item.GRNValue), "1", 0, "R", fill, 0, "")
		pdf.CellFormat(colWidths[10], 5.5, fmt.Sprintf("%.2f", math.Abs(item.ValueVariance)), "1", 0, "R", fill, 0, "")

		// Status badge with color
		switch item.Status {
		case "MATCHED":
			pdf.SetFillColor(220, 252, 231)
			pdf.SetTextColor(34, 120, 60)
		case "PARTIAL":
			pdf.SetFillColor(254, 226, 226)
			pdf.SetTextColor(185, 50, 50)
		case "EXCESS":
			pdf.SetFillColor(219, 234, 254)
			pdf.SetTextColor(40, 80, 160)
		case "PENDING":
			pdf.SetFillColor(254, 243, 199)
			pdf.SetTextColor(160, 120, 30)
		default:
			pdf.SetFillColor(245, 245, 245)
			pdf.SetTextColor(mutedR, mutedG, mutedB)
		}
		pdf.SetFont("Arial", "B", 6.5)
		pdf.CellFormat(colWidths[11], 5.5, item.Status, "1", 0, "C", true, 0, "")
		pdf.SetFont("Arial", "", 7)

		pdf.Ln(-1)
	}

	// Summary row at bottom of table — ensure it fits on current page
	if pdf.GetY()+6 > pageH-20 {
		drawPageFooter(pageNum)
		pageNum++
		pdf.AddPage()
		drawPageHeader()
		drawTableHeader()
	}
	pdf.SetFont("Arial", "B", 7)
	pdf.SetFillColor(headerR, headerG, headerB)
	pdf.SetTextColor(255, 255, 255)
	pdf.CellFormat(colWidths[0]+colWidths[1]+colWidths[2]+colWidths[3], 6, fmt.Sprintf("  TOTAL (%d items)", len(items)), "1", 0, "L", true, 0, "")
	var totalPOQty, totalGRNQty, totalVar, totalPOVal, totalGRNVal, totalValVar float64
	for _, item := range items {
		totalPOQty += item.POQty
		totalGRNQty += item.GRNQty
		totalVar += item.Variance
		totalPOVal += item.POValue
		totalGRNVal += item.GRNValue
		totalValVar += item.ValueVariance
	}
	pdf.CellFormat(colWidths[4], 6, fmt.Sprintf("%.2f", totalPOQty), "1", 0, "R", true, 0, "")
	pdf.CellFormat(colWidths[5], 6, fmt.Sprintf("%.2f", totalGRNQty), "1", 0, "R", true, 0, "")
	pdf.CellFormat(colWidths[6], 6, fmt.Sprintf("%.2f", totalVar), "1", 0, "R", true, 0, "")
	pdf.CellFormat(colWidths[7], 6, "", "1", 0, "C", true, 0, "")
	pdf.CellFormat(colWidths[8], 6, fmt.Sprintf("%.2f", totalPOVal), "1", 0, "R", true, 0, "")
	pdf.CellFormat(colWidths[9], 6, fmt.Sprintf("%.2f", totalGRNVal), "1", 0, "R", true, 0, "")
	pdf.CellFormat(colWidths[10], 6, fmt.Sprintf("%.2f", math.Abs(totalValVar)), "1", 0, "R", true, 0, "")
	pdf.CellFormat(colWidths[11], 6, "", "1", 0, "C", true, 0, "")
	pdf.Ln(-1)

	drawPageFooter(pageNum)

	// Write to buffer
	buf := new(bytes.Buffer)
	if err := pdf.Output(buf); err != nil {
		return nil, fmt.Errorf("failed to write PDF: %w", err)
	}

	return buf, nil
}

// fileExists checks if a file exists at the given path
func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
