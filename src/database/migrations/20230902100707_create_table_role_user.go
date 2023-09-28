package migrations

import (
	"gorm.io/gorm"
	"time"
	"yaza/src/database"
)

// / please replace or change &entities.Name{}
// / AutoMigrate will create tables, missing foreign keys, constraints, columns and indexes.
// It will change existing column’s type if its size, precision, nullable changed.
// It WON’T delete unused columns to protect your data.
type RoleUsers struct {
	UserID    string `gorm:"user_id;primaryKey;type:char(36)" json:"user_id"`
	RoleID    string `gorm:"role_id;primaryKey;type:char(36)" json:"role_id"`
	CreatedAt time.Time
}

func init() {
	database.MigrationRegister(func(db *gorm.DB) error {
		return db.AutoMigrate(&RoleUsers{})
	}, func(db *gorm.DB) error {
		return db.Migrator().DropTable(&RoleUsers{})
	})
}
