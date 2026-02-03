package repository

import (
	"context"
	"time"

	"github.com/yourusername/erp-system/shared/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type SerialNumberRepository struct {
	collection *mongo.Collection
}

func NewSerialNumberRepository(db *mongo.Database) *SerialNumberRepository {
	return &SerialNumberRepository{
		collection: db.Collection("serial_numbers"),
	}
}

func (r *SerialNumberRepository) Create(ctx context.Context, serial *models.SerialNumber) error {
	serial.ID = primitive.NewObjectID()
	if serial.CreatedAt.IsZero() {
		serial.CreatedAt = time.Now()
	}
	serial.UpdatedAt = time.Now()
	_, err := r.collection.InsertOne(ctx, serial)
	return err
}

func (r *SerialNumberRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*models.SerialNumber, error) {
	var serial models.SerialNumber
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&serial)
	if err != nil {
		return nil, err
	}
	return &serial, nil
}

func (r *SerialNumberRepository) SerialNoExists(ctx context.Context, orgID primitive.ObjectID, serialNo string, excludeID interface{}) (bool, error) {
	filter := bson.M{
		"organization_id": orgID,
		"serial_no":       serialNo,
	}
	// excludeID logic if we support updating serial no and checking uniqueness
	// The service just passes nil usually for create.

	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *SerialNumberRepository) FindBySerialNo(ctx context.Context, orgID primitive.ObjectID, serialNo string) (*models.SerialNumber, error) {
	filter := bson.M{
		"organization_id": orgID,
		"serial_no":       serialNo,
	}
	var serial models.SerialNumber
	err := r.collection.FindOne(ctx, filter).Decode(&serial)
	if err != nil {
		return nil, err
	}
	return &serial, nil
}

func (r *SerialNumberRepository) FindByProduct(ctx context.Context, productID primitive.ObjectID, filters map[string]interface{}) ([]*models.SerialNumber, error) {
	bsonFilters := bson.M{
		"product_id": productID,
	}
	for k, v := range filters {
		bsonFilters[k] = v
	}

	cursor, err := r.collection.Find(ctx, bsonFilters)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var serials []*models.SerialNumber
	if err = cursor.All(ctx, &serials); err != nil {
		return nil, err
	}
	return serials, nil
}

func (r *SerialNumberRepository) Update(ctx context.Context, id primitive.ObjectID, update primitive.M) error {
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	return err
}

func (r *SerialNumberRepository) Allocate(ctx context.Context, serialID, customerID, salesOrderID primitive.ObjectID) error {
	update := bson.M{
		"$set": bson.M{
			"status":         "allocated",
			"is_available":   false,
			"customer_id":    customerID,
			"sales_order_id": salesOrderID,
			"allocated_at":   time.Now(),
			"updated_at":     time.Now(),
		},
	}
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": serialID}, update)
	return err
}

func (r *SerialNumberRepository) MarkAsSold(ctx context.Context, serialID primitive.ObjectID) error {
	update := bson.M{
		"$set": bson.M{
			"status":       "sold",
			"is_available": false,
			"sold_at":      time.Now(),
			"updated_at":   time.Now(),
		},
	}
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": serialID}, update)
	return err
}

func (r *SerialNumberRepository) Delete(ctx context.Context, id, deletedBy primitive.ObjectID) error {
	update := bson.M{
		"$set": bson.M{
			"deleted_at": time.Now(),
			"deleted_by": deletedBy,
		},
	}
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	return err
}
