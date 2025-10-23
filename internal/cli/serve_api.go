package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/yourusername/youdu-app-mcp/internal/api"
	"github.com/yourusername/youdu-app-mcp/internal/config"
)

var (
	apiPort string
)

// serveAPICmd represents the serve-api command
var serveAPICmd = &cobra.Command{
	Use:   "serve-api",
	Short: "启动 HTTP API 服务器",
	Long: `启动 HTTP REST API 服务器，自动将所有 adapter 方法暴露为 HTTP endpoints。

所有业务 API 的路径格式为: POST /api/v1/{method_name}

示例:
  youdu-cli serve-api
  youdu-cli serve-api --port 8080
  youdu-cli serve-api --config config.yaml --port 9000

服务启动后可以访问:
  - GET /health - 健康检查
  - GET /api/v1/endpoints - 查看所有可用 API
  - POST /api/v1/* - 调用各种业务 API`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// 加载配置
		var cfg *config.Config
		var err error
		if cfgFile != "" {
			cfg, err = config.LoadFromFile(cfgFile)
		} else {
			cfg, err = config.Load()
		}
		if err != nil {
			return fmt.Errorf("加载配置失败: %w", err)
		}

		// 验证配置
		if err := cfg.Validate(); err != nil {
			return fmt.Errorf("配置无效: %w\n提示：请检查 config.yaml 文件或设置环境变量", err)
		}

		// 创建 API 服务器
		server, err := api.New(cfg)
		if err != nil {
			return fmt.Errorf("创建 API 服务器失败: %w", err)
		}
		defer server.Close()

		// 构建监听地址
		addr := fmt.Sprintf(":%s", apiPort)

		// 启动服务器
		fmt.Printf("\n🎉 YouDu HTTP API 服务器已就绪！\n\n")
		if err := server.Start(addr); err != nil {
			return fmt.Errorf("启动服务器失败: %w", err)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(serveAPICmd)

	// 添加端口参数
	serveAPICmd.Flags().StringVarP(&apiPort, "port", "p", "8080", "HTTP API 服务器监听端口")
}
