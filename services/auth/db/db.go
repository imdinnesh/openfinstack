package db

import (
	"fmt"
	"log"

	"github.com/imdinnesh/openfinstack/services/auth/config"
	"github.com/imdinnesh/openfinstack/services/auth/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB(cfg *config.Config) *gorm.DB {
	db, err := gorm.Open(postgres.Open(cfg.DBUrl), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto migrate user table
	err = db.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	fmt.Println("Database connected & migrated successfully")
	return db
}
