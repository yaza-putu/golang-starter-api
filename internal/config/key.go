package config

import "github.com/spf13/viper"

type key struct {
	Token      string
	Refresh    string
	Passphrase string
}

func Key() key {
	return key{
		Token:      viper.GetString("key_token"),
		Refresh:    viper.GetString("key_refresh"),
		Passphrase: viper.GetString("key_passphrase"),
	}
}
