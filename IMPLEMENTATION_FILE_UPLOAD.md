# 文件上传功能实现总结

## 概述

本次更新为 YouDu IM MCP Server 添加了文件上传功能，包括两个核心方法：

1. **UploadFile**: 上传文件到有度服务器并返回 media_id
2. **SendFileWithUpload**: 一键上传文件并发送文件消息（无需手动处理 media_id）

这两个方法遵循项目的架构原则，在适配器层定义一次，自动生成 CLI、MCP 和 HTTP API 三种接口。

## 实现细节

### 1. 代码修改

**文件**: `internal/adapter/message.go`

#### 新增方法

##### UploadFile
```go
func (a *Adapter) UploadFile(ctx context.Context, input UploadFileInput) (*UploadFileOutput, error)
```

**功能**: 
- 上传文件到有度服务器
- 返回 media_id 供后续使用

**参数**:
- `file_path`: 文件路径（必填）
- `file_name`: 文件名，不提供时自动从路径提取
- `file_type`: 文件类型（image/file/voice/video），默认为 file

**返回**:
- `media_id`: 上传成功后的媒体ID
- `success`: 是否成功

##### SendFileWithUpload
```go
func (a *Adapter) SendFileWithUpload(ctx context.Context, input SendFileWithUploadInput) (*SendFileWithUploadOutput, error)
```

**功能**:
- 上传文件并立即发送给指定用户或部门
- 自动处理 media_id，用户无需关心

**参数**:
- `to_user`: 目标用户ID（多个用 | 分隔）
- `to_dept`: 目标部门ID（多个用 | 分隔）
- `file_path`: 文件路径（必填）
- `file_name`: 文件名
- `file_type`: 文件类型

**返回**:
- `media_id`: 上传成功后的媒体ID
- `success`: 是否成功

#### 依赖导入
```go
import (
    "os"
    "path/filepath"
    "github.com/yourusername/youdu-app-mcp/internal/permission"
)
```

### 2. 自动生成的接口

得益于项目的反射自动化架构，添加方法后立即可用以下接口：

#### CLI 命令
```bash
# 上传文件
youdu-cli upload-file --file_path=/path/to/file --file_name="文件名.txt" --file_type=file

# 上传并发送
youdu-cli send-file-with-upload \
  --file_path=/path/to/file \
  --to_user="user123" \
  --file_name="文件名.txt"
```

#### HTTP API
```bash
# 上传文件
POST /api/v1/upload_file
{
  "file_path": "/path/to/file",
  "file_name": "文件名.txt",
  "file_type": "file"
}

# 上传并发送
POST /api/v1/send_file_with_upload
{
  "file_path": "/path/to/file",
  "file_name": "文件名.txt",
  "to_user": "user123"
}
```

#### MCP 工具
```json
// upload_file 工具
{
  "name": "upload_file",
  "arguments": {
    "file_path": "/path/to/file",
    "file_name": "文件名.txt"
  }
}

// send_file_with_upload 工具
{
  "name": "send_file_with_upload",
  "arguments": {
    "file_path": "/path/to/file",
    "to_user": "user123"
  }
}
```

### 3. YouDu SDK 集成

使用 YouDu SDK v2.6.0 的 UploadMedia 方法：

```go
resp, err := a.client.UploadMedia(ctx, youdu.UploadMediaRequest{
    File:     file,          // io.Reader
    FileName: fileName,      // 文件名
    FileType: youdu.FileType(fileType), // 文件类型
})

// 返回
type UploadMediaResponse struct {
    MediaID string `json:"mediaId"`
}
```

支持的文件类型：
- `youdu.FileTypeImage` = "image"
- `youdu.FileTypeFile` = "file"
- `youdu.FileTypeVoice` = "voice"
- `youdu.FileTypeVideo` = "video"

### 4. 权限控制

两个方法都集成了权限检查：

- **UploadFile**: 需要 `message.create` 权限
- **SendFileWithUpload**: 需要消息发送权限（受 `message.allowsend` 配置约束）

在 `permission.yaml` 中配置：
```yaml
permission:
  resources:
    message:
      create: true  # 允许上传文件
      allowsend:
        users: ["user1", "user2"]  # 限制发送目标
```

