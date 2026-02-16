import axiosInstance from "@/lib/axios";
import { ApiResponse } from "@/types/global.types";
import { Role } from "@/types/role.types";

class RoleService {
  async getRoles(orgId?: string): Promise<Role[]> {
    const params = orgId ? { organization_id: orgId } : {};
    const response = await axiosInstance.get<any>("/roles", {
      params,
    });

    // Handle paginated response structure { data: { data: Role[], pagination: ... } }
    if (response.data.data && Array.isArray(response.data.data.data)) {
      return response.data.data.data;
    }

    // Handle non-paginated structure { data: Role[] }
    if (Array.isArray(response.data.data)) {
      return response.data.data;
    }

    return [];
  }

  async getRoleById(roleId: string): Promise<Role> {
    const response = await axiosInstance.get(`/roles/${roleId}`);
    return response.data;
  }

  async createRole(role: Partial<Role>): Promise<Role> {
    const response = await axiosInstance.post("/roles", role);
    return response.data;
  }

  async updateRole(roleId: string, role: Partial<Role>): Promise<Role> {
    const response = await axiosInstance.put(`/roles/${roleId}`, role);
    return response.data;
  }

  async deleteRole(roleId: string): Promise<void> {
    await axiosInstance.delete(`/roles/${roleId}`);
  }

  async assignRole(userId: string, roleId: string): Promise<void> {
    await axiosInstance.post("/roles/assign", {
      user_id: userId,
      role_id: roleId,
    });
  }
}

export const roleService = new RoleService();
