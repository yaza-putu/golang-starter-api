package config

import "github.com/spf13/viper"

type db struct {
	Driver   string
	Host     string
	Port     int
	User     string
	Name     string
	Password string
}

func DB() db {
	return db{
		Driver:   viper.GetString("db_driver"),
		Host:     viper.GetString("db_host"),
		Port:     viper.GetInt("db_port"),
		User:     viper.GetString("db_user"),
		Name:     viper.GetString("db_name"),
		Password: viper.GetString("db_password"),
	}
}
