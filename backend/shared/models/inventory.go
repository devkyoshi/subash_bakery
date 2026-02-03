package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ProductType defines the type of product
type ProductType string

const (
	ProductTypeRaw          ProductType = "raw_material"
	ProductTypeFinished     ProductType = "finished_goods"
	ProductTypeSemiFinished ProductType = "semi_finished"
	ProductTypeConsumable   ProductType = "consumable"
	ProductTypeService      ProductType = "service"
)

// ProductStatus defines product status
type ProductStatus string

const (
	ProductStatusActive       ProductStatus = "active"
	ProductStatusInactive     ProductStatus = "inactive"
	ProductStatusDiscontinued ProductStatus = "discontinued"
)

// StockValuationMethod defines how stock is valued
type StockValuationMethod string

const (
	ValuationFIFO            StockValuationMethod = "fifo"
	ValuationLIFO            StockValuationMethod = "lifo"
	ValuationWeightedAverage StockValuationMethod = "weighted_average"
	ValuationStandard        StockValuationMethod = "standard"
)

// MovementType defines types of stock movements
type MovementType string

const (
	MovementIn         MovementType = "in"         // Stock received
	MovementOut        MovementType = "out"        // Stock issued/sold
	MovementTransfer   MovementType = "transfer"   // Transfer between locations
	MovementAdjustment MovementType = "adjustment" // Inventory adjustment
	MovementReturn     MovementType = "return"     // Customer/supplier return
	MovementScrap      MovementType = "scrap"      // Scrapped/damaged
)

// ReorderPolicy defines reorder rules
type ReorderPolicy struct {
	Enabled bool `bson:"enabled" json:"enabled"`

	ReOrderLevel    int `bson:"re_order_level" json:"re_order_level"`
	ReOrderQuantity int `bson:"re_order_quantity" json:"re_order_quantity"`
	MinStockLevel   int `bson:"min_stock_level" json:"min_stock_level"`
	MaxStockLevel   int `bson:"max_stock_level" json:"max_stock_level"`
	SafetyStock     int `bson:"safety_stock" json:"safety_stock"`
	LeadTimeDays    int `bson:"lead_time_days" json:"lead_time_days"`

	// Audit Fields
	CreatedAt  int64 `bson:"created_at" json:"created_at"`
	ModifiedAt int64 `bson:"modified_at" json:"modified_at"`

	// Deletion Fields
	IsDeleted   bool  `bson:"is_deleted" json:"is_deleted"`
	DeletedDate int64 `bson:"deleted_date" json:"deleted_date"`
}

// Brand represents a product brand
type Brand struct {
	BaseModel `bson:",inline"`

	OrganizationID primitive.ObjectID `bson:"organization_id" json:"organization_id" binding:"required"`

	Name        string `bson:"name" json:"name" binding:"required"`
	Code        string `bson:"code" json:"code"`
	Description string `bson:"description" json:"description"`

	LogoURL string `bson:"logo_url" json:"logo_url"`
	Website string `bson:"website" json:"website"`
	Country string `bson:"country" json:"country"`

	IsActive bool `bson:"is_active" json:"is_active"`

	Metadata map[string]interface{} `bson:"metadata" json:"metadata"`
}

