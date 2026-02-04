# Product Service API Documentation

Base URL: `/api/v1`

## Authentication

All endpoints require a valid JWT token in the `Authorization` header:
`Authorization: Bearer <token>`

## 1. Products

### Create Product

**Endpoint:** `POST /products`
**Description:** Create a new product.

**Request Body:**

```json
{
  "organization_id": "60d5ec49f1b2c8b1f8c8e1a0",
  "sku": "PROD-001",
  "barcode": "8901234567890",
  "name": "Widget X",
  "description": "High quality widget",
  "type": "finished_goods",
  "status": "active",
  "category_id": "60d5ec49f1b2c8b1f8c8e1b5",
  "subcategory_id": "60d5ec49f1b2c8b1f8c8e1b5",
  "brand_id": "60d5ec49f1b2c8b1f8c8e1b6",
  "manufacturer_id": "60d5ec49f1b2c8b1f8c8e1m1",
  "track_inventory": true,
  "track_batches": true,
  "track_serial_numbers": false,
  "valuation_method": "fifo",
  "base_unit_id": "60d5ec49f1b2c8b1f8c8e1b7",
  "allowed_unit_ids": ["60d5ec49f1b2c8b1f8c8e1b7", "60d5ec49f1b2c8b1f8c8e1b8"],
  "weight": 1.5,
  "weight_unit": "kg",
  "length": 10.0,
  "width": 5.0,
  "height": 2.0,
  "dimension_unit": "cm",
  "volume": 0.0001,
  "volume_unit": "m3",
  "location_prices": [
    {
      "location_id": "60d5ec49f1b2c8b1f8c8e1d1",
      "cost_price": 40.0,
      "selling_price": 60.0,
      "mrp": 70.0,
      "initial_stock": 100,
      "currency": "USD"
    }
  ],
  "tax_category_id": "60d5ec49f1b2c8b1f8c8e1t1",
  "hsn_code": "84713010",
  "sac_code": "998313",
  "reorder_level": 50,
  "reorder_quantity": 200,
  "min_stock_level": 20,
  "max_stock_level": 1000,
  "safety_stock": 30,
  "default_supplier_id": "60d5ec49f1b2c8b1f8c8e1s1",
  "supplier_ids": ["60d5ec49f1b2c8b1f8c8e1s1"],
  "lead_time_days": 15,
  "shelf_life_days": 365,
  "requires_qc": true,
  "perishable": false,
  "hazardous": false,
  "images": ["https://example.com/image1.jpg"],
  "thumbnail": "https://example.com/thumb.jpg",
  "specifications": {
    "color": "Red",
    "size": "L"
  },
  "tags": ["featured", "new"]
}
```

**Response (201 Created):**

```json
{
  "success": true,
  "message": "Product created successfully",
  "data": {
    "id": "60d5ec49f1b2c8b1f8c8e1p1",
    "organization_id": "60d5ec49f1b2c8b1f8c8e1a0",
    "sku": "PROD-001",
    "barcode": "8901234567890",
    "name": "Widget X",
    "description": "High quality widget",
    "type": "finished_goods",
    "status": "active",
    "category_id": "60d5ec49f1b2c8b1f8c8e1b5",
    "category": {
      "id": "60d5ec49f1b2c8b1f8c8e1b5",
      "name": "Electronics"
    },
    "subcategory_id": "60d5ec49f1b2c8b1f8c8e1b5",
    "brand_id": "60d5ec49f1b2c8b1f8c8e1b6",
    "brand": {
      "id": "60d5ec49f1b2c8b1f8c8e1b6",
      "name": "Acme"
    },
    "manufacturer_id": "60d5ec49f1b2c8b1f8c8e1m1",
    "track_inventory": true,
    "track_batches": true,
    "track_serial_numbers": false,
    "valuation_method": "fifo",
    "base_unit_id": "60d5ec49f1b2c8b1f8c8e1b7",
    "allowed_unit_ids": [
      "60d5ec49f1b2c8b1f8c8e1b7",
      "60d5ec49f1b2c8b1f8c8e1b8"
    ],
    "weight": 1.5,
    "weight_unit": "kg",
    "length": 10.0,
    "width": 5.0,
    "height": 2.0,
    "dimension_unit": "cm",
    "volume": 0.0001,
    "volume_unit": "m3",
    "location_prices": [
      {
        "location_id": "60d5ec49f1b2c8b1f8c8e1d1",
        "location_name": "Main Warehouse",
        "cost_price": 40.0,
        "selling_price": 60.0,
        "mrp": 70.0,
        "initial_stock": 100,
        "currency": "USD"
      }
    ],
    "tax_category_id": "60d5ec49f1b2c8b1f8c8e1t1",
    "hsn_code": "84713010",
    "sac_code": "998313",
    "reorder_level": 50,
    "reorder_quantity": 200,
    "min_stock_level": 20,
    "max_stock_level": 1000,
    "safety_stock": 30,
    "default_supplier_id": "60d5ec49f1b2c8b1f8c8e1s1",
    "supplier_ids": ["60d5ec49f1b2c8b1f8c8e1s1"],
    "lead_time_days": 15,
    "shelf_life_days": 365,
    "requires_qc": true,
    "perishable": false,
    "hazardous": false,
    "images": ["https://example.com/image1.jpg"],
    "thumbnail": "https://example.com/thumb.jpg",
    "specifications": {
      "color": "Red",
      "size": "L"
    },
    "tags": ["featured", "new"],
    "total_stock": 100,
    "available_stock": 100,
    "allocated_stock": 0,
    "in_transit_stock": 0,
    "stock_value": 4000.0,
    "created_at": "2023-10-27T10:00:00Z",
    "updated_at": "2023-10-27T10:00:00Z"
  }
}
```

