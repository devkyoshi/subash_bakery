import axiosInstance from "@/lib/axios";
import {
  Category,
  CreateCategoryRequest,
  UpdateCategoryRequest,
  CategoryFilters,
  CategoryListResponse,
} from "@/types/category.types";
import { ApiResponse } from "@/types/global.types";

class CategoryService {
  /**
   * Create a new category
   */
  async createCategory(data: CreateCategoryRequest): Promise<Category> {
    const response = await axiosInstance.post<ApiResponse<Category>>(
      "/categories",
      data,
    );
    return response.data.data;
  }

  /**
   * Get categories with filters and pagination
   */
  async getCategories(filters: CategoryFilters): Promise<CategoryListResponse> {
    const params = new URLSearchParams();

    // Required parameter
    params.append("organization_id", filters.organization_id);

    // Optional filters
    if (filters.is_active !== undefined) {
      params.append("is_active", filters.is_active.toString());
    }
    if (filters.q) params.append("q", filters.q);

    // Pagination
    params.append("page", (filters.page || 1).toString());
    params.append("limit", (filters.limit || 10).toString());

    const response = await axiosInstance.get<ApiResponse<CategoryListResponse>>(
      `/categories?${params.toString()}`,
    );
    return response.data.data;
  }

  /**
   * Get distinct root-level categories
   */
  async getRootCategories(organizationId: string): Promise<Category[]> {
    const response = await axiosInstance.get<ApiResponse<Category[]>>(
      `/categories/root?organization_id=${organizationId}`,
    );
    return response.data.data;
  }

  /**
   * Get a single category by ID
   */
  async getCategory(id: string): Promise<Category> {
    const response = await axiosInstance.get<ApiResponse<Category>>(
      `/categories/${id}`,
    );
    return response.data.data;
  }

  /**
   * Update a category
   */
  async updateCategory(
    id: string,
    data: UpdateCategoryRequest,
  ): Promise<Category> {
    const response = await axiosInstance.put<ApiResponse<Category>>(
      `/categories/${id}`,
      data,
    );
    return response.data.data;
  }

  /**
   * Delete a category
   */
  async deleteCategory(id: string): Promise<void> {
    await axiosInstance.delete(`/categories/${id}`);
  }
}

export const categoryService = new CategoryService();
