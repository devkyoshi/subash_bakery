package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/yourusername/erp-system/services/inventory-service/internal/service"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RPCHandler struct {
	stockLevelService *service.StockLevelService
}

func NewRPCHandler(stockLevelService *service.StockLevelService) *RPCHandler {
	return &RPCHandler{
		stockLevelService: stockLevelService,
	}
}

type RPCRequest struct {
	OrganizationID string                 `json:"organization_id"`
	Params         map[string]interface{} `json:"params,omitempty"`
}

func (h *RPCHandler) HandleDashboardStats(body []byte) (interface{}, error) {
	var req RPCRequest
	if err := json.Unmarshal(body, &req); err != nil {
		return nil, fmt.Errorf("invalid request body: %v", err)
	}

	orgID, err := primitive.ObjectIDFromHex(req.OrganizationID)
	if err != nil {
		return nil, fmt.Errorf("invalid organization ID: %v", err)
	}

	stats, err := h.stockLevelService.GetDashboardStats(context.Background(), orgID)
	if err != nil {
		log.Printf("Failed to get dashboard stats: %v", err)
		return nil, err
	}

	return stats, nil
}
