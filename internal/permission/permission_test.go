package permission

import (
	"testing"
)

func TestPermission_CheckWithID_AllowList(t *testing.T) {
	tests := []struct {
		name         string
		enabled      bool
		allowAll     bool
		policy       ResourcePolicy
		resource     Resource
		action       Action
		resourceID   string
		expectError  bool
		errorContain string
	}{
		{
			name:     "允许列表为空时，任何ID都可以访问",
			enabled:  true,
			allowAll: false,
			policy: ResourcePolicy{
				Read:      true,
				AllowList: []string{},
			},
			resource:    ResourceUser,
			action:      ActionRead,
			resourceID:  "user123",
			expectError: false,
		},
		{
			name:     "ID在允许列表中，应该允许访问",
			enabled:  true,
			allowAll: false,
			policy: ResourcePolicy{
				Read:      true,
				AllowList: []string{"10232", "10023"},
			},
			resource:    ResourceUser,
			action:      ActionRead,
			resourceID:  "10232",
			expectError: false,
		},
		{
			name:     "ID不在允许列表中，应该拒绝访问",
			enabled:  true,
			allowAll: false,
			policy: ResourcePolicy{
				Read:      true,
				AllowList: []string{"10232", "10023"},
			},
			resource:     ResourceUser,
			action:       ActionRead,
			resourceID:   "99999",
			expectError:  true,
			errorContain: "不在允许列表中",
		},
		{
			name:     "未提供resourceID时，只检查操作权限",
			enabled:  true,
			allowAll: false,
			policy: ResourcePolicy{
				Read:      true,
				AllowList: []string{"10232", "10023"},
			},
			resource:    ResourceUser,
			action:      ActionRead,
			resourceID:  "",
			expectError: false,
		},
		{
			name:     "操作权限被拒绝，即使ID在允许列表中",
			enabled:  true,
			allowAll: false,
			policy: ResourcePolicy{
				Read:      false,
				AllowList: []string{"10232", "10023"},
			},
			resource:     ResourceUser,
			action:       ActionRead,
			resourceID:   "10232",
			expectError:  true,
			errorContain: "不允许对资源",
		},
		{
			name:     "权限系统禁用时，应该允许所有访问",
			enabled:  false,
			allowAll: true,
			policy: ResourcePolicy{
				Read:      false,
				AllowList: []string{"10232"},
			},
			resource:    ResourceUser,
			action:      ActionRead,
			resourceID:  "99999",
			expectError: false,
		},
		{
			name:     "allow_all=true时，应该允许所有访问",
			enabled:  true,
			allowAll: true,
			policy: ResourcePolicy{
				Read:      false,
				AllowList: []string{"10232"},
			},
			resource:    ResourceUser,
			action:      ActionRead,
			resourceID:  "99999",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建权限实例
			resources := map[Resource]ResourcePolicy{
				tt.resource: tt.policy,
			}
			p := New(tt.enabled, tt.allowAll, resources)

			// 检查权限
			err := p.CheckWithID(tt.resource, tt.action, tt.resourceID)

			// 验证结果
			if tt.expectError {
				if err == nil {
					t.Errorf("期望错误，但没有返回错误")
					return
				}
				if tt.errorContain != "" {
					if !contains(err.Error(), tt.errorContain) {
						t.Errorf("错误消息不包含 '%s': %v", tt.errorContain, err)
					}
				}
			} else {
				if err != nil {
					t.Errorf("不期望错误，但返回了错误: %v", err)
				}
			}
		})
	}
}

func TestPermission_Check_BackwardCompatibility(t *testing.T) {
	// 测试向后兼容性：确保原有的 Check 方法仍然工作
	policy := ResourcePolicy{
		Read:      true,
		AllowList: []string{"10232", "10023"},
	}
	resources := map[Resource]ResourcePolicy{
		ResourceUser: policy,
	}
	p := New(true, false, resources)

	// 调用原有的 Check 方法（不传递 resourceID）
	err := p.Check(ResourceUser, ActionRead)
	if err != nil {
		t.Errorf("向后兼容性测试失败，Check 方法返回错误: %v", err)
	}
}

// contains 检查字符串是否包含子串
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 || 
		(len(s) > 0 && len(substr) > 0 && hasSubstring(s, substr)))
}

func hasSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
