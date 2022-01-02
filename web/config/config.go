package config

import (
	"time"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type Timeouts struct {
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

func Read() Timeouts {

	// these are defined in ../database/.env (file is not checked in, it needs to be created manually)
	viper.BindEnv("postgres_user")
	viper.BindEnv("postgres_password")
	viper.BindEnv("postgres_db")

	// look for environment variables that start with ROBOZ_
	viper.SetEnvPrefix("roboz")

	// environment variable ROBOZ_CONFIG specifies configuration file name
	viper.SetDefault("config", "roboz.yaml")
	path := viper.GetString("config")
	viper.SetConfigFile(path)

	if err := viper.ReadInConfig(); err != nil {
		zap.S().Fatalf("unable to load configuration %+v", err)
	}

	timeouts := &Timeouts{}
	if err := viper.Unmarshal(timeouts); err != nil {
		zap.S().Fatalf("unable to decode into timeouts struct, %+v", err)
	}

	return *timeouts
}
