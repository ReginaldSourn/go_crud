package migrations

import (
	"fmt"

	"github.com/reginaldsourn/go-crud/internal/core/domain"
	"gorm.io/gorm"
)

// Run applies all database migrations managed by GORM.
func Run(db *gorm.DB) error {
	if db == nil {
		return fmt.Errorf("db is nil")
	}

	if err := db.AutoMigrate(&domain.User{}); err != nil {
		return fmt.Errorf("auto migrate users: %w", err)
	}

	return nil
}
