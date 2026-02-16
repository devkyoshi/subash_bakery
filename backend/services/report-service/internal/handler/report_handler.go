package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/yourusername/erp-system/services/report-service/internal/models"
	"github.com/yourusername/erp-system/services/report-service/internal/service"
	"github.com/yourusername/erp-system/shared/middleware"
	"github.com/yourusername/erp-system/shared/utils"
)

type ReportHandler struct {
	povsGRNService       *service.POvsGRNService
	stockLevelService    *service.StockLevelService
	reorderStatusService *service.ReorderStatusService
	exportService        *service.ExportService
}

func NewReportHandler(povsGRNService *service.POvsGRNService, stockLevelService *service.StockLevelService, reorderStatusService *service.ReorderStatusService, exportService *service.ExportService) *ReportHandler {
	return &ReportHandler{
		povsGRNService:       povsGRNService,
		stockLevelService:    stockLevelService,
		reorderStatusService: reorderStatusService,
		exportService:        exportService,
	}
}

// RegisterRoutes registers all report routes
func (h *ReportHandler) RegisterRoutes(router *gin.RouterGroup, jwtManager *utils.JWTManager) {
	protected := router.Group("")
	protected.Use(middleware.AuthMiddleware(jwtManager))

	// PO vs GRN Comparison
	protected.GET("/organizations/:org_id/reports/po-vs-grn", h.GetPOvsGRNComparison)
	protected.GET("/organizations/:org_id/reports/po-vs-grn/export/excel", h.ExportPOvsGRNExcel)
	protected.GET("/organizations/:org_id/reports/po-vs-grn/export/pdf", h.ExportPOvsGRNPDF)

	// Stock Level Reports
	protected.GET("/organizations/:org_id/reports/stock-levels", h.GetStockLevelReport)
	protected.GET("/organizations/:org_id/reports/stock-levels/export/excel", h.ExportStockLevelExcel)
	protected.GET("/organizations/:org_id/reports/stock-levels/export/pdf", h.ExportStockLevelPDF)

	// Reorder Status Reports
	protected.GET("/organizations/:org_id/reports/reorder-status", h.GetReorderStatusReport)
}

// parseReportFilters extracts filter parameters from the request
func parseReportFilters(c *gin.Context) models.ReportFilters {
	filters := models.ReportFilters{}

	if startStr := c.Query("start_date"); startStr != "" {
		if t, err := time.Parse("2006-01-02", startStr); err == nil {
			filters.StartDate = &t
		}
	}

	if endStr := c.Query("end_date"); endStr != "" {
		if t, err := time.Parse("2006-01-02", endStr); err == nil {
			filters.EndDate = &t
		}
	}

	filters.SupplierID = c.Query("supplier_id")
	filters.Status = c.Query("status")
	filters.LocationID = c.Query("location_id")

	return filters
}

// GetPOvsGRNComparison returns the PO vs GRN comparison report data
func (h *ReportHandler) GetPOvsGRNComparison(c *gin.Context) {
	orgIDStr := c.Param("org_id")
	orgID, err := primitive.ObjectIDFromHex(orgIDStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "INVALID_ORG_ID", "Invalid organization ID", nil)
		return
	}

	page := utils.GetPageParam(c)
	limit := utils.GetLimitParam(c)
	filters := parseReportFilters(c)

	report, err := h.povsGRNService.GetPOvsGRNComparison(c.Request.Context(), orgID, filters, page, limit)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "REPORT_ERROR", err.Error(), nil)
		return
	}

	// Return the full report with pagination info
	c.JSON(http.StatusOK, utils.Response{
		Success: true,
		Data: utils.PaginationResponse{
			Data: report,
			Pagination: utils.PaginationMetadata{
				Page:       page,
				Limit:      limit,
				Total:      report.TotalItems,
				TotalPages: calculateTotalPages(report.TotalItems, limit),
			},
		},
		Message:   "PO vs GRN comparison report retrieved successfully",
		Timestamp: time.Now(),
	})
}

// ExportPOvsGRNExcel exports the PO vs GRN comparison as an Excel file
func (h *ReportHandler) ExportPOvsGRNExcel(c *gin.Context) {
	orgIDStr := c.Param("org_id")
	orgID, err := primitive.ObjectIDFromHex(orgIDStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "INVALID_ORG_ID", "Invalid organization ID", nil)
		return
	}

	filters := parseReportFilters(c)

	items, metrics, err := h.povsGRNService.GetAllPOvsGRNForExport(c.Request.Context(), orgID, filters)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "EXPORT_ERROR", err.Error(), nil)
		return
	}

	buf, err := h.exportService.GenerateExcel(items, metrics, filters)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "EXPORT_ERROR", fmt.Sprintf("Failed to generate Excel: %v", err), nil)
		return
	}

	filename := fmt.Sprintf("PO_vs_GRN_Report_%s.xlsx", time.Now().Format("20060102_150405"))
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	c.Header("Content-Length", fmt.Sprintf("%d", buf.Len()))
	c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", buf.Bytes())
}

// ExportPOvsGRNPDF exports the PO vs GRN comparison as a PDF file
func (h *ReportHandler) ExportPOvsGRNPDF(c *gin.Context) {
	orgIDStr := c.Param("org_id")
	orgID, err := primitive.ObjectIDFromHex(orgIDStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "INVALID_ORG_ID", "Invalid organization ID", nil)
		return
	}

	filters := parseReportFilters(c)

	items, metrics, err := h.povsGRNService.GetAllPOvsGRNForExport(c.Request.Context(), orgID, filters)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "EXPORT_ERROR", err.Error(), nil)
		return
	}

	buf, err := h.exportService.GeneratePDF(items, metrics, filters)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "EXPORT_ERROR", fmt.Sprintf("Failed to generate PDF: %v", err), nil)
		return
	}

	filename := fmt.Sprintf("PO_vs_GRN_Report_%s.pdf", time.Now().Format("20060102_150405"))
	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	c.Header("Content-Length", fmt.Sprintf("%d", buf.Len()))
	c.Data(http.StatusOK, "application/pdf", buf.Bytes())
}

