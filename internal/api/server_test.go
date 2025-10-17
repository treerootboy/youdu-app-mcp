package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/yourusername/youdu-app-mcp/internal/adapter/testdata"
	"github.com/yourusername/youdu-app-mcp/internal/config"
)

// setupTestServer 创建测试用的 HTTP 服务器
func setupTestServer(t *testing.T) *Server {
	t.Helper()

	// 从测试配置文件加载配置
	cfg, err := config.LoadFromFile("../../config_test.yaml")
	if err != nil {
		t.Fatalf("加载测试配置失败: %v", err)
	}

	// 启动 Mock YouDu Server，传入与配置相同的 AesKey 和 AppID
	mockServer := testdata.NewMockYouDuServer(cfg.Youdu.AesKey, cfg.Youdu.AppID)
	t.Cleanup(func() { mockServer.Close() })

	// 覆盖 Mock Server URL（测试时动态生成）
	cfg.Youdu.Addr = mockServer.URL()

	// 创建服务器（权限配置已通过 config_test.yaml 加载）
	server, err := New(cfg)
	if err != nil {
		t.Fatalf("创建服务器失败: %v", err)
	}

	return server
}

// TestHealthCheck 测试健康检查端点
func TestHealthCheck(t *testing.T) {
	server := setupTestServer(t)
	if server == nil {
		return
	}

	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()

	server.router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("期望状态码 %d, 得到 %d", http.StatusOK, w.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("解析响应失败: %v", err)
	}

	if response["status"] != "ok" {
		t.Errorf("期望 status=ok, 得到 %v", response["status"])
	}
}

// TestGetEndpoints 测试获取 API 列表
func TestGetEndpoints(t *testing.T) {
	server := setupTestServer(t)
	if server == nil {
		return
	}

	req := httptest.NewRequest("GET", "/api/v1/endpoints", nil)
	w := httptest.NewRecorder()

	server.router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("期望状态码 %d, 得到 %d", http.StatusOK, w.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("解析响应失败: %v", err)
	}

	count := int(response["count"].(float64))
	if count != 28 {
		t.Errorf("期望 28 个 endpoints, 得到 %d", count)
	}
}

// TestAPI_UserOperations 测试用户相关 API
func TestAPI_UserOperations(t *testing.T) {
	server := setupTestServer(t)
	if server == nil {
		return
	}

	for _, tc := range testdata.UserTestCases {
		t.Run(tc.Name, func(t *testing.T) {
			inputBytes, err := json.Marshal(tc.Input)
			if err != nil {
				t.Fatalf("序列化输入失败: %v", err)
			}

			methodName := toSnakeCase(tc.Method)
			req := httptest.NewRequest("POST", "/api/v1/"+methodName, bytes.NewBuffer(inputBytes))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			server.router.ServeHTTP(w, req)

			var response map[string]interface{}
			if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
				t.Fatalf("解析响应失败: %v", err)
			}

			if tc.ShouldError {
				if response["error"] != true {
					t.Errorf("期望返回错误，但没有错误")
					return
				}
				if tc.ErrorMsg != "" {
					if msg, ok := response["message"].(string); ok {
						if msg != tc.ErrorMsg {
							t.Errorf("错误消息不匹配\n期望: %s\n实际: %s", tc.ErrorMsg, msg)
						}
					}
				}
				return
			}

			if response["error"] == true {
				t.Errorf("不期望错误，但得到: %v", response["message"])
			}
		})
	}
}

// TestAPI_MessageOperations 测试消息相关 API
func TestAPI_MessageOperations(t *testing.T) {
	server := setupTestServer(t)
	if server == nil {
		return
	}

	for _, tc := range testdata.MessageTestCases {
		t.Run(tc.Name, func(t *testing.T) {
			inputBytes, err := json.Marshal(tc.Input)
			if err != nil {
				t.Fatalf("序列化输入失败: %v", err)
			}

			methodName := toSnakeCase(tc.Method)
			req := httptest.NewRequest("POST", "/api/v1/"+methodName, bytes.NewBuffer(inputBytes))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			server.router.ServeHTTP(w, req)

			var response map[string]interface{}
			if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
				t.Fatalf("解析响应失败: %v", err)
			}

			if tc.ShouldError {
				if response["error"] != true {
					t.Errorf("期望返回错误，但没有错误")
					return
				}
				return
			}

			if response["error"] == true {
				t.Errorf("不期望错误，但得到: %v", response["message"])
			}
		})
	}
}
