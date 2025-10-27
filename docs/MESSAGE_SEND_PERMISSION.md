# 消息发送权限控制功能文档

## 概述

消息发送权限控制是一项细粒度的权限管理功能，允许系统管理员精确控制应用可以向哪些用户和部门发送消息。这在需要限制消息发送范围的场景中非常有用。

## 功能特点

- **用户级别限制**：指定允许发送消息的目标用户ID列表
- **部门级别限制**：指定允许发送消息的目标部门ID列表
- **批量发送支持**：支持使用 `|` 分隔符同时向多个用户/部门发送
- **灵活配置**：可以单独配置用户或部门限制，也可以同时配置
- **向后兼容**：如果不配置 `allowsend`，则保持原有行为（允许向任何用户/部门发送）

## 配置方式

### 基本配置格式

在 `config.yaml` 文件中的 `permission.resources.message` 部分添加 `allowsend` 配置：

```yaml
permission:
  enabled: true
  allow_all: false
  
  resources:
    message:
      create: true   # 必须为 true 才能发送消息
      allowsend:
        users: ["10232", "8891"]  # 允许发送消息的目标用户ID列表
        dept: ["1"]               # 允许发送消息的目标部门ID列表
```

### 配置选项说明

#### `allowsend.users`
- **类型**：字符串数组
- **说明**：允许发送消息的目标用户ID列表
- **示例**：`["10232", "8891", "admin001"]`
- **默认值**：空（不配置则不限制）

#### `allowsend.dept`
- **类型**：字符串数组
- **说明**：允许发送消息的目标部门ID列表
- **示例**：`["1", "2", "100"]`
- **默认值**：空（不配置则不限制）

## 配置示例

### 示例 1：只限制用户

```yaml
permission:
  resources:
    message:
      create: true
      allowsend:
        users: ["10232", "8891"]
        # 不配置 dept，表示可以向任何部门发送
```

**效果**：
- ✓ 可以向用户 10232 和 8891 发送消息
- ✗ 不能向其他用户发送消息
- ✓ 可以向任何部门发送消息

### 示例 2：只限制部门

```yaml
permission:
  resources:
    message:
      create: true
      allowsend:
        dept: ["1", "2"]
        # 不配置 users，表示可以向任何用户发送
```

**效果**：
- ✓ 可以向任何用户发送消息
- ✓ 可以向部门 1 和 2 发送消息
- ✗ 不能向其他部门发送消息

### 示例 3：同时限制用户和部门

```yaml
permission:
  resources:
    message:
      create: true
      allowsend:
        users: ["10232", "8891"]
        dept: ["1"]
```

**效果**：
- ✓ 可以向用户 10232 和 8891 发送消息
- ✗ 不能向其他用户发送消息
- ✓ 可以向部门 1 发送消息
- ✗ 不能向其他部门发送消息

### 示例 4：不限制（默认行为）

```yaml
permission:
  resources:
    message:
      create: true
      # 不配置 allowsend
```

**效果**：
- ✓ 可以向任何用户发送消息
- ✓ 可以向任何部门发送消息

## 使用场景

### 场景 1：发送给允许的单个用户

```bash
# 配置：allowsend.users: ["10232", "8891"]
youdu-cli message send-text-message --to-user=10232 --content="Hello"
# 结果：✓ 成功
```

### 场景 2：发送给不允许的用户

```bash
# 配置：allowsend.users: ["10232", "8891"]
youdu-cli message send-text-message --to-user=99999 --content="Hello"
# 结果：✗ 权限拒绝：不允许向用户 '99999' 发送消息
```

### 场景 3：批量发送给多个允许的用户

```bash
# 配置：allowsend.users: ["10232", "8891"]
youdu-cli message send-text-message --to-user="10232|8891" --content="Hello"
# 结果：✓ 成功（所有用户都在允许列表中）
```

### 场景 4：批量发送中包含不允许的用户

```bash
# 配置：allowsend.users: ["10232", "8891"]
youdu-cli message send-text-message --to-user="10232|99999" --content="Hello"
# 结果：✗ 权限拒绝：不允许向用户 '99999' 发送消息
```

### 场景 5：发送给允许的部门

```bash
# 配置：allowsend.dept: ["1"]
youdu-cli message send-text-message --to-dept=1 --content="Hello"
# 结果：✓ 成功
```

### 场景 6：同时发送给用户和部门

```bash
# 配置：allowsend.users: ["10232"], allowsend.dept: ["1"]
youdu-cli message send-text-message --to-user=10232 --to-dept=1 --content="Hello"
# 结果：✓ 成功（用户和部门都在允许列表中）
```

## 支持的消息类型

此权限控制适用于所有类型的消息发送：

