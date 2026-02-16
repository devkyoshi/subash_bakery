package service

import (
	"context"
	"fmt"
	"time"

	"github.com/yourusername/erp-system/services/org-service/internal/repository"
	"github.com/yourusername/erp-system/shared/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OrganizationService struct {
	orgRepo          *repository.OrganizationRepository
	companyRepo      *repository.CompanyRepository
	locationRepo     *repository.LocationRepository
	locationUserRepo *repository.LocationUserRepository
}

func NewOrganizationService(
	orgRepo *repository.OrganizationRepository,
	companyRepo *repository.CompanyRepository,
	locationRepo *repository.LocationRepository,
	locationUserRepo *repository.LocationUserRepository,
) *OrganizationService {
	return &OrganizationService{
		orgRepo:          orgRepo,
		companyRepo:      companyRepo,
		locationRepo:     locationRepo,
		locationUserRepo: locationUserRepo,
	}
}

// CreateOrganization creates a new organization
func (s *OrganizationService) CreateOrganization(ctx context.Context, req CreateOrganizationRequest, createdBy primitive.ObjectID) (*models.Organization, error) {
	// Check domain uniqueness
	exists, err := s.orgRepo.DomainExists(ctx, req.Domain, nil)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, fmt.Errorf("domain already exists")
	}

	// Create organization with default settings
	org := &models.Organization{
		Name:               req.Name,
		LegalName:          req.LegalName,
		Domain:             req.Domain,
		Email:              req.Email,
		Phone:              req.Phone,
		Website:            req.Website,
		TaxID:              req.TaxID,
		RegistrationNumber: req.RegistrationNumber,
		Industry:           req.Industry,
		CompanySize:        req.CompanySize,
		Status:             models.OrganizationStatusActive,
		IsActive:           true,
		BillingEmail:       req.BillingEmail,
		MaxUsers:           req.MaxUsers,
		MaxCompanies:       req.MaxCompanies,
		MaxLocations:       req.MaxLocations,
		StorageLimitGB:     req.StorageLimitGB,
		Settings: models.OrganizationSettings{
			Timezone:                 "UTC",
			DateFormat:               "YYYY-MM-DD",
			TimeFormat:               "HH:mm:ss",
			Currency:                 "USD",
			Language:                 "en",
			EnabledModules:           []string{},
			AllowUserRegistration:    false,
			RequireEmailVerification: true,
			EnableMFA:                false,
			SessionTimeout:           30,
			PasswordPolicy: models.PasswordPolicy{
				MinLength:           8,
				RequireUppercase:    true,
				RequireLowercase:    true,
				RequireNumbers:      true,
				RequireSpecialChars: false,
				ExpiryDays:          90,
				PreventReuseCount:   5,
			},
		},
	}

	org.CreatedBy = createdBy
	now := time.Now()
	org.ActivatedAt = &now

	if err := s.orgRepo.Create(ctx, org); err != nil {
		return nil, err
	}

	return org, nil
}

// GetOrganization retrieves an organization by ID
func (s *OrganizationService) GetOrganization(ctx context.Context, id primitive.ObjectID) (*models.Organization, error) {
	org, err := s.orgRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if org == nil {
		return nil, fmt.Errorf("organization not found")
	}
	return org, nil
}

// ListOrganizations returns paginated organizations
func (s *OrganizationService) ListOrganizations(ctx context.Context, page, limit int, status string) ([]*models.Organization, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	return s.orgRepo.List(ctx, page, limit, status)
}

// GetOrganizationOptions returns a list of organizations with minimal details (ID, Name) for dropdowns
func (s *OrganizationService) GetOrganizationOptions(ctx context.Context) ([]*models.OrganizationOption, error) {
	return s.orgRepo.ListOptions(ctx)
}

// UpdateOrganization updates an organization
func (s *OrganizationService) UpdateOrganization(ctx context.Context, id primitive.ObjectID, req UpdateOrganizationRequest, updatedBy primitive.ObjectID) (*models.Organization, error) {
	org, err := s.orgRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if org == nil {
		return nil, fmt.Errorf("organization not found")
	}

	// Check domain uniqueness if changed
	if req.Domain != nil && *req.Domain != org.Domain {
		exists, err := s.orgRepo.DomainExists(ctx, *req.Domain, &id)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, fmt.Errorf("domain already exists")
		}
		org.Domain = *req.Domain
	}

	// Update fields
	if req.Name != nil {
		org.Name = *req.Name
	}
	if req.LegalName != nil {
		org.LegalName = *req.LegalName
	}
	if req.Email != nil {
		org.Email = *req.Email
	}
	if req.Phone != nil {
		org.Phone = *req.Phone
	}
	if req.Website != nil {
		org.Website = *req.Website
	}
	if req.Logo != nil {
		org.Logo = *req.Logo
	}
	if req.IsActive != nil {
		org.IsActive = *req.IsActive
	}

	org.UpdatedBy = updatedBy

	if err := s.orgRepo.Update(ctx, org); err != nil {
		return nil, err
	}

	return org, nil
}

