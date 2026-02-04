#!/bin/bash

# YouDu MCP Service - 一键启动脚本
# 用于 Docker Compose 启动服务，支持数据持久化和配置管理

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 打印带颜色的消息
print_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 打印标题
print_header() {
    echo ""
    echo -e "${BLUE}========================================${NC}"
    echo -e "${BLUE}   YouDu MCP Service - 一键启动${NC}"
    echo -e "${BLUE}========================================${NC}"
    echo ""
}

# 检查命令是否存在
check_command() {
    if ! command -v "$1" &> /dev/null; then
        print_error "$1 未安装"
        return 1
    fi
    return 0
}

# 检查依赖
check_dependencies() {
    print_info "检查依赖..."
    
    local has_error=0
    
    if ! check_command docker; then
        print_error "请先安装 Docker: https://docs.docker.com/get-docker/"
        has_error=1
    else
        print_success "Docker 已安装: $(docker --version)"
    fi
    
    if ! check_command docker-compose && ! docker compose version &> /dev/null; then
        print_error "请先安装 Docker Compose: https://docs.docker.com/compose/install/"
        has_error=1
    else
        if docker compose version &> /dev/null; then
            print_success "Docker Compose 已安装: $(docker compose version)"
        else
            print_success "Docker Compose 已安装: $(docker-compose --version)"
        fi
    fi
    
    if [ $has_error -eq 1 ]; then
        exit 1
    fi
    
    echo ""
}

# 创建必要的目录
create_directories() {
    print_info "创建必要的目录..."
    
    mkdir -p data
    mkdir -p config
    
    print_success "目录创建完成"
    echo ""
}

# 检查并创建配置文件
setup_config() {
    print_info "检查配置文件..."
    
    # 检查 .env 文件
    if [ ! -f .env ]; then
        print_warning ".env 文件不存在，正在创建..."
        
        # 复制示例文件
        cp .env.example .env
        
        print_warning "请编辑 .env 文件并填入您的有度服务器配置："
        print_info "  YOUDU_ADDR=http://your-youdu-server:7080"
        print_info "  YOUDU_BUIN=123456789"
        print_info "  YOUDU_APP_ID=your-app-id"
        print_info "  YOUDU_AES_KEY=your-aes-key"
        echo ""
        
        read -p "是否现在编辑配置文件？(y/n) " -n 1 -r
        echo ""
        if [[ $REPLY =~ ^[Yy]$ ]]; then
            ${EDITOR:-vi} .env
        else
            print_warning "请手动编辑 .env 文件后再次运行此脚本"
            exit 1
        fi
    else
        print_success ".env 文件已存在"
    fi
    
    # 检查 config.yaml
    if [ ! -f config/config.yaml ]; then
        print_info "创建默认配置文件 config/config.yaml..."
        cp config.yaml.example config/config.yaml
        print_success "配置文件已创建"
    else
        print_success "config/config.yaml 已存在"
    fi
    
    echo ""
}

# 构建 Docker 镜像
build_images() {
    print_info "构建 Docker 镜像..."
    
    if docker compose version &> /dev/null; then
        docker compose build
    else
        docker-compose build
    fi
    
    print_success "Docker 镜像构建完成"
    echo ""
}

# 启动服务
start_services() {
    print_info "启动服务..."
    
    if docker compose version &> /dev/null; then
        docker compose up -d
    else
        docker-compose up -d
    fi
    
    print_success "服务启动完成"
    echo ""
}

# 显示服务状态
show_status() {
    print_info "服务状态："
    echo ""
    
    if docker compose version &> /dev/null; then
        docker compose ps
    else
        docker-compose ps
    fi
    
    echo ""
}

# 显示访问信息
show_access_info() {
    # 从 .env 文件读取端口配置
    source .env 2>/dev/null || true
    API_PORT=${API_PORT:-8080}
    MCP_PORT=${MCP_PORT:-3000}
    
    print_success "🎉 YouDu MCP Service 已成功启动！"
    echo ""
    echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo -e "${GREEN}📌 访问信息：${NC}"
    echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo ""
    echo -e "  ${YELLOW}HTTP API 服务：${NC}"
    echo -e "    🌐 地址: http://localhost:$API_PORT"
    echo -e "    💚 健康检查: http://localhost:$API_PORT/health"
    echo -e "    📖 API 文档: http://localhost:$API_PORT/api/v1/endpoints"
    echo ""
    echo -e "  ${YELLOW}MCP 服务器：${NC}"
    echo -e "    🔌 端口: $MCP_PORT"
    echo -e "    📝 用于 Claude Desktop 等 MCP 客户端连接"
    echo ""
    echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo -e "${GREEN}🔧 常用命令：${NC}"
    echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo ""
    echo -e "  查看日志:"
    echo -e "    ${YELLOW}docker compose logs -f${NC}              # 查看所有服务日志"
    echo -e "    ${YELLOW}docker compose logs -f youdu-api${NC}   # 查看 API 服务日志"
    echo -e "    ${YELLOW}docker compose logs -f youdu-mcp${NC}   # 查看 MCP 服务日志"
    echo ""
    echo -e "  停止服务:"
    echo -e "    ${YELLOW}docker compose down${NC}                # 停止并删除容器"
    echo -e "    ${YELLOW}docker compose stop${NC}                # 停止容器（不删除）"
    echo ""
    echo -e "  重启服务:"
    echo -e "    ${YELLOW}docker compose restart${NC}             # 重启所有服务"
    echo -e "    ${YELLOW}docker compose restart youdu-api${NC}   # 重启 API 服务"
    echo ""
    echo -e "  管理 Token:"
    echo -e "    ${YELLOW}docker compose exec youdu-api /app/youdu-cli token generate --description \"My Token\"${NC}"
    echo -e "    ${YELLOW}docker compose exec youdu-api /app/youdu-cli token list${NC}"
    echo ""
    echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo -e "${GREEN}💾 数据持久化：${NC}"
    echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo ""
    echo -e "  数据库: ${YELLOW}./data/youdu.db${NC}"
    echo -e "  配置文件: ${YELLOW}./config/config.yaml${NC}"
    echo ""
    echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo ""
}

# 测试 API 连接
test_api() {
    print_info "等待服务启动（10 秒）..."
    sleep 10
    
    source .env 2>/dev/null || true
    API_PORT=${API_PORT:-8080}
    
    print_info "测试 API 连接..."
    
    if curl -s http://localhost:$API_PORT/health > /dev/null; then
        print_success "API 服务运行正常 ✓"
    else
        print_warning "API 服务可能还在启动中，请稍后访问"
    fi
    
    echo ""
}

# 主函数
main() {
    print_header
    
    # 检查依赖
    check_dependencies
    
    # 创建目录
    create_directories
    
    # 设置配置
    setup_config
    
    # 构建镜像
    build_images
    
    # 启动服务
    start_services
    
    # 显示状态
    show_status
    
    # 测试连接
    test_api
    
    # 显示访问信息
    show_access_info
}

# 运行主函数
main
