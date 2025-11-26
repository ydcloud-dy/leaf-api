# 多阶段构建 Dockerfile

# 第一阶段：构建
FROM swr.cn-north-4.myhuaweicloud.com/ddn-k8s/docker.io/golang:1.21-alpine AS builder

# 设置工作目录
WORKDIR /app
RUN echo "https://mirrors.aliyun.com/alpine/v3.20/main/" > /etc/apk/repositories && \
    echo "https://mirrors.aliyun.com/alpine/v3.20/community/" >> /etc/apk/repositories
# 安装必要的构建工具
RUN apk add --no-cache git make
# 3. 设置国内 Go 模块源
ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct \
    CGO_ENABLED=0 \
    GOOS=linux
# 复制 go mod 文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 构建应用
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o leaf-api .

# 第二阶段：运行
FROM swr.cn-north-4.myhuaweicloud.com/ddn-k8s/docker.io/library/alpine:latest
RUN echo "https://mirrors.aliyun.com/alpine/v3.20/main/" > /etc/apk/repositories && \
    echo "https://mirrors.aliyun.com/alpine/v3.20/community/" >> /etc/apk/repositories
# 安装必要的运行时依赖
RUN apk --no-cache add ca-certificates tzdata

# 设置时区
ENV TZ=Asia/Shanghai

# 创建非 root 用户
RUN addgroup -g 1000 appgroup && \
    adduser -D -u 1000 -G appgroup appuser

# 设置工作目录
WORKDIR /app

# 从构建阶段复制二进制文件
COPY --from=builder /app/leaf-api .

# 复制配置文件
COPY config.yaml .

# 创建必要的目录
RUN mkdir -p logs uploads && \
    chown -R appuser:appgroup /app

# 切换到非 root 用户
USER appuser

# 暴露端口
EXPOSE 8888

# 健康检查
HEALTHCHECK --interval=30s --timeout=3s --start-period=40s --retries=3 \
    CMD wget --quiet --tries=1 --spider http://localhost:8888/ping || exit 1

# 启动应用
CMD ["./leaf-api"]