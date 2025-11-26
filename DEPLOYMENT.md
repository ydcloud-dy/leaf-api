# Leaf Blog 完整部署指南

本文档提供 Leaf Blog 的完整部署方案，所有方法都经过测试，可以一键部署并直接访问。

## 快速开始

```bash
# 一键部署脚本（支持所有部署方式）
./deploy-all.sh
```

## 项目结构

```
leaf-api/               # 后端 Go API (端口 8888)
├── 路由: /blog/*       # 博客前端 API
├── 路由: /*            # 管理后台 API
└── ...

blog-frontend/          # 博客网站前端 (端口 3000/4173)
├── API路径: /blog/*    # 请求后端博客API
└── ...

web/                    # 管理后台前端 (端口 3001/4174)
├── API路径: /api/*     # 代理到后端根路径
└── ...
```

## 重要：API 路由说明

### 后端路由结构
- `/blog/*` - 博客前端API（如 `/blog/articles`, `/blog/stats`, `/blog/heartbeat`）
- `/*` - 管理后台API（如 `/articles`, `/users`, `/auth/login`）

### 前端请求路径
- **博客前端**: 请求 `/blog/*`，Nginx 直接代理到后端 `http://backend:8888/blog/*`
- **管理后台**: 请求 `/api/*`，Nginx 代理到后端 `http://backend:8888/*`（去掉 /api 前缀）

## 目录

