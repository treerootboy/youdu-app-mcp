#!/bin/bash
# MCP 服务器测试脚本

set -e

MCP_BIN="./bin/youdu-mcp"

echo "🧪 开始测试 MCP 服务器..."
echo ""

# 测试 1: 初始化并获取工具列表
echo "📋 测试 1: 获取工具列表"
TOOLS=$(echo '{"jsonrpc":"2.0","id":1,"method":"tools/list"}' | timeout 5 $MCP_BIN 2>&1)
echo "$TOOLS" | jq -r '.result.tools[] | .name' 2>/dev/null | head -10 || echo "获取工具列表失败"
TOOL_COUNT=$(echo "$TOOLS" | jq -r '.result.tools | length' 2>/dev/null || echo "0")
echo "总计: $TOOL_COUNT 个工具"
echo ""

# 测试 2: 调用允许的操作 - 获取用户信息
echo "✅ 测试 2: 调用允许的操作 - get_user (权限: read=true)"
RESPONSE=$(echo '{"jsonrpc":"2.0","id":2,"method":"tools/call","params":{"name":"get_user","arguments":{"user_id":"10232"}}}' | timeout 5 $MCP_BIN 2>&1)
if echo "$RESPONSE" | jq -e '.result' > /dev/null 2>&1; then
    echo "✓ 成功获取用户信息:"
    echo "$RESPONSE" | jq -r '.result.content[0].text' 2>/dev/null | jq '.user | {userId, name, email}' || echo "$RESPONSE" | jq '.result'
else
    echo "✗ 失败:"
    echo "$RESPONSE" | jq '.error.message' 2>/dev/null || echo "$RESPONSE"
fi
echo ""

# 测试 3: 调用被禁止的操作 - 创建用户
echo "❌ 测试 3: 调用被禁止的操作 - create_user (权限: create=false)"
RESPONSE=$(echo '{"jsonrpc":"2.0","id":3,"method":"tools/call","params":{"name":"create_user","arguments":{"user_id":"test999","name":"测试用户","dept_id":1}}}' | timeout 5 $MCP_BIN 2>&1)
if echo "$RESPONSE" | jq -e '.error' > /dev/null 2>&1; then
    echo "✓ 权限控制正常，操作被拒绝:"
    echo "$RESPONSE" | jq -r '.error.message'
else
    echo "✗ 意外：操作应该被拒绝但成功了"
    echo "$RESPONSE" | jq '.result'
fi
echo ""

# 测试 4: 调用允许的操作 - 发送消息
echo "✅ 测试 4: 调用允许的操作 - send_text_message (权限: create=true)"
RESPONSE=$(echo '{"jsonrpc":"2.0","id":4,"method":"tools/call","params":{"name":"send_text_message","arguments":{"to_user":"10232","content":"MCP 测试消息"}}}' | timeout 5 $MCP_BIN 2>&1)
if echo "$RESPONSE" | jq -e '.result' > /dev/null 2>&1; then
    echo "✓ 消息发送成功:"
    echo "$RESPONSE" | jq -r '.result.content[0].text' 2>/dev/null | jq '.'
else
    echo "✗ 失败:"
    echo "$RESPONSE" | jq '.error.message' 2>/dev/null || echo "$RESPONSE"
fi
echo ""

# 测试 5: 调用被禁止的操作 - 删除用户
echo "❌ 测试 5: 调用被禁止的操作 - delete_user (权限: delete=false)"
RESPONSE=$(echo '{"jsonrpc":"2.0","id":5,"method":"tools/call","params":{"name":"delete_user","arguments":{"user_id":"test999"}}}' | timeout 5 $MCP_BIN 2>&1)
if echo "$RESPONSE" | jq -e '.error' > /dev/null 2>&1; then
    echo "✓ 权限控制正常，操作被拒绝:"
    echo "$RESPONSE" | jq -r '.error.message'
else
    echo "✗ 意外：操作应该被拒绝但成功了"
    echo "$RESPONSE" | jq '.result'
fi
echo ""

echo "🎉 MCP 服务器测试完成！"
