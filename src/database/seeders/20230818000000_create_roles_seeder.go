package seeders

import (
	"gorm.io/gorm"
	"yaza/src/app/auth/model"
	"yaza/src/database"
)

// / please replace &ModelName{} and insert data
func init() {
	database.SeederRegister(func(db *gorm.DB) error {
		m := model.Roles{
			model.Role{ID: model.AdminId, Name: "Admin"},
			model.Role{ID: model.UserId, Name: "User"},
		}

		return db.Create(&m).Error
	})
}
