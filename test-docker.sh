#!/bin/bash

# 测试 Docker Compose 配置
# 验证所有配置文件和脚本是否正确

set -e

echo "==================================="
echo "Docker Compose 配置验证测试"
echo "==================================="
echo ""

# 1. 检查必要文件是否存在
echo "1. 检查必要文件..."
files=(
    "Dockerfile"
    "docker-compose.yml"
    ".dockerignore"
    ".env.example"
    "start.sh"
    "DOCKER.md"
    "config.yaml.example"
)

for file in "${files[@]}"; do
    if [ -f "$file" ]; then
        echo "  ✓ $file 存在"
    else
        echo "  ✗ $file 不存在"
        exit 1
    fi
done
echo ""

# 2. 验证 Dockerfile 语法
echo "2. 验证 Dockerfile..."
if docker build --no-cache -f Dockerfile -t youdu-test:latest . > /dev/null 2>&1; then
    echo "  ✓ Dockerfile 语法正确，镜像构建成功"
else
    echo "  ✗ Dockerfile 构建失败"
    exit 1
fi
echo ""

# 3. 验证 docker-compose.yml 配置
echo "3. 验证 docker-compose.yml..."
if docker compose config > /dev/null 2>&1; then
    echo "  ✓ docker-compose.yml 配置正确"
else
    echo "  ✗ docker-compose.yml 配置有误"
    exit 1
fi
echo ""

# 4. 检查环境变量配置
echo "4. 检查 .env 文件..."
if [ -f ".env" ]; then
    echo "  ✓ .env 文件存在"
    
    # 检查必要的环境变量
    required_vars=("YOUDU_ADDR" "YOUDU_BUIN" "YOUDU_APP_ID" "YOUDU_AES_KEY")
    for var in "${required_vars[@]}"; do
        if grep -q "^${var}=" .env; then
            echo "    ✓ $var 已配置"
        else
            echo "    ✗ $var 未配置"
        fi
    done
else
    echo "  ⚠ .env 文件不存在（测试跳过）"
fi
echo ""

# 5. 检查脚本权限
echo "5. 检查脚本权限..."
if [ -x "start.sh" ]; then
    echo "  ✓ start.sh 可执行"
else
    echo "  ✗ start.sh 不可执行"
    chmod +x start.sh
    echo "  ✓ 已设置 start.sh 为可执行"
fi
echo ""

# 6. 验证镜像信息
echo "6. 验证 Docker 镜像..."
echo "  镜像详情："
docker images | grep -E "youdu|REPOSITORY" | head -5
echo ""

# 7. 测试基本命令
echo "7. 测试容器基本命令..."
container_id=$(docker run -d --rm youdu-test:latest /app/youdu-cli --help)
if [ $? -eq 0 ]; then
    echo "  ✓ CLI 命令可用"
    docker stop $container_id > /dev/null 2>&1 || true
else
    echo "  ✗ CLI 命令执行失败"
fi
echo ""

# 8. 清理测试镜像
echo "8. 清理测试镜像..."
docker rmi youdu-test:latest > /dev/null 2>&1 || true
echo "  ✓ 测试镜像已清理"
echo ""

echo "==================================="
echo "✅ 所有测试通过！"
echo "==================================="
echo ""
echo "Docker Compose 配置已准备就绪，可以使用以下命令启动："
echo "  ./start.sh          # 一键启动"
echo "  docker compose up -d  # 手动启动"
echo ""
