import axiosInstance from "@/lib/axios";
import { ApiResponse } from "@/types/global.types";
import {
  POvsGRNReportResponse,
  ReportFilters,
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
}

export const reportService = new ReportService();
export default reportService;
