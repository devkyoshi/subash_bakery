// Backend API response wrapper
export interface ApiResponse<T = any> {
  success: boolean;
  message: string;
  data: T;
  error?: string;
  code?: string;
}

export interface PaginatedData<T> {
  data: T[];
  pagination: {
    page: number;
    limit: number;
    total: number;
    total_pages: number;
  };
}

export interface PaginatedResponse<T> extends ApiResponse<PaginatedData<T>> {}
