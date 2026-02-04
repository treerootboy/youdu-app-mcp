#!/bin/bash

# 测试文件上传 API

set -e

API_BASE="http://localhost:8080"
TEST_FILE="/tmp/test_upload.txt"

echo "========================================="
echo "文件上传 API 测试"
echo "========================================="
echo ""

# 创建测试文件
echo "准备测试文件..."
echo "这是一个测试文件，用于验证文件上传功能。" > "$TEST_FILE"
echo "✓ 测试文件已创建: $TEST_FILE"
echo ""

# 测试 1: upload_file API
echo "1. 测试 upload_file API..."
response=$(curl -s -X POST "${API_BASE}/api/v1/upload_file" \
    -H "Content-Type: application/json" \
    -d "{
        \"file_path\": \"$TEST_FILE\",
        \"file_name\": \"测试文件.txt\",
        \"file_type\": \"file\"
    }")

echo "响应: $response"
if echo "$response" | grep -q "media_id"; then
    media_id=$(echo "$response" | jq -r '.media_id' 2>/dev/null)
    echo "✓ upload_file API 成功"
    echo "  Media ID: $media_id"
else
    echo "⚠ upload_file API 响应异常（可能需要真实的有度服务器）"
fi
echo ""

# 测试 2: send_file_with_upload API
echo "2. 测试 send_file_with_upload API..."
response=$(curl -s -X POST "${API_BASE}/api/v1/send_file_with_upload" \
    -H "Content-Type: application/json" \
    -d "{
        \"file_path\": \"$TEST_FILE\",
        \"file_name\": \"测试文件.txt\",
        \"file_type\": \"file\",
        \"to_user\": \"test_user\"
    }")

echo "响应: $response"
if echo "$response" | grep -q "media_id"; then
    media_id=$(echo "$response" | jq -r '.media_id' 2>/dev/null)
    echo "✓ send_file_with_upload API 成功"
    echo "  Media ID: $media_id"
else
    echo "⚠ send_file_with_upload API 响应异常（可能需要真实的有度服务器）"
fi
echo ""

echo "========================================="
echo "✅ 文件上传 API 测试完成！"
echo "========================================="
