package migrations

import (
	"gorm.io/gorm"

	"github.com/Caknoooo/go-gin-clean-starter/entity"
)

// Migrate applies automatic database migrations for the specified models, creating or updating tables as needed.
func Migrate(db *gorm.DB) error {
	if err := db.AutoMigrate(
		&entity.User{},
		&entity.RefreshToken{},
	); err != nil {
		return err
	}

	return nil
}
