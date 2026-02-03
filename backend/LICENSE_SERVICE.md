# License Service Documentation

## Overview

The License Service manages ERP applications and their licensing at the location level. It supports multiple licensing models (user-based, device-based, usage-based, concurrent, and unlimited) and handles license assignment, activation, usage tracking, and compliance.

## Architecture

**Port**: 8004  
**Database**: MongoDB (`erp_license`)  
**Dependencies**: MongoDB, Redis

## Key Features

- ✅ Application catalog management
- ✅ Multiple license types (user, device, usage, concurrent, unlimited)
- ✅ User license assignments with role-based permissions
- ✅ Device activation and tracking
- ✅ Usage monitoring and quota enforcement
- ✅ License expiry and renewal management
- ✅ Trial period support
- ✅ Real-time usage tracking
- ✅ License suspension and revocation
- ✅ Compliance and audit tracking

## Data Models

### Application
Represents an ERP application/module that can be licensed.

**Key Fields**:
- `code`: Unique identifier (e.g., "CRM", "INV")
- `category`: Application category (inventory, sales, crm, etc.)
- `supported_license_types`: Array of supported license models
- `base_price`, `price_per_user`, `price_per_device`: Pricing structure
- `features`: Array of application features
- `is_public`: Whether app is publicly available for purchase

### LocationLicense
Represents a license assigned to a specific location.

**Key Fields**:
- `organization_id`, `location_id`, `application_id`: Ownership
- `license_type`: Type of license (user_based, device_based, etc.)
- `status`: License status (active, suspended, expired, exceeded, revoked)
- `license_key`, `activation_key`: Unique identifiers
- `max_users`, `max_devices`, `max_transactions`: Usage limits
- `current_users`, `current_devices`, `current_transactions`: Current usage
- `start_date`, `end_date`, `expires_at`: License period
- `is_trial`, `trial_end_date`: Trial support
- `is_perpetual`: Never expires flag
- `billing_cycle`, `price_per_cycle`: Billing information
- `enabled_features`, `disabled_features`: Feature control
- `allowed_ips`, `geo_restrictions`: Access restrictions

### UserLicenseAssignment
Tracks which users are assigned to which licenses.

**Key Fields**:
- `license_id`, `user_id`: Assignment mapping
- `role`: User's role within the application
- `permissions`: Specific permissions granted
- `assigned_at`, `last_accessed_at`: Activity tracking
- `is_active`: Assignment status

### DeviceLicenseAssignment
Tracks which devices are activated with licenses.

**Key Fields**:
- `license_id`, `device_id`: Assignment mapping
- `device_name`, `device_type`, `device_model`: Device info
- `os`, `os_version`, `mac_address`, `ip_address`: Technical details
- `activation_code`: Unique activation identifier
- `is_online`, `last_seen_at`: Online status
- `is_active`: Activation status

### LicenseUsageLog
Tracks detailed usage metrics for usage-based billing.

**Key Fields**:
- `period_start`, `period_end`: Usage period
- `transaction_count`, `api_call_count`: Activity metrics
- `storage_used_gb`, `data_processed_gb`: Storage metrics
- `active_users`, `active_devices`, `peak_concurrent`: User metrics
- `base_cost`, `usage_cost`, `overage_cost`, `total_cost`: Billing

### LicenseAlert
Tracks alerts for license expiry, usage exceeded, etc.

**Key Fields**:
- `alert_type`: Type of alert (expiry_warning, usage_exceeded)
- `severity`: Alert severity (info, warning, critical)
- `title`, `message`: Alert content
- `is_read`, `is_resolved`: Alert status

## API Endpoints

### Applications (Public)

#### List Applications (Public)
```http
GET /api/v1/applications
Query Parameters:
  - category: ApplicationCategory (optional)
  - public: boolean (optional, default: false)
  - page: int (optional, default: 1)
  - limit: int (optional, default: 20)
```

**Response**:
```json
{
  "success": true,
  "data": [
    {
      "id": "...",
      "name": "CRM System",
      "code": "CRM",
      "display_name": "Customer Relationship Management",
      "description": "Complete CRM solution",
      "category": "crm",
      "version": "2.1.0",
      "supported_license_types": ["user_based", "concurrent"],
      "default_license_type": "user_based",
      "base_price": 99.00,
      "price_per_user": 10.00,
      "features": [
        {
          "name": "Contact Management",
          "description": "Manage customer contacts",
          "is_advanced": false,
          "is_ai_powered": false
        },
        {
          "name": "AI Lead Scoring",
          "description": "AI-powered lead prioritization",
          "is_advanced": true,
          "is_ai_powered": true
        }
      ],
      "is_public": true,
      "is_active": true,
      "icon": "https://...",
      "total_licenses": 1250,
      "total_locations": 450
    }
  ],
  "message": "Applications retrieved successfully"
}
```

