package config

import (
	"github.com/spf13/viper"
)

type app struct {
	Name   string
	Lang   string
	Debug  bool
	Status string
}

func App() app {
	return app{
		Name:   viper.GetString("app_name"),
		Lang:   viper.GetString("app_lang"),
		Debug:  viper.GetBool("app_debug"),
		Status: viper.GetString("app_status"),
	}
}
