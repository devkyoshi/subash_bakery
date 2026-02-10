package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DeviceToken struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID         primitive.ObjectID `bson:"user_id" json:"user_id"`
	OrganizationID primitive.ObjectID `bson:"organization_id" json:"organization_id"`
	Token          string             `bson:"token" json:"token"`
	Platform       string             `bson:"platform" json:"platform"` // ios, android, web
	Name           string             `bson:"name" json:"name"`         // e.g. "iPhone 12"
	LastUsedAt     time.Time          `bson:"last_used_at" json:"last_used_at"`
	CreatedAt      time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt      time.Time          `bson:"updated_at" json:"updated_at"`
}
