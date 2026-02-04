# Docker 部署指南

本文档详细介绍如何使用 Docker 和 Docker Compose 一键启动 YouDu MCP Service。

## 目录

- [快速开始](#快速开始)
- [架构说明](#架构说明)
- [详细配置](#详细配置)
- [数据持久化](#数据持久化)
- [常用操作](#常用操作)
- [故障排查](#故障排查)
- [高级配置](#高级配置)

## 快速开始

### 前置要求

- Docker 20.10 或更高版本
- Docker Compose 2.0 或更高版本
- 可访问的有度 IM 服务器

### 一键启动

```bash
# 1. 克隆仓库
git clone https://github.com/yourusername/youdu-app-mcp.git
cd youdu-app-mcp

# 2. 运行启动脚本
./start.sh
```

启动脚本会自动：
1. 检查 Docker 和 Docker Compose 是否安装
2. 创建必要的目录（`data/`、`config/`）
3. 引导您配置 `.env` 文件
4. 构建 Docker 镜像
5. 启动服务
6. 显示访问信息

### 手动启动

如果您想手动控制每一步：

```bash
# 1. 复制并编辑配置文件
cp .env.example .env
# 编辑 .env 文件，填入您的有度服务器配置

# 2. 创建必要的目录
mkdir -p data config

# 3. 复制配置文件
cp config.yaml.example config/config.yaml
# 可选：编辑 config/config.yaml 进行详细配置

# 4. 构建并启动服务
docker compose up -d
```

## 架构说明

### 服务组件

Docker Compose 启动两个服务：

1. **youdu-api** - HTTP API 服务
   - 端口：8080（可配置）
   - 用途：供 LLM 通过 HTTP 协议调用
   - 健康检查：`GET /health`
   - API 文档：`GET /api/v1/endpoints`

2. **youdu-mcp** - MCP 服务器
   - 端口：3000（可配置）
   - 用途：供 Claude Desktop 等 MCP 客户端调用
   - 协议：Model Context Protocol (MCP)

### 数据持久化

所有数据通过 Docker 卷映射到宿主机：

```
./data/        → /app/data        # SQLite 数据库
./config/      → /app/config      # 配置文件
```

这确保了：
- 容器重启后数据不丢失
- Token 和权限配置持久化
- 可以直接在宿主机上备份数据

### 网络架构

```
┌─────────────────────────────────────────────────┐
│                   宿主机                        │
│                                                 │
│  ┌──────────────────────────────────────────┐  │
│  │         Docker Network (bridge)          │  │
│  │                                          │  │
│  │  ┌──────────────┐    ┌──────────────┐  │  │
│  │  │  youdu-api   │    │  youdu-mcp   │  │  │
│  │  │  (HTTP API)  │    │  (MCP 服务器)│  │  │
│  │  │  Port: 8080  │    │  Port: 3000  │  │  │
│  │  └──────┬───────┘    └──────┬───────┘  │  │
│  │         │                    │          │  │
│  │         └────────┬───────────┘          │  │
│  │                  │                      │  │
│  │         ┌────────▼────────┐            │  │
│  │         │  共享数据卷     │            │  │
│  │         │  ./data/        │            │  │
│  │         │  ./config/      │            │  │
│  │         └─────────────────┘            │  │
│  └──────────────────────────────────────────┘  │
│                                                 │
│  对外暴露端口：                                  │
│  - 8080 → HTTP API                             │
│  - 3000 → MCP Server                           │
└─────────────────────────────────────────────────┘
```

## 详细配置

### 环境变量（.env）

`.env` 文件包含服务的基本配置：

```bash
# 有度服务器配置（必填）
YOUDU_ADDR=http://your-youdu-server:7080
YOUDU_BUIN=123456789
YOUDU_APP_ID=your-app-id
YOUDU_AES_KEY=your-aes-key

# 服务端口（可选）
API_PORT=8080
MCP_PORT=3000

# Token 认证（可选）
TOKEN_ENABLED=false

# 权限控制（可选）
PERMISSION_ENABLED=true
PERMISSION_ALLOW_ALL=false
```

### 配置文件（config.yaml）

`config/config.yaml` 提供更详细的配置选项：

```yaml
# 有度服务器配置
youdu:
  addr: "http://your-youdu-server:7080"
  buin: 123456789
  app_id: "your-app-id"
  aes_key: "your-aes-key"

# 数据库配置
db:
  path: "/app/data/youdu.db"

# Token 认证
token:
  enabled: false

# 权限配置
permission:
  enabled: true
  allow_all: false
  resources:
    user:
      create: false
      read: true
      update: false
      delete: false
    message:
      create: true
      read: true
      update: false
      delete: false
    # ... 其他资源配置
```

## 数据持久化

### 数据库文件

SQLite 数据库文件位于：
```
./data/youdu.db
```

包含：
- Token 信息
- 会话数据
- 其他持久化数据

### 备份数据

```bash
# 备份数据库
cp data/youdu.db data/youdu.db.backup

# 或使用 tar 压缩备份
tar -czf youdu-backup-$(date +%Y%m%d).tar.gz data/ config/
```

### 恢复数据

```bash
# 停止服务
docker compose down

# 恢复数据库
cp data/youdu.db.backup data/youdu.db

# 重启服务
docker compose up -d
```

## 常用操作

### 查看日志

```bash
# 查看所有服务日志
docker compose logs -f

# 查看 API 服务日志
docker compose logs -f youdu-api

# 查看 MCP 服务日志
docker compose logs -f youdu-mcp

# 查看最近 100 行日志
docker compose logs --tail=100 youdu-api
```

### 重启服务

```bash
# 重启所有服务
docker compose restart

# 重启单个服务
docker compose restart youdu-api
docker compose restart youdu-mcp

# 重新构建并重启
docker compose up -d --build
```

### 停止服务

```bash
# 停止服务（保留容器）
docker compose stop

# 停止并删除容器
docker compose down

# 停止并删除容器及镜像
docker compose down --rmi all

# 停止并删除所有数据（危险！）
docker compose down -v
```

### 进入容器

```bash
# 进入 API 容器
docker compose exec youdu-api sh

# 进入 MCP 容器
docker compose exec youdu-mcp sh

# 执行 CLI 命令
docker compose exec youdu-api /app/youdu-cli --help
```

### Token 管理

```bash
# 生成新 Token
docker compose exec youdu-api /app/youdu-cli token generate --description "My Token"

# 列出所有 Token
docker compose exec youdu-api /app/youdu-cli token list

# 撤销 Token
docker compose exec youdu-api /app/youdu-cli token revoke --id <token-id>
```

### 查看服务状态

```bash
# 查看容器状态
docker compose ps

# 查看容器资源使用
docker stats

# 查看容器详细信息
docker compose inspect youdu-api
```

## 故障排查

### 服务无法启动

1. **检查 Docker 是否运行**
   ```bash
   docker info
   ```

2. **检查端口是否被占用**
   ```bash
   # Linux/Mac
   lsof -i :8080
   lsof -i :3000
   
   # Windows
   netstat -ano | findstr :8080
   netstat -ano | findstr :3000
   ```

3. **查看容器日志**
   ```bash
   docker compose logs
   ```

### 配置问题

1. **检查 .env 文件**
   ```bash
   cat .env
   ```
   确保所有必填字段都已填写

2. **检查 config.yaml**
   ```bash
   cat config/config.yaml
   ```
   确保 YAML 格式正确

3. **验证配置**
   ```bash
   docker compose config
   ```

### API 无法访问

1. **检查服务是否运行**
   ```bash
   docker compose ps
   ```

2. **测试健康检查**
   ```bash
   curl http://localhost:8080/health
   ```

3. **查看 API 日志**
   ```bash
   docker compose logs youdu-api
   ```

### 权限问题

1. **检查数据目录权限**
   ```bash
   ls -la data/
   ```
   确保容器用户有读写权限

2. **修复权限**
   ```bash
   sudo chown -R 1000:1000 data/
   ```

### 数据库问题

1. **检查数据库文件**
   ```bash
   ls -lh data/youdu.db
   ```

2. **重置数据库**（警告：会删除所有数据）
   ```bash
   docker compose down
   rm data/youdu.db
   docker compose up -d
   ```

### 网络问题

1. **检查 Docker 网络**
   ```bash
   docker network ls
   docker network inspect youdu-app-mcp_youdu-network
   ```

2. **重建网络**
   ```bash
   docker compose down
   docker network prune
   docker compose up -d
   ```

## 高级配置

### 自定义端口

在 `.env` 文件中修改：
```bash
API_PORT=9000
MCP_PORT=4000
```

或在启动时指定：
```bash
API_PORT=9000 MCP_PORT=4000 docker compose up -d
```

### 启用 Token 认证

1. 在 `.env` 中启用：
   ```bash
   TOKEN_ENABLED=true
   ```

2. 在 `config/config.yaml` 中启用：
   ```yaml
   token:
     enabled: true
   ```

3. 重启服务：
   ```bash
   docker compose restart
   ```

4. 生成 Token：
   ```bash
   docker compose exec youdu-api /app/youdu-cli token generate --description "API Token"
   ```

### 配置权限

编辑 `config/config.yaml`：

```yaml
permission:
  enabled: true
  allow_all: false
  resources:
    user:
      create: false
      read: true
      update: false
      delete: false
      # 行级权限：只允许访问指定用户
      allowlist: ["user001", "user002"]
    
    message:
      create: true
      # 消息发送权限：只允许向指定用户发送
      allowsend:
        users: ["user001"]
        dept: ["1"]
```

重启服务使配置生效：
```bash
docker compose restart
```

### 使用外部数据库

虽然默认使用 SQLite，但您可以扩展配置以支持其他数据库：

1. 在 `docker-compose.yml` 中添加数据库服务
2. 修改环境变量指向外部数据库
3. 更新配置文件

### 反向代理配置

使用 Nginx 作为反向代理：

```nginx
server {
    listen 80;
    server_name youdu-api.example.com;

    location / {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

### 生产环境部署

生产环境建议：

1. **启用 HTTPS**
   ```yaml
   # docker-compose.yml
   services:
     youdu-api:
       environment:
         - HTTPS_ENABLED=true
       volumes:
         - ./certs:/app/certs
   ```

2. **启用 Token 认证**
   ```yaml
   token:
     enabled: true
   ```

3. **限制权限**
   ```yaml
   permission:
     enabled: true
     allow_all: false
   ```

4. **配置日志级别**
   ```yaml
   logging:
     level: info
   ```

5. **定期备份数据**
   ```bash
   # 添加到 crontab
   0 2 * * * cd /path/to/youdu-app-mcp && tar -czf backup/youdu-$(date +\%Y\%m\%d).tar.gz data/ config/
   ```

### 监控和告警

使用 Docker 自带的监控：

```bash
# 实时监控资源使用
docker stats

# 设置资源限制
# docker-compose.yml
services:
  youdu-api:
    deploy:
      resources:
        limits:
          cpus: '1.0'
          memory: 512M
        reservations:
          cpus: '0.5'
          memory: 256M
```

### 多实例部署

需要负载均衡时：

```yaml
# docker-compose.yml
services:
  youdu-api-1:
    <<: *api-service
    container_name: youdu-api-1
    ports:
      - "8081:8080"
  
  youdu-api-2:
    <<: *api-service
    container_name: youdu-api-2
    ports:
      - "8082:8080"
  
  nginx:
    image: nginx:alpine
    ports:
      - "80:80"
    volumes:
      - ./nginx.conf:/etc/nginx/nginx.conf
    depends_on:
      - youdu-api-1
      - youdu-api-2
```

## 更新升级

### 拉取最新代码

```bash
# 停止服务
docker compose down

# 拉取最新代码
git pull origin main

# 重新构建并启动
docker compose up -d --build
```

### 版本管理

```bash
# 查看当前版本
docker images | grep youdu

# 标记版本
docker tag youdu-app-mcp:latest youdu-app-mcp:v1.0.0

# 清理旧镜像
docker image prune -a
```

## 常见问题 (FAQ)

### Q: 如何修改端口？
A: 在 `.env` 文件中修改 `API_PORT` 和 `MCP_PORT`，然后重启服务。

### Q: 数据会丢失吗？
A: 不会。数据存储在宿主机的 `./data/` 目录中，容器重启不影响数据。

### Q: 如何配置权限？
A: 编辑 `config/config.yaml` 中的 `permission` 部分，重启服务生效。

### Q: 能否只启动 API 服务？
A: 可以。使用 `docker compose up -d youdu-api`

### Q: 如何查看 API 文档？
A: 访问 `http://localhost:8080/api/v1/endpoints`

### Q: 支持 HTTPS 吗？
A: 建议使用 Nginx 等反向代理配置 HTTPS。

### Q: 如何备份数据？
A: 直接复制 `./data/` 目录即可。

## 相关资源

- [README.md](README.md) - 项目主文档
- [config.yaml.example](config.yaml.example) - 配置示例
- [.env.example](.env.example) - 环境变量示例
- [有度 IM 官网](https://youdu.cn) - 有度即时通讯
- [Docker 官方文档](https://docs.docker.com) - Docker 文档

## 技术支持

如有问题，请：
1. 查看 [故障排查](#故障排查) 部分
2. 查看 Docker 日志：`docker compose logs`
3. 提交 Issue 到 GitHub
4. 联系有度技术支持

---

**最后更新**: 2026-02-04
**维护者**: YouDu 开发团队
