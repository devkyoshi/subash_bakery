package service

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/yourusername/erp-system/services/procurement-service/internal/client"
	"github.com/yourusername/erp-system/services/procurement-service/internal/repository"
	"github.com/yourusername/erp-system/shared/models"
)

type ProcurementService struct {
	supplierRepo    *repository.SupplierRepository
	poRepo          *repository.PurchaseOrderRepository
	grnRepo         *repository.GRNRepository
	productClient   *client.ProductClient
	userClient      *client.UserClient
	inventoryClient *client.InventoryClient
}

func NewProcurementService(
	supplierRepo *repository.SupplierRepository,
	poRepo *repository.PurchaseOrderRepository,
	grnRepo *repository.GRNRepository,
	productClient *client.ProductClient,
	userClient *client.UserClient,
	inventoryClient *client.InventoryClient,
) *ProcurementService {
	return &ProcurementService{
		supplierRepo:    supplierRepo,
		poRepo:          poRepo,
		grnRepo:         grnRepo,
		productClient:   productClient,
		userClient:      userClient,
		inventoryClient: inventoryClient,
	}
}

// ============== Supplier Operations ==============

type CreateSupplierRequest struct {
	CompanyName   string          `json:"company_name" binding:"required"`
	ContactPerson string          `json:"contact_person" binding:"required"`
	Email         string          `json:"email" binding:"required,email"`
	Phone         string          `json:"phone,omitempty"`
	Mobile        string          `json:"mobile,omitempty"`
	Website       string          `json:"website,omitempty"`
	Address       *models.Address `json:"address,omitempty"`
	TaxID         string          `json:"tax_id,omitempty"`
	PaymentTerms  int             `json:"payment_terms,omitempty"`
	CreditLimit   float64         `json:"credit_limit,omitempty"`
	BankName      string          `json:"bank_name,omitempty"`
	AccountNumber string          `json:"account_number,omitempty"`
	SwiftCode     string          `json:"swift_code,omitempty"`
	Tags          []string        `json:"tags,omitempty"`
	Notes         string          `json:"notes,omitempty"`
}

func (s *ProcurementService) CreateSupplier(ctx context.Context, orgID primitive.ObjectID, req CreateSupplierRequest, createdBy primitive.ObjectID) (*models.Supplier, error) {
	// Generate supplier code
	supplierCode := s.generateSupplierCode(orgID)

	// Check if code exists
	exists, err := s.supplierRepo.CodeExists(ctx, orgID, supplierCode, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to check supplier code existence: %w", err)
	}
	if exists {
		supplierCode = fmt.Sprintf("%s-%d", supplierCode, time.Now().Unix())
	}

	// Create supplier
	supplier := &models.Supplier{
		BaseModel: models.BaseModel{
			CreatedBy: createdBy,
		},
		OrganizationID: orgID,
		SupplierCode:   supplierCode,
		CompanyName:    req.CompanyName,
		Status:         models.SupplierStatusActive,
		ContactPerson:  req.ContactPerson,
		Email:          req.Email,
		Phone:          req.Phone,
		Mobile:         req.Mobile,
		Website:        req.Website,
		TaxID:          req.TaxID,
		PaymentTerms:   req.PaymentTerms,
		CreditLimit:    req.CreditLimit,
		BankName:       req.BankName,
		AccountNumber:  req.AccountNumber,
		SwiftCode:      req.SwiftCode,
		Tags:           req.Tags,
		Notes:          req.Notes,
	}

	if req.Address != nil {
		supplier.Address = *req.Address
	}

	if err := s.supplierRepo.Create(ctx, supplier); err != nil {
		return nil, fmt.Errorf("failed to create supplier: %w", err)
	}

	return supplier, nil
}

func (s *ProcurementService) GetSupplier(ctx context.Context, id primitive.ObjectID) (*models.Supplier, error) {
	return s.supplierRepo.FindByID(ctx, id)
}

func (s *ProcurementService) ListSuppliers(ctx context.Context, orgID primitive.ObjectID, filters map[string]interface{}, page, limit int) ([]*models.Supplier, int64, error) {
	return s.supplierRepo.FindAll(ctx, orgID, filters, page, limit)
}

