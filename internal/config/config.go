package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
	"github.com/yourusername/youdu-app-mcp/internal/permission"
)

// Config 保存所有配置（YouDu + Permission）
type Config struct {
	Youdu      YouduConfig            `mapstructure:"youdu"`
	Permission *permission.Permission // 权限配置（由 config 包统一加载）

	viper *viper.Viper // 内部持有 viper 实例（不导出）
}

// YouduConfig 保存有度特定配置
type YouduConfig struct {
	Addr   string `mapstructure:"addr"`
	Buin   int    `mapstructure:"buin"`
	AppID  string `mapstructure:"app_id"`
	AesKey string `mapstructure:"aes_key"`
}

// LoadFromFile 从指定文件加载配置
// configPath 为空时使用默认搜索路径
func LoadFromFile(configPath string) (*Config, error) {
	v := newViper()

	// 设置配置文件
	if configPath != "" {
		v.SetConfigFile(configPath)
	} else {
		v.SetConfigName("config")
		v.AddConfigPath(".")
		v.AddConfigPath("$HOME/.youdu")
		v.AddConfigPath("/etc/youdu")
	}

	// 优先级 1: 读取配置文件（最高优先级）
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("读取配置失败: %w", err)
		}
		// 配置文件未找到；使用环境变量和默认值
	}

	// 优先级 2 & 3: 环境变量和默认值已在 newViper() 中设置

	// 解析主配置
	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("解析配置失败: %w", err)
	}

	// 加载权限配置
	perm, err := loadPermission(v)
	if err != nil {
		return nil, fmt.Errorf("加载权限配置失败: %w", err)
	}
	cfg.Permission = perm
	cfg.viper = v // 保存 viper 实例供后续使用

	return &cfg, nil
}

// Load 使用默认路径加载配置（向后兼容）
// 支持通过环境变量 YOUDU_CONFIG_FILE 指定配置文件路径
func Load() (*Config, error) {
	configFile := os.Getenv("YOUDU_CONFIG_FILE")
	return LoadFromFile(configFile)
}

// loadPermission 从 viper 加载权限配置（内部函数）
func loadPermission(v *viper.Viper) (*permission.Permission, error) {
	// 定义临时结构体用于解析配置文件
	type PermConfig struct {
		Enabled   bool                                      `mapstructure:"enabled"`
		AllowAll  bool                                      `mapstructure:"allow_all"`
		Resources map[string]permission.ResourcePolicy `mapstructure:"resources"`
	}

	var permCfg PermConfig

	// 从配置中读取（优先级: 配置文件 > 环境变量 > 默认值）
	if err := v.UnmarshalKey("permission", &permCfg); err != nil {
		// 使用默认配置（全部允许）
		permCfg = PermConfig{
			Enabled:  false,
			AllowAll: true,
			Resources: map[string]permission.ResourcePolicy{
				"dept":    {Create: true, Read: true, Update: true, Delete: true},
				"user":    {Create: true, Read: true, Update: true, Delete: true},
				"group":   {Create: true, Read: true, Update: true, Delete: true},
				"session": {Create: true, Read: true, Update: true, Delete: true},
				"message": {Create: true, Read: true, Update: true, Delete: true},
			},
		}
	}

	// 转换为 Permission 对象
	resources := make(map[permission.Resource]permission.ResourcePolicy)
	for k, v := range permCfg.Resources {
		resources[permission.Resource(k)] = v
	}

	// 使用 Permission 构造函数创建实例
	perm := permission.New(permCfg.Enabled, permCfg.AllowAll, resources)

	return perm, nil
}

// newViper 创建独立的 viper 实例
func newViper() *viper.Viper {
	v := viper.New()

	v.SetConfigType("yaml")

	// 环境变量配置（优先级低于配置文件）
	v.SetEnvPrefix("YOUDU")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	// 绑定环境变量
	bindEnvVars(v)

	// 设置默认值（优先级最低）
	setDefaults(v)

	return v
}

// bindEnvVars 绑定所有环境变量
func bindEnvVars(v *viper.Viper) {
	// YouDu 配置
	v.BindEnv("youdu.addr")
	v.BindEnv("youdu.buin")
	v.BindEnv("youdu.app_id")
	v.BindEnv("youdu.aes_key")

	// 权限配置
	v.BindEnv("permission.enabled")
	v.BindEnv("permission.allow_all")

	// 资源权限
	resources := []string{"user", "dept", "group", "session", "message"}
	actions := []string{"create", "read", "update", "delete"}
	for _, resource := range resources {
		for _, action := range actions {
			v.BindEnv(fmt.Sprintf("permission.resources.%s.%s", resource, action))
		}
	}
}

// setDefaults 设置默认值
func setDefaults(v *viper.Viper) {
	v.SetDefault("youdu.addr", "http://localhost:7080")
	v.SetDefault("youdu.buin", 0)

	// 权限默认值
	v.SetDefault("permission.enabled", false)
	v.SetDefault("permission.allow_all", true)
}

// GetPermission 获取权限配置
func (c *Config) GetPermission() *permission.Permission {
	return c.Permission
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
