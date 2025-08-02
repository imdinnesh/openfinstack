package clients

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Client struct {
	BaseURL    string
	HTTPClient *http.Client
}

type ProfileResponse struct {
	ID    uint   `json:"id"`
	Email string `json:"email"`
}

func NewClient(baseURL string) *Client {
	return &Client{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: 3 * time.Second,
		},
	}
}

func (c *Client) GetUserProfile(userID uint) (*ProfileResponse, error) {
	url := fmt.Sprintf("%s/api/v1/auth/profile", c.BaseURL)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("X-User-ID", fmt.Sprintf("%d", userID))
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	var profileResponse ProfileResponse
	if err := json.NewDecoder(resp.Body).Decode(&profileResponse); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	fmt.Println("Profile Response:", profileResponse)
	return &profileResponse, nil
}
