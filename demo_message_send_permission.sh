#!/bin/bash
# 消息发送权限功能演示脚本
# 此脚本演示如何使用 allowsend 配置实现消息发送权限控制

set -e

echo "============================================"
echo "YouDu IM MCP Server - 消息发送权限功能演示"
echo "============================================"
echo ""

CONFIG_FILE="/tmp/demo_message_send_permission.yaml"

# 创建测试配置文件
cat > $CONFIG_FILE << 'EOF'
youdu:
  addr: "http://localhost:7080"
  buin: 123456789
  app_id: "demo_app"
  aes_key: "dGVzdF9hZXNfa2V5XzMyX2NoYXJhY3RlcnNfbG9uZyE="

db:
  path: "/tmp/demo_youdu_message.db"

permission:
  enabled: true
  allow_all: false
  
  resources:
    user:
      create: false
      read: true
      update: false
      delete: false
    
    dept:
      create: false
      read: true
      update: false
      delete: false
    
    group:
      create: true
      read: true
      update: true
      delete: false
    
    session:
      create: true
      read: true
      update: true
      delete: false
    
    # 消息权限配置 - 重点演示
    message:
      create: true
      read: true
      update: false
      delete: false
      # 消息发送权限控制
      allowsend:
        users: ["10232", "8891"]  # 只允许向这些用户发送消息
        dept: ["1"]               # 只允许向这些部门发送消息

token:
  enabled: false
EOF

echo "1. 创建测试配置文件"
echo "   配置文件路径: $CONFIG_FILE"
echo ""

echo "2. 查看权限系统状态"
./bin/youdu-cli permission status --config=$CONFIG_FILE
echo ""

echo "3. 查看消息资源权限配置"
./bin/youdu-cli permission list --config=$CONFIG_FILE | grep -A 10 "message"
echo ""

echo "============================================"
echo "配置说明："
echo "============================================"
echo "✓ 消息发送权限："
echo "  - 允许发送消息（create: true）"
echo "  - 只允许向用户 ID: 10232, 8891 发送消息"
echo "  - 只允许向部门 ID: 1 发送消息"
echo "  - 尝试向其他用户或部门发送将被拒绝"
echo ""

echo "============================================"
echo "测试场景示例："
echo "============================================"
echo ""
echo "场景 1：向允许的用户发送消息"
echo "  命令: youdu-cli message send-text-message --to-user=10232 --content='Hello'"
echo "  结果: ✓ 成功（用户ID在 allowsend.users 中）"
echo ""
echo "场景 2：向不允许的用户发送消息"
echo "  命令: youdu-cli message send-text-message --to-user=99999 --content='Hello'"
echo "  结果: ✗ 权限拒绝（用户ID不在 allowsend.users 中）"
echo ""
echo "场景 3：向允许的部门发送消息"
echo "  命令: youdu-cli message send-text-message --to-dept=1 --content='Hello'"
echo "  结果: ✓ 成功（部门ID在 allowsend.dept 中）"
echo ""
echo "场景 4：向不允许的部门发送消息"
echo "  命令: youdu-cli message send-text-message --to-dept=999 --content='Hello'"
echo "  结果: ✗ 权限拒绝（部门ID不在 allowsend.dept 中）"
echo ""
echo "场景 5：同时向多个允许的用户发送"
echo "  命令: youdu-cli message send-text-message --to-user='10232|8891' --content='Hello'"
echo "  结果: ✓ 成功（所有用户ID都在 allowsend.users 中）"
echo ""
echo "场景 6：向部分不允许的用户发送"
echo "  命令: youdu-cli message send-text-message --to-user='10232|99999' --content='Hello'"
echo "  结果: ✗ 权限拒绝（99999 不在 allowsend.users 中）"
echo ""

echo "============================================"
echo "配置文件示例："
echo "============================================"
cat $CONFIG_FILE
echo ""

echo "============================================"
echo "使用说明："
echo "============================================"
echo ""
echo "1. 如果不配置 allowsend，则允许向任何用户/部门发送消息"
echo "   （仅受 create 权限限制）"
echo ""
echo "2. 配置 allowsend.users 后，只能向列表中的用户发送"
echo ""
echo "3. 配置 allowsend.dept 后，只能向列表中的部门发送"
echo ""
echo "4. 可以同时配置 users 和 dept，分别限制用户和部门"
echo ""
echo "5. 支持使用 | 分隔符同时向多个用户/部门发送"
echo ""
echo "6. 所有类型的消息（文本、图片、文件、链接、系统消息）"
echo "   都受相同的权限控制"
echo ""

echo "============================================"
echo "演示完成！"
echo "============================================"
echo ""
echo "提示：使用 --config=$CONFIG_FILE 参数来测试此配置"
