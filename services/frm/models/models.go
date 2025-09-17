// models.go
package models

import (
    "time"
)

type Rule struct {
    ID          uint      `gorm:"primaryKey"`
    Name        string    `gorm:"uniqueIndex;size:128"`
    Description string
    // Type: "velocity", "blacklist", "amount", "device_anomaly", ...
    Type        string    `gorm:"index;size:64"`
    // JSON conditions, example: {"field":"amount","op":">","value":10000}
    Conditions  string    `gorm:"type:jsonb"`
    Threshold   float64   // optional numeric threshold
    Enabled     bool
    CreatedAt   time.Time
    UpdatedAt   time.Time
}

type TransactionEvent struct {
    ID            uint      `gorm:"primaryKey"`
    EventID       string    `gorm:"uniqueIndex;size:128"`
    UserID        uint
    DeviceID      string    `gorm:"size:128;index"`
    IP            string    `gorm:"size:64;index"`
    Amount        float64
    Currency      string    `gorm:"size:8"`
    Type          string    `gorm:"size:64"` // e.g. "payment", "login", "transfer"
    Metadata      string    `gorm:"type:jsonb"`
    CreatedAt     time.Time
}

type Alert struct {
    ID         uint      `gorm:"primaryKey"`
    EventID    string    `gorm:"index;size:128"`
    UserID     uint
    RuleID     *uint
    RiskScore  float64
    Reason     string
    Payload    string    `gorm:"type:jsonb"`
    CreatedAt  time.Time
}
