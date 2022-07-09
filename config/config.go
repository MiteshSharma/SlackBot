package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	LoggerConfig LoggerConfig `mapstructure:"logger"`
	SlackConfig  SlackConfig  `mapstructure:"slack"`
}

// LoggerConfig has logger related configuration.
type LoggerConfig struct {
	LogFilePath string `mapstructure:"file"`
}

type SlackConfig struct {
	Token       string `mapstructure:"token"`
	ChannelName string `mapstructure:"channelName"`
}

func GetConfig() *Config {
	appConfig := &Config{}
	GetConfigFromFile1("default", appConfig)
	return appConfig
}

func GetConfigFromFile(file string) *Config {
	appConfig := &Config{}
	GetConfigFromFile1(file, appConfig)
	return appConfig
}

func GetConfigFromFile1(fileName string, config *Config) {
	if fileName == "" {
		fileName = "default"
	}
	viper.SetConfigName(fileName)
	viper.AddConfigPath("../conf/")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./conf/")
	viper.SetConfigType("toml")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Printf("couldn't load config: %s", err)
		os.Exit(1)
	}

	err = viper.Unmarshal(&config)
	fmt.Println(config)
	if err != nil {
		fmt.Printf("couldn't read config: %s", err)
		os.Exit(1)
	}
}
