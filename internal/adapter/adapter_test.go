package adapter

import (
	"context"
	"encoding/json"
	"reflect"
	"testing"

	"github.com/yourusername/youdu-app-mcp/internal/adapter/testdata"
	"github.com/yourusername/youdu-app-mcp/internal/config"
	"github.com/yourusername/youdu-app-mcp/internal/permission"
)

// MockAdapter 用于测试的模拟适配器
type MockAdapter struct {
	*Adapter
	mockMode bool
}

// setupTestAdapter 创建测试用的适配器
func setupTestAdapter(t *testing.T) *Adapter {
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

	// 创建适配器（权限配置已通过 config_test.yaml 加载）
	adapter, err := New(cfg)
	if err != nil {
		t.Fatalf("创建适配器失败: %v", err)
	}

	return adapter
}

// callMethod 使用反射调用适配器方法
func callMethod(adapter *Adapter, methodName string, input interface{}) (interface{}, error) {
	// 获取方法
	method := reflect.ValueOf(adapter).MethodByName(methodName)
	if !method.IsValid() {
		return nil, nil
	}

	// 获取方法类型
	methodType := method.Type()
	if methodType.NumIn() != 2 {
		return nil, nil
	}

	// 获取输入类型
	inputType := methodType.In(1)

	// 将 map 转换为对应的结构体
	inputBytes, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}

	inputValue := reflect.New(inputType)
	if err := json.Unmarshal(inputBytes, inputValue.Interface()); err != nil {
		return nil, err
	}

	// 调用方法
	ctx := context.Background()
	results := method.Call([]reflect.Value{
		reflect.ValueOf(ctx),
		inputValue.Elem(),
	})

	// 处理返回值
	var output interface{}
	var callErr error

	if len(results) > 0 && !results[0].IsNil() {
		output = results[0].Interface()
	}

	if len(results) > 1 && !results[1].IsNil() {
		callErr = results[1].Interface().(error)
	}

	return output, callErr
}

// TestAdapter_UserOperations 测试用户相关操作
func TestAdapter_UserOperations(t *testing.T) {
	adapter := setupTestAdapter(t)
	if adapter == nil {
		return
	}
	defer adapter.Close()

	for _, tc := range testdata.UserTestCases {
		t.Run(tc.Name, func(t *testing.T) {
			// 调用方法
			output, err := callMethod(adapter, tc.Method, tc.Input)

			// 验证错误
			if tc.ShouldError {
				if err == nil {
					t.Errorf("期望返回错误，但没有错误")
					return
				}
				if tc.ErrorMsg != "" && err.Error() != tc.ErrorMsg {
					t.Errorf("错误消息不匹配\n期望: %s\n实际: %s", tc.ErrorMsg, err.Error())
				}
				return
			}

			// 验证无错误
			if err != nil {
				t.Errorf("不期望错误，但得到: %v", err)
				return
			}

			// 验证输出（简单验证非空）
			if output == nil {
				t.Errorf("期望非空输出，但得到 nil")
			}
		})
	}
}

// TestAdapter_MessageOperations 测试消息相关操作
func TestAdapter_MessageOperations(t *testing.T) {
	adapter := setupTestAdapter(t)
	if adapter == nil {
		return
	}
	defer adapter.Close()

	for _, tc := range testdata.MessageTestCases {
		t.Run(tc.Name, func(t *testing.T) {
			// 调用方法
			output, err := callMethod(adapter, tc.Method, tc.Input)

			// 验证错误
			if tc.ShouldError {
				if err == nil {
					t.Errorf("期望返回错误，但没有错误")
					return
				}
				if tc.ErrorMsg != "" && err.Error() != tc.ErrorMsg {
					t.Errorf("错误消息不匹配\n期望: %s\n实际: %s", tc.ErrorMsg, err.Error())
				}
				return
			}

			// 验证无错误
			if err != nil {
				t.Errorf("不期望错误，但得到: %v", err)
				return
			}

			// 验证输出（简单验证非空）
			if output == nil {
				t.Errorf("期望非空输出，但得到 nil")
			}
		})
	}
}

