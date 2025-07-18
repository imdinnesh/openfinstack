package service

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/imdinnesh/openfinstack/services/wallet/internal/repository"
	"github.com/imdinnesh/openfinstack/services/wallet/models"
)

type Service struct {
	Repo *repository.Repository
}

func New(repo *repository.Repository) *Service {
	return &Service{Repo: repo}
}

func (s *Service) Transfer(from, to string, amount int64) error {
	if from == to {
		return errors.New("cannot transfer to self")
	}
	if amount <= 0 {
		return errors.New("amount must be positive")
	}

	fromWallet, err := s.Repo.GetWallet(from)
	if err != nil {
		return err
	}
	if fromWallet.Balance < amount {
		return errors.New("insufficient balance")
	}

	toWallet, err := s.Repo.GetWallet(to)
	if err != nil {
		return err
	}

	fromWallet.Balance -= amount
	toWallet.Balance += amount

	tx := &models.Transaction{
		FromUserID:    from,
		ToUserID:      to,
		Amount:        amount,
		CreatedAt:     time.Now(),
		TransactionID: uuid.New().String(),
	}

	if err := s.Repo.CreateTransaction(tx); err != nil {
		return err
	}

	if err := s.Repo.UpdateWallet(fromWallet); err != nil {
		return err
	}
	return s.Repo.UpdateWallet(toWallet)
}