### List Products

**Endpoint:** `GET /products`
**Description:** List products with filters.

**Query Parameters:**

- `organization_id` (Required)
- `category_id`
- `brand_id`
- `status`
- `search` (Name, SKU)
- `page`, `limit`
- `location_id` (Filter prices for specific location)

**Response (200 OK):**

```json
{
  "success": true,
  "page": 1,
  "limit": 10,
  "total": 50,
  "data": [
    {
      "id": "60d5ec49f1b2c8b1f8c8e1p1",
      "organization_id": "60d5ec49f1b2c8b1f8c8e1a0",
      "sku": "PROD-001",
      "barcode": "8901234567890",
      "name": "Widget X",
      "description": "High quality widget",
      "type": "finished_goods",
      "status": "active",
      "category_id": "60d5ec49f1b2c8b1f8c8e1b5",
      "category": {
        "id": "60d5ec49f1b2c8b1f8c8e1b5",
        "name": "Electronics"
      },
      "brand_id": "60d5ec49f1b2c8b1f8c8e1b6",
      "brand": {
        "id": "60d5ec49f1b2c8b1f8c8e1b6",
        "name": "Acme"
      },
      "location_prices": [
        {
          "location_id": "60d5ec49f1b2c8b1f8c8e1d1",
          "cost_price": 40.0,
          "selling_price": 60.0
        }
      ],
      "total_stock": 100,
      "stock_value": 4000.0,
      "created_at": "2023-10-27T10:00:00Z",
      "updated_at": "2023-10-27T10:00:00Z"
    }
  ]
}
```

### Get Product

**Endpoint:** `GET /products/:id`
**Description:** Get product details.

**Query Parameters:**

- `location_id` (Optional: Filter prices for this location)

**Response (200 OK):**

