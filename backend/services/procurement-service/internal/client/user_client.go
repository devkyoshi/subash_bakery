package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/yourusername/erp-system/services/procurement-service/config"
	"github.com/yourusername/erp-system/shared/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewUserClient(cfg *config.Config) *UserClient {
	return &UserClient{
		baseURL: cfg.AuthServiceURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// GetUser fetches a single user
func (c *UserClient) GetUser(ctx context.Context, id primitive.ObjectID, token string) (*models.User, error) {
	url := fmt.Sprintf("%s/api/v1/auth/users/%s", c.baseURL, id.Hex())
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
		return nil, fmt.Errorf("failed to fetch user: status %d", resp.StatusCode)
	}

	var response struct {
		Success bool         `json:"success"`
		Data    *models.User `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	if response.Success {
		return response.Data, nil
	}

	return nil, fmt.Errorf("failed to fetch user details")
}

// GetUsersBatch fetches users concurrently
func (c *UserClient) GetUsersBatch(ctx context.Context, ids []primitive.ObjectID, token string) (map[string]*models.User, error) {
	result := make(map[string]*models.User)
	var mu sync.Mutex
	var wg sync.WaitGroup
	errChan := make(chan error, len(ids))

	// Deduplicate IDs
	uniqueIDs := make(map[string]primitive.ObjectID)
	for _, id := range ids {
		uniqueIDs[id.Hex()] = id
	}

	for _, id := range uniqueIDs {
		wg.Add(1)
		go func(id primitive.ObjectID) {
			defer wg.Done()
			sem := make(chan struct{}, 5) // Limit concurrency to 5
			sem <- struct{}{}
			defer func() { <-sem }()

			user, err := c.GetUser(ctx, id, token)
			if err != nil {
				// Log error but continue
				fmt.Printf("Failed to fetch user %s: %v\n", id.Hex(), err)
				return
			}
			mu.Lock()
			result[id.Hex()] = user
			mu.Unlock()
		}(id)
	}

	wg.Wait()
	close(errChan)

	return result, nil
}
