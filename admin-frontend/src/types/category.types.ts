// Category related types and interfaces

export interface Subcategory {
  id: string;
  name: string;
  code?: string;
  description?: string;
  is_active: boolean;
  product_count: number;
  metadata?: Record<string, any>;
  created_at?: string;
  updated_at?: string;
}

export interface Category {
  id: string;
  organization_id: string;
  name: string;
  code?: string;
  description?: string;
  is_active: boolean;
  subcategories: Subcategory[];
  product_count: number;
  metadata?: Record<string, any>;
  created_at: string;
  updated_at: string;
  created_by?: string;
  updated_by?: string;
  is_deleted?: boolean;
}

export interface CreateCategoryRequest {
  organization_id: string;
  name: string;
  code?: string;
  description?: string;
  is_active?: boolean;
  subcategories?: Omit<
    Subcategory,
    "id" | "product_count" | "created_at" | "updated_at"
  >[];
  metadata?: Record<string, any>;
}

export interface CreateSubcategoryRequest {
  name: string;
  code?: string;
  description?: string;
  is_active?: boolean;
  metadata?: Record<string, any>;
}

export interface UpdateSubcategoryRequest {
  id: string;
  name?: string;
  code?: string;
  description?: string;
  is_active?: boolean;
  metadata?: Record<string, any>;
}

export interface UpdateCategoryRequest {
  name?: string;
  code?: string;
  description?: string;
  is_active?: boolean;
  metadata?: Record<string, any>;
  add_subcategories?: CreateSubcategoryRequest[];
  update_subcategories?: UpdateSubcategoryRequest[];
  remove_subcategories?: string[];
}

export interface CategoryFilters {
  organization_id: string;
  is_active?: boolean;
  q?: string;
  page?: number;
  limit?: number;
}

export interface CategoryListResponse {
  data: Category[];
  page: number;
  limit: number;
  total: number;
}
