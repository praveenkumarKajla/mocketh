package config

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var Config = viper.New()

func init() {

	Config.SetConfigName("server")
	Config.SetConfigType("yaml")
	Config.AddConfigPath(".")
	Config.SetConfigFile("config.yaml")

	if err := Config.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			logrus.Info("Added config file")
			// Config file not found; ignore error if desired

		} else {
			// Config file was found but another error was produced
			panic(fmt.Errorf("Fatal error config file: %s \n", err))
		}
	}

}
