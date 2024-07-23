// File: config/config.go
package config

import (
	"path/filepath"
	"runtime"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func Init() {
	viper.SetConfigName("config")
	viper.SetConfigType("yml")

	// Get the directory of the current file
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		logrus.Fatal("No caller information")
	}
	configDir := filepath.Join(filepath.Dir(filename), "..", "..", "config")

	viper.AddConfigPath(configDir)

	err := viper.ReadInConfig()
	if err != nil {
		logrus.Fatalf("Fatal error reading config file: %s", err)
	}
}

func GetImageProviderAccessKey() string {
	return viper.GetString("imageprovider.accessKey")
}
