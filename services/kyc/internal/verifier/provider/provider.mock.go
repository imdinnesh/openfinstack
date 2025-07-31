package provider

import (
	"context"
	"math/rand"
	"time"

	"github.com/imdinnesh/openfinstack/services/kyc/internal/verifier"
)

type MockVerifier struct{}

func NewMockVerifier() *MockVerifier {
	return &MockVerifier{}
}

func (m *MockVerifier) Verify(ctx context.Context, input verifier.VerificationInput) (verifier.VerificationResult, error) {
	// Simulate latency
	time.Sleep(500 * time.Millisecond)

	// Randomize success/failure
	if rand.Intn(100) < 80 {
		return verifier.VerificationResult{Verified: true}, nil
	}

	return verifier.VerificationResult{
		Verified:     false,
		RejectReason: "Mock: Invalid document",
	}, nil
}
