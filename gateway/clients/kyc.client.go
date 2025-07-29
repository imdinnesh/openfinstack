package clients
import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type Client struct {
	BaseURL    string
	HTTPClient *http.Client
}

type KYCStatusResponse struct {
	Status string `json:"status"`
}

func NewClient(baseURL string) *Client {
	return &Client{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: 3 * time.Second,
		},
	}
}

func (c *Client) GetStatus(userID uint) (string, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/kyc/status", c.BaseURL), nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("X-User-ID", fmt.Sprintf("%d", userID))

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("failed to get KYC status")
	}

	var result KYCStatusResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	return result.Status, nil
}
