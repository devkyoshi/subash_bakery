package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SupplierStatus string

const (
	SupplierStatusActive    SupplierStatus = "active"
	SupplierStatusInactive  SupplierStatus = "inactive"
	SupplierStatusBlacklist SupplierStatus = "blacklist"
)

type POStatus string

const (
	POStatusDraft             POStatus = "draft"
	POStatusSent              POStatus = "sent"
	POStatusConfirmed         POStatus = "confirmed"
	POStatusReceived          POStatus = "received"
	POStatusPartiallyReceived POStatus = "partial"
	POStatusCancelled         POStatus = "cancelled"
	POStatusRejected          POStatus = "rejected"
)

type GRNStatus string

const (
	GRNStatusDraft     GRNStatus = "draft"
	GRNStatusReceived  GRNStatus = "received"
	GRNStatusInspected GRNStatus = "inspected"
	GRNStatusAccepted  GRNStatus = "accepted"
	GRNStatusRejected  GRNStatus = "rejected"
)

type Supplier struct {
	BaseModel `bson:",inline"`

	OrganizationID primitive.ObjectID `bson:"organization_id" json:"organization_id" binding:"required"`
	SupplierCode   string             `bson:"supplier_code" json:"supplier_code" binding:"required"`
	CompanyName    string             `bson:"company_name" json:"company_name" binding:"required"`
	Status         SupplierStatus     `bson:"status" json:"status"`

	// Contact
	ContactPerson string `bson:"contact_person" json:"contact_person"`
	Email         string `bson:"email" json:"email"`
	Phone         string `bson:"phone" json:"phone"`
	Mobile        string `bson:"mobile" json:"mobile"`
	Website       string `bson:"website" json:"website"`

	// Address
	Address Address `bson:"address" json:"address"`

	// Financial
	TaxID        string  `bson:"tax_id" json:"tax_id"`
	PaymentTerms int     `bson:"payment_terms" json:"payment_terms"`
	CreditLimit  float64 `bson:"credit_limit" json:"credit_limit"`
	Currency     string  `bson:"currency" json:"currency"`

	// Bank Details
	BankName      string `bson:"bank_name" json:"bank_name"`
	AccountNumber string `bson:"account_number" json:"account_number"`
	AccountName   string `bson:"account_name" json:"account_name"`
	SwiftCode     string `bson:"swift_code" json:"swift_code"`

	// Analytics
	TotalOrders        int        `bson:"total_orders" json:"total_orders"`
	TotalPurchaseValue float64    `bson:"total_purchase_value" json:"total_purchase_value"`
	OutstandingBalance float64    `bson:"outstanding_balance" json:"outstanding_balance"`
	LastOrderDate      *time.Time `bson:"last_order_date" json:"last_order_date"`

	// Rating
	Rating   float64 `bson:"rating" json:"rating"`       // 0-5
	LeadTime int     `bson:"lead_time" json:"lead_time"` // Days

	Notes    string                 `bson:"notes" json:"notes"`
	Metadata map[string]interface{} `bson:"metadata" json:"metadata"`
	Tags     []string               `bson:"tags" json:"tags"`
}

type PurchaseOrder struct {
	BaseModel `bson:",inline"`

	OrganizationID primitive.ObjectID `bson:"organization_id" json:"organization_id" binding:"required"`
	PONumber       string             `bson:"po_number" json:"po_number" binding:"required"`
	Status         POStatus           `bson:"status" json:"status"`

	// Supplier
	SupplierID primitive.ObjectID `bson:"supplier_id" json:"supplier_id" binding:"required"`

	// Delivery
	DeliveryLocationID primitive.ObjectID `bson:"delivery_location_id" json:"delivery_location_id"`
	DeliveryAddress    Address            `bson:"delivery_address" json:"delivery_address"`

	// Items
	Items []PurchaseOrderItem `bson:"items" json:"items" binding:"required"`

	// Dates
	OrderDate    time.Time  `bson:"order_date" json:"order_date"`
	ExpectedDate *time.Time `bson:"expected_date" json:"expected_date"`
	ReceivedDate *time.Time `bson:"received_date" json:"received_date"`

	// Amounts
	Subtotal       float64 `bson:"subtotal" json:"subtotal"`
	TaxAmount      float64 `bson:"tax_amount" json:"tax_amount"`
	DiscountAmount float64 `bson:"discount_amount" json:"discount_amount"`
	ShippingCost   float64 `bson:"shipping_cost" json:"shipping_cost"`
	TotalAmount    float64 `bson:"total_amount" json:"total_amount"`
	Currency       string  `bson:"currency" json:"currency"`

	// Payment
	PaymentTerms  int    `bson:"payment_terms" json:"payment_terms"`
	PaymentMethod string `bson:"payment_method" json:"payment_method"`
	Terms         string `bson:"terms" json:"terms"`

	// Tracking
	ReferenceNumber string              `bson:"reference_number" json:"reference_number"`
	RequestedBy     primitive.ObjectID  `bson:"requested_by" json:"requested_by"`
	ApprovedBy      *primitive.ObjectID `bson:"approved_by" json:"approved_by"`
	ApprovedDate    *time.Time          `bson:"approved_date" json:"approved_date"`

	Notes    string                 `bson:"notes" json:"notes"`
	Metadata map[string]interface{} `bson:"metadata" json:"metadata"`
	Tags     []string               `bson:"tags" json:"tags"`

	// Enriched Fields
	SupplierName string `bson:"-" json:"supplier_name,omitempty"`
}

