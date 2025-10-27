# 行级权限（Row-Level Permission）功能说明

## 概述

从 v1.1.0 版本开始，YouDu IM MCP Server 支持行级权限控制（Row-Level Permission），允许管理员通过配置文件精确控制哪些资源 ID 可以被访问。

## 功能特性

- **细粒度控制**：在操作权限（create/read/update/delete）的基础上，进一步限制可访问的具体资源 ID
- **灵活配置**：通过 `allowlist` 配置项实现，支持为不同资源类型配置不同的允许列表
- **向后兼容**：未配置 `allowlist` 时，系统行为与之前版本完全一致
- **安全优先**：行级权限检查在操作权限检查之后进行，确保双重验证

## 配置方法

### 基本配置

在 `config.yaml` 中的资源权限配置节点下添加 `allowlist` 字段：

```yaml
permission:
  enabled: true
  allow_all: false
  
  resources:
    user:
      create: false
      read: true
      update: true
      delete: false
      # 只允许访问这些用户 ID
      allowlist: ["10232", "10023", "user001"]
```

### 配置示例

#### 示例 1：限制用户访问

只允许读取和更新特定用户：

```yaml
resources:
  user:
    create: false
    read: true
    update: true
    delete: false
    allowlist: ["10232", "10023"]
```

**效果**：
- ✅ 可以读取用户 `10232` 和 `10023` 的信息
- ✅ 可以更新用户 `10232` 和 `10023` 的信息
- ❌ 无法读取或更新其他用户（如 `99999`）
- ❌ 无法创建或删除任何用户（操作权限被禁用）

#### 示例 2：限制部门访问

只允许访问特定部门：

```yaml
resources:
  dept:
    create: false
    read: true
    update: false
    delete: false
    allowlist: ["1", "2", "100"]
```

**效果**：
- ✅ 可以读取部门 `1`、`2`、`100` 的信息
- ❌ 无法读取其他部门信息
- ❌ 无法进行创建、更新、删除操作

#### 示例 3：不限制群组访问

群组资源不配置 `allowlist`：

```yaml
resources:
  group:
    create: true
    read: true
    update: true
    delete: false
    # 未配置 allowlist，所有群组都可以访问
```

**效果**：
- ✅ 可以访问所有群组
- 仅受操作权限控制

## 权限检查逻辑

行级权限检查遵循以下流程：

```
1. 检查权限系统是否启用
   ↓ (未启用 → 允许所有操作)
2. 检查 allow_all 配置
   ↓ (allow_all=true → 允许所有操作)
3. 检查资源是否配置权限策略
   ↓ (未配置 → 拒绝)
4. 检查操作权限（create/read/update/delete）
   ↓ (未授权 → 拒绝)
5. 检查是否配置了 allowlist
   ↓ (未配置或为空 → 允许)
6. 检查资源 ID 是否在 allowlist 中
   ↓ (不在列表中 → 拒绝)
7. 允许操作
```

### 权限检查示例

假设配置如下：

```yaml
resources:
  user:
    read: true
    update: true
    allowlist: ["10232", "10023"]
```

各种操作的结果：

| 操作 | 用户ID | 结果 | 原因 |
|------|--------|------|------|
| GetUser | 10232 | ✅ 成功 | ID 在 allowlist 中，且有 read 权限 |
| GetUser | 10023 | ✅ 成功 | ID 在 allowlist 中，且有 read 权限 |
| GetUser | 99999 | ❌ 拒绝 | ID 不在 allowlist 中 |
| UpdateUser | 10232 | ✅ 成功 | ID 在 allowlist 中，且有 update 权限 |
| UpdateUser | 99999 | ❌ 拒绝 | ID 不在 allowlist 中 |
| CreateUser | - | ❌ 拒绝 | 没有 create 权限 |
| DeleteUser | 10232 | ❌ 拒绝 | 没有 delete 权限 |

## 适用范围

行级权限适用于以下资源和操作：

### 用户（User）
- `GetUser` - 获取用户信息 ✅
- `UpdateUser` - 更新用户信息 ✅
- `DeleteUser` - 删除用户 ✅

### 部门（Dept）
- `GetDeptList` - 获取部门列表 ✅
- `GetDeptUserList` - 获取部门用户列表 ✅
- `UpdateDept` - 更新部门信息 ✅
- `DeleteDept` - 删除部门 ✅

### 群组（Group）
- `GetGroupInfo` - 获取群组信息 ✅
- `UpdateGroup` - 更新群组 ✅
- `DeleteGroup` - 删除群组 ✅
- `AddGroupMember` - 添加群组成员 ✅
- `DelGroupMember` - 删除群组成员 ✅

### 会话（Session）
- `GetSession` - 获取会话信息 ✅
- `UpdateSession` - 更新会话 ✅
- `SendTextSessionMessage` - 发送文本会话消息 ✅
- `SendImageSessionMessage` - 发送图片会话消息 ✅
- `SendFileSessionMessage` - 发送文件会话消息 ✅

### 消息（Message）
- 消息操作主要用于发送，不涉及特定消息ID，因此不需要行级权限

## 实现细节

### 代码结构

