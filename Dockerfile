# 多阶段构建 Dockerfile for YouDu MCP Service

# 第一阶段：构建环境
FROM golang:1.24-alpine AS builder

# 安装构建依赖
RUN apk add --no-cache git gcc musl-dev

# 设置工作目录
WORKDIR /build

# 复制 go mod 文件并下载依赖
COPY go.mod go.sum ./
RUN go mod download

# 复制源代码
COPY . .

# 构建二进制文件
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -ldflags="-w -s" -o youdu-cli ./cmd/youdu-cli
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -ldflags="-w -s" -o youdu-mcp ./cmd/youdu-mcp

# 第二阶段：运行环境
FROM alpine:latest

# 安装运行时依赖
RUN apk add --no-cache ca-certificates tzdata

# 设置时区
ENV TZ=Asia/Shanghai

# 创建非 root 用户
RUN addgroup -g 1000 youdu && \
    adduser -D -u 1000 -G youdu youdu

# 创建必要的目录
RUN mkdir -p /app/data /app/config && \
    chown -R youdu:youdu /app

# 设置工作目录
WORKDIR /app

# 从构建阶段复制二进制文件
COPY --from=builder /build/youdu-cli /app/youdu-cli
COPY --from=builder /build/youdu-mcp /app/youdu-mcp

# 复制配置文件示例
COPY config.yaml.example /app/config.yaml.example

# 设置文件权限
RUN chmod +x /app/youdu-cli /app/youdu-mcp && \
    chown -R youdu:youdu /app

# 切换到非 root 用户
USER youdu

# 默认命令（可以在 docker-compose.yml 中覆盖）
CMD ["/app/youdu-cli", "serve-api", "--port", "8080"]