func (s *ProcurementService) UpdateSupplier(ctx context.Context, id primitive.ObjectID, updates map[string]interface{}) (*models.Supplier, error) {
	supplier, err := s.supplierRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Apply updates
	if companyName, ok := updates["company_name"].(string); ok {
		supplier.CompanyName = companyName
	}
	if contactPerson, ok := updates["contact_person"].(string); ok {
		supplier.ContactPerson = contactPerson
	}
	if email, ok := updates["email"].(string); ok {
		supplier.Email = email
	}
	if status, ok := updates["status"].(string); ok {
		supplier.Status = models.SupplierStatus(status)
	}
	if creditLimit, ok := updates["credit_limit"].(float64); ok {
		supplier.CreditLimit = creditLimit
	}

	if err := s.supplierRepo.Update(ctx, supplier); err != nil {
		return nil, err
	}

	return supplier, nil
}

func (s *ProcurementService) DeleteSupplier(ctx context.Context, id, deletedBy primitive.ObjectID) error {
	return s.supplierRepo.SoftDelete(ctx, id, deletedBy)
}

// ============== Purchase Order Operations ==============

type CreatePurchaseOrderRequest struct {
	SupplierID      primitive.ObjectID         `json:"supplier_id" binding:"required"`
	OrderDate       time.Time                  `json:"order_date" binding:"required"`
	ExpectedDate    *time.Time                 `json:"expected_date,omitempty"`
	Items           []models.PurchaseOrderItem `json:"items" binding:"required,min=1"`
	ShippingAddress *models.Address            `json:"shipping_address,omitempty"`
	Notes           string                     `json:"notes,omitempty"`
	Terms           string                     `json:"terms,omitempty"`
	ReferenceNumber string                     `json:"reference_number,omitempty"`
	TaxRate         float64                    `json:"tax_rate,omitempty"`
}

func (s *ProcurementService) CreatePurchaseOrder(ctx context.Context, orgID primitive.ObjectID, req CreatePurchaseOrderRequest, createdBy primitive.ObjectID) (*models.PurchaseOrder, error) {
	// Validate items
	if len(req.Items) == 0 {
		return nil, fmt.Errorf("purchase order must have at least one item")
	}

	// Verify supplier exists
	_, err := s.supplierRepo.FindByID(ctx, req.SupplierID)
	if err != nil {
		return nil, fmt.Errorf("supplier not found: %w", err)
	}

	// Generate PO number
	poNumber := s.generatePONumber(orgID)

	// Check if PO number exists
	exists, err := s.poRepo.PONumberExists(ctx, orgID, poNumber, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to check PO number existence: %w", err)
	}
	if exists {
		poNumber = fmt.Sprintf("%s-%d", poNumber, time.Now().Unix())
	}

	// Calculate totals
	subtotal := 0.0
	for i := range req.Items {
		req.Items[i].LineTotal = req.Items[i].Quantity * req.Items[i].UnitPrice
		subtotal += req.Items[i].LineTotal
	}

	taxAmount := subtotal * req.TaxRate / 100.0
	totalAmount := subtotal + taxAmount

	// Create purchase order
	po := &models.PurchaseOrder{
		BaseModel: models.BaseModel{
			CreatedBy: createdBy,
		},
		OrganizationID:  orgID,
		PONumber:        poNumber,
		SupplierID:      req.SupplierID,
		Status:          models.POStatusDraft,
		OrderDate:       req.OrderDate,
		ExpectedDate:    req.ExpectedDate,
		Items:           req.Items,
		Subtotal:        subtotal,
		TaxAmount:       taxAmount,
		TotalAmount:     totalAmount,
		Notes:           req.Notes,
		Terms:           req.Terms,
		ReferenceNumber: req.ReferenceNumber,
	}

	if req.ShippingAddress != nil {
		po.DeliveryAddress = *req.ShippingAddress
	}

	if err := s.poRepo.Create(ctx, po); err != nil {
		return nil, fmt.Errorf("failed to create purchase order: %w", err)
	}

	return po, nil
}

