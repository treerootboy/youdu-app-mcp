# Token 认证功能使用指南

本文档介绍如何使用 YouDu MCP Server 的 Token 认证功能。

## 功能概述

Token 认证功能为 HTTP API 提供安全访问控制，支持：

1. ✅ 从配置文件加载预定义的 token
2. ✅ 使用 CLI 命令生成新 token
3. ✅ 支持 token 过期时间设置
4. ✅ 支持 Bearer token 和直接 token 两种格式
5. ⏳ 动态重新加载 token（免重启）- 计划中

## 快速开始

### 1. 生成 Token

使用 CLI 命令生成新的 token：

```bash
# 生成永久 token
./bin/youdu-cli token generate --description "Production API Token"

# 生成有过期时间的 token（24小时后过期）
./bin/youdu-cli token generate --description "Temporary Token" --expires-in 24h

# 以 JSON 格式输出
./bin/youdu-cli token generate --description "Test Token" --json
```

输出示例：
```
✅ Token 生成成功！

📋 Token 信息:
  ID:          iqOliDQt34E=
  Value:       y6e5wrCnP1T5SU-R87DchBOlfIx2TJPRAayL8TyLCl4=
  Description: Production API Token
  Created At:  2025-10-23T13:23:59Z
  Expires At:  永不过期

⚠️  请将以下内容添加到 config.yaml 的 token.tokens 列表中:

created_at: "2025-10-23T13:23:59Z"
description: Production API Token
id: iqOliDQt34E=
value: y6e5wrCnP1T5SU-R87DchBOlfIx2TJPRAayL8TyLCl4=
```

### 2. 配置 Token

在 `config.yaml` 中添加生成的 token：

```yaml
token:
  # 启用 token 认证
  enabled: true
  
  # Token 列表
  tokens:
    - id: "iqOliDQt34E="
      value: "y6e5wrCnP1T5SU-R87DchBOlfIx2TJPRAayL8TyLCl4="
      description: "Production API Token"
      created_at: "2025-10-23T13:23:59Z"
    
    - id: "another-id"
      value: "another-token-value"
      description: "Test Token"
      created_at: "2025-10-23T00:00:00Z"
      expires_at: "2025-12-31T23:59:59Z"  # 可选：设置过期时间
```

### 3. 启动 API 服务器

```bash
./bin/youdu-cli serve-api --config config.yaml --port 8080
```

启动输出示例：
```
📋 正在注册 API Endpoints:
  ✓ POST /api/v1/get_user
  ✓ POST /api/v1/send_text_message
  ...

🚀 YouDu API Server 启动在 :8080
📖 API 文档: GET /api/v1/endpoints
💚 健康检查: GET /health
🔒 Token 认证: 已启用
   当前有效 token 数量: 2
```

### 4. 使用 Token 调用 API

#### 使用 Bearer Token 格式

```bash
curl -X POST http://localhost:8080/api/v1/send_text_message \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer y6e5wrCnP1T5SU-R87DchBOlfIx2TJPRAayL8TyLCl4=" \
  -d '{
    "to_user": "user123",
    "content": "Hello, World!"
  }'
```

#### 直接使用 Token

```bash
curl -X POST http://localhost:8080/api/v1/send_text_message \
  -H "Content-Type: application/json" \
  -H "Authorization: y6e5wrCnP1T5SU-R87DchBOlfIx2TJPRAayL8TyLCl4=" \
  -d '{
    "to_user": "user123",
    "content": "Hello, World!"
  }'
```

## Token 管理命令

### 列出所有 Token

```bash
./bin/youdu-cli token list --config config.yaml
```

输出示例：
```
📋 Token 列表 (共 2 个):

ID              Description            Created At            Expires At        Status
---             ---                    ---                   ---               ---
iqOliDQt34E=    Production API Token   2025-10-23 13:23:59   永不过期            ✅ 有效
another-id      Test Token             2025-10-23 00:00:00   2025-12-31 23:59  ✅ 有效
```

### 撤销 Token

