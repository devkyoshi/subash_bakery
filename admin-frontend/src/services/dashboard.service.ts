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

interface DashboardOverviewResponse {
  inventory: InventoryStats;
  procurement: ProcurementStats;
  activities: any[];
  errors?: string[];
}

export async function getDashboardData(): Promise<DashboardData> {
  try {
    const response = await api.get<ApiResponse<DashboardOverviewResponse>>(
      "/dashboard/overview",
    );
    const { inventory, procurement } = response.data.data;

    // Transform API data to Dashboard metrics
    const metrics: Metric[] = [
      {
        id: "critical_stock",
        label: "Critical Stock",
        value: inventory.critical_stock_count.toString(),
        deltaLabel: "Items",
        deltaVariant: "down", // high critical stock is bad
        tone: "warning",
      },
      {
        id: "pending_po",
        label: "Pending POs",
        value: procurement.pending_po_count.toString(),
        deltaLabel: "Orders",
        deltaVariant: "up",
        tone: "info",
      },
      {
        id: "pending_grn",
        label: "Pending GRNs",
        value: procurement.pending_grn_count.toString(),
        deltaLabel: "Receipts",
        deltaVariant: "up",
        tone: "brand",
      },
      {
        id: "pending_approvals",
        label: "Approvals Needed",
        value: procurement.pending_approvals.length.toString(),
        deltaLabel: "Requests",
        deltaVariant: "up",
        tone: "warning",
      },
    ];

    // Transform Pending Approvals to "Orders" list for widget
    const orders: OrderItem[] = procurement.pending_approvals.map(
      (po: any) => ({
        id: po.po_number,
        customer: po.supplier_name || "Unknown Supplier",
        amount: `$${po.total_amount.toFixed(2)}`,
        itemCount: po.items.length,
        date: new Date(po.order_date).toLocaleDateString(),
        mongoId: po.id,
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

    // Map Activities
    const activities: ActivityItem[] = (
      response.data.data.activities || []
    ).map((activity: any) => ({
      id: activity.id,
      title: activity.description || "Activity",
      description: `Action: ${activity.action} on ${activity.type.replace("_", " ")}`,
      time:
        new Date(activity.created_at).toLocaleTimeString() +
        ", " +
        new Date(activity.created_at).toLocaleDateString(),
      tone: "info",
    }));

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
