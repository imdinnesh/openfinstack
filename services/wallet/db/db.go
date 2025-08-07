package db

import (
	"fmt"
	"log"

	"github.com/imdinnesh/openfinstack/services/wallet/config"
	"github.com/imdinnesh/openfinstack/services/wallet/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB(cfg *config.Config) *gorm.DB {
	db, err := gorm.Open(postgres.Open(cfg.DBUrl), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto migrate wallet table
	db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp"`)
	err = db.AutoMigrate(&models.Wallet{})
	if err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	// Auto migrate transaction table
	db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp"`)
	err = db.AutoMigrate(&models.Transaction{})
	if err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	fmt.Println("Database connected & migrated successfully")
	return db
}
