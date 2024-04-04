package config

import (
	"github.com/spf13/viper"
)

// The top level object holding all configuraiton options set in the application.
type Configuration struct {
	App ApplicationConfiguration  `json:"app"`
	DB  DatabaseConfiguration     `json:"db"`
	MQ  MessageQueueConfiguration `json:"mq"`
	Rq  RequestConfiguration      `json:"rq"`
}

// The part of the configuration dedicated to manage the top level and shared
// configuration options.
type ApplicationConfiguration struct {
	Env  string `json:"env"`
	Port int    `json:"port"`
}

type DatabaseConfiguration struct {
	DSN          string `json:"-"`
	MaxOpenConns int    `json:"max-open-conns"`
	MaxIdleConns int    `json:"max-idle-conns"`
	MaxIdleTime  string `json:"max-idle-time"`
	Timeout      int    `json:"timeout"`
}

type MessageQueueConfiguration struct {
	DSN      string `json:"-"`
	MaxConns int    `json:"max-conns"`
}

type RequestConfiguration struct {
	URL string `json:"url"`
}

func New() (*Configuration, error) {
	viper.AutomaticEnv()
	viper.AllowEmptyEnv(false)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/app/")
	viper.AddConfigPath("$HOME/.config/modularmonolith")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	var config Configuration
	err = viper.Unmarshal(&config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
