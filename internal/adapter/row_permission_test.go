package adapter

import (
	"context"
	"strings"
	"testing"

	"github.com/yourusername/youdu-app-mcp/internal/adapter/testdata"
	"github.com/yourusername/youdu-app-mcp/internal/config"
	"github.com/yourusername/youdu-app-mcp/internal/permission"
)

// TestAdapter_RowLevelPermissions 测试行级权限功能
func TestAdapter_RowLevelPermissions(t *testing.T) {
	// 创建测试配置
	cfg, err := config.LoadFromFile("../../config_test.yaml")
	if err != nil {
		t.Fatalf("加载测试配置失败: %v", err)
	}

	// 启动 Mock Server
	mockServer := testdata.NewMockYouDuServer(cfg.Youdu.AesKey, cfg.Youdu.AppID)
	defer mockServer.Close()
	cfg.Youdu.Addr = mockServer.URL()

	// 创建适配器
	adapter, err := New(cfg)
	if err != nil {
		t.Fatalf("创建适配器失败: %v", err)
	}

	t.Run("未配置allowlist时，所有用户都可以访问", func(t *testing.T) {
		// GetUser 应该成功（read 权限为 true，allowlist 为空）
		_, err := adapter.GetUser(context.Background(), GetUserInput{
			UserID: "user123",
		})
		if err != nil {
			t.Errorf("未配置 allowlist 时应该允许访问: %v", err)
		}
	})

	t.Run("配置allowlist后，只允许列表中的用户", func(t *testing.T) {
		// 配置 allowlist
		policy := permission.ResourcePolicy{
			Create:    false,
			Read:      true,
			Update:    false,
			Delete:    false,
			AllowList: []string{"allowed_user", "another_user"},
		}
		adapter.permission.SetResourcePolicy(permission.ResourceUser, policy)

		// 允许列表中的用户应该可以访问
		_, err := adapter.GetUser(context.Background(), GetUserInput{
			UserID: "allowed_user",
		})
		if err != nil {
			t.Errorf("allowlist 中的用户应该允许访问: %v", err)
		}

		// 不在允许列表中的用户应该被拒绝
		_, err = adapter.GetUser(context.Background(), GetUserInput{
			UserID: "unauthorized_user",
		})
		if err == nil {
			t.Error("不在 allowlist 中的用户应该被拒绝访问")
		}
		if err != nil && !strings.Contains(err.Error(), "不在允许列表中") {
			t.Errorf("错误消息应该包含 '不在允许列表中': %v", err)
		}
	})

	t.Run("Update操作也受allowlist限制", func(t *testing.T) {
		// 配置 allowlist，并启用 update 权限
		policy := permission.ResourcePolicy{
			Create:    false,
			Read:      true,
			Update:    true,
			Delete:    false,
			AllowList: []string{"allowed_user"},
		}
		adapter.permission.SetResourcePolicy(permission.ResourceUser, policy)

		// 允许列表中的用户应该可以更新
		_, err := adapter.UpdateUser(context.Background(), UpdateUserInput{
			UserID: "allowed_user",
			Name:   "New Name",
		})
		if err != nil {
			t.Errorf("allowlist 中的用户应该允许更新: %v", err)
		}

		// 不在允许列表中的用户应该被拒绝
		_, err = adapter.UpdateUser(context.Background(), UpdateUserInput{
			UserID: "unauthorized_user",
			Name:   "New Name",
		})
		if err == nil {
			t.Error("不在 allowlist 中的用户应该被拒绝更新")
		}
	})

	t.Run("Delete操作也受allowlist限制", func(t *testing.T) {
		// 配置 allowlist，并启用 delete 权限
		policy := permission.ResourcePolicy{
			Create:    false,
			Read:      true,
			Update:    false,
			Delete:    true,
			AllowList: []string{"deletable_user"},
		}
		adapter.permission.SetResourcePolicy(permission.ResourceUser, policy)

		// 允许列表中的用户应该可以删除
		_, err := adapter.DeleteUser(context.Background(), DeleteUserInput{
			UserID: "deletable_user",
		})
		if err != nil {
			t.Errorf("allowlist 中的用户应该允许删除: %v", err)
		}

		// 不在允许列表中的用户应该被拒绝
		_, err = adapter.DeleteUser(context.Background(), DeleteUserInput{
			UserID: "unauthorized_user",
		})
		if err == nil {
			t.Error("不在 allowlist 中的用户应该被拒绝删除")
		}
	})

	t.Run("操作权限被禁用时，allowlist不生效", func(t *testing.T) {
		// 配置 allowlist，但禁用 read 权限
		policy := permission.ResourcePolicy{
			Create:    false,
			Read:      false, // 禁用 read
			Update:    false,
			Delete:    false,
			AllowList: []string{"allowed_user"},
		}
		adapter.permission.SetResourcePolicy(permission.ResourceUser, policy)

		// 即使在 allowlist 中，也应该被拒绝（因为 read 权限被禁用）
		_, err := adapter.GetUser(context.Background(), GetUserInput{
			UserID: "allowed_user",
		})
		if err == nil {
			t.Error("read 权限被禁用时，即使在 allowlist 中也应该被拒绝")
		}
		if err != nil && !strings.Contains(err.Error(), "不允许对资源") {
			t.Errorf("错误消息应该说明操作被禁止: %v", err)
		}
	})
}
