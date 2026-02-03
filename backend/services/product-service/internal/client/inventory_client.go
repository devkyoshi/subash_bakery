package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/yourusername/erp-system/services/product-service/config"
	"github.com/yourusername/erp-system/shared/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type InventoryClient struct {
	BaseURL    string
	HTTPClient *http.Client
}

func NewInventoryClient(cfg *config.Config) *InventoryClient {
	return &InventoryClient{
		BaseURL: cfg.InventoryServiceURL, // Ensure this exists in Config
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

type StockBatchResponse struct {
	Success bool                          `json:"success"`
	Data    map[string]*models.StockLevel `json:"data"`
	Message string                        `json:"message"`
}

func (c *InventoryClient) GetStockLevels(ctx context.Context, productIDs []primitive.ObjectID) (map[string]*models.StockLevel, error) {
	if len(productIDs) == 0 {
		return map[string]*models.StockLevel{}, nil
	}

	url := fmt.Sprintf("%s/api/v1/inventory/stock/bulk", c.BaseURL)

	idStrings := make([]string, len(productIDs))
	for i, id := range productIDs {
		idStrings[i] = id.Hex()
	}

	reqBody := map[string]interface{}{
		"product_ids": idStrings,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	// Add auth header if token is provided
	if req.Header.Get("Authorization") == "" {
		// Log warning or handle specific S2S auth if implemented
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("inventory service returned status: %d", resp.StatusCode)
	}

	var result StockBatchResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.Data, nil
}

type CreateAdjustmentRequest struct {
	LocationID string           `json:"location_id"`
	Reason     string           `json:"reason"`
	Items      []AdjustmentItem `json:"items"`
}

type AdjustmentItem struct {
	ProductID   string  `json:"product_id"`
	ExpectedQty float64 `json:"expected_qty"`
	ActualQty   float64 `json:"actual_qty"`
	UnitCost    float64 `json:"unit_cost"`
}

type AdjustmentResponse struct {
	Success bool                   `json:"success"`
	Data    map[string]interface{} `json:"data"` // We just need the ID really
}

func (c *InventoryClient) CreateStockAdjustment(ctx context.Context, token string, orgID primitive.ObjectID, reqBody CreateAdjustmentRequest) (string, error) {
	url := fmt.Sprintf("%s/api/v1/inventory/organizations/%s/stock-adjustments", c.BaseURL, orgID.Hex())

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", token)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("inventory service returned status: %d", resp.StatusCode)
	}

	var result AdjustmentResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	// Extract ID from map
	if id, ok := result.Data["id"].(string); ok {
		return id, nil
	}
	return "", fmt.Errorf("failed to parse adjustment ID")
}

func (c *InventoryClient) ApproveStockAdjustment(ctx context.Context, token string, adjustmentID string) error {
	url := fmt.Sprintf("%s/api/v1/inventory/stock-adjustments/%s/approve", c.BaseURL, adjustmentID)

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer([]byte("{}")))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", token)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("inventory service returned status: %d", resp.StatusCode)
	}

	return nil
}
