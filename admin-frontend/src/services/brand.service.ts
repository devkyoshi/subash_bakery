import axiosInstance from "@/lib/axios";
import {
  Brand,
  CreateBrandRequest,
  UpdateBrandRequest,
  BrandFilters,
  BrandListResponse,
} from "@/types/brand.types";
import { ApiResponse } from "@/types/global.types";

class BrandService {
  /**
   * Create a new brand
   */
  async createBrand(data: CreateBrandRequest): Promise<Brand> {
    const response = await axiosInstance.post<ApiResponse<Brand>>(
      "/brands",
      data,
    );
    return response.data.data;
  }

  /**
   * Get brands with filters and pagination
   */
  async getBrands(filters: BrandFilters): Promise<BrandListResponse> {
    const params = new URLSearchParams();

    // Required parameter
    params.append("org_id", filters.organization_id);

    // Optional filters
    if (filters.is_active !== undefined) {
      params.append("is_active", filters.is_active.toString());
    }
    if (filters.q) params.append("q", filters.q);

    // Pagination
    params.append("page", (filters.page || 1).toString());
    params.append("limit", (filters.limit || 10).toString());

    // Switch between search and list endpoints based on query presence
    // Actually the search endpoint is specialized, let's verify if list supports q
    // Handler has ListBrands (lines 96) which takes 'q'. Wait, line 121 says it uses it.
    // AND SearchBrands (lines 221) which is seemingly redundant or for more complex search?
    // ListBrands uses brandRepo.FindByOrganization which might not search.
    // SearchBrands uses brandRepo.Search.

    // Let's use /brands/search if q is present, otherwise /brands
    // But standard REST APIs usually support q on list endpoint.
    // The handler ListBrands logic (line 118) creates BrandFilter with Query.
    // Line 134 calls GetBrandsByOrganization.
    // Line 99 in Service: if filter.Query != "" call Search.

    // So calling /brands with q works!

    const response = await axiosInstance.get<ApiResponse<any>>(
      `/brands?${params.toString()}`,
    );

    // Handler returns { brands: [], pagination: {} } inside data
    // My BrandListResponse matches this shape naturally?
    // Wait, typical pattern in this project is response.data.data is the payload.
    // In ListBrands handler: response := gin.H{ "brands": brands, "pagination": ... }
    // success response wraps it in 'data'.

    return response.data.data as unknown as BrandListResponse;
  }

  /**
   * Get a single brand by ID
   */
  async getBrand(id: string): Promise<Brand> {
    const response = await axiosInstance.get<ApiResponse<Brand>>(
      `/brands/${id}`,
    );
    return response.data.data;
  }

  /**
   * Update a brand
   */
  async updateBrand(id: string, data: UpdateBrandRequest): Promise<Brand> {
    const response = await axiosInstance.put<ApiResponse<Brand>>(
      `/brands/${id}`,
      data,
    );
    return response.data.data;
  }

  /**
   * Delete a brand
   */
  async deleteBrand(id: string): Promise<void> {
    await axiosInstance.delete(`/brands/${id}`);
  }
}

export const brandService = new BrandService();
