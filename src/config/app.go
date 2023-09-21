package config

import (
	"github.com/spf13/viper"
)

type app struct {
	Name        string
	Debug       bool
	Environment string
	Lang        string
}

func App() app {
	return app{
		Name:        viper.GetString("app.name"),
		Debug:       viper.GetBool("app.debug"),
		Environment: viper.GetString("app.environment"),
		Lang:        viper.GetString("app.lang"),
	}
}
