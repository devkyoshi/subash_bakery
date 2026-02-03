import axiosInstance from "@/lib/axios";
import {
  LoginRequest,
  LoginResponse,
  RefreshTokenRequest,
  RefreshTokenResponse,
  RegisterRequest,
  UpdateProfileRequest,
  ChangePasswordRequest,
  User,
  Role,
} from "@/types/auth.types";
import { ApiResponse } from "@/types/global.types";

class AuthService {
  /**
   * Login user with email and password
   */
  async login(credentials: LoginRequest): Promise<LoginResponse> {
    const response = await axiosInstance.post<ApiResponse<LoginResponse>>(
      "/auth/login",
      credentials,
    );

    const authData = response.data.data;

    // Store tokens in localStorage
    if (authData.access_token) {
      localStorage.setItem("accessToken", authData.access_token);
    }
    if (authData.refresh_token) {
      localStorage.setItem("refreshToken", authData.refresh_token);
    }

    // Store user data
    if (authData.user) {
      localStorage.setItem("user", JSON.stringify(authData.user));

      // Store organization ID if available
      if (authData.user.organization_id) {
        localStorage.setItem("organizationId", authData.user.organization_id);
      }
    }

    // Store role data
    if (authData.role) {
      localStorage.setItem("role", JSON.stringify(authData.role));
    }

    return authData;
  }

  /**
   * Register new user
   */
  async register(data: RegisterRequest): Promise<LoginResponse> {
    const response = await axiosInstance.post<ApiResponse<LoginResponse>>(
      "/auth/register",
      data,
    );

    const authData = response.data.data;

    // Store tokens and user data
    if (authData.access_token) {
      localStorage.setItem("accessToken", authData.access_token);
    }
    if (authData.refresh_token) {
      localStorage.setItem("refreshToken", authData.refresh_token);
    }
    if (authData.user) {
      localStorage.setItem("user", JSON.stringify(authData.user));
      if (authData.user.organization_id) {
        localStorage.setItem("organizationId", authData.user.organization_id);
      }
    }
    if (authData.role) {
      localStorage.setItem("role", JSON.stringify(authData.role));
    }

    return authData;
  }

  /**
   * Admin creating a new user (does not store tokens)
   */
  async adminCreateUser(data: RegisterRequest): Promise<LoginResponse> {
    const response = await axiosInstance.post<ApiResponse<LoginResponse>>(
      "/auth/register",
      data,
    );
    return response.data.data;
  }

  /**
   * Logout user and clear stored data
   */
  async logout(): Promise<void> {
    try {
      const refreshToken = this.getRefreshToken();
      if (refreshToken) {
        // Call logout endpoint to invalidate tokens on server
        await axiosInstance.post("/auth/logout", {
          refresh_token: refreshToken,
        });
      }
    } catch (error) {
      console.error("Logout error:", error);
    } finally {
      // Clear all stored data
      localStorage.removeItem("accessToken");
      localStorage.removeItem("refreshToken");
      localStorage.removeItem("user");
      localStorage.removeItem("role");
      localStorage.removeItem("organizationId");
    }
  }

  /**
   * Refresh access token using refresh token
   */
  async refreshToken(
    refreshTokenRequest: RefreshTokenRequest,
  ): Promise<RefreshTokenResponse> {
    const response = await axiosInstance.post<
      ApiResponse<RefreshTokenResponse>
    >("/auth/refresh", refreshTokenRequest);

    const authData = response.data.data;

    // Update stored tokens
    if (authData.access_token) {
      localStorage.setItem("accessToken", authData.access_token);
    }
    if (authData.refresh_token) {
      localStorage.setItem("refreshToken", authData.refresh_token);
    }

    return authData;
  }

  /**
   * Get current user profile with role
   */
  async getCurrentUser(): Promise<{ user: User; role?: Role }> {
    const response =
      await axiosInstance.get<ApiResponse<{ user: User; role?: Role }>>(
        "/auth/me",
      );

    const data = response.data.data;

    // Update stored user data
    if (data.user) {
      localStorage.setItem("user", JSON.stringify(data.user));
    }

    // Update stored role data
    if (data.role) {
      localStorage.setItem("role", JSON.stringify(data.role));
    }

    return data;
  }

  /**
   * Get stored user from localStorage
   */
  getStoredUser(): User | null {
    const userStr = localStorage.getItem("user");
    if (!userStr) return null;

    try {
      return JSON.parse(userStr) as User;
    } catch (error) {
      console.error("Error parsing stored user:", error);
      return null;
    }
  }

  /**
   * Get stored role from localStorage
   */
  getStoredRole(): Role | null {
    const roleStr = localStorage.getItem("role");
    if (!roleStr) return null;

    try {
      return JSON.parse(roleStr) as Role;
    } catch (error) {
      console.error("Error parsing stored role:", error);
      return null;
    }
  }

  /**
   * Update user profile
   */
  async updateProfile(data: UpdateProfileRequest): Promise<User> {
    const response = await axiosInstance.put<ApiResponse<User>>(
      "/auth/me",
      data,
    );
    const updatedUser = response.data.data;

    // Update stored user if it's the current user
    if (updatedUser) {
      localStorage.setItem("user", JSON.stringify(updatedUser));
    }

    return updatedUser;
  }

  /**
   * Change user password
   */
  async changePassword(data: ChangePasswordRequest): Promise<void> {
    await axiosInstance.post("/auth/change-password", data);
  }

  /**
   * Get stored access token
   */
  getAccessToken(): string | null {
    return localStorage.getItem("accessToken");
  }

  /**
   * Get stored refresh token
   */
  getRefreshToken(): string | null {
    return localStorage.getItem("refreshToken");
  }

  /**
   * Check if user is authenticated
   */
  isAuthenticated(): boolean {
    return !!this.getAccessToken();
  }
}

export const authService = new AuthService();
