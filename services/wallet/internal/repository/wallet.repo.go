package repository

import (
	"github.com/imdinnesh/openfinstack/services/wallet/models"
	"gorm.io/gorm"
)

type Repository struct {
	DB *gorm.DB
}

func New(db *gorm.DB) *Repository {
	return &Repository{DB: db}
}

func (r *Repository) GetWallet(userID string) (*models.Wallet, error) {
	var w models.Wallet
	err := r.DB.FirstOrCreate(&w, models.Wallet{UserID: userID}).Error
	return &w, err
}

func (r *Repository) UpdateWallet(wallet *models.Wallet) error {
	return r.DB.Save(wallet).Error
}

func (r *Repository) CreateTransaction(tx *models.Transaction) error {
	return r.DB.Create(tx).Error
}

func (r *Repository) GetTransactions(userID string) ([]models.Transaction, error) {
	var txns []models.Transaction
	err := r.DB.Where("from_user_id = ? OR to_user_id = ?", userID, userID).Order("created_at desc").Find(&txns).Error
	return txns, err
}
