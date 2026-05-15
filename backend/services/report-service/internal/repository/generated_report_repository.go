package repository

import (
	"context"
	"fmt"

	"github.com/yourusername/erp-system/services/report-service/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type GeneratedReportRepository struct {
	collection *mongo.Collection
}

func NewGeneratedReportRepository(db *mongo.Database) *GeneratedReportRepository {
	return &GeneratedReportRepository{
		collection: db.Collection("generated_reports"),
	}
}

// Create stores a new generated report record
func (r *GeneratedReportRepository) Create(ctx context.Context, report *models.GeneratedReport) error {
	_, err := r.collection.InsertOne(ctx, report)
	if err != nil {
		return fmt.Errorf("failed to insert generated report: %w", err)
	}
	return nil
}

// FindByID retrieves a generated report by its ID
func (r *GeneratedReportRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*models.GeneratedReport, error) {
	var report models.GeneratedReport
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&report)
	if err != nil {
		return nil, fmt.Errorf("failed to find generated report: %w", err)
	}
	return &report, nil
}

// FindByOrganization retrieves all generated reports for an organization
func (r *GeneratedReportRepository) FindByOrganization(
	ctx context.Context,
	orgID primitive.ObjectID,
	page, limit int,
) ([]*models.GeneratedReport, int64, error) {
	filter := bson.M{"organization_id": orgID}

	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count generated reports: %w", err)
	}

	skip := int64((page - 1) * limit)
	opts := options.Find().
		SetSort(bson.M{"generated_at": -1}).
		SetSkip(skip).
		SetLimit(int64(limit))

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to find generated reports: %w", err)
	}
	defer cursor.Close(ctx)

	var reports []*models.GeneratedReport
	if err := cursor.All(ctx, &reports); err != nil {
		return nil, 0, fmt.Errorf("failed to decode generated reports: %w", err)
	}

	return reports, total, nil
}
