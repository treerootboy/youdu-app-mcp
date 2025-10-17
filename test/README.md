# 测试文档

本目录包含 YouDu IM MCP Server 项目的测试脚本和测试报告。

---

## 目录结构

```
test/
├── README.md              # 本文件
├── reports/              # 测试报告
│   ├── http_api_test.md  # HTTP API 测试报告
│   └── mcp_test.md       # MCP 服务器测试报告
└── scripts/              # 测试脚本
    ├── test_mcp_client.py    # MCP 测试客户端
    └── test_mcp.sh           # MCP Bash 测试脚本
```

---

## 测试脚本说明

### 1. MCP 测试客户端 (Python)

**文件**: `scripts/test_mcp_client.py`

**功能**:
- 完整的 MCP 协议测试客户端
- 自动化测试流程
- 测试工具列表、工具调用、权限控制

**使用方法**:
```bash
cd /path/to/youdu-app-mcp
python3 test/scripts/test_mcp_client.py
```

**测试覆盖**:
- MCP 初始化握手
- 获取工具列表（28 个工具）
- 调用允许的操作（get_user, send_text_message）
- 调用禁止的操作（create_user, delete_user）
- 权限系统验证

**输出示例**:
```
🧪 启动 MCP 服务器测试...

📋 测试 1: 初始化 MCP 连接
✓ 初始化成功

📋 测试 2: 获取工具列表
✓ 获取到 28 个工具
  - add_group_member: add group member
  - create_dept: create dept
  ...

✅ 测试 3: 调用允许的操作 - get_user (权限: read=true)
✓ 成功获取用户信息

❌ 测试 4: 调用被禁止的操作 - create_user (权限: create=false)
✓ 权限控制正常，操作被拒绝: 权限拒绝：不允许对资源 'user' 执行 'create' 操作

🎉 MCP 服务器测试完成！

📊 测试总结:
  ✅ MCP 协议实现正常
  ✅ 工具自动注册成功 (28 个)
  ✅ JSON Schema 正确生成
  ✅ 允许的操作执行成功
  ✅ 权限控制系统正常工作
```

---

### 2. MCP Bash 测试脚本

**文件**: `scripts/test_mcp.sh`

**功能**:
- 使用 Bash + curl 测试 MCP 服务器
- 适合快速验证

**使用方法**:
```bash
cd /path/to/youdu-app-mcp
./test/scripts/test_mcp.sh
```

**注意**: 需要安装 `jq` 工具用于 JSON 解析。

---

## 测试报告说明

### 1. HTTP API 测试报告

**文件**: `reports/http_api_test.md`

**内容**:
- 服务器启动测试
- 路由注册验证（28 个 endpoints）
- 健康检查测试
- API 列表获取测试
- 允许操作测试（发送消息、获取用户）
- 禁止操作测试（创建用户、删除用户）
- 权限系统验证
- 性能测试数据

**测试通过率**: 100% (8/8)

---

### 2. MCP 服务器测试报告

**文件**: `reports/mcp_test.md`

**内容**:
- MCP 协议实现验证
- 初始化握手测试
- 工具列表获取（28 个工具）
- JSON Schema 生成验证
- 工具调用测试
- 权限系统验证
- 与 Claude Desktop 集成说明

**测试通过率**: 100% (8/8)

---

## 测试环境要求

### 系统要求

- macOS / Linux / Windows
- Go 1.23.0+
- Python 3.x (用于 MCP 测试客户端)

### 依赖工具

- `jq` - JSON 命令行处理工具
- `curl` - HTTP 客户端

**安装方法** (macOS):
```bash
brew install jq curl
```

### 配置要求

测试前需要配置：

1. **YouDu 配置** (`config.yaml`):
```yaml
youdu:
  addr: "https://youdu.example.com"
  buin: 123456789
  app_id: "your_app_id"
  aes_key: "your_aes_key"
```

2. **权限配置** (`permission.yaml`):
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
    message:
      create: true
      read: true
      update: false
      delete: false
```

---

## 运行所有测试

### 自动化测试流程

```bash
# 1. 构建项目
go build -o bin/youdu-cli ./cmd/youdu-cli
go build -o bin/youdu-mcp ./cmd/youdu-mcp

# 2. 运行 MCP 测试
python3 test/scripts/test_mcp_client.py

# 3. 运行 HTTP API 测试（手动）
# 启动 API 服务器
./bin/youdu-cli serve-api --port 8080

# 在另一个终端测试
curl http://localhost:8080/health
curl http://localhost:8080/api/v1/endpoints
curl -X POST http://localhost:8080/api/v1/send_text_message \
  -H "Content-Type: application/json" \
  -d '{"to_user":"10232","content":"测试消息"}'
```

---

## 测试覆盖范围

### 已测试

- ✅ MCP 协议基础功能
- ✅ HTTP API 基础功能
- ✅ 工具/路由自动注册
- ✅ JSON Schema 生成
- ✅ 权限控制系统
- ✅ 实际业务操作（发送消息、获取用户）

### 待测试

- ⏳ 所有 28 个工具的详细测试
- ⏳ 并发性能测试
- ⏳ 边界条件测试
- ⏳ 异常情况处理
- ⏳ 与 Claude Desktop 实际集成测试

---

## 贡献测试

欢迎贡献新的测试用例和测试脚本！

### 添加新测试

1. 在 `scripts/` 目录添加测试脚本
2. 在 `reports/` 目录添加测试报告
3. 更新本 README.md

### 测试规范

- 测试脚本应包含清晰的注释
- 测试报告应包含详细的步骤和结果
- 测试应可重复执行
- 测试应有明确的通过/失败标准

---

## 测试结果总结

| 测试类别 | 测试用例数 | 通过 | 失败 | 通过率 |
|---------|----------|------|------|--------|
| HTTP API | 8 | 8 | 0 | 100% |
| MCP 服务器 | 8 | 8 | 0 | 100% |
| **总计** | **16** | **16** | **0** | **100%** |

---

## 问题反馈

如发现测试问题，请提交 Issue 并包含：
- 测试环境信息
- 复现步骤
- 错误日志
- 期望行为

---

**最后更新**: 2025-10-17
**维护者**: YouDu 开发团队