func (s *ProcurementService) GetPurchaseOrder(ctx context.Context, id primitive.ObjectID, token string) (*models.PurchaseOrder, error) {
	po, err := s.poRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Enrich with product names if missing
	productIDs := []primitive.ObjectID{}
	for _, item := range po.Items {
		productIDs = append(productIDs, item.ProductID)
	}

	if len(productIDs) > 0 {
		if products, err := s.productClient.GetProductsBatch(ctx, productIDs, token); err == nil {
			for i := range po.Items {
				if prod, ok := products[po.Items[i].ProductID.Hex()]; ok {
					if po.Items[i].Description == "" {
						po.Items[i].Description = prod.Name
					}
					// Also populate ProductName field for explicit requirement
					po.Items[i].ProductName = prod.Name
					// And SKU if missing
					if po.Items[i].SKU == "" {
						po.Items[i].SKU = prod.SKU
					}
				}
			}
		}
	}

	return po, nil
}

func (s *ProcurementService) ListPurchaseOrders(ctx context.Context, orgID primitive.ObjectID, filters map[string]interface{}, page, limit int) ([]*models.PurchaseOrder, int64, error) {
	orders, total, err := s.poRepo.FindAll(ctx, orgID, filters, page, limit)
	if err != nil {
		return nil, 0, err
	}

	// Collect all supplier IDs
	supplierIDSet := make(map[string]primitive.ObjectID)
	for _, order := range orders {
		if !order.SupplierID.IsZero() {
			supplierIDSet[order.SupplierID.Hex()] = order.SupplierID
		}
	}

	// Fetch suppliers and enrich
	for _, supplierID := range supplierIDSet {
		supplier, err := s.supplierRepo.FindByID(ctx, supplierID)
		if err == nil && supplier != nil {
			for _, order := range orders {
				if order.SupplierID == supplierID {
					order.SupplierName = supplier.CompanyName
				}
			}
		}
	}

	return orders, total, nil
}

func (s *ProcurementService) ApprovePurchaseOrder(ctx context.Context, id, approvedBy primitive.ObjectID) error {
	return s.poRepo.Approve(ctx, id, approvedBy)
}

func (s *ProcurementService) UpdatePOStatus(ctx context.Context, id primitive.ObjectID, status models.POStatus) error {
	return s.poRepo.UpdateStatus(ctx, id, status)
}

func (s *ProcurementService) DeletePurchaseOrder(ctx context.Context, id, deletedBy primitive.ObjectID) error {
	return s.poRepo.SoftDelete(ctx, id, deletedBy)
}

// ============== GRN Operations ==============

type CreateGRNRequest struct {
	PurchaseOrderID primitive.ObjectID `json:"purchase_order_id" binding:"required"`
	LocationID      primitive.ObjectID `json:"location_id" binding:"required"`
	ReceiptDate     time.Time          `json:"receipt_date" binding:"required"`
	Items           []models.GRNItem   `json:"items" binding:"required,min=1"`
	Notes           string             `json:"notes,omitempty"`
}

func (s *ProcurementService) CreateGRN(ctx context.Context, orgID primitive.ObjectID, req CreateGRNRequest, createdBy primitive.ObjectID) (*models.GoodsReceiptNote, error) {
	// Get the purchase order
	po, err := s.poRepo.FindByID(ctx, req.PurchaseOrderID)
	if err != nil {
		return nil, err
	}

	// Generate GRN number
	grnNumber := s.generateGRNNumber(orgID)

	// Check if GRN number exists
	exists, err := s.grnRepo.GRNNumberExists(ctx, orgID, grnNumber, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to check GRN number existence: %w", err)
	}
	if exists {
		grnNumber = fmt.Sprintf("%s-%d", grnNumber, time.Now().Unix())
	}

	// Create GRN
	grn := &models.GoodsReceiptNote{
		BaseModel: models.BaseModel{
			CreatedBy: createdBy,
		},
		OrganizationID:  orgID,
		GRNNumber:       grnNumber,
		PurchaseOrderID: req.PurchaseOrderID,
		PONumber:        po.PONumber,
		SupplierID:      po.SupplierID,
		LocationID:      req.LocationID,
		Status:          models.GRNStatusReceived,
		ReceiptDate:     req.ReceiptDate,
		ReceivedBy:      createdBy,
		Items:           req.Items,
		Notes:           req.Notes,
	}

	if err := s.grnRepo.Create(ctx, grn); err != nil {
		return nil, fmt.Errorf("failed to create GRN: %w", err)
	}

	// Update PO status
	allReceived := true
	for _, item := range req.Items {
		if item.ReceivedQuantity < item.OrderedQuantity {
			allReceived = false
			break
		}
	}

	if allReceived {
		s.poRepo.UpdateStatus(ctx, req.PurchaseOrderID, models.POStatusReceived)
	} else {
		s.poRepo.UpdateStatus(ctx, req.PurchaseOrderID, models.POStatusPartiallyReceived)
	}

	return grn, nil
}

