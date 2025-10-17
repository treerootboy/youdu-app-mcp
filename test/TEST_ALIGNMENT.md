# 测试脚本对齐说明

本文档说明三个测试脚本（MCP、HTTP API、CLI）的测试用例完全一致。

## 测试用例对齐

所有三个测试脚本都包含完全相同的 6 个测试用例：

### 测试 1: 初始化/健康检查
- **MCP**: 初始化 MCP 连接
- **HTTP API**: 健康检查 `/health`
- **CLI**: 查看权限系统状态

### 测试 2: 列表查询
- **MCP**: 获取工具列表 `tools/list`
- **HTTP API**: 获取 API 端点列表 `/endpoints`
- **CLI**: 查看权限配置列表

### 测试 3: 允许的操作 - 获取用户 ✅
- **权限**: `user.read = true`
- **MCP**: `get_user` 工具调用
- **HTTP API**: `POST /api/v1/get_user`
- **CLI**: `youdu-cli user get-user --user-id "10232"`
- **预期结果**: 成功获取用户信息

### 测试 4: 被禁止的操作 - 创建用户 ❌
- **权限**: `user.create = false`
- **MCP**: `create_user` 工具调用
- **HTTP API**: `POST /api/v1/create_user`
- **CLI**: `youdu-cli user create-user --user-id "test999" --name "测试用户" --dept-id 1`
- **预期结果**: 权限拒绝

### 测试 5: 允许的操作 - 发送消息 ✅
- **权限**: `message.create = true`
- **MCP**: `send_text_message` 工具调用
- **HTTP API**: `POST /api/v1/send_text_message`
- **CLI**: `youdu-cli message send-text-message --to-user "10232" --content "消息"`
- **预期结果**: 成功发送消息

### 测试 6: 被禁止的操作 - 删除用户 ❌
- **权限**: `user.delete = false`
- **MCP**: `delete_user` 工具调用
- **HTTP API**: `POST /api/v1/delete_user`
- **CLI**: `youdu-cli user delete-user --user-id "test999"`
- **预期结果**: 权限拒绝

## 测试脚本文件

| 接口类型 | 测试脚本文件 | 运行方法 |
|---------|------------|---------|
| MCP | `test/scripts/test_mcp_client.py` | `python3 test/scripts/test_mcp_client.py` |
| HTTP API | `test/scripts/test_http_api.py` | `python3 test/scripts/test_http_api.py` |
| CLI | `test/scripts/test_cli.sh` | `bash test/scripts/test_cli.sh` |

## 前置条件

所有测试脚本都需要：

1. **构建二进制文件**:
   ```bash
   go build -o bin/youdu-cli ./cmd/youdu-cli
   go build -o bin/youdu-mcp ./cmd/youdu-mcp
   ```

2. **配置文件**: `config.yaml` (或 `config_test.yaml` 用于单元测试)
   ```yaml
   youdu:
     addr: "https://youdu.example.com"
     buin: 123456789
     app_id: "your_app_id"
     aes_key: "your_aes_key"
   ```

3. **权限配置**: `permission.yaml`
   ```yaml
   permission:
     enabled: true
     allow_all: false
     resources:
       user:
         create: false  # 测试 4、6 验证
         read: true     # 测试 3 验证
         update: false
         delete: false
       message:
         create: true   # 测试 5 验证
   ```

## 验证结果

所有三个测试脚本应该产生一致的结果：

- ✅ 测试 3 成功（允许读取用户）
- ❌ 测试 4 被拒绝（禁止创建用户）
- ✅ 测试 5 成功（允许发送消息）
- ❌ 测试 6 被拒绝（禁止删除用户）

## 测试覆盖范围

这 6 个测试用例完整验证了：

1. **接口可用性**: 初始化和列表查询成功
2. **读权限**: 允许读取操作（get_user）
3. **写权限**: 允许写入操作（send_text_message）
4. **创建权限拒绝**: 禁止创建操作（create_user）
5. **删除权限拒绝**: 禁止删除操作（delete_user）
6. **权限系统**: 完整的权限控制系统正常工作

---

**创建日期**: 2025-10-17
**最后更新**: 2025-10-17
