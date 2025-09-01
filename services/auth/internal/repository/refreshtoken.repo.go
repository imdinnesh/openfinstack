package repository

import (
	"time"

	"github.com/imdinnesh/openfinstack/services/auth/models"
	"gorm.io/gorm"
)

type RefreshTokenRepository interface {
	Create(token *models.RefreshToken) error
	DeleteByUserID(userID uint) error
	IsExpired(userID uint) (bool, error)
}

type refreshTokenRepository struct {
	db *gorm.DB
}

func NewRefreshTokenRepository(db *gorm.DB) RefreshTokenRepository {
	return &refreshTokenRepository{db: db}
}

func (r *refreshTokenRepository) Create(token *models.RefreshToken) error {
	return r.db.Create(token).Error
}

func (r *refreshTokenRepository) DeleteByUserID(userID uint) error {
	return r.db.Where("user_id = ?", userID).Delete(&models.RefreshToken{}).Error
}

func (r *refreshTokenRepository) IsExpired(userID uint) (bool, error) {
	var refreshToken models.RefreshToken

	err := r.db.Where("user_id = ?", userID).First(&refreshToken).Error
	if err != nil {
		return false, err
	}

	// true if expired, false if still valid
	return refreshToken.ExpiresAt.Before(time.Now()), nil
}
