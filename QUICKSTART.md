# 快速开始指南

本文档提供 YouDu MCP Service 的快速上手指南。

## 🚀 一键启动（推荐）

### 前置条件

确保已安装：
- Docker 20.10+
- Docker Compose 2.0+

### 启动步骤

```bash
# 1. 克隆仓库
git clone https://github.com/yourusername/youdu-app-mcp.git
cd youdu-app-mcp

# 2. 运行启动脚本
./start.sh
```

脚本会自动：
1. ✅ 检查 Docker 环境
2. ✅ 创建必要目录
3. ✅ 引导配置 .env 文件
4. ✅ 构建 Docker 镜像
5. ✅ 启动服务
6. ✅ 显示访问信息

### 访问服务

启动成功后：

**HTTP API 服务**：
```bash
# 健康检查
curl http://localhost:8080/health

# 查看所有 API
curl http://localhost:8080/api/v1/endpoints

# 发送消息（示例）
curl -X POST http://localhost:8080/api/v1/send_text_message \
  -H "Content-Type: application/json" \
  -d '{
    "to_user": "user123",
    "content": "Hello from API!"
  }'
```

**MCP 服务器**：
- 端口：3000
- 用于 Claude Desktop 等 MCP 客户端

## 📝 配置说明

### 环境变量（.env）

```bash
# 必填：有度服务器配置
YOUDU_ADDR=http://your-youdu-server:7080
YOUDU_BUIN=123456789
YOUDU_APP_ID=your-app-id
YOUDU_AES_KEY=your-aes-key

# 可选：端口配置
API_PORT=8080
MCP_PORT=3000

# 可选：功能开关
TOKEN_ENABLED=false
PERMISSION_ENABLED=true
```

### 详细配置（config/config.yaml）

详细权限和资源配置请参考 `config.yaml.example`

## 🛠️ 常用操作

### 查看日志

```bash
# 所有服务
docker compose logs -f

# API 服务
docker compose logs -f youdu-api

# MCP 服务
docker compose logs -f youdu-mcp
```

### 重启服务

```bash
# 重启所有
docker compose restart

# 重启单个
docker compose restart youdu-api
```

### 停止服务

```bash
# 停止（保留数据）
docker compose down

# 停止并删除所有数据
docker compose down -v
```

### Token 管理

```bash
# 生成 Token
docker compose exec youdu-api /app/youdu-cli token generate --description "My Token"

# 查看 Token
docker compose exec youdu-api /app/youdu-cli token list

# 撤销 Token
docker compose exec youdu-api /app/youdu-cli token revoke --id <token-id>
```

## 🔍 测试验证

### 运行测试脚本

```bash
# 测试 Docker 配置
./test-docker.sh

# 测试 API 端点（需要先启动服务）
./test-api.sh
```

### 手动测试

```bash
# 1. 启动服务
docker compose up -d

# 2. 等待服务启动（约 10 秒）
sleep 10

# 3. 测试健康检查
curl http://localhost:8080/health

# 4. 查看 API 列表
curl http://localhost:8080/api/v1/endpoints | jq

# 5. 查看服务状态
docker compose ps

# 6. 查看日志
docker compose logs --tail=50
```

## 📊 数据持久化

所有数据保存在：

```
./data/youdu.db       # SQLite 数据库
./config/config.yaml  # 配置文件
```

### 备份数据

```bash
# 方式一：直接复制
cp -r data/ data.backup/

# 方式二：打包备份
tar -czf youdu-backup-$(date +%Y%m%d).tar.gz data/ config/
```

### 恢复数据

```bash
# 1. 停止服务
docker compose down

# 2. 恢复数据
cp -r data.backup/* data/

# 3. 重启服务
docker compose up -d
```

## 🔒 安全建议

生产环境部署：

1. **启用 Token 认证**
   ```bash
   TOKEN_ENABLED=true
   ```

2. **配置严格权限**
   编辑 `config/config.yaml`，只允许必要的操作

3. **使用 HTTPS**
   建议通过 Nginx 等反向代理配置 SSL

4. **限制网络访问**
   使用防火墙限制只允许必要的 IP 访问

5. **定期备份**
   设置自动备份任务

## ❓ 常见问题

### Q: 端口被占用怎么办？

A: 在 `.env` 文件中修改端口：
```bash
API_PORT=9000
MCP_PORT=4000
```

### Q: 如何更新到最新版本？

A:
```bash
git pull origin main
docker compose down
docker compose build --no-cache
docker compose up -d
```

### Q: 容器无法启动？

A: 查看日志排查：
```bash
docker compose logs
```

常见原因：
- 配置文件错误
- 端口冲突
- 权限问题

### Q: 如何只启动 API 服务？

A:
```bash
docker compose up -d youdu-api
```

### Q: 数据会丢失吗？

A: 不会。数据持久化在宿主机的 `./data/` 目录中。

## 📚 更多资源

- [完整文档](README.md)
- [Docker 详细指南](DOCKER.md)
- [配置示例](config.yaml.example)
- [有度 IM 官网](https://youdu.cn)

## 🆘 获取帮助

遇到问题？

1. 查看 [故障排查](DOCKER.md#故障排查)
2. 查看日志：`docker compose logs`
3. 提交 Issue 到 GitHub
4. 联系技术支持

---

**快速链接**：

- 🏠 [返回主页](README.md)
- 🐳 [Docker 详细文档](DOCKER.md)
- 📖 [API 文档](http://localhost:8080/api/v1/endpoints)

**最后更新**: 2026-02-04
