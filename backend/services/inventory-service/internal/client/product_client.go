package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/yourusername/erp-system/services/inventory-service/config"
	"github.com/yourusername/erp-system/shared/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProductClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewProductClient(cfg *config.Config) *ProductClient {
	return &ProductClient{
		baseURL: cfg.ProductServiceURL,
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

// GetProduct fetches a single product
func (c *ProductClient) GetProduct(ctx context.Context, id primitive.ObjectID, token string) (*models.Product, error) {
	url := fmt.Sprintf("%s/api/v1/products/%s", c.baseURL, id.Hex())
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", token)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch product: status %d", resp.StatusCode)
	}

	var response struct {
		Success bool            `json:"success"`
		Data    *models.Product `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	if response.Success {
		return response.Data, nil
	}

	return nil, fmt.Errorf("failed to fetch product details")
}

// FetchProducts details for a list of IDs
// This requires a valid JWT token to authenticate closely with Product Service
func (c *ProductClient) GetProductsBatch(ctx context.Context, productIDs []primitive.ObjectID, token string) (map[string]*models.Product, error) {
	var wg sync.WaitGroup
	var mu sync.Mutex
	results := make(map[string]*models.Product)

	// Limit concurrency
	sem := make(chan struct{}, 10)

	for _, pid := range productIDs {
		wg.Add(1)
		go func(id primitive.ObjectID) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			url := fmt.Sprintf("%s/api/v1/products/%s", c.baseURL, id.Hex())
			req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
			if err != nil {
				return
			}
			req.Header.Set("Authorization", token)
			// Add internal header if needed to bypass org check?
			// Handler checks `c.Get("organization_id")` which is set by Middleware from Token.
			// So Token is sufficient.

			resp, err := c.httpClient.Do(req)
			if err != nil {
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				return
			}

			// Decodes into utils.Response { success, data, ... }
			var response struct {
				Success bool            `json:"success"`
				Data    *models.Product `json:"data"`
			}

			if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
				return
			}

			if response.Success && response.Data != nil {
				mu.Lock()
				results[id.Hex()] = response.Data
				mu.Unlock()
			}
		}(pid)
	}

	wg.Wait()
	return results, nil
}
