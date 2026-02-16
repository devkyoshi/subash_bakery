import { Navigate, Route, Routes } from "react-router-dom";
import NotFound from "@/pages/NotFound";
import { AuthLayout } from "@/components/layout/AuthLayout";
import { DashboardLayout } from "@/components/layout/DashboardLayout";
import { ProtectedRoute } from "@/components/auth/ProtectedRoute";

import { WelcomePage } from "@/pages/auth/Welcome";
import { LoginPage } from "@/pages/auth/Login";
import { RegisterPage } from "@/pages/auth/Register";

import { DashboardPage } from "@/pages/dashboard/Dashboard";
import { CalendarPage } from "@/pages/dashboard/Calendar";
import { OrdersPage } from "@/pages/dashboard/Orders";
import { ReportsPage } from "@/pages/dashboard/Reports";

import POvsGRNPage from "@/pages/reports/po-vs-grn";
import ReorderStatusPage from "@/pages/reports/reorder-status";
import { UsersPage } from "@/pages/users/Users";
import { TransactionsPage } from "@/pages/dashboard/Transactions";
import { AnalyticsPage } from "@/pages/dashboard/Analytics";
import { SettingsPage } from "@/pages/settings/Settings";
import { ProfileSettings } from "@/pages/settings/ProfileSettings";
import { UserRoleList } from "@/pages/users/UserRoleList";
import { UserRoleFormPage } from "@/pages/users/UserRoleFormPage";
import { ProductsPage } from "@/pages/products/Products";
import { ProductDetailsPage } from "@/pages/products/ProductDetails";
import { ProductFormPage } from "@/pages/products/ProductForm";
import { CategoriesPage } from "@/pages/categories/CategoryList";
import { CategoryFormPage } from "@/pages/categories/CategoryForm";
import { BrandsPage } from "@/pages/brands/BrandList";
import { BrandFormPage } from "@/pages/brands/BrandForm";
import CompanyList from "@/pages/settings/companies/CompanyList";
import CompanyDetail from "@/pages/settings/companies/CompanyDetail";
import UnitList from "@/pages/settings/units/UnitList";
import UnitChartList from "@/pages/settings/unit-charts/UnitChartList";
import { SuppliersPage } from "@/pages/procurement/suppliers/SuppliersPage";
import { PurchaseOrdersPage } from "@/pages/procurement/orders/PurchaseOrdersPage";
import { CreatePurchaseOrderPage } from "@/pages/procurement/orders/CreatePurchaseOrderPage";
import { CreateSupplierPage } from "@/pages/procurement/suppliers/CreateSupplierPage";
import { EditSupplierPage } from "@/pages/procurement/suppliers/EditSupplierPage";
import { SupplierDetailsPage } from "@/pages/procurement/suppliers/SupplierDetailsPage";
import { PurchaseOrderDetailsPage } from "@/pages/procurement/orders/PurchaseOrderDetailsPage";
import { EditPurchaseOrderPage } from "@/pages/procurement/orders/EditPurchaseOrderPage";
import { GRNListPage } from "@/pages/procurement/grn/GRNListPage";
import { CreateGRNPage } from "@/pages/procurement/grn/CreateGRNPage";
import { GRNDetailsPage } from "@/pages/procurement/grn/GRNDetailsPage";
import { InventoryPage as StockLevelsPage } from "@/pages/inventory/InventoryPage";
import { StockAdjustmentsPage } from "@/pages/inventory/StockAdjustmentsPage";
import { CreateStockAdjustmentPage } from "@/pages/inventory/CreateStockAdjustmentPage";
import { StockAdjustmentDetailsPage } from "@/pages/inventory/StockAdjustmentDetailsPage";
import { StockMovementsPage } from "@/pages/inventory/StockMovementsPage";
import StockLevelReportPage from "@/pages/reports/stock-levels";

