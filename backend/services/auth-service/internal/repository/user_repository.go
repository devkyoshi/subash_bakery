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

type UserRepository struct {
	collection *mongo.Collection
}

func NewUserRepository(db *mongo.Database) *UserRepository {
	return &UserRepository{
		collection: db.Collection("users"),
	}
}

// Create creates a new user
func (r *UserRepository) Create(ctx context.Context, user *models.User) error {
	user.ID = primitive.NewObjectID()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	_, err := r.collection.InsertOne(ctx, user)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

// FindByEmail finds a user by email
func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := r.collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find user by email: %w", err)
	}
	return &user, nil
}

// FindByID finds a user by ID
func (r *UserRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*models.User, error) {
	var user models.User
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find user by ID: %w", err)
	}
	return &user, nil
}

// FindByGoogleID finds a user by Google ID
func (r *UserRepository) FindByGoogleID(ctx context.Context, googleID string) (*models.User, error) {
	var user models.User
	err := r.collection.FindOne(ctx, bson.M{"google_id": googleID}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find user by Google ID: %w", err)
	}
	return &user, nil
}

// Update updates a user
func (r *UserRepository) Update(ctx context.Context, user *models.User) error {
	user.UpdatedAt = time.Now()
	filter := bson.M{"_id": user.ID}
	update := bson.M{"$set": user}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	return nil
}

// UpdateLastLogin updates the last login timestamp
func (r *UserRepository) UpdateLastLogin(ctx context.Context, userID primitive.ObjectID) error {
	filter := bson.M{"_id": userID}
	update := bson.M{
		"$set": bson.M{
			"last_login": time.Now(),
			"updated_at": time.Now(),
		},
	}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update last login: %w", err)
	}
	return nil
}

// EmailExists checks if an email already exists
func (r *UserRepository) EmailExists(ctx context.Context, email string) (bool, error) {
	count, err := r.collection.CountDocuments(ctx, bson.M{"email": email})
	if err != nil {
		return false, fmt.Errorf("failed to check email existence: %w", err)
	}
	return count > 0, nil
}

// FindAll finds all users with pagination and filtering
func (r *UserRepository) FindAll(ctx context.Context, orgID primitive.ObjectID, filters map[string]interface{}, page, limit int) ([]*models.User, int64, error) {
	filter := bson.M{
		"organization_id": orgID,
	}

	// Apply filters
	if search, ok := filters["search"].(string); ok && search != "" {
		filter["$or"] = []bson.M{
			{"first_name": bson.M{"$regex": search, "$options": "i"}},
			{"last_name": bson.M{"$regex": search, "$options": "i"}},
			{"email": bson.M{"$regex": search, "$options": "i"}},
		}
	}
	if roleID, ok := filters["role_id"].(primitive.ObjectID); ok {
		filter["role_id"] = roleID
	}
	if status, ok := filters["status"].(string); ok && status != "" {
		if status == "active" {
			filter["is_active"] = true
		} else if status == "inactive" {
			filter["is_active"] = false
		}
	}

	// Count total documents
	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count users: %w", err)
	}

	// Find with pagination
	opts := options.Find().
		SetSort(bson.D{{Key: "created_at", Value: -1}}).
		SetSkip(int64((page - 1) * limit)).
		SetLimit(int64(limit))

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to find users: %w", err)
	}
	defer cursor.Close(ctx)

	var users []*models.User
	if err := cursor.All(ctx, &users); err != nil {
		return nil, 0, fmt.Errorf("failed to decode users: %w", err)
	}

	return users, total, nil
}

// Delete permanently deletes a user
func (r *UserRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	filter := bson.M{"_id": id}

	result, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}
