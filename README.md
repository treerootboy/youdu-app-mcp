# 有度多接口服务

为有度即时通讯提供 CLI、MCP（模型上下文协议）和 API 接口的综合服务，通过统一的适配器层实现。

## 架构

```
   CLI       MCP       HTTP API
    │         │         │
    └─────────┼─────────┘
              │
           适配器
              │
           有度SDK
```

## 特性

- **统一适配器层**：所有有度 SDK 操作都封装在简化的适配器层中
- **三种接口模式**：CLI 命令行、MCP 协议、HTTP REST API
- **自动接口生成**：CLI 命令、MCP 工具和 HTTP endpoints 通过反射自动从适配器方法生成
- **类型安全**：使用 Go 结构体和 JSON schema 注解实现完全类型安全
- **配置管理**：通过配置文件和环境变量灵活配置
- **权限控制**：内置细粒度的资源权限管理系统

## 安装

### 前置要求

- Go 1.23 或更高版本
- 可访问的有度 IM 服务器

### 构建

```bash
# 构建 MCP 服务器
go build -o bin/youdu-mcp ./cmd/youdu-mcp

# 构建 CLI（包含 HTTP API 服务器）
go build -o bin/youdu-cli ./cmd/youdu-cli
```

## 配置

在项目根目录或 `~/.youdu/config.yaml` 创建 `config.yaml` 文件：

```yaml
youdu:
  addr: "http://your-youdu-server:7080"
  buin: 123456789
  app_id: "your-app-id"
  aes_key: "your-aes-key"
```

或使用环境变量：

```bash
export YOUDU_ADDR="http://your-youdu-server:7080"
export YOUDU_BUIN=123456789
export YOUDU_APP_ID="your-app-id"
export YOUDU_AES_KEY="your-aes-key"
```

## 使用方法

### CLI

CLI 提供按功能组织的命令：

```bash
# 列出所有命令
./bin/youdu-cli --help

# 部门操作
./bin/youdu-cli dept get-list --dept-id=0
./bin/youdu-cli dept get-user-list --dept-id=1
./bin/youdu-cli dept create --name="技术部" --parent-id=0

# 用户操作
./bin/youdu-cli user get --user-id="user123"
./bin/youdu-cli user create --user-id="newuser" --name="新用户" --dept-id=1

# 消息操作
./bin/youdu-cli message send-text-message --to-user="user123" --content="你好！"

# 群组操作
./bin/youdu-cli group get-list --user-id="user123"
./bin/youdu-cli group create --name="项目组"

# 会话操作
./bin/youdu-cli session create --title="团队聊天" --creator="user123" --type="group"
```

### MCP 服务器

MCP 服务器将所有适配器方法作为 MCP 工具提供，可被 Claude Desktop 或其他 MCP 客户端调用。

#### 运行 MCP 服务器

```bash
./bin/youdu-mcp
```

#### Claude Desktop 集成

添加到 Claude Desktop 配置（macOS 上的 `~/Library/Application Support/Claude/claude_desktop_config.json`）：

```json
{
  "mcpServers": {
    "youdu": {
      "command": "/path/to/youdu-app-mcp/bin/youdu-mcp"
    }
  }
}
```

#### 可用的 MCP 工具

所有工具遵循 snake_case 命名规范：

- **部门**：`get_dept_list`、`get_dept_user_list`、`get_dept_alias_list`、`create_dept`、`update_dept`、`delete_dept`
- **用户**：`get_user`、`create_user`、`update_user`、`delete_user`
- **消息**：`send_text_message`、`send_image_message`、`send_file_message`、`send_link_message`、`send_sys_message`
- **群组**：`get_group_list`、`get_group_info`、`create_group`、`update_group`、`delete_group`、`add_group_member`、`del_group_member`
- **会话**：`create_session`、`get_session`、`update_session`、`send_text_session_message`、`send_image_session_message`、`send_file_session_message`

### HTTP API 服务器

HTTP API 服务器将所有适配器方法自动暴露为 RESTful API endpoints。

#### 启动 API 服务器

```bash
# 默认端口 8080
./bin/youdu-cli serve-api

# 指定端口
./bin/youdu-cli serve-api --port 9000

# 使用配置文件
./bin/youdu-cli serve-api --config config.yaml --port 8080
```

服务启动后可以访问：
- `GET /health` - 健康检查
- `GET /api/v1/endpoints` - 查看所有可用 API
- `POST /api/v1/*` - 调用各种业务 API

#### API 端点规范

所有业务 API：
- **方法**: `POST`
- **路径格式**: `/api/v1/{method_name}`（snake_case）
- **请求体**: JSON 格式（对应 Input 类型）
- **响应体**: JSON 格式（对应 Output 类型）
- **Content-Type**: `application/json`

#### 使用示例

```bash
# 健康检查
curl http://localhost:8080/health

# 查看所有 API
curl http://localhost:8080/api/v1/endpoints

# 发送文本消息
curl -X POST http://localhost:8080/api/v1/send_text_message \
  -H "Content-Type: application/json" \
  -d '{
    "to_user": "user123",
    "content": "你好，这是一条测试消息"
  }'

# 获取用户信息
curl -X POST http://localhost:8080/api/v1/get_user \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "zhangsan"
  }'

# 创建部门
curl -X POST http://localhost:8080/api/v1/create_dept \
  -H "Content-Type: application/json" \
  -d '{
    "name": "技术部",
    "parent_id": 0,
    "sort_id": 1
  }'

# 获取部门列表
curl -X POST http://localhost:8080/api/v1/get_dept_list \
  -H "Content-Type: application/json" \
  -d '{
    "dept_id": 0
  }'
```

