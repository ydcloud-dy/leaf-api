# Leaf API - 后端服务

基于 Go + Gin 框架开发的博客系统后端 API 服务。

## 📋 目录

- [技术栈](#技术栈)
- [项目架构](#项目架构)
- [目录结构](#目录结构)
- [快速开始](#快速开始)
- [部署方式](#部署方式)
  - [裸部署](#裸部署)
  - [Docker 部署](#docker-部署)
  - [Docker Compose 部署](#docker-compose-部署)
  - [Kubernetes 部署](#kubernetes-部署)
- [配置说明](#配置说明)
- [API 文档](#api-文档)

## 🛠 技术栈

- **语言**: Go 1.21+
- **Web 框架**: Gin
- **ORM**: GORM
- **数据库**: MySQL 8.0
- **缓存**: Redis 7.x
- **对象存储**: 阿里云 OSS
- **认证**: JWT
- **日志**: Logrus + Lumberjack
- **配置管理**: Viper
- **依赖注入**: Wire

## 🏗 项目架构

采用分层架构设计：

```
┌─────────────────┐
│   HTTP Layer    │  (Gin Router)
├─────────────────┤
│  Service Layer  │  (Business Logic)
├─────────────────┤
│  UseCase Layer  │  (Application Logic)
├─────────────────┤
│   Data Layer    │  (Repository Pattern)
├─────────────────┤
│  Model Layer    │  (PO/DTO)
└─────────────────┘
```

## 📂 目录结构

```
.
├── cmd/                    # 命令行入口
├── config/                 # 配置文件
│   └── config.yaml        # 默认配置
├── deploy/                 # 部署配置
│   ├── docker/            # Docker 相关
│   ├── k8s/               # Kubernetes 配置
│   └── scripts/           # 部署脚本
├── docs/                   # 文档
├── internal/              # 内部代码
│   ├── biz/              # 业务逻辑层
│   ├── data/             # 数据访问层
│   ├── model/            # 数据模型
│   │   ├── dto/         # 数据传输对象
│   │   └── po/          # 持久化对象
│   ├── server/           # 服务器配置
│   └── service/          # 服务层
├── pkg/                   # 公共包
│   ├── middleware/       # 中间件
│   ├── response/         # 响应封装
│   └── utils/            # 工具函数
├── config.yaml            # 配置文件
├── main.go               # 程序入口
├── Dockerfile            # Docker 镜像构建文件
├── docker-compose.yml    # Docker Compose 配置
└── README.md             # 项目文档
```

## 🚀 快速开始

### 环境要求

- Go 1.21+
- MySQL 8.0+
- Redis 7.x

### 本地开发

1. **克隆项目**

```bash
git clone https://github.com/ydcloud-dy/leaf-api.git
cd leaf-api
```

2. **安装依赖**

```bash
go mod download
```

3. **配置数据库**

修改 `config.yaml` 中的数据库配置：

```yaml
database:
  host: 127.0.0.1
  port: 3306
  user: root
  password: your_password
  dbname: leaf_admin
  charset: utf8mb4
```

4. **启动 MySQL 和 Redis**

```bash
# MySQL
mysql.server start

# Redis
redis-server
```

5. **运行应用**

```bash
# 开发模式
go run main.go

# 或使用编译后的二进制
go build -o leaf-api .
./leaf-api
```

6. **访问服务**

- API 服务: http://localhost:8888
- 健康检查: http://localhost:8888/ping

## 📦 部署方式

### 裸部署

使用自动化部署脚本：

```bash
# 运行部署脚本
chmod +x deploy/scripts/deploy.sh
./deploy/scripts/deploy.sh
```

或手动部署：

```bash
# 1. 构建
go build -o leaf-api .

# 2. 创建必要目录
mkdir -p logs uploads

# 3. 配置 config.yaml

# 4. 启动
./leaf-api

# 或后台运行
nohup ./leaf-api > server.log 2>&1 &
```

### Docker 部署

#### 构建镜像

```bash
docker build -t leaf-api:latest .
```

#### 运行容器

```bash
docker run -d \
  --name leaf-api \
  -p 8888:8888 \
  -v $(pwd)/config.yaml:/app/config.yaml \
  -v $(pwd)/uploads:/app/uploads \
  -v $(pwd)/logs:/app/logs \
  -e DB_HOST=your_mysql_host \
  -e REDIS_HOST=your_redis_host \
  leaf-api:latest
```

### Docker Compose 部署

**一键启动所有服务（推荐）**

```bash
# 启动所有服务（API + MySQL + Redis + 前端）
docker-compose up -d

# 查看日志
docker-compose logs -f

# 停止服务
docker-compose down

# 停止并清理数据
docker-compose down -v
```

访问地址：
- API 服务: http://localhost:8888
- 网站端: http://localhost:3000
- 管理端: http://localhost:3001

### Kubernetes 部署

#### 1. 创建命名空间和 PVC

```bash
kubectl apply -f deploy/k8s/pvc.yaml
```

#### 2. 部署应用

```bash
kubectl apply -f deploy/k8s/deployment.yaml
```

#### 3. 检查部署状态

```bash
# 查看 Pod
kubectl get pods -n leaf-blog

# 查看服务
kubectl get svc -n leaf-blog

# 查看日志
kubectl logs -f <pod-name> -n leaf-blog
```

#### 4. 访问服务

```bash
# 端口转发（用于测试）
kubectl port-forward svc/leaf-api-service 8888:8888 -n leaf-blog

# 或配置 Ingress 后通过域名访问
```

## ⚙️ 配置说明

### config.yaml 配置项

```yaml
server:
  port: 8888           # 服务端口
  mode: release        # 运行模式: debug, release, test

database:
  host: 127.0.0.1      # 数据库地址
  port: 3306           # 数据库端口
  user: root           # 数据库用户
  password: 123456     # 数据库密码
  dbname: leaf_admin   # 数据库名称
  charset: utf8mb4     # 字符集

jwt:
  secret: your_secret  # JWT 密钥
  expire: 24           # Token 过期时间（小时）

oss:                   # 阿里云 OSS 配置
  endpoint: oss-cn-hangzhou.aliyuncs.com
  access_key_id: your_key_id
  access_key_secret: your_key_secret
  bucket_name: your_bucket
  base_url: https://your-bucket.oss-cn-hangzhou.aliyuncs.com

redis:
  host: 127.0.0.1      # Redis 地址
  port: 6379           # Redis 端口
  password:            # Redis 密码
  db: 0                # Redis 数据库
  pool_size: 10        # 连接池大小

log:
  level: debug         # 日志级别: debug, info, warn, error
  format: text         # 日志格式: text, json
  output: stdout       # 输出: stdout, file
  file_path: logs/app.log
  max_size: 100        # 单个文件大小 (MB)
  max_backups: 3       # 保留文件数量
  max_age: 7           # 保留天数
```

### 环境变量

支持通过环境变量覆盖配置：

- `DB_HOST`: 数据库地址
- `DB_PORT`: 数据库端口
- `REDIS_HOST`: Redis 地址
- `REDIS_PORT`: Redis 端口

## 📖 API 文档

### 认证相关

- `POST /auth/register` - 用户注册
- `POST /auth/login` - 用户登录
- `POST /admin/login` - 管理员登录
- `POST /auth/refresh` - 刷新 Token

### 文章管理

- `GET /blog/articles` - 获取文章列表
- `GET /blog/articles/:id` - 获取文章详情
- `POST /articles` - 创建文章（需认证）
- `PUT /articles/:id` - 更新文章（需认证）
- `DELETE /articles/:id` - 删除文章（需认证）
- `POST /articles/import` - 批量导入 Markdown 文件

### 评论管理

- `GET /blog/articles/:id/comments` - 获取文章评论
- `POST /blog/comments` - 发表评论（需认证）
- `POST /guestbook` - 留言板消息（需认证）

### 用户管理

- `GET /users` - 获取用户列表（需管理员）
- `GET /users/:id` - 获取用户详情
- `PUT /users/:id` - 更新用户信息（需认证）

更多 API 详情请查看 [API 文档](docs/api.md)

## 🔧 开发相关

### 运行测试

```bash
go test ./...
```

### 代码格式化

```bash
go fmt ./...
```

### 代码检查

```bash
go vet ./...
```

## 📝 License

MIT License

## 👥 贡献

欢迎提交 Issue 和 Pull Request！

## 📞 联系方式

- 项目地址: https://github.com/ydcloud-dy/leaf-api
- 问题反馈: https://github.com/ydcloud-dy/leaf-api/issues
