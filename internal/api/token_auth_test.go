package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	_ "modernc.org/sqlite"

	"github.com/yourusername/youdu-app-mcp/internal/config"
	"github.com/yourusername/youdu-app-mcp/internal/permission"
	"github.com/yourusername/youdu-app-mcp/internal/token"
)

// createTestPermission 创建测试用的权限配置
func createTestPermission() *permission.Permission {
	resources := map[permission.Resource]permission.ResourcePolicy{
		permission.ResourceDept:    {Create: true, Read: true, Update: true, Delete: true},
		permission.ResourceUser:    {Create: true, Read: true, Update: true, Delete: true},
		permission.ResourceGroup:   {Create: true, Read: true, Update: true, Delete: true},
		permission.ResourceSession: {Create: true, Read: true, Update: true, Delete: true},
		permission.ResourceMessage: {Create: true, Read: true, Update: true, Delete: true},
	}
	return permission.New(false, true, resources)
}

// setupTestDB creates a temporary test database
func setupTestDB(t *testing.T) (*sql.DB, func()) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		t.Fatalf("打开测试数据库失败: %v", err)
	}

	// 初始化表结构
	schema := `
	CREATE TABLE IF NOT EXISTS tokens (
		id TEXT PRIMARY KEY,
		value TEXT UNIQUE NOT NULL,
		description TEXT NOT NULL,
		created_at DATETIME NOT NULL,
		expires_at DATETIME
	);
	`
	_, err = db.Exec(schema)
	if err != nil {
		db.Close()
		t.Fatalf("初始化数据库结构失败: %v", err)
	}

	cleanup := func() {
		db.Close()
		os.RemoveAll(tmpDir)
	}

	return db, cleanup
}

func TestTokenAuthMiddleware_NoToken(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	// 创建配置和 token 管理器
	cfg := &config.Config{
		Youdu: config.YouduConfig{
			Addr:   "http://test-server:7080",
			Buin:   12345678,
			AppID:  "test-app",
			AesKey: "MTIzNDU2Nzg5MDEyMzQ1Njc4OTAxMjM0NTY3ODkwMTI=",
		},
		Permission:   createTestPermission(),
		TokenManager: token.NewManager(db),
	}

	// 添加一个 token
	testToken := &token.Token{
		ID:          "test001",
		Value:       "test-token-value",
		Description: "Test token",
	}
	cfg.TokenManager.Add(testToken)

	// 创建服务器
	server, err := New(cfg)
	if err != nil {
		t.Fatalf("创建服务器失败: %v", err)
	}
	defer server.Close()

	// 创建请求（不带 token）
	reqBody := []byte(`{"dept_id": 0}`)
	req := httptest.NewRequest("POST", "/api/v1/get_dept_list", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	// 记录响应
	rr := httptest.NewRecorder()

	// 处理请求
	server.router.ServeHTTP(rr, req)

	// 检查状态码
	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("期望状态码 %v，得到 %v", http.StatusUnauthorized, status)
	}

	// 检查响应内容
	var response map[string]interface{}
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatalf("解析响应失败: %v", err)
	}

	if response["error"] != true {
		t.Error("期望 error 为 true")
	}

	if response["message"] != "缺少 Authorization header" {
		t.Errorf("期望错误消息为 '缺少 Authorization header'，得到 '%v'", response["message"])
	}
}

func TestTokenAuthMiddleware_InvalidToken(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	// 创建配置和 token 管理器
	cfg := &config.Config{
		Youdu: config.YouduConfig{
			Addr:   "http://test-server:7080",
			Buin:   12345678,
			AppID:  "test-app",
			AesKey: "MTIzNDU2Nzg5MDEyMzQ1Njc4OTAxMjM0NTY3ODkwMTI=",
		},
		Permission:   createTestPermission(),
		TokenManager: token.NewManager(db),
	}

	// 添加一个 token
	testToken := &token.Token{
		ID:          "test001",
		Value:       "test-token-value",
		Description: "Test token",
	}
	cfg.TokenManager.Add(testToken)

	// 创建服务器
	server, err := New(cfg)
	if err != nil {
		t.Fatalf("创建服务器失败: %v", err)
	}
	defer server.Close()

	// 创建请求（带无效 token）
	reqBody := []byte(`{"dept_id": 0}`)
	req := httptest.NewRequest("POST", "/api/v1/get_dept_list", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer invalid-token")

	// 记录响应
	rr := httptest.NewRecorder()

	// 处理请求
	server.router.ServeHTTP(rr, req)

	// 检查状态码
	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("期望状态码 %v，得到 %v", http.StatusUnauthorized, status)
	}

	// 检查响应内容
	var response map[string]interface{}
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatalf("解析响应失败: %v", err)
	}

	if response["error"] != true {
		t.Error("期望 error 为 true")
	}

	if response["message"] != "无效的 token" {
		t.Errorf("期望错误消息为 '无效的 token'，得到 '%v'", response["message"])
	}
}

