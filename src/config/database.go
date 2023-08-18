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
		Driver:   viper.GetString("database.driver"),
		Host:     viper.GetString("database.host"),
		Port:     viper.GetInt("database.port"),
		User:     viper.GetString("database.user"),
		Name:     viper.GetString("database.name"),
		Password: viper.GetString("database.password"),
	}
}
