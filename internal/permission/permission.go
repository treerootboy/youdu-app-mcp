package permission

import (
	"fmt"
	"sync"

	"github.com/spf13/viper"
)

// Action 操作类型
type Action string

const (
	ActionCreate Action = "create" // 创建
	ActionRead   Action = "read"   // 读取
	ActionUpdate Action = "update" // 更新
	ActionDelete Action = "delete" // 删除
)

// Resource 资源类型
type Resource string

const (
	ResourceDept    Resource = "dept"    // 部门
	ResourceUser    Resource = "user"    // 用户
	ResourceGroup   Resource = "group"   // 群组
	ResourceSession Resource = "session" // 会话
	ResourceMessage Resource = "message" // 消息
)

// Permission 权限配置
type Permission struct {
	Enabled   bool                        `mapstructure:"enabled"`   // 是否启用权限检查
	AllowAll  bool                        `mapstructure:"allow_all"` // 是否允许所有操作（调试用）
	Resources map[Resource]ResourcePolicy `mapstructure:"resources"` // 资源权限策略
	mu        sync.RWMutex
}

// ResourcePolicy 资源权限策略
type ResourcePolicy struct {
	Create bool `mapstructure:"create"` // 允许创建
	Read   bool `mapstructure:"read"`   // 允许读取
	Update bool `mapstructure:"update"` // 允许更新
	Delete bool `mapstructure:"delete"` // 允许删除
}

var (
	globalPermission *Permission
	once             sync.Once
)

// Load 加载权限配置
func Load() (*Permission, error) {
	var perm Permission

	// 从配置文件读取权限设置
	if err := viper.UnmarshalKey("permission", &perm); err != nil {
		// 如果没有配置，使用默认配置（全部允许）
		perm = Permission{
			Enabled:  false,
			AllowAll: true,
			Resources: map[Resource]ResourcePolicy{
				ResourceDept:    {Create: true, Read: true, Update: true, Delete: true},
				ResourceUser:    {Create: true, Read: true, Update: true, Delete: true},
				ResourceGroup:   {Create: true, Read: true, Update: true, Delete: true},
				ResourceSession: {Create: true, Read: true, Update: true, Delete: true},
				ResourceMessage: {Create: true, Read: true, Update: true, Delete: true},
			},
		}
	}

	// 如果未启用权限检查，设置为全部允许
	if !perm.Enabled {
		perm.AllowAll = true
	}

	return &perm, nil
}

// GetGlobal 获取全局权限配置
func GetGlobal() *Permission {
	once.Do(func() {
		perm, err := Load()
		if err != nil {
			// 如果加载失败，使用默认配置
			perm = &Permission{
				Enabled:  false,
				AllowAll: true,
			}
		}
		globalPermission = perm
	})
	return globalPermission
}

// Check 检查权限
func (p *Permission) Check(resource Resource, action Action) error {
	p.mu.RLock()
	defer p.mu.RUnlock()

	// 如果未启用权限检查或允许所有操作
	if !p.Enabled || p.AllowAll {
		return nil
	}

	// 检查资源权限
	policy, exists := p.Resources[resource]
	if !exists {
		return fmt.Errorf("权限拒绝：未配置资源 '%s' 的权限策略", resource)
	}

	// 检查具体操作权限
	var allowed bool
	switch action {
	case ActionCreate:
		allowed = policy.Create
	case ActionRead:
		allowed = policy.Read
	case ActionUpdate:
		allowed = policy.Update
	case ActionDelete:
		allowed = policy.Delete
	default:
		return fmt.Errorf("权限拒绝：未知的操作类型 '%s'", action)
	}

	if !allowed {
		return fmt.Errorf("权限拒绝：不允许对资源 '%s' 执行 '%s' 操作", resource, action)
	}

	return nil
}

// CheckGlobal 使用全局权限配置检查权限
func CheckGlobal(resource Resource, action Action) error {
	return GetGlobal().Check(resource, action)
}

// SetResourcePolicy 设置资源权限策略
func (p *Permission) SetResourcePolicy(resource Resource, policy ResourcePolicy) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.Resources == nil {
		p.Resources = make(map[Resource]ResourcePolicy)
	}
	p.Resources[resource] = policy
}

// GetResourcePolicy 获取资源权限策略
func (p *Permission) GetResourcePolicy(resource Resource) (ResourcePolicy, bool) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	policy, exists := p.Resources[resource]
	return policy, exists
}

// Enable 启用权限检查
func (p *Permission) Enable() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.Enabled = true
	p.AllowAll = false
}

// Disable 禁用权限检查
func (p *Permission) Disable() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.Enabled = false
	p.AllowAll = true
}

// IsEnabled 检查权限系统是否启用
func (p *Permission) IsEnabled() bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.Enabled
}
