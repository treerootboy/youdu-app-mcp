package permission

import (
	"strings"
	"testing"
)

// TestPermission_CheckMessageSend 测试消息发送权限功能
func TestPermission_CheckMessageSend(t *testing.T) {
	t.Run("权限系统禁用时，允许所有发送", func(t *testing.T) {
		perm := New(false, true, nil)
		err := perm.CheckMessageSend("user1", "dept1")
		if err != nil {
			t.Errorf("权限系统禁用时应该允许所有发送: %v", err)
		}
	})

	t.Run("allow_all=true时，允许所有发送", func(t *testing.T) {
		perm := New(true, true, nil)
		err := perm.CheckMessageSend("user1", "dept1")
		if err != nil {
			t.Errorf("allow_all=true时应该允许所有发送: %v", err)
		}
	})

	t.Run("未配置message资源时，拒绝发送", func(t *testing.T) {
		perm := New(true, false, map[Resource]ResourcePolicy{})
		err := perm.CheckMessageSend("user1", "dept1")
		if err == nil {
			t.Error("未配置message资源时应该拒绝发送")
		}
		if !strings.Contains(err.Error(), "未配置资源 'message'") {
			t.Errorf("错误消息应该说明未配置资源: %v", err)
		}
	})

	t.Run("create权限被禁用时，拒绝发送", func(t *testing.T) {
		perm := New(true, false, map[Resource]ResourcePolicy{
			ResourceMessage: {
				Create: false,
			},
		})
		err := perm.CheckMessageSend("user1", "dept1")
		if err == nil {
			t.Error("create权限被禁用时应该拒绝发送")
		}
		if !strings.Contains(err.Error(), "不允许发送消息") {
			t.Errorf("错误消息应该说明不允许发送: %v", err)
		}
	})

	t.Run("未配置allowsend限制时，允许发送给任何用户", func(t *testing.T) {
		perm := New(true, false, map[Resource]ResourcePolicy{
			ResourceMessage: {
				Create: true,
			},
		})
		err := perm.CheckMessageSend("any_user", "any_dept")
		if err != nil {
			t.Errorf("未配置allowsend限制时应该允许发送给任何用户: %v", err)
		}
	})

	t.Run("配置了allowsend.users后，只允许发送给列表中的用户", func(t *testing.T) {
		perm := New(true, false, map[Resource]ResourcePolicy{
			ResourceMessage: {
				Create: true,
				AllowSend: AllowSend{
					Users: []string{"10232", "8891"},
				},
			},
		})

		// 允许列表中的用户
		err := perm.CheckMessageSend("10232", "")
		if err != nil {
			t.Errorf("应该允许发送给列表中的用户: %v", err)
		}

		// 不在允许列表中的用户
		err = perm.CheckMessageSend("99999", "")
		if err == nil {
			t.Error("不应该允许发送给不在列表中的用户")
		}
		if !strings.Contains(err.Error(), "不允许向用户 '99999' 发送消息") {
			t.Errorf("错误消息应该说明不允许向该用户发送: %v", err)
		}
	})

	t.Run("配置了allowsend.dept后，只允许发送给列表中的部门", func(t *testing.T) {
		perm := New(true, false, map[Resource]ResourcePolicy{
			ResourceMessage: {
				Create: true,
				AllowSend: AllowSend{
					Dept: []string{"1", "2"},
				},
			},
		})

		// 允许列表中的部门
		err := perm.CheckMessageSend("", "1")
		if err != nil {
			t.Errorf("应该允许发送给列表中的部门: %v", err)
		}

		// 不在允许列表中的部门
		err = perm.CheckMessageSend("", "999")
		if err == nil {
			t.Error("不应该允许发送给不在列表中的部门")
		}
		if !strings.Contains(err.Error(), "不允许向部门 '999' 发送消息") {
			t.Errorf("错误消息应该说明不允许向该部门发送: %v", err)
		}
	})

	t.Run("支持多个用户ID（用|分隔）", func(t *testing.T) {
		perm := New(true, false, map[Resource]ResourcePolicy{
			ResourceMessage: {
				Create: true,
				AllowSend: AllowSend{
					Users: []string{"10232", "8891"},
				},
			},
		})

		// 所有用户都在允许列表中
		err := perm.CheckMessageSend("10232|8891", "")
		if err != nil {
			t.Errorf("所有用户都在列表中时应该允许: %v", err)
		}

		// 部分用户不在允许列表中
		err = perm.CheckMessageSend("10232|99999", "")
		if err == nil {
			t.Error("部分用户不在列表中时应该拒绝")
		}
		if !strings.Contains(err.Error(), "不允许向用户 '99999' 发送消息") {
			t.Errorf("错误消息应该说明哪个用户不允许: %v", err)
		}
	})

	t.Run("支持多个部门ID（用|分隔）", func(t *testing.T) {
		perm := New(true, false, map[Resource]ResourcePolicy{
			ResourceMessage: {
				Create: true,
				AllowSend: AllowSend{
					Dept: []string{"1", "2"},
				},
			},
		})

		// 所有部门都在允许列表中
		err := perm.CheckMessageSend("", "1|2")
		if err != nil {
			t.Errorf("所有部门都在列表中时应该允许: %v", err)
		}

		// 部分部门不在允许列表中
		err = perm.CheckMessageSend("", "1|999")
		if err == nil {
			t.Error("部分部门不在列表中时应该拒绝")
		}
		if !strings.Contains(err.Error(), "不允许向部门 '999' 发送消息") {
			t.Errorf("错误消息应该说明哪个部门不允许: %v", err)
		}
	})

	t.Run("同时配置users和dept限制", func(t *testing.T) {
		perm := New(true, false, map[Resource]ResourcePolicy{
			ResourceMessage: {
				Create: true,
				AllowSend: AllowSend{
					Users: []string{"10232", "8891"},
					Dept:  []string{"1"},
				},
			},
		})

		// 用户在允许列表中
		err := perm.CheckMessageSend("10232", "")
		if err != nil {
			t.Errorf("用户在列表中时应该允许: %v", err)
		}

		// 部门在允许列表中
		err = perm.CheckMessageSend("", "1")
		if err != nil {
			t.Errorf("部门在列表中时应该允许: %v", err)
		}

		// 同时发送给用户和部门
		err = perm.CheckMessageSend("10232", "1")
		if err != nil {
			t.Errorf("用户和部门都在列表中时应该允许: %v", err)
		}

		// 用户不在列表中
		err = perm.CheckMessageSend("99999", "1")
		if err == nil {
			t.Error("用户不在列表中时应该拒绝")
		}

		// 部门不在列表中
		err = perm.CheckMessageSend("10232", "999")
		if err == nil {
			t.Error("部门不在列表中时应该拒绝")
		}
	})

	t.Run("处理带空格的ID字符串", func(t *testing.T) {
		perm := New(true, false, map[Resource]ResourcePolicy{
			ResourceMessage: {
				Create: true,
				AllowSend: AllowSend{
					Users: []string{"10232", "8891"},
				},
			},
		})

		// 带空格的ID
		err := perm.CheckMessageSend("10232 | 8891", "")
		if err != nil {
			t.Errorf("应该正确处理带空格的ID: %v", err)
		}
	})
}

// TestSplitIDs 测试ID分割函数
func TestSplitIDs(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{"单个ID", "123", []string{"123"}},
		{"多个ID", "123|456|789", []string{"123", "456", "789"}},
		{"带空格的ID", "123 | 456 | 789", []string{"123", "456", "789"}},
		{"空字符串", "", nil},
		{"只有分隔符", "|", []string{}},
		{"开头和结尾的空格", " 123 | 456 ", []string{"123", "456"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := splitIDs(tt.input)
			if len(result) != len(tt.expected) {
				t.Errorf("期望 %d 个元素，得到 %d 个", len(tt.expected), len(result))
				return
			}
			for i, v := range result {
				if v != tt.expected[i] {
					t.Errorf("索引 %d: 期望 '%s'，得到 '%s'", i, tt.expected[i], v)
				}
			}
		})
	}
}
