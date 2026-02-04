# Organization Service API Documentation

Base URL: `/api/v1`

## Authentication

All endpoints require a valid JWT token in the `Authorization` header:
`Authorization: Bearer <token>`

## 1. Organizations

### Create Organization

**Endpoint:** `POST /organizations`
**Description:** Create a new organization.

**Request Body:**

```json
{
  "name": "Acme Holdings",
  "legal_name": "Acme Holdings Inc.",
  "domain": "acme.com",
  "email": "admin@acme.com",
  "billing_email": "billing@acme.com",
  "phone": "+1-555-0100",
  "website": "https://acme.com",
  "tax_id": "TAX-999",
  "registration_number": "REG-ORG-001",
  "company_size": "medium",
  "industry": "Retail",
  "max_users": 100,
  "max_companies": 5,
  "max_locations": 20,
  "storage_limit_gb": 10.0
}
```

**Response (201 Created):**

```json
{
  "success": true,
  "message": "Organization created successfully",
  "data": {
    "id": "60d5ec49f1b2c8b1f8c8e1z1",
    "name": "Acme Holdings",
    "legal_name": "Acme Holdings Inc.",
    "domain": "acme.com",
    "alternate_domains": [],
    "logo": "",
    "favicon": "",
    "email": "admin@acme.com",
    "phone": "+1-555-0100",
    "website": "https://acme.com",
    "tax_id": "TAX-999",
    "registration_number": "REG-ORG-001",
    "industry": "Retail",
    "company_size": "medium",
    "status": "active",
    "is_active": true,
    "billing_email": "billing@acme.com",
    "max_users": 100,
    "max_companies": 5,
    "max_locations": 20,
    "current_users": 1,
    "current_companies": 0,
    "current_locations": 0,
    "storage_used_gb": 0.0,
    "storage_limit_gb": 10.0,
    "settings": {
      "timezone": "UTC",
      "date_format": "YYYY-MM-DD",
      "time_format": "HH:mm:ss",
      "currency": "USD",
      "language": "en",
      "enabled_modules": [],
      "allow_user_registration": false,
      "require_email_verification": true,
      "enable_mfa": false,
      "session_timeout": 30,
      "password_policy": {
        "min_length": 8,
        "require_uppercase": true,
        "require_lowercase": true,
        "require_numbers": true,
        "require_special_chars": false,
        "expiry_days": 90,
        "prevent_reuse_count": 5
      }
    },
    "created_at": "2023-10-27T10:00:00Z",
    "updated_at": "2023-10-27T10:00:00Z"
  }
}
```

### List Organizations

**Endpoint:** `GET /organizations`
**Description:** List organizations.

**Query Parameters:**

- `page`, `limit`
- `status` (active, inactive, suspended)

**Response (200 OK):**

```json
{
  "success": true,
  "data": [
    {
      "id": "60d5ec49f1b2c8b1f8c8e1z1",
      "name": "Acme Holdings",
      "domain": "acme.com",
      "status": "active",
      "is_active": true,
      "created_at": "2023-10-27T10:00:00Z"
    }
  ]
}
```

### Get Organization

**Endpoint:** `GET /organizations/:id`
**Description:** Get organization details.

**Response (200 OK):**

```json
{
  "success": true,
  "data": {
    "id": "60d5ec49f1b2c8b1f8c8e1z1",
    "name": "Acme Holdings",
    "legal_name": "Acme Holdings Inc.",
    "domain": "acme.com",
    "email": "admin@acme.com",
    "phone": "+1-555-0100",
    "status": "active",
    "settings": {
      "timezone": "UTC",
      "currency": "USD"
    },
    "created_at": "2023-10-27T10:00:00Z",
    "updated_at": "2023-10-27T10:00:00Z"
  }
}
```

### Update Organization

**Endpoint:** `PUT /organizations/:id`
**Description:** Update organization details.

**Request Body:**

```json
{
  "name": "Acme Global",
  "legal_name": "Acme Global Inc.",
  "phone": "+1-555-9999",
  "website": "https://acmeglobal.com"
}
```

**Response (200 OK):**

```json
{
  "success": true,
  "message": "Organization updated successfully",
  "data": {
    "id": "60d5ec49f1b2c8b1f8c8e1z1",
    "name": "Acme Global",
    "legal_name": "Acme Global Inc.",
    "phone": "+1-555-9999",
    "website": "https://acmeglobal.com",
    "updated_at": "2023-10-28T10:00:00Z"
  }
}
```

