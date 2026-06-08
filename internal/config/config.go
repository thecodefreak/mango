package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type Config struct {
	Server string `mapstructure:"server"`
	ApiToken string `mapstructure:"api_token"`
}

func Load() (*Config, error) {
	v := viper.New()

	homedir, err := os.UserHomeDir()
	if err != nil {
		fmt.Errorf("Unable to initialize config")
		return nil, err
	}

	configDir := filepath.Join(homedir, ".mango")
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(configDir)

	return &Config{}, err
}