## 使用场景

### 场景 1: 批量发送同一文件

当需要将同一文件发送给多个用户时，先上传获取 media_id 可以避免重复上传：

```bash
# 1. 上传一次
MEDIA_ID=$(youdu-cli upload-file --file_path=/tmp/report.pdf | jq -r .media_id)

# 2. 多次发送
youdu-cli message send-file-message --media_id="$MEDIA_ID" --to_user="user1|user2|user3"
```

### 场景 2: 一次性发送

当只需发送一次时，使用 SendFileWithUpload 更简便：

```bash
youdu-cli send-file-with-upload \
  --file_path=/tmp/report.pdf \
  --to_user="user1|user2|user3"
```

## 测试

### 测试文件

1. **测试文档**: `test/FILE_UPLOAD_TEST.md`
   - 详细的功能说明
   - CLI、HTTP API、MCP 测试步骤
   - 错误处理说明
   - 权限配置指南

2. **测试脚本**: `test/scripts/test_upload_api.sh`
   - 自动化测试 HTTP API
   - 创建测试文件
   - 验证响应格式

### 测试方法

```bash
# CLI 测试
./bin/youdu-cli upload-file --file_path=/tmp/test.txt

# HTTP API 测试
./test/scripts/test_upload_api.sh

# MCP 测试
# 使用 MCP 客户端或测试工具
```

## 架构优势

本实现完全遵循项目的"单一数据源"原则：

```
         CLI          MCP        HTTP API
          │            │            │
          └────────────┼────────────┘
                       │
                    Adapter
                  (定义一次)
                       │
                  YouDu SDK
```

优势：
1. **代码复用**: 业务逻辑只写一次
2. **自动生成**: CLI、MCP、HTTP API 自动可用
3. **类型安全**: 使用 Go 结构体 + JSON schema
4. **易于维护**: 修改一处，所有接口同步更新

## 技术亮点

1. **智能文件名提取**: 如果不提供 file_name，自动从 file_path 提取
   ```go
   fileName := input.FileName
   if fileName == "" {
       fileName = filepath.Base(input.FilePath)
   }
   ```

2. **一体化操作**: SendFileWithUpload 内部调用 UploadFile + SendFileMessage
   ```go
   uploadOutput, err := a.UploadFile(ctx, uploadInput)
   // ...
   sendInput := SendFileMessageInput{
       MediaID: uploadOutput.MediaID,
       // ...
   }
   ```

3. **完整的错误处理**: 
   - 文件打开失败
   - 上传失败
   - 发送失败
   - 权限不足

4. **返回 media_id**: 用户可以获取 media_id 用于后续操作

## 兼容性

- **Go 版本**: 1.23+
- **YouDu SDK**: v2.6.0
- **现有功能**: 完全兼容，无破坏性变更

## 后续优化建议

1. **支持 base64 编码的文件内容**: 对于 HTTP API，可以接受 base64 编码的文件内容而不是文件路径
2. **批量上传**: 支持一次上传多个文件
3. **进度回调**: 对于大文件上传，提供进度信息
4. **文件校验**: 添加文件大小、类型校验
5. **缓存 media_id**: 对于相同文件，可以缓存 media_id 避免重复上传

## 提交信息

```
feat(adapter): add file upload methods - UploadFile and SendFileWithUpload

- 添加 UploadFile 方法：上传文件并返回 media_id
- 添加 SendFileWithUpload 方法：一键上传并发送文件
- 支持文件类型：image, file, voice, video
- 自动生成 CLI、MCP、HTTP API 接口
- 集成权限检查系统
- 返回 media_id 在响应中

遵循项目架构：单一数据源，反射自动生成
```

## 文档

- [测试文档](test/FILE_UPLOAD_TEST.md)
- [测试脚本](test/scripts/test_upload_api.sh)
- [本实现总结](IMPLEMENTATION_FILE_UPLOAD.md)

## 时间线

- **2026-02-04**: 初始实现
  - 添加 UploadFile 方法
  - 添加 SendFileWithUpload 方法
  - 创建测试文档和脚本
  - 验证三种接口自动生成

---

**作者**: Claude Code + Happy
**日期**: 2026-02-04
**版本**: v1.0.0
