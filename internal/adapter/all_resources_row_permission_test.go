package adapter

import (
	"context"
	"strings"
	"testing"

	"github.com/yourusername/youdu-app-mcp/internal/adapter/testdata"
	"github.com/yourusername/youdu-app-mcp/internal/config"
	"github.com/yourusername/youdu-app-mcp/internal/permission"
)

// TestAdapter_AllResourcesRowLevelPermissions 测试所有资源的行级权限功能
func TestAdapter_AllResourcesRowLevelPermissions(t *testing.T) {
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

	t.Run("部门资源行级权限拒绝测试", func(t *testing.T) {
		// 配置部门 allowlist - 只允许特定ID
		policy := permission.ResourcePolicy{
			Create:    false,
			Read:      true,
			Update:    true,
			Delete:    true,
			AllowList: []string{"1", "2", "100"},
		}
		adapter.permission.SetResourcePolicy(permission.ResourceDept, policy)

		// 测试不允许的部门ID - GetDeptList
		_, err := adapter.GetDeptList(context.Background(), DeptListInput{
			DeptID: 999,
		})
		if err == nil {
			t.Error("不在 allowlist 中的部门应该被拒绝")
		}
		if err != nil && !strings.Contains(err.Error(), "不在允许列表中") {
			t.Errorf("错误消息应该包含 '不在允许列表中': %v", err)
		}

		// 测试不允许的部门ID - UpdateDept
		_, err = adapter.UpdateDept(context.Background(), UpdateDeptInput{
			DeptID: 888,
			Name:   "New Name",
		})
		if err == nil {
			t.Error("不在 allowlist 中的部门应该无法更新")
		}

		// 测试不允许的部门ID - DeleteDept
		_, err = adapter.DeleteDept(context.Background(), DeleteDeptInput{
			DeptID: 777,
		})
		if err == nil {
			t.Error("不在 allowlist 中的部门应该无法删除")
		}
	})

	t.Run("群组资源行级权限拒绝测试", func(t *testing.T) {
		// 配置群组 allowlist
		policy := permission.ResourcePolicy{
			Create:    false,
			Read:      true,
			Update:    true,
			Delete:    true,
			AllowList: []string{"group1", "group2"},
		}
		adapter.permission.SetResourcePolicy(permission.ResourceGroup, policy)

		// 测试 GetGroupInfo - 拒绝未授权ID
		_, err := adapter.GetGroupInfo(context.Background(), GetGroupInfoInput{
			GroupID: "group999",
		})
		if err == nil {
			t.Error("不在 allowlist 中的群组应该被拒绝")
		}

		// 测试 UpdateGroup - 拒绝未授权ID
		_, err = adapter.UpdateGroup(context.Background(), UpdateGroupInput{
			GroupID: "group888",
			Name:    "New Group Name",
		})
		if err == nil {
			t.Error("不在 allowlist 中的群组应该无法更新")
		}

		// 测试 DeleteGroup - 拒绝未授权ID
		_, err = adapter.DeleteGroup(context.Background(), DeleteGroupInput{
			GroupID: "group777",
		})
		if err == nil {
			t.Error("不在 allowlist 中的群组应该无法删除")
		}

		// 测试 AddGroupMember - 拒绝未授权ID
		_, err = adapter.AddGroupMember(context.Background(), AddGroupMemberInput{
			GroupID: "group666",
			Members: []string{"user1"},
		})
		if err == nil {
			t.Error("不在 allowlist 中的群组应该无法添加成员")
		}

		// 测试 DelGroupMember - 拒绝未授权ID
		_, err = adapter.DelGroupMember(context.Background(), DelGroupMemberInput{
			GroupID: "group555",
			Members: []string{"user1"},
		})
		if err == nil {
			t.Error("不在 allowlist 中的群组应该无法删除成员")
		}
	})

	t.Run("会话资源行级权限拒绝测试", func(t *testing.T) {
		// 配置会话 allowlist
		policy := permission.ResourcePolicy{
			Create:    false,
			Read:      true,
			Update:    true,
			Delete:    false,
			AllowList: []string{"session1", "session2"},
		}
		adapter.permission.SetResourcePolicy(permission.ResourceSession, policy)

		// 测试 GetSession - 拒绝未授权ID
		_, err := adapter.GetSession(context.Background(), GetSessionInput{
			SessionID: "session999",
		})
		if err == nil {
			t.Error("不在 allowlist 中的会话应该被拒绝")
		}

		// 测试 UpdateSession - 拒绝未授权ID
		_, err = adapter.UpdateSession(context.Background(), UpdateSessionInput{
			SessionID: "session888",
			Title:     "New Title",
			OpUser:    "user1",
		})
		if err == nil {
			t.Error("不在 allowlist 中的会话应该无法更新")
		}

		// 测试 SendTextSessionMessage - 拒绝未授权ID
		_, err = adapter.SendTextSessionMessage(context.Background(), SendTextSessionMessageInput{
			SessionID: "session777",
			Content:   "Test message",
			Sender:    "user1",
		})
		if err == nil {
			t.Error("不在 allowlist 中的会话应该无法发送消息")
		}

		// 测试 SendImageSessionMessage - 拒绝未授权ID
		_, err = adapter.SendImageSessionMessage(context.Background(), SendImageSessionMessageInput{
			SessionID: "session666",
			MediaID:   "media123",
			Sender:    "user1",
		})
		if err == nil {
			t.Error("不在 allowlist 中的会话应该无法发送图片消息")
		}

		// 测试 SendFileSessionMessage - 拒绝未授权ID
		_, err = adapter.SendFileSessionMessage(context.Background(), SendFileSessionMessageInput{
			SessionID: "session555",
			MediaID:   "media123",
			Sender:    "user1",
		})
		if err == nil {
			t.Error("不在 allowlist 中的会话应该无法发送文件消息")
		}
	})

	t.Run("多种资源类型同时配置allowlist", func(t *testing.T) {
		// 同时为多种资源配置 allowlist
		adapter.permission.SetResourcePolicy(permission.ResourceUser, permission.ResourcePolicy{
			Read:      true,
			AllowList: []string{"user1", "user2"},
		})
		adapter.permission.SetResourcePolicy(permission.ResourceDept, permission.ResourcePolicy{
			Read:      true,
			AllowList: []string{"10", "20"},
		})
		adapter.permission.SetResourcePolicy(permission.ResourceGroup, permission.ResourcePolicy{
			Read:      true,
			AllowList: []string{"g1", "g2"},
		})
		adapter.permission.SetResourcePolicy(permission.ResourceSession, permission.ResourcePolicy{
			Read:      true,
			AllowList: []string{"s1", "s2"},
		})

		// 验证每个资源的 allowlist 独立工作 - 测试拒绝情况
		_, err := adapter.GetUser(context.Background(), GetUserInput{UserID: "user999"})
		if err == nil {
			t.Error("用户 allowlist 应该拒绝未授权的ID")
		}

		_, err = adapter.GetDeptList(context.Background(), DeptListInput{DeptID: 999})
		if err == nil {
			t.Error("部门 allowlist 应该拒绝未授权的ID")
		}

		_, err = adapter.GetGroupInfo(context.Background(), GetGroupInfoInput{GroupID: "g999"})
		if err == nil {
			t.Error("群组 allowlist 应该拒绝未授权的ID")
		}

		_, err = adapter.GetSession(context.Background(), GetSessionInput{SessionID: "s999"})
		if err == nil {
			t.Error("会话 allowlist 应该拒绝未授权的ID")
		}
	})

	t.Run("未配置allowlist时允许所有ID", func(t *testing.T) {
		// 为所有资源清除 allowlist
		adapter.permission.SetResourcePolicy(permission.ResourceDept, permission.ResourcePolicy{
			Read:      true,
			AllowList: []string{}, // 空 allowlist
		})
		adapter.permission.SetResourcePolicy(permission.ResourceGroup, permission.ResourcePolicy{
			Read:      true,
			AllowList: []string{}, // 空 allowlist
		})
		adapter.permission.SetResourcePolicy(permission.ResourceSession, permission.ResourcePolicy{
			Read:      true,
			AllowList: []string{}, // 空 allowlist
		})

		// 现在任何ID都应该被允许（只要有操作权限）
		// 这些调用可能会因为mock server的限制失败，但不应该因为权限被拒绝
		_, err := adapter.GetDeptList(context.Background(), DeptListInput{DeptID: 999})
		if err != nil && strings.Contains(err.Error(), "不在允许列表中") {
			t.Error("未配置 allowlist 时不应该因为ID被拒绝")
		}

		_, err = adapter.GetGroupInfo(context.Background(), GetGroupInfoInput{GroupID: "any_group"})
		if err != nil && strings.Contains(err.Error(), "不在允许列表中") {
			t.Error("未配置 allowlist 时不应该因为ID被拒绝")
		}

		_, err = adapter.GetSession(context.Background(), GetSessionInput{SessionID: "any_session"})
		if err != nil && strings.Contains(err.Error(), "不在允许列表中") {
			t.Error("未配置 allowlist 时不应该因为ID被拒绝")
		}
	})
}
