package cli

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/yourusername/youdu-app-mcp/internal/adapter"
)

func init() {
	rootCmd.AddCommand(testCmd)
}

var testCmd = &cobra.Command{
	Use:   "test",
	Short: "测试有度服务器连接和配置",
	Long:  `测试与有度服务器的连接，验证配置是否正确。`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := youduAdapter.GetConfig()
		fmt.Println("=== 有度连接测试 ===")
		fmt.Printf("服务器地址: %s\n", cfg.Youdu.Addr)
		fmt.Printf("企业总机号: %d\n", cfg.Youdu.Buin)
		fmt.Printf("应用ID: %s\n", cfg.Youdu.AppID)
		fmt.Println()

		// 测试获取部门列表
		fmt.Println("测试获取部门列表...")
		deptList, err := youduAdapter.GetDeptList(context.Background(), adapter.DeptListInput{DeptID: 0})
		if err != nil {
			return fmt.Errorf("获取部门列表失败: %w\n提示：请检查配置是否正确，服务器是否可访问", err)
		}

		fmt.Printf("✓ 成功获取 %d 个部门\n", len(deptList.Departments))
		if len(deptList.Departments) > 0 {
			fmt.Println("\n部门列表（前5个）：")
			for i, dept := range deptList.Departments {
				if i >= 5 {
					break
				}
				fmt.Printf("  - ID: %d, 名称: %s\n", dept.ID, dept.Name)
			}
		}

		fmt.Println("\n✓ 连接测试成功！配置正确。")
		return nil
	},
}
