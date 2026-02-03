package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/yourusername/erp-system/shared/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CompanyRepository struct {
	collection *mongo.Collection
}

func NewCompanyRepository(db *mongo.Database) *CompanyRepository {
	return &CompanyRepository{
		collection: db.Collection("companies"),
	}
}

// Create creates a new company
func (r *CompanyRepository) Create(ctx context.Context, company *models.Company) error {
	company.ID = primitive.NewObjectID()
	company.CreatedAt = time.Now()
	company.UpdatedAt = time.Now()

	_, err := r.collection.InsertOne(ctx, company)
	if err != nil {
		return fmt.Errorf("failed to create company: %w", err)
	}
	return nil
}

// FindByID finds a company by ID
func (r *CompanyRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*models.Company, error) {
	var company models.Company
	filter := bson.M{
		"_id":        id,
		"deleted_at": nil,
	}

	err := r.collection.FindOne(ctx, filter).Decode(&company)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find company: %w", err)
	}
	return &company, nil
}

// FindByOrganization finds all companies for an organization
func (r *CompanyRepository) FindByOrganization(ctx context.Context, orgID primitive.ObjectID, page, limit int) ([]*models.Company, int64, error) {
	filter := bson.M{
		"organization_id": orgID,
		"deleted_at":      nil,
	}

	// Get total count
	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count companies: %w", err)
	}

	// Get paginated results
	skip := (page - 1) * limit
	opts := options.Find().
		SetSkip(int64(skip)).
		SetLimit(int64(limit)).
		SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to find companies: %w", err)
	}
	defer cursor.Close(ctx)

	var companies []*models.Company
	if err := cursor.All(ctx, &companies); err != nil {
		return nil, 0, fmt.Errorf("failed to decode companies: %w", err)
	}

	return companies, total, nil
}

// Update updates a company
func (r *CompanyRepository) Update(ctx context.Context, company *models.Company) error {
	company.UpdatedAt = time.Now()
	company.Version++

	filter := bson.M{
		"_id":     company.ID,
		"version": company.Version - 1,
	}

	update := bson.M{"$set": company}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update company: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("company not found or version conflict")
	}

	return nil
}

// SoftDelete soft deletes a company
func (r *CompanyRepository) SoftDelete(ctx context.Context, id, deletedBy primitive.ObjectID) error {
	now := time.Now()
	filter := bson.M{"_id": id}
	update := bson.M{
		"$set": bson.M{
			"deleted_at": now,
			"deleted_by": deletedBy,
			"updated_at": now,
			"is_active":  false,
		},
	}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to soft delete company: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("company not found")
	}

	return nil
}

// CodeExists checks if a company code exists within an organization
func (r *CompanyRepository) CodeExists(ctx context.Context, orgID primitive.ObjectID, code string, excludeID *primitive.ObjectID) (bool, error) {
	filter := bson.M{
		"organization_id": orgID,
		"code":            code,
		"deleted_at":      nil,
	}

	if excludeID != nil {
		filter["_id"] = bson.M{"$ne": *excludeID}
	}

	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return false, fmt.Errorf("failed to check code existence: %w", err)
	}

	return count > 0, nil
}

// CountByOrganization counts companies for an organization
func (r *CompanyRepository) CountByOrganization(ctx context.Context, orgID primitive.ObjectID) (int, error) {
	filter := bson.M{
		"organization_id": orgID,
		"deleted_at":      nil,
	}

	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, fmt.Errorf("failed to count companies: %w", err)
	}

	return int(count), nil
}

// FindByIDs finds multiple companies by their IDs
func (r *CompanyRepository) FindByIDs(ctx context.Context, companyIDs []primitive.ObjectID) ([]*models.Company, error) {
	if len(companyIDs) == 0 {
		return []*models.Company{}, nil
	}

	filter := bson.M{
		"_id":        bson.M{"$in": companyIDs},
		"deleted_at": nil,
	}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to find companies: %w", err)
	}
	defer cursor.Close(ctx)

	var companies []*models.Company
	if err := cursor.All(ctx, &companies); err != nil {
		return nil, fmt.Errorf("failed to decode companies: %w", err)
	}

	return companies, nil
}
