package service

import (
	"context"
	"fmt"
	"time"

	"github.com/yourusername/erp-system/services/product-service/internal/repository"
	"github.com/yourusername/erp-system/shared/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BrandService struct {
	brandRepo *repository.BrandRepository
	orgRepo   *repository.OrganizationRepository
}

func NewBrandService(brandRepo *repository.BrandRepository, orgRepo *repository.OrganizationRepository) *BrandService {
	return &BrandService{
		brandRepo: brandRepo,
		orgRepo:   orgRepo,
	}
}

// CreateBrand creates a new brand
func (s *BrandService) CreateBrand(ctx context.Context, req CreateBrandRequest, createdBy primitive.ObjectID) (*models.Brand, error) {
	// Validate organization exists
	orgID, err := primitive.ObjectIDFromHex(req.OrganizationID)
	if err != nil {
		return nil, fmt.Errorf("invalid organization ID: %w", err)
	}

	exists, err := s.orgRepo.Exists(ctx, orgID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, fmt.Errorf("organization not found")
	}

	// Check if brand code already exists in organization
	if req.Code != "" {
		exists, err := s.brandRepo.CodeExistsInOrg(ctx, req.Code, orgID, nil)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, fmt.Errorf("brand code already exists in organization")
		}
	}

	// Check if brand name already exists in organization
	exists, err = s.brandRepo.NameExistsInOrg(ctx, req.Name, orgID, nil)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, fmt.Errorf("brand name already exists in organization")
	}

	brand := &models.Brand{
		OrganizationID: orgID,
		Name:           req.Name,
		Code:           req.Code,
		Description:    req.Description,
		LogoURL:        req.LogoURL,
		Website:        req.Website,
		Country:        req.Country,
		IsActive:       req.IsActive,
		Metadata:       req.Metadata,
	}

	brand.CreatedBy = createdBy
	brand.CreatedAt = time.Now()
	brand.UpdatedAt = time.Now()

	if err := s.brandRepo.Create(ctx, brand); err != nil {
		return nil, err
	}

	return brand, nil
}

// GetBrand retrieves a brand by ID
func (s *BrandService) GetBrand(ctx context.Context, id primitive.ObjectID) (*models.Brand, error) {
	brand, err := s.brandRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if brand == nil {
		return nil, fmt.Errorf("brand not found")
	}
	return brand, nil
}

// GetBrandsByOrganization retrieves all brands for an organization
func (s *BrandService) GetBrandsByOrganization(ctx context.Context, orgID primitive.ObjectID, filter BrandFilter) ([]*models.Brand, int64, error) {
	// If a search query is provided, use the search method
	if filter.Query != "" {
		return s.brandRepo.Search(ctx, orgID, filter.Query, filter.IsActive, filter.Page, filter.Limit)
	}
	// Otherwise, use the regular list method
	return s.brandRepo.FindByOrganization(ctx, orgID, filter.Page, filter.Limit, filter.IsActive)
}

// UpdateBrand updates an existing brand
func (s *BrandService) UpdateBrand(ctx context.Context, id primitive.ObjectID, req UpdateBrandRequest, updatedBy primitive.ObjectID) (*models.Brand, error) {
	// Check if brand exists
	brand, err := s.brandRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if brand == nil {
		return nil, fmt.Errorf("brand not found")
	}

	// Check if updating code and if new code already exists
	if req.Code != nil && *req.Code != brand.Code {
		exists, err := s.brandRepo.CodeExistsInOrg(ctx, *req.Code, brand.OrganizationID, &id)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, fmt.Errorf("brand code already exists in organization")
		}
		brand.Code = *req.Code
	}

	// Check if updating name and if new name already exists
	if req.Name != nil && *req.Name != brand.Name {
		exists, err := s.brandRepo.NameExistsInOrg(ctx, *req.Name, brand.OrganizationID, &id)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, fmt.Errorf("brand name already exists in organization")
		}
		brand.Name = *req.Name
	}

	// Update other fields
	if req.Description != nil {
		brand.Description = *req.Description
	}
	if req.LogoURL != nil {
		brand.LogoURL = *req.LogoURL
	}
	if req.Website != nil {
		brand.Website = *req.Website
	}
	if req.Country != nil {
		brand.Country = *req.Country
	}
	if req.IsActive != nil {
		brand.IsActive = *req.IsActive
	}
	if req.Metadata != nil {
		brand.Metadata = req.Metadata
	}

	brand.UpdatedBy = updatedBy
	brand.UpdatedAt = time.Now()

	if err := s.brandRepo.Update(ctx, brand); err != nil {
		return nil, err
	}

	return brand, nil
}

// DeleteBrand soft deletes a brand
func (s *BrandService) DeleteBrand(ctx context.Context, id primitive.ObjectID, deletedBy primitive.ObjectID) error {
	brand, err := s.brandRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if brand == nil {
		return fmt.Errorf("brand not found")
	}

	// Check if brand is being used by any products
	inUse, err := s.brandRepo.IsInUse(ctx, id)
	if err != nil {
		return err
	}
	if inUse {
		return fmt.Errorf("cannot delete brand that is in use by products")
	}

	return s.brandRepo.Delete(ctx, id, deletedBy)
}

// SearchBrands searches brands by name or code
func (s *BrandService) SearchBrands(ctx context.Context, orgID primitive.ObjectID, query string, isActive *bool, page, limit int) ([]*models.Brand, int64, error) {
	brands, total, err := s.brandRepo.Search(ctx, orgID, query, isActive, page, limit)
	if err != nil {
		return nil, 0, err
	}
	return brands, total, nil
}

// Request DTOs
type CreateBrandRequest struct {
	OrganizationID string                 `json:"organization_id" binding:"required"`
	Name           string                 `json:"name" binding:"required"`
	Code           string                 `json:"code"`
	Description    string                 `json:"description"`
	LogoURL        string                 `json:"logo_url"`
	Website        string                 `json:"website"`
	Country        string                 `json:"country"`
	IsActive       bool                   `json:"is_active"`
	Metadata       map[string]interface{} `json:"metadata"`
}

type UpdateBrandRequest struct {
	Name        *string                `json:"name"`
	Code        *string                `json:"code"`
	Description *string                `json:"description"`
	LogoURL     *string                `json:"logo_url"`
	Website     *string                `json:"website"`
	Country     *string                `json:"country"`
	IsActive    *bool                  `json:"is_active"`
	Metadata    map[string]interface{} `json:"metadata"`
}

type BrandFilter struct {
	Page     int
	Limit    int
	IsActive *bool
	Query    string // Optional search query
}
