package seeders

import (
	"gorm.io/gorm"
	"yaza/src/app/auth/entities"
	"yaza/src/database"
)

// / please replace &entities.Name{} and insert data
func init() {
	database.SeederRegister(func(db *gorm.DB) error {
		m := entities.Roles{
			entities.Role{ID: entities.AdminId, Name: "Admin"},
			entities.Role{ID: entities.UserId, Name: "User"},
		}

		return db.Create(&m).Error
	})
}
