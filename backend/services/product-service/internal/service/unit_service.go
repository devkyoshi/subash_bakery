package service

import (
	"context"

	"github.com/yourusername/erp-system/services/product-service/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UnitService struct {
	unitRepo *repository.UnitRepository
}

func NewUnitService(unitRepo *repository.UnitRepository) *UnitService {
	return &UnitService{
		unitRepo: unitRepo,
	}
}

type UnitConversionResponse struct {
	ToUnitID       primitive.ObjectID `json:"to_unit_id"`
	ToUnitName     string             `json:"to_unit_name"`
	ToUnitCode     string             `json:"to_unit_code"`
	ConversionRate float64            `json:"conversion_rate"`
}

type UnitResponse struct {
	ID          primitive.ObjectID       `json:"id"`
	Name        string                   `json:"name"`
	Code        string                   `json:"code"`
	Symbol      string                   `json:"symbol"`
	UnitType    string                   `json:"unit_type"`
	IsBaseUnit  bool                     `json:"is_base_unit"`
	IsActive    bool                     `json:"is_active"`
	Conversions []UnitConversionResponse `json:"conversions,omitempty"`
}

// ListUnits returns units without conversion rules
func (s *UnitService) ListUnits(ctx context.Context) ([]UnitResponse, error) {
	// Fetch units
	units, err := s.unitRepo.ListUnits(ctx)
	if err != nil {
		return nil, err
	}

	var response []UnitResponse
	for _, u := range units {
		response = append(response, UnitResponse{
			ID:         u.ID,
			Name:       u.Name,
			Code:       u.Code,
			Symbol:     u.Symbol,
			UnitType:   u.UnitType,
			IsBaseUnit: u.IsBaseUnit,
			IsActive:   u.IsActive,
		})
	}

	return response, nil
}

// ListUnitsWithConversion returns units with their conversion rules
func (s *UnitService) ListUnitsWithConversion(ctx context.Context) ([]UnitResponse, error) {
	// Fetch units
	units, err := s.unitRepo.ListUnits(ctx)
	if err != nil {
		return nil, err
	}

	// Fetch unit charts
	charts, err := s.unitRepo.ListUnitCharts(ctx)
	if err != nil {
		return nil, err
	}

	// Create a map of units for quick lookup
	unitMap := make(map[primitive.ObjectID]string)
	unitCodeMap := make(map[primitive.ObjectID]string)
	for _, u := range units {
		unitMap[u.ID] = u.Name
		unitCodeMap[u.ID] = u.Code
	}

	// Map conversions by FromUnitID
	conversionsMap := make(map[primitive.ObjectID][]UnitConversionResponse)
	for _, chart := range charts {
		if chart.FromUnitID.IsZero() || chart.ToUnitID.IsZero() {
			continue
		}

		conv := UnitConversionResponse{
			ToUnitID:       chart.ToUnitID,
			ToUnitName:     unitMap[chart.ToUnitID],
			ToUnitCode:     unitCodeMap[chart.ToUnitID],
			ConversionRate: chart.ConversionRate,
		}
		conversionsMap[chart.FromUnitID] = append(conversionsMap[chart.FromUnitID], conv)
	}

	var response []UnitResponse
	for _, u := range units {
		response = append(response, UnitResponse{
			ID:          u.ID,
			Name:        u.Name,
			Code:        u.Code,
			Symbol:      u.Symbol,
			UnitType:    u.UnitType,
			IsBaseUnit:  u.IsBaseUnit,
			IsActive:    u.IsActive,
			Conversions: conversionsMap[u.ID],
		})
	}

	return response, nil
}
