package db

import (
	"log"

	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/imdinnesh/openfinstack/services/kyc/migrations"
	"gorm.io/gorm"
)


func AllMigrations() []*gormigrate.Migration {
	return []*gormigrate.Migration{
		migrations.AddMoreKYCFields20250803(),
	}
}


func RunMigrations(db *gorm.DB) {
	m := gormigrate.New(db, gormigrate.DefaultOptions, AllMigrations())
	if err := m.Migrate(); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}
	log.Println("Migrations ran successfully")
}
