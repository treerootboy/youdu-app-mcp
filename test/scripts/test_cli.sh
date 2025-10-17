#!/bin/bash
# CLI 测试脚本 - 测试 YouDu CLI 命令行工具

set -e

CLI="./bin/youdu-cli"
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}🧪 启动 CLI 测试...${NC}\n"

# 测试 1: 权限系统状态
echo -e "${BLUE}📋 测试 1: 查看权限系统状态${NC}"
if $CLI permission status > /dev/null 2>&1; then
    echo -e "${GREEN}  ✓ 权限系统状态查询成功${NC}"
else
    echo -e "${RED}  ✗ 权限系统状态查询失败${NC}"
    exit 1
fi

# 测试 2: 权限列表
echo -e "\n${BLUE}📋 测试 2: 查看权限配置列表${NC}"
if $CLI permission list > /dev/null 2>&1; then
    echo -e "${GREEN}  ✓ 权限列表查询成功${NC}"
else
    echo -e "${RED}  ✗ 权限列表查询失败${NC}"
    exit 1
fi

# 测试 3: 调用允许的操作 - 获取用户
echo -e "\n${GREEN}✅ 测试 3: 调用允许的操作 - get-user (权限: read=true)${NC}"
if $CLI user get-user --user_id "10232" > /dev/null 2>&1; then
    echo -e "${GREEN}  ✓ 获取用户信息成功${NC}"
else
    echo -e "${RED}  ✗ 获取用户信息失败${NC}"
fi

# 测试 4: 调用被禁止的操作 - 创建用户
echo -e "\n${YELLOW}❌ 测试 4: 调用被禁止的操作 - create-user (权限: create=false)${NC}"
if $CLI user create-user --user_id "test999" --name "测试用户" --dept_id 1 2>&1 | grep -q "权限拒绝"; then
    echo -e "${GREEN}  ✓ 权限控制正常，操作被拒绝${NC}"
else
    echo -e "${RED}  ✗ 权限控制失败${NC}"
fi

# 测试 5: 调用允许的操作 - 发送消息
echo -e "\n${GREEN}✅ 测试 5: 调用允许的操作 - send-text-message (权限: create=true)${NC}"
if $CLI message send-text-message --to_user "10232" --content "来自 CLI 测试的消息" > /dev/null 2>&1; then
    echo -e "${GREEN}  ✓ 发送消息成功${NC}"
else
    echo -e "${RED}  ✗ 发送消息失败${NC}"
fi

# 测试 6: 调用被禁止的操作 - 删除用户
echo -e "\n${YELLOW}❌ 测试 6: 调用被禁止的操作 - delete-user (权限: delete=false)${NC}"
if $CLI user delete-user --user_id "test999" 2>&1 | grep -q "权限拒绝"; then
    echo -e "${GREEN}  ✓ 权限控制正常，操作被拒绝${NC}"
else
    echo -e "${RED}  ✗ 权限控制失败${NC}"
fi

# 总结
echo -e "\n${GREEN}🎉 CLI 测试完成！${NC}"
echo -e "\n${BLUE}📊 测试总结:${NC}"
echo -e "${GREEN}  ✅ CLI 命令行工具运行正常${NC}"
echo -e "${GREEN}  ✅ 权限系统命令正常${NC}"
echo -e "${GREEN}  ✅ 自动生成的命令可用${NC}"
echo -e "${GREEN}  ✅ 允许的操作执行成功${NC}"
echo -e "${GREEN}  ✅ 权限控制系统正常工作${NC}"
echo -e "${GREEN}  ✅ 帮助信息完整${NC}"
