package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
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
		Success bool               `json:"success"`
		Data    []*models.Location `json:"data"`
		Total   int                `json:"total"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	result := make(map[string]*models.Location)
	if response.Success {
		for _, loc := range response.Data {
			result[loc.ID.Hex()] = loc
		}
	}
	return result, nil
}
