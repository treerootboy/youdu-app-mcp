package adapter

import (
	"context"
	"strings"
	"testing"

	"github.com/yourusername/youdu-app-mcp/internal/adapter/testdata"
	"github.com/yourusername/youdu-app-mcp/internal/config"
	"github.com/yourusername/youdu-app-mcp/internal/permission"
)

// TestAdapter_MessageSendPermission 测试消息发送权限功能
func TestAdapter_MessageSendPermission(t *testing.T) {
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

	t.Run("未配置allowsend限制时，允许发送给任何用户", func(t *testing.T) {
		// 配置权限：允许发送消息，但不限制接收者
		policy := permission.ResourcePolicy{
			Create: true,
		}
		adapter.permission.SetResourcePolicy(permission.ResourceMessage, policy)

		// 应该允许发送给任何用户
		_, err := adapter.SendTextMessage(context.Background(), SendTextMessageInput{
			ToUser:  "any_user",
			Content: "Test message",
		})
		if err != nil {
			t.Errorf("未配置allowsend限制时应该允许发送: %v", err)
		}
	})

	t.Run("配置了allowsend.users后，只允许发送给列表中的用户", func(t *testing.T) {
		// 配置权限：只允许发送给指定用户
		policy := permission.ResourcePolicy{
			Create: true,
			AllowSend: permission.AllowSend{
				Users: []string{"10232", "8891"},
			},
		}
		adapter.permission.SetResourcePolicy(permission.ResourceMessage, policy)

		// 允许发送给列表中的用户
		_, err := adapter.SendTextMessage(context.Background(), SendTextMessageInput{
			ToUser:  "10232",
			Content: "Test message",
		})
		if err != nil {
			t.Errorf("应该允许发送给列表中的用户: %v", err)
		}

		// 不允许发送给不在列表中的用户
		_, err = adapter.SendTextMessage(context.Background(), SendTextMessageInput{
			ToUser:  "99999",
			Content: "Test message",
		})
		if err == nil {
			t.Error("不应该允许发送给不在列表中的用户")
		}
		if !strings.Contains(err.Error(), "不允许向用户 '99999' 发送消息") {
			t.Errorf("错误消息应该说明不允许向该用户发送: %v", err)
		}
	})

	t.Run("配置了allowsend.dept后，只允许发送给列表中的部门", func(t *testing.T) {
		// 配置权限：只允许发送给指定部门
		policy := permission.ResourcePolicy{
			Create: true,
			AllowSend: permission.AllowSend{
				Dept: []string{"1"},
			},
		}
		adapter.permission.SetResourcePolicy(permission.ResourceMessage, policy)

		// 允许发送给列表中的部门
		_, err := adapter.SendTextMessage(context.Background(), SendTextMessageInput{
			ToDept:  "1",
			Content: "Test message",
		})
		if err != nil {
			t.Errorf("应该允许发送给列表中的部门: %v", err)
		}

		// 不允许发送给不在列表中的部门
		_, err = adapter.SendTextMessage(context.Background(), SendTextMessageInput{
			ToDept:  "999",
			Content: "Test message",
		})
		if err == nil {
			t.Error("不应该允许发送给不在列表中的部门")
		}
		if !strings.Contains(err.Error(), "不允许向部门 '999' 发送消息") {
			t.Errorf("错误消息应该说明不允许向该部门发送: %v", err)
		}
	})

	t.Run("支持同时发送给多个用户（用|分隔）", func(t *testing.T) {
		// 配置权限：只允许发送给指定用户
		policy := permission.ResourcePolicy{
			Create: true,
			AllowSend: permission.AllowSend{
				Users: []string{"10232", "8891"},
			},
		}
		adapter.permission.SetResourcePolicy(permission.ResourceMessage, policy)

		// 所有用户都在列表中，应该允许
		_, err := adapter.SendTextMessage(context.Background(), SendTextMessageInput{
			ToUser:  "10232|8891",
			Content: "Test message",
		})
		if err != nil {
			t.Errorf("所有用户都在列表中时应该允许: %v", err)
		}

		// 部分用户不在列表中，应该拒绝
		_, err = adapter.SendTextMessage(context.Background(), SendTextMessageInput{
			ToUser:  "10232|99999",
			Content: "Test message",
		})
		if err == nil {
			t.Error("部分用户不在列表中时应该拒绝")
		}
	})

	t.Run("支持同时发送给用户和部门", func(t *testing.T) {
		// 配置权限：限制用户和部门
		policy := permission.ResourcePolicy{
			Create: true,
			AllowSend: permission.AllowSend{
				Users: []string{"10232"},
				Dept:  []string{"1"},
			},
		}
		adapter.permission.SetResourcePolicy(permission.ResourceMessage, policy)

		// 用户和部门都在列表中
		_, err := adapter.SendTextMessage(context.Background(), SendTextMessageInput{
			ToUser:  "10232",
			ToDept:  "1",
			Content: "Test message",
		})
		if err != nil {
			t.Errorf("用户和部门都在列表中时应该允许: %v", err)
		}

		// 用户不在列表中
		_, err = adapter.SendTextMessage(context.Background(), SendTextMessageInput{
			ToUser:  "99999",
			ToDept:  "1",
			Content: "Test message",
		})
		if err == nil {
			t.Error("用户不在列表中时应该拒绝")
		}

		// 部门不在列表中
		_, err = adapter.SendTextMessage(context.Background(), SendTextMessageInput{
			ToUser:  "10232",
			ToDept:  "999",
			Content: "Test message",
		})
		if err == nil {
			t.Error("部门不在列表中时应该拒绝")
		}
	})

	t.Run("其他消息类型也受限制", func(t *testing.T) {
		// 配置权限：只允许发送给指定用户
		policy := permission.ResourcePolicy{
			Create: true,
			AllowSend: permission.AllowSend{
				Users: []string{"10232"},
			},
		}
		adapter.permission.SetResourcePolicy(permission.ResourceMessage, policy)

		// 图片消息
		_, err := adapter.SendImageMessage(context.Background(), SendImageMessageInput{
			ToUser:  "99999",
			MediaID: "test_media_id",
		})
		if err == nil {
			t.Error("图片消息也应该受权限限制")
		}

		// 文件消息
		_, err = adapter.SendFileMessage(context.Background(), SendFileMessageInput{
			ToUser:  "99999",
			MediaID: "test_media_id",
		})
		if err == nil {
			t.Error("文件消息也应该受权限限制")
		}

		// 链接消息
		_, err = adapter.SendLinkMessage(context.Background(), SendLinkMessageInput{
			ToUser: "99999",
			Title:  "Test",
			URL:    "http://example.com",
		})
		if err == nil {
			t.Error("链接消息也应该受权限限制")
		}

		// 系统消息
		_, err = adapter.SendSysMessage(context.Background(), SendSysMessageInput{
			ToUser:  "99999",
			Title:   "Test",
			Content: "Test message",
		})
		if err == nil {
			t.Error("系统消息也应该受权限限制")
		}
	})

	t.Run("create权限被禁用时，不允许发送消息", func(t *testing.T) {
		// 配置权限：禁用消息发送
		policy := permission.ResourcePolicy{
			Create: false,
		}
		adapter.permission.SetResourcePolicy(permission.ResourceMessage, policy)

		_, err := adapter.SendTextMessage(context.Background(), SendTextMessageInput{
			ToUser:  "any_user",
			Content: "Test message",
		})
		if err == nil {
			t.Error("create权限被禁用时不应该允许发送消息")
		}
		if !strings.Contains(err.Error(), "不允许发送消息") {
			t.Errorf("错误消息应该说明不允许发送消息: %v", err)
		}
	})
}
