package seeders

import (
	"github.com/yaza-putu/golang-starter-api/internal/app/auth/entity"
	"github.com/yaza-putu/golang-starter-api/internal/database"
	"gorm.io/gorm"
)

// / please replace &entities.Name{} and insert data
func init() {
	database.SeederRegister(func(db *gorm.DB) error {
		m := entity.Roles{
			entity.Role{
				ID:   entity.ADM,
				Name: "adm",
			},
			entity.Role{
				ID:   entity.USR,
				Name: "usr",
			},
		}

		return db.Create(&m).Error
	})
}