#### Get Application by ID
```http
GET /api/v1/applications/:id
```

#### Get Application by Code
```http
GET /api/v1/applications/code/:code
```

### Applications (Protected - Admin Only)

#### Create Application
```http
POST /api/v1/applications
Authorization: Bearer <token>
Content-Type: application/json

{
  "name": "Inventory Management",
  "code": "INV",
  "display_name": "Inventory & Stock Management",
  "description": "Complete inventory management solution",
  "category": "inventory",
  "version": "1.0.0",
  "supported_license_types": ["user_based", "device_based", "usage_based"],
  "default_license_type": "user_based",
  "base_price": 149.00,
  "price_per_user": 15.00,
  "price_per_device": 25.00,
  "price_per_transaction": 0.02,
  "minimum_users": 5,
  "included_transactions": 10000,
  "features": [
    {
      "name": "Stock Tracking",
      "description": "Real-time stock level tracking",
      "is_advanced": false,
      "is_ai_powered": false
    }
  ],
  "api_endpoint": "http://inventory-service:8006",
  "service_url": "http://inventory-service:8006",
  "icon": "https://...",
  "is_public": true
}
```

#### Update Application
```http
PUT /api/v1/applications/:id
Authorization: Bearer <token>
```

#### Delete Application
```http
DELETE /api/v1/applications/:id
Authorization: Bearer <token>
```

### Licenses

#### Create License for Location
```http
POST /api/v1/organizations/:org_id/licenses
Authorization: Bearer <token>
Content-Type: application/json

{
  "application_id": "674a1b2c3d4e5f6g7h8i9j0k",
  "location_id": "675b2c3d4e5f6g7h8i9j0k1l",
  "license_type": "user_based",
  "max_users": 50,
  "max_storage_gb": 100,
  "max_api_calls_daily": 100000,
  "billing_cycle": "monthly",
  "price_per_cycle": 750.00,
  "is_trial": true,
  "trial_days": 30,
  "is_perpetual": false,
  "duration_months": 12,
  "enabled_features": ["stock_tracking", "purchase_orders", "reporting"],
  "disabled_features": ["ai_forecasting"]
}
```

**Response**:
```json
{
  "success": true,
  "data": {
    "id": "...",
    "organization_id": "...",
    "location_id": "...",
    "application_id": "...",
    "license_type": "user_based",
    "status": "active",
    "license_key": "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
    "activation_key": "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
    "max_users": 50,
    "current_users": 0,
    "start_date": "2024-12-09T00:00:00Z",
    "trial_end_date": "2025-01-08T00:00:00Z",
    "is_trial": true,
    "billing_cycle": "monthly",
    "price_per_cycle": 0.00,
    "currency": "USD",
    "enabled_features": ["stock_tracking", "purchase_orders", "reporting"]
  },
  "message": "License created successfully"
}
```

#### Get License by ID
```http
GET /api/v1/licenses/:id
Authorization: Bearer <token>
```

#### List Licenses by Organization
```http
GET /api/v1/organizations/:org_id/licenses
Authorization: Bearer <token>
Query Parameters:
  - page: int
  - limit: int
```

#### List Licenses by Location
```http
GET /api/v1/locations/:location_id/licenses
Authorization: Bearer <token>
Query Parameters:
  - active: boolean (optional, show only active licenses)
```

#### Update License Usage
```http
PUT /api/v1/licenses/:id/usage
Authorization: Bearer <token>
Content-Type: application/json

{
  "users": 42,
  "devices": 15,
  "transactions": 85432,
  "storage_gb": 78.5
}
```

**Notes**:
- Automatically checks if limits are exceeded
- Updates license status to "exceeded" if any limit is breached
- Used by applications to report usage in real-time

#### Suspend License
```http
POST /api/v1/licenses/:id/suspend
Authorization: Bearer <token>
Content-Type: application/json

{
  "reason": "Payment overdue"
}
```

#### Activate License
```http
POST /api/v1/licenses/:id/activate
Authorization: Bearer <token>
```

