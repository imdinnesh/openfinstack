package model

import (
	"time"
)

type User struct {
	ID           uint   `gorm:"primaryKey"`
	Email        string `gorm:"uniqueIndex;not null"`
	Phone        string `gorm:"uniqueIndex"`
	PasswordHash string `gorm:"not null"`
	IsVerified   bool   `gorm:"default:false"`
	Role         string `gorm:"default:'user'"`
	LastLoginAt  time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    *time.Time `gorm:"index"`
}
