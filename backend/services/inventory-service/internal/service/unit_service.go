package service

import (
	"context"
	"fmt"
	"time"

	"github.com/yourusername/erp-system/services/inventory-service/internal/repository"
	"github.com/yourusername/erp-system/shared/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UnitService struct {
	unitRepo      *repository.UnitRepository
	unitChartRepo *repository.UnitChartRepository
}

func NewUnitService(unitRepo *repository.UnitRepository, unitChartRepo *repository.UnitChartRepository) *UnitService {
	return &UnitService{
		unitRepo:      unitRepo,
		unitChartRepo: unitChartRepo,
	}
}

type CreateUnitRequest struct {
	Name        string `json:"name" binding:"required"`
	Code        string `json:"code" binding:"required"`
	Description string `json:"description"`
	UnitType    string `json:"unit_type" binding:"required"` // e.g., "weight", "volume", "quantity"
	IsBaseUnit  bool   `json:"is_base_unit"`
	IsActive    bool   `json:"is_active,omitempty"` // Default true if omitted/false? usually explicit
}

type UpdateUnitRequest struct {
	Name        *string `json:"name"`
	Code        *string `json:"code"`
	Description *string `json:"description"`
	UnitType    *string `json:"unit_type"`
	IsBaseUnit  *bool   `json:"is_base_unit"`
	IsActive    *bool   `json:"is_active"`
}

func (s *UnitService) CreateUnit(ctx context.Context, req CreateUnitRequest, createdBy primitive.ObjectID) (*models.Unit, error) {
	// TODO: Check if code already exists (uniqueness)

	unit := &models.Unit{
		Name:        req.Name,
		Code:        req.Code,
		Description: req.Description,
		UnitType:    req.UnitType,
		IsBaseUnit:  req.IsBaseUnit,
		IsActive:    true,
	}

	// If IsActive is explicitly provided in request, we could use it, but struct tag was omitempty.
	// Let's assume default is true.

	unit.ID = primitive.NewObjectID()
	unit.CreatedAt = time.Now()
	unit.UpdatedAt = time.Now()
	unit.CreatedBy = createdBy

	// We need a Create method in UnitRepository.
	// I only implemented FindByID and Find in UnitRepository.
	// I need to add Create, Update, Delete methods to UnitRepository as well!

	// Assuming UnitRepository has Create method (I need to add it)
	if err := s.unitRepo.Create(ctx, unit); err != nil {
		return nil, fmt.Errorf("failed to create unit: %w", err)
	}

	return unit, nil
}

func (s *UnitService) GetUnit(ctx context.Context, id primitive.ObjectID) (*models.Unit, error) {
	unit, err := s.unitRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("unit not found: %w", err)
	}
	return unit, nil
}

func (s *UnitService) GetUnits(ctx context.Context, unitType *string, activeOnly bool, ids []primitive.ObjectID) ([]*models.Unit, error) {
	filters := make(map[string]interface{})
	if unitType != nil {
		filters["unit_type"] = *unitType
	}
	if len(ids) > 0 {
		filters["ids"] = ids
	}

	units, err := s.unitRepo.Find(ctx, filters, activeOnly)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch units: %w", err)
	}
	return units, nil
}

func (s *UnitService) UpdateUnit(ctx context.Context, id primitive.ObjectID, req UpdateUnitRequest, updatedBy primitive.ObjectID) (*models.Unit, error) {
	unit, err := s.unitRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("unit not found: %w", err)
	}

	updates := make(map[string]interface{})
	if req.Name != nil {
		unit.Name = *req.Name
		updates["name"] = *req.Name
	}
	if req.Code != nil {
		unit.Code = *req.Code
		updates["code"] = *req.Code
	}
	if req.Description != nil {
		unit.Description = *req.Description
		updates["description"] = *req.Description
	}
	if req.UnitType != nil {
		unit.UnitType = *req.UnitType
		updates["unit_type"] = *req.UnitType
	}
	if req.IsBaseUnit != nil {
		unit.IsBaseUnit = *req.IsBaseUnit
		updates["is_base_unit"] = *req.IsBaseUnit
	}
	if req.IsActive != nil {
		unit.IsActive = *req.IsActive
		updates["is_active"] = *req.IsActive
	}

	if len(updates) > 0 {
		updates["updated_by"] = updatedBy
		// Assuming UnitRepository has Update method
		if err := s.unitRepo.Update(ctx, id, updates); err != nil {
			return nil, fmt.Errorf("failed to update unit: %w", err)
		}
	}

	return unit, nil
}

func (s *UnitService) DeleteUnit(ctx context.Context, id primitive.ObjectID) error {
	// Assuming UnitRepository has Delete (soft delete presumably) method
	if err := s.unitRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete unit: %w", err)
	}
	return nil
}
