package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/yourusername/youdu-app-mcp/internal/adapter"
	"github.com/yourusername/youdu-app-mcp/internal/config"
)

var (
	cfgFile     string
	youduAdapter *adapter.Adapter
)

// rootCmd represents the base command
var rootCmd = &cobra.Command{
	Use:   "youdu-cli",
	Short: "Youdu IM CLI tool",
	Long:  `A command-line interface for interacting with Youdu IM system.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// 跳过 token 命令的配置验证（token 命令不需要 YouDu 配置）
		if cmd.Parent() != nil && cmd.Parent().Name() == "token" {
			return nil
		}
		if cmd.Name() == "token" {
			return nil
		}

		// 加载配置
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("加载配置失败: %w", err)
		}

		// 验证配置
		if err := cfg.Validate(); err != nil {
			return fmt.Errorf("配置无效: %w\n提示：请检查 config.yaml 文件或设置环境变量", err)
		}

		// 创建适配器
		youduAdapter, err = adapter.New(cfg)
		if err != nil {
			return fmt.Errorf("创建适配器失败: %w", err)
		}

		return nil
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		if youduAdapter != nil {
			youduAdapter.Close()
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is ./config.yaml)")

	// Generate commands from adapter methods
	if err := generateCommands(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to generate commands: %v\n", err)
		os.Exit(1)
	}
}
