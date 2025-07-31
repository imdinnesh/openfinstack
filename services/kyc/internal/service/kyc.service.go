package service

import (
	"context"
	"errors"

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
}
type kycService struct {
	repo     repository.KYCRepository
	verifier verifier.Verifier
}

func NewKYCService(repo repository.KYCRepository, v verifier.Verifier) *kycService {
	return &kycService{
		repo:     repo,
		verifier: v,
	}
}
func (s *kycService) SubmitKYC(kyc *models.KYC) error {
	err := s.repo.Create(kyc)
	if err != nil {
		return err
	}
	ctx := context.Background()
	result, err := s.verifier.Verify(ctx, verifier.VerificationInput{
		DocumentType: kyc.DocumentType,
		DocumentURL:  kyc.DocumentURL,
	})
	if err != nil {
		return err  
	}

  if !result.Verified{
    s.repo.UpdateStatus(kyc.ID, "rejected", &result.RejectReason, 0)
    return errors.New("KYC verification failed: " + result.RejectReason)
  }

  // Update KYC status to verified
  err = s.repo.UpdateStatus(kyc.ID, "approved", nil, 0)
  if err != nil {
    return err
  }
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

	return s.repo.UpdateStatus(id, status, reason, adminID)
}

func (s *kycService) GetKYCStatusByUserID(userID uint) (string, error) {
	status, err := s.repo.GetKYCStatusByUserID(userID)
	if err != nil {
		return "", err
	}
	return status, nil
}