1. **文本消息** - `SendTextMessage`
2. **图片消息** - `SendImageMessage`
3. **文件消息** - `SendFileMessage`
4. **链接消息** - `SendLinkMessage`
5. **系统消息** - `SendSysMessage`

## 权限检查流程

```
1. 检查权限系统是否启用
   └─> 如果禁用，允许所有发送
   
2. 检查 message.create 权限
   └─> 如果为 false，拒绝所有发送
   
3. 检查是否配置了 allowsend
   └─> 如果未配置，允许所有发送
   
4. 检查目标用户ID（如果指定了 to_user）
   └─> 验证每个用户ID是否在 allowsend.users 中
   
5. 检查目标部门ID（如果指定了 to_dept）
   └─> 验证每个部门ID是否在 allowsend.dept 中
   
6. 所有检查通过，允许发送
```

## 错误消息

### 未配置消息资源

```
权限拒绝：未配置资源 'message' 的权限策略
```

**解决方法**：在配置文件中添加 `message` 资源配置

### 禁止发送消息

```
权限拒绝：不允许发送消息
```

**解决方法**：将 `permission.resources.message.create` 设置为 `true`

### 用户不在允许列表中

```
权限拒绝：不允许向用户 '99999' 发送消息
```

**解决方法**：将用户ID添加到 `allowsend.users` 列表中，或移除该配置以允许所有用户

### 部门不在允许列表中

```
权限拒绝：不允许向部门 '999' 发送消息
```

**解决方法**：将部门ID添加到 `allowsend.dept` 列表中，或移除该配置以允许所有部门

## API 使用

### CLI 命令

```bash
# 使用配置文件
youdu-cli message send-text-message \
  --config=/path/to/config.yaml \
  --to-user=10232 \
  --content="Hello"
```

### HTTP API

```bash
curl -X POST http://localhost:8080/api/v1/send_text_message \
  -H "Content-Type: application/json" \
  -d '{
    "to_user": "10232",
    "content": "Hello"
  }'
```

### MCP 协议

```json
{
  "method": "tools/call",
  "params": {
    "name": "send_text_message",
    "arguments": {
      "to_user": "10232",
      "content": "Hello"
    }
  }
}
```

## 测试和验证

### 运行演示脚本

```bash
./demo_message_send_permission.sh
```

### 运行单元测试

```bash
# 测试权限逻辑
go test -v ./internal/permission/ -run TestPermission_CheckMessageSend

# 测试适配器集成
go test -v ./internal/adapter/ -run TestAdapter_MessageSendPermission
```

## 注意事项

1. **配置优先级**：`allowsend` 配置优先级高于 `create` 权限。即使 `create: true`，如果配置了 `allowsend` 且目标不在列表中，仍会被拒绝。

2. **空列表行为**：如果 `allowsend.users` 或 `allowsend.dept` 配置为空数组 `[]`，则表示不允许向任何用户/部门发送。建议完全不配置该字段（或删除）以表示不限制。

3. **ID 格式**：用户ID和部门ID都是字符串类型，确保配置时使用引号，例如 `["1", "2"]` 而不是 `[1, 2]`。

4. **批量发送**：使用 `|` 分隔符时，系统会验证每一个ID是否在允许列表中，只要有一个不在列表中就会被拒绝。

5. **会话消息**：会话消息（Session Message）使用会话级别的权限控制，不受 `allowsend` 影响。

## 更新日志

### v1.2.0 (当前版本)

- ✨ 新增消息发送权限控制功能
- ✨ 支持用户级别和部门级别的发送限制
- ✨ 支持批量发送时的权限验证
- 📝 添加完整的测试用例和文档
- 🎉 创建演示配置和脚本

## 常见问题

**Q: 如果同时配置了 `allowsend.users` 和 `allowsend.dept`，发送消息时必须同时指定吗？**

A: 不需要。可以只指定 `to_user` 或只指定 `to_dept`，也可以同时指定。系统会分别验证用户和部门是否在各自的允许列表中。

**Q: 如果想禁止所有消息发送，应该怎么配置？**

A: 将 `permission.resources.message.create` 设置为 `false` 即可。

**Q: `allowsend` 配置是否影响会话消息？**

A: 不影响。会话消息（SendTextSessionMessage 等）使用会话级别的权限控制，与消息发送权限独立。

**Q: 可以使用通配符吗，比如允许所有以 "102" 开头的用户ID？**

A: 当前版本不支持通配符或正则表达式。如需此功能，请提交功能请求。

## 相关文档

- [权限系统总体介绍](../README.md#权限控制)
- [行级权限文档](ROW_LEVEL_PERMISSION.md)
- [配置文件示例](../config.yaml.example)
- [演示脚本](../demo_message_send_permission.sh)
