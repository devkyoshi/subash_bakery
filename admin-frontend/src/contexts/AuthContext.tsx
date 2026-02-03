import React, {
  createContext,
  useContext,
  useState,
  useEffect,
  ReactNode,
} from "react";
import { authService } from "@/services/auth.service";
import { AuthContextType, AuthState, User } from "@/types/auth.types";

const initialState: AuthState = {
  user: null,
  role: null,
  accessToken: null,
  refreshToken: null,
  isAuthenticated: false,
  isLoading: true,
  error: null,
};

const AuthContext = createContext<AuthContextType | undefined>(undefined);

interface AuthProviderProps {
  children: ReactNode;
}

export function AuthProvider({ children }: AuthProviderProps) {
  const [state, setState] = useState<AuthState>(initialState);

  // Initialize auth state from localStorage
  useEffect(() => {
    const initializeAuth = async () => {
      try {
        const storedUser = authService.getStoredUser();
        const storedRole = authService.getStoredRole();
        const accessToken = authService.getAccessToken();
        const refreshToken = authService.getRefreshToken();

        if (storedUser && accessToken) {
          // Try to fetch fresh user data
          try {
            const data = await authService.getCurrentUser();
            setState({
              user: data.user,
              role: data.role || null,
              accessToken,
              refreshToken,
              isAuthenticated: true,
              isLoading: false,
              error: null,
            });
          } catch (error) {
            // If fetching user fails, use stored user and role
            setState({
              user: storedUser,
              role: storedRole,
              accessToken,
              refreshToken,
              isAuthenticated: true,
              isLoading: false,
              error: null,
            });
          }
        } else {
          setState({
            ...initialState,
            isLoading: false,
          });
        }
      } catch (error) {
        console.error("Auth initialization error:", error);
        setState({
          ...initialState,
          isLoading: false,
          error: "Failed to initialize authentication",
        });
      }
    };

    initializeAuth();
  }, []);

  const login = async (email: string, password: string) => {
    try {
      setState((prev) => ({ ...prev, isLoading: true, error: null }));

      const response = await authService.login({ email, password });

      setState({
        user: response.user,
        role: response.role || null,
        accessToken: response.access_token,
        refreshToken: response.refresh_token,
        isAuthenticated: true,
        isLoading: false,
        error: null,
      });
    } catch (error: any) {
      const errorMessage =
        error.response?.data?.message || "Login failed. Please try again.";
      setState({
        ...initialState,
        isLoading: false,
        error: errorMessage,
      });
      throw error;
    }
  };

  const logout = async () => {
    try {
      await authService.logout();
    } catch (error) {
      console.error("Logout error:", error);
    } finally {
      setState({
        ...initialState,
        isLoading: false,
      });
    }
  };

  const refreshAccessToken = async () => {
    try {
      const refreshToken = authService.getRefreshToken();
      if (!refreshToken) {
        throw new Error("No refresh token available");
      }

      const response = await authService.refreshToken({
        refresh_token: refreshToken,
      });

      setState((prev) => ({
        ...prev,
        user: response.user,
        role: response.role || prev.role,
        accessToken: response.access_token,
        refreshToken: response.refresh_token,
      }));
    } catch (error) {
      console.error("Token refresh error:", error);
      // If refresh fails, logout user
      await logout();
      throw error;
    }
  };

  const hasRole = (roleName: string): boolean => {
    return state.role?.name === roleName;
  };

  const hasPermission = (permission: string): boolean => {
    return state.role?.permissions?.includes(permission) || false;
  };

  const value: AuthContextType = {
    ...state,
    login,
    logout,
    refreshAccessToken,
    hasRole,
    hasPermission,
  };

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}

export function useAuth() {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error("useAuth must be used within an AuthProvider");
  }
  return context;
}
