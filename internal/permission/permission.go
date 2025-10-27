package permission

import (
	"fmt"
	"sync"
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

// Permission 权限配置（纯数据结构 + 业务逻辑）
type Permission struct {
	Enabled   bool                        // 是否启用权限检查
	AllowAll  bool                        // 是否允许所有操作（调试用）
	Resources map[Resource]ResourcePolicy // 资源权限策略
	mu        sync.RWMutex                // 保护并发访问
}

// ResourcePolicy 资源权限策略
type ResourcePolicy struct {
	Create    bool     `mapstructure:"create"`    // 允许创建
	Read      bool     `mapstructure:"read"`      // 允许读取
	Update    bool     `mapstructure:"update"`    // 允许更新
	Delete    bool     `mapstructure:"delete"`    // 允许删除
	AllowList []string `mapstructure:"allowlist"` // 允许访问的资源ID列表（行级权限）
}

// New 创建新的 Permission 实例（构造函数）
// 配置加载逻辑已移到 config 包，此函数仅用于创建实例
func New(enabled, allowAll bool, resources map[Resource]ResourcePolicy) *Permission {
	// 如果未启用权限检查，自动设置为允许所有
	if !enabled {
		allowAll = true
	}

	return &Permission{
		Enabled:   enabled,
		AllowAll:  allowAll,
		Resources: resources,
	}
}

// Check 检查权限（业务逻辑）
func (p *Permission) Check(resource Resource, action Action) error {
	return p.CheckWithID(resource, action, "")
}

// CheckWithID 检查权限（包含行级权限）
func (p *Permission) CheckWithID(resource Resource, action Action, resourceID string) error {
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

	// 检查行级权限（如果配置了 allowlist 且提供了 resourceID）
	if len(policy.AllowList) > 0 && resourceID != "" {
		// 检查 resourceID 是否在 allowlist 中
		found := false
		for _, allowedID := range policy.AllowList {
			if allowedID == resourceID {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("权限拒绝：资源 ID '%s' 不在允许列表中", resourceID)
		}
	}

	return nil
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
