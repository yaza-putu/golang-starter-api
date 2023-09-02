package seeders

import (
	"gorm.io/gorm"
	"time"
	"yaza/src/app/auth/model"
	"yaza/src/database"
	"yaza/src/utils"
)

// / please replace &ModelName{} and insert data
func init() {
	database.SeederRegister(func(db *gorm.DB) error {
		m := model.Users{
			model.User{
				ID:       utils.Uid(13),
				Name:     "Admin",
				Email:    "admin@mail.com",
				Password: utils.Bcrypt("Password1"),
				Roles: model.Roles{
					model.Role{ID: model.AdminId, CreatedAt: time.Now()},
				},
			},
			model.User{
				ID:       utils.Uid(13),
				Name:     "User",
				Email:    "user@mail.com",
				Password: utils.Bcrypt("Password1"),
				Roles: model.Roles{
					model.Role{ID: model.UserId, CreatedAt: time.Now()},
				},
			},
		}

		return db.Create(&m).Error
	})
}
