package adapter

import (
	"context"

	"github.com/addcnos/youdu/v2"
	"github.com/yourusername/youdu-app-mcp/internal/config"
)

// Adapter 封装有度客户端并提供简化的方法
type Adapter struct {
	client *youdu.Client
	config *config.Config
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
		client: client,
		config: cfg,
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
