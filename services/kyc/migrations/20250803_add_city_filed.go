package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func AddCityFields20250803() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "20250803_add_city_fields",
		Migrate: func(tx *gorm.DB) error {
			return tx.Transaction(func(tx *gorm.DB) error {
				// Step 1: Add new columns as nullable
				err := tx.Exec(`
					ALTER TABLE kycs
					ADD COLUMN city TEXT,
					ADD COLUMN pincode TEXT,
					DROP COLUMN district,
					DROP COLUMN photo_id_type,
					DROP COLUMN photo_id_number
				`).Error
				if err != nil {
					return err
				}

				// Step 2: Optionally update old rows with dummy or empty data if needed
				err = tx.Exec(`
					UPDATE kycs
					SET
						city = 'Unknown',
						pincode = '000000'
					WHERE full_name IS NULL;
				`).Error
				if err != nil {
					return err
				}

				// Step 3: Set columns as NOT NULL
				return tx.Exec(`
					ALTER TABLE kycs
					ALTER COLUMN city SET NOT NULL,
					ALTER COLUMN pincode SET NOT NULL;
				`).Error
			})
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Transaction(func(tx *gorm.DB) error {
				// Drop all added columns
				if err := tx.Exec(`
					ALTER TABLE kycs
					DROP COLUMN full_name,
					DROP COLUMN date_of_birth,
					DROP COLUMN gender,
					DROP COLUMN address_line1,
					DROP COLUMN address_line2,
					DROP COLUMN city,
					DROP COLUMN state,
					DROP COLUMN pincode;
				`).Error; err != nil {
					return err
				}
				return nil
			})
		},
	}
}
