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
		fmt.Println("Unable to initialize config")
		return nil, err
	}

	configDir := filepath.Join(homedir, ".mango")
	if err := os.MkdirAll(configDir, 0755); err != nil && !os.IsExist(err) {
		fmt.Printf("Unable to create config directory: %s\n", configDir)
		return nil, err
	}

	v.SetDefault("server", "https://mango.example.com")
	v.SetDefault("api_token", "api_token_here")
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(configDir)

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			err := v.SafeWriteConfigAs(filepath.Join(configDir, "config.yaml"))
			if err != nil {

				fmt.Printf("Error creating config file: %s\n", err)
				return nil, err
			}
		} else {
			fmt.Printf("Error reading config file: %s\n", err)
			return nil, err
		}
	}

	var cfg Config

	if err := v.Unmarshal(&cfg); err != nil {
		fmt.Printf("Error parsing config file: %s\n", err)
		return nil, err
	}

	return &cfg, err
}