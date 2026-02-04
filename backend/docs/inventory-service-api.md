# Inventory Service API Documentation

Base URL: `/api/v1`

## Authentication

All endpoints require a valid JWT token in the `Authorization` header:
`Authorization: Bearer <token>`

## 1. Stock Levels

### List Stock Levels

**Endpoint:** `GET /inventory/stock-levels`
**Description:** List stock levels with filters.

**Query Parameters:**

- `organization_id`
- `product_id`
- `location_id`
- `page`, `limit`

**Response (200 OK):**

```json
{
  "success": true,
  "data": [
    {
      "id": "60d5ec49f1b2c8b1f8c8e1s1",
      "organization_id": "60d5ec49f1b2c8b1f8c8e1a0",
      "product_id": "60d5ec49f1b2c8b1f8c8e1p1",
      "location_id": "60d5ec49f1b2c8b1f8c8e1l1",
      "warehouse_zone": "A-01-05",
      "product_name": "Widget X",
      "sku": "PROD-001",
      "location_name": "Downtown Store",
      "quantity_on_hand": 100.0,
      "quantity_available": 80.0,
      "quantity_allocated": 20.0,
      "quantity_in_transit": 0.0,
      "quantity_reserved": 0.0,
      "average_cost": 40.0,
      "last_cost": 42.0,
      "total_value": 4000.0,
      "last_movement_date": "2023-10-27T10:00:00Z",
      "last_count_date": "2023-09-30T10:00:00Z",
      "metadata": {},
      "created_at": "2023-01-01T00:00:00Z",
      "updated_at": "2023-10-27T10:00:00Z"
    }
  ]
}
```

### Get Stock Level by Product

**Endpoint:** `GET /inventory/products/:product_id/stock`
**Description:** Get stock level for a product across all locations.

### Get Stock Level by Location

**Endpoint:** `GET /inventory/locations/:location_id/stock`
**Description:** Get stock levels for all products in a location.

### Allocate Stock

**Endpoint:** `POST /inventory/stock/allocate`
**Description:** Reserve stock for orders.

**Request Body:**

```json
{
  "product_id": "60d5ec49f1b2c8b1f8c8e1p1",
  "location_id": "60d5ec49f1b2c8b1f8c8e1l1",
  "quantity": 5.0
}
```

### Release Stock

**Endpoint:** `POST /inventory/stock/release`
**Description:** Release reserved stock.

---

## 2. Stock Movements

### Create Stock Movement

**Endpoint:** `POST /inventory/organizations/:org_id/stock-movements`
**Description:** Record a stock movement.

**Request Body:**

```json
{
  "reference_id": "PO-123",
  "reference_type": "purchase_order",
  "items": [
    {
      "product_id": "60d5ec49f1b2c8b1f8c8e1p1",
      "to_location_id": "60d5ec49f1b2c8b1f8c8e1l1",
      "quantity": 50.0,
      "batch_number": "BATCH-001",
      "unit_cost": 40.0,
      "uom": "pcs"
    }
  ],
  "movement_type": "in",
  "movement_date": "2023-10-27T10:00:00Z",
  "reason": "Purchase Order Receipt",
  "notes": "Received efficiently"
}
```

### List Stock Movements

**Endpoint:** `GET /inventory/organizations/:org_id/stock-movements`
**Description:** List movements history.

**Query Parameters:**

- `movement_type` (inbound, outbound, transfer, adjustment)
- `location_id`
- `product_id`

**Response (200 OK):**

```json
{
  "success": true,
  "data": [
    {
      "id": "60d5ec49f1b2c8b1f8c8e1m1",
      "organization_id": "60d5ec49f1b2c8b1f8c8e1a0",
      "product_id": "60d5ec49f1b2c8b1f8c8e1p1",
      "movement_type": "in",
      "to_location_id": "60d5ec49f1b2c8b1f8c8e1l1",
      "quantity": 50.0,
      "uom": "pcs",
      "unit_cost": 40.0,
      "total_cost": 2000.0,
      "reference_type": "purchase_order",
      "reference_no": "PO-123",
      "movement_date": "2023-10-27T10:00:00Z",
      "product_name": "Widget X",
      "to_location_name": "Downtown Store"
    }
  ]
}
```

