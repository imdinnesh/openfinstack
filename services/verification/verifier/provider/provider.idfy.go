package provider

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/imdinnesh/openfinstack/services/verifications/verifier"
)

type IDfyVerifier struct {
	APIKey     string
	BaseURL    string
	HttpClient *http.Client
}

func NewIDfyVerifier(apiKey, baseURL string) *IDfyVerifier {
	return &IDfyVerifier{
		APIKey:     apiKey,
		BaseURL:    baseURL,
		HttpClient: &http.Client{},
	}
}

func (v *IDfyVerifier) Verify(ctx context.Context, input verifier.VerificationInput) (verifier.VerificationResult, error) {
	reqBody := map[string]string{
		"document_type": input.DocumentType,
		"document_url":  input.DocumentURL,
	}

	bodyJSON, _ := json.Marshal(reqBody)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, v.BaseURL+"/verify", strings.NewReader(string(bodyJSON)))
	if err != nil {
		return verifier.VerificationResult{}, err
	}

	req.Header.Set("Authorization", "Bearer "+v.APIKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := v.HttpClient.Do(req)
	if err != nil {
		return verifier.VerificationResult{}, err
	}
	defer resp.Body.Close()

	var result struct {
		Status string `json:"status"`
		Reason string `json:"reason,omitempty"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return verifier.VerificationResult{}, err
	}

	switch result.Status {
	case "verified":
		return verifier.VerificationResult{Verified: true}, nil
	case "rejected":
		return verifier.VerificationResult{Verified: false, RejectReason: result.Reason}, nil
	default:
		return verifier.VerificationResult{}, errors.New("unexpected status from IDfy")
	}
}
