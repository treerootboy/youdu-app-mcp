package adapter

import (
	"context"

	"github.com/addcnos/youdu/v2"
	"github.com/yourusername/youdu-app-mcp/internal/config"
)

// Adapter wraps the Youdu client and provides simplified methods
type Adapter struct {
	client *youdu.Client
	config *config.Config
}

// New creates a new Adapter instance
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

// Close closes the adapter and releases resources
func (a *Adapter) Close() error {
	// Currently, the youdu client doesn't require cleanup
	return nil
}

// Context returns a default context
func (a *Adapter) Context() context.Context {
	return context.Background()
}
