package service

import (
	"context"
	"fmt"
	"time"

	"github.com/yourusername/erp-system/services/product-service/internal/repository"
	"github.com/yourusername/erp-system/shared/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CategoryService struct {
	categoryRepo *repository.CategoryRepository
	orgRepo      *repository.OrganizationRepository
}

func NewCategoryService(categoryRepo *repository.CategoryRepository, orgRepo *repository.OrganizationRepository) *CategoryService {
	return &CategoryService{
		categoryRepo: categoryRepo,
		orgRepo:      orgRepo,
	}
}

// CreateCategory creates a new category with optional embedded subcategories
func (s *CategoryService) CreateCategory(ctx context.Context, req CreateCategoryRequest, userOrgID primitive.ObjectID) (*models.ProductCategory, error) {
	var orgID primitive.ObjectID

	// For root categories, use organization ID from request
	if req.OrganizationID == "" {
		return nil, fmt.Errorf("organization_id is required for categories")
	}

	var err error
	orgID, err = primitive.ObjectIDFromHex(req.OrganizationID)
	if err != nil {
		return nil, fmt.Errorf("invalid organization ID: %w", err)
	}

	// Verify organization exists
	exists, err := s.orgRepo.Exists(ctx, orgID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, fmt.Errorf("organization not found")
	}

	// Check if user belongs to this organization
	if orgID != userOrgID {
		return nil, fmt.Errorf("unauthorized: cannot create category for different organization")
	}

	// Convert DTO to model
	category := &models.ProductCategory{
		OrganizationID: orgID,
		Name:           req.Name,
		Code:           req.Code,
		Description:    req.Description,
		IsActive:       req.IsActive,
		Metadata:       req.Metadata,
		Subcategories:  []models.ProductSubcategory{},
	}

	// Add subcategories if provided
	if len(req.Subcategories) > 0 {
		for _, subReq := range req.Subcategories {
			subcategory := models.ProductSubcategory{
				Name:        subReq.Name,
				Code:        subReq.Code,
				Description: subReq.Description,
				IsActive:    subReq.IsActive,
				Metadata:    subReq.Metadata,
			}
			category.Subcategories = append(category.Subcategories, subcategory)
		}
	}

	// Check if name already exists
	nameExists, err := s.categoryRepo.CheckNameExists(ctx, category.OrganizationID, category.Name, nil)
	if err != nil {
		return nil, err
	}
	if nameExists {
		return nil, fmt.Errorf("category name already exists")
	}

	if err := s.categoryRepo.Create(ctx, category); err != nil {
		return nil, err
	}

	return category, nil
}

// GetCategory retrieves a category by ID
func (s *CategoryService) GetCategory(ctx context.Context, id primitive.ObjectID, userOrgID primitive.ObjectID) (*models.ProductCategory, error) {
	category, err := s.categoryRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if category == nil {
		return nil, fmt.Errorf("category not found")
	}

	// Check if user belongs to this organization
	if category.OrganizationID != userOrgID {
		return nil, fmt.Errorf("unauthorized: category belongs to different organization")
	}

	return category, nil
}

// ListCategories retrieves categories with filters
func (s *CategoryService) ListCategories(ctx context.Context, orgID primitive.ObjectID, isActive *bool, query string, page, limit int) ([]*models.ProductCategory, int64, error) {
	// Verify organization exists
	exists, err := s.orgRepo.Exists(ctx, orgID)
	if err != nil {
		return nil, 0, err
	}
	if !exists {
		return nil, 0, fmt.Errorf("organization not found")
	}

	return s.categoryRepo.FindByOrganization(ctx, orgID, isActive, query, page, limit)
}

// GetCategoryTree retrieves the full category tree (categories with embedded subcategories)
func (s *CategoryService) GetCategoryTree(ctx context.Context, orgID primitive.ObjectID) ([]*models.ProductCategory, error) {
	return s.categoryRepo.GetCategoryTree(ctx, orgID)
}

