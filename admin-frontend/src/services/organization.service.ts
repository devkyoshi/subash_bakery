import axiosInstance from "@/lib/axios";
import {
  ApiResponse,
  PaginatedResponse,
  PaginatedData,
} from "@/types/global.types";
import { Location } from "@/types/product.types";

export interface OrganizationOption {
  id: string;
  name: string;
}

export interface Company {
  id: string;
  organization_id: string;
  name: string; // Display Name
  legal_name: string;
  code?: string;
  tax_id?: string;
  currency?: string;
  email: string; // Flattened as per error
  phone?: string;
  website?: string;
  address?: {
    street: string;
    city: string;
    state: string;
    country: string;
    postal_code: string;
  };
  locations?: Location[];
  is_active: boolean;
  created_at?: string;
  updated_at?: string;
}

export interface CreateCompanyRequest {
  name: string; // Display Name (Required)
  legal_name: string; // Legal Name (Required)
  code: string; // Required
  tax_id?: string;
  currency?: string;
  email: string; // Required
  address: {
    // Required
    street: string;
    city: string;
    state: string;
    country: string;
    postal_code: string;
  };
  phone?: string;
  website?: string;
}

class OrganizationService {
  /**
   * Get all organizations (options only)
   */
  async getOrganizations(): Promise<OrganizationOption[]> {
    const response = await axiosInstance.get<ApiResponse<OrganizationOption[]>>(
      "/organizations/options",
    );
    return response.data.data || [];
  }

  /**
   * List companies for an organization
   */
  async getCompanies(
    orgId: string,
    params?: { page?: number; limit?: number; q?: string; is_active?: boolean },
  ): Promise<ApiResponse<PaginatedData<Company>>> {
    const response = await axiosInstance.get<
      ApiResponse<PaginatedData<Company>>
    >(`/organizations/${orgId}/companies`, { params });
    return response.data;
  }

  /**
   * Get a company by ID
   */
  async getCompany(id: string): Promise<Company> {
    const response = await axiosInstance.get<ApiResponse<Company>>(
      `/companies/${id}`,
    );
    return response.data.data;
  }

  /**
   * Create a new company
   */
  async createCompany(
    orgId: string,
    data: CreateCompanyRequest,
  ): Promise<Company> {
    const response = await axiosInstance.post<ApiResponse<Company>>(
      `/organizations/${orgId}/companies`,
      data,
    );
    return response.data.data;
  }

  /**
   * Update a company
   */
  async updateCompany(
    id: string,
    data: CreateCompanyRequest,
  ): Promise<Company> {
    const response = await axiosInstance.put<ApiResponse<Company>>(
      `/companies/${id}`,
      data,
    );
    return response.data.data;
  }
}

export const organizationService = new OrganizationService();
