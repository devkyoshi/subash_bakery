import axiosInstance from "@/lib/axios";

export interface User {
  id: string;
  organization_id: string;
  email: string;
  first_name: string;
  last_name: string;
  phone?: string;
  avatar?: string;
  role_id: string;
  is_active: boolean;
  is_email_verified: boolean;
  google_id?: string;
  last_login_at?: string;
  created_at: string;
  updated_at: string;
}

export interface Pagination {
  page: number;
  limit: number;
  total: number;
  total_pages: number;
}

export interface GetUsersResponse {
  success: boolean;
  message: string;
  data: {
    data: User[];
    pagination: Pagination;
  };
  timestamp: string;
}

class UserService {
  /**
   * Get all users in an organization
   */
  async getUsers(
    organizationId: string,
    params?: {
      page?: number;
      limit?: number;
      search?: string;
      role?: string;
      status?: string;
      company_id?: string;
    },
  ): Promise<GetUsersResponse> {
    const response = await axiosInstance.get<GetUsersResponse>(
      `/organizations/${organizationId}/users`,
      { params },
    );
    return response.data;
  }

  /**
   * Get user by ID
   */
  async getUserById(userId: string): Promise<User> {
    const response = await axiosInstance.get(`/users/${userId}`);
    return response.data.data;
  }

  /**
   * Create new user
   */
  async createUser(userData: Partial<User>): Promise<User> {
    const response = await axiosInstance.post("/users", userData);
    return response.data.data;
  }

  /**
   * Update user
   */
  async updateUser(userId: string, userData: Partial<User>): Promise<User> {
    const response = await axiosInstance.put(`/users/${userId}`, userData);
    return response.data.data;
  }

  /**
   * Delete user
   */
  async deleteUser(userId: string): Promise<void> {
    await axiosInstance.delete(`/users/${userId}`);
  }
}

export const userService = new UserService();