// GetRootCategories is now synonymous with listing all categories (since all are top-level)
func (s *CategoryService) GetRootCategories(ctx context.Context, orgID primitive.ObjectID) ([]*models.ProductCategory, error) {
	categories, _, err := s.categoryRepo.FindByOrganization(ctx, orgID, nil, "", 1, 1000)
	return categories, err
}

// GetChildren retrieves subcategories for a specific category (from embedded array)
func (s *CategoryService) GetChildren(ctx context.Context, categoryID primitive.ObjectID, userOrgID primitive.ObjectID) ([]models.ProductSubcategory, error) {
	category, err := s.categoryRepo.FindByID(ctx, categoryID)
	if err != nil {
		return nil, err
	}
	if category == nil {
		return nil, fmt.Errorf("category not found")
	}
	if category.OrganizationID != userOrgID {
		return nil, fmt.Errorf("unauthorized: category belongs to different organization")
	}

	// Filter out soft-deleted subcategories
	activeSubcategories := make([]models.ProductSubcategory, 0)
	for _, sub := range category.Subcategories {
		if sub.DeletedAt == nil {
			activeSubcategories = append(activeSubcategories, sub)
		}
	}

	return activeSubcategories, nil
}

// UpdateCategory updates an existing category and its embedded subcategories
func (s *CategoryService) UpdateCategory(ctx context.Context, id primitive.ObjectID, req UpdateCategoryRequest, userOrgID primitive.ObjectID) (*models.ProductCategory, error) {
	// Get existing category
	existing, err := s.categoryRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, fmt.Errorf("category not found")
	}

	// Check if user belongs to this organization
	if existing.OrganizationID != userOrgID {
		return nil, fmt.Errorf("unauthorized: category belongs to different organization")
	}

	// Handle subcategory removal (mark as deleted in embedded array)
	if len(req.RemoveSubcategories) > 0 {
		for _, subIDStr := range req.RemoveSubcategories {
			subID, err := primitive.ObjectIDFromHex(subIDStr)
			if err != nil {
				return nil, fmt.Errorf("invalid subcategory ID '%s': %w", subIDStr, err)
			}

			found := false
			for i := range existing.Subcategories {
				if existing.Subcategories[i].ID == subID && existing.Subcategories[i].DeletedAt == nil {
					// Soft delete the subcategory
					now := time.Now()
					existing.Subcategories[i].DeletedAt = &now
					existing.Subcategories[i].UpdatedAt = now
					found = true
					break
				}
			}
			if !found {
				return nil, fmt.Errorf("subcategory with ID '%s' not found or already deleted", subIDStr)
			}
		}
	}

	// Handle subcategory updates
	if len(req.UpdateSubcategories) > 0 {
		for _, updateReq := range req.UpdateSubcategories {
			subID, err := primitive.ObjectIDFromHex(updateReq.ID)
			if err != nil {
				return nil, fmt.Errorf("invalid subcategory ID '%s': %w", updateReq.ID, err)
			}

			found := false
			for i := range existing.Subcategories {
				if existing.Subcategories[i].ID == subID && existing.Subcategories[i].DeletedAt == nil {
					// Apply updates
					if updateReq.Name != nil {
						existing.Subcategories[i].Name = *updateReq.Name
					}
					if updateReq.Code != nil {
						existing.Subcategories[i].Code = *updateReq.Code
					}
					if updateReq.Description != nil {
						existing.Subcategories[i].Description = *updateReq.Description
					}
					if updateReq.IsActive != nil {
						existing.Subcategories[i].IsActive = *updateReq.IsActive
					}
					if updateReq.Metadata != nil {
						existing.Subcategories[i].Metadata = updateReq.Metadata
					}
					existing.Subcategories[i].UpdatedAt = time.Now()
					found = true
					break
				}
			}
			if !found {
				return nil, fmt.Errorf("subcategory with ID '%s' not found or deleted", updateReq.ID)
			}
		}
	}

	// Handle adding new subcategories
	if len(req.AddSubcategories) > 0 {
		for _, subReq := range req.AddSubcategories {
			subcategory := models.ProductSubcategory{
				BaseModel: models.BaseModel{
					ID:        primitive.NewObjectID(),
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
				Name:        subReq.Name,
				Code:        subReq.Code,
				Description: subReq.Description,
				IsActive:    subReq.IsActive,
				Metadata:    subReq.Metadata,
			}
			existing.Subcategories = append(existing.Subcategories, subcategory)
		}
	}

	// Apply updates to main category
	if req.Name != nil {
		existing.Name = *req.Name
	}
	if req.Code != nil {
		existing.Code = *req.Code
	}
	if req.Description != nil {
		existing.Description = *req.Description
	}
	if req.IsActive != nil {
		existing.IsActive = *req.IsActive
	}
	if req.Metadata != nil {
		existing.Metadata = req.Metadata
	}

	// Check if name already exists (excluding current category)
	if req.Name != nil {
		exists, err := s.categoryRepo.CheckNameExists(ctx, existing.OrganizationID, existing.Name, &id)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, fmt.Errorf("category name already exists")
		}
	}

	if err := s.categoryRepo.Update(ctx, existing); err != nil {
		return nil, err
	}

	return existing, nil
}