### Delete Organization

**Endpoint:** `DELETE /organizations/:id`
**Description:** Soft delete an organization.

**Response (200 OK):**

```json
{
  "success": true,
  "message": "Organization deleted successfully"
}
```

### List All Locations (Org Level)

**Endpoint:** `GET /organizations/:id/locations`
**Description:** List all locations under an organization (across all companies).

**Query Parameters:**

- `page`, `limit`

**Response (200 OK):**

```json
{
  "success": true,
  "data": [
    {
      "id": "60d5ec49f1b2c8b1f8c8e1l1",
      "name": "Downtown Store",
      "code": "LOC-DT",
      "company_id": "60d5ec49f1b2c8b1f8c8e1c1",
      "organization_id": "60d5ec49f1b2c8b1f8c8e1z1"
    }
  ]
}
```

---

## 2. Companies

### Create Company

**Endpoint:** `POST /organizations/:id/companies`
**Description:** Create a new company under an organization.

**Request Body:**

```json
{
  "name": "Acme Retail",
  "legal_name": "Acme Retail LLC",
  "code": "ACME-RET",
  "email": "retail@acme.com",
  "phone": "+1-555-0200",
  "address": {
    "street": "123 Market St",
    "city": "San Francisco",
    "state": "CA",
    "country": "USA",
    "postal_code": "94105",
    "country_code": "US",
    "latitude": 37.7749,
    "longitude": -122.4194
  },
  "tax_id": "TAX-RET-001",
  "registration_number": "REG-001",
  "vat_number": "VAT-USA-001",
  "bank_accounts": [
    {
      "bank_name": "Major Bank",
      "account_number": "1234567890",
      "account_name": "Acme Retail Operating",
      "is_default": true
    }
  ],
  "is_default": true
}
```

**Response (201 Created):**

```json
{
  "success": true,
  "message": "Company created successfully",
  "data": {
    "id": "60d5ec49f1b2c8b1f8c8e1c1",
    "organization_id": "60d5ec49f1b2c8b1f8c8e1z1",
    "name": "Acme Retail",
    "legal_name": "Acme Retail LLC",
    "code": "ACME-RET",
    "email": "retail@acme.com",
    "phone": "+1-555-0200",
    "address": {
      "street": "123 Market St",
      "city": "San Francisco",
      "state": "CA",
      "country": "USA",
      "postal_code": "94105"
    },
    "tax_id": "TAX-RET-001",
    "registration_number": "REG-001",
    "vat_number": "VAT-USA-001",
    "is_active": true,
    "is_default": true,
    "settings": {
      "fiscal_year_start": "01-01",
      "currency": "USD",
      "timezone": "UTC",
      "enable_multi_currency": false
    },
    "created_at": "2023-10-27T10:00:00Z",
    "updated_at": "2023-10-27T10:00:00Z"
  }
}
```

### List Companies

**Endpoint:** `GET /organizations/:id/companies`
**Description:** List companies under an organization.

**Query Parameters:**

- `page`, `limit`

**Response (200 OK):**

```json
{
  "success": true,
  "data": [
    {
      "id": "60d5ec49f1b2c8b1f8c8e1c1",
      "name": "Acme Retail",
      "code": "ACME-RET",
      "is_active": true
    }
  ]
}
```

### Get Company

**Endpoint:** `GET /companies/:id`
**Description:** Get company details.

**Response (200 OK):**

```json
{
  "success": true,
  "data": {
    "id": "60d5ec49f1b2c8b1f8c8e1c1",
    "name": "Acme Retail",
    "legal_name": "Acme Retail LLC",
    "code": "ACME-RET",
    "email": "retail@acme.com",
    "address": {
      "street": "123 Market St",
      "city": "San Francisco"
    },
    "is_active": true,
    "created_at": "2023-10-27T10:00:00Z"
  }
}
```

### Update Company

**Endpoint:** `PUT /companies/:id`
**Description:** Update company details.

**Request Body:**

```json
{
  "name": "Acme Retail Inc.",
  "is_active": true
}
```

**Response (200 OK):**

