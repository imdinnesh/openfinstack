package service

import (
	"context"
	"log"

	clients "github.com/imdinnesh/openfinstack/services/verifications/client"
	repository "github.com/imdinnesh/openfinstack/services/verifications/repo"
	"github.com/imdinnesh/openfinstack/services/verifications/verifier"
)

type Service struct {
    Verifier verifier.Verifier
	KYCRepo repository.KYCRepository
	Client *clients.Client
}

type KYCDocumentSubmittedEvent struct {
    KYCID       uint   `json:"kyc_id"`
    DocumentType string `json:"document_type"`
	DocumentURL  string `json:"document_url"`
}

func NewService(verifier verifier.Verifier, kycRepo repository.KYCRepository, kycClient *clients.Client) *Service {
    return &Service{
		Verifier: verifier,
		KYCRepo: kycRepo,
		Client: kycClient,
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
	// Update the KYC document status
	if result.Verified{
		log.Printf("[KYCService] KYC document ID: %d verified successfully", kycDocument.KYCID)
		_,err:=s.Client.UpdateKYCStatus(kycDocument.KYCID,"approved", nil,0)
		if err != nil {
			log.Printf("[KYCService] Error updating KYC status for document ID: %d, error: %v", kycDocument.KYCID, err)
			return err
		}
		log.Printf("[KYCService] KYC status updated successfully for document ID: %d", kycDocument.KYCID)
		return nil
		
	} else {
		log.Printf("[KYCService] KYC document ID: %d verification failed", kycDocument.KYCID)
		reason := "Verification failed"
		_, err := s.Client.UpdateKYCStatus(kycDocument.KYCID, "rejected", &reason, 0)
		if err != nil {
			log.Printf("[KYCService] Error updating KYC status for document ID: %d, error: %v", kycDocument.KYCID, err)
			return err
		}
		log.Printf("[KYCService] KYC status updated successfully for document ID: %d", kycDocument.KYCID)
		return nil
	}
}