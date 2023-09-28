package seeders

import (
	"gorm.io/gorm"
	"time"
	"yaza/src/app/auth/entities"
	"yaza/src/database"
	"yaza/src/utils"
)

// / please replace &entities.Name{} and insert data
func init() {
	database.SeederRegister(func(db *gorm.DB) error {
		m := entities.Users{
			entities.User{
				ID:       utils.Uid(13),
				Name:     "Admin",
				Email:    "admin@mail.com",
				Password: utils.Bcrypt("Password1"),
				Roles: entities.Roles{
					entities.Role{ID: entities.AdminId, CreatedAt: time.Now()},
				},
			},
			entities.User{
				ID:       utils.Uid(13),
				Name:     "User",
				Email:    "user@mail.com",
				Password: utils.Bcrypt("Password1"),
				Roles: entities.Roles{
					entities.Role{ID: entities.UserId, CreatedAt: time.Now()},
				},
			},
		}

		return db.Create(&m).Error
	})
}
