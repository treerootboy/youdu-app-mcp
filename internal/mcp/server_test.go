package mcp

import (
	"bufio"
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/yourusername/youdu-app-mcp/internal/adapter/testdata"
	"github.com/yourusername/youdu-app-mcp/internal/config"
	"gopkg.in/yaml.v3"
)

// TestMCP_WithMockServer 使用 Mock 服务器进行真实的 MCP 协议测试
func TestMCP_WithMockServer(t *testing.T) {
	// 1. 从测试配置文件加载配置
	cfg, err := config.LoadFromFile("../../config_test.yaml")
	if err != nil {
		t.Fatalf("加载测试配置失败: %v", err)
	}

	// 2. 启动 Mock YouDu Server，传入与配置相同的 AesKey 和 AppID
	mockServer := testdata.NewMockYouDuServer(cfg.Youdu.AesKey, cfg.Youdu.AppID)
	defer mockServer.Close()

	t.Logf("Mock Server 启动在: %s", mockServer.URL())

	// 3. 确保二进制文件已构建
	buildMCPBinary(t)

	// 4. 创建临时配置文件（基于 config_test.yaml，但覆盖 youdu.addr）
	tempConfigPath := createTempTestConfig(t, mockServer.URL())
	defer os.Remove(tempConfigPath)

	// 4. 启动 MCP 服务器进程，使用临时配置文件
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	rootDir := "../.."
	cmd := exec.CommandContext(ctx, "./bin/youdu-mcp")
	cmd.Dir = rootDir

	// 通过环境变量指定临时配置文件
	cmd.Env = []string{
		"PATH=" + os.Getenv("PATH"),
		"YOUDU_CONFIG_FILE=" + tempConfigPath,
	}

	stdin, err := cmd.StdinPipe()
	if err != nil {
		t.Fatalf("创建 stdin pipe 失败: %v", err)
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		t.Fatalf("创建 stdout pipe 失败: %v", err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		t.Fatalf("创建 stderr pipe 失败: %v", err)
	}

	if err := cmd.Start(); err != nil {
		t.Fatalf("启动 MCP 服务器失败: %v", err)
	}
	defer cmd.Process.Kill()

	// 打印 stderr 以便调试
	go func() {
		stderrScanner := bufio.NewScanner(stderr)
		for stderrScanner.Scan() {
			t.Logf("[MCP stderr] %s", stderrScanner.Text())
		}
	}()

	scanner := bufio.NewScanner(stdout)
	encoder := json.NewEncoder(stdin)

	// 5. 等待 MCP 服务器启动
	time.Sleep(500 * time.Millisecond)

	// 6. MCP 协议初始化（必须在所有工具调用之前完成）
	// 注意：这不能放在 t.Run() 中，因为子测试过滤会跳过它
	initRequest := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "initialize",
		"params": map[string]interface{}{
			"protocolVersion": "2024-11-05",
			"capabilities":    map[string]interface{}{},
			"clientInfo": map[string]interface{}{
				"name":    "test-client",
				"version": "1.0.0",
			},
		},
	}

	if err := encoder.Encode(initRequest); err != nil {
		t.Fatalf("发送初始化请求失败: %v", err)
	}

	if !scanner.Scan() {
		t.Fatal("未收到初始化响应")
	}

	var initResponse map[string]interface{}
	if err := json.Unmarshal(scanner.Bytes(), &initResponse); err != nil {
		t.Fatalf("解析初始化响应失败: %v", err)
	}

	if _, ok := initResponse["result"]; !ok {
		t.Fatalf("初始化响应缺少 result 字段: %+v", initResponse)
	}

	t.Logf("✓ MCP 协议初始化成功")

	// 7. 测试 - 获取工具列表
	t.Run("ListTools", func(t *testing.T) {
		request := map[string]interface{}{
			"jsonrpc": "2.0",
			"id":      2,
			"method":  "tools/list",
		}

		if err := encoder.Encode(request); err != nil {
			t.Fatalf("发送工具列表请求失败: %v", err)
		}

		if !scanner.Scan() {
			t.Fatal("未收到工具列表响应")
		}

		var response map[string]interface{}
		if err := json.Unmarshal(scanner.Bytes(), &response); err != nil {
			t.Fatalf("解析工具列表响应失败: %v", err)
		}

		result, ok := response["result"].(map[string]interface{})
		if !ok {
			t.Fatal("工具列表响应格式错误")
		}

		tools, ok := result["tools"].([]interface{})
		if !ok {
			t.Fatal("工具列表格式错误")
		}

		if len(tools) != 28 {
			t.Errorf("期望 28 个工具，得到 %d 个", len(tools))
		}

		t.Logf("✓ 工具列表获取成功: %d 个工具", len(tools))
	})

	// 8. 使用 testdata.AllTestCases 进行完整测试
	requestID := 3
	for _, tc := range testdata.AllTestCases {
		t.Run(tc.Name, func(t *testing.T) {
			// 转换方法名为 snake_case
			toolName := toSnakeCase(tc.Method)

			request := map[string]interface{}{
				"jsonrpc": "2.0",
				"id":      requestID,
				"method":  "tools/call",
				"params": map[string]interface{}{
					"name":      toolName,
					"arguments": tc.Input,
				},
			}
			requestID++

			if err := encoder.Encode(request); err != nil {
				t.Fatalf("发送请求失败: %v", err)
			}

			if !scanner.Scan() {
				t.Fatal("未收到响应")
			}

			var response map[string]interface{}
			if err := json.Unmarshal(scanner.Bytes(), &response); err != nil {
				t.Fatalf("解析响应失败: %v", err)
			}

			// 验证响应
			if tc.ShouldError {
				// 期望返回错误
				if result, ok := response["result"].(map[string]interface{}); ok {
					if isError, _ := result["isError"].(bool); !isError {
						t.Errorf("期望操作被拒绝，但成功了\n完整响应: %+v\n测试用例: %+v", response, tc)
					} else {
						t.Logf("✓ 权限控制正常: %s", tc.ErrorMsg)
					}
				} else {
					t.Errorf("响应格式错误，无法获取 result 字段: %+v", response)
				}
			} else {
				// 期望成功
				if _, ok := response["result"]; !ok {
					t.Errorf("期望成功响应，但失败了: %+v", response)
				} else {
					t.Logf("✓ 操作成功: %s", tc.Method)
				}
			}
		})
	}
}

