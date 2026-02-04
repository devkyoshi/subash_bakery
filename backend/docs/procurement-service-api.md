# Procurement Service API Documentation

Base URL: `/api/v1`

## Authentication

All endpoints require a valid JWT token in the `Authorization` header:
`Authorization: Bearer <token>`

## 1. Suppliers

### Create Supplier

**Endpoint:** `POST /organizations/:org_id/suppliers`
**Description:** Create a new supplier for an organization.

**Request Body:**

```json
{
  "company_name": "Acme Corp",
  "contact_person": "John Doe",
  "email": "john@acme.com",
  "phone": "+1234567890",
  "mobile": "+0987654321",
  "website": "https://acme.com",
  "address": {
    "street": "123 Main St",
    "city": "Metropolis",
    "state": "NY",
    "zip_code": "10001",
    "country": "USA"
  },
  "tax_id": "TAX-12345",
  "payment_terms": 30,
  "credit_limit": 10000.0,
  "bank_name": "City Bank",
  "account_number": "123456789",
  "account_name": "Acme Corp Inc.",
  "swift_code": "CITYUS33",
  "currency": "USD",
  "tags": ["preferred", "electronics"],
  "notes": "Key supplier for parts"
}
```

**Response (201 Created):**

```json
{
  "success": true,
  "message": "Supplier created successfully",
  "data": {
    "id": "60d5ec49f1b2c8b1f8c8e1a1",
    "organization_id": "60d5ec49f1b2c8b1f8c8e1a0",
    "supplier_code": "SUP-123456",
    "company_name": "Acme Corp",
    "status": "active",
    "contact_person": "John Doe",
    "email": "john@acme.com",
    "phone": "+1234567890",
    "mobile": "+0987654321",
    "website": "https://acme.com",
    "address": {
      "street": "123 Main St",
      "city": "Metropolis",
      "state": "NY",
      "zip_code": "10001",
      "country": "USA"
    },
    "tax_id": "TAX-12345",
    "payment_terms": 30,
    "credit_limit": 10000.0,
    "currency": "USD",
    "bank_name": "City Bank",
    "account_number": "123456789",
    "account_name": "Acme Corp Inc.",
    "swift_code": "CITYUS33",
    "rating": 0.0,
    "lead_time": 0,
    "total_orders": 0,
    "total_purchase_value": 0.0,
    "outstanding_balance": 0.0,
    "tags": ["preferred", "electronics"],
    "notes": "Key supplier for parts",
    "created_at": "2023-10-27T10:00:00Z",
    "updated_at": "2023-10-27T10:00:00Z"
  }
}
```

### Get Supplier

**Endpoint:** `GET /suppliers/:id`
**Description:** Retrieve details of a specific supplier.

**Response (200 OK):**

```json
{
  "success": true,
  "message": "Supplier retrieved successfully",
  "data": {
    "id": "60d5ec49f1b2c8b1f8c8e1a1",
    "organization_id": "60d5ec49f1b2c8b1f8c8e1a0",
    "supplier_code": "SUP-123456",
    "company_name": "Acme Corp",
    "status": "active",
    "contact_person": "John Doe",
    "email": "john@acme.com",
    "phone": "+1234567890",
    "address": {
      "street": "123 Main St",
      "city": "Metropolis",
      "country": "USA"
    },
    "rating": 4.5,
    "created_at": "2023-10-27T10:00:00Z"
  }
}
```

### List Suppliers

**Endpoint:** `GET /organizations/:org_id/suppliers`
**Description:** List suppliers for an organization with pagination and filtering.

**Query Parameters:**

- `page`: Page number (default: 1)
- `limit`: Items per page (default: 20)
- `status`: Filter by status (e.g., "active")
- `search`: Search term for company name or contact person

**Response (200 OK):**

```json
{
  "success": true,
  "page": 1,
  "limit": 20,
  "total": 50,
  "data": [
    {
      "id": "60d5ec49f1b2c8b1f8c8e1a1",
      "company_name": "Acme Corp",
      "contact_person": "John Doe",
      "email": "john@acme.com",
      "status": "active"
    }
  ]
}
```

### Update Supplier

**Endpoint:** `PUT /suppliers/:id`
**Description:** Update supplier details.

**Request Body:** (Partial updates allowed)

```json
{
  "company_name": "Acme Corporation",
  "credit_limit": 15000.0
}
```

**Response (200 OK):**

