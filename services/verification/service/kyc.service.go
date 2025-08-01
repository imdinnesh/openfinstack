package service

import (
	"context"
	"log"

	repository "github.com/imdinnesh/openfinstack/services/verifications/repo"
	"github.com/imdinnesh/openfinstack/services/verifications/verifier"
)

type Service struct {
    Verifier verifier.Verifier
	KYCRepo repository.KYCRepository
}

type KYCDocumentSubmittedEvent struct {
    KYCID       uint   `json:"kyc_id"`
    DocumentType string `json:"document_type"`
	DocumentURL  string `json:"document_url"`
}

func NewService(verifier verifier.Verifier, kycRepo repository.KYCRepository) *Service {
    return &Service{
		Verifier: verifier,
		KYCRepo: kycRepo,
	}
}

func (s *Service) VerifyKYC(kycDocument *KYCDocumentSubmittedEvent) error {
	log.Printf("[KYCService] Verifying KYC for document ID: %d", kycDocument.KYCID)
	ctx := context.Background()
	result, err := s.Verifier.Verify(ctx, verifier.VerificationInput{
		DocumentType: kycDocument.DocumentType,
		DocumentURL:  kycDocument.DocumentURL,
	})
	if err != nil {
		log.Printf("[KYCService] Error verifying KYC for document ID: %d, error: %v", kycDocument.KYCID, err)
		return err
	}
	log.Printf("[KYCService] Successfully verified KYC for document ID: %d, result: %v", kycDocument.KYCID, result)
	// Update the KYC document status
	return nil

}