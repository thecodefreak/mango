package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type Config struct {
	Server   string `mapstructure:"server"`
	ApiToken string `mapstructure:"api_token"`
}

type ServerConfig struct {
	Addr         string `mapstructure:"server_addr"`
	ServerToken  string `mapstructure:"server_token"`
	DocumentRoot string `mapstructure:"document_root"`
}

func initConfig(fileName string) (v *viper.Viper, confDir string, err error) {
	viper := viper.New()

	homedir, err := os.UserHomeDir()
	if err != nil {
		return nil, "", err
	}

	configDir := filepath.Join(homedir, ".mango")
	if err := os.MkdirAll(configDir, 0755); err != nil && !os.IsExist(err) {
		return nil, "", fmt.Errorf("Unable to create config directory: %s", configDir)
	}

	viper.AddConfigPath(configDir)
	viper.SetConfigName(fileName)
	viper.SetConfigType("yaml")

	return viper, configDir, nil
}

func checkConfigFile(v *viper.Viper, cd string, fn string) error {
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			err := v.SafeWriteConfigAs(filepath.Join(cd, fn+".yaml"))
			if err != nil {
				return fmt.Errorf("Error creating config file: %s", err)
			}
		} else {
			return fmt.Errorf("Error reading config file: %s", err)
		}
	}
	return nil
}

func Load() (*Config, error) {
	v, configDir, err := initConfig("config")
	if err != nil {
		return nil, fmt.Errorf("Unable to initialise config, %w", err)
	}

	v.SetDefault("server", "https://mango.example.com")
	v.SetDefault("api_token", "api_token_here")

	err = checkConfigFile(v, configDir, "config")
	if err != nil {
		fmt.Printf("Error checking config file: %s", err)
		return nil, err
	}

	var cfg Config

	if err := v.Unmarshal(&cfg); err != nil {
		fmt.Printf("Error parsing config file: %s\n", err)
		return nil, err
	}

	return &cfg, err
}

func LoadServerConf() (*ServerConfig, error) {
	v, configDir, err := initConfig("server")
	if err != nil {
		return nil, fmt.Errorf("Unable to initialise config, %w", err)
	}

	userHome, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("Unable to get user home directory, %w", err)
	}

	v.SetDefault("server_addr", ":3000")
	v.SetDefault("server_token", "")
	v.SetDefault("document_root", userHome+"/.mango/static_pages")

	err = checkConfigFile(v, configDir, "server")
	if err != nil {
		fmt.Printf("Error checking config file: %s", err)
		return nil, err
	}

	var serverCfg ServerConfig
	if err := v.Unmarshal(&serverCfg); err != nil {
		fmt.Printf("Error parsing config file: %s\n", err)
		return nil, err
	}

	return &serverCfg, nil
}
