package repository

import (
	"github.com/imdinnesh/openfinstack/services/auth/models"
	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(user *models.User) error
	FindByEmail(email string) (*models.User, error)
	UpdatePassword(userID uint, newHash string) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) CreateUser(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) UpdatePassword(userID uint, newHash string) error {
	return r.db.Model(&models.User{}).
		Where("id = ?", userID).
		Update("password_hash", newHash).Error
}