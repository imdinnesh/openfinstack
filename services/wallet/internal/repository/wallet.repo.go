package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	clients "github.com/imdinnesh/openfinstack/services/wallet/client"
	"github.com/imdinnesh/openfinstack/services/wallet/internal/events"
	"github.com/imdinnesh/openfinstack/services/wallet/models"
)

type WalletRepository interface {
	CreateWallet(userID uint) error
	GetWalletByUserID(userID uint) (*models.Wallet, error)
	AddFunds(ctx context.Context, userID uint, amount int64, desc string) error
	WithdrawFunds(ctx context.Context, userID uint, amount int64, desc string) error
	Transfer(ctx context.Context, fromUserID, toUserID uint, amount int64) error
	GetTransactions(userID uint, limit, offset int) ([]models.Transaction, error)
}

type walletRepo struct {
	db        *gorm.DB
	publisher *events.WalletEventPublisher
	ledgerClient *clients.LedgerClient // New field
}

func New(db *gorm.DB, publisher *events.WalletEventPublisher, ledgerClient *clients.LedgerClient) WalletRepository {
	return &walletRepo{db: db, publisher: publisher, ledgerClient: ledgerClient}
}

func (r *walletRepo) CreateWallet(userID uint) error {
	return r.db.Create(&models.Wallet{
		UserID:  userID,
		Balance: 0,
	}).Error
}

func (r *walletRepo) GetWalletByUserID(userID uint) (*models.Wallet, error) {
	var wallet models.Wallet
	err := r.db.Where("user_id = ?", userID).First(&wallet).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &wallet, err
}

func (r *walletRepo) AddFunds(ctx context.Context, userID uint, amount int64, desc string) error {
	var wallet models.Wallet
	
	// Create a balanced transaction in the Ledger service first
	ledgerReq := map[string]interface{}{
		"externalRef": fmt.Sprintf("wallet_credit_%d_%d", userID, amount),
		"type":        "WALLET_TOPUP",
		"description": desc,
		"currency":    "USD",
		"entries": []map[string]interface{}{
			{
				"accountId": fmt.Sprintf("user_wallet_%d", userID),
				"entryType": "CREDIT",
				"amount":    amount,
			},
			{
				"accountId": "internal_bank_account",
				"entryType": "DEBIT",
				"amount":    amount,
			},
		},
	}
	ledgerTxnID, err := r.ledgerClient.CreateTransaction(ctx, ledgerReq)
	if err != nil {
		return fmt.Errorf("failed to create ledger transaction: %w", err)
	}

	err = r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
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

		// Create transaction record and link it to the Ledger transaction
		transaction := models.Transaction{
			WalletID:    wallet.ID,
			LedgerTxnID: ledgerTxnID, // Store the new Ledger ID
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
			fmt.Printf("Failed to publish add funds event: %v\n", publishErr)
		}
	}

	return err
}

func (r *walletRepo) WithdrawFunds(ctx context.Context, userID uint, amount int64, desc string) error {
	var wallet models.Wallet
	
	// Check balance before attempting to create the ledger transaction
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&wallet).Error; err != nil {
		return fmt.Errorf("failed to find wallet: %w", err)
	}
	if wallet.Balance < amount {
		return errors.New("insufficient balance")
	}

	// Create a balanced transaction in the Ledger service first
	ledgerReq := map[string]interface{}{
		"externalRef": fmt.Sprintf("wallet_debit_%d_%d", userID, amount),
		"type":        "WITHDRAWAL",
		"description": desc,
		"currency":    "USD",
		"entries": []map[string]interface{}{
			{
				"accountId": fmt.Sprintf("user_wallet_%d", userID),
				"entryType": "DEBIT",
				"amount":    amount,
			},
			{
				"accountId": "external_bank_account",
				"entryType": "CREDIT",
				"amount":    amount,
			},
		},
	}
	ledgerTxnID, err := r.ledgerClient.CreateTransaction(ctx, ledgerReq)
	if err != nil {
		return fmt.Errorf("failed to create ledger transaction: %w", err)
	}
	
	err = r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Re-fetch with lock for balance update
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("user_id = ?", userID).
			First(&wallet).Error; err != nil {
			return fmt.Errorf("failed to find wallet: %w", err)
		}

		// Check sufficient balance again to be safe
		if wallet.Balance < amount {
			return fmt.Errorf("insufficient balance: required %d, available %d", amount, wallet.Balance)
		}

		// Update balance
		wallet.Balance -= amount
		if err := tx.Save(&wallet).Error; err != nil {
			return fmt.Errorf("failed to update wallet balance: %w", err)
		}

		// Create transaction record and link it to the Ledger transaction
		transaction := models.Transaction{
			WalletID:    wallet.ID,
			LedgerTxnID: ledgerTxnID, // Store the new Ledger ID
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
			fmt.Printf("Failed to publish debit funds event: %v\n", publishErr)
		}
	}

	return err
}

func (r *walletRepo) Transfer(ctx context.Context, fromUserID, toUserID uint, amount int64) error {
	var fromWallet, toWallet models.Wallet
	
	if fromUserID == toUserID {
		return errors.New("cannot transfer to self")
	}
	if amount <= 0 {
		return errors.New("amount must be positive")
	}

	// Lock both wallets in consistent order
	firstID, secondID := fromUserID, toUserID
	if fromUserID > toUserID {
		firstID, secondID = toUserID, fromUserID
	}
	
	// Create the Ledger transaction first, before locking wallets
	ledgerReq := map[string]interface{}{
		"externalRef": fmt.Sprintf("wallet_transfer_%d_%d_%d", fromUserID, toUserID, amount),
		"type":        "WALLET_TRANSFER",
		"description": "Funds transfer",
		"currency":    "USD",
		"entries": []map[string]interface{}{
			{
				"accountId": fmt.Sprintf("user_wallet_%d", fromUserID),
				"entryType": "DEBIT",
				"amount":    amount,
			},
			{
				"accountId": fmt.Sprintf("user_wallet_%d", toUserID),
				"entryType": "CREDIT",
				"amount":    amount,
			},
		},
	}
	ledgerTxnID, err := r.ledgerClient.CreateTransaction(ctx, ledgerReq)
	if err != nil {
		return fmt.Errorf("failed to create ledger transaction: %w", err)
	}

	err = r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Lock first wallet
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("user_id = ?", firstID).First(&fromWallet).Error; err != nil {
			return fmt.Errorf("failed to find wallet for user %d: %w", firstID, err)
		}

		// Lock second wallet
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("user_id = ?", secondID).First(&toWallet).Error; err != nil {
			return fmt.Errorf("failed to find wallet for user %d: %w", secondID, err)
		}

		// Assign wallets based on original order
		if fromUserID != firstID {
			fromWallet, toWallet = toWallet, fromWallet
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
		transactions := []models.Transaction{
			{
				WalletID:    fromWallet.ID,
				LedgerTxnID: ledgerTxnID, // Link to the Ledger transaction
				Type:        models.TRANSFER_OUT,
				Amount:      amount,
				Description: fmt.Sprintf("Transfer to user %d", toUserID),
				ReferenceID: &refID,
				CreatedAt:   time.Now(),
			},
			{
				WalletID:    toWallet.ID,
				LedgerTxnID: ledgerTxnID, // Link to the Ledger transaction
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
			fmt.Printf("Failed to publish transfer event: %v\n", publishErr)
		}
	}

	return err
}

func (r *walletRepo) GetTransactions(userID uint, limit, offset int) ([]models.Transaction, error) {
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