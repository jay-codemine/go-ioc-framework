package config

import (
	"fmt"
	"github.com/spf13/viper"
)

var Config *viper.Viper

func Load(path string) *viper.Viper {
	v := viper.New()
	v.SetConfigFile(path)
	v.SetConfigType("yaml")

	if err := v.ReadInConfig(); err != nil {
		panic(fmt.Errorf("read config failed: %w", err))
	}

	Config = v
	return v
}
