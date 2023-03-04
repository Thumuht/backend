/*
Package config configs thumuht using viper.
*/
package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// use viper to config this project.
// todo(wj, low): make this thing more human friendly
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
