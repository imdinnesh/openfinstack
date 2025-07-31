package verifier

import "context"

type VerificationInput struct {
	DocumentType string
	DocumentID   string
	DocumentURL  string
}

type VerificationResult struct {
	Verified     bool
	RejectReason string
}

type Verifier interface {
	Verify(ctx context.Context, input VerificationInput) (VerificationResult, error)
}
