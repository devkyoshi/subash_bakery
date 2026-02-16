import {
  BarChart3,
  Boxes,
  CalendarDays,
  CreditCard,
  FileText,
  LayoutDashboard,
  Settings,
  Users,
  Users2,
  ShoppingBag,
  Package,
  Layers,
  Tag,
  Building2,
  Ruler,
  ArrowRightLeft,
  PackageCheck,
  ShieldCheck,
  Monitor,
} from "lucide-react";

export type SidebarNavItem = {
  label: string;
  to: string;
  icon: React.ComponentType<{ className?: string }>;
  badge?: string;
  submenu?: SidebarNavItem[];
};

export const SIDEBAR_NAV: SidebarNavItem[] = [
  { label: "Dashboard", to: "/app/dashboard", icon: LayoutDashboard },
  {
    label: "Users",
    to: "/app/users",
    icon: Users,
    submenu: [
      { label: "All Users", to: "/app/users/all", icon: Users },
      { label: "Roles", to: "/app/users/roles", icon: ShieldCheck },
    ],
  },
  { label: "Products", to: "/app/products", icon: Package },
  { label: "Categories", to: "/app/categories", icon: Layers },
  { label: "Brands", to: "/app/brands", icon: Tag },
  { label: "Companies", to: "/app/companies", icon: Building2 },
  { label: "Devices", to: "/app/devices", icon: Monitor },
  {
    label: "Procurement",
    to: "/app/procurement",
    icon: ShoppingBag,
    submenu: [
      { label: "Suppliers", to: "/app/procurement/suppliers", icon: Users2 },
      { label: "Orders", to: "/app/procurement/orders", icon: FileText },
      {
        label: "Goods Receipt",
        to: "/app/procurement/grn",
        icon: PackageCheck,
      },
    ],
  },
  {
    label: "Inventory",
    to: "/app/inventory",
    icon: Boxes,
    submenu: [
      { label: "Stock Levels", to: "/app/inventory/stock-levels", icon: Boxes },
      {
        label: "Adjustments",
        to: "/app/inventory/adjustments",
        icon: FileText,
      },
      {
        label: "Movements",
        to: "/app/inventory/movements",
        icon: ArrowRightLeft,
      },
    ],
  },
  {
    label: "Reports",
    to: "/app/reports",
    icon: FileText,
    submenu: [
      { label: "PO vs GRN", to: "/app/reports/po-vs-grn", icon: Boxes },
      {
        label: "Reorder Status",
        to: "/app/reports/reorder-status",
        icon: Boxes,
      },
      { label: "Stock Levels", to: "/app/reports/stock-levels", icon: Boxes },
    ],
  },
  {
    label: "Settings",
    to: "/app/settings",
    icon: Settings,
    submenu: [
      { label: "Units", to: "/app/settings/units", icon: Ruler },
      {
        label: "Unit Conversions",
        to: "/app/settings/unit-charts",
        icon: ArrowRightLeft,
      },
    ],
  },
];
