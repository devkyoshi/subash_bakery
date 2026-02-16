import axiosInstance from "@/lib/axios";
import {
  POvsGRNReportResponse,
  ReportFilters,
  StockLevelReportResponse,
  StockLevelFilters,
  ReorderStatusReportResponse,
  ReorderStatusFilters,
} from "@/types/report.types";

interface PaginatedReportResponse {
  success: boolean;
  message: string;
  data: {
    data: POvsGRNReportResponse;
    pagination: {
      page: number;
      limit: number;
      total: number;
      total_pages: number;
    };
  };
}

interface PaginatedStockLevelResponse {
  success: boolean;
  message: string;
  data: {
    data: StockLevelReportResponse;
    pagination: {
      page: number;
      limit: number;
      total: number;
      total_pages: number;
    };
  };
}

interface PaginatedReorderStatusResponse {
  success: boolean;
  message: string;
  data: {
    data: ReorderStatusReportResponse;
    pagination: {
      page: number;
      limit: number;
      total: number;
      total_pages: number;
    };
  };
}

class ReportService {
  /**
   * Get PO vs GRN comparison report data
   */
  async getPOvsGRNComparison(
    orgId: string,
    params?: ReportFilters & { page?: number; limit?: number }
  ): Promise<PaginatedReportResponse> {
    const response = await axiosInstance.get<PaginatedReportResponse>(
      `/organizations/${orgId}/reports/po-vs-grn`,
      { params }
    );
    return response.data;
  }

  /**
   * Export PO vs GRN comparison as Excel
   */
  async exportPOvsGRNExcel(
    orgId: string,
    filters?: ReportFilters
  ): Promise<Blob> {
    const response = await axiosInstance.get(
      `/organizations/${orgId}/reports/po-vs-grn/export/excel`,
      {
        params: filters,
        responseType: "blob",
      }
    );
    return response.data;
  }

  /**
   * Export PO vs GRN comparison as PDF
   */
  async exportPOvsGRNPDF(
    orgId: string,
    filters?: ReportFilters
  ): Promise<Blob> {
    const response = await axiosInstance.get(
      `/organizations/${orgId}/reports/po-vs-grn/export/pdf`,
      {
        params: filters,
        responseType: "blob",
      }
    );
    return response.data;
  }

  // ============================================================
  // Stock Level Report
  // ============================================================

  async getStockLevelReport(
    orgId: string,
    params?: StockLevelFilters & { page?: number; limit?: number }
  ): Promise<PaginatedStockLevelResponse> {
    const response = await axiosInstance.get<PaginatedStockLevelResponse>(
      `/organizations/${orgId}/reports/stock-levels`,
      { params }
    );
    return response.data;
  }

  async exportStockLevelExcel(
    orgId: string,
    filters?: StockLevelFilters
  ): Promise<Blob> {
    const response = await axiosInstance.get(
      `/organizations/${orgId}/reports/stock-levels/export/excel`,
      {
        params: filters,
        responseType: "blob",
      }
    );
    return response.data;
  }

  async exportStockLevelPDF(
    orgId: string,
    filters?: StockLevelFilters
  ): Promise<Blob> {
    const response = await axiosInstance.get(
      `/organizations/${orgId}/reports/stock-levels/export/pdf`,
      {
        params: filters,
        responseType: "blob",
      }
    );
    return response.data;
  }

  // ============================================================
  // Reorder Status Report
  // ============================================================

  async getReorderStatusReport(
    orgId: string,
    params?: ReorderStatusFilters & { page?: number; limit?: number }
  ): Promise<PaginatedReorderStatusResponse> {
    const response = await axiosInstance.get<PaginatedReorderStatusResponse>(
      `/organizations/${orgId}/reports/reorder-status`,
      { params }
    );
    return response.data;
  }
}

export const reportService = new ReportService();
export default reportService;
