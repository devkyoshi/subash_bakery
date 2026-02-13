# PO vs GRN Comparison Report

## What Is It?

The **PO vs GRN Comparison Report** answers a simple question: **"Did we receive what we ordered?"**

It compares every line item on a **Purchase Order (PO)** against the corresponding **Goods Receipt Notes (GRNs)** to identify quantity mismatches — shortages, excess deliveries, or perfect matches.

---

## Data Sources

| Collection | Database | Purpose |
|---|---|---|
| `purchase_orders` | `erp_procurement` | What was ordered (items, quantities, prices) |
| `goods_receipt_notes` | `erp_procurement` | What was actually received |
| `suppliers` | `erp_procurement` | Supplier name resolution |

The report service **reads** from the procurement database — it does not write to it.

---

## How Comparison Items Are Built

The core unit is a **line-item comparison** — one row per product per PO.

### Step-by-step

1. **Fetch POs** matching the filters (date range, supplier, status) for the current page.
2. **Fetch all GRNs** linked to those POs via `purchase_order_id`.
3. **Group GRNs** by PO: `PO ID → [GRN₁, GRN₂, ...]`
4. **For each PO item** (a specific product on a PO):
   - Loop through all GRNs for that PO.
   - Find GRN items where `grnItem.ProductID == poItem.ProductID`.
   - **Sum up** `ReceivedQuantity`, `AcceptedQuantity`, and `RejectedQuantity` across all matching GRN items.

> A single PO item can have quantities spread across multiple GRNs (e.g. partial deliveries on different dates). The report aggregates all of them.

### Quantity Calculations

| Field | Formula |
|---|---|
| `variance` | `GRN Qty − PO Qty` |
| `variance_pct` | `(variance / PO Qty) × 100` |
| `po_value` | `PO Qty × Unit Price` |
| `grn_value` | `GRN Qty × Unit Price` |
| `value_variance` | `GRN Value − PO Value` |

### Status Determination

Each line item gets one of four statuses:

| Status | Condition |
|---|---|
| **PENDING** | `GRN Qty == 0` and PO is not in `received` status |
| **MATCHED** | `|variance| < 0.01` (effectively zero) |
| **EXCESS** | `variance > 0` (received more than ordered) |
| **PARTIAL** | `variance < 0` (received less than ordered) |

---

## Summary Metrics

Calculated across **all** matching POs (not just the current page).

| Metric | How It's Calculated |
|---|---|
| `total_pos` | Sum of all PO status counts from aggregation |
| `completed_pos` | Count of POs with status `received` |
| `partial_pos` | Count of POs with status `partial` |
| `pending_pos` | Count of POs with status `draft` + `sent` + `confirmed` |
| `excess_pos` | Count of line items on current page with status `EXCESS` |
| `total_po_value` | Sum of `po_value` across all items on current page |
| `total_grn_value` | Sum of `grn_value` across all items on current page |
| `total_variance` | Sum of `|value_variance|` across all items (absolute values) |
| `variance_percent` | `(total_variance / total_po_value) × 100` |
| `completed_percent` | `(completed_pos / total_pos) × 100` |

---

## Variance Distribution (Pie Chart)

Shows the percentage breakdown of PO statuses for charting:

| Segment | Value | Color |
|---|---|---|
| Matched | `(completed_pos / total_pos) × 100` | 🟢 `#22c55e` |
| Shortage | `(partial_pos / total_pos) × 100` | 🟡 `#eab308` |
| Excess | `(excess_pos / total_pos) × 100` | 🔵 `#3b82f6` |
| Pending | `(pending_pos / total_pos) × 100` | 🔴 `#ef4444` |

---

## Action Items

Auto-generated alerts for the **top 5 largest variances** (by absolute quantity).

Only items with status `PARTIAL` or `EXCESS` are considered (not `MATCHED` or `PENDING`).

### Severity Rules

| Condition | Type | Title |
|---|---|---|
| Shortage with `|variance_%| > 10%` | `critical` | High Discrepancy |
| Shortage with `|variance_%| ≤ 10%` | `warning` | Review Required |
| Excess delivery | `info` | Excess Delivery |

Each action item includes the supplier name, variance amount, percentage, PO number, and product name for quick context.

---

## API Endpoints

All endpoints require JWT authentication.

### Get Report Data (Paginated)

```
GET /api/v1/organizations/:org_id/reports/po-vs-grn
```

**Query Parameters:**

| Param | Type | Description |
|---|---|---|
| `start_date` | `YYYY-MM-DD` | Filter POs from this order date |
| `end_date` | `YYYY-MM-DD` | Filter POs until this order date |
| `supplier_id` | `ObjectID hex` | Filter by supplier |
| `status` | `string` | Filter by PO status (`draft`, `sent`, `confirmed`, `received`, `partial`) |
| `page` | `int` | Page number (default: 1) |
| `limit` | `int` | Items per page (default: 10) |

**Response structure:**

```json
{
  "success": true,
  "data": {
    "data": {
      "metrics": { ... },
      "variance_distribution": [ ... ],
      "items": [ ... ],
      "action_items": [ ... ],
      "total_items": 42
    },
    "pagination": {
      "page": 1,
      "limit": 10,
      "total": 42,
      "total_pages": 5
    }
  }
}
```

### Export as Excel

```
GET /api/v1/organizations/:org_id/reports/po-vs-grn/export/excel
```

Returns an `.xlsx` file with two sheets:
- **Summary** — metrics table (Total POs, values, variance %, completion rate)
- **PO vs GRN Detail** — full item-level comparison with conditional formatting (red for negative variance, green for positive)

### Export as PDF

```
GET /api/v1/organizations/:org_id/reports/po-vs-grn/export/pdf
```

Returns a landscape A4 PDF with:
- Summary metrics table on page 1
- Detail table with color-coded variance columns, auto page breaks

> Both export endpoints fetch **all** matching data (no pagination) so the full report is included.

---

## Example Scenario

**PO-2026-001** orders 100 kg of flour from Supplier A at ₹50/kg.

Two GRNs are received:
- GRN-001: 60 kg received, 58 accepted, 2 rejected
- GRN-002: 30 kg received, 30 accepted, 0 rejected

The report row would show:

| Field | Value |
|---|---|
| PO Qty | 100 |
| GRN Qty | 90 (60 + 30) |
| Accepted Qty | 88 (58 + 30) |
| Rejected Qty | 2 (2 + 0) |
| Variance | −10 |
| Variance % | −10.00% |
| PO Value | ₹5,000 |
| GRN Value | ₹4,500 |
| Value Variance | −₹500 |
| Status | **PARTIAL** |

This would generate a **critical** action item ("High Discrepancy") because the variance is exactly 10%.