```bash
./bin/youdu-cli token revoke --id iqOliDQt34E= --config config.yaml
```

输出：
```
✅ Token iqOliDQt34E= 已撤销

⚠️  请记得从 config.yaml 中删除此 token
```

**注意**：撤销命令只从运行时内存中删除 token，需要手动从配置文件中删除以永久撤销。

## 错误处理

### 缺少 Token

请求：
```bash
curl -X POST http://localhost:8080/api/v1/get_user \
  -H "Content-Type: application/json" \
  -d '{"user_id": "test"}'
```

响应（401 Unauthorized）：
```json
{
  "error": true,
  "message": "缺少 Authorization header"
}
```

### 无效的 Token

请求：
```bash
curl -X POST http://localhost:8080/api/v1/get_user \
  -H "Content-Type: application/json" \
  -H "Authorization: invalid-token" \
  -d '{"user_id": "test"}'
```

响应（401 Unauthorized）：
```json
{
  "error": true,
  "message": "无效的 token"
}
```

### Token 已过期

当使用已过期的 token 时，会收到相同的 "无效的 token" 错误响应。

## 免认证端点

以下端点不需要 token 认证：

- `GET /health` - 健康检查
- `GET /api/v1/endpoints` - API 端点列表

示例：
```bash
# 健康检查（不需要 token）
curl http://localhost:8080/health

# 查看所有 API（不需要 token）
curl http://localhost:8080/api/v1/endpoints
```

## 安全建议

1. **保护 Token 安全**
   - 不要将 token 提交到版本控制系统
   - 使用环境变量或密钥管理服务存储 token
   - 定期轮换 token

2. **使用 HTTPS**
   - 在生产环境中始终使用 HTTPS 保护 API
   - 防止 token 在传输过程中被窃取

3. **设置过期时间**
   - 为临时访问设置 token 过期时间
   - 定期审查和清理过期的 token

4. **最小权限原则**
   - 为不同的服务创建不同的 token
   - 结合权限系统限制 token 可执行的操作

## Token 过期时间格式

支持的时间格式：

- `24h` - 24 小时
- `7d` - 7 天（注意：需要写成 `168h`）
- `30d` - 30 天（注意：需要写成 `720h`）
- `1h30m` - 1 小时 30 分钟
- `2h45m30s` - 2 小时 45 分钟 30 秒

示例：
```bash
# 24 小时后过期
./bin/youdu-cli token generate --description "24h token" --expires-in 24h

# 7 天后过期
./bin/youdu-cli token generate --description "7d token" --expires-in 168h

# 30 天后过期
./bin/youdu-cli token generate --description "30d token" --expires-in 720h
```

## 禁用 Token 认证

如果不需要 token 认证，可以在配置文件中禁用：

```yaml
token:
  enabled: false
```

或者不配置任何 token（如果 token 列表为空，认证会自动禁用）。

## 故障排查

### Token 认证未启用

**问题**：配置了 token 但 API 不需要认证

**原因**：
- `token.enabled` 设置为 `false`
- token 列表为空
- 配置文件未正确加载

**解决方法**：
1. 检查配置文件路径是否正确
2. 确认 `token.enabled: true`
3. 确认至少有一个有效的 token
4. 重启 API 服务器

### Token 总是无效

**问题**：使用正确的 token 仍然收到 "无效的 token" 错误

**可能原因**：
- Token 值复制错误（包含多余的空格或换行）
- Token 已过期
- 配置文件未正确加载

**解决方法**：
1. 使用 `token list` 命令查看配置的 token
2. 检查 token 状态（是否过期）
3. 重新生成并配置 token
4. 确保 token 值完全匹配（区分大小写）

## 下一步计划

未来将添加以下功能：

- [ ] 动态重新加载 token（免重启）
- [ ] Token 使用统计和审计日志
- [ ] 基于 IP 地址的访问控制
- [ ] Token 权限范围（scope）限制
- [ ] gRPC API 的 token 认证支持

---

**更新日期**: 2025-10-23
**版本**: v1.0.0
