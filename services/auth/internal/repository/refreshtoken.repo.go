package repository

import (
	"github.com/imdinnesh/openfinstack/services/auth/models"
	"gorm.io/gorm"
)

type RefreshTokenRepository interface {
	Create(token *models.RefreshToken) error
	DeleteByUserID(userID uint) error
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
