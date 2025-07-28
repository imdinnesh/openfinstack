package db

import (
	"fmt"
	"log"

	"github.com/imdinnesh/openfinstack/gateway/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB(cfg *config.ConfigVariables) *gorm.DB {
	db, err := gorm.Open(postgres.Open(cfg.DBUrl), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}


	fmt.Println("Database connected")
	return db
}
