import axiosInstance from "@/lib/axios";
import {
  Supplier,
  SupplierStatus,
  CreateSupplierRequest,
  PurchaseOrder,
  CreatePurchaseOrderRequest,
  POStatus,
  GoodsReceiptNote,
  GoodsReceiptNoteSummary,
  CreateGRNRequest,
} from "@/types/procurement.types";
import { ApiResponse, PaginatedResponse } from "@/types/global.types";

class ProcurementService {
  // ==================== Suppliers ====================

  /**
   * Create a new supplier
   */
  async createSupplier(
    orgId: string,
    data: CreateSupplierRequest,
  ): Promise<Supplier> {
    const response = await axiosInstance.post<ApiResponse<Supplier>>(
      `/organizations/${orgId}/suppliers`,
      data,
    );
    return response.data.data;
  }

  /**
   * List suppliers
   */
  async getSuppliers(
    orgId: string,
    params?: {
      page?: number;
      limit?: number;
      search?: string;
      status?: SupplierStatus;
    },
  ): Promise<PaginatedResponse<Supplier>> {
    const response = await axiosInstance.get<PaginatedResponse<Supplier>>(
      `/organizations/${orgId}/suppliers`,
      { params },
    );
    return response.data;
  }

  /**
   * Get a supplier by ID
   */
  async getSupplier(id: string): Promise<Supplier> {
    const response = await axiosInstance.get<ApiResponse<Supplier>>(
      `/suppliers/${id}`,
    );
    return response.data.data;
  }

  /**
   * Update a supplier
   */
  async updateSupplier(
    id: string,
    data: Partial<CreateSupplierRequest>,
  ): Promise<Supplier> {
    const response = await axiosInstance.put<ApiResponse<Supplier>>(
      `/suppliers/${id}`,
      data,
    );
    return response.data.data;
  }

  /**
   * Delete a supplier
   */
  async deleteSupplier(id: string): Promise<void> {
    await axiosInstance.delete<ApiResponse<void>>(`/suppliers/${id}`);
  }

  // ==================== Purchase Orders ====================

  /**
   * Create a new purchase order
   */
  async createPurchaseOrder(
    orgId: string,
    data: CreatePurchaseOrderRequest,
  ): Promise<PurchaseOrder> {
    const response = await axiosInstance.post<ApiResponse<PurchaseOrder>>(
      `/organizations/${orgId}/purchase-orders`,
      data,
    );
    return response.data.data;
  }

  /**
   * List purchase orders
   */
  async getPurchaseOrders(
    orgId: string,
    params?: {
      page?: number;
      limit?: number;
      search?: string;
      status?: string;
      supplier_id?: string;
    },
  ): Promise<PaginatedResponse<PurchaseOrder>> {
    const response = await axiosInstance.get<PaginatedResponse<PurchaseOrder>>(
      `/organizations/${orgId}/purchase-orders`,
      { params },
    );
    return response.data;
  }

  /**
   * Get a purchase order by ID
   */
  async getPurchaseOrder(id: string): Promise<PurchaseOrder> {
    const response = await axiosInstance.get<ApiResponse<PurchaseOrder>>(
      `/purchase-orders/${id}`,
    );
    return response.data.data;
  }

  /**
   * Update purchase order status
   */
  async updatePOStatus(id: string, status: POStatus): Promise<void> {
    await axiosInstance.put<ApiResponse<void>>(
      `/purchase-orders/${id}/status`,
      { status },
    );
  }

  /**
   * Approve a purchase order
   */
  async approvePurchaseOrder(id: string): Promise<void> {
    await axiosInstance.post<ApiResponse<void>>(
      `/purchase-orders/${id}/approve`,
    );
  }

  /**
   * Delete a purchase order
   */
  async deletePurchaseOrder(id: string): Promise<void> {
    await axiosInstance.delete<ApiResponse<void>>(`/purchase-orders/${id}`);
  }

  // ==================== Goods Receipt Notes (GRN) ====================

  /**
   * Create a new GRN
   */
  async createGRN(
    orgId: string,
    data: CreateGRNRequest,
  ): Promise<GoodsReceiptNote> {
    const response = await axiosInstance.post<ApiResponse<GoodsReceiptNote>>(
      `/organizations/${orgId}/grns`,
      data,
    );
    return response.data.data;
  }

  /**
   * List GRNs
   */
  async getGRNs(
    orgId: string,
    params?: {
      page?: number;
      limit?: number;
      search?: string;
      status?: string;
      purchase_order_id?: string;
    },
  ): Promise<PaginatedResponse<GoodsReceiptNoteSummary>> {
    const response = await axiosInstance.get<
      PaginatedResponse<GoodsReceiptNoteSummary>
    >(`/organizations/${orgId}/grns`, { params });
    return response.data;
  }

  /**
   * Get a GRN by ID
   */
  async getGRN(id: string): Promise<GoodsReceiptNote> {
    const response = await axiosInstance.get<ApiResponse<GoodsReceiptNote>>(
      `/grns/${id}`,
    );
    return response.data.data;
  }

  /**
   * Complete inspection for a GRN
   */
  async inspectGRN(
    id: string,
    data: { qc_status: string; qc_notes?: string },
  ): Promise<void> {
    await axiosInstance.post<ApiResponse<void>>(`/grns/${id}/inspect`, data);
  }
}

export const procurementService = new ProcurementService();
