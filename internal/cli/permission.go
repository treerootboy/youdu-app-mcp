package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/yourusername/youdu-app-mcp/internal/permission"
)

func init() {
	rootCmd.AddCommand(permissionCmd)
	permissionCmd.AddCommand(permStatusCmd)
	permissionCmd.AddCommand(permListCmd)
}

var permissionCmd = &cobra.Command{
	Use:   "permission",
	Short: "权限管理命令",
	Long:  `查看和管理 API 操作权限配置。`,
}

var permStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "查看权限系统状态",
	Long:  `显示权限系统是否启用以及当前配置。`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// 从 Adapter 获取权限配置
		perm := youduAdapter.GetPermission()

		fmt.Println("=== 权限系统状态 ===")
		if perm.IsEnabled() {
			fmt.Println("状态: ✓ 已启用")
		} else {
			fmt.Println("状态: ✗ 未启用（所有操作都允许）")
		}
		fmt.Println()

		return nil
	},
}

var permListCmd = &cobra.Command{
	Use:   "list",
	Short: "列出所有资源的权限配置",
	Long:  `显示部门、用户、群组、会话、消息的 CRUD 权限配置。`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// 从 Adapter 获取权限配置
		perm := youduAdapter.GetPermission()

		fmt.Println("=== 资源权限配置 ===")
		fmt.Println()

		resources := []permission.Resource{
			permission.ResourceDept,
			permission.ResourceUser,
			permission.ResourceGroup,
			permission.ResourceSession,
			permission.ResourceMessage,
		}

		resourceNames := map[permission.Resource]string{
			permission.ResourceDept:    "部门 (dept)",
			permission.ResourceUser:    "用户 (user)",
			permission.ResourceGroup:   "群组 (group)",
			permission.ResourceSession: "会话 (session)",
			permission.ResourceMessage: "消息 (message)",
		}

		for _, res := range resources {
			policy, exists := perm.GetResourcePolicy(res)
			fmt.Printf("【%s】\n", resourceNames[res])

			if !exists {
				fmt.Println("  未配置（默认拒绝所有操作）")
			} else {
				fmt.Printf("  创建 (create): %s\n", formatPermission(policy.Create))
				fmt.Printf("  读取 (read):   %s\n", formatPermission(policy.Read))
				fmt.Printf("  更新 (update): %s\n", formatPermission(policy.Update))
				fmt.Printf("  删除 (delete): %s\n", formatPermission(policy.Delete))
			}
			fmt.Println()
		}

		if !perm.IsEnabled() {
			fmt.Println("提示：权限检查未启用，以上配置不生效。")
			fmt.Println("要启用权限检查，请在 config.yaml 中设置 permission.enabled: true")
		}

		return nil
	},
}

func formatPermission(allowed bool) string {
	if allowed {
		return "✓ 允许"
	}
	return "✗ 拒绝"
}
