package db

import (
	"database/sql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/imdinnesh/openfinstack/services/ledger/models"
)

func Connect(dsn string) (*gorm.DB, error) {
	return gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
}

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(&models.LedgerTransaction{}, &models.LedgerEntry{})
}

var sqlTxOptions = &sql.TxOptions{
	Isolation: sql.LevelReadCommitted,
	ReadOnly:  false,
}

// WithTx runs fn within a DB transaction.
func WithTx(db *gorm.DB, fn func(tx *gorm.DB) error) error {
	return db.Transaction(func(tx *gorm.DB) error {
		tx.Statement.Settings.Store("gorm:insert_option", "ON CONFLICT DO NOTHING")
		return fn(tx)
	}, sqlTxOptions)
}



