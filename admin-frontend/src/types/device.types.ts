export interface Device {
  id: string;
  organization_id: string;
  name: string;
  mac_address: string;
  device_type: DeviceType;
  description?: string;
  location?: string;
  is_active: boolean;
  last_seen_at?: number;
  metadata?: Record<string, any>;
  created_at: string;
  updated_at: string;
  created_by?: string;
  updated_by?: string;
}

export type DeviceType =
  | "pos"
  | "tablet"
  | "mobile"
  | "desktop"
  | "kiosk"
  | "other";

export const DEVICE_TYPES: { value: DeviceType; label: string }[] = [
  { value: "pos", label: "POS Terminal" },
  { value: "tablet", label: "Tablet" },
  { value: "mobile", label: "Mobile" },
  { value: "desktop", label: "Desktop" },
  { value: "kiosk", label: "Kiosk" },
  { value: "other", label: "Other" },
];

export interface CreateDeviceRequest {
  organization_id: string;
  name: string;
  mac_address: string;
  device_type: DeviceType;
  description?: string;
  location?: string;
}

export interface UpdateDeviceRequest {
  name?: string;
  mac_address?: string;
  device_type?: DeviceType;
  description?: string;
  location?: string;
  is_active?: boolean;
}

export interface DeviceFilters {
  organization_id: string;
  search?: string;
  is_active?: boolean;
  device_type?: DeviceType;
  page?: number;
  limit?: number;
}

export interface DeviceListResponse {
  data: Device[];
  pagination: {
    page: number;
    limit: number;
    total: number;
    total_pages: number;
  };
}
