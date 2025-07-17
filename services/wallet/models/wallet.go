// internal/model/wallet.go
package model

import "time"

type Wallet struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    string    `gorm:"uniqueIndex"`
	Balance   int64    
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Transaction struct {
	ID            uint      `gorm:"primaryKey"`
	FromUserID    string
	ToUserID      string
	Amount        int64     
	CreatedAt     time.Time
	TransactionID string    `gorm:"uniqueIndex"`
}