// Product represents an inventory item
type Product struct {
	BaseModel `bson:",inline"`

	// Basic Information
	OrganizationID primitive.ObjectID `bson:"organization_id" json:"organization_id" binding:"required"`

	//Identifications
	SKU         string `bson:"sku" json:"sku" binding:"required"` // Stock Keeping Unit
	Barcode     string `bson:"barcode" json:"barcode"`
	Name        string `bson:"name" json:"name" binding:"required"`
	Description string `bson:"description" json:"description"`

	Type   ProductType   `bson:"type" json:"type"`
	Status ProductStatus `bson:"status" json:"status"`

	// Classification
	CategoryID     primitive.ObjectID `bson:"category_id" json:"category_id"`
	SubcategoryID  primitive.ObjectID `bson:"subcategory_id" json:"subcategory_id"`
	BrandID        primitive.ObjectID `bson:"brand_id" json:"brand_id"`
	ManufacturerID primitive.ObjectID `bson:"manufacturer_id" json:"manufacturer_id"`

	// Inventory Settings
	TrackInventory     bool                 `bson:"track_inventory" json:"track_inventory"`
	TrackBatches       bool                 `bson:"track_batches" json:"track_batches"`
	TrackSerialNumbers bool                 `bson:"track_serial_numbers" json:"track_serial_numbers"`
	ValuationMethod    StockValuationMethod `bson:"valuation_method" json:"valuation_method"`

	// Unit of Measure
	BaseUnitID     primitive.ObjectID   `bson:"base_unit_id" json:"base_unit_id"`
	AllowedUnitIDs []primitive.ObjectID `bson:"allowed_unit_ids" json:"allowed_unit_ids"`

	// Dimensions & Weight
	Weight        float64 `bson:"weight" json:"weight"`
	WeightUnit    string  `bson:"weight_unit" json:"weight_unit"`
	Length        float64 `bson:"length" json:"length"`
	Width         float64 `bson:"width" json:"width"`
	Height        float64 `bson:"height" json:"height"`
	DimensionUnit string  `bson:"dimension_unit" json:"dimension_unit"`
	Volume        float64 `bson:"volume" json:"volume"`
	VolumeUnit    string  `bson:"volume_unit" json:"volume_unit"`

	// Pricing - Location-wise
	LocationPrices []LocationPrice `bson:"location_prices" json:"location_prices"`

	// Tax & Accounting
	TaxCategoryID primitive.ObjectID `bson:"tax_category_id" json:"tax_category_id"`
	HSNCode       string             `bson:"hsn_code" json:"hsn_code"` // Harmonized System Nomenclature
	SACCode       string             `bson:"sac_code" json:"sac_code"` // Service Accounting Code

	// Reorder Settings
	ReorderLevel    int `bson:"reorder_level" json:"reorder_level"`
	ReorderQuantity int `bson:"reorder_quantity" json:"reorder_quantity"`
	MinStockLevel   int `bson:"min_stock_level" json:"min_stock_level"`
	MaxStockLevel   int `bson:"max_stock_level" json:"max_stock_level"`
	SafetyStock     int `bson:"safety_stock" json:"safety_stock"`

	// Supplier Info
	DefaultSupplierID primitive.ObjectID   `bson:"default_supplier_id" json:"default_supplier_id"`
	SupplierIDs       []primitive.ObjectID `bson:"supplier_ids" json:"supplier_ids"`
	LeadTimeDays      int                  `bson:"lead_time_days" json:"lead_time_days"`

	// Quality & Expiry
	ShelfLifeDays int  `bson:"shelf_life_days" json:"shelf_life_days"`
	RequiresQC    bool `bson:"requires_qc" json:"requires_qc"` // Quality Check
	Perishable    bool `bson:"perishable" json:"perishable"`
	Hazardous     bool `bson:"hazardous" json:"hazardous"`

	// Images & Attachments
	Images         []string          `bson:"images" json:"images"`
	Thumbnail      string            `bson:"thumbnail" json:"thumbnail"`
	Specifications map[string]string `bson:"specifications" json:"specifications"`
	Attachments    []Attachment      `bson:"attachments" json:"attachments"`

	// Current Stock Summary (denormalized for quick access)
	TotalStock     float64 `bson:"total_stock" json:"total_stock"`
	AvailableStock float64 `bson:"available_stock" json:"available_stock"`
	AllocatedStock float64 `bson:"allocated_stock" json:"allocated_stock"`
	InTransitStock float64 `bson:"in_transit_stock" json:"in_transit_stock"`
	StockValue     float64 `bson:"stock_value" json:"stock_value"`

	// Analytics
	TotalSold        float64    `bson:"total_sold" json:"total_sold"`
	TotalPurchased   float64    `bson:"total_purchased" json:"total_purchased"`
	LastSoldDate     *time.Time `bson:"last_sold_date" json:"last_sold_date"`
	LastPurchaseDate *time.Time `bson:"last_purchase_date" json:"last_purchase_date"`

	// Metadata
	Metadata map[string]interface{} `bson:"metadata" json:"metadata"`
	Tags     []string               `bson:"tags" json:"tags"`
}