#### Revoke License
```http
POST /api/v1/licenses/:id/revoke
Authorization: Bearer <token>
Content-Type: application/json

{
  "reason": "Contract terminated"
}
```

### User Assignments

#### Assign User to License
```http
POST /api/v1/licenses/:id/users
Authorization: Bearer <token>
Content-Type: application/json

{
  "user_id": "676c3d4e5f6g7h8i9j0k1l2m",
  "role": "admin",
  "permissions": ["read", "write", "delete", "manage_users"]
}
```

**Response**:
```json
{
  "success": true,
  "message": "User assigned to license successfully"
}
```

**Business Rules**:
- For user-based licenses: Checks if `max_users` limit is reached
- Creates assignment record
- Updates license `current_users` count
- User can now access the application

#### Revoke User from License
```http
DELETE /api/v1/licenses/:id/users/:assignment_id
Authorization: Bearer <token>
Content-Type: application/json

{
  "reason": "User left company"
}
```

### Device Assignments

#### Activate Device for License
```http
POST /api/v1/licenses/:id/devices
Authorization: Bearer <token>
Content-Type: application/json

{
  "device_id": "DEVICE-12345-ABCDE",
  "device_name": "Store POS Terminal 1",
  "device_type": "pos",
  "device_model": "HP EliteOne 800",
  "os": "Windows",
  "os_version": "11 Pro",
  "mac_address": "00:1B:44:11:3A:B7",
  "ip_address": "192.168.1.105"
}
```

**Response**:
```json
{
  "success": true,
  "message": "Device activated successfully"
}
```

**Business Rules**:
- For device-based licenses: Checks if `max_devices` limit is reached
- Generates unique activation code
- Device can now access the application
- Tracks device online status

#### Deactivate Device
```http
DELETE /api/v1/licenses/:id/devices/:assignment_id
Authorization: Bearer <token>
Content-Type: application/json

{
  "reason": "Device replaced"
}
```

## License Types Explained

### 1. User-Based Licensing
- License is tied to number of users
- Each user gets individual access
- `max_users` defines limit
- Common for productivity apps (CRM, Project Management)

**Example**: CRM license for 50 users @ $10/user/month = $500/month

### 2. Device-Based Licensing
- License is tied to physical devices
- Device must be activated with activation code
- `max_devices` defines limit
- Common for POS systems, factory equipment

**Example**: POS system license for 10 terminals @ $25/terminal/month = $250/month

### 3. Usage-Based Licensing
- License based on consumption metrics
- Tracked: transactions, API calls, storage, data processed
- `max_transactions` defines included amount
- Overage charges apply
- Common for transaction processing, API services

**Example**: Payment processing @ $0.02/transaction, 100K included, overage @ $0.025/transaction

### 4. Concurrent Licensing
- Based on simultaneous users
- Users share a pool of licenses
- `max_concurrent` defines pool size
- First-come, first-served access
- Common for expensive professional tools

**Example**: CAD software with 10 concurrent seats for 50 users

### 5. Unlimited Licensing
- No limits on users, devices, or usage
- Typically enterprise tier
- Flat monthly/annual fee
- Common for large organizations

**Example**: Enterprise plan @ $10,000/month unlimited everything

## Workflows

### Complete License Purchase Flow

1. **Customer browses available applications**
   ```http
   GET /api/v1/applications?public=true
   ```

2. **Customer selects application and creates license**
   ```http
   POST /api/v1/organizations/:org_id/licenses
   {
     "application_id": "...",
     "location_id": "...",
     "license_type": "user_based",
     "max_users": 50,
     "is_trial": true,
     "trial_days": 30
   }
   ```

3. **System creates license with trial status**
   - `status`: "active"
   - `is_trial`: true
   - `price_per_cycle`: 0 (no charge during trial)
   - Generates unique `license_key` and `activation_key`

4. **Admin assigns users to license**
   ```http
   POST /api/v1/licenses/:id/users
   {
     "user_id": "...",
     "role": "admin",
     "permissions": ["read", "write"]
   }
   ```

5. **Users can now access the application**
   - Application validates license via `license_key`
   - Checks user assignment
   - Grants access based on permissions

6. **Application reports usage**
   ```http
   PUT /api/v1/licenses/:id/usage
   {
     "users": 42,
     "transactions": 5234
   }
   ```

7. **Trial expires, license converts to paid**
   - Subscription service handles billing
   - License `price_per_cycle` updated to actual price
   - `is_trial` set to false

