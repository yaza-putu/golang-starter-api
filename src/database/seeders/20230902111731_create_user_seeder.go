package seeders

import (
	"gorm.io/gorm"
	"yaza/src/app/auth/model"
	"yaza/src/database"
	"yaza/src/utils"
)

// / please replace &ModelName{} and insert data
func init() {
	database.SeederRegister(func(db *gorm.DB) error {
		m := model.Users{
			model.User{ID: utils.Uid(13), Name: "Admin", Email: "admin@mail.com", Password: utils.Bcrypt("Password1")},
			model.User{ID: utils.Uid(13), Name: "User", Email: "auser@mail.com", Password: utils.Bcrypt("Password1")},
		}

		return db.Create(&m).Error
	})
}
