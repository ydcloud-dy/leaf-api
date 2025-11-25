# Leaf Blog 部署指南

本文档详细介绍了 Leaf Blog 系统的多种部署方式，包括裸部署、Docker 部署、Docker Compose 部署和 Kubernetes 部署。

## 项目结构

```
leaf-api/               # 后端 API 服务
├── Dockerfile
├── .dockerignore
├── deploy/
│   ├── docker/         # Docker 相关配置
│   ├── k8s/            # Kubernetes 配置
│   └── scripts/        # 部署脚本
└── ...

blog-frontend/          # 博客网站前端
├── Dockerfile
├── .dockerignore
├── deploy/
│   ├── k8s/            # Kubernetes 配置
│   ├── nginx/          # Nginx 配置
│   └── scripts/        # 部署脚本
└── ...

web/                    # 管理后台前端
├── Dockerfile
├── .dockerignore
├── deploy/
│   ├── k8s/            # Kubernetes 配置
│   ├── nginx/          # Nginx 配置
│   └── scripts/        # 部署脚本
└── ...
```

## 目录

- [1. 裸部署（Bare Metal）](#1-裸部署bare-metal)
- [2. Docker 部署](#2-docker-部署)
- [3. Docker Compose 部署](#3-docker-compose-部署)
- [4. Kubernetes 部署](#4-kubernetes-部署)

---

## 1. 裸部署（Bare Metal）

裸部署适合开发环境或简单的生产环境。

### 前置要求

- Go 1.21+
- Node.js 18+
- MySQL 8.0+
- Redis 7+
- Nginx（可选，用于前端部署）

### 1.1 后端 API 部署

```bash
# 进入后端目录
cd leaf-api

# 运行部署脚本
chmod +x deploy/scripts/deploy.sh
./deploy/scripts/deploy.sh
```

部署脚本会自动：
- 检查 Go 环境
- 检查 MySQL 和 Redis 连接
- 安装依赖
- 构建应用
- 创建必要的目录
- 启动应用

手动部署步骤：
```bash
# 1. 安装依赖
go mod download

# 2. 构建
go build -o leaf-api .

# 3. 创建必要的目录
mkdir -p logs uploads

# 4. 配置 config.yaml（根据实际环境修改）
cp config.yaml.example config.yaml
vim config.yaml

# 5. 启动
./leaf-api
```

### 1.2 博客前端部署

```bash
# 进入博客前端目录
cd blog-frontend

# 运行部署脚本
chmod +x deploy/scripts/deploy.sh
./deploy/scripts/deploy.sh
```

部署脚本会自动：
- 检查 Node.js 环境
- 安装依赖
- 构建应用
- 部署到 Nginx（如果有）或启动预览服务器

手动部署步骤：
```bash
# 1. 安装依赖
npm install

# 2. 构建
npm run build

# 3. 部署到 Nginx
sudo cp -r dist/* /usr/share/nginx/html/
sudo cp deploy/nginx/nginx.conf /etc/nginx/conf.d/blog-frontend.conf
sudo nginx -t && sudo systemctl restart nginx

# 或者使用预览服务器
npm run preview
```

### 1.3 管理后台部署

```bash
# 进入管理后台目录
cd web

# 运行部署脚本
chmod +x deploy/scripts/deploy.sh
./deploy/scripts/deploy.sh
```

部署步骤与博客前端类似。

---

## 2. Docker 部署

使用 Docker 单独部署各个服务。

### 2.1 后端 API

```bash
# 构建镜像
docker build -t leaf-api:latest .

# 运行容器
docker run -d \
  --name leaf-api \
  -p 8888:8888 \
  -v $(pwd)/config.yaml:/app/config.yaml:ro \
  -v $(pwd)/uploads:/app/uploads \
  -v $(pwd)/logs:/app/logs \
  -e DB_HOST=your-mysql-host \
  -e DB_PORT=3306 \
  -e REDIS_HOST=your-redis-host \
  -e REDIS_PORT=6379 \
  leaf-api:latest
```

### 2.2 博客前端

```bash
# 构建镜像
cd blog-frontend
docker build -t blog-frontend:latest .

# 运行容器
docker run -d \
  --name blog-frontend \
  -p 3000:80 \
  blog-frontend:latest
```

### 2.3 管理后台

```bash
# 构建镜像
cd web
docker build -t admin-frontend:latest .

# 运行容器
docker run -d \
  --name admin-frontend \
  -p 3001:80 \
  admin-frontend:latest
```

---

## 3. Docker Compose 部署

Docker Compose 是推荐的部署方式，可以一键启动所有服务。

### 3.1 快速开始

```bash
# 1. 克隆项目
git clone <repository-url>
cd leaf-api

# 2. 配置环境
cp config.yaml.example config.yaml
# 根据需要修改 config.yaml 和 docker-compose.yml

# 3. 启动所有服务
docker-compose up -d

# 4. 查看日志
docker-compose logs -f

# 5. 停止所有服务
docker-compose down

# 6. 停止并删除数据卷（谨慎使用）
docker-compose down -v
```

### 3.2 服务端口

- MySQL: `3306`
- Redis: `6379`
- 后端 API: `8888`
- 博客前端: `3000`
- 管理后台: `3001`

### 3.3 访问地址

- 博客网站: http://localhost:3000
- 管理后台: http://localhost:3001
- 后端 API: http://localhost:8888

### 3.4 数据持久化

Docker Compose 会自动创建以下数据卷：
- `mysql_data`: MySQL 数据
- `redis_data`: Redis 数据
- `./uploads`: 上传文件
- `./logs`: 日志文件

### 3.5 常用命令

```bash
# 查看服务状态
docker-compose ps

# 查看某个服务的日志
docker-compose logs -f api
docker-compose logs -f blog-frontend

# 重启某个服务
docker-compose restart api

# 重新构建并启动
docker-compose up -d --build

# 进入容器
docker-compose exec api sh
docker-compose exec mysql mysql -uroot -p123456
```

---

## 4. Kubernetes 部署

适合大规模生产环境，支持高可用和自动扩缩容。

### 4.1 前置要求

- Kubernetes 集群（1.20+）
- kubectl 命令行工具
- 容器镜像仓库（用于存储镜像）
- Ingress Controller（如 nginx-ingress）
- 可选：cert-manager（用于自动管理 TLS 证书）

### 4.2 准备镜像

```bash
# 1. 构建后端镜像
docker build -t your-registry/leaf-api:latest .
docker push your-registry/leaf-api:latest

# 2. 构建博客前端镜像
cd blog-frontend
docker build -t your-registry/blog-frontend:latest .
docker push your-registry/blog-frontend:latest

# 3. 构建管理后台镜像
cd ../web
docker build -t your-registry/admin-frontend:latest .
docker push your-registry/admin-frontend:latest
```

### 4.3 修改配置

修改以下文件中的镜像地址和域名：

1. **后端配置** (`deploy/k8s/deployment.yaml`):
   - 修改 `image: your-registry/leaf-api:latest`
   - 修改 `host: api.yourdomain.com`
   - 根据需要修改 ConfigMap 中的配置

2. **博客前端配置** (`blog-frontend/deploy/k8s/deployment.yaml`):
   - 修改 `image: your-registry/blog-frontend:latest`
   - 修改 `host: blog.yourdomain.com`

3. **管理后台配置** (`web/deploy/k8s/deployment.yaml`):
   - 修改 `image: your-registry/admin-frontend:latest`
   - 修改 `host: admin.yourdomain.com`

### 4.4 部署步骤

```bash
# 1. 创建命名空间和 PVC
kubectl apply -f deploy/k8s/pvc.yaml

# 2. 部署后端服务（包括 MySQL、Redis、API）
kubectl apply -f deploy/k8s/deployment.yaml

# 3. 部署博客前端
kubectl apply -f blog-frontend/deploy/k8s/deployment.yaml

# 4. 部署管理后台
kubectl apply -f web/deploy/k8s/deployment.yaml

# 5. 查看部署状态
kubectl get pods -n leaf-blog
kubectl get svc -n leaf-blog
kubectl get ingress -n leaf-blog
```

### 4.5 一键部署脚本

```bash
# 创建一键部署脚本
cat > deploy-k8s.sh << 'EOF'
#!/bin/bash
set -e

echo "🚀 开始部署到 Kubernetes..."

# 应用所有配置
kubectl apply -f deploy/k8s/pvc.yaml
kubectl apply -f deploy/k8s/deployment.yaml
kubectl apply -f blog-frontend/deploy/k8s/deployment.yaml
kubectl apply -f web/deploy/k8s/deployment.yaml

echo "⏳ 等待 Pod 就绪..."
kubectl wait --for=condition=ready pod -l app=leaf-api -n leaf-blog --timeout=300s
kubectl wait --for=condition=ready pod -l app=blog-frontend -n leaf-blog --timeout=300s
kubectl wait --for=condition=ready pod -l app=admin-frontend -n leaf-blog --timeout=300s

echo "✅ 部署完成！"
echo ""
echo "查看状态："
kubectl get pods -n leaf-blog
echo ""
echo "访问地址："
kubectl get ingress -n leaf-blog
EOF

chmod +x deploy-k8s.sh
./deploy-k8s.sh
```

### 4.6 常用 K8s 命令

```bash
# 查看所有资源
kubectl get all -n leaf-blog

# 查看 Pod 日志
kubectl logs -f <pod-name> -n leaf-blog

# 查看 Pod 详情
kubectl describe pod <pod-name> -n leaf-blog

# 进入容器
kubectl exec -it <pod-name> -n leaf-blog -- sh

# 扩容/缩容
kubectl scale deployment leaf-api --replicas=3 -n leaf-blog

# 更新镜像
kubectl set image deployment/leaf-api leaf-api=your-registry/leaf-api:v2 -n leaf-blog

# 查看配置
kubectl get configmap leaf-api-config -n leaf-blog -o yaml

# 删除所有资源
kubectl delete namespace leaf-blog
```

### 4.7 高可用配置

对于生产环境，建议：

1. **数据库高可用**：使用云数据库服务或 MySQL 集群
2. **Redis 高可用**：使用 Redis Sentinel 或 Redis Cluster
3. **应用多副本**：API 和前端至少 2 个副本
4. **资源限制**：根据实际负载调整 resources 配置
5. **健康检查**：配置合适的 livenessProbe 和 readinessProbe
6. **日志收集**：集成 ELK 或其他日志系统
7. **监控告警**：集成 Prometheus + Grafana

### 4.8 存储类配置

根据云平台修改 `storageClassName`：

- **阿里云**: `alicloud-disk-ssd`
- **腾讯云**: `cbs`
- **AWS**: `gp2` 或 `gp3`
- **GCP**: `standard` 或 `ssd`
- **本地**: `local-path` 或 `nfs`

---

## 5. 环境变量说明

### 后端 API 环境变量

| 变量名 | 说明 | 默认值 |
|--------|------|--------|
| `DB_HOST` | MySQL 主机地址 | localhost |
| `DB_PORT` | MySQL 端口 | 3306 |
| `DB_USER` | MySQL 用户名 | root |
| `DB_PASSWORD` | MySQL 密码 | 123456 |
| `DB_NAME` | 数据库名 | leaf_admin |
| `REDIS_HOST` | Redis 主机地址 | localhost |
| `REDIS_PORT` | Redis 端口 | 6379 |
| `REDIS_PASSWORD` | Redis 密码 | - |
| `JWT_SECRET` | JWT 密钥 | - |
| `TZ` | 时区 | Asia/Shanghai |

---

## 6. 故障排查

### 6.1 后端 API 无法启动

```bash
# 检查日志
tail -f logs/app.log

# 检查端口占用
lsof -i :8888

# 检查数据库连接
mysql -h127.0.0.1 -P3306 -uroot -p123456

# 检查 Redis 连接
redis-cli -h 127.0.0.1 -p 6379 ping
```

### 6.2 前端无法访问

```bash
# 检查 Nginx 配置
nginx -t

# 查看 Nginx 日志
tail -f /var/log/nginx/error.log

# 检查构建产物
ls -la dist/
```

### 6.3 Docker Compose 问题

```bash
# 查看所有服务状态
docker-compose ps

# 查看特定服务日志
docker-compose logs api

# 重新构建
docker-compose build --no-cache

# 清理并重新启动
docker-compose down -v
docker-compose up -d
```

### 6.4 Kubernetes 问题

```bash
# 查看 Pod 状态
kubectl get pods -n leaf-blog

# 查看 Pod 日志
kubectl logs <pod-name> -n leaf-blog

# 查看 Pod 事件
kubectl describe pod <pod-name> -n leaf-blog

# 查看 Service
kubectl get svc -n leaf-blog

# 检查 Ingress
kubectl describe ingress -n leaf-blog
```

---

## 7. 安全建议

1. **更改默认密码**：修改 MySQL、Redis 的默认密码
2. **JWT 密钥**：使用强随机密钥
3. **HTTPS**：生产环境启用 HTTPS
4. **防火墙**：限制不必要的端口访问
5. **定期更新**：及时更新依赖和系统补丁
6. **备份**：定期备份数据库和上传文件
7. **监控**：配置日志和性能监控

---

## 8. 性能优化

1. **数据库优化**：
   - 添加合适的索引
   - 定期分析和优化查询
   - 配置连接池

2. **Redis 缓存**：
   - 缓存热点数据
   - 设置合理的过期时间

3. **前端优化**：
   - 启用 Gzip 压缩
   - 配置静态资源缓存
   - 使用 CDN

4. **负载均衡**：
   - 使用 Nginx 或云负载均衡
   - 配置多个后端实例

---

## 9. 备份和恢复

### 9.1 数据库备份

```bash
# 备份
docker-compose exec mysql mysqldump -uroot -p123456 leaf_admin > backup.sql

# 恢复
docker-compose exec -T mysql mysql -uroot -p123456 leaf_admin < backup.sql
```

### 9.2 文件备份

```bash
# 备份上传文件
tar -czf uploads-backup.tar.gz uploads/

# 恢复
tar -xzf uploads-backup.tar.gz
```

---

## 10. 联系和支持

如有问题，请提交 Issue 或联系维护团队。
