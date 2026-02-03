export interface StockLevel {
  id: string;
  created_at: string;
  updated_at: string;
  organization_id: string;
  product_id: string;
  location_id: string;
  warehouse_zone?: string;
  quantity_on_hand: number;
  quantity_available: number;
  quantity_allocated: number;
  quantity_in_transit: number;
  quantity_reserved: number;
  average_cost: number;
  last_cost: number;
  total_value: number;
  last_movement_date?: string;
  last_count_date?: string;
}

export enum MovementType {
  In = "in",
  Out = "out",
  Transfer = "transfer",
  Adjustment = "adjustment",
  Return = "return",
  Scrap = "scrap",
}

export interface StockMovement {
  id: string;
  created_at: string;
  updated_at: string;
  created_by: string;
  organization_id: string;
  product_id: string;
  movement_type: MovementType;
  from_location_id?: string;
  to_location_id?: string;
  quantity: number;
  uom: string;
  unit_cost: number;
  total_cost: number;
  reference_type?: string;
  reference_id?: string;
  reference_no?: string;
  batch_id?: string;
  serial_numbers?: string[];
  reason?: string;
  notes?: string;
  movement_date: string;
  is_reversed: boolean;
}

export enum AdjustmentStatus {
  Draft = "draft",
  Pending = "pending",
  Approved = "approved",
  Rejected = "rejected",
}

export interface StockAdjustmentItem {
  product_id: string;
  expected_qty: number;
  actual_qty: number;
  difference_qty: number;
  uom: string;
  unit_cost: number;
  total_cost: number;
  batch_id?: string;
  reason?: string;
}

export interface StockAdjustment {
  id: string;
  created_at: string;
  updated_at: string;
  created_by: string;
  organization_id: string;
  location_id: string;
  adjustment_no: string;
  adjustment_date: string;
  reason: string;
  reason_details?: string;
  items: StockAdjustmentItem[];
  status: AdjustmentStatus;
  approved_by?: string;
  approved_at?: string;
  rejected_reason?: string;
  notes?: string;
}

export interface CreateStockAdjustmentRequest {
  location_id: string;
  adjustment_no?: string;
  adjustment_date?: string;
  reason: string;
  reason_details?: string;
  items: CreateStockAdjustmentItemRequest[];
  notes?: string;
}

export interface CreateStockAdjustmentItemRequest {
  product_id: string;
  expected_qty: number;
  actual_qty: number;
  uom?: string;
  unit_cost?: number;
  batch_id?: string;
  reason?: string;
}

export interface CreateStockMovementRequest {
  product_id: string;
  movement_type: MovementType;
  from_location_id?: string;
  to_location_id?: string;
  quantity: number;
  uom?: string;
  unit_cost?: number;
  reference_type?: string;
  reference_id?: string;
  reference_no?: string;
  batch_id?: string;
  serial_numbers?: string[];
  reason?: string;
  notes?: string;
  movement_date?: string;
}

export interface Batch {
  id: string;
  created_at: string;
  updated_at: string;
  organization_id: string;
  product_id: string;
  location_id: string;
  batch_number: string;
  manufacture_date: string;
  expiry_date?: string;
  receive_date: string;
  initial_quantity: number;
  current_quantity: number;
  allocated_quantity: number;
  unit_cost: number;
  total_cost: number;
  supplier_id?: string;
  purchase_order_id?: string;
  qc_status?: string;
  qc_date?: string;
  qc_notes?: string;
  is_active: boolean;
  is_expired: boolean;
  is_depleted: boolean;
}

export interface InventoryCountItem {
  product_id: string;
  system_qty: number;
  counted_qty: number;
  variance_qty: number;
  uom: string;
  unit_cost: number;
  variance_value: number;
  batch_id?: string;
  counted_by: string;
  counted_at: string;
  notes?: string;
}

export interface InventoryCount {
  id: string;
  created_at: string;
  updated_at: string;
  organization_id: string;
  location_id: string;
  count_no: string;
  count_date: string;
  count_type: "full" | "cycle" | "spot";
  items: InventoryCountItem[];
  status: "in_progress" | "completed" | "cancelled";
  started_by: string;
  started_at: string;
  completed_by?: string;
  completed_at?: string;
  total_items_counted: number;
  total_variance: number;
  variance_value: number;
  notes?: string;
}
