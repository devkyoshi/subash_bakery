package handlers

import (
	"context"
	"encoding/json"
	"log"

	"github.com/yourusername/erp-system/services/product-service/internal/service"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RPCHandler struct {
	productService *service.ProductService
}

func NewRPCHandler(productService *service.ProductService) *RPCHandler {
	return &RPCHandler{
		productService: productService,
	}
}

type RPCRequest struct {
	OrganizationID string                 `json:"organization_id"`
	Params         map[string]interface{} `json:"params,omitempty"`
}

func (h *RPCHandler) HandleDashboardStats(payload []byte) (interface{}, error) {
	var req RPCRequest
	if err := json.Unmarshal(payload, &req); err != nil {
		return nil, err
	}

	orgID, err := primitive.ObjectIDFromHex(req.OrganizationID)
	if err != nil {
		return nil, err
	}

	// stats is map[string]int64
	stats, err := h.productService.GetDashboardStats(context.Background(), orgID)
	if err != nil {
		log.Printf("Error fetching dashboard stats: %v", err)
		return nil, err
	}

	return stats, nil
}
