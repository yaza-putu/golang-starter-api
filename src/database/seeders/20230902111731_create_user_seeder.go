package seeders

import (
	"github.com/yaza-putu/golang-starter-api/src/app/auth/entity"
	"github.com/yaza-putu/golang-starter-api/src/database"
	"github.com/yaza-putu/golang-starter-api/src/utils"
	"gorm.io/gorm"
)

// / please replace &entities.Name{} and insert data
func init() {
	database.SeederRegister(func(db *gorm.DB) error {
		m := entity.Users{
			entity.User{
				ID:       utils.Uid(13),
				Name:     "User",
				Email:    "user@mail.com",
				Password: utils.Bcrypt("Password1"),
			},
		}

		return db.Create(&m).Error
	})
}
