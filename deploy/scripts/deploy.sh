#!/bin/bash

set -e

echo "🚀 开始部署 Leaf API..."

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 检查 Go 环境
check_go() {
    if ! command -v go &> /dev/null; then
        echo -e "${RED}❌ Go 未安装，请先安装 Go 1.21+${NC}"
        exit 1
    fi
    echo -e "${GREEN}✓ Go 环境检查通过${NC}"
}

# 检查 MySQL
check_mysql() {
    if ! command -v mysql &> /dev/null; then
        echo -e "${YELLOW}⚠ MySQL 客户端未安装${NC}"
        return
    fi

    echo -e "${YELLOW}正在检查 MySQL 连接...${NC}"
    if mysql -h127.0.0.1 -P3306 -uroot -p123456 -e "SELECT 1" &> /dev/null; then
        echo -e "${GREEN}✓ MySQL 连接成功${NC}"
    else
        echo -e "${RED}❌ 无法连接到 MySQL，请确保 MySQL 已启动${NC}"
        exit 1
    fi
}

# 检查 Redis
check_redis() {
    if ! command -v redis-cli &> /dev/null; then
        echo -e "${YELLOW}⚠ Redis 客户端未安装${NC}"
        return
    fi

    echo -e "${YELLOW}正在检查 Redis 连接...${NC}"
    if redis-cli -h 127.0.0.1 -p 6379 ping &> /dev/null; then
        echo -e "${GREEN}✓ Redis 连接成功${NC}"
    else
        echo -e "${RED}❌ 无法连接到 Redis，请确保 Redis 已启动${NC}"
        exit 1
    fi
}

# 安装依赖
install_deps() {
    echo -e "${YELLOW}📦 正在安装依赖...${NC}"
    go mod download
    echo -e "${GREEN}✓ 依赖安装完成${NC}"
}

# 构建应用
build_app() {
    echo -e "${YELLOW}🔨 正在构建应用...${NC}"
    go build -o leaf-api .
    echo -e "${GREEN}✓ 构建完成${NC}"
}

# 创建必要的目录
create_dirs() {
    echo -e "${YELLOW}📁 正在创建必要的目录...${NC}"
    mkdir -p logs uploads
    echo -e "${GREEN}✓ 目录创建完成${NC}"
}

# 启动应用
start_app() {
    echo -e "${YELLOW}🎯 正在启动应用...${NC}"

    # 检查端口是否被占用
    if lsof -i :8888 &> /dev/null; then
        echo -e "${YELLOW}端口 8888 已被占用，正在停止旧进程...${NC}"
        pkill -f leaf-api || true
        sleep 2
    fi

    # 启动应用
    nohup ./leaf-api > server.log 2>&1 &

    # 等待启动
    sleep 3

    # 检查是否启动成功
    if lsof -i :8888 &> /dev/null; then
        echo -e "${GREEN}✅ 应用启动成功！${NC}"
        echo -e "${GREEN}访问地址: http://localhost:8888${NC}"
    else
        echo -e "${RED}❌ 应用启动失败，请查看 server.log${NC}"
        tail -n 20 server.log
        exit 1
    fi
}

# 主流程
main() {
    echo "========================================"
    echo "  Leaf API 部署脚本"
    echo "========================================"

    check_go
    check_mysql
    check_redis
    install_deps
    build_app
    create_dirs
    start_app

    echo -e "${GREEN}========================================"
    echo "  🎉 部署完成！"
    echo "========================================${NC}"
}

main
