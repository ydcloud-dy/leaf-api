#!/bin/bash

set -e

echo "🚀 开始部署 Admin Frontend..."

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 检查 Node.js 环境
check_node() {
    if ! command -v node &> /dev/null; then
        echo -e "${RED}❌ Node.js 未安装，请先安装 Node.js 16+${NC}"
        exit 1
    fi

    NODE_VERSION=$(node -v | cut -d'v' -f2 | cut -d'.' -f1)
    if [ "$NODE_VERSION" -lt 16 ]; then
        echo -e "${RED}❌ Node.js 版本过低，需要 16+${NC}"
        exit 1
    fi

    echo -e "${GREEN}✓ Node.js 环境检查通过 ($(node -v))${NC}"
}

# 检查 npm
check_npm() {
    if ! command -v npm &> /dev/null; then
        echo -e "${RED}❌ npm 未安装${NC}"
        exit 1
    fi
    echo -e "${GREEN}✓ npm 环境检查通过 ($(npm -v))${NC}"
}

# 安装依赖
install_deps() {
    echo -e "${YELLOW}📦 正在安装依赖...${NC}"
    npm install
    echo -e "${GREEN}✓ 依赖安装完成${NC}"
}

# 构建应用
build_app() {
    echo -e "${YELLOW}🔨 正在构建应用...${NC}"
    npm run build
    echo -e "${GREEN}✓ 构建完成${NC}"
}

# 检查 Nginx
check_nginx() {
    if ! command -v nginx &> /dev/null; then
        echo -e "${YELLOW}⚠ Nginx 未安装，将跳过部署步骤${NC}"
        echo -e "${YELLOW}请手动将 dist 目录复制到 Web 服务器${NC}"
        return 1
    fi
    echo -e "${GREEN}✓ Nginx 检查通过${NC}"
    return 0
}

# 部署到 Nginx
deploy_to_nginx() {
    if check_nginx; then
        echo -e "${YELLOW}📂 正在部署到 Nginx...${NC}"

        # 备份旧版本
        if [ -d "/usr/share/nginx/html-admin.backup" ]; then
            rm -rf /usr/share/nginx/html-admin.backup
        fi

        if [ -d "/usr/share/nginx/html-admin" ]; then
            sudo mv /usr/share/nginx/html-admin /usr/share/nginx/html-admin.backup
        fi

        # 部署新版本
        sudo mkdir -p /usr/share/nginx/html-admin
        sudo cp -r dist/* /usr/share/nginx/html-admin/

        # 复制 Nginx 配置
        if [ -f "deploy/nginx/nginx.conf" ]; then
            sudo cp deploy/nginx/nginx.conf /etc/nginx/conf.d/admin-frontend.conf
        fi

        # 重启 Nginx
        sudo nginx -t && sudo systemctl restart nginx

        echo -e "${GREEN}✓ 部署到 Nginx 完成${NC}"
    fi
}

# 启动预览服务器（如果没有 Nginx）
start_preview() {
    if ! check_nginx; then
        echo -e "${YELLOW}🎯 正在启动预览服务器...${NC}"

        # 检查端口是否被占用
        if lsof -i :4174 &> /dev/null; then
            echo -e "${YELLOW}端口 4174 已被占用，正在停止旧进程...${NC}"
            pkill -f "vite preview" || true
            sleep 2
        fi

        # 启动预览服务器
        nohup npm run preview -- --port 4174 > admin-frontend.log 2>&1 &

        sleep 3

        # 检查是否启动成功
        if lsof -i :4174 &> /dev/null; then
            echo -e "${GREEN}✅ 预览服务器启动成功！${NC}"
            echo -e "${GREEN}访问地址: http://localhost:4174${NC}"
        else
            echo -e "${RED}❌ 预览服务器启动失败，请查看 admin-frontend.log${NC}"
            exit 1
        fi
    fi
}

# 主流程
main() {
    echo "========================================"
    echo "  Admin Frontend 部署脚本"
    echo "========================================"

    check_node
    check_npm
    install_deps
    build_app

    # 尝试部署到 Nginx，如果没有则启动预览服务器
    deploy_to_nginx || start_preview

    echo -e "${GREEN}========================================"
    echo "  🎉 部署完成！"
    echo "========================================${NC}"
}

main
