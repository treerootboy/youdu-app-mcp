# Claude Code 协作开发规范

本文档为 YouDu IM MCP Server 项目的研发规范和流程指南，适用于人工开发者和 AI 助手（如 Claude Code）协作开发。

---

## 项目概述

**项目名称**: YouDu IM MCP Server
**版本**: v1.0.0
**描述**: 为有度即时通讯提供 CLI、MCP（模型上下文协议）和 HTTP REST API 三种接口的综合服务

**核心特性**:
- 统一适配器层，单一数据源
- 三种接口模式自动生成（CLI、MCP、HTTP API）
- 完整的权限控制系统
- 反射驱动的自动化架构

---

## 架构设计原则

### 1. 单一数据源原则

所有业务逻辑在 `internal/adapter/` 中只定义一次，三种接口自动生成：

```
   CLI       MCP       HTTP API
    │         │         │
    └─────────┼─────────┘
              │
           Adapter (统一业务逻辑)
              │
           YouDu SDK
```

### 2. 反射自动化

- CLI 命令通过反射自动生成 (`internal/cli/generator.go`)
- MCP 工具通过反射自动注册 (`internal/mcp/server.go`)
- HTTP endpoints 通过反射自动映射 (`internal/api/server.go`)

### 3. 类型安全

使用 Go 结构体 + JSON schema 注解实现完全类型安全：

```go
type MethodInput struct {
    Field string `json:"field" jsonschema:"description=字段说明,required"`
}
```

---

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
│   │   ├── permission.go   # 权限管理命令
│   │   └── test.go         # 测试命令
│   ├── mcp/                # MCP 服务器实现
│   │   └── server.go       # 自动注册工具
│   ├── permission/         # 权限控制
│   │   └── permission.go   # 权限管理系统
│   └── config/             # 配置管理
│       └── config.go       # Viper 配置
├── test/                   # 测试文件
│   ├── reports/            # 测试报告
│   └── scripts/            # 测试脚本
├── bin/                    # 编译后的二进制文件
├── config.yaml.example     # 配置文件模板
├── permission.yaml.example # 权限配置模板
├── README.md               # 项目说明
├── CHANGELOG.md            # 更新日志
└── CLAUDE.md              # 本文件
```

---

## 开发工作流

### 添加新功能的标准流程

#### 1. 在 Adapter 层添加方法

**位置**: `internal/adapter/{resource}.go`

**方法签名**:
```go
func (a *Adapter) MethodName(ctx context.Context, input MethodNameInput) (*MethodNameOutput, error)
```

**Input 定义**:
```go
type MethodNameInput struct {
    Field1 string `json:"field1" jsonschema:"description=字段1说明,required"`
    Field2 int    `json:"field2" jsonschema:"description=字段2说明,default=0"`
}
```

**Output 定义**:
```go
type MethodNameOutput struct {
    Result bool `json:"result" jsonschema:"description=操作结果"`
}
```

**方法实现**:
```go
func (a *Adapter) MethodName(ctx context.Context, input MethodNameInput) (*MethodNameOutput, error) {
    // 1. 权限检查
    if err := a.checkPermission(permission.ResourceXxx, permission.ActionXxx); err != nil {
        return nil, err
    }

    // 2. 业务逻辑
    // ...

    // 3. 返回结果
    return &MethodNameOutput{Result: true}, nil
}
```

#### 2. 自动生成接口

无需手动编写，系统会自动生成：

- **CLI 命令**: `youdu-cli {resource} method-name --field1=value`
- **MCP 工具**: `method_name`
- **HTTP API**: `POST /api/v1/method_name`

#### 3. 测试验证

```bash
# CLI 测试
./bin/youdu-cli {resource} method-name --field1=value

# HTTP API 测试
curl -X POST http://localhost:8080/api/v1/method_name \
  -H "Content-Type: application/json" \
  -d '{"field1":"value"}'

