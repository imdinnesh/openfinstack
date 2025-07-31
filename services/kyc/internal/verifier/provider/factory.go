package provider

import (
	"github.com/imdinnesh/openfinstack/services/kyc/config"
	"github.com/imdinnesh/openfinstack/services/kyc/internal/verifier"
)

func NewVerifier(cfg *config.Config) verifier.Verifier{
	switch cfg.KYCVerifier {
		case "idfy":
			return NewIDfyVerifier(cfg.IDFYApiKey,cfg.IDFYBaseURL)
		default:
			return NewMockVerifier()
	}
}