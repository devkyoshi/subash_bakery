import * as z from "zod";

// Define the schema based on CreateProductRequest
export const productSchema = z.object({
  name: z.string().min(2, "Name is required"),
  sku: z.string().min(2, "SKU is required"),
  barcode: z.string().optional(),
  description: z.string().optional(),
  type: z.string().default("finished_goods"),
  status: z.string().default("active"),

  category_id: z.string().optional(),
  subcategory_id: z.string().optional(),
  brand_id: z.string().optional(),

  // Inventory
  track_inventory: z.boolean().default(true),
  track_batches: z.boolean().default(false),
  track_serial_numbers: z.boolean().default(false),
  reorder_level: z.coerce.number().min(0).default(0),

  // Dimensions & Weight
  weight: z.coerce.number().min(0).default(0),
  weight_unit: z.string().default("kg"),
  length: z.coerce.number().min(0).default(0),
  width: z.coerce.number().min(0).default(0),
  height: z.coerce.number().min(0).default(0),
  dimension_unit: z.string().default("cm"),
  volume: z.coerce.number().min(0).default(0),
  volume_unit: z.string().default("m3"),
  tags: z.array(z.string()).optional(),

  // Media
  images: z.array(z.string()).optional(),

  // Pricing
  base_unit_id: z.string().optional(),
  location_prices: z
    .array(
      z.object({
        location_id: z.string(),
        location_name: z.string().optional(),

        // Purchase Details
        cost_price: z.coerce.number().min(0).default(0),
        purchase_unit_id: z.string().optional(),

        // Selling Details
        selling_price: z.coerce.number().min(0).default(0),
        selling_unit_id: z.string().optional(),

        mrp: z.coerce.number().min(0).default(0),
        initial_stock: z.coerce.number().min(0).default(0),
        currency: z.string().default("LKR"),
        is_active: z.boolean().default(true),
        // Units for this location price might be needed if they differ
      }),
    )
    .default([]),
});

export type ProductFormValues = z.infer<typeof productSchema>;
