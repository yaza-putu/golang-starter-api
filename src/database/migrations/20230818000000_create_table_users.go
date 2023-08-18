package migrations

import (
	"gorm.io/gorm"
	"yaza/src/app/auth/model"
	"yaza/src/database"
)

/// please replace or change &ModelName{}
/// AutoMigrate will create tables, missing foreign keys, constraints, columns and indexes.
// It will change existing column’s type if its size, precision, nullable changed.
// It WON’T delete unused columns to protect your data.

func init() {
	database.MigrationRegister(func(db *gorm.DB) error {
		return db.AutoMigrate(&model.User{})
	}, func(db *gorm.DB) error {
		return db.Migrator().DropTable(&model.User{})
	})
}
