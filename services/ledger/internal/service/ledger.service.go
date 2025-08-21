package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/imdinnesh/openfinstack/services/ledger/db"
	"github.com/imdinnesh/openfinstack/services/ledger/internal/repository"
	"github.com/imdinnesh/openfinstack/services/ledger/models"
	"github.com/imdinnesh/openfinstack/services/ledger/pkg/apierror"
)

type CreateEntry struct {
	AccountID string             `json:"accountId"`
	EntryType models.DebitCredit `json:"entryType"`
	Amount    int64              `json:"amount"`
	Meta      map[string]any     `json:"meta"`
}

type CreateTxnRequest struct {
	ExternalRef string               `json:"externalRef"` // idempotency key
	Type        models.LedgerTxnType `json:"type"`
	Description string               `json:"description"`
	Currency    string               `json:"currency"`
	Entries     []CreateEntry        `json:"entries"`
}

type LedgerService interface {
	CreateAndPost(req CreateTxnRequest) (*models.LedgerTransaction, error)
	Get(id uuid.UUID) (*models.LedgerTransaction, error)
	GetByExternalRef(ref string) (*models.LedgerTransaction, error)
	ListEntries(accountID string, limit, offset int) ([]models.LedgerEntry, error)
	Reverse(original uuid.UUID, reason string) (*models.LedgerTransaction, error)
}

type ledgerService struct {
	db   *gorm.DB
	repo repository.LedgerRepository
}

func NewLedgerService(dbConn *gorm.DB, repo repository.LedgerRepository) LedgerService {
	return &ledgerService{db: dbConn, repo: repo}
}

func (s *ledgerService) CreateAndPost(req CreateTxnRequest) (*models.LedgerTransaction, error) {
	if len(req.Entries) < 2 {
		return nil, apierror.BadRequest("at least two entries required")
	}
	if req.Currency == "" {
		return nil, apierror.BadRequest("currency is required")
	}

	var debitSum, creditSum int64
	entries := make([]models.LedgerEntry, 0, len(req.Entries))
	for i, e := range req.Entries {
		if e.Amount <= 0 {
			return nil, apierror.BadRequest(fmt.Sprintf("entry %d: amount must be > 0", i))
		}
		if e.AccountID == "" {
			return nil, apierror.BadRequest(fmt.Sprintf("entry %d: accountId required", i))
		}
		if e.EntryType != models.DEBIT && e.EntryType != models.CREDIT {
			return nil, apierror.BadRequest(fmt.Sprintf("entry %d: entryType must be DEBIT or CREDIT", i))
		}
		metaBytes, _ := json.Marshal(e.Meta)
		entries = append(entries, models.LedgerEntry{
			AccountID: e.AccountID,
			EntryType: e.EntryType,
			Amount:    e.Amount,
			Currency:  req.Currency,
			Meta:      metaBytes,
		})
		if e.EntryType == models.DEBIT {
			debitSum += e.Amount
		} else {
			creditSum += e.Amount
		}
	}
	if debitSum != creditSum {
		return nil, apierror.BadRequest("debits and credits must balance")
	}

	var created *models.LedgerTransaction
	err := db.WithTx(s.db, func(tx *gorm.DB) error {
		// Idempotency via externalRef unique index
		if req.ExternalRef != "" {
			if existing, err := s.repo.GetByExternalRef(tx, req.ExternalRef); err == nil && existing != nil {
				created = existing
				return nil // idempotent replay
			}
		}

		lt := &models.LedgerTransaction{
			ExternalRef: req.ExternalRef,
			Type:        req.Type,
			Status:      models.TxnPending,
			Description: req.Description,
			Currency:    req.Currency,
		}
		if err := s.repo.CreatePending(tx, lt); err != nil {
			return err
		}
		for i := range entries {
			entries[i].TransactionID = lt.ID
		}
		if err := s.repo.AddEntries(tx, entries); err != nil {
			return err
		}
		postedAt := time.Now().UTC()
		if err := s.repo.MarkPosted(tx, lt.ID, [2]int64{debitSum, creditSum}, postedAt); err != nil {
			return err
		}
		// fetch with entries
		var fetchErr error
		created, fetchErr = s.repo.GetByID(tx, lt.ID)
		return fetchErr
	})
	if err != nil {
		return nil, err
	}
	return created, nil
}

func (s *ledgerService) Get(id uuid.UUID) (*models.LedgerTransaction, error) {
	return s.repo.GetByID(s.db, id)
}

func (s *ledgerService) GetByExternalRef(ref string) (*models.LedgerTransaction, error) {
	return s.repo.GetByExternalRef(s.db, ref)
}

func (s *ledgerService) ListEntries(accountID string, limit, offset int) ([]models.LedgerEntry, error) {
	if limit <= 0 || limit > 200 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}
	return s.repo.ListEntriesByAccount(s.db, accountID, limit, offset)
}

func (s *ledgerService) Reverse(original uuid.UUID, reason string) (*models.LedgerTransaction, error) {
	orig, err := s.repo.GetByID(s.db, original)
	if err != nil {
		return nil, err
	}
	if orig.Status != models.TxnPosted {
		return nil, errors.New("only posted txns can be reversed")
	}

	// Build opposite entries
	req := CreateTxnRequest{
		ExternalRef: "REV-" + uuid.New().String(),
		Type:        models.TxnReversal,
		Description: "Reversal of " + original.String() + ": " + reason,
		Currency:    orig.Currency,
	}
	for _, e := range orig.Entries {
		oppType := models.DEBIT
		if e.EntryType == models.DEBIT {
			oppType = models.CREDIT
		}
		var meta map[string]any
		_ = json.Unmarshal(e.Meta, &meta)
		req.Entries = append(req.Entries, CreateEntry{
			AccountID: e.AccountID,
			EntryType: oppType,
			Amount:    e.Amount,
			Meta:      meta,
		})
	}
	return s.CreateAndPost(req)
}