- [方式一：Docker Compose 部署（推荐）](#方式一docker-compose-部署推荐)
- [方式二：裸部署](#方式二裸部署)
- [方式三：Docker 部署](#方式三docker-部署)
- [方式四：Kubernetes 部署](#方式四kubernetes-部署)
- [故障排查](#故障排查)

---

## 方式一：Docker Compose 部署（推荐）

这是最简单的部署方式，一键启动所有服务。

### 前置要求

- Docker 20.10+
- Docker Compose 2.0+

### 快速部署

```bash
# 1. 克隆项目
git clone <repository-url>
cd leaf-api

# 2. 确保配置文件存在
cp config.yaml.example config.yaml
# 根据需要修改 config.yaml

# 3. 一键启动
docker-compose up -d

# 4. 查看日志
docker-compose logs -f

# 5. 访问服务
# 博客网站: http://localhost:3000
# 管理后台: http://localhost:3001
# 后端 API: http://localhost:8888
```

### 服务端口映射

| 服务 | 容器端口 | 主机端口 | 说明 |
|------|---------|---------|------|
| MySQL | 3306 | 3306 | 数据库 |
| Redis | 6379 | 6379 | 缓存 |
| 后端 API | 8888 | 8888 | Go API |
| 博客前端 | 80 | 3000 | Nginx + Vue |
| 管理后台 | 80 | 3001 | Nginx + Vue |

### 常用命令

```bash
# 查看服务状态
docker-compose ps

# 查看日志
docker-compose logs -f [service-name]

# 重启服务
docker-compose restart [service-name]

# 停止服务
docker-compose stop

# 删除服务（不删除数据卷）
docker-compose down

# 删除服务和数据卷（慎用）
docker-compose down -v

# 重新构建并启动
docker-compose up -d --build
```

### 数据持久化

数据卷会自动创建：
- `mysql_data`: MySQL 数据
- `redis_data`: Redis 数据
- `./uploads`: 上传文件
- `./logs`: 日志文件

---

## 方式二：裸部署

直接在Linux主机上运行，适合开发环境。

### 前置要求

- Go 1.21+
- Node.js 18+
- MySQL 8.0+
- Redis 7+
- Nginx（可选，用于生产环境）

### 一键部署

```bash
./deploy-all.sh
# 选择 1) 裸部署
```

### 手动部署步骤

#### 1. 准备数据库

```bash
# 启动 MySQL
sudo systemctl start mysql

# 创建数据库
mysql -uroot -p
CREATE DATABASE leaf_admin CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

# 导入初始化脚本（如果有）
mysql -uroot -p leaf_admin < deploy/docker/mysql/init.sql

# 启动 Redis
sudo systemctl start redis
```

#### 2. 部署后端

```bash
cd leaf-api

# 安装依赖
go mod download

# 配置文件
cp config.yaml.example config.yaml
vim config.yaml  # 修改数据库、Redis 等配置

# 构建
go build -o leaf-api .

# 创建目录
mkdir -p logs uploads

# 启动
nohup ./leaf-api > logs/server.log 2>&1 &

# 验证
curl http://localhost:8888/ping
```

#### 3. 部署博客前端

```bash
cd blog-frontend

# 安装依赖
npm install

# 构建
npm run build

# 使用预览服务器（开发）
nohup npm run preview > ../logs/blog-frontend.log 2>&1 &

# 或者部署到 Nginx（生产）
sudo cp -r dist/* /var/www/blog-frontend/
```

#### 4. 部署管理后台

```bash
cd web

# 安装依赖
npm install

# 构建
npm run build

# 使用预览服务器（开发）
nohup npm run preview -- --port 4174 > ../logs/admin-frontend.log 2>&1 &

# 或者部署到 Nginx（生产）
sudo cp -r dist/* /var/www/admin-frontend/
```

#### 5. 配置 Nginx（生产环境）

```bash
# 复制配置文件
sudo cp deploy/nginx/production.conf /etc/nginx/sites-available/leaf-blog.conf

# 修改域名
sudo vim /etc/nginx/sites-available/leaf-blog.conf

# 创建软链接
sudo ln -s /etc/nginx/sites-available/leaf-blog.conf /etc/nginx/sites-enabled/

# 测试配置
sudo nginx -t

# 重启 Nginx
sudo systemctl restart nginx
```

### 访问地址

- **开发环境**:
  - 博客网站: http://localhost:4173
  - 管理后台: http://localhost:4174
  - 后端 API: http://localhost:8888

- **生产环境（Nginx）**:
  - 博客网站: http://yourdomain.com
  - 管理后台: http://admin.yourdomain.com
  - 后端 API: http://api.yourdomain.com

---

## 方式三:Docker 部署

使用 Docker 单独部署各个服务。

### 一键部署

```bash
./deploy-all.sh
# 选择 2) Docker 部署
```

### 手动部署步骤

```bash
# 1. 创建网络
docker network create leaf-network

# 2. 启动 MySQL
docker run -d \
  --name leaf-mysql \
  --network leaf-network \
  -e MYSQL_ROOT_PASSWORD=123456 \
  -e MYSQL_DATABASE=leaf_admin \
  -e TZ=Asia/Shanghai \
  -p 3306:3306 \
  -v leaf-mysql-data:/var/lib/mysql \
  mysql:8.0

# 3. 启动 Redis
docker run -d \
  --name leaf-redis \
  --network leaf-network \
  -p 6379:6379 \
  -v leaf-redis-data:/data \
  redis:7-alpine \
  redis-server --appendonly yes

# 4. 构建并启动后端
docker build -t leaf-api:latest .
docker run -d \
  --name leaf-api \
  --network leaf-network \
  -p 8888:8888 \
  -v $(pwd)/config.yaml:/app/config.yaml:ro \
  -v $(pwd)/uploads:/app/uploads \
  -e DB_HOST=leaf-mysql \
  -e REDIS_HOST=leaf-redis \
  leaf-api:latest

# 5. 构建并启动博客前端
cd blog-frontend
docker build -t blog-frontend:latest .
docker run -d \
  --name leaf-blog-frontend \
  --network leaf-network \
  -p 3000:80 \
  blog-frontend:latest

# 6. 构建并启动管理后台
cd ../web
docker build -t admin-frontend:latest .
docker run -d \
  --name leaf-admin-frontend \
  --network leaf-network \
  -p 3001:80 \
  admin-frontend:latest
```

---

## 方式四：Kubernetes 部署

适合生产环境，支持高可用和自动扩缩容。

### 前置要求

- Kubernetes 集群 1.20+
- kubectl
- 容器镜像仓库

### 一键部署

```bash
./deploy-k8s.sh
```

### 手动部署步骤

#### 1. 准备镜像

```bash
# 构建后端镜像
docker build -t your-registry/leaf-api:latest .
docker push your-registry/leaf-api:latest

# 构建博客前端镜像
cd blog-frontend
docker build -t your-registry/blog-frontend:latest .
docker push your-registry/blog-frontend:latest

# 构建管理后台镜像
cd ../web
docker build -t your-registry/admin-frontend:latest .
docker push your-registry/admin-frontend:latest
```

#### 2. 修改配置

修改以下文件中的镜像地址和域名：

1. `deploy/k8s/deployment.yaml` - 修改 API 镜像和域名
2. `blog-frontend/deploy/k8s/deployment.yaml` - 修改博客前端镜像和域名
3. `web/deploy/k8s/deployment.yaml` - 修改管理后台镜像和域名

#### 3. 部署

```bash
# 创建命名空间和 PVC
kubectl apply -f deploy/k8s/pvc.yaml

# 部署后端服务
kubectl apply -f deploy/k8s/deployment.yaml

# 部署博客前端
kubectl apply -f blog-frontend/deploy/k8s/deployment.yaml

# 部署管理后台
kubectl apply -f web/deploy/k8s/deployment.yaml

# 查看状态
kubectl get pods -n leaf-blog
kubectl get svc -n leaf-blog
kubectl get ingress -n leaf-blog
```

---

## 故障排查

### 1. API 接口 404 错误

**问题**: 访问 `/api/heartbeat` 或 `/api/visit` 返回 404

**原因**:
- 博客前端应该请求 `/blog/heartbeat` 和 `/blog/visit`
- Nginx 配置不正确

**解决方案**:

1. 检查前端 API 配置:
```javascript
// blog-frontend/src/api/request.js
baseURL: '/blog'  // 应该是 /blog 而不是 /api
```

2. 检查 Nginx 配置:
```nginx
# 博客前端 Nginx
location /blog/ {
    proxy_pass http://leaf-api:8888;  # 正确
}

# 管理后台 Nginx
location /api/ {
    proxy_pass http://leaf-api:8888/;  # 注意末尾的 /
}
```

### 2. 前端无法访问后端

**检查步骤**:

```bash
# Docker Compose
docker-compose ps  # 查看所有服务状态
docker-compose logs api  # 查看后端日志
docker-compose logs blog-frontend  # 查看前端日志

# 测试后端连接
curl http://localhost:8888/ping
curl http://localhost:8888/blog/stats

# 测试前端Nginx配置
docker exec leaf-blog-frontend nginx -t
```

### 3. 数据库连接失败

```bash
# 检查 MySQL 是否启动
docker ps | grep mysql

# 进入 MySQL 容器
docker exec -it leaf-mysql mysql -uroot -p123456

# 检查数据库
SHOW DATABASES;
USE leaf_admin;
SHOW TABLES;

# 查看后端日志
docker logs leaf-api
```

### 4. Redis 连接失败

```bash
# 检查 Redis
docker exec -it leaf-redis redis-cli ping

# 测试连接
docker exec -it leaf-redis redis-cli
> PING
PONG
```

### 5. Nginx 配置错误

```bash
# Docker 环境
docker exec leaf-blog-frontend nginx -t
docker exec leaf-blog-frontend cat /etc/nginx/conf.d/default.conf

# 裸部署
sudo nginx -t
sudo cat /etc/nginx/sites-enabled/leaf-blog.conf
```

### 6. 端口占用

```bash
# 查看端口占用
lsof -i :8888
lsof -i :3000
lsof -i :3001

# 停止占用进程
kill -9 <PID>
```

---

## 性能优化建议

### 1. 后端优化
- 配置合适的数据库连接池
- 启用 Redis 缓存热点数据
- 添加数据库索引

### 2. 前端优化
- 启用 Gzip 压缩
- 配置静态资源缓存
- 使用 CDN

### 3. 数据库优化
```sql
-- 添加索引示例
CREATE INDEX idx_article_status ON articles(status);
CREATE INDEX idx_article_created ON articles(created_at);
```

---

## 安全建议

1. **更改默认密码**
   - MySQL root 密码
   - Redis 密码
   - JWT 密钥

2. **启用 HTTPS**
```nginx
server {
    listen 443 ssl http2;
    ssl_certificate /path/to/cert.pem;
    ssl_certificate_key /path/to/key.pem;
}
```

3. **配置防火墙**
```bash
# 只开放必要端口
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp
sudo ufw enable
```

4. **定期备份**
```bash
# 数据库备份
docker exec leaf-mysql mysqldump -uroot -p123456 leaf_admin > backup-$(date +%Y%m%d).sql

# 上传文件备份
tar -czf uploads-backup-$(date +%Y%m%d).tar.gz uploads/
```

---

## 常见问题

**Q: 如何修改端口？**

A: 修改 `docker-compose.yml` 中的端口映射：
```yaml
ports:
  - "8080:8888"  # 将 8888 改为 8080
```

**Q: 如何查看日志？**

A:
```bash
# Docker Compose
docker-compose logs -f [service-name]

# 裸部署
tail -f logs/server.log
tail -f logs/blog-frontend.log
```

**Q: 如何重启服务？**

A:
```bash
# Docker Compose
docker-compose restart

# 裸部署
pkill -f leaf-api && nohup ./leaf-api > logs/server.log 2>&1 &
```

**Q: 如何清理所有数据？**

A:
```bash
./deploy-all.sh
# 选择 5) 清理所有容器和数据
```

---

## 联系支持

如有问题，请：
1. 查看日志文件
2. 检查本文档的故障排查部分
3. 提交 Issue

祝部署顺利！
