# Leaf Blog 部署文件说明

本目录包含 Leaf Blog 后端服务的所有部署相关文件。

## 目录结构

```
deploy/
├── docker/             # Docker 相关配置
│   └── mysql/
│       └── init.sql    # MySQL 初始化脚本
├── k8s/                # Kubernetes 配置
│   ├── deployment.yaml # 后端、MySQL、Redis 部署配置
│   └── pvc.yaml        # 持久化存储配置
└── scripts/            # 部署脚本
    └── deploy.sh       # 裸部署脚本
```

## 使用说明

详细的部署说明请查看项目根目录的 [DEPLOYMENT.md](../DEPLOYMENT.md)。

### 快速开始

**裸部署：**
```bash
./scripts/deploy.sh
```

**Docker 部署：**
```bash
docker build -t leaf-api:latest ..
docker run -d -p 8888:8888 leaf-api:latest
```

**Docker Compose 部署：**
```bash
cd ..
docker-compose up -d
```

**Kubernetes 部署：**
```bash
cd ..
./deploy-k8s.sh
```
