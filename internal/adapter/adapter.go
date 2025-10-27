package adapter

import (
	"context"

	"github.com/addcnos/youdu/v2"
	"github.com/yourusername/youdu-app-mcp/internal/config"
	"github.com/yourusername/youdu-app-mcp/internal/permission"
)

// Adapter 封装有度客户端并提供简化的方法
type Adapter struct {
	client     *youdu.Client
	config     *config.Config
	permission *permission.Permission // 权限实例（从 Config 获取）
}

// New 创建一个新的 Adapter 实例
func New(cfg *config.Config) (*Adapter, error) {
	client := youdu.NewClient(&youdu.Config{
		Addr:   cfg.Youdu.Addr,
		Buin:   cfg.Youdu.Buin,
		AppID:  cfg.Youdu.AppID,
		AesKey: cfg.Youdu.AesKey,
	})

	return &Adapter{
		client:     client,
		config:     cfg,
		permission: cfg.GetPermission(), // 从 Config 获取权限配置
	}, nil
}

// Close 关闭适配器并释放资源
func (a *Adapter) Close() error {
	// 目前，有度客户端不需要清理
	return nil
}

// Context 返回默认上下文
func (a *Adapter) Context() context.Context {
	return context.Background()
}

// GetConfig 返回配置信息
func (a *Adapter) GetConfig() *config.Config {
	return a.config
}

// GetPermission 返回权限配置（供 CLI 命令使用）
func (a *Adapter) GetPermission() *permission.Permission {
	return a.permission
}

// checkPermission 检查操作权限（使用实例方法）
func (a *Adapter) checkPermission(resource permission.Resource, action permission.Action) error {
	return a.permission.Check(resource, action)
}

// checkPermissionWithID 检查操作权限（包含行级权限）
func (a *Adapter) checkPermissionWithID(resource permission.Resource, action permission.Action, resourceID string) error {
	return a.permission.CheckWithID(resource, action, resourceID)
}

// checkMessageSendPermission 检查消息发送权限
func (a *Adapter) checkMessageSendPermission(toUser, toDept string) error {
	return a.permission.CheckMessageSend(toUser, toDept)
}