// parseStockLevelFilters extracts stock level filter parameters from the request
func parseStockLevelFilters(c *gin.Context) models.StockLevelFilters {
	return models.StockLevelFilters{
		CategoryID:  c.Query("category_id"),
		LocationID:  c.Query("location_id"),
		StockStatus: c.Query("stock_status"),
		Search:      c.Query("search"),
	}
}

// GetStockLevelReport returns the stock level comparison report data
func (h *ReportHandler) GetStockLevelReport(c *gin.Context) {
	orgIDStr := c.Param("org_id")
	orgID, err := primitive.ObjectIDFromHex(orgIDStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "INVALID_ORG_ID", "Invalid organization ID", nil)
		return
	}

	page := utils.GetPageParam(c)
	limit := utils.GetLimitParam(c)
	filters := parseStockLevelFilters(c)

	report, err := h.stockLevelService.GetStockLevelReport(c.Request.Context(), orgID, filters, page, limit)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "REPORT_ERROR", err.Error(), nil)
		return
	}

	c.JSON(http.StatusOK, utils.Response{
		Success: true,
		Data: utils.PaginationResponse{
			Data: report,
			Pagination: utils.PaginationMetadata{
				Page:       page,
				Limit:      limit,
				Total:      report.TotalItems,
				TotalPages: calculateTotalPages(report.TotalItems, limit),
			},
		},
		Message:   "Stock level report retrieved successfully",
		Timestamp: time.Now(),
	})
}

// ExportStockLevelExcel exports the stock level report as an Excel file
func (h *ReportHandler) ExportStockLevelExcel(c *gin.Context) {
	orgIDStr := c.Param("org_id")
	orgID, err := primitive.ObjectIDFromHex(orgIDStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "INVALID_ORG_ID", "Invalid organization ID", nil)
		return
	}

	filters := parseStockLevelFilters(c)

	items, metrics, err := h.stockLevelService.GetAllStockLevelsForExport(c.Request.Context(), orgID, filters)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "EXPORT_ERROR", err.Error(), nil)
		return
	}

	buf, err := h.exportService.GenerateStockLevelExcel(items, metrics, filters)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "EXPORT_ERROR", fmt.Sprintf("Failed to generate Excel: %v", err), nil)
		return
	}

	filename := fmt.Sprintf("Stock_Level_Report_%s.xlsx", time.Now().Format("20060102_150405"))
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	c.Header("Content-Length", fmt.Sprintf("%d", buf.Len()))
	c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", buf.Bytes())
}

// ExportStockLevelPDF exports the stock level report as a PDF file
func (h *ReportHandler) ExportStockLevelPDF(c *gin.Context) {
	orgIDStr := c.Param("org_id")
	orgID, err := primitive.ObjectIDFromHex(orgIDStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "INVALID_ORG_ID", "Invalid organization ID", nil)
		return
	}

	filters := parseStockLevelFilters(c)

	items, metrics, err := h.stockLevelService.GetAllStockLevelsForExport(c.Request.Context(), orgID, filters)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "EXPORT_ERROR", err.Error(), nil)
		return
	}

	buf, err := h.exportService.GenerateStockLevelPDF(items, metrics, filters)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "EXPORT_ERROR", fmt.Sprintf("Failed to generate PDF: %v", err), nil)
		return
	}

	filename := fmt.Sprintf("Stock_Level_Report_%s.pdf", time.Now().Format("20060102_150405"))
	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	c.Header("Content-Length", fmt.Sprintf("%d", buf.Len()))
	c.Data(http.StatusOK, "application/pdf", buf.Bytes())
}

func calculateTotalPages(total int64, limit int) int {
	pages := int(total) / limit
	if int(total)%limit > 0 {
		pages++
	}
	return pages
}

// parseReorderStatusFilters extracts reorder status filter parameters from the request
func parseReorderStatusFilters(c *gin.Context) models.ReorderStatusFilters {
	includePending := c.Query("include_pending") == "true"
	return models.ReorderStatusFilters{
		CategoryID:     c.Query("category_id"),
		LocationID:     c.Query("location_id"),
		Priority:       c.Query("priority"),
		Search:         c.Query("search"),
		IncludePending: includePending,
	}
}

// GetReorderStatusReport returns the reorder status report data
func (h *ReportHandler) GetReorderStatusReport(c *gin.Context) {
	orgIDStr := c.Param("org_id")
	orgID, err := primitive.ObjectIDFromHex(orgIDStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "INVALID_ORG_ID", "Invalid organization ID", nil)
		return
	}

	page := utils.GetPageParam(c)
	limit := utils.GetLimitParam(c)
	filters := parseReorderStatusFilters(c)

	report, err := h.reorderStatusService.GetReorderStatusReport(c.Request.Context(), orgID, filters, page, limit)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "REPORT_ERROR", err.Error(), nil)
		return
	}

	c.JSON(http.StatusOK, utils.Response{
		Success: true,
		Data: utils.PaginationResponse{
			Data: report,
			Pagination: utils.PaginationMetadata{
				Page:       page,
				Limit:      limit,
				Total:      report.TotalItems,
				TotalPages: calculateTotalPages(report.TotalItems, limit),
			},
		},
		Message:   "Reorder status report retrieved successfully",
		Timestamp: time.Now(),
	})
}
