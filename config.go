package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	CoreHost     string   `mapstructure:"CORE_HOST"`
	PrivateToken string   `mapstructure:"PRIVATE_TOKEN"`
	ReplayPaths  []string `mapstructure:"REPLAY_PATHS"`

	OverWriteReplay bool `mapstructure:"OverWriteReplay"`
}

func LoadConfig(cfgPath string) *Config {
	cfg := Config{
		CoreHost:     "http://localhost:8080",
		PrivateToken: "",
		ReplayPaths:  []string{},
	}
	loadConfigFromFile(cfgPath, &cfg)
	slog.Info(fmt.Sprintf("Config: %+v", cfg))
	return &cfg
}

func loadConfigFromFile(path string, conf *Config) {
	var err error
	if have(path) {
		fileViper := viper.New()
		fileViper.SetConfigFile(path)
		if err = fileViper.ReadInConfig(); err == nil {
			if err = fileViper.Unmarshal(conf); err == nil {
				slog.Info(fmt.Sprintf("Load config from %s success", path))
				return
			}
		}
	}
	if err != nil {
		slog.Error(fmt.Sprintf("Load config from %s failed: %s", path, err))
		os.Exit(1)
	}
}

func have(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
