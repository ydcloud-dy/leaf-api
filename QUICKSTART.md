# Leaf Blog 快速部署参考

## 一键部署（推荐）

```bash
./deploy-all.sh
```

选择部署方式：
1. 裸部署 - 直接在主机运行
2. Docker - 单独容器部署
3. Docker Compose - 一键启动（最简单）
4. Kubernetes - 生产环境

## 核心配置说明

### API 路由结构
```
后端 (localhost:8888)
├── /blog/*      → 博客前端 API
└── /*           → 管理后台 API
```

### 前端请求路径
```
博客前端 → /blog/* → Nginx → backend:8888/blog/*
管理后台 → /api/*  → Nginx → backend:8888/*  (去掉/api)
```

## Docker Compose 快速开始

```bash
# 启动
docker-compose up -d

# 访问
http://localhost:3000  # 博客网站
http://localhost:3001  # 管理后台
http://localhost:8888  # API

# 查看日志
docker-compose logs -f

# 停止
docker-compose down
```

## 常见问题修复

### 404 错误
检查:
- blog-frontend/src/api/request.js → baseURL: '/blog'
- Nginx location /blog/ → proxy_pass http://leaf-api:8888

### 连接失败
```bash
# 检查服务
docker-compose ps

# 测试后端
curl http://localhost:8888/ping
curl http://localhost:8888/blog/stats

# 查看日志
docker-compose logs api
docker-compose logs blog-frontend
```

## 详细文档

查看 [DEPLOYMENT.md](./DEPLOYMENT.md) 获取完整部署指南。
