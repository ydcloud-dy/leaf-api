#!/bin/bash

##############################################################################
# Leaf Blog 一键部署脚本
# 支持: 裸部署、Docker、Docker Compose、Kubernetes
##############################################################################

set -e

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 项目根目录
PROJECT_ROOT=$(cd "$(dirname "$0")" && pwd)

# 显示欢迎信息
show_banner() {
    echo -e "${BLUE}"
    cat << "EOF"
╔══════════════════════════════════════════════════════════╗
║                                                          ║
║              Leaf Blog 一键部署脚本                      ║
║                                                          ║
║  支持: 裸部署 | Docker | Docker Compose | Kubernetes     ║
║                                                          ║
╚══════════════════════════════════════════════════════════╝
EOF
    echo -e "${NC}"
}

# 显示菜单
show_menu() {
    echo -e "${YELLOW}请选择部署方式：${NC}"
    echo "1) 裸部署 (Bare Metal) - 直接在主机上运行"
    echo "2) Docker 部署 - 单独使用 Docker 容器"
    echo "3) Docker Compose 部署 - 一键启动所有服务 (推荐)"
    echo "4) Kubernetes 部署 - K8s 集群部署"
    echo "5) 清理所有容器和数据"
    echo "0) 退出"
    echo ""
}

# 检查命令是否存在
command_exists() {
    command -v "$1" &> /dev/null
}

##############################################################################
# 1. 裸部署
##############################################################################
deploy_bare_metal() {
    echo -e "${GREEN}========================================${NC}"
    echo -e "${GREEN}  开始裸部署...${NC}"
    echo -e "${GREEN}========================================${NC}"
    echo ""

    # 检查必需的环境
    echo -e "${YELLOW}[1/7] 检查环境...${NC}"

    if ! command_exists go; then
        echo -e "${RED}❌ Go 未安装，请先安装 Go 1.21+${NC}"
        exit 1
    fi
    echo -e "${GREEN}✓ Go 环境已安装 ($(go version))${NC}"

    if ! command_exists node; then
        echo -e "${RED}❌ Node.js 未安装，请先安装 Node.js 18+${NC}"
        exit 1
    fi
    echo -e "${GREEN}✓ Node.js 环境已安装 ($(node -v))${NC}"

    if ! command_exists mysql; then
        echo -e "${YELLOW}⚠ MySQL 客户端未安装，跳过数据库检查${NC}"
    else
        if mysql -h127.0.0.1 -P3306 -uroot -p123456 -e "SELECT 1" &> /dev/null; then
            echo -e "${GREEN}✓ MySQL 连接正常${NC}"
        else
            echo -e "${YELLOW}⚠ 无法连接到 MySQL，请确保 MySQL 已启动并配置正确${NC}"
        fi
    fi

    if ! command_exists redis-cli; then
        echo -e "${YELLOW}⚠ Redis 客户端未安装，跳过 Redis 检查${NC}"
    else
        if redis-cli ping &> /dev/null; then
            echo -e "${GREEN}✓ Redis 连接正常${NC}"
        else
            echo -e "${YELLOW}⚠ 无法连接到 Redis，请确保 Redis 已启动${NC}"
        fi
    fi

    # 部署后端
    echo ""
    echo -e "${YELLOW}[2/7] 部署后端 API...${NC}"
    cd "$PROJECT_ROOT"

    if [ ! -f "config.yaml" ]; then
        echo -e "${YELLOW}配置文件不存在，请创建 config.yaml${NC}"
        exit 1
    fi

    echo "安装依赖..."
    go mod download
    echo "构建应用..."
    go build -o leaf-api .

    # 创建必要的目录
    mkdir -p logs uploads

    # 停止旧进程
    if lsof -i :8888 &> /dev/null; then
        echo "停止旧进程..."
        pkill -f leaf-api || true
        sleep 2
    fi

    # 启动后端
    nohup ./leaf-api > logs/server.log 2>&1 &
    sleep 3

    if lsof -i :8888 &> /dev/null; then
        echo -e "${GREEN}✓ 后端 API 启动成功 (http://localhost:8888)${NC}"
    else
        echo -e "${RED}❌ 后端 API 启动失败，请查看 logs/server.log${NC}"
        tail -n 20 logs/server.log
        exit 1
    fi

    # 部署博客前端
    echo ""
    echo -e "${YELLOW}[3/7] 部署博客前端...${NC}"
    cd "$PROJECT_ROOT/blog-frontend"

    echo "安装依赖..."
    npm install
    echo "构建应用..."
    npm run build

    # 使用预览服务器
    if lsof -i :4173 &> /dev/null; then
        pkill -f "vite preview" || true
        sleep 2
    fi

    nohup npm run preview > ../logs/blog-frontend.log 2>&1 &
    sleep 3

    if lsof -i :4173 &> /dev/null; then
        echo -e "${GREEN}✓ 博客前端启动成功 (http://localhost:4173)${NC}"
    else
        echo -e "${RED}❌ 博客前端启动失败${NC}"
        exit 1
    fi

    # 部署管理后台
    echo ""
    echo -e "${YELLOW}[4/7] 部署管理后台...${NC}"
    cd "$PROJECT_ROOT/web"

    echo "安装依赖..."
    npm install
    echo "构建应用..."
    npm run build

    # 使用预览服务器
    if lsof -i :4174 &> /dev/null; then
        pkill -f "vite preview.*4174" || true
        sleep 2
    fi

    nohup npm run preview -- --port 4174 > ../logs/admin-frontend.log 2>&1 &
    sleep 3

    if lsof -i :4174 &> /dev/null; then
        echo -e "${GREEN}✓ 管理后台启动成功 (http://localhost:4174)${NC}"
    else
        echo -e "${RED}❌ 管理后台启动失败${NC}"
        exit 1
    fi

    echo ""
    echo -e "${GREEN}========================================${NC}"
    echo -e "${GREEN}  🎉 裸部署完成！${NC}"
    echo -e "${GREEN}========================================${NC}"
    echo ""
    echo -e "${YELLOW}访问地址：${NC}"
    echo -e "  博客网站: ${BLUE}http://localhost:4173${NC}"
    echo -e "  管理后台: ${BLUE}http://localhost:4174${NC}"
    echo -e "  后端 API: ${BLUE}http://localhost:8888${NC}"
    echo ""
    echo -e "${YELLOW}日志文件：${NC}"
    echo "  后端: logs/server.log"
    echo "  博客前端: logs/blog-frontend.log"
    echo "  管理后台: logs/admin-frontend.log"
}

