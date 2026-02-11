export enum SupplierStatus {
  Active = "active",
  Inactive = "inactive",
  Blacklist = "blacklist",
}

export enum POStatus {
  Draft = "draft",
  Sent = "sent",
  Confirmed = "confirmed",
  Received = "received",
  PartiallyReceived = "partial",
  Cancelled = "cancelled",
}

export enum GRNStatus {
  Draft = "draft",
  Received = "received",
  Inspected = "inspected",
  Accepted = "accepted",
  Rejected = "rejected",
}

export interface Address {
  street: string;
  city: string;
  state: string;
  country: string;
  postal_code: string;
}

export interface Supplier {
  id: string;
  organization_id: string;
  supplier_code: string;
  company_name: string;
  status: SupplierStatus;
  contact_person: string;
  email: string;
  phone: string;
  mobile?: string;
  website?: string;
  address?: Address;
  tax_id?: string;
  payment_terms?: number;
  credit_limit?: number;
  currency?: string;
  bank_name?: string;
  account_number?: string;
  account_name?: string;
  swift_code?: string;
  total_orders: number;
  total_purchase_value: number;
  outstanding_balance: number;
  last_order_date?: string;
  rating?: number;
  lead_time?: number;
  notes?: string;
  tags?: string[];
  created_at: string;
  updated_at: string;
}

export interface CreateSupplierRequest {
  company_name: string;
  contact_person: string;
  email: string;
  phone?: string;
  mobile?: string;
  website?: string;
  address?: Address;
  tax_id?: string;
  payment_terms?: number;
  credit_limit?: number;
  bank_name?: string;
  account_number?: string;
  swift_code?: string;
  tags?: string[];
  notes?: string;
}

export interface PurchaseOrderItem {
  product_id: string;
  sku?: string;
  description?: string;
  quantity: number;
  quantity_received?: number;
  unit_price: number;
  tax_rate?: number;
  discount_percent?: number;
  line_total: number;
}

export interface PurchaseOrder {
  id: string;
  organization_id: string;
  po_number: string;
  status: POStatus;
  supplier_id: string;
  delivery_location_id?: string;
  delivery_address?: Address;
  items: PurchaseOrderItem[];
  order_date: string;
  expected_date?: string;
  received_date?: string;
  subtotal: number;
  tax_amount: number;
  discount_amount: number;
  shipping_cost: number;
  total_amount: number;
  currency: string;
  payment_terms?: number;
  payment_method?: string;
  terms?: string;
  reference_number?: string;
  requested_by?: string;
  approved_by?: string;
  approved_date?: string;
  notes?: string;
  tags?: string[];
  created_at: string;
  updated_at: string;
  supplier_name?: string;
}

export interface CreatePurchaseOrderRequest {
  supplier_id: string;
  order_date: string;
  expected_date?: string;
  items: PurchaseOrderItem[];
  shipping_address?: Address;
  notes?: string;
  terms?: string;
  reference_number?: string;
  tax_rate?: number;
}

export interface GRNItem {
  product_id: string;
  product_name?: string;
  sku?: string;
  description?: string;
  ordered_quantity?: number;
  received_quantity: number;
  accepted_quantity?: number;
  rejected_quantity?: number;
  unit_cost?: number;
  batch_number?: string;
  expiry_date?: string;
  condition?: string;
  rejection_reason?: string;
}

export interface GoodsReceiptNote {
  id: string;
  organization_id: string;
  grn_number: string;
  status: GRNStatus;
  purchase_order_id: string;
  po_number?: string;
  supplier_id?: string;
  location_id: string;
  receipt_date: string;
  received_by: string;
  received_by_name?: string;
  items: GRNItem[];
  inspected_by?: string;
  inspected_by_name?: string;
  inspected_date?: string;
  qc_status?: string;
  qc_notes?: string;
  invoice_number?: string;
  delivery_note?: string;
  notes?: string;
  created_at: string;
  updated_at: string;
  // Enriched fields
  supplier_name?: string;
  location_name?: string;
}

export interface GoodsReceiptNoteSummary {
  id: string;
  grn_number: string;
  status: GRNStatus;
  purchase_order_id: string;
  po_number?: string;
  supplier: { id: string; name: string };
  location_id: string;
  receipt_date: string;
  received_by: { id: string; name: string };
  inspected_by?: { id: string; name: string };
  inspected_date?: string;
  qc_status?: string;
  invoice_number?: string;
  delivery_note?: string;
  total_value: number;
  items: GRNItem[];
  po_unit?: string;
  ordered_unit?: string;
  received_unit?: string;
}

export interface CreateGRNRequest {
  purchase_order_id: string;
  location_id: string;
  receipt_date: string;
  items: GRNItem[];
  notes?: string;
}
