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

type OrgClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewOrgClient(cfg *config.Config) *OrgClient {
	return &OrgClient{
		baseURL: cfg.OrgServiceURL,
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

// GetLocationsByOrganization fetches all locations for an organization
func (c *OrgClient) GetLocationsByOrganization(ctx context.Context, orgID primitive.ObjectID, token string) (map[string]*models.Location, error) {
	url := fmt.Sprintf("%s/api/v1/organizations/%s/locations?limit=1000", c.baseURL, orgID.Hex())
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
		return nil, fmt.Errorf("failed to fetch locations: status %d", resp.StatusCode)
	}

	var response struct {
		Success bool `json:"success"`
		Data    struct {
			Data []*models.Location `json:"data"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	result := make(map[string]*models.Location)
	if response.Success {
		for _, loc := range response.Data.Data {
			result[loc.ID.Hex()] = loc
		}
	}
	return result, nil
}

// GetLocationsBatch fetches multiple locations by IDs
func (c *OrgClient) GetLocationsBatch(ctx context.Context, ids []primitive.ObjectID, token string) (map[string]*models.Location, error) {
	if len(ids) == 0 {
		return map[string]*models.Location{}, nil
	}

	// For now, we'll fetch one by one concurrently as OrgService might not have a bulk endpoint yet
	// Or we can use GetLocationsByOrganization if they belong to same org (likely) but that fetches ALL.
	// Optimally, we should implement a bulk fetch in OrgService.
	// But given constraints, I'll use concurrency here similar to UserClient.

	result := make(map[string]*models.Location)
	var mu sync.Mutex
	var wg sync.WaitGroup
	errChan := make(chan error, len(ids))

	// Deduplicate
	uniqueIDs := make(map[string]primitive.ObjectID)
	for _, id := range ids {
		uniqueIDs[id.Hex()] = id
	}

	for _, id := range uniqueIDs {
		wg.Add(1)
		go func(id primitive.ObjectID) {
			defer wg.Done()
			loc, err := c.GetLocation(ctx, id, token)
			if err != nil {
				// Log error but continue
				fmt.Printf("Failed to fetch location %s: %v\n", id.Hex(), err)
				return
			}
			mu.Lock()
			result[id.Hex()] = loc
			mu.Unlock()
		}(id)
	}

	wg.Wait()
	close(errChan)

	return result, nil
}

// GetLocation fetches a single location
func (c *OrgClient) GetLocation(ctx context.Context, id primitive.ObjectID, token string) (*models.Location, error) {
	url := fmt.Sprintf("%s/api/v1/locations/%s", c.baseURL, id.Hex())
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
		return nil, fmt.Errorf("failed to fetch location: status %d", resp.StatusCode)
	}

	var response struct {
		Success bool             `json:"success"`
		Data    *models.Location `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	if response.Success {
		return response.Data, nil
	}

	return nil, fmt.Errorf("failed to fetch location details")
}
