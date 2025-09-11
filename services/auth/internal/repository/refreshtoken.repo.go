package repository

import (
	"time"

	"github.com/imdinnesh/openfinstack/services/auth/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type RefreshTokenRepository interface {
	Upsert(token *models.RefreshToken) error
	FindByUserID(userID uint) (*models.RefreshToken, error)
	DeleteByUserID(userID uint) error
	IsExpired(userID uint) (bool, error)
}

type refreshTokenRepository struct {
	db *gorm.DB
}

func NewRefreshTokenRepository(db *gorm.DB) RefreshTokenRepository {
	return &refreshTokenRepository{db: db}
}

// Upsert ensures only one refresh token per user
func (r *refreshTokenRepository) Upsert(token *models.RefreshToken) error {
	return r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}}, // conflict on user_id
		UpdateAll: true,                               // overwrite old record
	}).Create(token).Error
}

// FindByUserID retrieves the refresh token for a user
func (r *refreshTokenRepository) FindByUserID(userID uint) (*models.RefreshToken, error) {
	var rt models.RefreshToken
	err := r.db.Where("user_id = ?", userID).First(&rt).Error
	if err != nil {
		return nil, err
	}
	return &rt, nil
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

	return refreshToken.ExpiresAt.Before(time.Now()), nil
}
