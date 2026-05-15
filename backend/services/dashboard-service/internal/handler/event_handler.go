package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/yourusername/erp-system/services/dashboard-service/internal/models"
	"github.com/yourusername/erp-system/services/dashboard-service/internal/service"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type EventHandler struct {
	activityService *service.ActivityService
}

func NewEventHandler(activityService *service.ActivityService) *EventHandler {
	return &EventHandler{activityService: activityService}
}

type ActivityEvent struct {
	OrganizationID string                 `json:"organization_id"`
	Type           string                 `json:"type"`
	Action         string                 `json:"action"`
	EntityID       string                 `json:"entity_id"`
	EntityCode     string                 `json:"entity_code"`
	Description    string                 `json:"description"`
	CreatedBy      string                 `json:"created_by"`
	CreatedByName  string                 `json:"created_by_name"`
	Metadata       map[string]interface{} `json:"metadata,omitempty"`
}

func (h *EventHandler) HandleActivityEvent(body []byte) error {
	var event ActivityEvent
	if err := json.Unmarshal(body, &event); err != nil {
		return fmt.Errorf("failed to unmarshal event: %w", err)
	}

	log.Printf("Received activity event: %s %s", event.Type, event.Action)

	// Convert string IDs to ObjectIDs
	orgID, err := primitive.ObjectIDFromHex(event.OrganizationID)
	if err != nil {
		return fmt.Errorf("invalid organization_id: %w", err)
	}

	entityID, err := primitive.ObjectIDFromHex(event.EntityID)
	if err != nil {
		return fmt.Errorf("invalid entity_id: %w", err)
	}

	createdBy, err := primitive.ObjectIDFromHex(event.CreatedBy)
	if err != nil {
		return fmt.Errorf("invalid created_by: %w", err)
	}

	// Record activity
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return h.activityService.RecordActivity(
		ctx,
		orgID,
		models.ActivityType(event.Type),
		models.ActivityAction(event.Action),
		entityID,
		event.EntityCode,
		event.Description,
		createdBy,
		event.CreatedByName,
		event.Metadata,
	)
}
