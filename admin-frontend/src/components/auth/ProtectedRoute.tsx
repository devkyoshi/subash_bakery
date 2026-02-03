import { Navigate, useLocation } from "react-router-dom";
import { useAuth } from "@/contexts/AuthContext";
import { AlertCircle } from "lucide-react";
import { Button } from "@/components/ui/button";

interface ProtectedRouteProps {
  children: React.ReactNode;
  requiredRole?: string;
  requiredPermissions?: string[];
}

export function ProtectedRoute({
  children,
  requiredRole,
  requiredPermissions,
}: ProtectedRouteProps) {
  const { isAuthenticated, isLoading, hasRole, hasPermission, logout } =
    useAuth();
  const location = useLocation();

  if (isLoading) {
    return (
      <div className="flex h-screen items-center justify-center">
        <div className="text-center">
          <div className="h-8 w-8 animate-spin rounded-full border-4 border-brand border-t-transparent"></div>
          <p className="mt-4 text-sm text-muted-foreground">Loading...</p>
        </div>
      </div>
    );
  }

  if (!isAuthenticated) {
    // Redirect to login page with return url
    return <Navigate to="/auth/login" state={{ from: location }} replace />;
  }

  // Check role if required
  if (requiredRole && !hasRole(requiredRole)) {
    return (
      <div className="flex h-screen items-center justify-center bg-background">
        <div className="text-center max-w-md px-6">
          <div className="mb-6 flex justify-center">
            <div className="rounded-full bg-destructive/10 p-4">
              <AlertCircle className="h-12 w-12 text-destructive" />
            </div>
          </div>
          <h1 className="text-3xl font-bold mb-3">Access Denied</h1>
          <p className="text-muted-foreground mb-6">
            You don't have permission to access this page. This application
            requires admin privileges.
          </p>
          <div className="flex gap-3 justify-center">
            <Button onClick={() => logout()} variant="outline">
              Logout
            </Button>
            <Button onClick={() => window.history.back()} variant="default">
              Go Back
            </Button>
          </div>
        </div>
      </div>
    );
  }

  // Check permissions if required
  if (requiredPermissions && requiredPermissions.length > 0) {
    const hasAllPermissions = requiredPermissions.every((permission) =>
      hasPermission(permission),
    );

    if (!hasAllPermissions) {
      return (
        <div className="flex h-screen items-center justify-center bg-background">
          <div className="text-center max-w-md px-6">
            <div className="mb-6 flex justify-center">
              <div className="rounded-full bg-destructive/10 p-4">
                <AlertCircle className="h-12 w-12 text-destructive" />
              </div>
            </div>
            <h1 className="text-3xl font-bold mb-3">Access Denied</h1>
            <p className="text-muted-foreground mb-6">
              You don't have the required permissions to access this page.
            </p>
            <div className="flex gap-3 justify-center">
              <Button onClick={() => logout()} variant="outline">
                Logout
              </Button>
              <Button onClick={() => window.history.back()} variant="default">
                Go Back
              </Button>
            </div>
          </div>
        </div>
      );
    }
  }

  return <>{children}</>;
}