8. **License renewal**
   - Auto-renewal if `auto_renew` is true
   - Payment processed via subscription service
   - `end_date` extended

### Device Activation Flow

1. **Admin activates device**
   ```http
   POST /api/v1/licenses/:id/devices
   {
     "device_id": "UNIQUE-DEVICE-ID",
     "device_name": "POS Terminal 1"
   }
   ```

2. **System generates activation code**
   - Returns activation code to admin
   - Device record created with `is_online`: false

3. **Device client uses activation code to authenticate**
   - Device sends activation code to application
   - Application validates with license service
   - Device granted access

4. **Device reports heartbeat**
   - Periodic heartbeat updates `last_seen_at`
   - `is_online` set to true
   - If heartbeat stops, device marked offline

5. **Admin can view active devices**
   ```http
   GET /api/v1/licenses/:id
   ```
   - Shows all activated devices
   - Online/offline status
   - Last seen timestamp

## Usage Monitoring & Enforcement

### Real-Time Limit Checking

When usage is updated, the system checks:
```go
if license.MaxUsers > 0 && usage.Users > license.MaxUsers {
    license.Status = "exceeded"
}
```

Limit types checked:
- `max_users` vs `current_users`
- `max_devices` vs `current_devices`
- `max_transactions` vs `current_transactions`
- `max_storage_gb` vs `current_storage_gb`
- `max_api_calls_daily` vs `current_api_calls_today`

### Automatic Status Updates

License status automatically changes based on:
- **Active**: All limits within bounds, not expired
- **Exceeded**: Any usage limit breached
- **Expired**: `expires_at` < current time (for non-perpetual)
- **Suspended**: Manual suspension (e.g., payment issue)
- **Revoked**: Manual revocation (e.g., contract terminated)

### Alert Generation

System generates `LicenseAlert` records for:
- Approaching expiry (e.g., 30, 7, 1 days before)
- Usage approaching limits (e.g., 80%, 90%, 100%)
- License exceeded
- Device offline for extended period

## Database Indexes

### Applications
- `code` (unique)
- `category + is_active`
- `is_public + is_active`
- `deleted_at`

### Location Licenses
- `organization_id + deleted_at`
- `location_id + status`
- `location_id + application_id`
- `application_id + status`
- `status + expires_at`
- `license_key` (unique)
- `deleted_at`

### User License Assignments
- `license_id + is_active`
- `user_id + is_active`
- `user_id + application_id + is_active`
- `organization_id`

### Device License Assignments
- `license_id + is_active`
- `device_id + license_id`
- `organization_id`
- `is_online + last_seen_at`

## Integration with Other Services

### Organization Service
- Validates `organization_id` and `location_id` before creating license
- Checks location exists and belongs to organization

### Subscription Service
- Links licenses to subscription plans
- Subscription determines which applications are available
- Handles billing for licenses

### Individual Application Services
- Each ERP application validates user/device access via license service
- Reports usage metrics back to license service
- Checks feature flags (`enabled_features`, `disabled_features`)

## Security Features

- ✅ JWT authentication for all protected endpoints
- ✅ Unique license keys (UUID)
- ✅ Activation codes for device security
- ✅ IP whitelisting support (`allowed_ips`)
- ✅ Geo-restrictions support (`geo_restrictions`)
- ✅ Email domain restrictions (`allowed_domains`)
- ✅ Audit trail (created_by, updated_by, deleted_by)
- ✅ Soft deletes for data recovery

## Compliance & Audit

Each license can track:
- `compliance_level`: Regulatory compliance (HIPAA, SOC2, GDPR, etc.)
- `last_audit_date`: When license was last audited
- `certification_expiry`: When compliance certification expires

All license modifications tracked:
- Who created the license
- Who modified usage/status
- Who revoked/suspended license
- Complete timestamp history

## Summary

The License Service provides comprehensive application licensing management for the ERP system with:
- **15+ API endpoints**
- **6 core data models** (Application, LocationLicense, UserLicenseAssignment, DeviceLicenseAssignment, LicenseUsageLog, LicenseAlert)
- **5 license types** (user, device, usage, concurrent, unlimited)
- **Real-time usage tracking** and enforcement
- **Trial period** support
- **Multi-tenant** architecture
- **Production-ready** with indexes, validation, error handling
- **~2500 lines of code**

This service enables the flexible, scalable licensing model required for a modern SaaS ERP platform.
