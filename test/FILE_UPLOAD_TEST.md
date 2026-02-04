# 文件上传功能测试文档

## 功能概述

本文档描述新增的文件上传功能，包括两个主要方法：

1. **UploadFile**: 上传文件到有度服务器并返回 media_id
2. **SendFileWithUpload**: 一键上传文件并发送文件消息（无需手动获取 media_id）

## 方法详情

### 1. UploadFile - 文件上传

上传文件到有度服务器，返回 media_id 供后续使用。

**输入参数**:
- `file_path` (必填): 要上传的文件路径
- `file_name` (可选): 文件名（含扩展名），如果不提供则从 file_path 自动提取
- `file_type` (可选): 文件类型，可选值：
  - `image` - 图片文件
  - `file` - 普通文件（默认）
  - `voice` - 语音文件
  - `video` - 视频文件

**输出结果**:
- `media_id`: 上传成功后返回的媒体ID
- `success`: 是否上传成功

**权限要求**:
- 需要 `message.create` 权限

### 2. SendFileWithUpload - 一键发送文件

上传文件并立即发送给指定用户或部门，一步到位。

**输入参数**:
- `to_user` (可选): 目标用户ID（多个用户用 | 分隔）
- `to_dept` (可选): 目标部门ID（多个部门用 | 分隔）
- `file_path` (必填): 要上传的文件路径
- `file_name` (可选): 文件名（含扩展名），如果不提供则从 file_path 自动提取
- `file_type` (可选): 文件类型（同 UploadFile）

注意：`to_user` 和 `to_dept` 至少需要填写一个

**输出结果**:
- `media_id`: 上传成功后返回的媒体ID
- `success`: 是否上传并发送成功

**权限要求**:
- 需要消息发送权限（根据 `permission.yaml` 中的 `message.allowsend` 配置）

## 接口测试

### CLI 命令测试

#### 测试 1: 上传文件

```bash
# 创建测试文件
echo "测试文件内容" > /tmp/test_upload.txt

# 上传文件
./bin/youdu-cli upload-file \
  --file_path=/tmp/test_upload.txt \
  --file_name="测试文件.txt" \
  --file_type=file
```

**预期输出**:
```json
{
  "media_id": "xxxxxxxxxxxxxxxxxxxxx",
  "success": true
}
```

#### 测试 2: 一键上传并发送文件

```bash
# 上传并发送文件给用户
./bin/youdu-cli send-file-with-upload \
  --file_path=/tmp/test_upload.txt \
  --file_name="测试文件.txt" \
  --file_type=file \
  --to_user="user123"

# 或者发送给部门
./bin/youdu-cli send-file-with-upload \
  --file_path=/tmp/test_upload.txt \
  --file_name="测试文件.txt" \
  --file_type=file \
  --to_dept="1"
```

**预期输出**:
```json
{
  "media_id": "xxxxxxxxxxxxxxxxxxxxx",
  "success": true
}
```

### HTTP API 测试

#### 测试 1: 上传文件

```bash
curl -X POST http://localhost:8080/api/v1/upload_file \
  -H "Content-Type: application/json" \
  -d '{
    "file_path": "/tmp/test_upload.txt",
    "file_name": "测试文件.txt",
    "file_type": "file"
  }'
```

**预期响应**:
```json
{
  "media_id": "xxxxxxxxxxxxxxxxxxxxx",
  "success": true
}
```

#### 测试 2: 一键上传并发送文件

```bash
curl -X POST http://localhost:8080/api/v1/send_file_with_upload \
  -H "Content-Type: application/json" \
  -d '{
    "file_path": "/tmp/test_upload.txt",
    "file_name": "测试文件.txt",
    "file_type": "file",
    "to_user": "user123"
  }'
```

**预期响应**:
```json
{
  "media_id": "xxxxxxxxxxxxxxxxxxxxx",
  "success": true
}
```

### MCP 工具测试

#### 测试 1: upload_file 工具

使用 MCP 客户端调用 `upload_file` 工具：

```json
{
  "method": "tools/call",
  "params": {
    "name": "upload_file",
    "arguments": {
      "file_path": "/tmp/test_upload.txt",
      "file_name": "测试文件.txt",
      "file_type": "file"
    }
  }
}
```

**预期返回**:
```json
{
  "content": [
    {
      "type": "text",
      "text": "{\"media_id\":\"xxxxxxxxxxxxxxxxxxxxx\",\"success\":true}"
    }
  ]
}
```

#### 测试 2: send_file_with_upload 工具

