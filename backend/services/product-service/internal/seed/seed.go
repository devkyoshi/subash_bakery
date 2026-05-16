package seed

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/yourusername/erp-system/shared/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Seeder struct {
	db *mongo.Database
}

func NewSeeder(db *mongo.Database) *Seeder {
	return &Seeder{db: db}
}

func (s *Seeder) SeedAll(ctx context.Context) error {
	log.Println("Starting product-service seeding...")

	if err := s.SeedUnits(ctx); err != nil {
		return fmt.Errorf("failed to seed units: %w", err)
	}

	if err := s.SeedUnitCharts(ctx); err != nil {
		return fmt.Errorf("failed to seed unit charts: %w", err)
	}

	log.Println("Product-service seeding completed successfully")
	return nil
}

func (s *Seeder) SeedUnits(ctx context.Context) error {
	log.Println("Seeding units...")

	col := s.db.Collection("units")

	count, err := col.CountDocuments(ctx, bson.M{})
	if err != nil {
		return fmt.Errorf("failed to count units: %w", err)
	}
	if count > 0 {
		log.Printf("Units already seeded (%d found), skipping", count)
		return nil
	}

	now := time.Now()
	base := func(id primitive.ObjectID, code, name, symbol, unitType string, isBase bool) models.Unit {
		return models.Unit{
			BaseModel: models.BaseModel{
				ID:        id,
				CreatedAt: now,
				UpdatedAt: now,
				Version:   0,
			},
			Code:       code,
			Name:       name,
			Symbol:     symbol,
			UnitType:   unitType,
			IsBaseUnit: isBase,
			IsActive:   true,
		}
	}

	units := []interface{}{
		// Quantity
		base(pcsID, "PCS", "Pieces", "pcs", "quantity", true),
		base(boxID, "BOX", "Box", "box", "quantity", false),
		base(ctnID, "CTN", "Carton", "ctn", "quantity", false),
		base(dznID, "DZN", "Dozen", "doz", "quantity", false),
		base(pckID, "PCK", "Pack", "pck", "quantity", false),

		// Weight
		base(kgID, "KG", "Kilogram", "kg", "weight", true),
		base(gID, "G", "Gram", "g", "weight", false),
		base(mgID, "MG", "Milligram", "mg", "weight", false),
		base(lbID, "LB", "Pound", "lb", "weight", false),
		base(ozID, "OZ", "Ounce", "oz", "weight", false),

		// Volume
		base(lID, "L", "Liter", "L", "volume", true),
		base(mlID, "ML", "Milliliter", "mL", "volume", false),
		base(galID, "GAL", "Gallon", "gal", "volume", false),

		// Length
		base(mID, "M", "Meter", "m", "length", true),
		base(cmID, "CM", "Centimeter", "cm", "length", false),
		base(mmID, "MM", "Millimeter", "mm", "length", false),
		base(ftID, "FT", "Foot", "ft", "length", false),
		base(inID, "IN", "Inch", "in", "length", false),
	}

	if _, err := col.InsertMany(ctx, units); err != nil {
		return fmt.Errorf("failed to insert units: %w", err)
	}

	log.Printf("Seeded %d units", len(units))
	return nil
}

func (s *Seeder) SeedUnitCharts(ctx context.Context) error {
	log.Println("Seeding unit charts...")

	col := s.db.Collection("unit_charts")

	count, err := col.CountDocuments(ctx, bson.M{})
	if err != nil {
		return fmt.Errorf("failed to count unit charts: %w", err)
	}
	if count > 0 {
		log.Printf("Unit charts already seeded (%d found), skipping", count)
		return nil
	}

	now := time.Now()
	chart := func(from, to primitive.ObjectID, rate float64) models.UnitChart {
		return models.UnitChart{
			BaseModel: models.BaseModel{
				ID:        primitive.NewObjectID(),
				CreatedAt: now,
				UpdatedAt: now,
				Version:   0,
			},
			FromUnitID:     from,
			ToUnitID:       to,
			ConversionRate: rate,
			IsActive:       true,
		}
	}

	charts := []interface{}{
		// Quantity conversions (base: PCS)
		chart(boxID, pcsID, 12),
		chart(ctnID, pcsID, 144),
		chart(ctnID, boxID, 12),
		chart(dznID, pcsID, 12),
		chart(pckID, pcsID, 6),

		// Weight conversions (base: KG)
		chart(gID, kgID, 0.001),
		chart(mgID, kgID, 0.000001),
		chart(mgID, gID, 0.001),
		chart(lbID, kgID, 0.453592),
		chart(ozID, kgID, 0.0283495),
		chart(ozID, lbID, 0.0625),

		// Volume conversions (base: L)
		chart(mlID, lID, 0.001),
		chart(galID, lID, 3.78541),

		// Length conversions (base: M)
		chart(cmID, mID, 0.01),
		chart(mmID, mID, 0.001),
		chart(mmID, cmID, 0.1),
		chart(ftID, mID, 0.3048),
		chart(inID, mID, 0.0254),
		chart(inID, ftID, 0.0833333),
	}

	if _, err := col.InsertMany(ctx, charts); err != nil {
		return fmt.Errorf("failed to insert unit charts: %w", err)
	}

	log.Printf("Seeded %d unit charts", len(charts))
	return nil
}

// Fixed ObjectIDs so unit charts can reference units by known IDs.
var (
	// Quantity
	pcsID = mustID("000000000000000000000001")
	boxID = mustID("000000000000000000000002")
	ctnID = mustID("000000000000000000000003")
	dznID = mustID("000000000000000000000004")
	pckID = mustID("000000000000000000000005")

	// Weight
	kgID = mustID("000000000000000000000010")
	gID  = mustID("000000000000000000000011")
	mgID = mustID("000000000000000000000012")
	lbID = mustID("000000000000000000000013")
	ozID = mustID("000000000000000000000014")

	// Volume
	lID   = mustID("000000000000000000000020")
	mlID  = mustID("000000000000000000000021")
	galID = mustID("000000000000000000000022")

	// Length
	mID  = mustID("000000000000000000000030")
	cmID = mustID("000000000000000000000031")
	mmID = mustID("000000000000000000000032")
	ftID = mustID("000000000000000000000033")
	inID = mustID("000000000000000000000034")
)

func mustID(hex string) primitive.ObjectID {
	id, err := primitive.ObjectIDFromHex(hex)
	if err != nil {
		panic(err)
	}
	return id
}