func (s *ProcurementService) GetGRN(ctx context.Context, id primitive.ObjectID, token string) (*models.GoodsReceiptNote, error) {
	grn, err := s.grnRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Enrich Product Details
	productIDs := []primitive.ObjectID{}
	for _, item := range grn.Items {
		productIDs = append(productIDs, item.ProductID)
	}
	if len(productIDs) > 0 {
		if products, err := s.productClient.GetProductsBatch(ctx, productIDs, token); err == nil {
			for i := range grn.Items {
				if prod, ok := products[grn.Items[i].ProductID.Hex()]; ok {
					if grn.Items[i].Description == "" {
						grn.Items[i].Description = prod.Name
					}
					if grn.Items[i].SKU == "" {
						grn.Items[i].SKU = prod.SKU
					}
				}
			}
		}
	}

	// Enrich User Names
	userIDs := []primitive.ObjectID{grn.ReceivedBy}
	if grn.InspectedBy != nil {
		userIDs = append(userIDs, *grn.InspectedBy)
	}

	users, err := s.userClient.GetUsersBatch(ctx, userIDs, token)
	if err != nil {
		fmt.Printf("Warning: Failed to fetch users for GRN %s: %v\n", grn.GRNNumber, err)
	} else {
		if user, ok := users[grn.ReceivedBy.Hex()]; ok {
			name := user.FullName
			if name == "" {
				name = user.FirstName + " " + user.LastName
			}
			grn.ReceivedByName = name
		}
		if grn.InspectedBy != nil {
			if user, ok := users[grn.InspectedBy.Hex()]; ok {
				name := user.FullName
				if name == "" {
					name = user.FirstName + " " + user.LastName
				}
				grn.InspectedByName = name
			}
		}
	}

	return grn, nil
}

