import axiosInstance from "@/lib/axios";
import { Unit } from "@/types/product.types";
import { ApiResponse } from "@/types/global.types";

export interface CreateUnitRequest {
  code: string;
  name: string;
  symbol: string;
  unit_type: string; // Changed from type
  is_base_unit: boolean; // Changed from base_unit
  is_active?: boolean;
}

export interface UpdateUnitRequest {
  name?: string;
  symbol?: string;
  unit_type?: string; // Changed from type
  is_base_unit?: boolean; // Changed from base_unit
  is_active?: boolean;
}

export interface UnitChart {
  id: string;
  from_unit_id: string;
  to_unit_id: string;
  conversion_rate: number; // Changed from factor
  organization_id: string;
  is_active: boolean;
  created_at: number;
  updated_at: number;
}

export interface CreateUnitChartRequest {
  from_unit_id: string;
  to_unit_id: string;
  conversion_rate: number; // Changed from factor
  is_active?: boolean;
}

export interface UpdateUnitChartRequest {
  from_unit_id?: string;
  to_unit_id?: string;
  conversion_rate?: number; // Changed from factor
  is_active?: boolean;
}

export interface UnitConversionRate {
  from_unit_id: string;
  to_unit_id: string;
  conversion_rate: number;
}

class UnitService {
  /**
   * Get all units
   */
  async getUnits(params?: {
    unit_type?: string;
    active_only?: boolean;
  }): Promise<Unit[]> {
    const response = await axiosInstance.get<ApiResponse<Unit[]>>("/units", {
      params,
    });
    // Handle both wrapped and unwrapped responses if necessary, but assuming standard ApiResponse
    return response.data.data;
  }

  /**
   * Get unit by ID
   */
  async getUnit(id: string): Promise<Unit> {
    const response = await axiosInstance.get<ApiResponse<Unit>>(`/units/${id}`);
    return response.data.data;
  }

  /**
   * Create a new unit
   */
  async createUnit(data: CreateUnitRequest): Promise<Unit> {
    const response = await axiosInstance.post<ApiResponse<Unit>>(
      "/units",
      data,
    );
    return response.data.data;
  }

  /**
   * Update a unit
   */
  async updateUnit(id: string, data: UpdateUnitRequest): Promise<Unit> {
    const response = await axiosInstance.put<ApiResponse<Unit>>(
      `/units/${id}`,
      data,
    );
    return response.data.data;
  }

  /**
   * Delete a unit
   */
  async deleteUnit(id: string): Promise<void> {
    await axiosInstance.delete(`/units/${id}`);
  }

  // --- Unit Chart Methods ---

  /**
   * Get all unit charts
   */
  async getUnitCharts(params?: {
    active_only?: boolean;
  }): Promise<UnitChart[]> {
    const response = await axiosInstance.get<ApiResponse<UnitChart[]>>(
      "/unit-charts",
      { params },
    );
    return response.data.data;
  }

  /**
   * Get unit chart by ID
   */
  async getUnitChart(id: string): Promise<UnitChart> {
    const response = await axiosInstance.get<ApiResponse<UnitChart>>(
      `/unit-charts/${id}`,
    );
    return response.data.data;
  }

  /**
   * Create a new unit chart
   */
  async createUnitChart(data: CreateUnitChartRequest): Promise<UnitChart> {
    const response = await axiosInstance.post<ApiResponse<UnitChart>>(
      "/unit-charts",
      data,
    );
    return response.data.data;
  }

  /**
   * Update a unit chart
   */
  async updateUnitChart(
    id: string,
    data: UpdateUnitChartRequest,
  ): Promise<UnitChart> {
    const response = await axiosInstance.put<ApiResponse<UnitChart>>(
      `/unit-charts/${id}`,
      data,
    );
    return response.data.data;
  }

  /**
   * Delete a unit chart
   */
  async deleteUnitChart(id: string): Promise<void> {
    await axiosInstance.delete(`/unit-charts/${id}`);
  }

  /**
   * Get conversion rate
   */
  async getConversionRate(
    fromUnitId: string,
    toUnitId: string,
  ): Promise<UnitConversionRate> {
    const response = await axiosInstance.get<ApiResponse<UnitConversionRate>>(
      "/unit-charts/conversion-rate",
      {
        params: { from_unit_id: fromUnitId, to_unit_id: toUnitId },
      },
    );
    return response.data.data;
  }
}

export const unitService = new UnitService();
