package core

import (
	"database/sql"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"yaza/src/config"
	"yaza/src/database"
)

func mysqlDriver() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", config.DB().User, config.DB().Password, config.DB().Host, config.DB().Port, config.DB().Name)

	sqlDB, err := sql.Open("mysql", dsn)
	db, err := gorm.Open(mysql.New(mysql.Config{
		Conn: sqlDB,
	}), &gorm.Config{})

	if err != nil {
		if config.App().Debug == true {
			log.Panic(err)
		} else {
			log.Panic("Database connection error, please enable debug mode to view error")
		}
	}

	database.Instance = db
}

// Database load instance
func Database() {
	switch config.DB().Driver {
	case "mysql":
		mysqlDriver()
		break

	default:
		log.Panic("Database Driver Not Found")
	}
}