// UnitOfMeasure represents alternative units and their conversion rates
type UnitOfMeasure struct {
	UOM            string  `bson:"uom" json:"uom"`                         // e.g., "BOX"
	ConversionRate float64 `bson:"conversion_rate" json:"conversion_rate"` // e.g., 1 BOX = 12 PCS
	Description    string  `bson:"description" json:"description"`
}

// Attachment represents a file attachment
type Attachment struct {
	Name string `bson:"name" json:"name"`
	URL  string `bson:"url" json:"url"`
	Type string `bson:"type" json:"type"` // e.g., "pdf", "image"
	Size int64  `bson:"size" json:"size"`
}

// LocationPrice represents pricing for a specific location
type LocationPrice struct {
	LocationID     primitive.ObjectID `bson:"location_id" json:"location_id" binding:"required"`
	LocationName   string             `bson:"location_name" json:"location_name"`
	PurchaseUnitID primitive.ObjectID `bson:"purchase_unit_id" json:"purchase_unit_id"`
	PurchaseUnit   *Unit              `bson:"-" json:"purchase_unit,omitempty"`
	SellingUnitID  primitive.ObjectID `bson:"selling_unit_id" json:"selling_unit_id"`
	SellingUnit    *Unit              `bson:"-" json:"selling_unit,omitempty"`
	CostPrice      float64            `bson:"cost_price" json:"cost_price"`
	SellingPrice   float64            `bson:"selling_price" json:"selling_price"`
	MRP            float64            `bson:"mrp" json:"mrp"` // Maximum Retail Price
	InitialStock   float64            `bson:"initial_stock" json:"initial_stock"`
	CurrentStock   float64            `bson:"-" json:"current_stock"`
	AvailableStock float64            `bson:"-" json:"available_stock"`
	AllocatedStock float64            `bson:"-" json:"allocated_stock"`
	Currency       string             `bson:"currency" json:"currency"`
	IsActive       bool               `bson:"is_active" json:"is_active"`
	CreatedAt      int64              `bson:"created_at" json:"created_at"`
	ModifiedAt     int64              `bson:"modified_at" json:"modified_at"`
}

// StockLevel represents stock at a specific location
type StockLevel struct {
	BaseModel `bson:",inline"`

	OrganizationID primitive.ObjectID `bson:"organization_id" json:"organization_id" binding:"required"`
	ProductID      primitive.ObjectID `bson:"product_id" json:"product_id" binding:"required"`
	LocationID     primitive.ObjectID `bson:"location_id" json:"location_id" binding:"required"`
	WarehouseZone  string             `bson:"warehouse_zone" json:"warehouse_zone"` // e.g., "A-01-05"

	// Quantities
	QuantityOnHand    float64 `bson:"quantity_on_hand" json:"quantity_on_hand"`
	QuantityAvailable float64 `bson:"quantity_available" json:"quantity_available"` // On hand - allocated
	QuantityAllocated float64 `bson:"quantity_allocated" json:"quantity_allocated"` // Reserved for orders
	QuantityInTransit float64 `bson:"quantity_in_transit" json:"quantity_in_transit"`
	QuantityReserved  float64 `bson:"quantity_reserved" json:"quantity_reserved"` // Reserved for production

	// Valuation
	AverageCost float64 `bson:"average_cost" json:"average_cost"`
	LastCost    float64 `bson:"last_cost" json:"last_cost"`
	TotalValue  float64 `bson:"total_value" json:"total_value"`

	// Tracking
	LastMovementDate *time.Time `bson:"last_movement_date" json:"last_movement_date"`
	LastCountDate    *time.Time `bson:"last_count_date" json:"last_count_date"`

	Metadata map[string]interface{} `bson:"metadata" json:"metadata"`
}

