// Package database opens the GORM connection and keeps the schema in sync with the models.
package database

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"tripo-backend/internal/config"
	"tripo-backend/internal/models"
)

func Connect(cfg config.DBConfig) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(cfg.DSN()), &gorm.Config{TranslateError: true})
	if err != nil {
		return nil, fmt.Errorf("connect database: %w", err)
	}
	return db, nil
}

// Migrate creates/updates tables from the model definitions.
// Order does not matter; GORM resolves foreign-key dependencies.
func Migrate(db *gorm.DB) error {
	err := db.AutoMigrate(
		&models.User{},
		&models.RefreshToken{},
		&models.Itinerary{},
		&models.ItineraryDay{},
		&models.Place{},
		&models.ItineraryStop{},
		&models.RouteLeg{},
		&models.Memory{},
		&models.Activity{},
		&models.ItineraryShare{},
		&models.ChatConversation{},
		&models.ChatMessage{},
	)
	if err != nil {
		return fmt.Errorf("auto migrate: %w", err)
	}
	return nil
}
