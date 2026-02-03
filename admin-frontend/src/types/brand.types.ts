export interface Brand {
  id: string;
  organization_id: string;
  name: string;
  code?: string;
  description?: string;
  logo_url?: string;
  website?: string;
  country?: string;
  is_active: boolean;
  metadata?: Record<string, any>;
  created_at: string;
  updated_at: string;
  created_by?: string;
  updated_by?: string;
}

export interface CreateBrandRequest {
  organization_id: string;
  name: string;
  code?: string;
  description?: string;
  logo_url?: string;
  website?: string;
  country?: string;
  is_active?: boolean;
  metadata?: Record<string, any>;
}

export interface UpdateBrandRequest {
  name?: string;
  code?: string;
  description?: string;
  logo_url?: string;
  website?: string;
  country?: string;
  is_active?: boolean;
  metadata?: Record<string, any>;
}

export interface BrandFilters {
  organization_id: string;
  is_active?: boolean;
  q?: string;
  page?: number;
  limit?: number;
}

export interface BrandListResponse {
  brands: Brand[];
  pagination: {
    page: number;
    limit: number;
    total: number;
  };
}
