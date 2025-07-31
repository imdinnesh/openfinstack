package provider

import (
	"github.com/imdinnesh/openfinstack/services/verifications/config"
	"github.com/imdinnesh/openfinstack/services/verifications/verifier"
)

func NewVerifier(cfg *config.Config) verifier.Verifier{
	switch cfg.KYCVerifier {
		case "idfy":
			return NewIDfyVerifier(cfg.IDFYApiKey,cfg.IDFYBaseURL)
		default:
			return NewMockVerifier()
	}
}