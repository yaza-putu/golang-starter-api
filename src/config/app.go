package config

import (
	"github.com/spf13/viper"
)

type app struct {
	Name        string
	Key         string
	Debug       bool
	Environment string
}

func App() app {
	return app{
		Name:        viper.GetString("app.name"),
		Key:         viper.GetString("app.key"),
		Debug:       viper.GetBool("app.debug"),
		Environment: viper.GetString("app.environment"),
	}
}