#### 错误响应格式

```json
{
  "error": true,
  "message": "错误详细信息"
}
```

#### 可用的 HTTP API（28 个）

**部门管理**：
- `POST /api/v1/get_dept_list` - 获取部门列表
- `POST /api/v1/get_dept_user_list` - 获取部门用户列表
- `POST /api/v1/get_dept_alias_list` - 获取部门别名列表
- `POST /api/v1/create_dept` - 创建部门
- `POST /api/v1/update_dept` - 更新部门
- `POST /api/v1/delete_dept` - 删除部门

**用户管理**：
- `POST /api/v1/get_user` - 获取用户信息
- `POST /api/v1/create_user` - 创建用户
- `POST /api/v1/update_user` - 更新用户
- `POST /api/v1/delete_user` - 删除用户

**消息管理**：
- `POST /api/v1/send_text_message` - 发送文本消息
- `POST /api/v1/send_image_message` - 发送图片消息
- `POST /api/v1/send_file_message` - 发送文件消息
- `POST /api/v1/send_link_message` - 发送链接消息
- `POST /api/v1/send_sys_message` - 发送系统消息

**群组管理**：
- `POST /api/v1/get_group_list` - 获取群组列表
- `POST /api/v1/get_group_info` - 获取群组信息
- `POST /api/v1/create_group` - 创建群组
- `POST /api/v1/update_group` - 更新群组
- `POST /api/v1/delete_group` - 删除群组
- `POST /api/v1/add_group_member` - 添加群组成员
- `POST /api/v1/del_group_member` - 删除群组成员

**会话管理**：
- `POST /api/v1/create_session` - 创建会话
- `POST /api/v1/get_session` - 获取会话信息
- `POST /api/v1/update_session` - 更新会话
- `POST /api/v1/send_text_session_message` - 发送会话文本消息
- `POST /api/v1/send_image_session_message` - 发送会话图片消息
- `POST /api/v1/send_file_session_message` - 发送会话文件消息

## 项目结构

```
youdu-app-mcp/
├── cmd/
│   ├── youdu-cli/          # CLI 入口
│   └── youdu-mcp/          # MCP 服务器入口
├── internal/
│   ├── adapter/            # 适配器层（核心业务逻辑）
│   │   ├── adapter.go      # 基础适配器
│   │   ├── dept.go         # 部门方法
│   │   ├── user.go         # 用户方法
│   │   ├── message.go      # 消息方法
│   │   ├── group.go        # 群组方法
│   │   └── session.go      # 会话方法
│   ├── api/                # HTTP API 服务器
│   │   └── server.go       # 自动路由注册
│   ├── cli/                # CLI 实现
│   │   ├── root.go         # 根命令
│   │   ├── generator.go    # 自动生成命令
│   │   └── serve_api.go    # API 服务器命令
│   ├── mcp/                # MCP 服务器实现
│   │   └── server.go       # 自动注册工具
│   ├── permission/         # 权限控制
│   │   └── permission.go   # 权限管理系统
│   └── config/             # 配置管理
│       └── config.go       # Viper 配置
├── bin/                    # 编译后的二进制文件
├── config.yaml.example     # 配置示例
├── go.mod
├── go.sum
└── README.md
```

## 开发

### 添加新方法

要添加新的有度 API 方法：

1. 在相应的适配器文件中添加方法（例如 `internal/adapter/dept.go`）
2. 遵循以下模式：
   ```go
   type MethodNameInput struct {
       Field string `json:"field" jsonschema:"description=字段描述,required"`
   }

   type MethodNameOutput struct {
       Result string `json:"result" jsonschema:"description=结果描述"`
   }

   func (a *Adapter) MethodName(ctx context.Context, input MethodNameInput) (*MethodNameOutput, error) {
       // 实现代码
   }
   ```
3. 该方法将自动作为以下形式可用：
   - CLI 命令：`youdu-cli category method-name --field=value`
   - MCP 工具：`method_name`
   - HTTP API：`POST /api/v1/method_name`

### 关键设计原则

1. **单一数据源**：所有 API 在适配器层只定义一次
2. **自动生成**：CLI 命令、MCP 工具和 HTTP endpoints 使用反射自动生成
3. **类型安全**：使用 JSON schema 注解的输入/输出结构体
4. **简洁性**：适配器方法具有简单、直观的名称和参数
5. **统一接口**：三种接口模式（CLI、MCP、HTTP）共享同一套业务逻辑

## 依赖项

- [github.com/modelcontextprotocol/go-sdk](https://github.com/modelcontextprotocol/go-sdk) - 官方 MCP SDK
- [github.com/addcnos/youdu/v2](https://github.com/addcnos/youdu) - 有度 IM SDK
- [github.com/spf13/cobra](https://github.com/spf13/cobra) - CLI 框架
- [github.com/spf13/viper](https://github.com/spf13/viper) - 配置管理
- [github.com/go-chi/chi/v5](https://github.com/go-chi/chi) - 轻量级 HTTP 路由

## 许可证

MIT License

## 贡献

欢迎贡献！请随时提交 Pull Request。
