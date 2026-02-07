package service

import (
	"context"
	"fmt"
	"time"

	"github.com/yourusername/erp-system/services/inventory-service/internal/repository"
	"github.com/yourusername/erp-system/shared/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UnitChartService struct {
	unitChartRepo *repository.UnitChartRepository
	unitRepo      *repository.UnitRepository
}

func NewUnitChartService(
	unitChartRepo *repository.UnitChartRepository,
	unitRepo *repository.UnitRepository,
) *UnitChartService {
	return &UnitChartService{
		unitChartRepo: unitChartRepo,
		unitRepo:      unitRepo,
	}
}

type CreateUnitChartRequest struct {
	FromUnitID     primitive.ObjectID     `json:"from_unit_id" binding:"required"`
	ToUnitID       primitive.ObjectID     `json:"to_unit_id" binding:"required"`
	ConversionRate float64                `json:"conversion_rate" binding:"required,gt=0"`
	Metadata       map[string]interface{} `json:"metadata"`
}

type UpdateUnitChartRequest struct {
	ConversionRate *float64               `json:"conversion_rate" binding:"omitempty,gt=0"`
	IsActive       *bool                  `json:"is_active"`
	Metadata       map[string]interface{} `json:"metadata"`
}

func (s *UnitChartService) CreateUnitChart(
	ctx context.Context,
	req CreateUnitChartRequest,
	createdBy primitive.ObjectID,
) (*models.UnitChart, error) {

	// Validate units exist
	fromUnit, err := s.unitRepo.FindByID(ctx, req.FromUnitID)
	if err != nil {
		return nil, fmt.Errorf("from_unit not found: %w", err)
	}

	toUnit, err := s.unitRepo.FindByID(ctx, req.ToUnitID)
	if err != nil {
		return nil, fmt.Errorf("to_unit not found: %w", err)
	}

	// Validate same unit type
	if fromUnit.UnitType != toUnit.UnitType {
		return nil, fmt.Errorf("units must be of the same type (from: %s, to: %s)", fromUnit.UnitType, toUnit.UnitType)
	}

	// Prevent self-conversion
	if req.FromUnitID == req.ToUnitID {
		return nil, fmt.Errorf("cannot create conversion from a unit to itself")
	}

	// Check for existing conversion
	existing, _ := s.unitChartRepo.FindByUnits(ctx, req.FromUnitID, req.ToUnitID)
	if existing != nil {
		return nil, fmt.Errorf("conversion already exists between these units")
	}

	// Check for circular path
	pathExists, err := s.unitChartRepo.PathExists(ctx, req.ToUnitID, req.FromUnitID)
	if err != nil {
		return nil, fmt.Errorf("failed to check circular path: %w", err)
	}
	if pathExists {
		return nil, fmt.Errorf("circular conversion path detected")
	}

	chart := &models.UnitChart{
		FromUnitID:     req.FromUnitID,
		ToUnitID:       req.ToUnitID,
		ConversionRate: req.ConversionRate,
		IsActive:       true,
		Metadata:       req.Metadata,
	}

	chart.BaseModel.ID = primitive.NewObjectID()
	chart.BaseModel.CreatedAt = time.Now()
	chart.BaseModel.UpdatedAt = time.Now()
	chart.BaseModel.Version = 1
	chart.BaseModel.CreatedBy = createdBy

	if err := s.unitChartRepo.Create(ctx, chart); err != nil {
		return nil, fmt.Errorf("failed to create unit chart: %w", err)
	}

	return chart, nil
}

func (s *UnitChartService) GetUnitChart(
	ctx context.Context,
	chartID primitive.ObjectID,
) (*models.UnitChart, error) {

	chart, err := s.unitChartRepo.FindByID(ctx, chartID)
	if err != nil {
		return nil, fmt.Errorf("unit chart not found: %w", err)
	}

	return chart, nil
}

type UnitConversionResponse struct {
	ToUnitID       primitive.ObjectID `json:"to_unit_id"`
	ToUnitName     string             `json:"to_unit_name"`
	ToUnitCode     string             `json:"to_unit_code"`
	ConversionRate float64            `json:"conversion_rate"`
}

