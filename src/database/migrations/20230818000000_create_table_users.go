package migrations

import (
	"gorm.io/gorm"
	"yaza/src/app/auth/entities"
	"yaza/src/database"
)

/// please replace or change &entities.Name{}
/// AutoMigrate will create tables, missing foreign keys, constraints, columns and indexes.
// It will change existing column’s type if its size, precision, nullable changed.
// It WON’T delete unused columns to protect your data.

func init() {
	database.MigrationRegister(func(db *gorm.DB) error {
		return db.AutoMigrate(&entities.User{})
	}, func(db *gorm.DB) error {
		return db.Migrator().DropTable(&entities.User{})
	})
}
