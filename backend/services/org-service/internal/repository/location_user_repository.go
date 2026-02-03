package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/yourusername/erp-system/shared/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type LocationUserRepository struct {
	collection *mongo.Collection
}

func NewLocationUserRepository(db *mongo.Database) *LocationUserRepository {
	return &LocationUserRepository{
		collection: db.Collection("location_users"),
	}
}

// Create creates a new location user
func (r *LocationUserRepository) Create(ctx context.Context, locationUser *models.LocationUser) error {
	locationUser.ID = primitive.NewObjectID()
	locationUser.CreatedAt = time.Now()
	locationUser.UpdatedAt = time.Now()

	_, err := r.collection.InsertOne(ctx, locationUser)
	if err != nil {
		return fmt.Errorf("failed to create location user: %w", err)
	}
	return nil
}

// FindLocationsByUserID finds all locations accessible by a user
func (r *LocationUserRepository) FindLocationsByUserID(ctx context.Context, userID primitive.ObjectID) ([]primitive.ObjectID, error) {
	filter := bson.M{
		"user_id":    userID,
		"is_active":  true,
		"deleted_at": nil,
	}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to find location users: %w", err)
	}
	defer cursor.Close(ctx)

	var locationUsers []*models.LocationUser
	if err := cursor.All(ctx, &locationUsers); err != nil {
		return nil, fmt.Errorf("failed to decode location users: %w", err)
	}

	// Extract location IDs
	locationIDs := make([]primitive.ObjectID, 0, len(locationUsers))
	for _, lu := range locationUsers {
		locationIDs = append(locationIDs, lu.LocationID)
	}

	return locationIDs, nil
}

// FindByUserIDAndLocationID finds a location user by user ID and location ID
func (r *LocationUserRepository) FindByUserIDAndLocationID(ctx context.Context, userID, locationID primitive.ObjectID) (*models.LocationUser, error) {
	var locationUser models.LocationUser
	filter := bson.M{
		"user_id":     userID,
		"location_id": locationID,
		"is_active":   true,
		"deleted_at":  nil,
	}

	err := r.collection.FindOne(ctx, filter).Decode(&locationUser)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find location user: %w", err)
	}
	return &locationUser, nil
}

// FindByUserID finds all location user records for a specific user
func (r *LocationUserRepository) FindByUserID(ctx context.Context, userID primitive.ObjectID) ([]*models.LocationUser, error) {
	filter := bson.M{
		"user_id":    userID,
		"is_active":  true,
		"deleted_at": nil,
	}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to find location users: %w", err)
	}
	defer cursor.Close(ctx)

	var locationUsers []*models.LocationUser
	if err := cursor.All(ctx, &locationUsers); err != nil {
		return nil, fmt.Errorf("failed to decode location users: %w", err)
	}

	return locationUsers, nil
}
