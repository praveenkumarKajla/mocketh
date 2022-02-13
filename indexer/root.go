package indexer

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const (
	erc20tokenCfgFileName = "erc20"
	erc20tokenCfgFileType = "yaml"
	erc20tokenCfgFilePath = "."
)

var (
	list      map[string]string
	addresses []string
	blocks    []int64
)

var Config = viper.New()

func init() {

	Config.SetConfigName(erc20tokenCfgFileName)
	Config.SetConfigType(erc20tokenCfgFileType)
	Config.AddConfigPath(erc20tokenCfgFilePath)

	if err := Config.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			logrus.Info("Added config file")
			// Config file not found; ignore error if desired

		} else {
			// Config file was found but another error was produced
			panic(fmt.Errorf("Fatal error config file: %s \n", err))
		}
	}
	list = Config.GetStringMapString(erc20tokenCfgFileName)

}

func LoadTokensFromConfig() map[string]string {
	return list
}
