package service

import (
	"errors"

	"github.com/imdinnesh/openfinstack/services/kyc/internal/repository"
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
  repo repository.KYCRepository
}

func NewKYCService(repo repository.KYCRepository) KYCService {
  return &kycService{repo}
}

func (s *kycService) SubmitKYC(kyc *models.KYC) error {
  return s.repo.Create(kyc)
}

func (s *kycService) GetUserKYC(userID uint) ([]models.KYC, error) {
  return s.repo.GetByUserID(userID)
}

func (s *kycService) ListPending() ([]models.KYC, error) {
  return s.repo.GetPending()
}

func (s *kycService) VerifyKYC(id uint, status string, reason *string, adminID uint) error {
  kycRecord,err:=s.repo.GetKYCByID(id)
  if err != nil {
    return err
  }

  if kycRecord == nil {
    return errors.New("KYC record not found")
  }
  
  return s.repo.UpdateStatus(id,status, reason, adminID)
}

func (s *kycService) GetKYCStatusByUserID(userID uint) (string, error) {
  status, err := s.repo.GetKYCStatusByUserID(userID)
  if err != nil {
    return "", err
  }
  return status, nil
}