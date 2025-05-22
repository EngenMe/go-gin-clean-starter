package migrations

import (
	"gorm.io/gorm"

	"github.com/Caknoooo/go-gin-clean-starter/migrations/seeds"
)

// Seeder seeds the database by executing predefined seeder functions, such as ListUserSeeder, and returns any errors encountered.
func Seeder(db *gorm.DB) error {
	if err := seeds.ListUserSeeder(db); err != nil {
		return err
	}

	return nil
}
