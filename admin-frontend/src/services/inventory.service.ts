import axiosInstance from "@/lib/axios";
import type {
  StockLevel,
  StockMovement,
  StockAdjustment,
  CreateStockAdjustmentRequest,
  CreateStockMovementRequest,
  Batch,
  InventoryCount,
} from "@/types/inventory.types";
import { PaginatedResponse, PaginatedData } from "@/types/global.types";

class InventoryService {
  // ============== Stock Levels ==============

  async getStockLevels(filters?: {
    product_id?: string;
    location_id?: string;
    organization_id?: string;
    search?: string;
    page?: number;
    limit?: number;
  }) {
    const params = new URLSearchParams();
    if (filters?.organization_id)
      params.append("organization_id", filters.organization_id);
    if (filters?.product_id) params.append("product_id", filters.product_id);
    if (filters?.location_id) params.append("location_id", filters.location_id);
    if (filters?.search) params.append("search", filters.search);
    if (filters?.page) params.append("page", filters.page.toString());
    if (filters?.limit) params.append("limit", filters.limit.toString());

    return axiosInstance.get<PaginatedResponse<StockLevel>>(
      `/inventory/stock-levels?${params}`,
    );
  }

  async getStockByProduct(productId: string) {
    return axiosInstance.get<{ data: StockLevel[] }>(
      `/inventory/products/${productId}/stock`,
    );
  }

  async getStockByLocation(locationId: string) {
    return axiosInstance.get<{ data: StockLevel[] }>(
      `/inventory/locations/${locationId}/stock`,
    );
  }

  async allocateStock(data: {
    product_id: string;
    location_id: string;
    quantity: number;
  }) {
    return axiosInstance.post("/inventory/stock/allocate", data);
  }

  async releaseStock(data: {
    product_id: string;
    location_id: string;
    quantity: number;
  }) {
    return axiosInstance.post("/inventory/stock/release", data);
  }

  // ============== Stock Adjustments ==============

  async createStockAdjustment(
    orgId: string,
    data: CreateStockAdjustmentRequest,
  ) {
    return axiosInstance.post<{ data: StockAdjustment }>(
      `/inventory/organizations/${orgId}/stock-adjustments`,
      data,
    );
  }

  async getStockAdjustments(
    orgId: string,
    filters?: {
      status?: string;
      location_id?: string;
      page?: number;
      limit?: number;
    },
  ) {
    const params = new URLSearchParams();
    if (filters?.status) params.append("status", filters.status);
    if (filters?.location_id) params.append("location_id", filters.location_id);
    if (filters?.page) params.append("page", filters.page.toString());
    if (filters?.limit) params.append("limit", filters.limit.toString());

    return axiosInstance.get<PaginatedResponse<StockAdjustment>>(
      `/inventory/organizations/${orgId}/stock-adjustments?${params}`,
    );
  }

  async getStockAdjustment(id: string) {
    return axiosInstance.get<{ data: StockAdjustment }>(
      `/inventory/stock-adjustments/${id}`,
    );
  }

  async updateStockAdjustment(id: string, data: CreateStockAdjustmentRequest) {
    return axiosInstance.put<{ data: StockAdjustment }>(
      `/inventory/stock-adjustments/${id}`,
      data,
    );
  }

  async approveStockAdjustment(id: string) {
    return axiosInstance.post(`/inventory/stock-adjustments/${id}/approve`);
  }

  async rejectStockAdjustment(id: string, reason: string) {
    return axiosInstance.post(`/inventory/stock-adjustments/${id}/reject`, {
      reason,
    });
  }

  async deleteStockAdjustment(id: string) {
    return axiosInstance.delete(`/inventory/stock-adjustments/${id}`);
  }

  // ============== Stock Movements ==============

  async createStockMovement(orgId: string, data: CreateStockMovementRequest) {
    return axiosInstance.post<{ data: StockMovement }>(
      `/inventory/organizations/${orgId}/stock-movements`,
      data,
    );
  }

  async getStockMovements(
    orgId: string,
    filters?: {
      product_id?: string;
      location_id?: string;
      movement_type?: string;
      page?: number;
      limit?: number;
    },
  ) {
    const params = new URLSearchParams();
    if (filters?.product_id) params.append("product_id", filters.product_id);
    if (filters?.location_id) params.append("location_id", filters.location_id);
    if (filters?.movement_type)
      params.append("movement_type", filters.movement_type);
    if (filters?.page) params.append("page", filters.page.toString());
    if (filters?.limit) params.append("limit", filters.limit.toString());

    return axiosInstance.get<PaginatedResponse<StockMovement>>(
      `/inventory/organizations/${orgId}/stock-movements?${params}`,
    );
  }

  async getStockMovement(id: string) {
    return axiosInstance.get<{ data: StockMovement }>(
      `/inventory/stock-movements/${id}`,
    );
  }

  async getStockMovementsByProduct(productId: string) {
    return axiosInstance.get<{ data: StockMovement[] }>(
      `/inventory/products/${productId}/movements`,
    );
  }

  async getStockMovementsByLocation(locationId: string) {
    return axiosInstance.get<{ data: StockMovement[] }>(
      `/inventory/locations/${locationId}/movements`,
    );
  }

  // ============== Batches ==============

  async createBatch(orgId: string, data: any) {
    return axiosInstance.post<{ data: Batch }>(
      `/inventory/organizations/${orgId}/batches`,
      data,
    );
  }

  async getBatch(id: string) {
    return axiosInstance.get<{ data: Batch }>(`/inventory/batches/${id}`);
  }

  async getBatchesByProduct(productId: string) {
    return axiosInstance.get<{ data: Batch[] }>(
      `/inventory/products/${productId}/batches`,
    );
  }

  async updateBatchQuantity(id: string, quantity: number) {
    return axiosInstance.put(`/inventory/batches/${id}/quantity`, { quantity });
  }

  // ============== Inventory Counts ==============

  async createInventoryCount(orgId: string, data: any) {
    return axiosInstance.post<{ data: InventoryCount }>(
      `/inventory/organizations/${orgId}/inventory-counts`,
      data,
    );
  }

  async getInventoryCounts(
    orgId: string,
    filters?: {
      status?: string;
      location_id?: string;
      page?: number;
      limit?: number;
    },
  ) {
    const params = new URLSearchParams();
    if (filters?.status) params.append("status", filters.status);
    if (filters?.location_id) params.append("location_id", filters.location_id);
    if (filters?.page) params.append("page", filters.page.toString());
    if (filters?.limit) params.append("limit", filters.limit.toString());

    return axiosInstance.get<{ data: InventoryCount[] }>(
      `/inventory/organizations/${orgId}/inventory-counts?${params}`,
    );
  }

  async getInventoryCount(id: string) {
    return axiosInstance.get<{ data: InventoryCount }>(
      `/inventory/inventory-counts/${id}`,
    );
  }

  async updateCountItem(id: string, data: any) {
    return axiosInstance.post(`/inventory/inventory-counts/${id}/items`, data);
  }

  async completeInventoryCount(id: string) {
    return axiosInstance.post(`/inventory/inventory-counts/${id}/complete`);
  }

  async cancelInventoryCount(id: string) {
    return axiosInstance.post(`/inventory/inventory-counts/${id}/cancel`);
  }

  async deleteInventoryCount(id: string) {
    return axiosInstance.delete(`/inventory/inventory-counts/${id}`);
  }
}

export const inventoryService = new InventoryService();
