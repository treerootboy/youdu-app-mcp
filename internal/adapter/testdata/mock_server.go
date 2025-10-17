package testdata

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/addcnos/youdu/v2"
)

// MockYouDuServer Mock 有度服务器，基于 AllTestCases 自动生成响应
type MockYouDuServer struct {
	server     *httptest.Server
	apiMapping map[string][]*TestCase // 改为切片，支持同一 API 的多个测试用例
	encryptor  *youdu.Encryptor        // 加密器，用于加密响应
}

// NewMockYouDuServer 创建 Mock 有度服务器
// 参数：
//   - aesKey: AES密钥（base64编码的字符串）
//   - appID: 应用ID
func NewMockYouDuServer(aesKey, appID string) *MockYouDuServer {
	// 将 base64 编码的 AES Key 解码
	aesKeyBytes, err := base64.StdEncoding.DecodeString(aesKey)
	if err != nil {
		panic("invalid aes key: " + err.Error())
	}

	mock := &MockYouDuServer{
		apiMapping: make(map[string][]*TestCase),
		encryptor:  youdu.NewEncryptor(aesKeyBytes, appID),
	}

	// 从 AllTestCases 构建 API 路径到测试用例的映射
	for i := range AllTestCases {
		tc := &AllTestCases[i]
		if tc.MockAPI != "" {
			mock.apiMapping[tc.MockAPI] = append(mock.apiMapping[tc.MockAPI], tc)
		}
	}

	// 创建 HTTP 服务器
	mux := http.NewServeMux()

	// 注册 token 获取接口（所有请求都需要先获取token）
	mux.HandleFunc("/cgi/gettoken", mock.handleRequest)

	// 为每个唯一的 MockAPI 注册处理器
	registeredAPIs := make(map[string]bool)
	for _, tc := range AllTestCases {
		if tc.MockAPI != "" && !registeredAPIs[tc.MockAPI] {
			registeredAPIs[tc.MockAPI] = true
			mux.HandleFunc(tc.MockAPI, mock.handleRequest)
		}
	}

	mock.server = httptest.NewServer(mux)
	return mock
}

// URL 返回 Mock 服务器 URL
func (m *MockYouDuServer) URL() string {
	return m.server.URL
}

// Close 关闭 Mock 服务器
func (m *MockYouDuServer) Close() {
	m.server.Close()
}

// handleRequest 处理所有 API 请求，根据 test_cases.go 中的 Expected 返回加密响应
func (m *MockYouDuServer) handleRequest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// 调试：打印收到的请求（已禁用）
	// fmt.Printf("Mock Server 收到请求: %s %s\n", r.Method, r.URL.Path)

	// 特殊处理：token 获取请求
	if r.URL.Path == "/cgi/gettoken" {
		m.handleTokenRequest(w, r)
		return
	}

	// 查找匹配的测试用例列表
	testCases, ok := m.apiMapping[r.URL.Path]
	if !ok {
		// 未找到匹配的测试用例，返回默认成功响应
		m.writeEncryptedResponse(w, map[string]interface{}{
			"errcode": 0,
			"errmsg":  "ok",
		})
		return
	}

	// 对于同一 API 的多个测试用例，返回第一个不期望错误的用例的响应
	// 因为权限检查已经在 adapter 层完成，Mock Server 只需要返回成功响应
	for _, tc := range testCases {
		if !tc.ShouldError {
			response := m.generateResponse(tc)
			m.writeEncryptedResponse(w, response)
			return
		}
	}

	// 如果所有测试用例都期望错误（不应该发生，因为权限检查在 adapter 层），返回默认成功响应
	m.writeEncryptedResponse(w, map[string]interface{}{
		"errcode": 0,
		"errmsg":  "ok",
	})
}

// generateResponse 根据测试用例的 Expected 字段生成 Mock 响应
// 返回符合有度IM API规范的响应格式
func (m *MockYouDuServer) generateResponse(tc *TestCase) map[string]interface{} {
	// 基础响应 - 符合有度IM API规范
	// 根据文档，成功响应必须包含 errcode=0 和 errmsg="ok"
	response := map[string]interface{}{
		"errcode": 0,
		"errmsg":  "ok",
	}

	// 将 Expected 中的数据合并到响应中
	// 这些数据会作为响应的额外字段返回
	if tc.Expected != nil {
		if expectedMap, ok := tc.Expected.(map[string]interface{}); ok {
			for k, v := range expectedMap {
				response[k] = v
			}
		}
	}

	return response
}

// handleTokenRequest 处理 access token 获取请求
// 根据有度IM API规范，返回一个有效的 access token
func (m *MockYouDuServer) handleTokenRequest(w http.ResponseWriter, r *http.Request) {
	// Token 响应格式
	tokenResponse := map[string]interface{}{
		"errcode":      0,
		"errmsg":       "ok",
		"accessToken":  "mock_access_token_for_testing",
		"expiresIn":    7200, // 2小时有效期
	}

	m.writeEncryptedResponse(w, tokenResponse)
}

// writeEncryptedResponse 将响应数据加密后写入 HTTP 响应
// 根据有度IM API规范，响应格式为 {"encrypt": "加密后的内容"}
func (m *MockYouDuServer) writeEncryptedResponse(w http.ResponseWriter, data map[string]interface{}) {
	// 1. 将响应数据转换为 JSON
	plaintext, err := json.Marshal(data)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	// 2. 使用 YouDu SDK 的 Encryptor 加密响应
	encrypted, err := m.encryptor.Encrypt(plaintext)
	if err != nil {
		http.Error(w, "encryption error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 3. 构造加密响应格式：{"encrypt": "加密后的内容"}
	encryptedResponse := map[string]interface{}{
		"encrypt": encrypted,
	}

	// 4. 写入响应
	responseData, _ := json.Marshal(encryptedResponse)
	w.WriteHeader(http.StatusOK)
	w.Write(responseData)
}
