package core

import (
	"fmt"
	"net"
	"path/filepath"
	"runtime"

	"github.com/spf13/viper"
	"github.com/yaza-putu/golang-starter-api/internal/config"
	"github.com/yaza-putu/golang-starter-api/internal/database"
	_ "github.com/yaza-putu/golang-starter-api/internal/database/migrations"
	_ "github.com/yaza-putu/golang-starter-api/internal/database/seeders"
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
	// force auto migrate true
	viper.Set("db_auto_migrate", true)

	// call database
	Database()

	// run seeder data
	database.SeederUp()
	// run server
	if !isPortActive() {
		go HttpServe()
	}

	return err
}

func EnvRollback() {
	database.MigrationDown()
}

func isPortActive() bool {
	addr := fmt.Sprintf(":%d", config.Host().Port)
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return true
	}
	ln.Close()
	return false
}
