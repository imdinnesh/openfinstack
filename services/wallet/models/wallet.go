package models

import (
	"time"

	"github.com/google/uuid"
)

type Wallet struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	UserID    uint      `gorm:"not null"`
	Balance   int64     // Store in paisa/cents
	CreatedAt time.Time
	UpdatedAt time.Time
}

type TransactionType string

const (
	CREDIT       TransactionType = "CREDIT"
	DEBIT        TransactionType = "DEBIT"
	TRANSFER_IN  TransactionType = "TRANSFER_IN"
	TRANSFER_OUT TransactionType = "TRANSFER_OUT"
)

type Transaction struct {
	ID            uuid.UUID       `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	WalletID      uuid.UUID       `gorm:"type:uuid;index"`
	LedgerTxnID   uuid.UUID       `gorm:"type:uuid;index;unique"` // New field
	Type          TransactionType `gorm:"type:varchar(20)"`
	Amount        int64
	Description   string
	ReferenceID   *uuid.UUID `gorm:"type:uuid"` // e.g., transfer id
	CreatedAt     time.Time
}