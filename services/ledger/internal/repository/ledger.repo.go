package repository

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/imdinnesh/openfinstack/services/ledger/models"
	"gorm.io/gorm"
)

type LedgerRepository interface {
	CreatePending(tx *gorm.DB, t *models.LedgerTransaction) error
	AddEntries(tx *gorm.DB, entries []models.LedgerEntry) error
	MarkPosted(tx *gorm.DB, tID uuid.UUID, totals [2]int64, postedAt time.Time) error
	GetByID(db *gorm.DB, id uuid.UUID) (*models.LedgerTransaction, error)
	GetByExternalRef(db *gorm.DB, ref string) (*models.LedgerTransaction, error)
	ListEntriesByAccount(db *gorm.DB, accountID string, limit, offset int) ([]models.LedgerEntry, error)
}

type ledgerRepository struct{ db *gorm.DB }

func NewLedgerRepository(db *gorm.DB) LedgerRepository { return &ledgerRepository{db: db} }

func (r *ledgerRepository) CreatePending(tx *gorm.DB, t *models.LedgerTransaction) error {
	return tx.Create(t).Error
}

func (r *ledgerRepository) AddEntries(tx *gorm.DB, entries []models.LedgerEntry) error {
	return tx.Create(&entries).Error
}

func (r *ledgerRepository) MarkPosted(tx *gorm.DB, tID uuid.UUID, totals [2]int64, postedAt time.Time) error {
	res := tx.Model(&models.LedgerTransaction{}).
		Where("id = ? AND status = ?", tID, models.TxnPending).
		Updates(map[string]any{
			"status":       models.TxnPosted,
			"total_debit":  totals[0],
			"total_credit": totals[1],
			"posted_at":    postedAt,
		})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return errors.New("transaction not found or not pending")
	}
	return nil
}

func (r *ledgerRepository) GetByID(db *gorm.DB, id uuid.UUID) (*models.LedgerTransaction, error) {
	var t models.LedgerTransaction
	if err := db.Preload("Entries").First(&t, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *ledgerRepository) GetByExternalRef(db *gorm.DB, ref string) (*models.LedgerTransaction, error) {
	var t models.LedgerTransaction
	err := db.Preload("Entries").First(&t, "external_ref = ?", ref).Error
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *ledgerRepository) ListEntriesByAccount(db *gorm.DB, accountID string, limit, offset int) ([]models.LedgerEntry, error) {
	var entries []models.LedgerEntry
	err := db.Where("account_id = ?", accountID).Order("created_at DESC").Limit(limit).Offset(offset).Find(&entries).Error
	return entries, err
}