// StockMovement records all stock transactions
type StockMovement struct {
	BaseModel `bson:",inline"`

	OrganizationID primitive.ObjectID `bson:"organization_id" json:"organization_id" binding:"required"`
	ProductID      primitive.ObjectID `bson:"product_id" json:"product_id" binding:"required"`
	MovementType   MovementType       `bson:"movement_type" json:"movement_type" binding:"required"`

	// Locations
	FromLocationID *primitive.ObjectID `bson:"from_location_id" json:"from_location_id"`
	ToLocationID   *primitive.ObjectID `bson:"to_location_id" json:"to_location_id"`

	// Quantity & Cost
	Quantity  float64 `bson:"quantity" json:"quantity" binding:"required"`
	UOM       string  `bson:"uom" json:"uom"`
	UnitCost  float64 `bson:"unit_cost" json:"unit_cost"`
	TotalCost float64 `bson:"total_cost" json:"total_cost"`

	// Reference
	ReferenceType string             `bson:"reference_type" json:"reference_type"` // e.g., "purchase_order", "sales_order"
	ReferenceID   primitive.ObjectID `bson:"reference_id" json:"reference_id"`
	ReferenceNo   string             `bson:"reference_no" json:"reference_no"`

	// Batch & Serial
	BatchID       *primitive.ObjectID `bson:"batch_id" json:"batch_id"`
	SerialNumbers []string            `bson:"serial_numbers" json:"serial_numbers"`

	// Details
	Reason       string             `bson:"reason" json:"reason"`
	Notes        string             `bson:"notes" json:"notes"`
	MovementDate time.Time          `bson:"movement_date" json:"movement_date"`
	IsReversed   bool               `bson:"is_reversed" json:"is_reversed"`
	ReversedBy   primitive.ObjectID `bson:"reversed_by" json:"reversed_by"`
	ReversedAt   *time.Time         `bson:"reversed_at" json:"reversed_at"`

	Metadata map[string]interface{} `bson:"metadata" json:"metadata"`
}

// Batch represents a batch of products
type Batch struct {
	BaseModel `bson:",inline"`

	OrganizationID primitive.ObjectID `bson:"organization_id" json:"organization_id" binding:"required"`
	ProductID      primitive.ObjectID `bson:"product_id" json:"product_id" binding:"required"`
	LocationID     primitive.ObjectID `bson:"location_id" json:"location_id" binding:"required"`

	// Batch Information
	BatchNumber     string     `bson:"batch_number" json:"batch_number" binding:"required"`
	ManufactureDate time.Time  `bson:"manufacture_date" json:"manufacture_date"`
	ExpiryDate      *time.Time `bson:"expiry_date" json:"expiry_date"`
	ReceiveDate     time.Time  `bson:"receive_date" json:"receive_date"`

	// Quantity & Cost
	InitialQuantity   float64 `bson:"initial_quantity" json:"initial_quantity"`
	CurrentQuantity   float64 `bson:"current_quantity" json:"current_quantity"`
	AllocatedQuantity float64 `bson:"allocated_quantity" json:"allocated_quantity"`
	UnitCost          float64 `bson:"unit_cost" json:"unit_cost"`
	TotalCost         float64 `bson:"total_cost" json:"total_cost"`

	// Supplier
	SupplierID      primitive.ObjectID `bson:"supplier_id" json:"supplier_id"`
	PurchaseOrderID primitive.ObjectID `bson:"purchase_order_id" json:"purchase_order_id"`

	// Quality
	QCStatus string     `bson:"qc_status" json:"qc_status"` // pending, passed, failed
	QCDate   *time.Time `bson:"qc_date" json:"qc_date"`
	QCNotes  string     `bson:"qc_notes" json:"qc_notes"`

	// Status
	IsActive   bool `bson:"is_active" json:"is_active"`
	IsExpired  bool `bson:"is_expired" json:"is_expired"`
	IsDepleted bool `bson:"is_depleted" json:"is_depleted"`

	Metadata map[string]interface{} `bson:"metadata" json:"metadata"`
}

