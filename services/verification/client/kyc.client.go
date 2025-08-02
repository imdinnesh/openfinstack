package clients

import (
	"bytes"
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

type KYCUpdateResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func NewClient(baseURL string) *Client {
	return &Client{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: 3 * time.Second,
		},
	}
}

func (c *Client) UpdateKYCStatus(id uint, status string, reason *string, adminID uint) (*KYCUpdateResponse, error) {
	url:= fmt.Sprintf("%s/api/v1/kyc/update/%d", c.BaseURL, id)
	
	type UpdateRequest struct {
		Status string  `json:"status" binding:"required"`
		Reason *string `json:"reason"`
	}

	reqBody := UpdateRequest{
		Status: status,
		Reason: reason,
	}
	reqBodyJSON, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBodyJSON))
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-Admin-ID", fmt.Sprintf("%d", adminID))

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("failed to update KYC status")
	}

	var updateResp KYCUpdateResponse
	if err := json.NewDecoder(resp.Body).Decode(&updateResp); err != nil {
		return nil, err
	}

	return &updateResp, nil
}
