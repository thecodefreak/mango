package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
	"github.com/thecodefreak/mango/internal/helpers"
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

func InitConfig() *viper.Viper {
	viper := viper.New()
	viper.SetConfigType("yaml")
	return viper
}

func CreateConfig(v *viper.Viper, configFileName string) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("Unable to get user home directory for config init")
	}

	configDir := filepath.Join(homeDir, ".mango")
	configPath := filepath.Join(configDir, configFileName+".yaml")

	if !helpers.IsFileExist(configDir) {
		if err := os.MkdirAll(configDir, 0700); err != nil {
			return fmt.Errorf("Unable to create config directory: %s", configDir)
		}

		err = v.SafeWriteConfigAs(configPath)
		if err != nil {
			return fmt.Errorf("Error creating config file: %s", err)
		}

		fmt.Printf("Config file successfully created at: %s", configPath)
	}

	v.SetConfigFile(configPath)

	return nil
}

func readConfig(v *viper.Viper) error {
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil
		} else {
			return fmt.Errorf("Error reading config file: %s", err)
		}
	}
	return nil
}

func Load() (*Config, error) {
	v := InitConfig()

	v.SetDefault("server", "https://mango.example.com")
	v.SetDefault("api_token", "api_token_here")
	err := CreateConfig(v, "config")
	if err != nil {
		return nil, fmt.Errorf("Config init failed: %w", err)
	}

	err = readConfig(v)
	if err != nil {
		return nil, fmt.Errorf("Failed to read config: %w", err)
	}

	var cfg Config

	if err := v.Unmarshal(&cfg); err != nil {
		fmt.Printf("Error parsing config file: %s\n", err)
		return nil, err
	}

	return &cfg, err
}

func LoadServerConf() (*ServerConfig, error) {
	v := InitConfig()

	v.SetDefault("server_addr", ":3000")
	v.SetDefault("server_token", "")
	v.SetDefault("document_root", "")
	v.SetEnvPrefix("MANGO")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	envConfFile := os.Getenv("MANGO_CONFIG_FILE")
	if envConfFile != "" {
		if !helpers.IsFileExist(envConfFile) {
			return nil, fmt.Errorf("Config file %s does not exist", envConfFile)
		}
		v.SetConfigFile(envConfFile)
	} else {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("Config init failed: %w", err)
		}
		v.SetConfigFile(homeDir)
		v.SetConfigName("server")
	}

	keys := []string{
		"server_addr",
		"server_token",
		"document_root",
	}

	for _, key := range keys {
		if err := v.BindEnv(key); err != nil {
			return nil, fmt.Errorf("failed to bind env %s: %w", key, err)
		}
	}

	err := readConfig(v)
	if err != nil {
		return nil, fmt.Errorf("Failed to read config: %w", err)
	}

	var serverCfg ServerConfig
	if err := v.Unmarshal(&serverCfg); err != nil {
		fmt.Printf("Error parsing config file: %s\n", err)
		return nil, err
	}

	return &serverCfg, nil
}