// SerialNumber represents individual serialized items
type SerialNumber struct {
	BaseModel `bson:",inline"`

	OrganizationID primitive.ObjectID  `bson:"organization_id" json:"organization_id" binding:"required"`
	ProductID      primitive.ObjectID  `bson:"product_id" json:"product_id" binding:"required"`
	LocationID     *primitive.ObjectID `bson:"location_id" json:"location_id"`

	// Serial Information
	SerialNo        string     `bson:"serial_no" json:"serial_no" binding:"required"`
	ManufactureDate *time.Time `bson:"manufacture_date" json:"manufacture_date"`
	WarrantyExpiry  *time.Time `bson:"warranty_expiry" json:"warranty_expiry"`

	// Batch
	BatchID *primitive.ObjectID `bson:"batch_id" json:"batch_id"`

	// Status
	Status      string `bson:"status" json:"status"` // available, allocated, sold, returned, scrapped
	IsAvailable bool   `bson:"is_available" json:"is_available"`

	// Tracking
	CustomerID   *primitive.ObjectID `bson:"customer_id" json:"customer_id"`
	SalesOrderID *primitive.ObjectID `bson:"sales_order_id" json:"sales_order_id"`
	SoldDate     *time.Time          `bson:"sold_date" json:"sold_date"`

	// Cost
	UnitCost float64 `bson:"unit_cost" json:"unit_cost"`

	Metadata map[string]interface{} `bson:"metadata" json:"metadata"`
}

// ProductSubcategory represents a child classification within a category.
type ProductSubcategory struct {
	BaseModel `bson:",inline"`

	Name         string                 `bson:"name" json:"name" binding:"required"`
	Code         string                 `bson:"code" json:"code"`
	Description  string                 `bson:"description" json:"description"`
	IsActive     bool                   `bson:"is_active" json:"is_active"`
	ProductCount int                    `bson:"product_count" json:"product_count"`
	Metadata     map[string]interface{} `bson:"metadata" json:"metadata"`
}

// ProductCategory represents product categories with embedded subcategories (two-level hierarchy).
type ProductCategory struct {
	BaseModel `bson:",inline"`

	OrganizationID primitive.ObjectID `bson:"organization_id" json:"organization_id" binding:"required"`

	Name        string `bson:"name" json:"name" binding:"required"`
	Code        string `bson:"code" json:"code"`
	Description string `bson:"description" json:"description"`

	// Embedded subcategories (max depth = 1)
	Subcategories []ProductSubcategory `bson:"subcategories" json:"subcategories"`

	// Inventory Defaults (Applied if product doesn’t override)
	DefaultReorderPolicy *ReorderPolicy `bson:"default_reorder_policy" json:"default_reorder_policy"`

	// Status
	IsActive     bool `bson:"is_active" json:"is_active"`
	ProductCount int  `bson:"product_count" json:"product_count"`

	Metadata map[string]interface{} `bson:"metadata" json:"metadata"`
}

