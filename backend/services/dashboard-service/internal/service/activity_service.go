package service

import (
	"context"
	"time"

	"github.com/yourusername/erp-system/services/dashboard-service/internal/models"
	"github.com/yourusername/erp-system/services/dashboard-service/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ActivityService struct {
	repo *repository.ActivityRepository
}

func NewActivityService(repo *repository.ActivityRepository) *ActivityService {
	return &ActivityService{repo: repo}
}

func (s *ActivityService) RecordActivity(
	ctx context.Context,
	orgID primitive.ObjectID,
	activityType models.ActivityType,
	action models.ActivityAction,
	entityID primitive.ObjectID,
	entityCode string,
	description string,
	createdBy primitive.ObjectID,
	createdByName string,
	metadata map[string]interface{},
) error {
	activity := &models.Activity{
		OrganizationID: orgID,
		Type:           activityType,
		Action:         action,
		EntityID:       entityID,
		EntityCode:     entityCode,
		Description:    description,
		CreatedBy:      createdBy,
		CreatedByName:  createdByName,
		CreatedAt:      time.Now(),
		Metadata:       metadata,
	}

	return s.repo.Create(ctx, activity)
}

func (s *ActivityService) GetRecentActivities(ctx context.Context, orgID primitive.ObjectID, limit int64) ([]*models.Activity, error) {
	return s.repo.GetRecent(ctx, orgID, limit)
}
