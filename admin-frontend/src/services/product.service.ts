import axiosInstance from "@/lib/axios";
import {
  Product,
  CreateProductRequest,
  UpdateProductRequest,
  ProductFilters,
  ProductListResponse,
} from "@/types/product.types";
import { ApiResponse } from "@/types/global.types";

class ProductService {
  /**
   * Create a new product
   */
  async createProduct(data: CreateProductRequest): Promise<Product> {
    const response = await axiosInstance.post<ApiResponse<Product>>(
      "/products",
      data,
    );
    return response.data.data;
  }

  /**
   * Get products with filters and pagination
   */
  async getProducts(filters: ProductFilters): Promise<ProductListResponse> {
    const params = new URLSearchParams();

    // Required parameter
    params.append("organization_id", filters.organization_id);

    // Optional filters
    if (filters.category_id) params.append("category_id", filters.category_id);
    if (filters.subcategory_id)
      params.append("subcategory_id", filters.subcategory_id);
    if (filters.brand_id) params.append("brand_id", filters.brand_id);
    if (filters.status) params.append("status", filters.status);
    if (filters.type) params.append("type", filters.type);
    if (filters.track_inventory !== undefined) {
      params.append("track_inventory", filters.track_inventory.toString());
    }
    if (filters.location_id) params.append("location_id", filters.location_id);
    if (filters.search) params.append("search", filters.search);

    // Pagination
    params.append("page", (filters.page || 1).toString());
    params.append("limit", (filters.limit || 10).toString());

    const response = await axiosInstance.get<ApiResponse<ProductListResponse>>(
      `/products?${params.toString()}`,
    );
    return response.data.data;
  }

  /**
   * Get a single product by ID
   */
  async getProduct(id: string, locationId?: string): Promise<Product> {
    const params = locationId ? `?location_id=${locationId}` : "";
    const response = await axiosInstance.get<ApiResponse<Product>>(
      `/products/${id}${params}`,
    );
    return response.data.data;
  }

  /**
   * Get product by SKU
   */
  async getProductBySKU(sku: string, organizationId: string): Promise<Product> {
    const response = await axiosInstance.get<ApiResponse<Product>>(
      `/products/sku/${sku}?organization_id=${organizationId}`,
    );
    return response.data.data;
  }

  /**
   * Get low stock products
   */
  async getLowStockProducts(organizationId: string): Promise<Product[]> {
    const response = await axiosInstance.get<ApiResponse<Product[]>>(
      `/products/low-stock?organization_id=${organizationId}`,
    );
    return response.data.data;
  }

  /**
   * Update a product
   */
  async updateProduct(
    id: string,
    data: UpdateProductRequest,
  ): Promise<Product> {
    const response = await axiosInstance.put<ApiResponse<Product>>(
      `/products/${id}`,
      data,
    );
    return response.data.data;
  }

  /**
   * Delete a product
   */
  async deleteProduct(id: string): Promise<void> {
    await axiosInstance.delete(`/products/${id}`);
  }
}

export const productService = new ProductService();
