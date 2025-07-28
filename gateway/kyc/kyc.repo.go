package kyc

import (
	"errors"

	"github.com/imdinnesh/openfinstack/services/kyc/models"
	"gorm.io/gorm"
)

type Repository interface {
    // returns “active”, “pending”, “rejected”, etc.
    GetStatus(userID uint) (string, error)
}

type repo struct {
    db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
    return &repo{db: db}
}

func (r *repo) GetStatus(userID uint) (string, error) {
	var kyc models.KYC
	err := r.db.Where("user_id = ?", userID).First(&kyc).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", nil // no KYC record found
		}
		return "", err // other error
	}

	return kyc.Status, nil // return the status of the KYC record
}