func (s *ProcurementService) ListGRNs(ctx context.Context, orgID primitive.ObjectID, filters map[string]interface{}, page, limit int, token string) ([]*models.GoodsReceiptNote, int64, error) {
	grns, total, err := s.grnRepo.FindAll(ctx, orgID, filters, page, limit)
	if err != nil {
		return nil, 0, err
	}

	// Collect all user IDs and product IDs
	userIDSet := make(map[string]primitive.ObjectID)
	productIDSet := make(map[string]primitive.ObjectID)

	for _, grn := range grns {
		userIDSet[grn.ReceivedBy.Hex()] = grn.ReceivedBy
		if grn.InspectedBy != nil {
			userIDSet[grn.InspectedBy.Hex()] = *grn.InspectedBy
		}
		for _, item := range grn.Items {
			productIDSet[item.ProductID.Hex()] = item.ProductID
		}
	}

	// Fetch users
	var users map[string]*models.User
	if len(userIDSet) > 0 {
		var userIDs []primitive.ObjectID
		for _, id := range userIDSet {
			userIDs = append(userIDs, id)
		}
		var err error
		users, err = s.userClient.GetUsersBatch(ctx, userIDs, token)
		if err != nil {
			fmt.Printf("Warning: Failed to fetch users for GRN list: %v\n", err)
		}
	}

	// Fetch products and collect unit IDs
	var products map[string]*models.Product
	unitIDSet := make(map[string]string)

	if len(productIDSet) > 0 {
		var productIDs []primitive.ObjectID
		for _, id := range productIDSet {
			productIDs = append(productIDs, id)
		}
		var err error
		products, err = s.productClient.GetProductsBatch(ctx, productIDs, token)
		if err != nil {
			fmt.Printf("Warning: Failed to fetch products for GRN list: %v\n", err)
		} else {
			for _, prod := range products {
				if !prod.BaseUnitID.IsZero() {
					unitIDSet[prod.BaseUnitID.Hex()] = prod.BaseUnitID.Hex()
				} else {
					fmt.Printf("Debug: Product %s has zero BaseUnitID\n", prod.ID.Hex())
				}
			}
		}
	}

	// Fetch units
	var units map[string]*models.Unit
	if len(unitIDSet) > 0 {
		var unitIDs []string
		for id := range unitIDSet {
			unitIDs = append(unitIDs, id)
		}
		fmt.Printf("Debug: Fetching units with IDs: %v\n", unitIDs)
		var err error
		units, err = s.inventoryClient.GetUnitsBatch(ctx, unitIDs, token)
		if err != nil {
			fmt.Printf("Warning: Failed to fetch units for GRN list: %v\n", err)
		} else {
			fmt.Printf("Debug: Fetched %d units\n", len(units))
		}
	} else {
		fmt.Printf("Debug: No unit IDs collected from products\n")
	}

	// Collect supplier IDs and PO IDs
	supplierIDSet := make(map[string]primitive.ObjectID)
	poIDSet := make(map[string]primitive.ObjectID)
	for _, grn := range grns {
		if !grn.SupplierID.IsZero() {
			supplierIDSet[grn.SupplierID.Hex()] = grn.SupplierID
		}
		if !grn.PurchaseOrderID.IsZero() {
			poIDSet[grn.PurchaseOrderID.Hex()] = grn.PurchaseOrderID
		}
	}

	// Fetch suppliers
	var suppliers map[string]*models.Supplier
	if len(supplierIDSet) > 0 {
		suppliers = make(map[string]*models.Supplier)
		for _, supplierID := range supplierIDSet {
			supplier, err := s.supplierRepo.FindByID(ctx, supplierID)
			if err == nil && supplier != nil {
				suppliers[supplierID.Hex()] = supplier
			}
		}
	}

	// Fetch Purchase Orders (to get UnitPrice for items)
	purchaseOrders := make(map[string]*models.PurchaseOrder)
	for _, poID := range poIDSet {
		po, err := s.poRepo.FindByID(ctx, poID)
		if err == nil && po != nil {
			purchaseOrders[poID.Hex()] = po
		}
	}

	// Enrich GRNs
	for _, grn := range grns {
		// Enrich Users
		if users != nil {
			if user, ok := users[grn.ReceivedBy.Hex()]; ok {
				grn.ReceivedByName = user.FullName
				if grn.ReceivedByName == "" {
					grn.ReceivedByName = user.FirstName + " " + user.LastName
				}
			}
			if grn.InspectedBy != nil {
				if user, ok := users[grn.InspectedBy.Hex()]; ok {
					grn.InspectedByName = user.FullName
					if grn.InspectedByName == "" {
						grn.InspectedByName = user.FirstName + " " + user.LastName
					}
				}
			}
		}

		// Enrich Supplier
		if suppliers != nil {
			if supplier, ok := suppliers[grn.SupplierID.Hex()]; ok {
				grn.SupplierName = supplier.CompanyName
			}
		}

		// Build PO item price map for this GRN's PO
		poItemPrices := make(map[string]float64)
		if po, ok := purchaseOrders[grn.PurchaseOrderID.Hex()]; ok {
			for _, poItem := range po.Items {
				poItemPrices[poItem.ProductID.Hex()] = poItem.UnitPrice
			}
		}

		// Enrich Products, Units, and Calculate Total
		var totalValue float64
		var unitName string

		for i := range grn.Items {
			// Fill UnitCost from PO if it is zero
			if grn.Items[i].UnitCost == 0 {
				if price, ok := poItemPrices[grn.Items[i].ProductID.Hex()]; ok {
					grn.Items[i].UnitCost = price
				}
			}

			// Total Value
			itemTotal := grn.Items[i].ReceivedQuantity * grn.Items[i].UnitCost
			totalValue += itemTotal

			// Product Name, SKU, Description and Unit
			if products != nil {
				if prod, ok := products[grn.Items[i].ProductID.Hex()]; ok {
					grn.Items[i].ProductName = prod.Name

					if grn.Items[i].SKU == "" {
						grn.Items[i].SKU = prod.SKU
					}
					if grn.Items[i].Description == "" {
						grn.Items[i].Description = prod.Description
					}

					// Unit
					if unitName == "" && units != nil && !prod.BaseUnitID.IsZero() {
						if unit, ok := units[prod.BaseUnitID.Hex()]; ok {
							unitName = unit.Name
						}
					}
				}
			}
		}

		grn.TotalValue = totalValue
		grn.POUnitName = unitName
		grn.OrderedUnitName = unitName
		grn.ReceivedUnitName = unitName
	}

	return grns, total, nil
}

