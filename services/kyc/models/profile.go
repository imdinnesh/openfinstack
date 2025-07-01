package models

import (
	"time"
)

type Profile struct {
	ID           uint      `gorm:"primaryKey"`
	UserID       uint      `gorm:"uniqueIndex"`
	FullName     string    `gorm:"not null"`
	DOB          time.Time
	Address      string
	PAN          string    `gorm:"uniqueIndex"`
	Aadhaar      string    `gorm:"uniqueIndex"`
	KYCStatus    string    `gorm:"default:'pending'"` // pending, verified, rejected
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    *time.Time `gorm:"index"`
}