// createTempTestConfig 创建临时测试配置文件
// 基于 config_test.yaml，但覆盖 youdu.addr 为 Mock Server URL
func createTempTestConfig(t *testing.T, mockServerURL string) string {
	t.Helper()

	// 读取 config_test.yaml
	configTestPath := "../../config_test.yaml"
	data, err := os.ReadFile(configTestPath)
	if err != nil {
		t.Fatalf("读取 config_test.yaml 失败: %v", err)
	}

	// 解析 YAML
	var configData map[string]interface{}
	if err := yaml.Unmarshal(data, &configData); err != nil {
		t.Fatalf("解析 config_test.yaml 失败: %v", err)
	}

	// 覆盖 youdu.addr
	if youdu, ok := configData["youdu"].(map[string]interface{}); ok {
		youdu["addr"] = mockServerURL
	} else {
		t.Fatal("config_test.yaml 中缺少 youdu 配置")
	}

	// 写入临时文件
	tempDir := t.TempDir()
	tempConfigPath := filepath.Join(tempDir, "config_mcp_test.yaml")

	modifiedData, err := yaml.Marshal(configData)
	if err != nil {
		t.Fatalf("序列化配置失败: %v", err)
	}

	if err := os.WriteFile(tempConfigPath, modifiedData, 0600); err != nil {
		t.Fatalf("写入临时配置文件失败: %v", err)
	}

	absPath, err := filepath.Abs(tempConfigPath)
	if err != nil {
		t.Fatalf("获取临时配置文件绝对路径失败: %v", err)
	}

	t.Logf("✓ 创建临时配置文件: %s", absPath)
	return absPath
}

// buildMCPBinary 构建 MCP 二进制文件
func buildMCPBinary(t *testing.T) {
	t.Helper()

	// 从项目根目录构建
	cmd := exec.Command("sh", "-c", "cd ../.. && go build -o bin/youdu-mcp ./cmd/youdu-mcp")
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("构建 MCP 二进制文件失败: %v\n%s", err, output)
	}

	t.Log("✓ MCP 二进制文件构建成功")
}

// TestMCP_UnitTests 单元测试 - 不启动真实进程
func TestMCP_UnitTests(t *testing.T) {
	// 从测试配置文件加载配置
	cfg, err := config.LoadFromFile("../../config_test.yaml")
	if err != nil {
		t.Fatalf("加载测试配置失败: %v", err)
	}

	// 启动 Mock Server，传入与配置相同的 AesKey 和 AppID
	mockServer := testdata.NewMockYouDuServer(cfg.Youdu.AesKey, cfg.Youdu.AppID)
	defer mockServer.Close()

	// 覆盖 Mock Server URL
	cfg.Youdu.Addr = mockServer.URL()

	// 创建 MCP 服务器
	server, err := New(cfg)
	if err != nil {
		t.Fatalf("创建 MCP 服务器失败: %v", err)
	}

	// 验证服务器创建成功
	if server == nil {
		t.Fatal("MCP 服务器为 nil")
	}

	if server.server == nil {
		t.Fatal("MCP server 实例为 nil")
	}

	if server.adapter == nil {
		t.Fatal("Adapter 实例为 nil")
	}

	t.Log("✓ MCP 服务器创建成功")
}

