package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/imdinnesh/openfinstack/services/wallet/internal/repository"
	"github.com/imdinnesh/openfinstack/services/wallet/models"
)

type WalletService interface {
	CreateWallet(userID uint) error
	GetWallet(userID uint) (*models.Wallet, error)
	AddFunds(ctx context.Context, userID uint, amount int64, desc string) error
	WithdrawFunds(ctx context.Context, userID uint, amount int64, desc string) error
	Transfer(ctx context.Context, fromUserID, toUserID uint, amount int64) error
	GetTransactions(userID uint, limit, offset int) ([]models.Transaction, error)
}

type walletService struct {
	repo repository.WalletRepository
}

func New(repo repository.WalletRepository) WalletService {
	return &walletService{repo: repo}
}

func (s *walletService) CreateWallet(userID uint) error {
	// In future, check if wallet already exists
	return s.repo.CreateWallet(userID)
}

func (s *walletService) GetWallet(userID uint) (*models.Wallet, error) {
	return s.repo.GetWalletByUserID(userID)
}

func (s *walletService) AddFunds(ctx context.Context, userID uint, amount int64, desc string) error {
	if amount <= 0 {
		return errors.New("amount must be positive")
	}

	wallet, err := s.repo.GetWalletByUserID(userID)
	if err != nil {
		return err
	}
	if wallet == nil {
		return fmt.Errorf("wallet not found for user %d", userID)
	}

	// Optional: Check wallet is active
	// if wallet.Status == "frozen" { return errors.New("wallet is frozen") }

	return s.repo.AddFunds(ctx, userID, amount, desc)
}

func (s *walletService) WithdrawFunds(ctx context.Context, userID uint, amount int64, desc string) error {
	if amount <= 0 {
		return errors.New("amount must be positive")
	}

	wallet, err := s.repo.GetWalletByUserID(userID)
	if err != nil {
		return err
	}
	if wallet == nil {
		return fmt.Errorf("wallet not found for user %d", userID)
	}

	// Optional: Check wallet is active
	return s.repo.WithdrawFunds(ctx, userID, amount, desc)
}

func (s *walletService) Transfer(ctx context.Context, fromUserID, toUserID uint, amount int64) error {
	if fromUserID == toUserID {
		return errors.New("cannot transfer to self")
	}
	if amount <= 0 {
		return errors.New("amount must be positive")
	}

	// Optional: Add KYC check or daily limit validation here

	return s.repo.Transfer(ctx, fromUserID, toUserID, amount)
}

func (s *walletService) GetTransactions(userID uint, limit, offset int) ([]models.Transaction, error) {
	if limit <= 0 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}
	return s.repo.GetTransactions(userID, limit, offset)
}