```json
{
  "success": true,
  "message": "Supplier updated successfully",
  "data": {
    "id": "60d5ec49f1b2c8b1f8c8e1a1",
    "company_name": "Acme Corporation",
    "credit_limit": 15000.0,
    "updated_at": "2023-10-28T10:00:00Z"
  }
}
```

### Delete Supplier

**Endpoint:** `DELETE /suppliers/:id`
**Description:** Soft delete a supplier.

**Response (200 OK):**

```json
{
  "success": true,
  "message": "Supplier deleted successfully",
  "data": null
}
```

---

## 2. Purchase Orders (PO)

### Create Purchase Order

**Endpoint:** `POST /organizations/:org_id/purchase-orders`
**Description:** Create a new purchase order.

**Request Body:**

```json
{
  "supplier_id": "60d5ec49f1b2c8b1f8c8e1a1",
  "order_date": "2023-10-27T10:00:00Z",
  "expected_date": "2023-11-10T10:00:00Z",
  "payment_terms": 30,
  "payment_method": "Bank Transfer",
  "items": [
    {
      "product_id": "60d5ec49f1b2c8b1f8c8e1b1",
      "sku": "PROD-001",
      "description": "Widget A",
      "quantity": 100.0,
      "unit_price": 50.0,
      "tax_rate": 10.0,
      "discount_percent": 5.0
    }
  ],
  "shipping_address": {
    "street": "Warehouse 1",
    "city": "Metropolis",
    "country": "USA"
  },
  "notes": "Urgent delivery",
  "terms": "Net 30",
  "reference_number": "REF-2023-001",
  "tax_amount": 500.0,
  "shipping_cost": 100.0,
  "discount_amount": 250.0
}
```

**Response (201 Created):**

```json
{
  "success": true,
  "message": "Purchase order created successfully",
  "data": {
    "id": "60d5ec49f1b2c8b1f8c8e1c1",
    "organization_id": "60d5ec49f1b2c8b1f8c8e1a0",
    "po_number": "PO-123456",
    "status": "draft",
    "supplier_id": "60d5ec49f1b2c8b1f8c8e1a1",
    "supplier_name": "Acme Corp",
    "order_date": "2023-10-27T10:00:00Z",
    "expected_date": "2023-11-10T10:00:00Z",
    "subtotal": 4750.0,
    "tax_amount": 500.0,
    "shipping_cost": 100.0,
    "discount_amount": 250.0,
    "total_amount": 5350.0,
    "currency": "USD",
    "items": [
      {
        "product_id": "60d5ec49f1b2c8b1f8c8e1b1",
        "product_name": "Widget A",
        "sku": "PROD-001",
        "description": "Widget A",
        "quantity": 100.0,
        "quantity_received": 0.0,
        "unit_price": 50.0,
        "tax_rate": 10.0,
        "discount_percent": 5.0,
        "line_total": 4750.0
      }
    ],
    "created_at": "2023-10-27T10:00:00Z",
    "updated_at": "2023-10-27T10:00:00Z"
  }
}
```

### Get Purchase Order

**Endpoint:** `GET /purchase-orders/:id`
**Description:** Retrieve a purchase order by ID. Enriched with product names.

**Response (200 OK):**

```json
{
  "success": true,
  "message": "Purchase order retrieved successfully",
  "data": {
    "id": "60d5ec49f1b2c8b1f8c8e1c1",
    "po_number": "PO-123456",
    "status": "sent",
    "supplier_id": "60d5ec49f1b2c8b1f8c8e1a1",
    "supplier_name": "Acme Corp",
    "total_amount": 5350.0,
    "items": [
      {
        "product_id": "60d5ec49f1b2c8b1f8c8e1b1",
        "product_name": "Widget A",
        "quantity": 100.0,
        "quantity_received": 50.0
      }
    ]
  }
}
```

### List Purchase Orders

**Endpoint:** `GET /organizations/:org_id/purchase-orders`
**Description:** List purchase orders.

**Query Parameters:**

- `page`, `limit`
- `status` (draft, sent, confirmed, received, partial, cancelled)
- `supplier_id`
- `search`

**Response (200 OK):**

```json
{
  "success": true,
  "page": 1,
  "limit": 20,
  "total": 15,
  "data": [
    {
      "id": "60d5ec49f1b2c8b1f8c8e1c1",
      "po_number": "PO-123456",
      "supplier_name": "Acme Corp",
      "status": "sent",
      "total_amount": 5350.0,
      "order_date": "2023-10-27T10:00:00Z"
    }
  ]
}
```

### Update PO Status

**Endpoint:** `PUT /purchase-orders/:id/status`
**Description:** Update the status of a PO manually.

**Request Body:**

