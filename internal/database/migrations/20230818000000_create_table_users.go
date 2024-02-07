package migrations

import (
	"github.com/yaza-putu/golang-starter-api/internal/app/auth/entity"
	"github.com/yaza-putu/golang-starter-api/internal/database"
	"gorm.io/gorm"
)

/// please replace or change &entities.Name{}
/// AutoMigrate will create tables, missing foreign keys, constraints, columns and indexes.
// It will change existing column’s type if its size, precision, nullable changed.
// It WON’T delete unused columns to protect your data.

func init() {
	database.MigrationRegister(func(db *gorm.DB) error {
		return db.AutoMigrate(&entity.User{})
	}, func(db *gorm.DB) error {
		return db.Migrator().DropTable(&entity.User{})
	})
}