// StockAdjustment represents manual stock adjustments
type StockAdjustment struct {
	BaseModel `bson:",inline"`

	OrganizationID primitive.ObjectID `bson:"organization_id" json:"organization_id" binding:"required"`
	LocationID     primitive.ObjectID `bson:"location_id" json:"location_id" binding:"required"`
	AdjustmentNo   string             `bson:"adjustment_no" json:"adjustment_no"`
	AdjustmentDate time.Time          `bson:"adjustment_date" json:"adjustment_date"`

	// Reason
	Reason        string `bson:"reason" json:"reason"` // damaged, theft, count_correction, etc.
	ReasonDetails string `bson:"reason_details" json:"reason_details"`

	// Items
	Items []StockAdjustmentItem `bson:"items" json:"items"`

	// Approval
	Status         string             `bson:"status" json:"status"` // draft, pending, approved, rejected
	ApprovedBy     primitive.ObjectID `bson:"approved_by" json:"approved_by"`
	ApprovedAt     *time.Time         `bson:"approved_at" json:"approved_at"`
	RejectedReason string             `bson:"rejected_reason" json:"rejected_reason"`

	Notes    string                 `bson:"notes" json:"notes"`
	Metadata map[string]interface{} `bson:"metadata" json:"metadata"`
}

// StockAdjustmentItem represents an item in stock adjustment
type StockAdjustmentItem struct {
	ProductID     primitive.ObjectID  `bson:"product_id" json:"product_id"`
	ExpectedQty   float64             `bson:"expected_qty" json:"expected_qty"`
	ActualQty     float64             `bson:"actual_qty" json:"actual_qty"`
	DifferenceQty float64             `bson:"difference_qty" json:"difference_qty"`
	UOM           string              `bson:"uom" json:"uom"`
	UnitCost      float64             `bson:"unit_cost" json:"unit_cost"`
	TotalCost     float64             `bson:"total_cost" json:"total_cost"`
	BatchID       *primitive.ObjectID `bson:"batch_id" json:"batch_id"`
	Reason        string              `bson:"reason" json:"reason"`
}

// InventoryCount represents physical stock count
type InventoryCount struct {
	BaseModel `bson:",inline"`

	OrganizationID primitive.ObjectID `bson:"organization_id" json:"organization_id" binding:"required"`
	LocationID     primitive.ObjectID `bson:"location_id" json:"location_id" binding:"required"`
	CountNo        string             `bson:"count_no" json:"count_no"`
	CountDate      time.Time          `bson:"count_date" json:"count_date"`
	CountType      string             `bson:"count_type" json:"count_type"` // full, cycle, spot

	// Items
	Items []InventoryCountItem `bson:"items" json:"items"`

	// Status
	Status      string             `bson:"status" json:"status"` // in_progress, completed, cancelled
	StartedBy   primitive.ObjectID `bson:"started_by" json:"started_by"`
	StartedAt   time.Time          `bson:"started_at" json:"started_at"`
	CompletedBy primitive.ObjectID `bson:"completed_by" json:"completed_by"`
	CompletedAt *time.Time         `bson:"completed_at" json:"completed_at"`

	// Summary
	TotalItemsCounted int     `bson:"total_items_counted" json:"total_items_counted"`
	TotalVariance     float64 `bson:"total_variance" json:"total_variance"`
	VarianceValue     float64 `bson:"variance_value" json:"variance_value"`

	Notes    string                 `bson:"notes" json:"notes"`
	Metadata map[string]interface{} `bson:"metadata" json:"metadata"`
}

// InventoryCountItem represents an item in inventory count
type InventoryCountItem struct {
	ProductID     primitive.ObjectID  `bson:"product_id" json:"product_id"`
	SystemQty     float64             `bson:"system_qty" json:"system_qty"`
	CountedQty    float64             `bson:"counted_qty" json:"counted_qty"`
	VarianceQty   float64             `bson:"variance_qty" json:"variance_qty"`
	UOM           string              `bson:"uom" json:"uom"`
	UnitCost      float64             `bson:"unit_cost" json:"unit_cost"`
	VarianceValue float64             `bson:"variance_value" json:"variance_value"`
	BatchID       *primitive.ObjectID `bson:"batch_id" json:"batch_id"`
	CountedBy     primitive.ObjectID  `bson:"counted_by" json:"counted_by"`
	CountedAt     time.Time           `bson:"counted_at" json:"counted_at"`
	Notes         string              `bson:"notes" json:"notes"`
}
