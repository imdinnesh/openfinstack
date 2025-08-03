package models

import (
	"time"
)

type KYC struct {
	ID           uint   `gorm:"primaryKey"`
	UserID       uint   `gorm:"not null"`
	FullName        string     `gorm:"not null" json:"full_name"`
	DateOfBirth     string     `gorm:"not null" json:"date_of_birth"`
	Gender          string     `gorm:"not null" json:"gender"`
	AddressLine1    string     `gorm:"not null" json:"address_line1"`
	AddressLine2    string     `json:"address_line2"` // optional
	City            string     `gorm:"not null" json:"city"`
	State           string     `gorm:"not null" json:"state"`
	Pincode         string     `gorm:"not null" json:"pincode"`
	DocumentType string `gorm:"not null"` // AADHAAR, PAN, etc.
	DocumentURL  string `gorm:"not null"`
	Status       string `gorm:"default:'pending'"` // pending, approved, rejected
	RejectReason *string
	VerifiedAt   *time.Time
	VerifiedBy   *uint // Admin ID
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
