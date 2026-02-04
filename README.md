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
- **Token 认证**：HTTP API 支持基于 Token 的认证，保护 API 安全

## 安装

### 方式一：Docker 部署（推荐）🐳

使用 Docker Compose 一键启动服务（支持数据持久化）：

```bash
# 克隆仓库
git clone https://github.com/yourusername/youdu-app-mcp.git
cd youdu-app-mcp

# 一键启动
./start.sh
```

启动后可访问：
- HTTP API: http://localhost:8080
- MCP Server: localhost:3000

详细文档请查看 [DOCKER.md](DOCKER.md)

### 方式二：下载预编译二进制文件（推荐）

从 [Releases 页面](https://github.com/treerootboy/youdu-app-mcp/releases)下载适合您平台的预编译二进制文件，无需安装 Go 环境。

下载后添加执行权限（Linux/macOS）：
```bash
chmod +x youdu-cli-*
chmod +x youdu-mcp-*
```

### 方式三：从源码构建

#### 前置要求

- Go 1.23 或更高版本
- 可访问的有度 IM 服务器

#### 构建

```bash
# 构建 MCP 服务器
go build -o bin/youdu-mcp ./cmd/youdu-mcp

# 构建 CLI（包含 HTTP API 服务器）
go build -o bin/youdu-cli ./cmd/youdu-cli
```

## 配置

在项目根目录或 `~/.youdu/config.yaml` 创建 `config.yaml` 文件：

```yaml
# 有度服务器配置
youdu:
  addr: "http://your-youdu-server:7080"
  buin: 123456789
  app_id: "your-app-id"
  aes_key: "your-aes-key"

# 数据库配置（用于 Token 存储）
db:
  path: "./youdu.db"  # SQLite 数据库文件路径

# Token 认证配置
token:
  enabled: false  # 是否启用 token 认证

# 权限配置
permission:
  enabled: true
  allow_all: false
  resources:
    user:
      create: false
      read: true
      update: false
      delete: false
      # 可选：行级权限，只允许访问指定用户ID
      # allowlist: ["10232", "10023"]
```

### 行级权限（AllowList）

从 v1.1.0 开始，支持对资源进行行级权限控制。通过配置 `allowlist`，可以限制只能访问特定 ID 的资源：

```yaml
permission:
  enabled: true
  allow_all: false
  resources:
    user:
      read: true
      update: true
      # 只允许访问这些用户ID
      allowlist: ["10232", "10023", "user001"]
    
    dept:
      read: true
      # 只允许访问这些部门ID
      allowlist: ["1", "2", "100"]
```

**行级权限说明**：
- 当配置了 `allowlist` 时，只有列表中的资源 ID 可以被访问
- 如果未配置 `allowlist` 或列表为空，则不限制资源 ID（仍受操作权限控制）
- 行级权限检查在操作权限检查通过后进行
- **支持所有资源类型**：User、Dept、Group、Session（共 24 个操作方法）

**支持的资源操作**：
- **用户（User）**：GetUser、UpdateUser、DeleteUser
- **部门（Dept）**：GetDeptList、GetDeptUserList、UpdateDept、DeleteDept
- **群组（Group）**：GetGroupInfo、UpdateGroup、DeleteGroup、AddGroupMember、DelGroupMember
- **会话（Session）**：GetSession、UpdateSession、SendTextSessionMessage、SendImageSessionMessage、SendFileSessionMessage

### 消息发送权限（AllowSend）

从 v1.2.0 开始，支持对消息发送进行细粒度的权限控制。通过配置 `allowsend`，可以限制只能向特定用户和部门发送消息：

```yaml
permission:
  resources:
    message:
      create: true
      # 消息发送权限控制
      allowsend:
        users: ["10232", "8891"]  # 只允许向这些用户发送消息
        dept: ["1"]               # 只允许向这些部门发送消息
```

**消息发送权限说明**：
- 可以单独配置 `users` 或 `dept`，也可以同时配置
- 如果不配置 `allowsend`，则允许向任何用户/部门发送消息
- 支持使用 `|` 分隔符同时向多个用户/部门发送
- 适用于所有消息类型：文本、图片、文件、链接、系统消息
- 详细文档请参考：[docs/MESSAGE_SEND_PERMISSION.md](docs/MESSAGE_SEND_PERMISSION.md)

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

# 文件上传和发送
./bin/youdu-cli upload-file --file-path="/path/to/file.pdf" --file-name="文档.pdf"
./bin/youdu-cli send-file-with-upload --file-path="/path/to/file.pdf" --to-user="user123"

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
- **文件**：`upload_file`、`send_file_with_upload`
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

#### Token 认证

HTTP API 支持 Token 认证功能，可以保护 API 不被未授权访问。Token 使用 SQLite 数据库持久化存储。

##### 启用 Token 认证

1. 在 `config.yaml` 中配置数据库和启用 token 认证：

```yaml
# 数据库配置
db:
  path: "./youdu.db"  # SQLite 数据库文件路径

# Token 认证配置
token:
  enabled: true  # 启用 token 认证
```

2. 重启 API 服务器

##### 生成 Token

使用 CLI 命令生成新的 token，自动保存到数据库：

```bash
# 生成永久 token
./bin/youdu-cli token generate --description "Production API Token"

# 生成有过期时间的 token
./bin/youdu-cli token generate --description "Temporary Token" --expires-in 24h

# JSON 格式输出
./bin/youdu-cli token generate --description "Test Token" --json
```

生成的 token 会自动保存到 SQLite 数据库中，无需手动添加到配置文件。

##### 管理 Token

```bash
# 列出所有 token
./bin/youdu-cli token list

# 撤销 token（从数据库中永久删除）
./bin/youdu-cli token revoke --id token001
```

##### 使用 Token 调用 API

在请求中添加 `Authorization` header：

```bash
# 使用 Bearer 格式
curl -X POST http://localhost:8080/api/v1/send_text_message \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer your-token-value" \
  -d '{"to_user": "user123", "content": "Hello"}'

# 或直接使用 token
curl -X POST http://localhost:8080/api/v1/send_text_message \
  -H "Content-Type: application/json" \
  -H "Authorization: your-token-value" \
  -d '{"to_user": "user123", "content": "Hello"}'
```

**注意**：
- 健康检查 (`/health`) 和 API 列表 (`/api/v1/endpoints`) 不需要 token
- 所有业务 API 调用都需要有效的 token
- Token 存储在 SQLite 数据库中，持久化保存
- 修改 token（添加/删除）后无需重启服务器（动态生效）

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
│   │   ├── serve_api.go    # API 服务器命令
│   │   └── token.go        # Token 管理命令
│   ├── mcp/                # MCP 服务器实现
│   │   └── server.go       # 自动注册工具
│   ├── permission/         # 权限控制
│   │   └── permission.go   # 权限管理系统
│   ├── token/              # Token 管理
│   │   └── token.go        # Token 管理器（SQLite 存储）
│   ├── database/           # 数据库管理
│   │   └── database.go     # SQLite 数据库封装
│   └── config/             # 配置管理
│       └── config.go       # Viper 配置
├── bin/                    # 编译后的二进制文件
├── config.yaml.example     # 配置示例
├── youdu.db                # SQLite 数据库（自动创建）
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

## 发布

项目使用 GitHub Actions 自动构建和发布多平台二进制文件。

### 创建新版本

1. **更新版本信息**
   - 更新 CHANGELOG.md，记录新版本的变更内容
   - 确保所有代码已提交并推送到 main 分支

2. **创建版本标签**
   ```bash
   # 创建带注释的标签
   git tag -a v1.0.0 -m "Release v1.0.0"
   
   # 推送标签到远程仓库
   git push origin v1.0.0
   ```

3. **自动构建**
   - 推送标签后，GitHub Actions 会自动触发构建流程
   - 构建过程会为以下平台生成二进制文件：
     - Linux (amd64, arm64)
     - Windows (amd64, arm64)
     - macOS (amd64, arm64)
   - 每个平台会生成两个可执行文件：
     - `youdu-cli-{platform}-{arch}` - CLI 工具（包含 HTTP API 功能）
     - `youdu-mcp-{platform}-{arch}` - MCP 服务器

4. **发布到 GitHub Releases**
   - 构建完成后，所有二进制文件会自动上传到 GitHub Releases
   - Release 会自动生成更新说明
   - 用户可以直接从 Releases 页面下载对应平台的可执行文件

### 下载已发布版本

访问 [Releases 页面](https://github.com/treerootboy/youdu-app-mcp/releases)下载最新版本的二进制文件。

选择适合您平台的文件：
- **Linux 用户**: 下载 `youdu-cli-linux-amd64` 或 `youdu-mcp-linux-amd64`
- **Windows 用户**: 下载 `youdu-cli-windows-amd64.exe` 或 `youdu-mcp-windows-amd64.exe`
- **macOS 用户**: 
  - Intel 芯片: 下载 `youdu-cli-darwin-amd64` 或 `youdu-mcp-darwin-amd64`
  - Apple Silicon (M1/M2): 下载 `youdu-cli-darwin-arm64` 或 `youdu-mcp-darwin-arm64`

下载后需要添加执行权限（Linux/macOS）：
```bash
chmod +x youdu-cli-linux-amd64
chmod +x youdu-mcp-linux-amd64
```

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