```json
{
  "method": "tools/call",
  "params": {
    "name": "send_file_with_upload",
    "arguments": {
      "file_path": "/tmp/test_upload.txt",
      "file_name": "测试文件.txt",
      "file_type": "file",
      "to_user": "user123"
    }
  }
}
```

**预期返回**:
```json
{
  "content": [
    {
      "type": "text",
      "text": "{\"media_id\":\"xxxxxxxxxxxxxxxxxxxxx\",\"success\":true}"
    }
  ]
}
```

## 使用场景

### 场景 1: 上传文件后多次使用

当需要将同一个文件发送给多个用户时：

1. 先使用 `UploadFile` 上传文件获取 media_id
2. 然后使用 `SendFileMessage` 多次发送（使用同一个 media_id）

```bash
# 步骤1: 上传文件
MEDIA_ID=$(./bin/youdu-cli upload-file \
  --file_path=/tmp/report.pdf \
  --file_name="月度报告.pdf" | jq -r .media_id)

# 步骤2: 发送给多个用户
./bin/youdu-cli message send-file-message \
  --to_user="user1|user2|user3" \
  --media_id="$MEDIA_ID"
```

### 场景 2: 一键发送

当只需要发送一次文件时，使用 `SendFileWithUpload` 更简便：

```bash
./bin/youdu-cli send-file-with-upload \
  --file_path=/tmp/report.pdf \
  --file_name="月度报告.pdf" \
  --to_user="user1|user2|user3"
```

## 错误处理

### 常见错误

1. **文件路径不存在**
   ```
   错误: 打开文件失败: open /tmp/nonexistent.txt: no such file or directory
   ```
   解决方法：检查文件路径是否正确

2. **权限不足**
   ```
   错误: 权限不足: 没有消息创建权限
   ```
   解决方法：在 `permission.yaml` 中设置 `message.create: true`

3. **未指定接收者**（SendFileWithUpload）
   ```
   错误: 必须指定接收者：to_user 或 to_dept 至少填写一个
   ```
   解决方法：提供 `to_user` 或 `to_dept` 参数

4. **文件类型错误**
   ```
   错误: 不支持的文件类型
   ```
   解决方法：确保 file_type 为 image/file/voice/video 之一

## 技术实现

### 代码位置

- **适配器方法**: `internal/adapter/message.go`
  - `UploadFile()` - 行 239-283
  - `SendFileWithUpload()` - 行 304-346

### 自动生成的接口

由于项目使用反射自动生成接口，添加上述两个方法后，以下接口自动可用：

1. **CLI 命令**:
   - `youdu-cli upload-file`
   - `youdu-cli send-file-with-upload`

2. **MCP 工具**:
   - `upload_file`
   - `send_file_with_upload`

3. **HTTP API**:
   - `POST /api/v1/upload_file`
   - `POST /api/v1/send_file_with_upload`

### 使用的 YouDu SDK 方法

```go
// 上传文件
resp, err := client.UploadMedia(ctx, youdu.UploadMediaRequest{
    File:     file,          // io.Reader
    FileName: fileName,      // 文件名
    FileType: fileType,      // 文件类型
})

// 返回
type UploadMediaResponse struct {
    MediaID string `json:"mediaId"`
}
```

## 权限配置

在 `permission.yaml` 中配置文件上传权限：

```yaml
permission:
  enabled: true
  allow_all: false
  
  resources:
    message:
      create: true  # 必须为 true 才能上传和发送文件
      read: true
      update: false
      delete: false
      
      # 可选：限制可以发送消息的目标
      allowsend:
        users: ["user1", "user2"]  # 只能发送给这些用户
        dept: ["1", "2"]           # 只能发送给这些部门
```

## 测试清单

- [ ] CLI 命令 `upload-file` 测试
- [ ] CLI 命令 `send-file-with-upload` 测试
- [ ] HTTP API `/api/v1/upload_file` 测试
- [ ] HTTP API `/api/v1/send_file_with_upload` 测试
- [ ] MCP 工具 `upload_file` 测试
- [ ] MCP 工具 `send_file_with_upload` 测试
- [ ] 权限检查测试（enabled=true）
- [ ] 权限检查测试（message.create=false）
- [ ] 文件类型测试（image, file, voice, video）
- [ ] 错误处理测试（文件不存在、无权限等）
- [ ] media_id 返回值验证

## 更新日志

### 2026-02-04
- ✅ 新增 `UploadFile` 方法
- ✅ 新增 `SendFileWithUpload` 方法
- ✅ CLI、MCP、HTTP API 自动生成
- ✅ 添加权限检查
- ✅ 返回 media_id 在响应中
