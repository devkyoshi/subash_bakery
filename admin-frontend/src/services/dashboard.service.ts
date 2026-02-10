import api from "@/lib/axios";
import type {
  ActivityItem,
  Metric,
  OrderItem,
  ReportItem,
  SalesPoint,
} from "@/types/dashboard";

export type DashboardData = {
  metrics: Metric[];
  reports: ReportItem[];
  orders: OrderItem[];
  activities: ActivityItem[];
  sales: SalesPoint[];
};

interface ApiResponse<T> {
  success: boolean;
  data: T;
  message: string;
  timestamp: string;
}

interface InventoryStats {
  critical_stock_count: number;
  low_stock_items: any[]; // refined type would be better but keeping simple for now
}

interface ProcurementStats {
  pending_po_count: number;
  pending_grn_count: number;
  pending_approvals: any[];
}

export async function getDashboardData(): Promise<DashboardData> {
  try {
    const [inventoryRes, procurementRes] = await Promise.all([
      api.get<ApiResponse<InventoryStats>>("/inventory/dashboard/stats"),
      api.get<ApiResponse<ProcurementStats>>("/procurement/dashboard/stats"),
    ]);

    const inventoryData = inventoryRes.data.data;
    const procurementData = procurementRes.data.data;

    // Transform API data to Dashboard metrics
    const metrics: Metric[] = [
      {
        id: "critical_stock",
        label: "Critical Stock",
        value: inventoryData.critical_stock_count.toString(),
        deltaLabel: "Items",
        deltaVariant: "down", // high critical stock is bad
        tone: "warning",
      },
      {
        id: "pending_po",
        label: "Pending POs",
        value: procurementData.pending_po_count.toString(),
        deltaLabel: "Orders",
        deltaVariant: "up",
        tone: "info",
      },
      {
        id: "pending_grn",
        label: "Pending GRNs",
        value: procurementData.pending_grn_count.toString(),
        deltaLabel: "Receipts",
        deltaVariant: "up",
        tone: "brand",
      },
      {
        id: "pending_approvals",
        label: "Approvals Needed",
        value: procurementData.pending_approvals.length.toString(),
        deltaLabel: "Requests",
        deltaVariant: "up",
        tone: "warning",
      },
    ];

    // Transform Pending Approvals to "Orders" list for widget
    const orders: OrderItem[] = procurementData.pending_approvals.map(
      (po: any) => ({
        id: po.po_number,
        customer: "Unknown Supplier", // Supplier name enrichment happens on frontend or is missing
        amount: `$${po.total_amount.toFixed(2)}`,
        status: "Pending", // Mapped to Pending to match type
        flag: "orange",
      }),
    );

    // Fallback/Mock data for missing pieces (Reports, Sales, Activities)
    // As per plan, these services might not exist yet.
    const reports: ReportItem[] = [
      {
        id: "sales",
        title: "Sales Report",
        subtitle: "Unavailable",
        tone: "info",
      },
      {
        id: "inventory",
        title: "Inventory Report",
        subtitle: "Available",
        tone: "success",
      },
    ];

    const sales: SalesPoint[] = [
      { month: "Jan", invoiced: 0, cashed: 0, cashedPct: 0 },
      // ... minimal mock data
    ];

    const activities: ActivityItem[] = [
      {
        id: "a1",
        title: "System Ready",
        description: "Dashboard backend integrated",
        time: "Just now",
        tone: "success",
      },
    ];

    return {
      metrics,
      reports,
      orders,
      activities,
      sales,
    };
  } catch (error) {
    console.error("Failed to fetch dashboard data:", error);
    // Return empty/safe data on error
    return {
      metrics: [],
      reports: [],
      orders: [],
      activities: [],
      sales: [],
    };
  }
}
