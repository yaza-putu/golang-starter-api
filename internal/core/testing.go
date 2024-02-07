package core

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"github.com/yaza-putu/golang-starter-api/internal/config"
	"github.com/yaza-putu/golang-starter-api/internal/database"
	"github.com/yaza-putu/golang-starter-api/internal/pkg/logger"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"path/filepath"
	"runtime"
)

func EnvTesting() error {
	_, b, _, _ := runtime.Caller(0)

	// Root folder of this project
	Root := filepath.Join(filepath.Dir(b), "../..")
	viper.SetConfigName(".env.test")
	viper.SetConfigType("env")
	viper.AddConfigPath(Root)
	viper.AutomaticEnv()
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	viper.Set("app_debug", false)
	viper.Set("app_status", "test")

	return err
}

func DatabaseTesting() error {
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
	return err
}
