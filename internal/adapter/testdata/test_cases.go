package testdata

// TestCase 定义统一的测试用例结构
type TestCase struct {
	Name        string      // 测试用例名称
	Description string      // 测试描述
	Method      string      // 要测试的方法名
	MockAPI     string      // Mock 服务器对应的 API 路径（例如：/cgi/user/get）
	Input       interface{} // 输入参数
	Expected    interface{} // 期望输出
	ShouldError bool        // 是否应该返回错误
	ErrorMsg    string      // 期望的错误消息
	Permission  struct {
		Resource string // 资源类型
		Action   string // 操作类型
		Allowed  bool   // 是否允许
	}
}

// 用户管理测试用例
var UserTestCases = []TestCase{
	{
		Name:        "GetUser_Success",
		Description: "成功获取用户信息",
		Method:      "GetUser",
		MockAPI:     "/cgi/user/get",
		Input: map[string]interface{}{
			"user_id": "10232",
		},
		Expected: map[string]interface{}{
			"user_id": "10232",
			"name":    "Tc-黎明",
			"gender":  0,
			"mobile":  "13728758403",
			"phone":   "02-2999-5691#10232",
			"email":   "liming@addcn.com",
		},
		ShouldError: false,
		Permission: struct {
			Resource string
			Action   string
			Allowed  bool
		}{
			Resource: "user",
			Action:   "read",
			Allowed:  true,
		},
	},
	{
		Name:        "CreateUser_PermissionDenied",
		Description: "创建用户被权限拒绝",
		Method:      "CreateUser",
		MockAPI:     "/cgi/user/create",
		Input: map[string]interface{}{
			"user_id":  "test999",
			"name":     "测试用户",
			"dept_id":  1,
			"gender":   1,
			"mobile":   "13800138999",
			"email":    "test999@example.com",
			"password": "Welcome123",
		},
		Expected: map[string]interface{}{
			"id": "test999",
		},
		ShouldError: true,
		ErrorMsg:    "权限拒绝：不允许对资源 'user' 执行 'create' 操作",
		Permission: struct {
			Resource string
			Action   string
			Allowed  bool
		}{
			Resource: "user",
			Action:   "create",
			Allowed:  false,
		},
	},
	{
		Name:        "DeleteUser_PermissionDenied",
		Description: "删除用户被权限拒绝",
		Method:      "DeleteUser",
		MockAPI:     "/cgi/user/delete",
		Input: map[string]interface{}{
			"user_id": "test999",
		},
		Expected:    map[string]interface{}{},
		ShouldError: true,
		ErrorMsg:    "权限拒绝：不允许对资源 'user' 执行 'delete' 操作",
		Permission: struct {
			Resource string
			Action   string
			Allowed  bool
		}{
			Resource: "user",
			Action:   "delete",
			Allowed:  false,
		},
	},
	{
		Name:        "UpdateUser_PermissionDenied",
		Description: "更新用户被权限拒绝",
		Method:      "UpdateUser",
		MockAPI:     "/cgi/user/update",
		Input: map[string]interface{}{
			"user_id": "10232",
			"name":    "新名字",
		},
		Expected:    map[string]interface{}{},
		ShouldError: true,
		ErrorMsg:    "权限拒绝：不允许对资源 'user' 执行 'update' 操作",
		Permission: struct {
			Resource string
			Action   string
			Allowed  bool
		}{
			Resource: "user",
			Action:   "update",
			Allowed:  false,
		},
	},
}

