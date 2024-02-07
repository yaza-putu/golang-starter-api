package seeders

import (
	"github.com/yaza-putu/golang-starter-api/internal/app/auth/entity"
	"github.com/yaza-putu/golang-starter-api/internal/database"
	"github.com/yaza-putu/golang-starter-api/internal/pkg/encrypt"
	"github.com/yaza-putu/golang-starter-api/pkg/unique"
	"gorm.io/gorm"
)

// / please replace &entities.Name{} and insert data
func init() {
	key := unique.Uid(13)
	database.SeederRegister(func(db *gorm.DB) error {
		m := entity.Users{
			entity.User{
				ID:       key,
				Name:     "admin",
				Email:    "admin@mail.com",
				Password: encrypt.Bcrypt("Password1"),
				Avatar:   "assets/images/avatar/default.png",
			},
		}

		return db.Create(&m).Error
	})

	database.SeederRegister(func(db *gorm.DB) error {
		m := entity.RoleUser{
			ID:     unique.Key(13),
			UserId: key,
			RoleId: entity.ADM,
		}

		return db.Create(&m).Error
	})
}