type UnitResponse struct {
	ID         primitive.ObjectID     `json:"id"`
	ChartID    primitive.ObjectID     `json:"chart_id"`
	Name       string                 `json:"name"`
	Code       string                 `json:"code"`
	UnitType   string                 `json:"unit_type"`
	IsBaseUnit bool                   `json:"is_base_unit"`
	IsActive   bool                   `json:"is_active"`
	Conversion UnitConversionResponse `json:"conversion"`
}

func (s *UnitChartService) GetUnitCharts(
	ctx context.Context,
	activeOnly bool,
) ([]UnitResponse, error) {

	// Fetch all units
	units, err := s.unitRepo.Find(ctx, nil, activeOnly)
	if err != nil {
		return nil, fmt.Errorf("failed to get units: %w", err)
	}

	// Fetch all charts
	charts, err := s.unitChartRepo.Find(ctx, activeOnly)
	if err != nil {
		return nil, fmt.Errorf("failed to get unit charts: %w", err)
	}

	// Create a map of units for quick lookup (by ID)
	unitMap := make(map[primitive.ObjectID]*models.Unit)
	for _, u := range units {
		unitMap[u.ID] = u
	}

	response := make([]UnitResponse, 0)

	// Iterate over charts to build the response
	// This naturally filters out units with no conversions, as we drive off the charts list
	for _, chart := range charts {
		if chart.FromUnitID.IsZero() || chart.ToUnitID.IsZero() {
			continue
		}

		fromUnit, fromOk := unitMap[chart.FromUnitID]
		toUnit, toOk := unitMap[chart.ToUnitID]

		// Ensure both units exist
		if !fromOk || !toOk {
			continue
		}

		conv := UnitConversionResponse{
			ToUnitID:       toUnit.ID,
			ToUnitName:     toUnit.Name,
			ToUnitCode:     toUnit.Code,
			ConversionRate: chart.ConversionRate,
		}

		response = append(response, UnitResponse{
			ID:         fromUnit.ID,
			ChartID:    chart.BaseModel.ID, // Map the Chart ID
			Name:       fromUnit.Name,
			Code:       fromUnit.Code,
			UnitType:   fromUnit.UnitType,
			IsBaseUnit: fromUnit.IsBaseUnit,
			IsActive:   chart.IsActive,
			Conversion: conv,
		})
	}

	return response, nil
}

func (s *UnitChartService) UpdateUnitChart(
	ctx context.Context,
	chartID primitive.ObjectID,
	req UpdateUnitChartRequest,
	updatedBy primitive.ObjectID,
) (*models.UnitChart, error) {

	chart, err := s.unitChartRepo.FindByID(ctx, chartID)
	if err != nil {
		return nil, fmt.Errorf("unit chart not found: %w", err)
	}

	// Update fields if provided
	if req.ConversionRate != nil {
		chart.ConversionRate = *req.ConversionRate
	}

	if req.IsActive != nil {
		chart.IsActive = *req.IsActive
	}

	if req.Metadata != nil {
		chart.Metadata = req.Metadata
	}

	chart.UpdatedBy = updatedBy
	chart.UpdatedAt = time.Now()

	if err := s.unitChartRepo.Update(ctx, chart); err != nil {
		return nil, fmt.Errorf("failed to update unit chart: %w", err)
	}

	return chart, nil
}

func (s *UnitChartService) DeleteUnitChart(
	ctx context.Context,
	chartID primitive.ObjectID,
) error {

	chart, err := s.unitChartRepo.FindByID(ctx, chartID)
	if err != nil {
		return fmt.Errorf("unit chart not found: %w", err)
	}

	if err := s.unitChartRepo.SoftDelete(ctx, chart.ID); err != nil {
		return fmt.Errorf("failed to delete unit chart: %w", err)
	}

	return nil
}

// GetConversionRate retrieves the conversion rate between two units
func (s *UnitChartService) GetConversionRate(
	ctx context.Context,
	fromUnitID, toUnitID primitive.ObjectID,
) (float64, error) {

	if fromUnitID == toUnitID {
		return 1.0, nil
	}

	chart, err := s.unitChartRepo.FindByUnits(ctx, fromUnitID, toUnitID)
	if err != nil {
		return 0, fmt.Errorf("conversion not found: %w", err)
	}

	if !chart.IsActive {
		return 0, fmt.Errorf("conversion is inactive")
	}

	return chart.ConversionRate, nil
}
