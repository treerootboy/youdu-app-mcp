#!/bin/bash
# 行级权限功能演示脚本
# 此脚本演示如何使用 allowlist 配置实现行级权限控制

set -e

echo "============================================"
echo "YouDu IM MCP Server - 行级权限功能演示"
echo "============================================"
echo ""

# 创建测试配置文件
cat > /tmp/demo_row_permission.yaml << 'EOF'
youdu:
  addr: "http://localhost:7080"
  buin: 123456789
  app_id: "demo_app"
  aes_key: "dGVzdF9hZXNfa2V5XzMyX2NoYXJhY3RlcnNfbG9uZyE="

db:
  path: "/tmp/demo_youdu.db"

permission:
  enabled: true
  allow_all: false
  
  resources:
    # 用户权限：只允许访问特定用户ID
    user:
      create: false
      read: true
      update: true
      delete: false
      allowlist: ["10232", "10023", "admin001"]
    
    # 部门权限：只允许访问特定部门ID
    dept:
      create: false
      read: true
      update: false
      delete: false
      allowlist: ["1", "2", "100"]
    
    # 群组权限：未配置allowlist，允许访问所有群组
    group:
      create: true
      read: true
      update: true
      delete: false
    
    # 会话权限
    session:
      create: true
      read: true
      update: true
      delete: false
    
    # 消息权限
    message:
      create: true
      read: true
      update: false
      delete: false

token:
  enabled: false
EOF

echo "1. 创建测试配置文件"
echo "   配置文件路径: /tmp/demo_row_permission.yaml"
echo ""

echo "2. 查看权限系统状态"
./bin/youdu-cli permission status --config=/tmp/demo_row_permission.yaml
echo ""

echo "3. 查看资源权限配置（包含 allowlist）"
./bin/youdu-cli permission list --config=/tmp/demo_row_permission.yaml
echo ""

echo "============================================"
echo "配置说明："
echo "============================================"
echo "✓ 用户资源："
echo "  - 只允许读取和更新用户 ID: 10232, 10023, admin001"
echo "  - 尝试访问其他用户将被拒绝"
echo ""
echo "✓ 部门资源："
echo "  - 只允许读取部门 ID: 1, 2, 100"
echo "  - 尝试访问其他部门将被拒绝"
echo ""
echo "✓ 群组资源："
echo "  - 未配置 allowlist"
echo "  - 允许访问所有群组（仅受操作权限限制）"
echo ""
echo "============================================"
echo "测试场景示例："
echo "============================================"
echo ""
echo "场景 1：允许的用户 ID"
echo "  命令: youdu-cli user get-user --user-id=10232"
echo "  结果: ✓ 成功（ID 在 allowlist 中）"
echo ""
echo "场景 2：不允许的用户 ID"
echo "  命令: youdu-cli user get-user --user-id=99999"
echo "  结果: ✗ 权限拒绝（ID 不在 allowlist 中）"
echo ""
echo "场景 3：更新允许的用户"
echo "  命令: youdu-cli user update-user --user-id=10232 --name='New Name'"
echo "  结果: ✓ 成功（ID 在 allowlist 中，且有 update 权限）"
echo ""
echo "场景 4：删除用户（操作权限被禁用）"
echo "  命令: youdu-cli user delete-user --user-id=10232"
echo "  结果: ✗ 权限拒绝（delete 操作被禁用）"
echo ""
echo "============================================"
echo "配置文件示例："
echo "============================================"
cat /tmp/demo_row_permission.yaml
echo ""
echo "============================================"
echo "演示完成！"
echo "============================================"