```json
{
  "success": true,
  "message": "Company updated successfully",
  "data": {
    "id": "60d5ec49f1b2c8b1f8c8e1c1",
    "name": "Acme Retail Inc.",
    "updated_at": "2023-10-28T10:00:00Z"
  }
}
```

### Delete Company

**Endpoint:** `DELETE /companies/:id`
**Description:** Soft delete a company.

**Response (200 OK):**

```json
{
  "success": true,
  "message": "Company deleted successfully"
}
```

### Assign User to Company

**Endpoint:** `POST /companies/:id/users`
**Description:** Assign a user to a company (activates access to all locations in company).

**Request Body:**

```json
{
  "user_id": "60d5ec49f1b2c8b1f8c8e1u1",
  "role_id": "60d5ec49f1b2c8b1f8c8e1r1"
}
```

**Response (200 OK):**

```json
{
  "success": true,
  "message": "User assigned to company successfully"
}
```

---

## 3. Locations

### Create Location

**Endpoint:** `POST /companies/:id/locations`
**Description:** Create a new location under a company.

**Request Body:**

```json
{
  "name": "Downtown Store",
  "code": "LOC-DT",
  "type": "store",
  "email": "downtown@acme.com",
  "phone": "+1-555-0300",
  "address": {
    "street": "100 Main St",
    "city": "San Francisco",
    "state": "CA",
    "country": "USA",
    "postal_code": "94105"
  },
  "store_info": {
    "floor_area": 1500.0,
    "pos_count": 3,
    "parking_spaces": 20,
    "has_online_ordering": true
  },
  "warehouse_info": {
    "total_area": 5000.0,
    "storage_capacity": 1000,
    "docking_bays": 2,
    "refrigerated_area": 500.0,
    "has_cold_storage": true
  },
  "is_default": false
}
```

**Response (201 Created):**

```json
{
  "success": true,
  "message": "Location created successfully",
  "data": {
    "id": "60d5ec49f1b2c8b1f8c8e1l1",
    "company_id": "60d5ec49f1b2c8b1f8c8e1c1",
    "organization_id": "60d5ec49f1b2c8b1f8c8e1z1",
    "name": "Downtown Store",
    "code": "LOC-DT",
    "type": "store",
    "email": "downtown@acme.com",
    "address": {
      "street": "100 Main St",
      "city": "San Francisco"
    },
    "store_info": {
      "floor_area": 1500.0,
      "pos_count": 3
    },
    "is_active": true,
    "settings": {
      "timezone": "UTC",
      "allow_backdated_transactions": false,
      "require_approval": true
    },
    "created_at": "2023-10-27T10:00:00Z",
    "updated_at": "2023-10-27T10:00:00Z"
  }
}
```

### List Locations (Company Level)

**Endpoint:** `GET /companies/:id/locations`
**Description:** List locations under a company.

**Query Parameters:**

- `page`, `limit`

**Response (200 OK):**

```json
{
  "success": true,
  "data": [
    {
      "id": "60d5ec49f1b2c8b1f8c8e1l1",
      "name": "Downtown Store",
      "code": "LOC-DT",
      "type": "store",
      "is_active": true
    }
  ]
}
```

### Get Location

**Endpoint:** `GET /locations/:id`
**Description:** Get location details.

**Response (200 OK):**

```json
{
  "success": true,
  "data": {
    "id": "60d5ec49f1b2c8b1f8c8e1l1",
    "name": "Downtown Store",
    "code": "LOC-DT",
    "type": "store",
    "address": {
      "street": "100 Main St",
      "city": "San Francisco"
    },
    "store_info": {
      "floor_area": 1500.0
    },
    "is_active": true
  }
}
```

---

## 4. User Access

### Get User Access

**Endpoint:** `GET /users/me/access`
**Description:** Get all organizations, companies, and locations accessible by the logged-in user.

**Response (200 OK):**

```json
{
  "success": true,
  "data": {
    "organization": {
      "id": "60d5ec49f1b2c8b1f8c8e1z1",
      "name": "Acme Holdings"
    },
    "companies": [
      {
        "company": {
          "id": "60d5ec49f1b2c8b1f8c8e1c1",
          "name": "Acme Retail"
        },
        "locations": [
          {
            "id": "60d5ec49f1b2c8b1f8c8e1l1",
            "name": "Downtown Store"
          }
        ]
      }
    ]
  }
}
```
