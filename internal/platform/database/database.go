package database

import (
	"fmt"

	"multicliente-backend/internal/platform/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Connect establishes a connection to PostgreSQL using GORM.
func Connect(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBSSLMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return db, nil
}

// Migrate runs GORM AutoMigrate for the provided models after ensuring the required schema exists.
func Migrate(db *gorm.DB, models ...interface{}) error {
	// Create administrative schema if not exists
	if err := db.Exec("CREATE SCHEMA IF NOT EXISTS administrative").Error; err != nil {
		return fmt.Errorf("failed to create schema administrative: %w", err)
	}
	return db.AutoMigrate(models...)
}
