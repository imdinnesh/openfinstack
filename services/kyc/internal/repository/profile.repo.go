// File: repository/profile_repo.go

package repository

import (
	"github.com/imdinnesh/openfinstack/services/kyc/models"
	"gorm.io/gorm"
)

type ProfileRepository interface {
	CreateProfile(profile *models.Profile) error
	GetProfileByUserID(userID uint) (*models.Profile, error)
	UpdateProfile(profile *models.Profile) error
	UpdateKYCStatus(userID uint, status string) error
}

type profileRepository struct {
	db *gorm.DB
}

func NewProfileRepository(db *gorm.DB) ProfileRepository {
	return &profileRepository{db}
}

func (r *profileRepository) CreateProfile(profile *models.Profile) error {
	return r.db.Create(profile).Error
}

func (r *profileRepository) GetProfileByUserID(userID uint) (*models.Profile, error) {
	var profile models.Profile
	if err := r.db.Where("user_id = ?", userID).First(&profile).Error; err != nil {
		return nil, err
	}
	return &profile, nil
}

func (r *profileRepository) UpdateProfile(profile *models.Profile) error {
	return r.db.Save(profile).Error
}

func (r *profileRepository) UpdateKYCStatus(userID uint, status string) error {
	return r.db.Model(&models.Profile{}).
		Where("user_id = ?", userID).
		Update("kyc_status", status).Error
}
