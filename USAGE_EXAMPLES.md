# 使用示例

本文档提供 YouDu MCP Service 的实际使用示例。

## 快速开始

### 1. 启动服务

```bash
# 方式一：使用一键启动脚本（推荐）
./start.sh

# 方式二：手动启动
docker compose up -d
```

### 2. 验证服务

```bash
# 检查服务状态
docker compose ps

# 测试健康检查
curl http://localhost:8080/health

# 查看所有 API
curl http://localhost:8080/api/v1/endpoints | jq
```

## HTTP API 使用示例

### 消息管理

#### 发送文本消息

```bash
curl -X POST http://localhost:8080/api/v1/send_text_message \
  -H "Content-Type: application/json" \
  -d '{
    "to_user": "zhangsan",
    "content": "你好，这是一条测试消息"
  }'
```

#### 发送图片消息

```bash
curl -X POST http://localhost:8080/api/v1/send_image_message \
  -H "Content-Type: application/json" \
  -d '{
    "to_user": "zhangsan",
    "image": "base64_encoded_image_data"
  }'
```

#### 发送文件消息

```bash
curl -X POST http://localhost:8080/api/v1/send_file_message \
  -H "Content-Type: application/json" \
  -d '{
    "to_user": "zhangsan",
    "file": "base64_encoded_file_data",
    "file_name": "document.pdf"
  }'
```

#### 发送系统消息

```bash
curl -X POST http://localhost:8080/api/v1/send_sys_message \
  -H "Content-Type: application/json" \
  -d '{
    "to_user": "zhangsan",
    "title": "系统通知",
    "content": "您有一条新消息"
  }'
```

### 用户管理

#### 获取用户信息

```bash
curl -X POST http://localhost:8080/api/v1/get_user \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "zhangsan"
  }'
```

#### 创建用户

```bash
curl -X POST http://localhost:8080/api/v1/create_user \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "lisi",
    "name": "李四",
    "gender": 1,
    "mobile": "13800138000",
    "email": "lisi@example.com",
    "dept_id": 1
  }'
```

#### 更新用户

```bash
curl -X POST http://localhost:8080/api/v1/update_user \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "lisi",
    "name": "李四（更新）",
    "mobile": "13900139000"
  }'
```

### 部门管理

#### 获取部门列表

```bash
curl -X POST http://localhost:8080/api/v1/get_dept_list \
  -H "Content-Type: application/json" \
  -d '{
    "dept_id": 0
  }'
```

#### 获取部门用户列表

```bash
curl -X POST http://localhost:8080/api/v1/get_dept_user_list \
  -H "Content-Type: application/json" \
  -d '{
    "dept_id": 1
  }'
```

#### 创建部门

```bash
curl -X POST http://localhost:8080/api/v1/create_dept \
  -H "Content-Type: application/json" \
  -d '{
    "name": "技术部",
    "parent_id": 0,
    "sort_id": 1
  }'
```

### 群组管理

#### 获取群组列表

```bash
curl -X POST http://localhost:8080/api/v1/get_group_list \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "zhangsan"
  }'
```

#### 创建群组

```bash
curl -X POST http://localhost:8080/api/v1/create_group \
  -H "Content-Type: application/json" \
  -d '{
    "name": "项目讨论组",
    "owner": "zhangsan"
  }'
```

#### 添加群组成员

```bash
curl -X POST http://localhost:8080/api/v1/add_group_member \
  -H "Content-Type: application/json" \
  -d '{
    "group_id": "group123",
    "member": "lisi"
  }'
```

### 会话管理

#### 创建会话

```bash
curl -X POST http://localhost:8080/api/v1/create_session \
  -H "Content-Type: application/json" \
  -d '{
    "title": "项目讨论",
    "creator": "zhangsan",
    "type": "group"
  }'
```

#### 发送会话消息

```bash
curl -X POST http://localhost:8080/api/v1/send_text_session_message \
  -H "Content-Type: application/json" \
  -d '{
    "session_id": "session123",
    "content": "会话消息内容"
  }'
```

## Token 认证示例

### 启用 Token 认证

1. 在 `.env` 中设置：
```bash
TOKEN_ENABLED=true
```

2. 重启服务：
```bash
docker compose restart
```

### 生成 Token

```bash
# 生成永久 Token
docker compose exec youdu-api /app/youdu-cli token generate \
  --description "Production API Token"

# 生成临时 Token（24小时）
docker compose exec youdu-api /app/youdu-cli token generate \
  --description "Temporary Token" \
  --expires-in 24h

# JSON 格式输出
docker compose exec youdu-api /app/youdu-cli token generate \
  --description "API Token" \
  --json
```

### 使用 Token 调用 API

```bash
# 方式一：Bearer 格式
curl -X POST http://localhost:8080/api/v1/send_text_message \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -d '{
    "to_user": "zhangsan",
    "content": "认证消息"
  }'

# 方式二：直接使用 token
curl -X POST http://localhost:8080/api/v1/send_text_message \
  -H "Content-Type: application/json" \
  -H "Authorization: YOUR_TOKEN_HERE" \
  -d '{
    "to_user": "zhangsan",
    "content": "认证消息"
  }'
```

### 管理 Token

```bash
# 列出所有 Token
docker compose exec youdu-api /app/youdu-cli token list

# 撤销 Token
docker compose exec youdu-api /app/youdu-cli token revoke --id TOKEN_ID
```

## Python 使用示例

### 基本调用