---

## 3. Stock Adjustments

### Create Stock Adjustment

**Endpoint:** `POST /inventory/organizations/:org_id/stock-adjustments`
**Description:** Create a stock adjustment request (requires approval).

**Request Body:**

```json
{
  "location_id": "60d5ec49f1b2c8b1f8c8e1l1",
  "reason": "Damaged Goods",
  "reason_details": "Water damage during transit",
  "items": [
    {
      "product_id": "60d5ec49f1b2c8b1f8c8e1p1",
      "actual_qty": 45.0,
      "adjustment_qty": -5.0,
      "unit_cost": 10.0,
      "difference_qty": -5.0,
      "reason": "Damaged"
    }
  ],
  "adjustment_date": "2023-10-28T10:00:00Z",
  "notes": "Immediate attention required"
}
```

### Approve Stock Adjustment

**Endpoint:** `POST /inventory/stock-adjustments/:id/approve`
**Description:** Approve and apply inventory changes.

### Reject Stock Adjustment

**Endpoint:** `POST /inventory/stock-adjustments/:id/reject`
**Description:** Reject adjustment request.

---

## 4. Inventory Counts (Stocktaking)

### Create Inventory Count

**Endpoint:** `POST /inventory/organizations/:org_id/inventory-counts`
**Description:** Start a stocktaking session.

**Request Body:**

```json
{
  "location_id": "60d5ec49f1b2c8b1f8c8e1l1",
  "count_type": "full",
  "notes": "Year-end count",
  "count_date": "2023-12-31T00:00:00Z"
}
```

### Update Count Item

**Endpoint:** `POST /inventory/inventory-counts/:id/items`
**Description:** Record counted quantity for an item.

**Request Body:**

```json
{
  "product_id": "60d5ec49f1b2c8b1f8c8e1p1",
  "counted_qty": 100.0,
  "notes": "Verified twice"
}
```

### Complete Inventory Count

**Endpoint:** `POST /inventory/inventory-counts/:id/complete`
**Description:** Finalize count. Optionally auto-create adjustment for discrepancies.

**Request Body:**

```json
{
  "create_adjustment": true
}
```

---

## 5. Serial Numbers

### Create Serial Number

**Endpoint:** `POST /inventory/organizations/:org_id/serial-numbers`
**Description:** Register individual serial numbers for tracked products.

### List Serial Numbers

**Endpoint:** `GET /inventory/products/:product_id/serial-numbers`
**Description:** List serials for a product.

### Allocate Serial Number

**Endpoint:** `POST /inventory/serial-numbers/:id/allocate`
**Description:** Assign serial number to a customer/order.

### Mark Serial as Sold

**Endpoint:** `POST /inventory/serial-numbers/:id/sold`
**Description:** Mark as sold (outbound).

---

## 6. Batches

### Create Batch

**Endpoint:** `POST /inventory/organizations/:org_id/batches`
**Description:** Create a new batch.

### Get Batches by Product

**Endpoint:** `GET /inventory/products/:product_id/batches`
**Description:** List batches for a product.

---

## 7. Units & Conversions

### Create Unit

**Endpoint:** `POST /api/v1/units`
**Description:** Create a unit of measure.

### Create Unit Conversion (Chart)

**Endpoint:** `POST /api/v1/unit-charts`
**Description:** Define conversion rule.

**Request Body:**

```json
{
  "from_unit_id": "60d5ec49f1b2c8b1f8c8e1u3",
  "to_unit_id": "60d5ec49f1b2c8b1f8c8e1u2",
  "conversion_factor": 12.0,
  "operator": "multiply",
  "is_active": true
}
```

### Get Conversion Rate

**Endpoint:** `GET /api/v1/unit-charts/conversion-rate`
**Query Parameters:** `from_unit_id`, `to_unit_id`