// DeleteCategory deletes a category
func (s *CategoryService) DeleteCategory(ctx context.Context, id primitive.ObjectID, userOrgID primitive.ObjectID) error {
	// Get category
	category, err := s.categoryRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if category == nil {
		return fmt.Errorf("category not found")
	}

	// Check if user belongs to this organization
	if category.OrganizationID != userOrgID {
		return fmt.Errorf("unauthorized: category belongs to different organization")
	}

	// Check if category has children (non-deleted subcategories)
	hasChildren, err := s.categoryRepo.HasChildren(ctx, id)
	if err != nil {
		return err
	}
	if hasChildren {
		return fmt.Errorf("cannot delete category with subcategories")
	}

	// Check if category has products
	hasProducts, err := s.categoryRepo.HasProducts(ctx, id)
	if err != nil {
		return err
	}
	if hasProducts {
		return fmt.Errorf("cannot delete category with products")
	}

	return s.categoryRepo.Delete(ctx, id)
}

// UpdateProductCount updates the product count for a category
func (s *CategoryService) UpdateProductCount(ctx context.Context, categoryID primitive.ObjectID) error {
	return s.categoryRepo.UpdateProductCount(ctx, categoryID)
}

// Request DTOs
type CreateCategoryRequest struct {
	OrganizationID string                     `json:"organization_id" binding:"required"`
	Name           string                     `json:"name" binding:"required"`
	Code           string                     `json:"code"`
	Description    string                     `json:"description"`
	IsActive       bool                       `json:"is_active"`
	Metadata       map[string]interface{}     `json:"metadata"`
	Subcategories  []CreateSubcategoryRequest `json:"subcategories"` // Embedded subcategories
}

type CreateSubcategoryRequest struct {
	Name        string                 `json:"name" binding:"required"`
	Code        string                 `json:"code"`
	Description string                 `json:"description"`
	IsActive    bool                   `json:"is_active"`
	Metadata    map[string]interface{} `json:"metadata"`
}

type UpdateCategoryRequest struct {
	Name        *string                `json:"name"`
	Code        *string                `json:"code"`
	Description *string                `json:"description"`
	IsActive    *bool                  `json:"is_active"`
	Metadata    map[string]interface{} `json:"metadata"`
	// Subcategory management
	AddSubcategories    []CreateSubcategoryRequest `json:"add_subcategories"`    // New subcategories to create
	UpdateSubcategories []UpdateSubcategoryRequest `json:"update_subcategories"` // Existing subcategories to update
	RemoveSubcategories []string                   `json:"remove_subcategories"` // Subcategory IDs to soft delete
}

type UpdateSubcategoryRequest struct {
	ID          string                 `json:"id" binding:"required"` // Subcategory ID to update
	Name        *string                `json:"name"`
	Code        *string                `json:"code"`
	Description *string                `json:"description"`
	IsActive    *bool                  `json:"is_active"`
	Metadata    map[string]interface{} `json:"metadata"`
}

type CategoryFilter struct {
	IsActive *bool
	Page     int
	Limit    int
}
