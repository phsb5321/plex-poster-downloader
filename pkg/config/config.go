// File: config/config.go
package config

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func Init() {
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath("./../config")

	err := viper.ReadInConfig()
	if err != nil {
		logrus.Fatalf("Fatal error reading config file: %s", err)
	}
}

func GetBaseUrl() string {
	return viper.GetString("baseUrl")
}
