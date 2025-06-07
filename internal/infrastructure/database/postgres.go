package database

import (
	"fmt"

	"github.com/abhi9s-realm/lean-backend-boilerplate-golang/config"
	"github.com/abhi9s-realm/lean-backend-boilerplate-golang/internal/domain/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPostgresDB(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		cfg.DBHost,
		cfg.DBUser,
		cfg.DBPass,
		cfg.DBName,
		cfg.DBPort,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Auto migrate the schemas
	if err := db.AutoMigrate(&models.User{}); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	return db, nil
}
