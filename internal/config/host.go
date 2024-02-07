package config

import "github.com/spf13/viper"

type host struct {
	Name string
	Port int
}

// Host configuration
func Host() host {
	return host{
		Name: viper.GetString("host_name"),
		Port: viper.GetInt("host_port"),
	}
}
