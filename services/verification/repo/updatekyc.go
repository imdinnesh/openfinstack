package repository

import (
	"github.com/imdinnesh/openfinstack/services/kyc/models"
	"gorm.io/gorm"
)

type KYCRepository interface {
	UpdateStatus(id uint, status string, rejectReason *string, verifiedBy uint) error
}

type kycRepository struct {
	db *gorm.DB
}

func NewKYCRepository(db *gorm.DB) KYCRepository {
	return &kycRepository{db}
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
