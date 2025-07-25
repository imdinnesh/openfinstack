package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/imdinnesh/openfinstack/services/wallet/internal/events"
	"github.com/imdinnesh/openfinstack/services/wallet/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type WalletRepository interface {
	CreateWallet(userID string) error
	GetWalletByUserID(userID string) (*models.Wallet, error)
	AddFunds(ctx context.Context, userID uint, amount int64, desc string) error
	WithdrawFunds(ctx context.Context, userID uint, amount int64, desc string) error
	Transfer(ctx context.Context, fromUserID, toUserID uint, amount int64) error
	GetTransactions(userID string, limit, offset int) ([]models.Transaction, error)
}

type walletRepo struct {
	db        *gorm.DB
	publisher events.WalletEventPublisher
}

func New(db *gorm.DB, publisher events.WalletEventPublisher) WalletRepository {
	return &walletRepo{db: db, publisher: publisher}
}

func (r *walletRepo) CreateWallet(userID string) error {
	// Parse UUID properly with error handling
	uid, err := uuid.Parse(userID)
	if err != nil {
		return fmt.Errorf("invalid userID format: %w", err)
	}
	
	return r.db.Create(&models.Wallet{
		UserID:  uid,
		Balance: 0,
	}).Error
}

func (r *walletRepo) GetWalletByUserID(userID string) (*models.Wallet, error) {
	var wallet models.Wallet
	err := r.db.Where("user_id = ?", userID).First(&wallet).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &wallet, err
}

func (r *walletRepo) AddFunds(ctx context.Context, userID uint, amount int64, desc string) error {
	var wallet models.Wallet
	var transaction models.Transaction
	
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Use proper SELECT FOR UPDATE locking
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("user_id = ?", userID).
			First(&wallet).Error; err != nil {
			return fmt.Errorf("failed to find wallet: %w", err)
		}

		// Update balance
		wallet.Balance += amount
		if err := tx.Save(&wallet).Error; err != nil {
			return fmt.Errorf("failed to update wallet balance: %w", err)
		}

		// Create transaction record
		transaction = models.Transaction{
			WalletID:    wallet.ID,
			Type:        models.CREDIT,
			Amount:      amount,
			Description: desc,
			CreatedAt:   time.Now(),
		}
		
		if err := tx.Create(&transaction).Error; err != nil {
			return fmt.Errorf("failed to create transaction: %w", err)
		}

		return nil
	})

	// Publish event only after successful transaction commit
	if err == nil {
		if publishErr := r.publisher.PublishAddFunds(ctx, userID, amount); publishErr != nil {
			// Log the error but don't fail the transaction
			// You might want to use a proper logger here
			fmt.Printf("Failed to publish add funds event: %v\n", publishErr)
		}
	}

	return err
}

func (r *walletRepo) WithdrawFunds(ctx context.Context, userID uint, amount int64, desc string) error {
	var wallet models.Wallet
	var transaction models.Transaction
	
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Use proper SELECT FOR UPDATE locking
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("user_id = ?", userID).
			First(&wallet).Error; err != nil {
			return fmt.Errorf("failed to find wallet: %w", err)
		}

		// Check sufficient balance
		if wallet.Balance < amount {
			return fmt.Errorf("insufficient balance: required %d, available %d", amount, wallet.Balance)
		}

		// Update balance
		wallet.Balance -= amount
		if err := tx.Save(&wallet).Error; err != nil {
			return fmt.Errorf("failed to update wallet balance: %w", err)
		}

		// Create transaction record
		transaction = models.Transaction{
			WalletID:    wallet.ID,
			Type:        models.DEBIT,
			Amount:      amount,
			Description: desc,
			CreatedAt:   time.Now(),
		}
		
		if err := tx.Create(&transaction).Error; err != nil {
			return fmt.Errorf("failed to create transaction: %w", err)
		}

		return nil
	})

	// Publish event only after successful transaction commit
	if err == nil {
		if publishErr := r.publisher.PublishDebitFunds(ctx, userID, amount); publishErr != nil {
			// Log the error but don't fail the transaction
			fmt.Printf("Failed to publish debit funds event: %v\n", publishErr)
		}
	}

	return err
}

func (r *walletRepo) Transfer(ctx context.Context, fromUserID, toUserID uint, amount int64) error {
	var fromWallet, toWallet models.Wallet
	var transactions []models.Transaction
	
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Lock both wallets in consistent order to prevent deadlocks
		// Always lock the wallet with smaller ID first
		firstID, secondID := fromUserID, toUserID
		if fromUserID > toUserID {
			firstID, secondID = toUserID, fromUserID
		}

		// Lock first wallet
		var firstWallet models.Wallet
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("user_id = ?", firstID).First(&firstWallet).Error; err != nil {
			return fmt.Errorf("failed to find wallet for user %d: %w", firstID, err)
		}

		// Lock second wallet
		var secondWallet models.Wallet
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("user_id = ?", secondID).First(&secondWallet).Error; err != nil {
			return fmt.Errorf("failed to find wallet for user %d: %w", secondID, err)
		}

		// Assign wallets based on original order
		if fromUserID == firstID {
			fromWallet, toWallet = firstWallet, secondWallet
		} else {
			fromWallet, toWallet = secondWallet, firstWallet
		}

		// Check sufficient balance
		if fromWallet.Balance < amount {
			return fmt.Errorf("insufficient balance: required %d, available %d", amount, fromWallet.Balance)
		}

		// Update balances
		fromWallet.Balance -= amount
		toWallet.Balance += amount

		if err := tx.Save(&fromWallet).Error; err != nil {
			return fmt.Errorf("failed to update sender wallet: %w", err)
		}
		if err := tx.Save(&toWallet).Error; err != nil {
			return fmt.Errorf("failed to update receiver wallet: %w", err)
		}

		// Generate reference ID for linking transactions
		refID := uuid.New()

		// Create transaction records
		transactions = []models.Transaction{
			{
				WalletID:    fromWallet.ID,
				Type:        models.TRANSFER_OUT,
				Amount:      amount,
				Description: fmt.Sprintf("Transfer to user %d", toUserID),
				ReferenceID: &refID,
				CreatedAt:   time.Now(),
			},
			{
				WalletID:    toWallet.ID,
				Type:        models.TRANSFER_IN,
				Amount:      amount,
				Description: fmt.Sprintf("Transfer from user %d", fromUserID),
				ReferenceID: &refID,
				CreatedAt:   time.Now(),
			},
		}

		if err := tx.Create(&transactions).Error; err != nil {
			return fmt.Errorf("failed to create transaction records: %w", err)
		}

		return nil
	})

	// Publish event only after successful transaction commit
	if err == nil {
		if publishErr := r.publisher.PublishTransaction(ctx, fromUserID, toUserID, amount); publishErr != nil {
			// Log the error but don't fail the transaction
			fmt.Printf("Failed to publish transfer event: %v\n", publishErr)
		}
	}

	return err
}

func (r *walletRepo) GetTransactions(userID string, limit, offset int) ([]models.Transaction, error) {
	var wallet models.Wallet
	if err := r.db.Where("user_id = ?", userID).First(&wallet).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []models.Transaction{}, nil
		}
		return nil, fmt.Errorf("failed to find wallet: %w", err)
	}

	var txs []models.Transaction
	err := r.db.Where("wallet_id = ?", wallet.ID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&txs).Error
	
	if err != nil {
		return nil, fmt.Errorf("failed to fetch transactions: %w", err)
	}
	
	return txs, nil
}