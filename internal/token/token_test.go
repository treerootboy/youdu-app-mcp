package token

import (
	"testing"
	"time"
)

func TestManager_Generate(t *testing.T) {
	m := NewManager()

	// 生成不过期的 token
	token, err := m.Generate("test token", nil)
	if err != nil {
		t.Fatalf("生成 token 失败: %v", err)
	}

	if token.Value == "" {
		t.Error("token value 为空")
	}

	if token.ID == "" {
		t.Error("token ID 为空")
	}

	if token.Description != "test token" {
		t.Errorf("期望描述为 'test token'，得到 '%s'", token.Description)
	}

	if token.ExpiresAt != nil {
		t.Error("期望 ExpiresAt 为 nil")
	}

	// 生成有过期时间的 token
	duration := 1 * time.Hour
	token2, err := m.Generate("expiring token", &duration)
	if err != nil {
		t.Fatalf("生成 token 失败: %v", err)
	}

	if token2.ExpiresAt == nil {
		t.Error("期望 ExpiresAt 不为 nil")
	}

	// 验证 token 数量
	if m.Count() != 2 {
		t.Errorf("期望 token 数量为 2，得到 %d", m.Count())
	}
}

func TestManager_Add(t *testing.T) {
	m := NewManager()

	token := &Token{
		Value:       "test-value",
		Description: "manually added",
	}

	err := m.Add(token)
	if err != nil {
		t.Fatalf("添加 token 失败: %v", err)
	}

	// 验证 token 被添加
	if m.Count() != 1 {
		t.Errorf("期望 token 数量为 1，得到 %d", m.Count())
	}

	// 验证 ID 被自动生成
	if token.ID == "" {
		t.Error("期望 ID 被自动生成")
	}

	// 验证 CreatedAt 被设置
	if token.CreatedAt.IsZero() {
		t.Error("期望 CreatedAt 被设置")
	}
}

func TestManager_Add_EmptyValue(t *testing.T) {
	m := NewManager()

	token := &Token{
		Description: "no value",
	}

	err := m.Add(token)
	if err == nil {
		t.Error("期望添加空 value 的 token 失败")
	}
}

func TestManager_Validate(t *testing.T) {
	m := NewManager()

	// 添加有效的 token
	token, _ := m.Generate("valid token", nil)

	// 验证有效的 token
	if !m.Validate(token.Value) {
		t.Error("期望 token 有效")
	}

	// 验证无效的 token
	if m.Validate("invalid-token") {
		t.Error("期望 token 无效")
	}

	// 添加过期的 token
	past := time.Now().Add(-1 * time.Hour)
	expiredToken := &Token{
		Value:     "expired-token",
		ExpiresAt: &past,
	}
	m.Add(expiredToken)

	// 验证过期的 token
	if m.Validate("expired-token") {
		t.Error("期望过期的 token 无效")
	}
}

func TestManager_Revoke(t *testing.T) {
	m := NewManager()

	token, _ := m.Generate("to be revoked", nil)

	// 撤销 token
	err := m.Revoke(token.Value)
	if err != nil {
		t.Fatalf("撤销 token 失败: %v", err)
	}

	// 验证 token 已被撤销
	if m.Validate(token.Value) {
		t.Error("期望 token 已被撤销")
	}

	// 再次撤销应该失败
	err = m.Revoke(token.Value)
	if err == nil {
		t.Error("期望撤销不存在的 token 失败")
	}
}

func TestManager_RevokeByID(t *testing.T) {
	m := NewManager()

	token, _ := m.Generate("to be revoked by ID", nil)

	// 通过 ID 撤销 token
	err := m.RevokeByID(token.ID)
	if err != nil {
		t.Fatalf("通过 ID 撤销 token 失败: %v", err)
	}

	// 验证 token 已被撤销
	if m.Validate(token.Value) {
		t.Error("期望 token 已被撤销")
	}

	// 撤销不存在的 ID 应该失败
	err = m.RevokeByID("non-existent-id")
	if err == nil {
		t.Error("期望撤销不存在的 ID 失败")
	}
}

func TestManager_List(t *testing.T) {
	m := NewManager()

	m.Generate("token 1", nil)
	m.Generate("token 2", nil)
	m.Generate("token 3", nil)

	tokens := m.List()

	if len(tokens) != 3 {
		t.Errorf("期望列表长度为 3，得到 %d", len(tokens))
	}
}

func TestManager_Get(t *testing.T) {
	m := NewManager()

	token, _ := m.Generate("get test", nil)

	// 通过 value 获取
	retrieved, exists := m.Get(token.Value)
	if !exists {
		t.Error("期望 token 存在")
	}

	if retrieved.Description != "get test" {
		t.Errorf("期望描述为 'get test'，得到 '%s'", retrieved.Description)
	}

	// 获取不存在的 token
	_, exists = m.Get("non-existent")
	if exists {
		t.Error("期望 token 不存在")
	}
}

func TestManager_GetByID(t *testing.T) {
	m := NewManager()

	token, _ := m.Generate("get by ID test", nil)

	// 通过 ID 获取
	retrieved, exists := m.GetByID(token.ID)
	if !exists {
		t.Error("期望 token 存在")
	}

	if retrieved.Description != "get by ID test" {
		t.Errorf("期望描述为 'get by ID test'，得到 '%s'", retrieved.Description)
	}

	// 获取不存在的 token
	_, exists = m.GetByID("non-existent-id")
	if exists {
		t.Error("期望 token 不存在")
	}
}

func TestManager_Clear(t *testing.T) {
	m := NewManager()

	m.Generate("token 1", nil)
	m.Generate("token 2", nil)

	m.Clear()

	if m.Count() != 0 {
		t.Errorf("期望 token 数量为 0，得到 %d", m.Count())
	}
}

func TestManager_Count(t *testing.T) {
	m := NewManager()

	if m.Count() != 0 {
		t.Error("期望初始 token 数量为 0")
	}

	m.Generate("token 1", nil)
	if m.Count() != 1 {
		t.Error("期望 token 数量为 1")
	}

	m.Generate("token 2", nil)
	if m.Count() != 2 {
		t.Error("期望 token 数量为 2")
	}
}
