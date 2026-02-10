package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/yourusername/erp-system/services/procurement-service/config"
	"github.com/yourusername/erp-system/shared/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type InventoryClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewInventoryClient(cfg *config.Config) *InventoryClient {
	return &InventoryClient{
		baseURL: cfg.InventoryServiceURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// StockMovementRequest represents the request to create a stock movement
type StockMovementRequest struct {
	ProductID      string  `json:"product_id"`
	MovementType   string  `json:"movement_type"`
	FromLocationID string  `json:"from_location_id,omitempty"`
	ToLocationID   string  `json:"to_location_id,omitempty"`
	Quantity       float64 `json:"quantity"`
	UnitCost       float64 `json:"unit_cost"`
	ReferenceType  string  `json:"reference_type"`
	ReferenceNo    string  `json:"reference_no"`
	Reason         string  `json:"reason,omitempty"`
	Notes          string  `json:"notes,omitempty"`
	BatchNumber    string  `json:"batch_number,omitempty"`
}

// CreateStockMovement calls the inventory service to create a stock movement
func (c *InventoryClient) CreateStockMovement(ctx context.Context, orgID primitive.ObjectID, req StockMovementRequest, token string) error {
	url := fmt.Sprintf("%s/api/v1/inventory/organizations/%s/stock-movements", c.baseURL, orgID.Hex())

	body, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	if token != "" {
		httpReq.Header.Set("Authorization", token)
	}

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("failed to call inventory service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("inventory service returned status %d", resp.StatusCode)
	}

	return nil
}

// CreateStockMovementsBatch creates multiple stock movements
func (c *InventoryClient) CreateStockMovementsBatch(ctx context.Context, orgID primitive.ObjectID, requests []StockMovementRequest, token string) error {
	for _, req := range requests {
		if err := c.CreateStockMovement(ctx, orgID, req, token); err != nil {
			// Log error but continue with other movements
			fmt.Printf("Failed to create stock movement for product %s: %v\n", req.ProductID, err)
		}
	}
	return nil
}

// GetUnitsBatch fetches details for a list of unit IDs
func (c *InventoryClient) GetUnitsBatch(ctx context.Context, unitIDs []string, token string) (map[string]*models.Unit, error) {
	if len(unitIDs) == 0 {
		return make(map[string]*models.Unit), nil
	}

	idsParam := ""
	for i, id := range unitIDs {
		if i > 0 {
			idsParam += ","
		}
		idsParam += id
	}

	url := fmt.Sprintf("%s/api/v1/inventory/units?ids=%s", c.baseURL, idsParam)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	if token != "" {
		req.Header.Set("Authorization", token)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch units: status %d", resp.StatusCode)
	}

	var response struct {
		Success bool           `json:"success"`
		Data    []*models.Unit `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	results := make(map[string]*models.Unit)
	if response.Success {
		for _, unit := range response.Data {
			results[unit.ID.Hex()] = unit
		}
	}

	return results, nil
}
