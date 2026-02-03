// Authentication related types and interfaces

export interface Role {
  id: string;
  name: string;
  permissions: string[];
}

export interface User {
  id: string;
  organization_id?: string;
  email: string;
  first_name: string;
  last_name: string;
  phone?: string;
  avatar?: string;
  role_id?: string;
  is_active: boolean;
  is_email_verified: boolean;
  google_id?: string;
  last_login_at?: string;
  created_at: string;
  updated_at: string;
}

export interface LoginRequest {
  email: string;
  password: string;
}

export interface LoginResponse {
  access_token: string;
  refresh_token: string;
  user: User;
  role?: Role;
}

export interface RefreshTokenRequest {
  refresh_token: string;
}

export interface RefreshTokenResponse {
  access_token: string;
  refresh_token: string;
  user: User;
  role?: Role;
}

export interface RegisterRequest {
  email: string;
  password: string;
  first_name: string;
  last_name: string;
  phone?: string;
  organization_id?: string;
}

export interface AuthState {
  user: User | null;
  role: Role | null;
  accessToken: string | null;
  refreshToken: string | null;
  isAuthenticated: boolean;
  isLoading: boolean;
  error: string | null;
}

export interface AuthContextType extends AuthState {
  login: (email: string, password: string) => Promise<void>;
  logout: () => void;
  refreshAccessToken: () => Promise<void>;
  hasRole: (roleName: string) => boolean;
  hasPermission: (permission: string) => boolean;
}

export interface UpdateProfileRequest {
  first_name: string;
  last_name: string;
  phone?: string;
  avatar?: string;
}

export interface ChangePasswordRequest {
  current_password?: string; // Optional if admin resets, but usually required for user
  new_password: string;
}