##############################################################################
# 2. Docker 部署
##############################################################################
deploy_docker() {
    echo -e "${GREEN}========================================${NC}"
    echo -e "${GREEN}  开始 Docker 部署...${NC}"
    echo -e "${GREEN}========================================${NC}"
    echo ""

    if ! command_exists docker; then
        echo -e "${RED}❌ Docker 未安装${NC}"
        exit 1
    fi

    cd "$PROJECT_ROOT"

    # 创建 Docker 网络
    echo -e "${YELLOW}[1/6] 创建 Docker 网络...${NC}"
    docker network create leaf-network 2>/dev/null || echo "网络已存在"

    # 启动 MySQL
    echo -e "${YELLOW}[2/6] 启动 MySQL...${NC}"
    docker run -d \
        --name leaf-mysql \
        --network leaf-network \
        -e MYSQL_ROOT_PASSWORD=123456 \
        -e MYSQL_DATABASE=leaf_admin \
        -e TZ=Asia/Shanghai \
        -p 3306:3306 \
        -v leaf-mysql-data:/var/lib/mysql \
        mysql:8.0 \
        || echo "MySQL 容器已存在"

    sleep 5

    # 启动 Redis
    echo -e "${YELLOW}[3/6] 启动 Redis...${NC}"
    docker run -d \
        --name leaf-redis \
        --network leaf-network \
        -p 6379:6379 \
        -v leaf-redis-data:/data \
        redis:7-alpine \
        redis-server --appendonly yes \
        || echo "Redis 容器已存在"

    sleep 3

    # 构建并启动后端
    echo -e "${YELLOW}[4/6] 构建并启动后端...${NC}"
    docker build -t leaf-api:latest .
    docker run -d \
        --name leaf-api \
        --network leaf-network \
        -p 8888:8888 \
        -v "$PROJECT_ROOT/config.yaml:/app/config.yaml:ro" \
        -v "$PROJECT_ROOT/uploads:/app/uploads" \
        -v "$PROJECT_ROOT/logs:/app/logs" \
        -e DB_HOST=leaf-mysql \
        -e REDIS_HOST=leaf-redis \
        -e TZ=Asia/Shanghai \
        leaf-api:latest \
        || echo "API 容器已存在"

    sleep 3

    # 构建并启动博客前端
    echo -e "${YELLOW}[5/6] 构建并启动博客前端...${NC}"
    cd "$PROJECT_ROOT/blog-frontend"
    docker build -t blog-frontend:latest .
    docker run -d \
        --name leaf-blog-frontend \
        --network leaf-network \
        -p 3000:80 \
        blog-frontend:latest \
        || echo "博客前端容器已存在"

    # 构建并启动管理后台
    echo -e "${YELLOW}[6/6] 构建并启动管理后台...${NC}"
    cd "$PROJECT_ROOT/web"
    docker build -t admin-frontend:latest .
    docker run -d \
        --name leaf-admin-frontend \
        --network leaf-network \
        -p 3001:80 \
        admin-frontend:latest \
        || echo "管理后台容器已存在"

    echo ""
    echo -e "${GREEN}========================================${NC}"
    echo -e "${GREEN}  🎉 Docker 部署完成！${NC}"
    echo -e "${GREEN}========================================${NC}"
    echo ""
    echo -e "${YELLOW}访问地址：${NC}"
    echo -e "  博客网站: ${BLUE}http://localhost:3000${NC}"
    echo -e "  管理后台: ${BLUE}http://localhost:3001${NC}"
    echo -e "  后端 API: ${BLUE}http://localhost:8888${NC}"
    echo ""
    echo -e "${YELLOW}常用命令：${NC}"
    echo "  查看日志: docker logs -f <container-name>"
    echo "  停止服务: docker stop <container-name>"
    echo "  删除容器: docker rm <container-name>"
}

