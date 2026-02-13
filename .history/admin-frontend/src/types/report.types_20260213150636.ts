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
