package service

import (
	"context"
	"errors"
	"fmt"

	clients "github.com/imdinnesh/openfinstack/services/kyc/client"
	"github.com/imdinnesh/openfinstack/services/kyc/config"
	"github.com/imdinnesh/openfinstack/services/kyc/internal/events"
	"github.com/imdinnesh/openfinstack/services/kyc/internal/repository"
	"github.com/imdinnesh/openfinstack/services/kyc/internal/verifier"
	"github.com/imdinnesh/openfinstack/services/kyc/models"
)

type KYCService interface {
	SubmitKYC(kyc *models.KYC) error
	GetUserKYC(userID uint) ([]models.KYC, error)
	ListPending() ([]models.KYC, error)
	VerifyKYC(id uint, status string, reason *string, adminID uint) error
	GetKYCStatusByUserID(userID uint) (string, error)
	UpdateKYCStatus(id uint, status string, reason *string, adminID uint) error
}
type kycService struct {
	repo     repository.KYCRepository
	verifier verifier.Verifier
	events   *events.KYCEventPublisher
	client   *clients.Client
}

func NewKYCService(repo repository.KYCRepository, v verifier.Verifier, e *events.KYCEventPublisher, c *clients.Client) *kycService {
	return &kycService{
		repo:     repo,
		verifier: v,
		events:   e,
		client:   c,
	}
}
func (s *kycService) SubmitKYC(kyc *models.KYC) error {
	err := s.repo.Create(kyc)
	if err != nil {
		return err
	}

	VerificationMethod := config.Load().KYCVerifier
	fmt.Println("KYC Verification Method:", VerificationMethod)
	// If Manual Verification is enabled, skip automatic verification
	if VerificationMethod == config.KYCVerifierManual {
		fmt.Println("Manual verification enabled, skipping automatic verification for KYC ID:", kyc.ID)
		return nil
	}

	if err := s.events.PublishKYCDocumentSubmitted(context.Background(), kyc.ID, kyc.DocumentType, kyc.DocumentURL); err != nil {
		return err
	}
	fmt.Println("KYC document submitted event published for KYC ID:", kyc.ID)
	return nil
}

func (s *kycService) GetUserKYC(userID uint) ([]models.KYC, error) {
	return s.repo.GetByUserID(userID)
}

func (s *kycService) ListPending() ([]models.KYC, error) {
	return s.repo.GetPending()
}

func (s *kycService) VerifyKYC(id uint, status string, reason *string, adminID uint) error {
	kycRecord, err := s.repo.GetKYCByID(id)
	if err != nil {
		return err
	}

	if kycRecord == nil {
		return errors.New("KYC record not found")
	}

	err = s.repo.UpdateStatus(id, status, reason, adminID)
	if err != nil {
		return err
	}

	// Get profile details for the user
	userProfile, err := s.client.GetUserProfile(kycRecord.UserID)
	if err != nil {
		return fmt.Errorf("failed to get user profile: %w", err)
	}

	// Safe dereference of reason
	var reasonStr string
	if reason != nil {
		reasonStr = *reason
	}

	// Notify via event
	if err := s.events.PublishKYCStatus(context.Background(), kycRecord.UserID, userProfile.Email, status, reasonStr); err != nil {
		return fmt.Errorf("failed to publish KYC status event: %w", err)
	}
	fmt.Println("KYC status event published for user ID:", kycRecord.UserID)

	return nil
}

func (s *kycService) GetKYCStatusByUserID(userID uint) (string, error) {
	status, err := s.repo.GetKYCStatusByUserID(userID)
	if err != nil {
		return "", err
	}
	return status, nil
}

func (s *kycService) UpdateKYCStatus(id uint, status string, reason *string, adminID uint) error {
	kycRecord, err := s.repo.GetKYCByID(id)
	if err != nil {
		return err
	}

	if kycRecord == nil {
		return errors.New("KYC record not found")
	}

	if err := s.repo.UpdateStatus(id, status, reason, adminID); err != nil {
		return err
	}

	userProfile, err := s.client.GetUserProfile(kycRecord.UserID)
	if err != nil {
		return fmt.Errorf("failed to get user profile: %w", err)
	}
	
	var reasonStr string
	if reason != nil {
		reasonStr = *reason
	}

	// Notify via event
	if err := s.events.PublishKYCStatus(context.Background(), kycRecord.UserID, userProfile.Email, status, reasonStr); err != nil {
		return fmt.Errorf("failed to publish KYC status event: %w", err)
	}
	fmt.Println("KYC status event published for user ID:", kycRecord.UserID)

	return nil
}
