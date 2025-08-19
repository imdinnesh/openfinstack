package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type DebitCredit string

const (
	DEBIT  DebitCredit = "DEBIT"
	CREDIT DebitCredit = "CREDIT"
)

type LedgerTxnStatus string

const (
	TxnPending LedgerTxnStatus = "PENDING"
	TxnPosted  LedgerTxnStatus = "POSTED"
	TxnFailed  LedgerTxnStatus = "FAILED"
)

type LedgerTxnType string

const (
	TxnWalletTopup    LedgerTxnType = "WALLET_TOPUP"
	TxnWalletTransfer LedgerTxnType = "WALLET_TRANSFER"
	TxnWithdrawal     LedgerTxnType = "WITHDRAWAL"
	TxnAdjustment     LedgerTxnType = "ADJUSTMENT"
	TxnReversal       LedgerTxnType = "REVERSAL"
)

// LedgerTransaction groups multiple entries that must balance.
type LedgerTransaction struct {
	ID          uuid.UUID       `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	ExternalRef string          `gorm:"uniqueIndex;size:128"` // idempotency key from caller
	Type        LedgerTxnType   `gorm:"type:text;not null"`
	Status      LedgerTxnStatus `gorm:"type:text;not null;default:PENDING"`
	Description string          `gorm:"type:text"`
	Currency    string          `gorm:"size:8;not null"`
	TotalDebit  int64           `gorm:"not null;default:0"`
	TotalCredit int64           `gorm:"not null;default:0"`
	PostedAt    *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Entries     []LedgerEntry `gorm:"constraint:OnDelete:RESTRICT;"`
}

// LedgerEntry is immutable once the parent transaction is POSTED.
type LedgerEntry struct {
	ID            uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	TransactionID uuid.UUID      `gorm:"type:uuid;index;not null"`
	AccountID     string         `gorm:"size:128;index;not null"`
	EntryType     DebitCredit    `gorm:"type:text;not null"`
	Amount        int64          `gorm:"not null;check:amount_positive,amount>0"`
	Currency      string         `gorm:"size:8;not null"`
	Meta          datatypes.JSON `gorm:"type:jsonb"`
	CreatedAt     time.Time
}
