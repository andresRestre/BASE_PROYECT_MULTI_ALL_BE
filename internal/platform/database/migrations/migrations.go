package migrations

import (
	"fmt"
	"gorm.io/gorm"
)

// Migrate runs GORM AutoMigrate for the provided models after ensuring the required schema exists.
func Migrate(db *gorm.DB, models ...interface{}) error {
	// Create administrative schema if not exists
	if err := db.Exec("CREATE SCHEMA IF NOT EXISTS administrative").Error; err != nil {
		return fmt.Errorf("failed to create schema administrative: %w", err)
	}
	// Create app schema if not exists
	if err := db.Exec("CREATE SCHEMA IF NOT EXISTS app").Error; err != nil {
		return fmt.Errorf("failed to create schema app: %w", err)
	}
	return db.AutoMigrate(models...)
}
