package config

import "github.com/spf13/viper"

type redis struct {
	Host     string
	Port     int
	Password string
	DB       int
}

func Redis() redis {
	return redis{
		Host:     viper.GetString("redis_host"),
		Port:     viper.GetInt("redis_port"),
		Password: viper.GetString("redis_password"),
		DB:       viper.GetInt("redis_db"),
	}
}