func (s *ProcurementService) CompleteInspection(ctx context.Context, id, inspectedBy primitive.ObjectID, qcStatus, qcNotes, token string) error {
	// First, fetch the GRN to get location and item details
	grn, err := s.grnRepo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to fetch GRN: %w", err)
	}

	// Complete inspection in repository
	if err := s.grnRepo.CompleteInspection(ctx, id, inspectedBy, qcStatus, qcNotes); err != nil {
		return err
	}

	// If QC passed, create stock movements
	if qcStatus == "passed" {
		var stockRequests []client.StockMovementRequest

		for _, item := range grn.Items {
			// Use AcceptedQuantity if available, otherwise ReceivedQuantity
			qty := item.AcceptedQuantity
			if qty == 0 {
				qty = item.ReceivedQuantity
			}

			req := client.StockMovementRequest{
				ProductID:     item.ProductID.Hex(),
				MovementType:  "in",
				ToLocationID:  grn.LocationID.Hex(),
				Quantity:      qty,
				UnitCost:      item.UnitCost,
				ReferenceType: "grn",
				ReferenceNo:   grn.GRNNumber,
				BatchNumber:   item.BatchNumber,
				Notes:         fmt.Sprintf("Stock received from GRN %s", grn.GRNNumber),
			}
			stockRequests = append(stockRequests, req)
		}

		// Create stock movements
		if len(stockRequests) > 0 {
			if err := s.inventoryClient.CreateStockMovementsBatch(ctx, grn.OrganizationID, stockRequests, token); err != nil {
				// Log error but don't fail the inspection
				fmt.Printf("Warning: Failed to create stock movements for GRN %s: %v\n", grn.GRNNumber, err)
			}
		}
	}

	return nil
}

// ============== Helper Functions ==============

func (s *ProcurementService) generateSupplierCode(orgID primitive.ObjectID) string {
	timestamp := time.Now().Unix()
	return fmt.Sprintf("SUP-%06d", timestamp%1000000)
}

func (s *ProcurementService) generatePONumber(orgID primitive.ObjectID) string {
	timestamp := time.Now().Unix()
	return fmt.Sprintf("PO-%06d", timestamp%1000000)
}

func (s *ProcurementService) generateGRNNumber(orgID primitive.ObjectID) string {
	timestamp := time.Now().Unix()
	return fmt.Sprintf("GRN-%06d", timestamp%1000000)
}

func (s *ProcurementService) GetDashboardStats(ctx context.Context, orgID primitive.ObjectID) (map[string]interface{}, error) {
	// Pending POs
	pendingPOCount, pendingApprovals, err := s.poRepo.GetDashboardStats(ctx, orgID)
	if err != nil {
		return nil, fmt.Errorf("failed to get PO stats: %w", err)
	}

	// Pending GRNs
	pendingGRNCount, err := s.grnRepo.GetPendingCount(ctx, orgID)
	if err != nil {
		return nil, fmt.Errorf("failed to get GRN stats: %w", err)
	}

	return map[string]interface{}{
		"pending_po_count":  pendingPOCount,
		"pending_grn_count": pendingGRNCount,
		"pending_approvals": pendingApprovals,
	}, nil
}