export function AppRoutes() {
  return (
    <Routes>
      <Route path="/" element={<Navigate to="/app/dashboard" replace />} />

      <Route path="/auth" element={<AuthLayout />}>
        <Route index element={<WelcomePage />} />
        <Route path="login" element={<LoginPage />} />
        <Route path="register" element={<RegisterPage />} />
      </Route>

      <Route
        path="/app"
        element={
          <ProtectedRoute requiredRole="ADMIN">
            <DashboardLayout />
          </ProtectedRoute>
        }
      >
        <Route path="dashboard" element={<DashboardPage />} />
        <Route path="reports" element={<ReportsPage />} />
        <Route path="reports/po-vs-grn" element={<POvsGRNPage />} />
        <Route path="reports/stock-levels" element={<StockLevelReportPage />} />
        <Route path="reports/reorder-status" element={<ReorderStatusPage />} />
        <Route path="orders" element={<OrdersPage />} />
        <Route path="calendar" element={<CalendarPage />} />

        <Route path="users/all" element={<UsersPage />} />
        <Route path="users/roles" element={<UserRoleList />} />
        <Route path="users/roles/new" element={<UserRoleFormPage />} />
        <Route path="users/roles/:id/edit" element={<UserRoleFormPage />} />

        <Route path="products" element={<ProductsPage />} />
        <Route path="products/new" element={<ProductFormPage />} />
        <Route path="products/:id" element={<ProductDetailsPage />} />
        <Route path="products/:id/edit" element={<ProductFormPage />} />

        <Route path="categories" element={<CategoriesPage />} />
        <Route path="categories/new" element={<CategoryFormPage />} />
        <Route path="categories/:id/edit" element={<CategoryFormPage />} />

        <Route path="brands" element={<BrandsPage />} />
        <Route path="brands/new" element={<BrandFormPage />} />
        <Route path="brands/:id/edit" element={<BrandFormPage />} />

        <Route path="transactions" element={<TransactionsPage />} />
        <Route path="analytics" element={<AnalyticsPage />} />
        <Route path="settings" element={<SettingsPage />} />
        <Route path="settings/units" element={<UnitList />} />
        <Route path="settings/unit-charts" element={<UnitChartList />} />
        <Route path="companies" element={<CompanyList />} />
        <Route path="companies/:id" element={<CompanyDetail />} />

        <Route path="procurement/suppliers" element={<SuppliersPage />} />
        <Route
          path="procurement/suppliers/new"
          element={<CreateSupplierPage />}
        />
        <Route
          path="procurement/suppliers/:id/edit"
          element={<EditSupplierPage />}
        />
        <Route
          path="procurement/suppliers/:id"
          element={<SupplierDetailsPage />}
        />

        <Route path="procurement/orders" element={<PurchaseOrdersPage />} />
        <Route
          path="procurement/orders/new"
          element={<CreatePurchaseOrderPage />}
        />
        <Route
          path="procurement/orders/:id/edit"
          element={<EditPurchaseOrderPage />}
        />
        <Route
          path="procurement/orders/:id"
          element={<PurchaseOrderDetailsPage />}
        />

        <Route path="procurement/grn" element={<GRNListPage />} />
        <Route path="procurement/grn/new" element={<CreateGRNPage />} />
        <Route path="procurement/grn/:id" element={<GRNDetailsPage />} />

        <Route path="inventory/stock-levels" element={<StockLevelsPage />} />
        <Route path="inventory/movements" element={<StockMovementsPage />} />

        <Route
          path="inventory/adjustments"
          element={<StockAdjustmentsPage />}
        />
        <Route
          path="inventory/adjustments/new"
          element={<CreateStockAdjustmentPage />}
        />
        <Route
          path="inventory/adjustments/:id"
          element={<StockAdjustmentDetailsPage />}
        />

        <Route path="profile" element={<ProfileSettings />} />
        <Route index element={<Navigate to="/app/dashboard" replace />} />
      </Route>

      {/* ADD ALL CUSTOM ROUTES ABOVE THE CATCH-ALL "*" ROUTE */}
      <Route path="*" element={<NotFound />} />
    </Routes>
  );
}
