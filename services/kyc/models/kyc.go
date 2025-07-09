package models

import (
	"time"
)

type KYC struct {
	ID           uint   `gorm:"primaryKey"`
	UserID       uint   `gorm:"not null"`
	DocumentType string `gorm:"not null"` // AADHAAR, PAN, etc.
	DocumentURL  string `gorm:"not null"`
	Status       string `gorm:"default:'pending'"` // pending, approved, rejected
	RejectReason *string
	VerifiedAt   *time.Time
	VerifiedBy   *uint // Admin ID
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
