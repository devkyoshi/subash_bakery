import axiosInstance from "@/lib/axios";
import {
  Device,
  CreateDeviceRequest,
  UpdateDeviceRequest,
  DeviceFilters,
  DeviceListResponse,
} from "@/types/device.types";
import { ApiResponse } from "@/types/global.types";

class DeviceService {
  async createDevice(data: CreateDeviceRequest): Promise<Device> {
    const response = await axiosInstance.post<ApiResponse<Device>>(
      "/org-devices",
      data,
    );
    return response.data.data;
  }

  async getDevices(filters: DeviceFilters): Promise<DeviceListResponse> {
    const params = new URLSearchParams();
    params.append("organization_id", filters.organization_id);
    if (filters.is_active !== undefined)
      params.append("is_active", filters.is_active.toString());
    if (filters.search) params.append("search", filters.search);
    if (filters.device_type) params.append("device_type", filters.device_type);
    params.append("page", (filters.page || 1).toString());
    params.append("limit", (filters.limit || 10).toString());

    const response = await axiosInstance.get<ApiResponse<any>>(
      `/org-devices?${params.toString()}`,
    );
    return response.data.data as unknown as DeviceListResponse;
  }

  async getDevice(id: string): Promise<Device> {
    const response = await axiosInstance.get<ApiResponse<Device>>(
      `/org-devices/${id}`,
    );
    return response.data.data;
  }

  async updateDevice(id: string, data: UpdateDeviceRequest): Promise<Device> {
    const response = await axiosInstance.put<ApiResponse<Device>>(
      `/org-devices/${id}`,
      data,
    );
    return response.data.data;
  }

  async deleteDevice(id: string): Promise<void> {
    await axiosInstance.delete(`/org-devices/${id}`);
  }
}

export const deviceService = new DeviceService();