```python
import requests
import json

# API 基础 URL
API_BASE = "http://localhost:8080"

# 发送消息
def send_message(to_user, content):
    url = f"{API_BASE}/api/v1/send_text_message"
    payload = {
        "to_user": to_user,
        "content": content
    }
    response = requests.post(url, json=payload)
    return response.json()

# 获取用户信息
def get_user(user_id):
    url = f"{API_BASE}/api/v1/get_user"
    payload = {"user_id": user_id}
    response = requests.post(url, json=payload)
    return response.json()

# 使用示例
result = send_message("zhangsan", "Hello from Python!")
print(result)
```

### 带 Token 认证

```python
import requests

class YouDuClient:
    def __init__(self, base_url, token=None):
        self.base_url = base_url
        self.token = token
        self.headers = {
            "Content-Type": "application/json"
        }
        if token:
            self.headers["Authorization"] = f"Bearer {token}"
    
    def send_message(self, to_user, content):
        url = f"{self.base_url}/api/v1/send_text_message"
        payload = {
            "to_user": to_user,
            "content": content
        }
        response = requests.post(url, json=payload, headers=self.headers)
        return response.json()
    
    def get_user(self, user_id):
        url = f"{self.base_url}/api/v1/get_user"
        payload = {"user_id": user_id}
        response = requests.post(url, json=payload, headers=self.headers)
        return response.json()

# 使用示例
client = YouDuClient("http://localhost:8080", token="YOUR_TOKEN")
result = client.send_message("zhangsan", "Authenticated message")
print(result)
```

## JavaScript/Node.js 示例

```javascript
const axios = require('axios');

const API_BASE = 'http://localhost:8080';

// 发送消息
async function sendMessage(toUser, content) {
  const response = await axios.post(
    `${API_BASE}/api/v1/send_text_message`,
    {
      to_user: toUser,
      content: content
    }
  );
  return response.data;
}

// 获取用户信息
async function getUser(userId) {
  const response = await axios.post(
    `${API_BASE}/api/v1/get_user`,
    {
      user_id: userId
    }
  );
  return response.data;
}

// 带 Token 的客户端
class YouDuClient {
  constructor(baseUrl, token) {
    this.baseUrl = baseUrl;
    this.headers = {
      'Content-Type': 'application/json'
    };
    if (token) {
      this.headers['Authorization'] = `Bearer ${token}`;
    }
  }

  async sendMessage(toUser, content) {
    const response = await axios.post(
      `${this.baseUrl}/api/v1/send_text_message`,
      { to_user: toUser, content: content },
      { headers: this.headers }
    );
    return response.data;
  }

  async getUser(userId) {
    const response = await axios.post(
      `${this.baseUrl}/api/v1/get_user`,
      { user_id: userId },
      { headers: this.headers }
    );
    return response.data;
  }
}

// 使用示例
const client = new YouDuClient('http://localhost:8080', 'YOUR_TOKEN');
client.sendMessage('zhangsan', 'Hello from Node.js!')
  .then(result => console.log(result))
  .catch(error => console.error(error));
```

## 批量操作示例

### 批量发送消息

```bash
#!/bin/bash

# 批量发送消息给多个用户
users=("user1" "user2" "user3")
message="重要通知：请及时查看"

for user in "${users[@]}"; do
  curl -X POST http://localhost:8080/api/v1/send_text_message \
    -H "Content-Type: application/json" \
    -d "{
      \"to_user\": \"$user\",
      \"content\": \"$message\"
    }"
  echo "已发送给 $user"
done
```

### Python 批量操作

```python
import requests
import time

API_BASE = "http://localhost:8080"

def send_batch_messages(users, message):
    """批量发送消息"""
    results = []
    for user in users:
        try:
            response = requests.post(
                f"{API_BASE}/api/v1/send_text_message",
                json={"to_user": user, "content": message}
            )
            results.append({
                "user": user,
                "success": response.status_code == 200,
                "response": response.json()
            })
            time.sleep(0.1)  # 避免请求过快
        except Exception as e:
            results.append({
                "user": user,
                "success": False,
                "error": str(e)
            })
    return results

# 使用示例
users = ["user1", "user2", "user3"]
message = "重要通知"
results = send_batch_messages(users, message)

for result in results:
    print(f"User: {result['user']}, Success: {result['success']}")
```

## 日常维护

### 查看日志

```bash
# 实时查看所有日志
docker compose logs -f

# 查看最近 100 行
docker compose logs --tail=100

# 只看 API 服务日志
docker compose logs -f youdu-api
```

### 重启服务

```bash
# 重启所有服务
docker compose restart

# 重启单个服务
docker compose restart youdu-api

# 重新构建并重启
docker compose up -d --build
```

### 备份数据

```bash
# 备份数据库
cp data/youdu.db "data/youdu.db.backup.$(date +%Y%m%d)"

# 完整备份
tar -czf "youdu-backup-$(date +%Y%m%d-%H%M%S).tar.gz" data/ config/
```

### 清理

```bash
# 停止并删除容器
docker compose down

# 删除镜像
docker compose down --rmi all

# 删除所有数据（危险！）
docker compose down -v
```

## 故障排查

### 检查服务状态

```bash
# 查看容器状态
docker compose ps

# 查看容器详细信息
docker inspect youdu-api

# 查看资源使用
docker stats
```

### 测试连接

```bash
# 测试健康检查
curl -v http://localhost:8080/health

# 测试 API 可用性
curl -v http://localhost:8080/api/v1/endpoints

# 检查端口监听
netstat -tlnp | grep 8080
```

### 进入容器调试

```bash
# 进入容器
docker compose exec youdu-api sh

# 查看进程
ps aux

# 查看环境变量
env | grep YOUDU

# 测试配置
/app/youdu-cli --help
```

## 更多资源

- [快速开始指南](QUICKSTART.md)
- [Docker 详细文档](DOCKER.md)
- [完整 README](README.md)
- [API 文档](http://localhost:8080/api/v1/endpoints)

---

**最后更新**: 2026-02-04
