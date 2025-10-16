package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// Config holds the Youdu configuration
type Config struct {
	Youdu YouduConfig `mapstructure:"youdu"`
}

// YouduConfig holds Youdu-specific settings
type YouduConfig struct {
	Addr   string `mapstructure:"addr"`
	Buin   int    `mapstructure:"buin"`
	AppID  string `mapstructure:"app_id"`
	AesKey string `mapstructure:"aes_key"`
}

// Load reads configuration from file and environment variables
func Load() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME/.youdu")
	viper.AddConfigPath("/etc/youdu")

	// Environment variables
	viper.SetEnvPrefix("YOUDU")
	viper.AutomaticEnv()

	// Set defaults
	viper.SetDefault("youdu.addr", "http://localhost:7080")
	viper.SetDefault("youdu.buin", 0)

	// Read config file if exists
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read config: %w", err)
		}
		// Config file not found; using defaults and env vars
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &cfg, nil
}

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
	if c.Youdu.Addr == "" {
		return fmt.Errorf("youdu.addr is required")
	}
	if c.Youdu.Buin == 0 {
		return fmt.Errorf("youdu.buin is required")
	}
	if c.Youdu.AppID == "" {
		return fmt.Errorf("youdu.app_id is required")
	}
	if c.Youdu.AesKey == "" {
		return fmt.Errorf("youdu.aes_key is required")
	}
	return nil
}
