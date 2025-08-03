package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func AddMoreKYCFields20250803() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "20250803_add_kyc_fields",
		Migrate: func(tx *gorm.DB) error {
			return tx.Transaction(func(tx *gorm.DB) error {
				// Add initial KYC fields to kycs table
				if err := tx.Exec(`
					ALTER TABLE kycs 
						ADD COLUMN full_name TEXT,
						ADD COLUMN date_of_birth TEXT,
						ADD COLUMN gender TEXT,
						ADD COLUMN address_line1 TEXT,
						ADD COLUMN address_line2 TEXT,
						ADD COLUMN district TEXT,
						ADD COLUMN state TEXT,
						ADD COLUMN photo_id_type TEXT,
						ADD COLUMN photo_id_number TEXT
				`).Error; err != nil {
					return err
				}
				return nil
			})
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Transaction(func(tx *gorm.DB) error {
				// Drop all added columns
				if err := tx.Exec(`ALTER TABLE kycs 
					DROP COLUMN full_name,
					DROP COLUMN date_of_birth,
					DROP COLUMN gender,
					DROP COLUMN address_line1,
					DROP COLUMN address_line2,
					DROP COLUMN district,
					DROP COLUMN state,
					DROP COLUMN photo_id_type,
					DROP COLUMN photo_id_number
				`).Error; err != nil {
					return err
				}
				return nil
			})
		},
	}
}
