package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Notification struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID         primitive.ObjectID `bson:"user_id" json:"user_id"`
	OrganizationID primitive.ObjectID `bson:"organization_id" json:"organization_id"`
	Title          string             `bson:"title" json:"title"`
	Body           string             `bson:"body" json:"body"`
	Type           string             `bson:"type" json:"type"`
	Data           map[string]string  `bson:"data" json:"data"`
	IsRead         bool               `bson:"is_read" json:"is_read"`
	CreatedAt      time.Time          `bson:"created_at" json:"created_at"`
}
