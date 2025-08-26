package clients

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
)

// LedgerClient is the interface for interacting with the Ledger service.
type LedgerClient struct {
	BaseURL    string
	HTTPClient *http.Client
}

// LedgerTransactionResponse represents the response from a successful ledger transaction creation.
type LedgerTransactionResponse struct {
	ID uuid.UUID `json:"id"`
}

// NewLedgerClient creates a new LedgerClient instance.
func NewLedgerClient(baseURL string) *LedgerClient {
	return &LedgerClient{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

// CreateTransaction sends a request to the Ledger service to create a new transaction.
func (c *LedgerClient) CreateTransaction(ctx context.Context, reqBody map[string]interface{}) (uuid.UUID, error) {
	url := fmt.Sprintf("%s/ledger/transaction", c.BaseURL)
	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to make API call: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return uuid.Nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var ledgerResp LedgerTransactionResponse
	if err := json.NewDecoder(resp.Body).Decode(&ledgerResp); err != nil {
		return uuid.Nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return ledgerResp.ID, nil
}