```json
{
  "success": true,
  "data": {
    "id": "60d5ec49f1b2c8b1f8c8e1p1",
    "organization_id": "60d5ec49f1b2c8b1f8c8e1a0",
    "sku": "PROD-001",
    "barcode": "8901234567890",
    "name": "Widget X",
    "description": "High quality widget",
    "type": "finished_goods",
    "status": "active",
    "category_id": "60d5ec49f1b2c8b1f8c8e1b5",
    "category": {
      "id": "60d5ec49f1b2c8b1f8c8e1b5",
      "name": "Electronics"
    },
    "subcategory_id": "60d5ec49f1b2c8b1f8c8e1b5",
    "brand_id": "60d5ec49f1b2c8b1f8c8e1b6",
    "brand": {
      "id": "60d5ec49f1b2c8b1f8c8e1b6",
      "name": "Acme"
    },
    "manufacturer_id": "60d5ec49f1b2c8b1f8c8e1m1",
    "track_inventory": true,
    "track_batches": true,
    "track_serial_numbers": false,
    "valuation_method": "fifo",
    "base_unit_id": "60d5ec49f1b2c8b1f8c8e1b7",
    "allowed_unit_ids": [
      "60d5ec49f1b2c8b1f8c8e1b7",
      "60d5ec49f1b2c8b1f8c8e1b8"
    ],
    "weight": 1.5,
    "weight_unit": "kg",
    "length": 10.0,
    "width": 5.0,
    "height": 2.0,
    "dimension_unit": "cm",
    "volume": 0.0001,
    "volume_unit": "m3",
    "location_prices": [
      {
        "location_id": "60d5ec49f1b2c8b1f8c8e1d1",
        "location_name": "Main Warehouse",
        "cost_price": 40.0,
        "selling_price": 60.0,
        "mrp": 70.0,
        "initial_stock": 100,
        "currency": "USD"
      }
    ],
    "tax_category_id": "60d5ec49f1b2c8b1f8c8e1t1",
    "hsn_code": "84713010",
    "sac_code": "998313",
    "reorder_level": 50,
    "reorder_quantity": 200,
    "min_stock_level": 20,
    "max_stock_level": 1000,
    "safety_stock": 30,
    "default_supplier_id": "60d5ec49f1b2c8b1f8c8e1s1",
    "supplier_ids": ["60d5ec49f1b2c8b1f8c8e1s1"],
    "lead_time_days": 15,
    "shelf_life_days": 365,
    "requires_qc": true,
    "perishable": false,
    "hazardous": false,
    "images": ["https://example.com/image1.jpg"],
    "thumbnail": "https://example.com/thumb.jpg",
    "specifications": {
      "color": "Red",
      "size": "L"
    },
    "tags": ["featured", "new"],
    "total_stock": 100,
    "available_stock": 100,
    "allocated_stock": 0,
    "in_transit_stock": 0,
    "stock_value": 4000.0,
    "created_at": "2023-10-27T10:00:00Z",
    "updated_at": "2023-10-27T10:00:00Z"
  }
}
```

### Get Product by SKU

**Endpoint:** `GET /products/sku/:sku`
**Description:** Get product by SKU.

**Query Parameters:**

- `organization_id` (Required)

**Response (200 OK):**

```json
{
  "success": true,
  "data": {
    "id": "60d5ec49f1b2c8b1f8c8e1p1",
    "organization_id": "60d5ec49f1b2c8b1f8c8e1a0",
    "sku": "PROD-001",
    "barcode": "8901234567890",
    "name": "Widget X",
    "description": "High quality widget",
    "type": "finished_goods",
    "status": "active",
    "category_id": "60d5ec49f1b2c8b1f8c8e1b5",
    "brand_id": "60d5ec49f1b2c8b1f8c8e1b6",
    "track_inventory": true,
    "total_stock": 100,
    "available_stock": 100,
    "created_at": "2023-10-27T10:00:00Z",
    "updated_at": "2023-10-27T10:00:00Z"
  }
}
```

### Update Product

**Endpoint:** `PUT /products/:id`
**Description:** Update product details.

**Request Body:** (Partial updates allowed)

```json
{
  "name": "Widget X Pro",
  "selling_price": 65.0,
  "description": "Updated description",
  "status": "active"
}
```

**Response (200 OK):**

```json
{
  "success": true,
  "message": "Product updated successfully",
  "data": {
    "id": "60d5ec49f1b2c8b1f8c8e1p1",
    "organization_id": "60d5ec49f1b2c8b1f8c8e1a0",
    "sku": "PROD-001",
    "name": "Widget X Pro",
    "description": "Updated description",
    "status": "active",
    "created_at": "2023-10-27T10:00:00Z",
    "updated_at": "2023-10-28T10:00:00Z"
  }
}
```

### Delete Product

**Endpoint:** `DELETE /products/:id`
**Description:** Soft delete a product (only allowed if no stock).

**Response (200 OK):**

```json
{
  "success": true,
  "message": "Product deleted successfully"
}
```

---

## 2. Brands

### Create Brand

**Endpoint:** `POST /brands`
**Description:** Create a new brand.

**Request Body:**

