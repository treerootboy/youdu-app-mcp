package main

import (
	"context"
	"fmt"
	"os"

	"github.com/yourusername/youdu-app-mcp/internal/config"
	"github.com/yourusername/youdu-app-mcp/internal/mcp"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		return fmt.Errorf("invalid config: %w", err)
	}

	// Create MCP server
	server, err := mcp.New(cfg)
	if err != nil {
		return fmt.Errorf("failed to create MCP server: %w", err)
	}

	// Run server
	ctx := context.Background()
	return server.Run(ctx)
}
