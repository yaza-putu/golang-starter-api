package core

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/yaza-putu/golang-starter-api/src/config"
	"github.com/yaza-putu/golang-starter-api/src/database"
	"github.com/yaza-putu/golang-starter-api/src/logger"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func mysqlDriver() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", config.DB().User, config.DB().Password, config.DB().Host, config.DB().Port, config.DB().Name)

	sqlDB, err := sql.Open("mysql", dsn)
	db, err := gorm.Open(mysql.New(mysql.Config{
		Conn: sqlDB,
	}), &gorm.Config{})

	if err != nil {
		if config.App().Debug == true {
			logger.New(err, logger.SetType(logger.FATAL))
		} else {
			logger.New(
				errors.New("Database connection error, please enable debug mode to view error"),
				logger.SetType(logger.FATAL),
			)
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
		logger.New(
			errors.New("Database Driver Not Found"),
			logger.SetType(logger.FATAL),
		)
	}
}