// DeleteOrganization soft deletes an organization
func (s *OrganizationService) DeleteOrganization(ctx context.Context, id, deletedBy primitive.ObjectID) error {
	org, err := s.orgRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if org == nil {
		return fmt.Errorf("organization not found")
	}

	return s.orgRepo.SoftDelete(ctx, id, deletedBy)
}

// CreateCompany creates a new company under an organization
func (s *OrganizationService) CreateCompany(ctx context.Context, orgID primitive.ObjectID, req CreateCompanyRequest, createdBy primitive.ObjectID) (*models.Company, error) {
	// Verify organization exists and has capacity
	org, err := s.orgRepo.FindByID(ctx, orgID)
	if err != nil {
		return nil, err
	}
	if org == nil {
		return nil, fmt.Errorf("organization not found")
	}

	// Check company limit
	currentCount, err := s.companyRepo.CountByOrganization(ctx, orgID)
	if err != nil {
		return nil, err
	}
	if org.MaxCompanies > 0 && currentCount >= org.MaxCompanies {
		return nil, fmt.Errorf("maximum company limit reached")
	}

	// Check code uniqueness
	exists, err := s.companyRepo.CodeExists(ctx, orgID, req.Code, nil)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, fmt.Errorf("company code already exists")
	}

	company := &models.Company{
		OrganizationID:     orgID,
		Name:               req.Name,
		LegalName:          req.LegalName,
		Code:               req.Code,
		TaxID:              req.TaxID,
		RegistrationNumber: req.RegistrationNumber,
		VATNumber:          req.VATNumber,
		Email:              req.Email,
		Phone:              req.Phone,
		Address:            req.Address,
		BankAccounts:       req.BankAccounts,
		Settings: models.CompanySettings{
			FiscalYearStart:     "01-01",
			Currency:            org.Settings.Currency,
			Timezone:            org.Settings.Timezone,
			EnableMultiCurrency: false,
		},
		IsActive:  true,
		IsDefault: req.IsDefault,
	}

	company.CreatedBy = createdBy

	if err := s.companyRepo.Create(ctx, company); err != nil {
		return nil, err
	}

	// Update organization usage
	s.orgRepo.UpdateUsage(ctx, orgID, org.CurrentUsers, currentCount+1, org.CurrentLocations, org.StorageUsedGB)

	return company, nil
}

// GetCompany retrieves a company by ID
func (s *OrganizationService) GetCompany(ctx context.Context, id primitive.ObjectID) (*models.Company, error) {
	company, err := s.companyRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if company == nil {
		return nil, fmt.Errorf("company not found")
	}
	return company, nil
}

// ListCompanies returns paginated companies for an organization
func (s *OrganizationService) ListCompanies(ctx context.Context, orgID primitive.ObjectID, page, limit int) ([]*models.Company, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	return s.companyRepo.FindByOrganization(ctx, orgID, page, limit)
}

// UpdateCompany updates a company
func (s *OrganizationService) UpdateCompany(ctx context.Context, id primitive.ObjectID, req UpdateCompanyRequest, updatedBy primitive.ObjectID) (*models.Company, error) {
	company, err := s.companyRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if company == nil {
		return nil, fmt.Errorf("company not found")
	}

	// Check code uniqueness if changed
	if req.Code != nil && *req.Code != company.Code {
		exists, err := s.companyRepo.CodeExists(ctx, company.OrganizationID, *req.Code, &id)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, fmt.Errorf("company code already exists")
		}
		company.Code = *req.Code
	}

	// Update fields
	if req.Name != nil {
		company.Name = *req.Name
	}
	if req.LegalName != nil {
		company.LegalName = *req.LegalName
	}
	if req.Email != nil {
		company.Email = *req.Email
	}
	if req.Phone != nil {
		company.Phone = *req.Phone
	}
	if req.IsActive != nil {
		company.IsActive = *req.IsActive
	}

	company.UpdatedBy = updatedBy

	if err := s.companyRepo.Update(ctx, company); err != nil {
		return nil, err
	}

	return company, nil
}

// DeleteCompany soft deletes a company
func (s *OrganizationService) DeleteCompany(ctx context.Context, id, deletedBy primitive.ObjectID) error {
	company, err := s.companyRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if company == nil {
		return fmt.Errorf("company not found")
	}

	return s.companyRepo.SoftDelete(ctx, id, deletedBy)
}