type PurchaseOrderItem struct {
	ProductID        primitive.ObjectID `bson:"product_id" json:"product_id" binding:"required"`
	ProductName      string             `bson:"-" json:"product_name,omitempty"`
	SKU              string             `bson:"sku" json:"sku"`
	Description      string             `bson:"description" json:"description"`
	Quantity         float64            `bson:"quantity" json:"quantity" binding:"required,gt=0"`
	QuantityReceived float64            `bson:"quantity_received" json:"quantity_received"`
	UnitPrice        float64            `bson:"unit_price" json:"unit_price" binding:"required"`
	TaxRate          float64            `bson:"tax_rate" json:"tax_rate"`
	DiscountPercent  float64            `bson:"discount_percent" json:"discount_percent"`
	LineTotal        float64            `bson:"line_total" json:"line_total"`
}

type GoodsReceiptNote struct {
	BaseModel `bson:",inline"`

	OrganizationID primitive.ObjectID `bson:"organization_id" json:"organization_id" binding:"required"`
	GRNNumber      string             `bson:"grn_number" json:"grn_number" binding:"required"`
	Status         GRNStatus          `bson:"status" json:"status"`

	// Reference
	PurchaseOrderID primitive.ObjectID `bson:"purchase_order_id" json:"purchase_order_id" binding:"required"`
	PONumber        string             `bson:"po_number" json:"po_number"`
	SupplierID      primitive.ObjectID `bson:"supplier_id" json:"supplier_id"`

	// Receiving
	LocationID  primitive.ObjectID `bson:"location_id" json:"location_id" binding:"required"`
	ReceiptDate time.Time          `bson:"receipt_date" json:"receipt_date"`
	ReceivedBy  primitive.ObjectID `bson:"received_by" json:"received_by"`

	// Items
	Items []GRNItem `bson:"items" json:"items" binding:"required"`

	// Inspection
	InspectedBy   *primitive.ObjectID `bson:"inspected_by" json:"inspected_by"`
	InspectedDate *time.Time          `bson:"inspected_date" json:"inspected_date"`
	QCStatus      string              `bson:"qc_status" json:"qc_status"`
	QCNotes       string              `bson:"qc_notes" json:"qc_notes"`

	// Document
	InvoiceNumber string `bson:"invoice_number" json:"invoice_number"`
	DeliveryNote  string `bson:"delivery_note" json:"delivery_note"`

	Notes    string                 `bson:"notes" json:"notes"`
	Metadata map[string]interface{} `bson:"metadata" json:"metadata"`

	// Enriched Fields
	ReceivedByName  string `bson:"-" json:"received_by_name,omitempty"`
	InspectedByName string `bson:"-" json:"inspected_by_name,omitempty"`
	SupplierName    string `bson:"-" json:"supplier_name,omitempty"`
	LocationName    string `bson:"-" json:"location_name,omitempty"`

	// Calculated/Enriched for List View
	TotalValue       float64 `bson:"-" json:"total_value,omitempty"`
	POUnitName       string  `bson:"-" json:"po_unit_name,omitempty"`
	OrderedUnitName  string  `bson:"-" json:"ordered_unit_name,omitempty"`
	ReceivedUnitName string  `bson:"-" json:"received_unit_name,omitempty"`
}

type GRNItem struct {
	ProductID        primitive.ObjectID `bson:"product_id" json:"product_id" binding:"required"`
	ProductName      string             `bson:"-" json:"product_name,omitempty"`
	SKU              string             `bson:"sku" json:"sku"`
	Description      string             `bson:"description" json:"description"`
	OrderedQuantity  float64            `bson:"ordered_quantity" json:"ordered_quantity"`
	ReceivedQuantity float64            `bson:"received_quantity" json:"received_quantity" binding:"required,gt=0"`
	AcceptedQuantity float64            `bson:"accepted_quantity" json:"accepted_quantity"`
	RejectedQuantity float64            `bson:"rejected_quantity" json:"rejected_quantity"`
	UnitCost         float64            `bson:"unit_cost" json:"unit_cost"`
	BatchNumber      string             `bson:"batch_number" json:"batch_number"`
	ExpiryDate       *time.Time         `bson:"expiry_date" json:"expiry_date"`
	Condition        string             `bson:"condition" json:"condition"`
	RejectionReason  string             `bson:"rejection_reason" json:"rejection_reason"`
}
