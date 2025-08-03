package db

import (
	"fmt"
	"log"

	"github.com/imdinnesh/openfinstack/services/kyc/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB(cfg *config.Config) *gorm.DB {
	db, err := gorm.Open(postgres.Open(cfg.DBUrl), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto migrate kyc table
	// err = db.AutoMigrate(&models.KYC{})
	// if err != nil {
	// 	log.Fatal("Failed to run migrations:", err)
	// }

	fmt.Println("Database connected & migrated successfully")
	return db
}