// TestAdapter_DeptOperations 测试部门相关操作
func TestAdapter_DeptOperations(t *testing.T) {
	adapter := setupTestAdapter(t)
	if adapter == nil {
		return
	}
	defer adapter.Close()

	for _, tc := range testdata.DeptTestCases {
		t.Run(tc.Name, func(t *testing.T) {
			// 调用方法
			output, err := callMethod(adapter, tc.Method, tc.Input)

			// 验证错误
			if tc.ShouldError {
				if err == nil {
					t.Errorf("期望返回错误，但没有错误")
					return
				}
				if tc.ErrorMsg != "" && err.Error() != tc.ErrorMsg {
					t.Errorf("错误消息不匹配\n期望: %s\n实际: %s", tc.ErrorMsg, err.Error())
				}
				return
			}

			// 验证无错误
			if err != nil {
				t.Errorf("不期望错误，但得到: %v", err)
				return
			}

			// 验证输出（简单验证非空）
			if output == nil {
				t.Errorf("期望非空输出，但得到 nil")
			}
		})
	}
}

// TestAdapter_GroupOperations 测试群组相关操作
func TestAdapter_GroupOperations(t *testing.T) {
	adapter := setupTestAdapter(t)
	if adapter == nil {
		return
	}
	defer adapter.Close()

	for _, tc := range testdata.GroupTestCases {
		t.Run(tc.Name, func(t *testing.T) {
			// 调用方法
			output, err := callMethod(adapter, tc.Method, tc.Input)

			// 验证错误
			if tc.ShouldError {
				if err == nil {
					t.Errorf("期望返回错误，但没有错误")
					return
				}
				if tc.ErrorMsg != "" && err.Error() != tc.ErrorMsg {
					t.Errorf("错误消息不匹配\n期望: %s\n实际: %s", tc.ErrorMsg, err.Error())
				}
				return
			}

			// 验证无错误
			if err != nil {
				t.Errorf("不期望错误，但得到: %v", err)
				return
			}

			// 验证输出（简单验证非空）
			if output == nil {
				t.Errorf("期望非空输出，但得到 nil")
			}
		})
	}
}

// TestAdapter_SessionOperations 测试会话相关操作
func TestAdapter_SessionOperations(t *testing.T) {
	adapter := setupTestAdapter(t)
	if adapter == nil {
		return
	}
	defer adapter.Close()

	for _, tc := range testdata.SessionTestCases {
		t.Run(tc.Name, func(t *testing.T) {
			// 调用方法
			output, err := callMethod(adapter, tc.Method, tc.Input)

			// 验证错误
			if tc.ShouldError {
				if err == nil {
					t.Errorf("期望返回错误，但没有错误")
					return
				}
				if tc.ErrorMsg != "" && err.Error() != tc.ErrorMsg {
					t.Errorf("错误消息不匹配\n期望: %s\n实际: %s", tc.ErrorMsg, err.Error())
				}
				return
			}

			// 验证无错误
			if err != nil {
				t.Errorf("不期望错误，但得到: %v", err)
				return
			}

			// 验证输出（简单验证非空）
			if output == nil {
				t.Errorf("期望非空输出，但得到 nil")
			}
		})
	}
}

// TestAdapter_PermissionSystem 测试权限系统
func TestAdapter_PermissionSystem(t *testing.T) {
	adapter := setupTestAdapter(t)
	if adapter == nil {
		return
	}
	defer adapter.Close()

	tests := []struct {
		resource permission.Resource
		action   permission.Action
		expected bool
	}{
		{permission.ResourceUser, permission.ActionRead, true},
		{permission.ResourceUser, permission.ActionCreate, false},
		{permission.ResourceUser, permission.ActionUpdate, false},
		{permission.ResourceUser, permission.ActionDelete, false},
		{permission.ResourceMessage, permission.ActionCreate, true},
		{permission.ResourceDept, permission.ActionRead, true},
		{permission.ResourceDept, permission.ActionCreate, false},
	}

	for _, tt := range tests {
		t.Run(string(tt.resource)+"_"+string(tt.action), func(t *testing.T) {
			err := adapter.checkPermission(tt.resource, tt.action)
			hasPermission := (err == nil)

			if hasPermission != tt.expected {
				t.Errorf("权限检查结果不符合预期\n资源: %s\n操作: %s\n期望: %v\n实际: %v",
					tt.resource, tt.action, tt.expected, hasPermission)
			}
		})
	}
}
