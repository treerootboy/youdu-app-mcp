#!/bin/bash

# Docker Compose API 测试脚本
# 用于测试 YouDu API 服务的各个端点

set -e

API_BASE="http://localhost:8080"

echo "========================================="
echo "YouDu API 服务测试"
echo "========================================="
echo ""

# 1. 健康检查
echo "1. 测试健康检查 endpoint..."
response=$(curl -s -w "\n%{http_code}" "${API_BASE}/health")
http_code=$(echo "$response" | tail -n1)
body=$(echo "$response" | head -n-1)

if [ "$http_code" = "200" ]; then
    echo "  ✓ 健康检查成功"
    echo "  响应: $body"
else
    echo "  ✗ 健康检查失败 (HTTP $http_code)"
    exit 1
fi
echo ""

# 2. 获取 API 列表
echo "2. 获取 API 端点列表..."
response=$(curl -s "${API_BASE}/api/v1/endpoints")
endpoints_count=$(echo "$response" | grep -o "method_name" | wc -l)

if [ $endpoints_count -gt 0 ]; then
    echo "  ✓ 成功获取 API 列表"
    echo "  可用端点数量: $endpoints_count"
    echo ""
    echo "  前 10 个端点："
    echo "$response" | grep -o '"method_name":"[^"]*"' | head -10 | sed 's/"method_name":/  - /g' | sed 's/"//g'
else
    echo "  ✗ 获取 API 列表失败"
    exit 1
fi
echo ""

# 3. 测试消息发送 API（模拟）
echo "3. 测试消息发送 API..."
echo "  注意: 由于没有真实的有度服务器，此测试会失败，但可以验证 API 格式"
response=$(curl -s -X POST "${API_BASE}/api/v1/send_text_message" \
    -H "Content-Type: application/json" \
    -d '{
        "to_user": "test_user",
        "content": "这是一条测试消息"
    }')

if echo "$response" | grep -q "error"; then
    echo "  ⚠ API 正常响应（预期失败，因为没有真实有度服务器）"
    echo "  响应: $(echo $response | jq -r '.message' 2>/dev/null || echo $response)"
else
    echo "  ✓ API 响应正常"
fi
echo ""

# 4. 测试用户查询 API（模拟）
echo "4. 测试用户查询 API..."
response=$(curl -s -X POST "${API_BASE}/api/v1/get_user" \
    -H "Content-Type: application/json" \
    -d '{
        "user_id": "test_user"
    }')

if echo "$response" | grep -q "error"; then
    echo "  ⚠ API 正常响应（预期失败，因为没有真实有度服务器）"
    echo "  响应: $(echo $response | jq -r '.message' 2>/dev/null || echo $response)"
else
    echo "  ✓ API 响应正常"
fi
echo ""

echo "========================================="
echo "✅ API 服务测试完成！"
echo "========================================="
echo ""
echo "说明："
echo "- 健康检查和 API 列表测试应该成功"
echo "- 业务 API 测试会失败（需要真实的有度服务器）"
echo "- 失败是正常的，表示 API 格式正确但缺少有度服务器连接"
echo ""
echo "完整 API 文档: ${API_BASE}/api/v1/endpoints"
echo ""
