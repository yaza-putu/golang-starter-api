package config

import "github.com/spf13/viper"

type db struct {
	Driver       string
	Host         string
	Port         int
	User         string
	Name         string
	Password     string
	Idle         int
	MaxConn      int
	ConnLifetime int
	AutoMigrate  bool
}

func DB() db {
	return db{
		Driver:       viper.GetString("db_driver"),
		Host:         viper.GetString("db_host"),
		Port:         viper.GetInt("db_port"),
		User:         viper.GetString("db_user"),
		Name:         viper.GetString("db_name"),
		Password:     viper.GetString("db_password"),
		Idle:         viper.GetInt("db_max_idle_conn"),
		MaxConn:      viper.GetInt("db_max_open_conn"),
		ConnLifetime: viper.GetInt("db_conn_lifetime"),
		AutoMigrate:  viper.GetBool("db_auto_migrate"),
	}
}