func TestTokenAuthMiddleware_ValidToken(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	// 创建配置和 token 管理器
	cfg := &config.Config{
		Youdu: config.YouduConfig{
			Addr:   "http://test-server:7080",
			Buin:   12345678,
			AppID:  "test-app",
			AesKey: "MTIzNDU2Nzg5MDEyMzQ1Njc4OTAxMjM0NTY3ODkwMTI=",
		},
		Permission:   createTestPermission(),
		TokenManager: token.NewManager(db),
	}

	// 添加一个 token
	testToken := &token.Token{
		ID:          "test001",
		Value:       "test-token-value",
		Description: "Test token",
	}
	cfg.TokenManager.Add(testToken)

	// 创建服务器
	server, err := New(cfg)
	if err != nil {
		t.Fatalf("创建服务器失败: %v", err)
	}
	defer server.Close()

	// 创建请求（带有效 token）
	reqBody := []byte(`{"dept_id": 0}`)
	req := httptest.NewRequest("POST", "/api/v1/get_dept_list", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer test-token-value")

	// 记录响应
	rr := httptest.NewRecorder()

	// 处理请求
	server.router.ServeHTTP(rr, req)

	// 注意: 因为没有实际的 YouDu 服务器，会得到连接错误，但不会是 401 Unauthorized
	// 这证明 token 验证通过了
	if status := rr.Code; status == http.StatusUnauthorized {
		t.Errorf("不应该返回 401 Unauthorized，但得到了")
	}
}

func TestTokenAuthMiddleware_ValidTokenWithoutBearer(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	// 创建配置和 token 管理器
	cfg := &config.Config{
		Youdu: config.YouduConfig{
			Addr:   "http://test-server:7080",
			Buin:   12345678,
			AppID:  "test-app",
			AesKey: "MTIzNDU2Nzg5MDEyMzQ1Njc4OTAxMjM0NTY3ODkwMTI=",
		},
		Permission:   createTestPermission(),
		TokenManager: token.NewManager(db),
	}

	// 添加一个 token
	testToken := &token.Token{
		ID:          "test001",
		Value:       "test-token-value",
		Description: "Test token",
	}
	cfg.TokenManager.Add(testToken)

	// 创建服务器
	server, err := New(cfg)
	if err != nil {
		t.Fatalf("创建服务器失败: %v", err)
	}
	defer server.Close()

	// 创建请求（不带 Bearer 前缀）
	reqBody := []byte(`{"dept_id": 0}`)
	req := httptest.NewRequest("POST", "/api/v1/get_dept_list", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "test-token-value")

	// 记录响应
	rr := httptest.NewRecorder()

	// 处理请求
	server.router.ServeHTTP(rr, req)

	// 注意: 因为没有实际的 YouDu 服务器，会得到连接错误，但不会是 401 Unauthorized
	// 这证明 token 验证通过了
	if status := rr.Code; status == http.StatusUnauthorized {
		t.Errorf("不应该返回 401 Unauthorized，但得到了")
	}
}

func TestHealthEndpoint_NoTokenRequired(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	// 创建配置和 token 管理器
	cfg := &config.Config{
		Youdu: config.YouduConfig{
			Addr:   "http://test-server:7080",
			Buin:   12345678,
			AppID:  "test-app",
			AesKey: "MTIzNDU2Nzg5MDEyMzQ1Njc4OTAxMjM0NTY3ODkwMTI=",
		},
		Permission:   createTestPermission(),
		TokenManager: token.NewManager(db),
	}

	// 添加一个 token
	testToken := &token.Token{
		ID:          "test001",
		Value:       "test-token-value",
		Description: "Test token",
	}
	cfg.TokenManager.Add(testToken)

	// 创建服务器
	server, err := New(cfg)
	if err != nil {
		t.Fatalf("创建服务器失败: %v", err)
	}
	defer server.Close()

	// 创建请求（不带 token）
	req := httptest.NewRequest("GET", "/health", nil)

	// 记录响应
	rr := httptest.NewRecorder()

	// 处理请求
	server.router.ServeHTTP(rr, req)

	// 健康检查应该不需要 token
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("期望状态码 %v，得到 %v", http.StatusOK, status)
	}

	// 检查响应内容
	var response map[string]interface{}
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatalf("解析响应失败: %v", err)
	}

	if response["status"] != "ok" {
		t.Errorf("期望 status 为 'ok'，得到 '%v'", response["status"])
	}
}

func TestEndpointsListing_NoTokenRequired(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	// 创建配置和 token 管理器
	cfg := &config.Config{
		Youdu: config.YouduConfig{
			Addr:   "http://test-server:7080",
			Buin:   12345678,
			AppID:  "test-app",
			AesKey: "MTIzNDU2Nzg5MDEyMzQ1Njc4OTAxMjM0NTY3ODkwMTI=",
		},
		Permission:   createTestPermission(),
		TokenManager: token.NewManager(db),
	}

	// 添加一个 token
	testToken := &token.Token{
		ID:          "test001",
		Value:       "test-token-value",
		Description: "Test token",
	}
	cfg.TokenManager.Add(testToken)

	// 创建服务器
	server, err := New(cfg)
	if err != nil {
		t.Fatalf("创建服务器失败: %v", err)
	}
	defer server.Close()

	// 创建请求（不带 token）
	req := httptest.NewRequest("GET", "/api/v1/endpoints", nil)

	// 记录响应
	rr := httptest.NewRecorder()

	// 处理请求
	server.router.ServeHTTP(rr, req)

	// endpoints 列表应该不需要 token
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("期望状态码 %v，得到 %v", http.StatusOK, status)
	}

	// 检查响应内容
	var response map[string]interface{}
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatalf("解析响应失败: %v", err)
	}

	if _, ok := response["endpoints"]; !ok {
		t.Error("响应中缺少 endpoints 字段")
	}
}
