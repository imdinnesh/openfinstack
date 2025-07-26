package repository

import (
	"github.com/imdinnesh/openfinstack/services/kyc/models"
	"gorm.io/gorm"
)

type KYCRepository interface {
	Create(kyc *models.KYC) error
	GetByUserID(userID uint) ([]models.KYC, error)
	GetPending() ([]models.KYC, error)
	UpdateStatus(id uint, status string, rejectReason *string, verifiedBy uint) error
	GetKYCByID(id uint) (*models.KYC, error)
}

type kycRepository struct {
	db *gorm.DB
}

func NewKYCRepository(db *gorm.DB) KYCRepository {
	return &kycRepository{db}
}

func (r *kycRepository) Create(kyc *models.KYC) error {
	return r.db.Create(kyc).Error
}

func (r *kycRepository) GetByUserID(userID uint) ([]models.KYC, error) {
	var kycs []models.KYC
	err := r.db.Where("user_id = ?", userID).Find(&kycs).Error
	return kycs, err
}

func (r *kycRepository) GetPending() ([]models.KYC, error) {
	var kycs []models.KYC
	err := r.db.Where("status = ?", "pending").Find(&kycs).Error
	return kycs, err
}

func (r *kycRepository) UpdateStatus(id uint, status string, rejectReason *string, verifiedBy uint) error {
	return r.db.Model(&models.KYC{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":        status,
			"reject_reason": rejectReason,
			"verified_at":   gorm.Expr("CURRENT_TIMESTAMP"),
			"verified_by":   verifiedBy,
		}).Error
}

func (r *kycRepository) GetKYCByID(id uint) (*models.KYC, error) {
	var kyc models.KYC
	err := r.db.First(&kyc, id).Error
	if err != nil {
		return nil, err
	}
	return &kyc, nil
}