##############################################################################
# 3. Docker Compose 部署
##############################################################################
deploy_docker_compose() {
    echo -e "${GREEN}========================================${NC}"
    echo -e "${GREEN}  开始 Docker Compose 部署...${NC}"
    echo -e "${GREEN}========================================${NC}"
    echo ""

    if ! command_exists docker-compose && ! docker compose version &> /dev/null; then
        echo -e "${RED}❌ Docker Compose 未安装${NC}"
        exit 1
    fi

    cd "$PROJECT_ROOT"

    # 检查配置文件
    if [ ! -f "config.yaml" ]; then
        echo -e "${YELLOW}配置文件不存在，请创建 config.yaml${NC}"
        exit 1
    fi

    # 使用 docker-compose 或 docker compose
    DOCKER_COMPOSE_CMD="docker-compose"
    if ! command_exists docker-compose; then
        DOCKER_COMPOSE_CMD="docker compose"
    fi

    echo -e "${YELLOW}[1/3] 停止旧服务...${NC}"
    $DOCKER_COMPOSE_CMD down 2>/dev/null || true

    echo -e "${YELLOW}[2/3] 构建镜像...${NC}"
    $DOCKER_COMPOSE_CMD build

    echo -e "${YELLOW}[3/3] 启动所有服务...${NC}"
    $DOCKER_COMPOSE_CMD up -d

    echo ""
    echo -e "${YELLOW}等待服务启动...${NC}"
    sleep 10

    echo ""
    echo -e "${GREEN}========================================${NC}"
    echo -e "${GREEN}  🎉 Docker Compose 部署完成！${NC}"
    echo -e "${GREEN}========================================${NC}"
    echo ""
    echo -e "${YELLOW}访问地址：${NC}"
    echo -e "  博客网站: ${BLUE}http://localhost:3000${NC}"
    echo -e "  管理后台: ${BLUE}http://localhost:3001${NC}"
    echo -e "  后端 API: ${BLUE}http://localhost:8888${NC}"
    echo ""
    echo -e "${YELLOW}常用命令：${NC}"
    echo "  查看状态: $DOCKER_COMPOSE_CMD ps"
    echo "  查看日志: $DOCKER_COMPOSE_CMD logs -f [service-name]"
    echo "  停止服务: $DOCKER_COMPOSE_CMD stop"
    echo "  重启服务: $DOCKER_COMPOSE_CMD restart"
    echo "  删除服务: $DOCKER_COMPOSE_CMD down"
    echo "  删除数据: $DOCKER_COMPOSE_CMD down -v"
}

