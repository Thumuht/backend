package config

import (
	"fmt"

	"github.com/spf13/viper"
)

func ConfigProject() error {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		return fmt.Errorf("read config file failed. Please Config your App before Using it.\n\t%w", err)
	}
	return nil
}