1. **权限数据结构**（`internal/permission/permission.go`）：
   ```go
   type ResourcePolicy struct {
       Create    bool     `mapstructure:"create"`
       Read      bool     `mapstructure:"read"`
       Update    bool     `mapstructure:"update"`
       Delete    bool     `mapstructure:"delete"`
       AllowList []string `mapstructure:"allowlist"` // 新增
   }
   ```

2. **权限检查方法**（`internal/permission/permission.go`）：
   ```go
   // 原有方法，向后兼容
   func (p *Permission) Check(resource Resource, action Action) error
   
   // 新增方法，支持行级权限
   func (p *Permission) CheckWithID(resource Resource, action Action, resourceID string) error
   ```

3. **适配器方法**（`internal/adapter/adapter.go`）：
   ```go
   // 原有方法
   func (a *Adapter) checkPermission(resource permission.Resource, action permission.Action) error
   
   // 新增方法
   func (a *Adapter) checkPermissionWithID(resource permission.Resource, action permission.Action, resourceID string) error
   ```

### 使用示例

在适配器方法中使用行级权限：

```go
func (a *Adapter) GetUser(ctx context.Context, input GetUserInput) (*GetUserOutput, error) {
    // 使用行级权限检查
    if err := a.checkPermissionWithID(permission.ResourceUser, permission.ActionRead, input.UserID); err != nil {
        return nil, err
    }
    
    // 执行实际操作
    resp, err := a.client.GetUser(ctx, input.UserID)
    // ...
}
```

## 测试

项目包含完整的测试用例：

1. **单元测试**（`internal/permission/permission_test.go`）：
   - 测试 allowlist 为空时的行为
   - 测试 ID 在 allowlist 中的情况
   - 测试 ID 不在 allowlist 中的情况
   - 测试向后兼容性

2. **用户资源集成测试**（`internal/adapter/row_permission_test.go`）：
   - 测试未配置 allowlist 时的行为
   - 测试配置 allowlist 后的访问控制
   - 测试不同操作的 allowlist 限制
   - 测试操作权限与行级权限的交互

3. **所有资源集成测试**（`internal/adapter/all_resources_row_permission_test.go`）：
   - 测试部门资源的行级权限
   - 测试群组资源的行级权限
   - 测试会话资源的行级权限
   - 测试多种资源类型同时配置 allowlist
   - 测试未配置 allowlist 时的行为

运行测试：

```bash
# 运行所有权限相关测试
go test ./internal/permission/... -v

# 运行用户资源行级权限集成测试
go test ./internal/adapter/... -v -run TestAdapter_RowLevelPermissions

# 运行所有资源行级权限集成测试
go test ./internal/adapter/... -v -run TestAdapter_AllResourcesRowLevelPermissions

# 运行所有测试
go test ./...
```

## 故障排查

### 问题：配置了 allowlist 但仍然可以访问所有资源

**可能原因**：
1. `permission.enabled` 设置为 `false`
2. `permission.allow_all` 设置为 `true`
3. 配置文件路径不正确

**解决方法**：
```bash
# 检查权限状态
./bin/youdu-cli permission status

# 查看当前权限配置
./bin/youdu-cli permission list
```

### 问题：明明 ID 在 allowlist 中，但仍然被拒绝

**可能原因**：
1. ID 类型不匹配（字符串 vs 数字）
2. 操作权限未授予
3. 配置文件格式错误

**解决方法**：
1. 确保 allowlist 中的 ID 是字符串格式
2. 检查对应的操作权限（read/update/delete）是否为 true
3. 验证 YAML 配置文件格式正确

### 问题：错误消息"不在允许列表中"

这是正常的行级权限拒绝消息，表示您尝试访问的资源 ID 不在 allowlist 中。

**解决方法**：
- 将需要访问的资源 ID 添加到 allowlist
- 或者移除 allowlist 配置以允许访问所有资源

## 最佳实践

1. **最小权限原则**：只配置必需的资源 ID
   ```yaml
   allowlist: ["10232", "10023"]  # 仅必要的 ID
   ```

2. **分离开发和生产配置**：
   - 开发环境：不配置 allowlist 或设置 `allow_all: true`
   - 生产环境：严格配置 allowlist

3. **定期审查**：定期检查和更新 allowlist，移除不再需要的 ID

4. **文档化**：在配置文件中注释说明每个 ID 的用途
   ```yaml
   allowlist:
     - "10232"  # 管理员账号
     - "10023"  # 测试账号
   ```

5. **日志监控**：监控权限拒绝日志，发现潜在的配置问题

## 版本兼容性

- **向后兼容**：未配置 `allowlist` 时，系统行为与旧版本完全一致
- **平滑升级**：可以逐步为不同资源添加 allowlist 配置
- **配置验证**：系统启动时会验证配置格式，避免运行时错误

## 未来计划

- [ ] 支持 allowlist 正则表达式匹配
- [ ] 支持 denylist（黑名单）
- [ ] 支持基于用户角色的动态 allowlist
- [ ] 提供 CLI 命令管理 allowlist
- [ ] 支持运行时动态更新 allowlist

## 参考

- [权限系统设计文档](./internal/permission/README.md)
- [配置文件示例](./config.yaml.example)
- [测试用例](./internal/permission/permission_test.go)