// 消息管理测试用例
var MessageTestCases = []TestCase{
	{
		Name:        "SendTextMessage_Success",
		Description: "成功发送文本消息",
		Method:      "SendTextMessage",
		MockAPI:     "/cgi/msg/send",
		Input: map[string]interface{}{
			"to_user": "10232",
			"content": "这是一条测试消息",
		},
		Expected: map[string]interface{}{
			"success": true,
		},
		ShouldError: false,
		Permission: struct {
			Resource string
			Action   string
			Allowed  bool
		}{
			Resource: "message",
			Action:   "create",
			Allowed:  true,
		},
	},
	{
		Name:        "SendImageMessage_Success",
		Description: "成功发送图片消息",
		Method:      "SendImageMessage",
		MockAPI:     "/cgi/msg/send",
		Input: map[string]interface{}{
			"to_user":  "10232",
			"media_id": "test_media_id",
		},
		Expected: map[string]interface{}{
			"success": true,
		},
		ShouldError: false,
		Permission: struct {
			Resource string
			Action   string
			Allowed  bool
		}{
			Resource: "message",
			Action:   "create",
			Allowed:  true,
		},
	},
}

// 部门管理测试用例
var DeptTestCases = []TestCase{
	{
		Name:        "GetDeptList_Success",
		Description: "成功获取部门列表",
		Method:      "GetDeptList",
		MockAPI:     "/cgi/dept/list",
		Input: map[string]interface{}{
			"dept_id": 0,
		},
		Expected: map[string]interface{}{
			"dept_list": []map[string]interface{}{
				{"id": 1, "name": "研发部", "parent_id": 0},
				{"id": 8, "name": "技术中心", "parent_id": 1},
			},
		},
		ShouldError: false,
		Permission: struct {
			Resource string
			Action   string
			Allowed  bool
		}{
			Resource: "dept",
			Action:   "read",
			Allowed:  true,
		},
	},
	{
		Name:        "CreateDept_PermissionDenied",
		Description: "创建部门被权限拒绝",
		Method:      "CreateDept",
		MockAPI:     "/cgi/dept/create",
		Input: map[string]interface{}{
			"name":      "测试部门",
			"parent_id": 0,
		},
		Expected: map[string]interface{}{
			"id": 100,
		},
		ShouldError: true,
		ErrorMsg:    "权限拒绝：不允许对资源 'dept' 执行 'create' 操作",
		Permission: struct {
			Resource string
			Action   string
			Allowed  bool
		}{
			Resource: "dept",
			Action:   "create",
			Allowed:  false,
		},
	},
}

// 群组管理测试用例
var GroupTestCases = []TestCase{
	{
		Name:        "CreateGroup_PermissionDenied",
		Description: "创建群组被权限拒绝",
		Method:      "CreateGroup",
		MockAPI:     "/cgi/group/create",
		Input: map[string]interface{}{
			"name": "测试群组",
		},
		Expected: map[string]interface{}{
			"id": "group_123",
		},
		ShouldError: true,
		ErrorMsg:    "权限拒绝：不允许对资源 'group' 执行 'create' 操作",
		Permission: struct {
			Resource string
			Action   string
			Allowed  bool
		}{
			Resource: "group",
			Action:   "create",
			Allowed:  false,
		},
	},
}

// 会话管理测试用例
var SessionTestCases = []TestCase{
	{
		Name:        "CreateSession_PermissionDenied",
		Description: "创建会话被权限拒绝",
		Method:      "CreateSession",
		MockAPI:     "/cgi/session/create",
		Input: map[string]interface{}{
			"title":   "测试会话",
			"creator": "10232",
			"type":    "single",
		},
		Expected: map[string]interface{}{
			"session_id": "session_123",
		},
		ShouldError: true,
		ErrorMsg:    "权限拒绝：不允许对资源 'session' 执行 'create' 操作",
		Permission: struct {
			Resource string
			Action   string
			Allowed  bool
		}{
			Resource: "session",
			Action:   "create",
			Allowed:  false,
		},
	},
}

// AllTestCases 包含所有测试用例
var AllTestCases = []TestCase{}

func init() {
	AllTestCases = append(AllTestCases, UserTestCases...)
	AllTestCases = append(AllTestCases, MessageTestCases...)
	AllTestCases = append(AllTestCases, DeptTestCases...)
	AllTestCases = append(AllTestCases, GroupTestCases...)
	AllTestCases = append(AllTestCases, SessionTestCases...)
}
