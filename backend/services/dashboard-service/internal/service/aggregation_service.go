package service

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/yourusername/erp-system/services/dashboard-service/internal/models"
	"github.com/yourusername/erp-system/shared/rabbitmq"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AggregationService struct {
	rabbitClient    *rabbitmq.RabbitMQClient
	activityService *ActivityService
}

func NewAggregationService(rabbitClient *rabbitmq.RabbitMQClient, activityService *ActivityService) *AggregationService {
	return &AggregationService{
		rabbitClient:    rabbitClient,
		activityService: activityService,
	}
}

type DashboardOverview struct {
	Inventory   interface{}        `json:"inventory"`
	Procurement interface{}        `json:"procurement"`
	Product     interface{}        `json:"product"`
	Activities  []*models.Activity `json:"activities"`
	Errors      []string           `json:"errors,omitempty"`
}

type RPCRequest struct {
	OrganizationID string                 `json:"organization_id"`
	Params         map[string]interface{} `json:"params,omitempty"`
}

func (s *AggregationService) GetDashboardOverview(ctx context.Context, orgID primitive.ObjectID, token string) (*DashboardOverview, error) {
	var wg sync.WaitGroup
	overview := &DashboardOverview{}
	errors := make([]string, 0)
	var mu sync.Mutex

	rpcReq := RPCRequest{
		OrganizationID: orgID.Hex(),
		Params: map[string]interface{}{
			"auth_token": token,
		},
	}

	// Fetch Inventory Stats
	wg.Add(1)
	go func() {
		defer wg.Done()
		timeoutCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		resp, err := s.rabbitClient.RPCRequest(timeoutCtx, "inventory.dashboard.stats", rpcReq)
		mu.Lock()
		defer mu.Unlock()
		if err != nil {
			errors = append(errors, fmt.Sprintf("Inventory service error: %v", err))
		} else {
			overview.Inventory = resp
		}
	}()

	// Fetch Procurement Stats
	wg.Add(1)
	go func() {
		defer wg.Done()
		timeoutCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		resp, err := s.rabbitClient.RPCRequest(timeoutCtx, "procurement.dashboard.stats", rpcReq)
		mu.Lock()
		defer mu.Unlock()
		if err != nil {
			errors = append(errors, fmt.Sprintf("Procurement service error: %v", err))
		} else {
			overview.Procurement = resp
		}
	}()

	// Fetch Product Stats
	wg.Add(1)
	go func() {
		defer wg.Done()
		timeoutCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		resp, err := s.rabbitClient.RPCRequest(timeoutCtx, "product.dashboard.stats", rpcReq)
		mu.Lock()
		defer mu.Unlock()
		if err != nil {
			errors = append(errors, fmt.Sprintf("Product service error: %v", err))
		} else {
			overview.Product = resp
		}
	}()

	// Fetch Recent Activities (Local DB)
	wg.Add(1)
	go func() {
		defer wg.Done()
		activities, err := s.activityService.GetRecentActivities(ctx, orgID, 10)
		mu.Lock()
		defer mu.Unlock()
		if err != nil {
			errors = append(errors, fmt.Sprintf("Activity service error: %v", err))
		} else {
			overview.Activities = activities
		}
	}()

	wg.Wait()
	overview.Errors = errors
	return overview, nil
}