##############################################################################
# 4. Kubernetes 部署
##############################################################################
deploy_kubernetes() {
    echo -e "${GREEN}========================================${NC}"
    echo -e "${GREEN}  开始 Kubernetes 部署...${NC}"
    echo -e "${GREEN}========================================${NC}"
    echo ""

    if ! command_exists kubectl; then
        echo -e "${RED}❌ kubectl 未安装${NC}"
        exit 1
    fi

    if ! kubectl cluster-info &> /dev/null; then
        echo -e "${RED}❌ 无法连接到 Kubernetes 集群${NC}"
        exit 1
    fi

    cd "$PROJECT_ROOT"

    echo -e "${YELLOW}[1/5] 创建命名空间和 PVC...${NC}"
    kubectl apply -f deploy/k8s/pvc.yaml

    echo -e "${YELLOW}[2/5] 部署后端服务...${NC}"
    kubectl apply -f deploy/k8s/deployment.yaml

    echo -e "${YELLOW}[3/5] 部署博客前端...${NC}"
    kubectl apply -f blog-frontend/deploy/k8s/deployment.yaml

    echo -e "${YELLOW}[4/5] 部署管理后台...${NC}"
    kubectl apply -f web/deploy/k8s/deployment.yaml

    echo -e "${YELLOW}[5/5] 等待 Pod 就绪...${NC}"
    kubectl wait --for=condition=ready pod -l app=mysql -n leaf-blog --timeout=300s || true
    kubectl wait --for=condition=ready pod -l app=redis -n leaf-blog --timeout=300s || true
    kubectl wait --for=condition=ready pod -l app=leaf-api -n leaf-blog --timeout=300s || true
    kubectl wait --for=condition=ready pod -l app=blog-frontend -n leaf-blog --timeout=300s || true
    kubectl wait --for=condition=ready pod -l app=admin-frontend -n leaf-blog --timeout=300s || true

    echo ""
    echo -e "${GREEN}========================================${NC}"
    echo -e "${GREEN}  🎉 Kubernetes 部署完成！${NC}"
    echo -e "${GREEN}========================================${NC}"
    echo ""

    kubectl get pods -n leaf-blog
    echo ""
    kubectl get svc -n leaf-blog
    echo ""
    kubectl get ingress -n leaf-blog
    echo ""

    echo -e "${YELLOW}访问地址（通过 Ingress）：${NC}"
    echo -e "  API: $(kubectl get ingress leaf-api-ingress -n leaf-blog -o jsonpath='{.spec.rules[0].host}' 2>/dev/null || echo '未配置')"
    echo -e "  博客: $(kubectl get ingress blog-frontend-ingress -n leaf-blog -o jsonpath='{.spec.rules[0].host}' 2>/dev/null || echo '未配置')"
    echo -e "  管理后台: $(kubectl get ingress admin-frontend-ingress -n leaf-blog -o jsonpath='{.spec.rules[0].host}' 2>/dev/null || echo '未配置')"
}

##############################################################################
# 5. 清理
##############################################################################
cleanup() {
    echo -e "${YELLOW}========================================${NC}"
    echo -e "${YELLOW}  开始清理...${NC}"
    echo -e "${YELLOW}========================================${NC}"
    echo ""

    read -p "确认要清理所有容器和数据吗? (y/N): " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        echo "取消清理"
        exit 0
    fi

    # 清理 Docker Compose
    if command_exists docker-compose || docker compose version &> /dev/null; then
        echo -e "${YELLOW}清理 Docker Compose...${NC}"
        DOCKER_COMPOSE_CMD="docker-compose"
        if ! command_exists docker-compose; then
            DOCKER_COMPOSE_CMD="docker compose"
        fi
        cd "$PROJECT_ROOT"
        $DOCKER_COMPOSE_CMD down -v 2>/dev/null || true
    fi

    # 清理 Docker 容器
    if command_exists docker; then
        echo -e "${YELLOW}清理 Docker 容器...${NC}"
        docker stop leaf-mysql leaf-redis leaf-api leaf-blog-frontend leaf-admin-frontend 2>/dev/null || true
        docker rm leaf-mysql leaf-redis leaf-api leaf-blog-frontend leaf-admin-frontend 2>/dev/null || true
        docker network rm leaf-network 2>/dev/null || true
        docker volume rm leaf-mysql-data leaf-redis-data 2>/dev/null || true
    fi

    # 清理 Kubernetes
    if command_exists kubectl && kubectl cluster-info &> /dev/null; then
        echo -e "${YELLOW}清理 Kubernetes 资源...${NC}"
        read -p "确认要删除 K8s namespace leaf-blog 吗? (y/N): " -n 1 -r
        echo
        if [[ $REPLY =~ ^[Yy]$ ]]; then
            kubectl delete namespace leaf-blog 2>/dev/null || true
        fi
    fi

    # 清理裸部署进程
    echo -e "${YELLOW}清理裸部署进程...${NC}"
    pkill -f leaf-api 2>/dev/null || true
    pkill -f "vite preview" 2>/dev/null || true

    echo ""
    echo -e "${GREEN}✓ 清理完成${NC}"
}

##############################################################################
# 主菜单
##############################################################################
main() {
    show_banner

    while true; do
        show_menu
        read -p "请输入选项 [0-5]: " choice

        case $choice in
            1)
                deploy_bare_metal
                break
                ;;
            2)
                deploy_docker
                break
                ;;
            3)
                deploy_docker_compose
                break
                ;;
            4)
                deploy_kubernetes
                break
                ;;
            5)
                cleanup
                break
                ;;
            0)
                echo "退出部署脚本"
                exit 0
                ;;
            *)
                echo -e "${RED}无效选项，请重新选择${NC}"
                ;;
        esac
    done
}

# 运行主程序
main
