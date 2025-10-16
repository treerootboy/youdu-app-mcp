package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// Config 保存有度配置
type Config struct {
	Youdu YouduConfig `mapstructure:"youdu"`
}

// YouduConfig 保存有度特定配置
type YouduConfig struct {
	Addr   string `mapstructure:"addr"`
	Buin   int    `mapstructure:"buin"`
	AppID  string `mapstructure:"app_id"`
	AesKey string `mapstructure:"aes_key"`
}

// Load 从配置文件和环境变量读取配置
func Load() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME/.youdu")
	viper.AddConfigPath("/etc/youdu")

	// 环境变量
	viper.SetEnvPrefix("YOUDU")
	viper.AutomaticEnv()

	// 设置默认值
	viper.SetDefault("youdu.addr", "http://localhost:7080")
	viper.SetDefault("youdu.buin", 0)

	// 读取配置文件（如果存在）
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("读取配置失败: %w", err)
		}
		// 配置文件未找到；使用默认值和环境变量
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("解析配置失败: %w", err)
	}

	return &cfg, nil
}

// Validate 检查配置是否有效
func (c *Config) Validate() error {
	if c.Youdu.Addr == "" {
		return fmt.Errorf("youdu.addr 为必填项")
	}
	if c.Youdu.Buin == 0 {
		return fmt.Errorf("youdu.buin 为必填项")
	}
	if c.Youdu.AppID == "" {
		return fmt.Errorf("youdu.app_id 为必填项")
	}
	if c.Youdu.AesKey == "" {
		return fmt.Errorf("youdu.aes_key 为必填项")
	}
	return nil
}
