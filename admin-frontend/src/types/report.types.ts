// PO vs GRN Comparison Types

export interface POvsGRNComparisonItem {
  po_id: string;
  po_number: string;
  order_date: string;
  supplier_id: string;
  supplier_name: string;
  product_id: string;
  sku: string;
  product_name: string;
  po_qty: number;
  grn_qty: number;
  accepted_qty: number;
  rejected_qty: number;
  variance: number;
  variance_pct: number;
  unit_price: number;
  po_value: number;
  grn_value: number;
  value_variance: number;
  status: "MATCHED" | "PARTIAL" | "EXCESS" | "PENDING";
}

export interface POvsGRNMetrics {
  total_pos: number;
  completed_pos: number;
  partial_pos: number;
  pending_pos: number;
  excess_pos: number;
  total_variance: number;
  total_po_value: number;
  total_grn_value: number;
  variance_percent: number;
  completed_percent: number;
}

export interface VarianceDistribution {
  name: string;
  value: number;
  color: string;
}

export interface ActionItem {
  id: string;
  type: "critical" | "warning" | "info";
  title: string;
  description: string;
  po_id: string;
  po_number: string;
}

export interface POvsGRNReportResponse {
  metrics: POvsGRNMetrics;
  variance_distribution: VarianceDistribution[];
  items: POvsGRNComparisonItem[];
  action_items: ActionItem[];
  total_items: number;
}

export interface ReportFilters {
  start_date?: string;
  end_date?: string;
  supplier_id?: string;
  status?: string;
  location_id?: string;
}

// ============================================================
// Stock Level Comparison Types
// ============================================================

export interface StockLevelComparisonItem {
  product_id: string;
  sku: string;
  product_name: string;
  category_id: string;
  category_name: string;
  location_id: string;
  location_name: string;
  unit: string;
  system_qty: number;
  available_qty: number;
  allocated_qty: number;
  in_transit_qty: number;
  reorder_level: number;
  min_stock: number;
  max_stock: number;
  average_cost: number;
  total_value: number;
  stock_status: "OPTIMAL" | "LOW" | "CRITICAL" | "OVERSTOCK" | "OUT_OF_STOCK";
}

export interface StockLevelMetrics {
  total_products: number;
  optimal_count: number;
  low_stock_count: number;
  critical_count: number;
  overstock_count: number;
  out_of_stock_count: number;
  total_stock_value: number;
  total_on_hand: number;
  total_allocated: number;
  total_available: number;
}

export interface StockStatusDistribution {
  name: string;
  value: number;
  color: string;
}

export interface StockLevelReportResponse {
  metrics: StockLevelMetrics;
  status_distribution: StockStatusDistribution[];
  items: StockLevelComparisonItem[];
  total_items: number;
}

export interface StockLevelFilters {
  category_id?: string;
  location_id?: string;
  stock_status?: string;
  search?: string;
}

// ============================================================
// Reorder Status Report Types
// ============================================================

export interface ReorderItem {
  id: string;
  name: string;
  unit: string;
  priority: "CRITICAL" | "WARNING" | "NORMAL";
  currentStock: number;
  minLevel: number;
  remainingDays: number;
  pending: string;
  sugQty: number;
  leadTime: string;
}

export interface ConsumptionRow {
  category: string;
  avgDaily: string;
  trend: string;
  trendDir: "up" | "down" | "neutral";
  forecast: string;
}

export interface ReorderMetrics {
  critical_count: number;
  warning_count: number;
  normal_count: number;
}

export interface ReorderStatusReportResponse {
  metrics: ReorderMetrics;
  items: ReorderItem[];
  consumption_data: ConsumptionRow[];
  total_items: number;
}

export interface ReorderStatusFilters {
  category_id?: string;
  location_id?: string;
  priority?: string;
  search?: string;
  include_pending?: boolean;
}
