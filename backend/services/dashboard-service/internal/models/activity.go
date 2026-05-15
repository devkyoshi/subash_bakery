package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ActivityType string

const (
	ActivityTypePO  ActivityType = "purchase_order"
	ActivityTypeGRN ActivityType = "grn"
	ActivityTypeInv ActivityType = "inventory"
)

type ActivityAction string

const (
	ActionCreate  ActivityAction = "create"
	ActionUpdate  ActivityAction = "update"
	ActionDelete  ActivityAction = "delete"
	ActionApprove ActivityAction = "approve"
	ActionReject  ActivityAction = "reject"
)

type Activity struct {
	ID             primitive.ObjectID     `bson:"_id,omitempty" json:"id"`
	OrganizationID primitive.ObjectID     `bson:"organization_id" json:"organization_id"`
	Type           ActivityType           `bson:"type" json:"type"`
	Action         ActivityAction         `bson:"action" json:"action"`
	EntityID       primitive.ObjectID     `bson:"entity_id" json:"entity_id"`
	EntityCode     string                 `bson:"entity_code" json:"entity_code"` // e.g., PO-123
	Description    string                 `bson:"description" json:"description"`
	Metadata       map[string]interface{} `bson:"metadata,omitempty" json:"metadata,omitempty"`
	CreatedBy      primitive.ObjectID     `bson:"created_by" json:"created_by"`
	CreatedByName  string                 `bson:"created_by_name" json:"created_by_name"`
	CreatedAt      time.Time              `bson:"created_at" json:"created_at"`
}