# MCP 测试
# 使用 test/scripts/test_mcp_client.py
```

---

## 代码风格规范

### Go 代码规范

1. **命名规范**
   - 导出类型/函数: PascalCase
   - 私有类型/函数: camelCase
   - 常量: PascalCase 或 UPPER_CASE

2. **注释规范**
   - 所有导出函数必须有注释
   - 注释以函数名开头
   - 使用中文注释

3. **错误处理**
   ```go
   if err != nil {
       return nil, fmt.Errorf("操作失败: %w", err)
   }
   ```

4. **结构体 tag**
   ```go
   type Example struct {
       Field string `json:"field" jsonschema:"description=说明,required"`
   }
   ```

### Git 提交规范

遵循 Commitizen 规范：

**格式**:
```
<type>(<scope>): <subject>

<body>

🤖 Generated with [Claude Code](https://claude.com/claude-code)
via [Happy](https://happy.engineering)

Co-Authored-By: Claude <noreply@anthropic.com>
Co-Authored-By: Happy <yesreply@happy.engineering>
```

**Type 类型**:
- `feat`: 新功能
- `fix`: Bug 修复
- `docs`: 文档更新
- `style`: 代码格式
- `refactor`: 重构
- `test`: 测试
- `chore`: 构建/工具

**示例**:
```
feat(api): 添加 HTTP REST API 支持

- 实现反射自动路由注册
- 集成 Chi 框架
- 自动生成 28 个 API endpoints
- 权限系统完美集成

🤖 Generated with [Claude Code](https://claude.com/claude-code)
via [Happy](https://happy.engineering)

Co-Authored-By: Claude <noreply@anthropic.com>
Co-Authored-By: Happy <yesreply@happy.engineering>
```

---

## 权限系统

### 权限配置

**文件**: `permission.yaml`

```yaml
permission:
  enabled: true
  allow_all: false

  resources:
    user:
      create: false
      read: true
      update: false
      delete: false
```

### 权限检查

在 adapter 方法中：

```go
if err := a.checkPermission(permission.ResourceUser, permission.ActionCreate); err != nil {
    return nil, err
}
```

### 资源类型

- `ResourceDept`: 部门
- `ResourceUser`: 用户
- `ResourceGroup`: 群组
- `ResourceSession`: 会话
- `ResourceMessage`: 消息

### 操作类型

- `ActionCreate`: 创建
- `ActionRead`: 读取
- `ActionUpdate`: 更新
- `ActionDelete`: 删除

---

## 测试规范

### 测试文件组织

```
test/
├── reports/
│   ├── http_api_test.md       # HTTP API 测试报告
│   ├── mcp_test.md            # MCP 测试报告
│   └── permission_test.md     # 权限测试报告
└── scripts/
    ├── test_mcp_client.py     # MCP 测试客户端
    └── test_http_api.sh       # HTTP API 测试脚本
```

### 测试覆盖范围

1. **单元测试** (待完善)
   - Adapter 方法测试
   - 权限检查测试
   - 配置加载测试

2. **集成测试**
   - CLI 命令测试
   - HTTP API 端点测试
   - MCP 工具调用测试

3. **权限测试**
   - 允许的操作
   - 禁止的操作
   - 边界条件

---

## 配置管理

### 配置优先级

1. 命令行参数（最��）
2. 环境变量
3. 配置文件
4. 默认值（最低）

### 配置文件

**YouDu 配置**: `config.yaml`
```yaml
youdu:
  addr: "https://youdu.example.com"
  buin: 123456789
  app_id: "your_app_id"
  aes_key: "your_aes_key"
```

**权限配置**: `permission.yaml`
```yaml
permission:
  enabled: true
  allow_all: false
  resources:
    # 资源权限配置
```

### 环境变量

```bash
export YOUDU_ADDR="https://youdu.example.com"
export YOUDU_BUIN=123456789
export YOUDU_APP_ID="your_app_id"
export YOUDU_AES_KEY="your_aes_key"
```

---

## 构建和部署

### 本地构建

```bash
# 构建 CLI（包含 HTTP API）
go build -o bin/youdu-cli ./cmd/youdu-cli

# 构建 MCP 服务器
go build -o bin/youdu-mcp ./cmd/youdu-mcp
```

### 依赖管理

```bash
# 下载依赖
go mod download

# 更新依赖
go get -u ./...

# 整理依赖
go mod tidy
```

### Docker 部署

```bash
# 构建镜像
docker build -t youdu-mcp:latest .

# 运行容器
docker run -d \
  -e YOUDU_ADDR="https://youdu.example.com" \
  -e YOUDU_BUIN=123456789 \
  -e YOUDU_APP_ID="app_id" \
  -e YOUDU_AES_KEY="aes_key" \
  youdu-mcp:latest
```

---

## 版本发布流程

### 1. 准备发布

- [ ] 确保所有测试通过
- [ ] 更新 README.md
- [ ] 编写 CHANGELOG.md
- [ ] 更新版本号

### 2. 创建 Tag

```bash
# 检查未提交文件
git status

# 提交所有更改
git add .
git commit -m "chore: prepare for v1.0.0 release"

# 创建标签
git tag -a v1.0.0 -m "Release v1.0.0

主要功能:
- CLI 命令行工具
- MCP 协议服务器
- HTTP REST API
- 权限控制系统
- 28 个自动生成的接口

🤖 Generated with [Claude Code](https://claude.com/claude-code)
via [Happy](https://happy.engineering)"

# 推送到远程
git push origin main --tags
```

### 3. 发布说明

在 GitHub/GitLab 创建 Release，附带：
- 更新日志
- 二进制文件
- 使用文档

---

## AI 协作指南

### 与 Claude Code 协作

1. **明确任务范围**
   - 具体说明要实现的功能
   - 提供必要的上下文信息

2. **遵循项目规范**
   - 使用反射自动化
   - 保持单一数据源
   - 遵循命名规范

3. **代码审查**
   - AI 生成的代码需要人工审查
   - 特别注意安全性和性能

4. **测试验证**
   - 所有新功能必须测试
   - 提供测试报告

### 常见任务模板

#### 添加新的 Adapter 方法

```
请在 internal/adapter/{resource}.go 中添加一个新方法：
- 方法名: MethodName
- 功能: [描述]
- 输入参数: [列出参数]
- 返回结果: [列出结果]
- 权限要求: ResourceXxx / ActionXxx
```

#### 修复 Bug

```
发现 Bug:
- 位置: [文件:行号]
- 现象: [描述问题]
- 预期行为: [描述期望]
- 复现步骤: [列出步骤]
```

#### 优化性能

```
性能优化需求:
- 目标: [优化目标]
- 当前性能: [性能指标]
- 期望性能: [性能指标]
- 瓶颈分析: [分析结果]
```

---

## 安全考虑

### 1. 配置文件安全

```bash
# 设置正确的文件权限
chmod 600 config.yaml
chmod 600 permission.yaml
```

### 2. 敏感信息管理

- 不要将敏感信息提交到 Git
- 使用环境变量或密钥管理服务
- 提供 `.example` 模板文件

### 3. 权限控制

- 默认拒绝所有操作
- 明确配置允许的操作
- 定期审查权限配置

### 4. API 安全

- HTTPS 通信
- 输入验证
- 错误信息脱敏
- 请求限流（待实现）

---

## 故障排查

### 常见问题

#### 1. 配置加载失败

```bash
# 检查配置文件
cat config.yaml

# 验证配置
./bin/youdu-cli test
```

#### 2. 权限被拒绝

```bash
# 查看权限配置
./bin/youdu-cli permission list

# 检查权限状态
./bin/youdu-cli permission status
```

#### 3. API 调用失败

```bash
# 检查服务器日志
# 验证请求格式
# 确认权限配置
```

---

## 参考资源

### 官方文档

- [有度 IM 官网](https://youdu.cn)
- [有度 Go SDK](https://github.com/addcnos/youdu)
- [Model Context Protocol](https://modelcontextprotocol.io)
- [Claude Code 文档](https://docs.claude.com/claude-code)

### 依赖库文档

- [Cobra CLI 框架](https://github.com/spf13/cobra)
- [Viper 配置库](https://github.com/spf13/viper)
- [Chi HTTP 路由](https://github.com/go-chi/chi)
- [MCP Go SDK](https://github.com/modelcontextprotocol/go-sdk)

---

## 更新历史

| 版本 | 日期 | 说明 |
|------|------|------|
| v1.0.0 | 2025-10-17 | 初始版本，包含完整的开发规范 |

---

## 许可证

MIT License

---

**创建日期**: 2025-10-17
**最后更新**: 2025-10-17
**维护者**: 有度开发团队 + Claude Code
