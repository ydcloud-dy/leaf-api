#!/bin/bash

set -e

echo "🚀 开始部署 Leaf Blog 到 Kubernetes..."

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 检查 kubectl
check_kubectl() {
    if ! command -v kubectl &> /dev/null; then
        echo -e "${RED}❌ kubectl 未安装${NC}"
        exit 1
    fi
    echo -e "${GREEN}✓ kubectl 检查通过${NC}"
}

# 检查集群连接
check_cluster() {
    if ! kubectl cluster-info &> /dev/null; then
        echo -e "${RED}❌ 无法连接到 Kubernetes 集群${NC}"
        exit 1
    fi
    echo -e "${GREEN}✓ 集群连接正常${NC}"
}

# 部署基础设施（命名空间、PVC）
deploy_infrastructure() {
    echo -e "${YELLOW}📦 部署基础设施...${NC}"
    kubectl apply -f deploy/k8s/pvc.yaml
    echo -e "${GREEN}✓ 基础设施部署完成${NC}"
}

# 部署后端服务
deploy_backend() {
    echo -e "${YELLOW}🔧 部署后端服务...${NC}"
    kubectl apply -f deploy/k8s/deployment.yaml
    echo -e "${GREEN}✓ 后端服务部署完成${NC}"
}

# 部署博客前端
deploy_blog_frontend() {
    echo -e "${YELLOW}🌐 部署博客前端...${NC}"
    kubectl apply -f blog-frontend/deploy/k8s/deployment.yaml
    echo -e "${GREEN}✓ 博客前端部署完成${NC}"
}

# 部署管理后台
deploy_admin_frontend() {
    echo -e "${YELLOW}🎛 部署管理后台...${NC}"
    kubectl apply -f web/deploy/k8s/deployment.yaml
    echo -e "${GREEN}✓ 管理后台部署完成${NC}"
}

# 等待 Pod 就绪
wait_for_pods() {
    echo -e "${YELLOW}⏳ 等待 Pod 就绪...${NC}"

    echo -e "${YELLOW}等待 MySQL...${NC}"
    kubectl wait --for=condition=ready pod -l app=mysql -n leaf-blog --timeout=300s || true

    echo -e "${YELLOW}等待 Redis...${NC}"
    kubectl wait --for=condition=ready pod -l app=redis -n leaf-blog --timeout=300s || true

    echo -e "${YELLOW}等待 API...${NC}"
    kubectl wait --for=condition=ready pod -l app=leaf-api -n leaf-blog --timeout=300s || true

    echo -e "${YELLOW}等待博客前端...${NC}"
    kubectl wait --for=condition=ready pod -l app=blog-frontend -n leaf-blog --timeout=300s || true

    echo -e "${YELLOW}等待管理后台...${NC}"
    kubectl wait --for=condition=ready pod -l app=admin-frontend -n leaf-blog --timeout=300s || true

    echo -e "${GREEN}✓ 所有 Pod 已就绪${NC}"
}

# 显示部署信息
show_info() {
    echo ""
    echo -e "${GREEN}========================================${NC}"
    echo -e "${GREEN}  🎉 部署完成！${NC}"
    echo -e "${GREEN}========================================${NC}"
    echo ""

    echo -e "${YELLOW}命名空间：${NC}"
    kubectl get namespace leaf-blog
    echo ""

    echo -e "${YELLOW}Pod 状态：${NC}"
    kubectl get pods -n leaf-blog
    echo ""

    echo -e "${YELLOW}服务状态：${NC}"
    kubectl get svc -n leaf-blog
    echo ""

    echo -e "${YELLOW}Ingress 信息：${NC}"
    kubectl get ingress -n leaf-blog
    echo ""

    echo -e "${GREEN}访问地址：${NC}"
    echo -e "- API: $(kubectl get ingress leaf-api-ingress -n leaf-blog -o jsonpath='{.spec.rules[0].host}' 2>/dev/null || echo '未配置')"
    echo -e "- 博客: $(kubectl get ingress blog-frontend-ingress -n leaf-blog -o jsonpath='{.spec.rules[0].host}' 2>/dev/null || echo '未配置')"
    echo -e "- 管理后台: $(kubectl get ingress admin-frontend-ingress -n leaf-blog -o jsonpath='{.spec.rules[0].host}' 2>/dev/null || echo '未配置')"
    echo ""

    echo -e "${YELLOW}常用命令：${NC}"
    echo "  查看日志: kubectl logs -f <pod-name> -n leaf-blog"
    echo "  查看详情: kubectl describe pod <pod-name> -n leaf-blog"
    echo "  进入容器: kubectl exec -it <pod-name> -n leaf-blog -- sh"
    echo "  删除部署: kubectl delete namespace leaf-blog"
}

# 主流程
main() {
    echo "========================================"
    echo "  Leaf Blog K8s 部署脚本"
    echo "========================================"
    echo ""

    check_kubectl
    check_cluster
    deploy_infrastructure
    deploy_backend
    deploy_blog_frontend
    deploy_admin_frontend
    wait_for_pods
    show_info
}

main
