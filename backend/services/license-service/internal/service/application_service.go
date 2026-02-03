package service

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"github.com/yourusername/erp-system/services/license-service/internal/repository"
	"github.com/yourusername/erp-system/shared/models"
)

type ApplicationService struct {
	appRepo *repository.ApplicationRepository
}

func NewApplicationService(appRepo *repository.ApplicationRepository) *ApplicationService {
	return &ApplicationService{
		appRepo: appRepo,
	}
}

// CreateApplicationRequest represents a request to create an application
type CreateApplicationRequest struct {
	Name                  string                      `json:"name" binding:"required"`
	Code                  string                      `json:"code" binding:"required"`
	DisplayName           string                      `json:"display_name"`
	Description           string                      `json:"description"`
	Category              models.ApplicationCategory  `json:"category" binding:"required"`
	Version               string                      `json:"version"`
	SupportedLicenseTypes []models.LicenseType        `json:"supported_license_types"`
	DefaultLicenseType    models.LicenseType          `json:"default_license_type"`
	BasePrice             float64                     `json:"base_price"`
	PricePerUser          float64                     `json:"price_per_user"`
	PricePerDevice        float64                     `json:"price_per_device"`
	PricePerTransaction   float64                     `json:"price_per_transaction"`
	MinimumUsers          int                         `json:"minimum_users"`
	IncludedTransactions  int64                       `json:"included_transactions"`
	Features              []models.ApplicationFeature `json:"features"`
	APIEndpoint           string                      `json:"api_endpoint"`
	ServiceURL            string                      `json:"service_url"`
	Icon                  string                      `json:"icon"`
	IsPublic              bool                        `json:"is_public"`
}

// CreateApplication creates a new application
func (s *ApplicationService) CreateApplication(ctx context.Context, req CreateApplicationRequest, createdBy primitive.ObjectID) (*models.Application, error) {
	// Check if code already exists
	exists, err := s.appRepo.CodeExists(ctx, req.Code, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to check code existence: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("application code '%s' already exists", req.Code)
	}

	// Create application
	app := &models.Application{
		Name:                  req.Name,
		Code:                  req.Code,
		DisplayName:           req.DisplayName,
		Description:           req.Description,
		Category:              req.Category,
		Version:               req.Version,
		SupportedLicenseTypes: req.SupportedLicenseTypes,
		DefaultLicenseType:    req.DefaultLicenseType,
		BasePrice:             req.BasePrice,
		PricePerUser:          req.PricePerUser,
		PricePerDevice:        req.PricePerDevice,
		PricePerTransaction:   req.PricePerTransaction,
		MinimumUsers:          req.MinimumUsers,
		IncludedTransactions:  req.IncludedTransactions,
		Features:              req.Features,
		APIEndpoint:           req.APIEndpoint,
		ServiceURL:            req.ServiceURL,
		Icon:                  req.Icon,
		IsPublic:              req.IsPublic,
		IsActive:              true,
	}

	app.BaseModel.CreatedBy = createdBy

	if err := s.appRepo.Create(ctx, app); err != nil {
		return nil, fmt.Errorf("failed to create application: %w", err)
	}

	return app, nil
}

// GetApplication retrieves an application by ID
func (s *ApplicationService) GetApplication(ctx context.Context, id primitive.ObjectID) (*models.Application, error) {
	app, err := s.appRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("application not found: %w", err)
	}
	return app, nil
}

// GetApplicationByCode retrieves an application by code
func (s *ApplicationService) GetApplicationByCode(ctx context.Context, code string) (*models.Application, error) {
	app, err := s.appRepo.FindByCode(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("application not found: %w", err)
	}
	return app, nil
}

// ListApplications returns all applications
func (s *ApplicationService) ListApplications(ctx context.Context, category models.ApplicationCategory, publicOnly bool, page, limit int) ([]*models.Application, error) {
	apps, err := s.appRepo.FindAll(ctx, category, publicOnly, page, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to list applications: %w", err)
	}
	return apps, nil
}

// UpdateApplication updates an application
func (s *ApplicationService) UpdateApplication(ctx context.Context, id primitive.ObjectID, req CreateApplicationRequest, updatedBy primitive.ObjectID) (*models.Application, error) {
	// Find existing application
	app, err := s.appRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("application not found: %w", err)
	}

	// Check if code is being changed and if it exists
	if req.Code != app.Code {
		exists, err := s.appRepo.CodeExists(ctx, req.Code, &id)
		if err != nil {
			return nil, fmt.Errorf("failed to check code existence: %w", err)
		}
		if exists {
			return nil, fmt.Errorf("application code '%s' already exists", req.Code)
		}
	}

	// Update fields
	app.Name = req.Name
	app.Code = req.Code
	app.DisplayName = req.DisplayName
	app.Description = req.Description
	app.Category = req.Category
	app.Version = req.Version
	app.SupportedLicenseTypes = req.SupportedLicenseTypes
	app.DefaultLicenseType = req.DefaultLicenseType
	app.BasePrice = req.BasePrice
	app.PricePerUser = req.PricePerUser
	app.PricePerDevice = req.PricePerDevice
	app.PricePerTransaction = req.PricePerTransaction
	app.MinimumUsers = req.MinimumUsers
	app.IncludedTransactions = req.IncludedTransactions
	app.Features = req.Features
	app.APIEndpoint = req.APIEndpoint
	app.ServiceURL = req.ServiceURL
	app.Icon = req.Icon
	app.IsPublic = req.IsPublic
	app.BaseModel.UpdatedBy = updatedBy

	if err := s.appRepo.Update(ctx, app); err != nil {
		return nil, fmt.Errorf("failed to update application: %w", err)
	}

	return app, nil
}

// DeleteApplication soft deletes an application
func (s *ApplicationService) DeleteApplication(ctx context.Context, id, deletedBy primitive.ObjectID) error {
	if err := s.appRepo.SoftDelete(ctx, id, deletedBy); err != nil {
		return fmt.Errorf("failed to delete application: %w", err)
	}
	return nil
}
