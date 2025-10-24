package token

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"fmt"
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
	db *sql.DB // SQLite 数据库连接
}

// NewManager 创建新的 token 管理器
// db 参数可以为 nil，此时将使用内存模式（仅用于测试）
func NewManager(db *sql.DB) *Manager {
	return &Manager{
		db: db,
	}
}

// Generate 生成新的 token
func (m *Manager) Generate(description string, expiresIn *time.Duration) (*Token, error) {
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

	// 存储 token 到数据库
	if m.db != nil {
		if err := m.saveToken(token); err != nil {
			return nil, fmt.Errorf("保存 token 失败: %w", err)
		}
	}

	return token, nil
}

// Add 添加已存在的 token (从配置文件加载)
func (m *Manager) Add(token *Token) error {
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

	// 保存到数据库
	if m.db != nil {
		return m.saveToken(token)
	}

	return nil
}

// Validate 验证 token 是否有效
func (m *Manager) Validate(tokenValue string) bool {
	if m.db == nil {
		return false
	}

	var expiresAt sql.NullTime
	err := m.db.QueryRow(`
		SELECT expires_at FROM tokens WHERE value = ?
	`, tokenValue).Scan(&expiresAt)

	if err != nil {
		return false
	}

	// 检查是否过期
	if expiresAt.Valid && time.Now().After(expiresAt.Time) {
		return false
	}

	return true
}

// Revoke 撤销 token
func (m *Manager) Revoke(tokenValue string) error {
	if m.db == nil {
		return fmt.Errorf("数据库未初始化")
	}

	result, err := m.db.Exec(`DELETE FROM tokens WHERE value = ?`, tokenValue)
	if err != nil {
		return fmt.Errorf("撤销 token 失败: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("检查删除结果失败: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("token 不存在")
	}

	return nil
}

// RevokeByID 通过 ID 撤销 token
func (m *Manager) RevokeByID(tokenID string) error {
	if m.db == nil {
		return fmt.Errorf("数据库未初始化")
	}

	result, err := m.db.Exec(`DELETE FROM tokens WHERE id = ?`, tokenID)
	if err != nil {
		return fmt.Errorf("撤销 token 失败: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("检查删除结果失败: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("token ID %s 不存在", tokenID)
	}

	return nil
}

// List 列出所有 token
func (m *Manager) List() []*Token {
	if m.db == nil {
		return []*Token{}
	}

	rows, err := m.db.Query(`
		SELECT id, value, description, created_at, expires_at
		FROM tokens
		ORDER BY created_at DESC
	`)
	if err != nil {
		return []*Token{}
	}
	defer rows.Close()

	var tokens []*Token
	for rows.Next() {
		var token Token
		var expiresAt sql.NullTime

		err := rows.Scan(
			&token.ID,
			&token.Value,
			&token.Description,
			&token.CreatedAt,
			&expiresAt,
		)
		if err != nil {
			continue
		}

		if expiresAt.Valid {
			token.ExpiresAt = &expiresAt.Time
		}

		tokens = append(tokens, &token)
	}

	return tokens
}

// Get 通过 value 获取 token
func (m *Manager) Get(tokenValue string) (*Token, bool) {
	if m.db == nil {
		return nil, false
	}

	var token Token
	var expiresAt sql.NullTime

	err := m.db.QueryRow(`
		SELECT id, value, description, created_at, expires_at
		FROM tokens
		WHERE value = ?
	`, tokenValue).Scan(
		&token.ID,
		&token.Value,
		&token.Description,
		&token.CreatedAt,
		&expiresAt,
	)

	if err != nil {
		return nil, false
	}

	if expiresAt.Valid {
		token.ExpiresAt = &expiresAt.Time
	}

	return &token, true
}

// GetByID 通过 ID 获取 token
func (m *Manager) GetByID(tokenID string) (*Token, bool) {
	if m.db == nil {
		return nil, false
	}

	var token Token
	var expiresAt sql.NullTime

	err := m.db.QueryRow(`
		SELECT id, value, description, created_at, expires_at
		FROM tokens
		WHERE id = ?
	`, tokenID).Scan(
		&token.ID,
		&token.Value,
		&token.Description,
		&token.CreatedAt,
		&expiresAt,
	)

	if err != nil {
		return nil, false
	}

	if expiresAt.Valid {
		token.ExpiresAt = &expiresAt.Time
	}

	return &token, true
}

// Clear 清除所有 token
func (m *Manager) Clear() {
	if m.db != nil {
		m.db.Exec(`DELETE FROM tokens`)
	}
}

// Count 返回 token 数量
func (m *Manager) Count() int {
	if m.db == nil {
		return 0
	}

	var count int
	err := m.db.QueryRow(`SELECT COUNT(*) FROM tokens`).Scan(&count)
	if err != nil {
		return 0
	}

	return count
}

// saveToken 保存 token 到数据库
func (m *Manager) saveToken(token *Token) error {
	var expiresAt interface{}
	if token.ExpiresAt != nil {
		expiresAt = token.ExpiresAt
	}

	_, err := m.db.Exec(`
		INSERT OR REPLACE INTO tokens (id, value, description, created_at, expires_at)
		VALUES (?, ?, ?, ?, ?)
	`, token.ID, token.Value, token.Description, token.CreatedAt, expiresAt)

	return err
}

// generateRandomToken 生成指定字节长度的随机 token
func generateRandomToken(byteLength int) (string, error) {
	b := make([]byte, byteLength)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
