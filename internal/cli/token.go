package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/spf13/cobra"
	"github.com/yourusername/youdu-app-mcp/internal/config"
)

var (
	tokenDescription string
	tokenExpiresIn   string
	tokenID          string
	tokenOutputJSON  bool
)

// tokenCmd represents the token command
var tokenCmd = &cobra.Command{
	Use:   "token",
	Short: "Token 管理",
	Long:  `管理 API token，包括生成、列出、撤销等操作。`,
}

// tokenGenerateCmd generates a new token
var tokenGenerateCmd = &cobra.Command{
	Use:   "generate",
	Short: "生成新的 token",
	Long: `生成一个新的 API token。

生成的 token 会输出到控制台，需要手动添加到配置文件中。

示例:
  youdu-cli token generate --description "API token for service A"
  youdu-cli token generate --description "Temporary token" --expires-in 24h
  youdu-cli token generate --description "Test token" --json`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// 加载配置以获取数据库连接
		var cfg *config.Config
		var err error

		// 如果指定了配置文件，从文件加载
		if cfgFile != "" {
			cfg, err = config.LoadFromFile(cfgFile)
		} else {
			cfg, err = config.Load()
		}

		if err != nil {
			return fmt.Errorf("加载配置失败: %w", err)
		}

		// 解析过期时间
		var expiresIn *time.Duration
		if tokenExpiresIn != "" {
			duration, err := time.ParseDuration(tokenExpiresIn)
			if err != nil {
				return fmt.Errorf("无效的过期时间格式: %w\n示例: 24h, 7d, 30d", err)
			}
			expiresIn = &duration
		}

		// 生成 token
		token, err := cfg.TokenManager.Generate(tokenDescription, expiresIn)
		if err != nil {
			return fmt.Errorf("生成 token 失败: %w", err)
		}

		// 输出 token
		if tokenOutputJSON {
			output, _ := json.MarshalIndent(token, "", "  ")
			fmt.Println(string(output))
		} else {
			fmt.Println("\n✅ Token 生成成功！")
			fmt.Println("\n📋 Token 信息:")
			fmt.Printf("  ID:          %s\n", token.ID)
			fmt.Printf("  Value:       %s\n", token.Value)
			fmt.Printf("  Description: %s\n", token.Description)
			fmt.Printf("  Created At:  %s\n", token.CreatedAt.Format(time.RFC3339))
			if token.ExpiresAt != nil {
				fmt.Printf("  Expires At:  %s\n", token.ExpiresAt.Format(time.RFC3339))
			} else {
				fmt.Printf("  Expires At:  永不过期\n")
			}

			fmt.Println("\n💡 提示:")
			fmt.Println("  Token 已保存到数据库中。")
			fmt.Println("  确保配置文件中 token.enabled 设置为 true 以启用认证。")
		}

		return nil
	},
}

// tokenListCmd lists all tokens
var tokenListCmd = &cobra.Command{
	Use:   "list",
	Short: "列出所有 token",
	Long: `列出配置中的所有 token。

示例:
  youdu-cli token list
  youdu-cli token list --json`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// 加载配置以获取 token 管理器
		var cfg *config.Config
		var err error

		// 如果指定了配置文件，从文件加载
		if cfgFile != "" {
			cfg, err = config.LoadFromFile(cfgFile)
		} else {
			cfg, err = config.Load()
		}

		if err != nil {
			return fmt.Errorf("加载配置失败: %w", err)
		}

		tokens := cfg.TokenManager.List()

		if len(tokens) == 0 {
			fmt.Println("📭 没有配置任何 token")
			fmt.Println("\n💡 提示: 使用 'youdu-cli token generate' 生成新 token")
			return nil
		}

		if tokenOutputJSON {
			output, _ := json.MarshalIndent(tokens, "", "  ")
			fmt.Println(string(output))
		} else {
			fmt.Printf("\n📋 Token 列表 (共 %d 个):\n\n", len(tokens))

			w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
			fmt.Fprintln(w, "ID\tDescription\tCreated At\tExpires At\tStatus")
			fmt.Fprintln(w, "---\t---\t---\t---\t---")

			for _, token := range tokens {
				expiresAt := "永不过期"
				status := "✅ 有效"

				if token.ExpiresAt != nil {
					expiresAt = token.ExpiresAt.Format("2006-01-02 15:04:05")
					if time.Now().After(*token.ExpiresAt) {
						status = "❌ 已过期"
					}
				}

				fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
					token.ID,
					token.Description,
					token.CreatedAt.Format("2006-01-02 15:04:05"),
					expiresAt,
					status,
				)
			}

			w.Flush()
			fmt.Println()
		}

		return nil
	},
}

// tokenRevokeCmd revokes a token
var tokenRevokeCmd = &cobra.Command{
	Use:   "revoke",
	Short: "撤销 token",
	Long: `通过 ID 撤销一个 token。

Token 将从数据库中永久删除。

示例:
  youdu-cli token revoke --id abc123`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if tokenID == "" {
			return fmt.Errorf("请使用 --id 参数指定要撤销的 token ID")
		}

		// 加载配置以获取 token 管理器
		var cfg *config.Config
		var err error

		// 如果指定了配置文件，从文件加载
		if cfgFile != "" {
			cfg, err = config.LoadFromFile(cfgFile)
		} else {
			cfg, err = config.Load()
		}

		if err != nil {
			return fmt.Errorf("加载配置失败: %w", err)
		}

		// 撤销 token
		if err := cfg.TokenManager.RevokeByID(tokenID); err != nil {
			return fmt.Errorf("撤销 token 失败: %w", err)
		}

		fmt.Printf("✅ Token %s 已撤销\n", tokenID)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(tokenCmd)

	// token generate
	tokenCmd.AddCommand(tokenGenerateCmd)
	tokenGenerateCmd.Flags().StringVarP(&tokenDescription, "description", "d", "", "Token 描述")
	tokenGenerateCmd.Flags().StringVar(&tokenExpiresIn, "expires-in", "", "过期时间 (例如: 24h, 7d, 30d)")
	tokenGenerateCmd.Flags().BoolVar(&tokenOutputJSON, "json", false, "以 JSON 格式输出")
	tokenGenerateCmd.MarkFlagRequired("description")

	// token list
	tokenCmd.AddCommand(tokenListCmd)
	tokenListCmd.Flags().BoolVar(&tokenOutputJSON, "json", false, "以 JSON 格式输出")

	// token revoke
	tokenCmd.AddCommand(tokenRevokeCmd)
	tokenRevokeCmd.Flags().StringVar(&tokenID, "id", "", "要撤销的 token ID")
	tokenRevokeCmd.MarkFlagRequired("id")
}
