// Product related types and interfaces

export enum ProductType {
  RAW_MATERIAL = "raw_material",
  FINISHED_GOODS = "finished_goods",
  SEMI_FINISHED = "semi_finished",
  CONSUMABLE = "consumable",
  SERVICE = "service",
}

export enum ProductStatus {
  ACTIVE = "active",
  INACTIVE = "inactive",
  DISCONTINUED = "discontinued",
}

export enum ValuationMethod {
  FIFO = "fifo",
  LIFO = "lifo",
  WEIGHTED_AVERAGE = "weighted_average",
  STANDARD = "standard",
}

export interface LocationPrice {
  location_id: string;
  location_name?: string;
  purchase_unit_id?: string;
  purchase_unit?: any;
  selling_unit_id?: string;
  selling_unit?: any;
  cost_price: number;
  selling_price: number;
  mrp: number;
  initial_stock?: number;
  currency: string;
  is_active: boolean;
  created_at?: number;
  modified_at?: number;
  // Populated fields
  current_stock?: number;
  available_stock?: number;
  allocated_stock?: number;
}

export interface Product {
  id: string;
  organization_id: string;

  // Identifications
  sku: string;
  barcode?: string;
  name: string;
  description?: string;

  // Type and Status
  type: ProductType;
  status: ProductStatus;

  // Classification
  category_id?: string;
  category?: { id: string; name: string; code: string }; // Populated
  subcategory_id?: string;
  subcategory?: { id: string; name: string; code: string }; // Populated
  brand_id?: string;
  brand?: { id: string; name: string; code: string }; // Populated
  manufacturer_id?: string;

  // Inventory Settings
  track_inventory: boolean;
  track_batches: boolean;
  track_serial_numbers: boolean;
  valuation_method: ValuationMethod;

  // Units
  base_unit_id?: string;
  allowed_unit_ids?: string[];

  // Dimensions & Weight
  weight?: number;
  weight_unit?: string;
  length?: number;
  width?: number;
  height?: number;
  dimension_unit?: string;
  volume?: number;
  volume_unit?: string;

  // Pricing - Location-wise
  location_prices: LocationPrice[];

  // Tax & Accounting
  tax_category_id?: string;
  hsn_code?: string;
  sac_code?: string;

  // Reorder Settings
  reorder_level: number;
  reorder_quantity: number;
  min_stock_level: number;
  max_stock_level: number;
  safety_stock: number;

  // Supplier Info
  default_supplier_id?: string;
  supplier_ids?: string[];
  lead_time_days: number;

  // Quality & Expiry
  shelf_life_days?: number;
  requires_qc: boolean;
  perishable: boolean;
  hazardous: boolean;

  // Images & Attachments
  images?: string[];
  thumbnail?: string;
  specifications?: Record<string, string>;
  attachments?: Attachment[];

  // Current Stock Summary
  total_stock: number;
  available_stock: number;
  allocated_stock: number;
  in_transit_stock: number;
  stock_value: number;

  // Analytics
  total_sold?: number;
  total_purchased?: number;
  last_sold_date?: string;
  last_purchase_date?: string;

  // Metadata
  metadata?: Record<string, any>;
  tags?: string[];

  // Timestamps
  created_at: string;
  updated_at: string;
  created_by?: string;
  updated_by?: string;
  is_deleted?: boolean;
}

export interface Attachment {
  name: string;
  url: string;
  type: string;
  size: number;
}

export interface CreateProductRequest {
  sku: string;
  barcode?: string;
  name: string;
  description?: string;
  type: ProductType;
  status: ProductStatus;
  category_id?: string;
  subcategory_id?: string;
  brand_id?: string;
  track_inventory?: boolean;
  track_batches?: boolean;
  track_serial_numbers?: boolean;
  valuation_method?: ValuationMethod;
  location_prices: LocationPrice[];
  reorder_level?: number;
  reorder_quantity?: number;
  min_stock_level?: number;
  max_stock_level?: number;
  safety_stock?: number;
  lead_time_days?: number;
  hsn_code?: string;
  sac_code?: string;
  requires_qc?: boolean;
  perishable?: boolean;
  hazardous?: boolean;
  shelf_life_days?: number;
  images?: string[];
  tags?: string[];
}

export interface UpdateProductRequest {
  sku?: string;
  barcode?: string;
  name?: string;
  description?: string;
  type?: ProductType;
  status?: ProductStatus;
  category_id?: string;
  subcategory_id?: string;
  brand_id?: string;
  track_inventory?: boolean;
  track_batches?: boolean;
  track_serial_numbers?: boolean;
  valuation_method?: ValuationMethod;
  location_prices?: LocationPrice[];
  reorder_level?: number;
  reorder_quantity?: number;
  min_stock_level?: number;
  max_stock_level?: number;
  safety_stock?: number;
  lead_time_days?: number;
  hsn_code?: string;
  sac_code?: string;
  requires_qc?: boolean;
  perishable?: boolean;
  hazardous?: boolean;
  shelf_life_days?: number;
  images?: string[];
  tags?: string[];
}

export interface ProductFilters {
  organization_id: string;
  category_id?: string;
  subcategory_id?: string;
  brand_id?: string;
  status?: ProductStatus;
  type?: ProductType;
  track_inventory?: boolean;
  location_id?: string;
  search?: string;
  page?: number;
  limit?: number;
}

// ... (previous content)

export interface ProductListResponse {
  data: Product[];
  page: number;
  limit: number;
  total: number;
}

export interface Unit {
  id: string;
  code: string;
  name: string;
  symbol: string;
  unit_type: string; // Changed from type
  is_base_unit: boolean; // Changed from base_unit
  is_active: boolean;
}

export interface Location {
  id: string;
  organization_id: string;
  company_id: string;
  name: string;
  code: string;
  type: string;
  address?: {
    street: string;
    city: string;
    state: string;
    country: string;
    postal_code: string;
  };
  is_active: boolean;
}

export interface LocationPriceRequest {
  location_id: string;
  location_name?: string;
  purchase_unit_id?: string;
  selling_unit_id?: string;
  cost_price: number;
  selling_price: number;
  mrp: number;
  initial_stock?: number;
  currency: string;
  is_active: boolean;
}
