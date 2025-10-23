package token

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"sync"
	"time"
)

// Token 代表一个访问令牌
type Token struct {
	ID          string    `json:"id" yaml:"id"`                     // Token ID
	Value       string    `json:"value" yaml:"value"`               // Token 值
	Description string    `json:"description" yaml:"description"`   // 描述
	CreatedAt   time.Time `json:"created_at" yaml:"created_at"`     // 创建时间
	ExpiresAt   *time.Time `json:"expires_at,omitempty" yaml:"expires_at,omitempty"` // 过期时间 (可选)
}

// Manager 管理所有 token
type Manager struct {
	mu     sync.RWMutex
	tokens map[string]*Token // key 为 token value
}

// NewManager 创建新的 token 管理器
func NewManager() *Manager {
	return &Manager{
		tokens: make(map[string]*Token),
	}
}

// Generate 生成新的 token
func (m *Manager) Generate(description string, expiresIn *time.Duration) (*Token, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// 生成随机 token
	tokenValue, err := generateRandomToken(32)
	if err != nil {
		return nil, fmt.Errorf("生成 token 失败: %w", err)
	}

	// 生成 token ID (短一些，用于识别)
	tokenID, err := generateRandomToken(8)
	if err != nil {
		return nil, fmt.Errorf("生成 token ID 失败: %w", err)
	}

	token := &Token{
		ID:          tokenID,
		Value:       tokenValue,
		Description: description,
		CreatedAt:   time.Now(),
	}

	// 设置过期时间
	if expiresIn != nil {
		expiresAt := time.Now().Add(*expiresIn)
		token.ExpiresAt = &expiresAt
	}

	// 存储 token
	m.tokens[token.Value] = token

	return token, nil
}

// Add 添加已存在的 token (从配置文件加载)
func (m *Manager) Add(token *Token) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if token.Value == "" {
		return fmt.Errorf("token value 不能为空")
	}

	// 如果没有 ID，生成一个
	if token.ID == "" {
		tokenID, err := generateRandomToken(8)
		if err != nil {
			return fmt.Errorf("生成 token ID 失败: %w", err)
		}
		token.ID = tokenID
	}

	// 如果没有创建时间，设置为当前时间
	if token.CreatedAt.IsZero() {
		token.CreatedAt = time.Now()
	}

	m.tokens[token.Value] = token
	return nil
}

// Validate 验证 token 是否有效
func (m *Manager) Validate(tokenValue string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	token, exists := m.tokens[tokenValue]
	if !exists {
		return false
	}

	// 检查是否过期
	if token.ExpiresAt != nil && time.Now().After(*token.ExpiresAt) {
		return false
	}

	return true
}

// Revoke 撤销 token
func (m *Manager) Revoke(tokenValue string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.tokens[tokenValue]; !exists {
		return fmt.Errorf("token 不存在")
	}

	delete(m.tokens, tokenValue)
	return nil
}

// RevokeByID 通过 ID 撤销 token
func (m *Manager) RevokeByID(tokenID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for value, token := range m.tokens {
		if token.ID == tokenID {
			delete(m.tokens, value)
			return nil
		}
	}

	return fmt.Errorf("token ID %s 不存在", tokenID)
}

// List 列出所有 token
func (m *Manager) List() []*Token {
	m.mu.RLock()
	defer m.mu.RUnlock()

	tokens := make([]*Token, 0, len(m.tokens))
	for _, token := range m.tokens {
		tokens = append(tokens, token)
	}

	return tokens
}

// Get 通过 value 获取 token
func (m *Manager) Get(tokenValue string) (*Token, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	token, exists := m.tokens[tokenValue]
	return token, exists
}

// GetByID 通过 ID 获取 token
func (m *Manager) GetByID(tokenID string) (*Token, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, token := range m.tokens {
		if token.ID == tokenID {
			return token, true
		}
	}

	return nil, false
}

// Clear 清除所有 token
func (m *Manager) Clear() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.tokens = make(map[string]*Token)
}

// Count 返回 token 数量
func (m *Manager) Count() int {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return len(m.tokens)
}

// generateRandomToken 生成指定字节长度的随机 token
func generateRandomToken(byteLength int) (string, error) {
	b := make([]byte, byteLength)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