```json
{
  "organization_id": "60d5ec49f1b2c8b1f8c8e1a0",
  "name": "Acme",
  "code": "ACME",
  "description": "Best Brand",
  "website": "https://acme.com",
  "logo_url": "https://acme.com/logo.png",
  "country": "USA",
  "is_active": true,
  "metadata": {
    "key": "value"
  }
}
```

**Response (201 Created):**

```json
{
  "success": true,
  "message": "Brand created successfully",
  "data": {
    "id": "60d5ec49f1b2c8b1f8c8e1b6",
    "organization_id": "60d5ec49f1b2c8b1f8c8e1a0",
    "name": "Acme",
    "code": "ACME",
    "description": "Best Brand",
    "website": "https://acme.com",
    "logo_url": "https://acme.com/logo.png",
    "country": "USA",
    "is_active": true,
    "created_at": "2023-10-27T10:00:00Z",
    "updated_at": "2023-10-27T10:00:00Z"
  }
}
```

### List Brands

**Endpoint:** `GET /brands`
**Description:** List brands.

**Query Parameters:**

- `org_id` (Required)
- `q` (Search query)
- `is_active` (true/false)
- `page`, `limit`

**Response (200 OK):**

```json
{
  "success": true,
  "data": {
    "brands": [
      {
        "id": "60d5ec49f1b2c8b1f8c8e1b6",
        "name": "Acme",
        "code": "ACME",
        "is_active": true
      }
    ],
    "pagination": {
      "page": 1,
      "limit": 20,
      "total": 1
    }
  }
}
```

### Search Brands

**Endpoint:** `GET /brands/search`
**Description:** Search brands by name or code.

**Query Parameters:**

- `org_id` (Required)
- `q` (Required)

---

## 3. Categories

### Create Category

**Endpoint:** `POST /categories`
**Description:** Create a product category.

**Request Body:**

```json
{
  "organization_id": "60d5ec49f1b2c8b1f8c8e1a0",
  "name": "Electronics",
  "code": "ELEC",
  "description": "Electronic Items",
  "is_active": true,
  "subcategories": [
    {
      "name": "Mobile Phones",
      "code": "MOB",
      "description": "Smartphones",
      "is_active": true
    }
  ]
}
```

**Response (201 Created):**

```json
{
  "success": true,
  "message": "Category created successfully",
  "data": {
    "id": "60d5ec49f1b2c8b1f8c8e1b5",
    "organization_id": "60d5ec49f1b2c8b1f8c8e1a0",
    "name": "Electronics",
    "code": "ELEC",
    "description": "Electronic Items",
    "is_active": true,
    "product_count": 0,
    "subcategories": [
      {
        "name": "Mobile Phones",
        "code": "MOB",
        "description": "Smartphones",
        "is_active": true,
        "product_count": 0
      }
    ],
    "created_at": "2023-10-27T10:00:00Z",
    "updated_at": "2023-10-27T10:00:00Z"
  }
}
```

### List Categories

**Endpoint:** `GET /categories`
**Description:** List categories.

**Query Parameters:**

- `organization_id` (Required)
- `q` (Search)
- `is_active`
- `page`, `limit`

### Get Category Tree

**Endpoint:** `GET /categories/tree`
**Description:** Get full category hierarchy.

**Query Parameters:**

- `organization_id` (Required)

**Response (200 OK):**

```json
{
  "success": true,
  "data": [
    {
      "id": "60d5ec49f1b2c8b1f8c8e1b5",
      "name": "Electronics",
      "code": "ELEC",
      "subcategories": [
        {
          "id": "60d5ec49f1b2c8b1f8c8e1s5",
          "name": "Mobile Phones",
          "code": "MOB"
        }
      ]
    }
  ]
}
```

---

## 4. Units

### List Units

**Endpoint:** `GET /units`
**Description:** List all available units of measure.

**Response (200 OK):**

```json
{
  "success": true,
  "data": [
    {
      "id": "60d5ec49f1b2c8b1f8c8e1u1",
      "name": "Kilogram",
      "symbol": "kg",
      "type": "mass",
      "is_base_unit": true,
      "is_active": true
    },
    {
      "id": "60d5ec49f1b2c8b1f8c8e1u2",
      "name": "Piece",
      "symbol": "pcs",
      "type": "quantity",
      "is_base_unit": true,
      "is_active": true
    }
  ]
}
```
