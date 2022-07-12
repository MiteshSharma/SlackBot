package config

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	LoggerConfig   LoggerConfig   `mapstructure:"logger"`
	SlackConfig    SlackConfig    `mapstructure:"slack"`
	ServerConfig   ServerConfig   `mapstructure:"server"`
	DatabaseConfig DatabaseConfig `mapstructure:"database"`
}

// LoggerConfig has logger related configuration.
type LoggerConfig struct {
	LogFilePath string `mapstructure:"file"`
}

type SlackConfig struct {
	Token         string `mapstructure:"token"`
	ChannelName   string `mapstructure:"channelName"`
	SigningSecret string `mapstructure:"signingSecret"`
	ClientID      string `mapstructure:"clientId"`
	ClientSecret  string `mapstructure:"clientSecret"`
}

type ServerConfig struct {
	Port         string        `mapstructure:"port"`
	ReadTimeout  time.Duration `mapstructure:"readTimeout"`
	WriteTimeout time.Duration `mapstructure:"writeTimeout"`
}

// DatabaseConfig has database related configuration.
type DatabaseConfig struct {
	Type             string `mapstructure:"type"`
	Host             string `mapstructure:"host"`
	DbName           string `mapstructure:"dbName"`
	UserName         string `mapstructure:"userName"`
	Password         string `mapstructure:"password"`
	ConnectionString string `mapstructure:"connectionString"`
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
