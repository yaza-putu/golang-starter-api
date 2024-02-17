package core

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/yaza-putu/golang-starter-api/internal/config"
	"github.com/yaza-putu/golang-starter-api/internal/database"
	"github.com/yaza-putu/golang-starter-api/internal/pkg/logger"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
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

	sqlDb, err := db.DB()
	if err != nil {
		logger.New(err, logger.SetType(logger.FATAL))
	}

	sqlDb.SetMaxIdleConns(config.DB().Idle)
	sqlDb.SetMaxOpenConns(config.DB().MaxConn)
	sqlDb.SetConnMaxLifetime(time.Hour * time.Duration(config.DB().ConnLifetime))

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