// CreateLocation creates a new location under a company
func (s *OrganizationService) CreateLocation(ctx context.Context, companyID primitive.ObjectID, req CreateLocationRequest, createdBy primitive.ObjectID) (*models.Location, error) {
	// Verify company exists
	company, err := s.companyRepo.FindByID(ctx, companyID)
	if err != nil {
		return nil, err
	}
	if company == nil {
		return nil, fmt.Errorf("company not found")
	}

	// Check organization location limit
	org, err := s.orgRepo.FindByID(ctx, company.OrganizationID)
	if err != nil {
		return nil, err
	}

	currentCount, err := s.locationRepo.CountByOrganization(ctx, company.OrganizationID)
	if err != nil {
		return nil, err
	}
	if org.MaxLocations > 0 && currentCount >= org.MaxLocations {
		return nil, fmt.Errorf("maximum location limit reached")
	}

	// Check code uniqueness
	exists, err := s.locationRepo.CodeExists(ctx, companyID, req.Code, nil)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, fmt.Errorf("location code already exists")
	}

	location := &models.Location{
		CompanyID:      companyID,
		OrganizationID: company.OrganizationID,
		Name:           req.Name,
		Code:           req.Code,
		Type:           req.Type,
		Email:          req.Email,
		Phone:          req.Phone,
		Address:        req.Address,
		WarehouseInfo:  req.WarehouseInfo,
		StoreInfo:      req.StoreInfo,
		Settings: models.LocationSettings{
			Timezone:                   org.Settings.Timezone,
			AllowBackdatedTransactions: false,
			RequireApproval:            true,
		},
		IsActive:  true,
		IsDefault: req.IsDefault,
	}

	location.CreatedBy = createdBy

	if err := s.locationRepo.Create(ctx, location); err != nil {
		return nil, err
	}

	// Update organization usage
	s.orgRepo.UpdateUsage(ctx, company.OrganizationID, org.CurrentUsers, org.CurrentCompanies, currentCount+1, org.StorageUsedGB)

	return location, nil
}

// GetLocation retrieves a location by ID
func (s *OrganizationService) GetLocation(ctx context.Context, id primitive.ObjectID) (*models.Location, error) {
	location, err := s.locationRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if location == nil {
		return nil, fmt.Errorf("location not found")
	}
	return location, nil
}

// ListLocations returns paginated locations for a company
func (s *OrganizationService) ListLocations(ctx context.Context, companyID primitive.ObjectID, page, limit int) ([]*models.Location, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	return s.locationRepo.FindByCompany(ctx, companyID, page, limit)
}

// Request/Response DTOs
type CreateOrganizationRequest struct {
	Name               string  `json:"name" binding:"required"`
	LegalName          string  `json:"legal_name" binding:"required"`
	Domain             string  `json:"domain" binding:"required"`
	Email              string  `json:"email" binding:"required,email"`
	Phone              string  `json:"phone"`
	Website            string  `json:"website"`
	TaxID              string  `json:"tax_id"`
	RegistrationNumber string  `json:"registration_number"`
	Industry           string  `json:"industry"`
	CompanySize        string  `json:"company_size"`
	BillingEmail       string  `json:"billing_email" binding:"required,email"`
	MaxUsers           int     `json:"max_users"`
	MaxCompanies       int     `json:"max_companies"`
	MaxLocations       int     `json:"max_locations"`
	StorageLimitGB     float64 `json:"storage_limit_gb"`
}

// ListLocationsByOrganization returns paginated locations for an organization
func (s *OrganizationService) ListLocationsByOrganization(ctx context.Context, orgID primitive.ObjectID, page, limit int) ([]*models.Location, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	return s.locationRepo.FindByOrganization(ctx, orgID, "", page, limit)
}

type UpdateOrganizationRequest struct {
	Name      *string `json:"name"`
	LegalName *string `json:"legal_name"`
	Domain    *string `json:"domain"`
	Email     *string `json:"email"`
	Phone     *string `json:"phone"`
	Website   *string `json:"website"`
	Logo      *string `json:"logo"`
	IsActive  *bool   `json:"is_active"`
}

type CreateCompanyRequest struct {
	Name               string               `json:"name" binding:"required"`
	LegalName          string               `json:"legal_name" binding:"required"`
	Code               string               `json:"code" binding:"required"`
	TaxID              string               `json:"tax_id"`
	RegistrationNumber string               `json:"registration_number"`
	VATNumber          string               `json:"vat_number"`
	Email              string               `json:"email" binding:"required,email"`
	Phone              string               `json:"phone"`
	Address            models.Address       `json:"address" binding:"required"`
	BankAccounts       []models.BankAccount `json:"bank_accounts"`
	IsDefault          bool                 `json:"is_default"`
}

