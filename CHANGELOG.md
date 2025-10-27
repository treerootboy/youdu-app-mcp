# 更新日志

本文档记录 YouDu IM MCP Server 项目的所有重要更新。

格式基于 [Keep a Changelog](https://keepachangelog.com/zh-CN/1.0.0/)，
版本号遵循 [语义化版本](https://semver.org/lang/zh-CN/)。

---

## [v1.1.0] - 2025-10-27

### 新增功能

#### 行级权限控制 (Row-Level Permission)
- ✨ **AllowList 配置**: 支持在资源权限配置中添加 `allowlist` 字段
- ✨ **细粒度控制**: 可以精确控制哪些资源 ID 可以被访问
- ✨ **灵活配置**: 支持为不同资源类型配置不同的允许列表
- ✨ **向后兼容**: 未配置 `allowlist` 时，系统行为与之前版本完全一致

#### API 改进
- ✨ 新增 `CheckWithID` 方法，支持行级权限检查
- ✨ 新增 `checkPermissionWithID` 适配器方法
- ✨ 用户操作（GetUser、UpdateUser、DeleteUser）现在支持行级权限

#### 文档
- 📝 新增行级权限功能文档 (`docs/ROW_LEVEL_PERMISSION.md`)
- 📝 更新 `config.yaml.example` 包含 allowlist 配置示例
- 📝 更新 README 说明行级权限功能

#### 测试
- ✅ 新增行级权限单元测试 (`internal/permission/permission_test.go`)
- ✅ 新增行级权限集成测试 (`internal/adapter/row_permission_test.go`)
- ✅ 新增测试配置文件 (`config_row_permission_test.yaml`)

### 配置示例

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
```

---

## [v1.0.0] - 2025-10-17

### 首次发布 🎉

这是 YouDu IM MCP Server 的第一个正式版本，提供完整的三种接口模式：CLI、MCP 协议和 HTTP REST API。

### 新增功能

#### 核心架构
- ✨ **统一适配器层**: 所有业务逻辑在 `internal/adapter/` 中统一实现
- ✨ **反射自动化**: CLI 命令、MCP 工具和 HTTP endpoints 通过反射自动生成
- ✨ **三种接口模式**: 支持 CLI、MCP 协议和 HTTP REST API
- ✨ **类型安全**: 使用 Go 结构体 + JSON schema 注解实现完全类型安全

#### CLI 命令行工具
- ✨ 实现基于 Cobra 的 CLI 框架
- ✨ 自动生成命令（`internal/cli/generator.go`）
- ✨ 支持部门、用户、消息、群组、会话管理
- ✨ 内置权限管理命令
- ✨ 内置连接测试命令
- ✨ **28 个自动生成的 CLI 命令**

#### MCP 协议服务器
- ✨ 完整的 MCP (Model Context Protocol) 实现
- ✨ JSON-RPC 2.0 协议支持
- ✨ stdio 传输层
- ✨ 自动工具注册（反射驱动）
- ✨ JSON Schema 自动生成
- ✨ **28 个自动注册的 MCP 工具**
- ✨ 可与 Claude Desktop 集成

#### HTTP REST API
- ✨ 基于 Chi 框架的轻量级 HTTP 服务器
- ✨ 自动路由注册（反射驱动）
- ✨ **28 个自动生成的 REST API endpoints**
- ✨ CORS 支持
- ✨ 统一的 JSON 请求/响应格式
- ✨ 健康检查端点 (`GET /health`)
- ✨ API 列表查询端点 (`GET /api/v1/endpoints`)
- ✨ 错误处理中间件
- ✨ 请求日志中间件

#### 权限控制系统
- ✨ 细粒度的资源权限管理
- ✨ 支持 5 种资源类型（dept, user, group, session, message）
- ✨ 支持 4 种操作类型（create, read, update, delete）
- ✨ 基于 YAML 的权限配置文件
- ✨ 权限检查在 adapter 层统一执行
- ✨ 三种接口自动继承权限控制

#### 配置管理
- ✨ 支持 YAML 配置文件
- ✨ 支持环境变量
- ✨ 支持命令行参数
- ✨ 配置优先级: CLI 参数 > 环境变量 > 配置文件 > 默认值
- ✨ 配置验证和错误提示

#### 有度 IM 功能
- ✨ **部门管理** (6 个方法)
  - 获取部门列表
  - 获取部门用户列表
  - 获取部门别名列表
  - 创建部门
  - 更新部门
  - 删除部门

- ✨ **用户管理** (4 个方法)
  - 获取用户信息
  - 创建用户
  - 更新用户
  - 删除用户

- ✨ **消息管理** (5 个方法)
  - 发送文本消息
  - 发送图片消息
  - 发送文件消息
  - 发送链接消息
  - 发送系统消息

- ✨ **群组管理** (7 个方法)
  - 获取群组列表
  - 获取群组信息
  - 创建群组
  - 更新群组
  - 删除群组
  - 添加群组成员
  - 删除群组成员

- ✨ **会话管理** (6 个方法)
  - 创建会话
  - 获取会话信息
  - 更新会话
  - 发送会话文本消息
  - 发送会话图片消息
  - 发送会话文件消息

### 测试
- ✅ HTTP API 测试（通过率 100%）
- ✅ MCP 服务器测试（通过率 100%）
- ✅ 权限系统测试（通过率 100%）
- ✅ 实际业务操作测试（发送消息、获取用户）
- 📝 完整的测试报告和测试脚本

### 文档
- 📖 完整的 README.md
- 📖 开发规范 (CLAUDE.md)
- 📖 配置示例文件
- 📖 权限配置示例
- 📖 HTTP API 测试报告
- 📖 MCP 测试报告
- 📖 测试脚本说明

### 依赖
- Go 1.23.0+
- github.com/addcnos/youdu/v2 v2.6.0
- github.com/modelcontextprotocol/go-sdk v1.0.0
- github.com/spf13/cobra v1.10.1
- github.com/spf13/viper v1.21.0
- github.com/go-chi/chi/v5 v5.2.3

### 技术亮点
- 🚀 **零配置维护**: 新增 adapter 方法自动暴露为三种接口
- 🚀 **单一数据源**: 所有业务逻辑只定义一次
- 🚀 **反射驱动**: CLI、MCP、HTTP API 全部自动生成
- 🚀 **类型安全**: 编译时类型检查
- 🚀 **权限集成**: 统一的权限控制系统
- 🚀 **生产就绪**: 完整的错误处理和日志记录

### 性能
- HTTP API 响应时间: < 20ms
- 支持并发请求
- 轻量级路由框架 (Chi)

### 安全
- 配置文件权限控制
- 敏感信息环境变量存储
- 细粒度权限管理
- HTTPS 支持（通过反向代理）
- 输入验证

### 已知限制
- 暂无批量操作支持
- 暂无请求限流功能
- 暂无认证中间件（计划 v1.1.0）
- 暂无 API 版本控制（计划 v1.1.0）

### 致谢
- 感谢 [有度](https://youdu.cn) 提供优秀的企业 IM 平台
- 感谢 [addcnos/youdu](https://github.com/addcnos/youdu) 提供 Go SDK
- 感谢 [Model Context Protocol](https://modelcontextprotocol.io) 提供优秀的协议规范
- 感谢 Claude Code 协助开发

---

## 版本说明

### 版本号规则

本项目遵循 [语义化版本 2.0.0](https://semver.org/lang/zh-CN/)：

- **主版本号 (MAJOR)**: 不兼容的 API 修改
- **次版本号 (MINOR)**: 向下兼容的功能性新增
- **修订号 (PATCH)**: 向下兼容的问题修正

### 发布周期

- 主版本: 根据需要发布
- 次版本: 每季度发布
- 修订版本: 根据 Bug 修复需要发布

---

## 下一步计划

### v1.1.0 (计划 2025-Q1)

**功能增强**:
- [ ] 添加 API 认证中间件
- [ ] 实现请求限流
- [ ] API 版本控制
- [ ] 批量操作支持
- [ ] WebSocket 实时消息推送

**性能优化**:
- [ ] 连接池
- [ ] 请求缓存
- [ ] 数据库连接优化

**测试完善**:
- [ ] 所有 28 个工具的详细测试
- [ ] 并发性能测试
- [ ] 压力测试
- [ ] 边界条件测试

**文档完善**:
- [ ] API 文档（OpenAPI/Swagger）
- [ ] 架构图
- [ ] 部署指南
- [ ] 故障排查指南

### v1.2.0 (计划 2025-Q2)

**新功能**:
- [ ] 监控和指标（Prometheus）
- [ ] 分布式追踪
- [ ] 日志聚合
- [ ] 健康检查增强

**集成**:
- [ ] Docker Compose 示例
- [ ] Kubernetes 部署模板
- [ ] CI/CD 流水线

---

## 链接

- [GitHub 仓库](https://github.com/yourusername/youdu-app-mcp)
- [问题追踪](https://github.com/yourusername/youdu-app-mcp/issues)
- [拉取请求](https://github.com/yourusername/youdu-app-mcp/pulls)
- [发布页面](https://github.com/yourusername/youdu-app-mcp/releases)

---

**格式**: 遵循 [Keep a Changelog](https://keepachangelog.com/zh-CN/1.0.0/)
**许可证**: MIT License
