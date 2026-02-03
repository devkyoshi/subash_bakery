package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// Unit represents a unit of measure (UOM) - Global/System-wide
type Unit struct {
	BaseModel `bson:",inline"`

	Code        string `bson:"code" json:"code" binding:"required"` // PCS, KG, L
	Name        string `bson:"name" json:"name"`                    // Pieces, Kilogram
	Symbol      string `bson:"symbol" json:"symbol"`                // pcs, kg
	Description string `bson:"description" json:"description"`

	// Classification
	UnitType string `bson:"unit_type" json:"unit_type"` // quantity, weight, volume, length

	IsBaseUnit bool `bson:"is_base_unit" json:"is_base_unit"`
	IsActive   bool `bson:"is_active" json:"is_active"`

	Metadata map[string]interface{} `bson:"metadata" json:"metadata"`
}

// UnitChart defines unit conversion rules - Global/System-wide
type UnitChart struct {
	BaseModel `bson:",inline"`

	FromUnitID primitive.ObjectID `bson:"from_unit_id" json:"from_unit_id"`
	ToUnitID   primitive.ObjectID `bson:"to_unit_id" json:"to_unit_id"`

	// Conversion
	ConversionRate float64 `bson:"conversion_rate" json:"conversion_rate"`
	// Example: 1 BOX = 12 PCS → rate = 12

	IsActive bool `bson:"is_active" json:"is_active"`

	Metadata map[string]interface{} `bson:"metadata" json:"metadata"`
}