type UpdateCompanyRequest struct {
	Name      *string `json:"name"`
	LegalName *string `json:"legal_name"`
	Code      *string `json:"code"`
	Email     *string `json:"email"`
	Phone     *string `json:"phone"`
	IsActive  *bool   `json:"is_active"`
}

type CreateLocationRequest struct {
	Name          string                `json:"name" binding:"required"`
	Code          string                `json:"code" binding:"required"`
	Type          models.LocationType   `json:"type" binding:"required"`
	Email         string                `json:"email" binding:"required,email"`
	Phone         string                `json:"phone"`
	Address       models.Address        `json:"address" binding:"required"`
	WarehouseInfo *models.WarehouseInfo `json:"warehouse_info"`
	StoreInfo     *models.StoreInfo     `json:"store_info"`
	IsDefault     bool                  `json:"is_default"`
}

// GetUserAccess retrieves all organizations, companies, and locations accessible by a user
func (s *OrganizationService) GetUserAccess(ctx context.Context, userID primitive.ObjectID, orgID primitive.ObjectID) (*UserAccessResponse, error) {
	// Get the user's organization
	org, err := s.orgRepo.FindByID(ctx, orgID)
	if err != nil {
		return nil, err
	}
	if org == nil {
		return nil, fmt.Errorf("organization not found")
	}

	// Get all locations the user has access to
	locationIDs, err := s.locationUserRepo.FindLocationsByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Get location details
	locations, err := s.locationRepo.FindByIDs(ctx, locationIDs)
	if err != nil {
		return nil, err
	}

	// Extract unique company IDs from locations
	companyIDMap := make(map[primitive.ObjectID]bool)
	for _, loc := range locations {
		companyIDMap[loc.CompanyID] = true
	}

	companyIDs := make([]primitive.ObjectID, 0, len(companyIDMap))
	for companyID := range companyIDMap {
		companyIDs = append(companyIDs, companyID)
	}

	// Get company details
	companies, err := s.companyRepo.FindByIDs(ctx, companyIDs)
	if err != nil {
		return nil, err
	}

	// Organize data by company
	companyMap := make(map[primitive.ObjectID]*CompanyWithLocations)
	for _, company := range companies {
		companyMap[company.ID] = &CompanyWithLocations{
			Company:   company,
			Locations: []*models.Location{},
		}
	}

	// Add locations to their respective companies
	for _, loc := range locations {
		if companyData, exists := companyMap[loc.CompanyID]; exists {
			companyData.Locations = append(companyData.Locations, loc)
		}
	}

	// Convert map to slice
	companiesWithLocations := make([]*CompanyWithLocations, 0, len(companyMap))
	for _, companyData := range companyMap {
		companiesWithLocations = append(companiesWithLocations, companyData)
	}

	return &UserAccessResponse{
		Organization: org,
		Companies:    companiesWithLocations,
	}, nil
}

// AssignUserToCompany assigns a user to all locations within a company.
// This effectively gives the user access to the company.
func (s *OrganizationService) AssignUserToCompany(ctx context.Context, companyID primitive.ObjectID, userID primitive.ObjectID, roleID primitive.ObjectID) error {
	// Verify company exists
	company, err := s.companyRepo.FindByID(ctx, companyID)
	if err != nil {
		return err
	}
	if company == nil {
		return fmt.Errorf("company not found")
	}

	// Get all locations for this company
	locations, _, err := s.locationRepo.FindByCompany(ctx, companyID, 1, 1000) // Assume max 1000 locations for now
	if err != nil {
		return err
	}

	for _, location := range locations {
		// Check if user is already assigned
		existing, err := s.locationUserRepo.FindByUserIDAndLocationID(ctx, userID, location.ID)
		if err != nil {
			return err
		}
		if existing != nil {
			continue // Skip if already assigned
		}

		locationUser := &models.LocationUser{
			LocationID:  location.ID,
			UserID:      userID,
			RoleID:      roleID,
			IsPrimary:   false,
			IsActive:    true,
			AssignedAt:  time.Now(),
			AccessLevel: "full", // Default access level
		}

		if err := s.locationUserRepo.Create(ctx, locationUser); err != nil {
			return err
		}
	}

	return nil
}

// Response structures for user access
type UserAccessResponse struct {
	Organization *models.Organization    `json:"organization"`
	Companies    []*CompanyWithLocations `json:"companies"`
}

type CompanyWithLocations struct {
	Company   *models.Company    `json:"company"`
	Locations []*models.Location `json:"locations"`
}