```json
{
  "status": "confirmed"
}
```

**Response (200 OK):**

```json
{
  "success": true,
  "message": "Purchase order status updated successfully",
  "data": null
}
```

### Approve Purchase Order

**Endpoint:** `POST /purchase-orders/:id/approve`
**Description:** Mark a PO as approved.

**Response (200 OK):**

```json
{
  "success": true,
  "message": "Purchase order approved successfully",
  "data": null
}
```

### Delete Purchase Order

**Endpoint:** `DELETE /purchase-orders/:id`
**Description:** Soft delete a PO.

**Response (200 OK):**

```json
{
  "success": true,
  "message": "Purchase order deleted successfully",
  "data": null
}
```

---

## 3. Goods Receipt Notes (GRN)

### Create GRN

**Endpoint:** `POST /organizations/:org_id/grns`
**Description:** Create a GRN to record received goods against a PO.

**Request Body:**

```json
{
  "purchase_order_id": "60d5ec49f1b2c8b1f8c8e1c1",
  "location_id": "60d5ec49f1b2c8b1f8c8e1d1",
  "receipt_date": "2023-11-01T09:00:00Z",
  "invoice_number": "INV-999",
  "delivery_note": "DN-888",
  "items": [
    {
      "product_id": "60d5ec49f1b2c8b1f8c8e1b1",
      "received_quantity": 100.0,
      "ordered_quantity": 100.0,
      "unit_cost": 50.0,
      "batch_number": "BATCH-001",
      "expiry_date": "2024-11-01T00:00:00Z",
      "condition": "Good"
    }
  ],
  "notes": "Received in good condition"
}
```

**Response (201 Created):**

```json
{
  "success": true,
  "message": "GRN created successfully",
  "data": {
    "id": "60d5ec49f1b2c8b1f8c8e1e1",
    "organization_id": "60d5ec49f1b2c8b1f8c8e1a0",
    "grn_number": "GRN-123456",
    "status": "received",
    "purchase_order_id": "60d5ec49f1b2c8b1f8c8e1c1",
    "po_number": "PO-123456",
    "supplier_id": "60d5ec49f1b2c8b1f8c8e1a1",
    "location_id": "60d5ec49f1b2c8b1f8c8e1d1",
    "receipt_date": "2023-11-01T09:00:00Z",
    "invoice_number": "INV-999",
    "items": [
      {
        "product_id": "60d5ec49f1b2c8b1f8c8e1b1",
        "sku": "PROD-001",
        "description": "Widget A",
        "received_quantity": 100.0,
        "accepted_quantity": 0.0,
        "rejected_quantity": 0.0,
        "batch_number": "BATCH-001"
      }
    ],
    "created_at": "2023-11-01T09:00:00Z",
    "updated_at": "2023-11-01T09:00:00Z"
  }
}
```

### Get GRN

**Endpoint:** `GET /grns/:id`
**Description:** Retrieve GRN details. Enriched with product, user, supplier, and location names.

**Response (200 OK):**

```json
{
  "success": true,
  "message": "GRN retrieved successfully",
  "data": {
    "id": "60d5ec49f1b2c8b1f8c8e1e1",
    "grn_number": "GRN-123456",
    "status": "received",
    "supplier_name": "Acme Corp",
    "location_name": "Warehouse 1",
    "received_by_name": "John Doe",
    "items": [
      {
        "product_id": "60d5ec49f1b2c8b1f8c8e1b1",
        "product_name": "Widget A",
        "received_quantity": 100.0
      }
    ]
  }
}
```

### List GRNs

**Endpoint:** `GET /organizations/:org_id/grns`
**Description:** List GRNs.

**Query Parameters:**

- `page`, `limit`
- `status` (received, inspected, accepted, rejected)
- `purchase_order_id`
- `search`

**Response (200 OK):**

```json
{
  "success": true,
  "page": 1,
  "limit": 20,
  "total": 5,
  "data": [
    {
      "id": "60d5ec49f1b2c8b1f8c8e1e1",
      "grn_number": "GRN-123456",
      "status": "received",
      "receipt_date": "2023-11-01T09:00:00Z"
    }
  ]
}
```

### Complete Inspection

**Endpoint:** `POST /grns/:id/inspect`
**Description:** Complete QC inspection for a GRN. This triggers stock movement if passed.

**Request Body:**

```json
{
  "qc_status": "passed",
  "qc_notes": "All items checked and verified."
}
```

**Response (200 OK):**

```json
{
  "success": true,
  "message": "Inspection completed successfully",
  "data": null
}
```
