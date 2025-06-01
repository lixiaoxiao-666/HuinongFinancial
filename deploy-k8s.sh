#!/bin/bash

# 惠农金融项目 Kubernetes 部署脚本
# 作者: AI Assistant
# 日期: $(date +%Y-%m-%d)

set -e

echo "🚀 开始部署惠农金融项目到 Kubernetes..."

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 配置
NAMESPACE="huinong-financial"
KUBECTL_TIMEOUT="300s"

# 函数：打印带颜色的消息
print_message() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 检查kubectl是否安装
check_kubectl() {
    if ! command -v kubectl &> /dev/null; then
        print_error "kubectl 未安装，请先安装 kubectl"
        exit 1
    fi
    print_success "kubectl 已安装"
}

# 检查Kubernetes集群连接
check_k8s_connection() {
    print_message "检查 Kubernetes 集群连接..."
    if kubectl cluster-info &> /dev/null; then
        print_success "Kubernetes 集群连接正常"
    else
        print_error "无法连接到 Kubernetes 集群"
        exit 1
    fi
}

# 完整部署
deploy_complete() {
    print_message "使用完整部署配置..."
    kubectl apply -f kubernetes/complete-deployment.yaml
    
    print_message "等待所有服务启动..."
    kubectl wait --for=condition=ready pod -l app=huinong-backend -n $NAMESPACE --timeout=$KUBECTL_TIMEOUT
    kubectl wait --for=condition=ready pod -l app=huinong-admin -n $NAMESPACE --timeout=$KUBECTL_TIMEOUT
    kubectl wait --for=condition=ready pod -l app=huinong-users -n $NAMESPACE --timeout=$KUBECTL_TIMEOUT
    
    print_success "完整部署完成"
}

# 部署HPA
deploy_hpa() {
    print_message "部署HPA自动扩缩容..."
    kubectl apply -f kubernetes/hpa.yaml
    print_success "HPA部署完成"
}

# 分步部署
deploy_step_by_step() {
    print_message "开始分步部署流程..."
    
    print_message "创建命名空间..."
    kubectl apply -f kubernetes/namespace.yaml
    
    print_message "部署后端服务..."
    kubectl apply -f kubernetes/backend-deployment.yaml
    
    print_message "等待后端服务启动..."
    kubectl wait --for=condition=ready pod -l app=huinong-backend -n $NAMESPACE --timeout=$KUBECTL_TIMEOUT
    
    print_message "部署前端服务..."
    kubectl apply -f kubernetes/frontend-admin-deployment.yaml
    kubectl apply -f kubernetes/frontend-users-deployment.yaml
    
    print_message "等待前端服务启动..."
    kubectl wait --for=condition=ready pod -l app=huinong-admin -n $NAMESPACE --timeout=$KUBECTL_TIMEOUT
    kubectl wait --for=condition=ready pod -l app=huinong-users -n $NAMESPACE --timeout=$KUBECTL_TIMEOUT
    
    print_message "部署Ingress配置..."
    kubectl apply -f kubernetes/ingress.yaml
    
    print_success "分步部署完成"
}

# 显示部署状态
show_status() {
    print_message "显示部署状态:"
    echo ""
    print_message "Pods 状态:"
    kubectl get pods -n $NAMESPACE -o wide
    echo ""
    print_message "Services 状态:"
    kubectl get services -n $NAMESPACE
    echo ""
    print_message "Ingress 状态:"
    kubectl get ingress -n $NAMESPACE
    echo ""
    print_message "HPA 状态:"
    kubectl get hpa -n $NAMESPACE 2>/dev/null || echo "HPA未部署"
}

# 显示访问信息
show_access_info() {
    echo ""
    print_success "🎉 部署完成！"
    echo ""
    print_message "访问信息:"
    echo "----------------------------------------"
    
    # 获取NodePort端口
    ADMIN_PORT=$(kubectl get service huinong-admin-nodeport -n $NAMESPACE -o jsonpath='{.spec.ports[0].nodePort}' 2>/dev/null || echo "30081")
    USERS_PORT=$(kubectl get service huinong-users-nodeport -n $NAMESPACE -o jsonpath='{.spec.ports[0].nodePort}' 2>/dev/null || echo "30080")
    BACKEND_PORT=$(kubectl get service huinong-backend-nodeport -n $NAMESPACE -o jsonpath='{.spec.ports[0].nodePort}' 2>/dev/null || echo "30082")
    
    # 获取集群节点IP
    NODE_IP=$(kubectl get nodes -o jsonpath='{.items[0].status.addresses[?(@.type=="ExternalIP")].address}' 2>/dev/null)
    if [ -z "$NODE_IP" ]; then
        NODE_IP=$(kubectl get nodes -o jsonpath='{.items[0].status.addresses[?(@.type=="InternalIP")].address}' 2>/dev/null)
    fi
    if [ -z "$NODE_IP" ]; then
        NODE_IP="<NODE_IP>"
    fi
    
    echo "📱 前端用户端: http://$NODE_IP:$USERS_PORT"
    echo "🔧 前端管理端: http://$NODE_IP:$ADMIN_PORT"
    echo "🔌 后端API:   http://$NODE_IP:$BACKEND_PORT"
    echo ""
    echo "如果使用 Ingress (需要配置 hosts):"
    echo "📱 前端用户端: http://huinong-users.local"
    echo "🔧 前端管理端: http://huinong-admin.local"
    echo "🔌 后端API:   http://huinong-api.local"
    echo ""
    echo "配置 hosts 文件 (Windows: C:\\Windows\\System32\\drivers\\etc\\hosts):"
    echo "$NODE_IP huinong-users.local"
    echo "$NODE_IP huinong-admin.local"
    echo "$NODE_IP huinong-api.local"
    echo "----------------------------------------"
}

# 清理部署
cleanup() {
    print_warning "清理所有部署资源..."
    kubectl delete namespace $NAMESPACE --ignore-not-found=true
    print_success "清理完成"
}

# 重启服务
restart_service() {
    local service=$1
    print_message "重启服务: $service"
    kubectl rollout restart deployment/$service -n $NAMESPACE
    kubectl rollout status deployment/$service -n $NAMESPACE --timeout=$KUBECTL_TIMEOUT
    print_success "服务 $service 重启完成"
}

# 查看日志
show_logs() {
    local service=$1
    print_message "显示服务日志: $service"
    kubectl logs -l app=$service -n $NAMESPACE --tail=50
}

# 扩缩容
scale_service() {
    local service=$1
    local replicas=$2
    print_message "扩缩容服务: $service 到 $replicas 个副本"
    kubectl scale deployment/$service --replicas=$replicas -n $NAMESPACE
    kubectl rollout status deployment/$service -n $NAMESPACE --timeout=$KUBECTL_TIMEOUT
    print_success "服务 $service 扩缩容完成"
}

# 主函数
main() {
    print_message "开始完整部署流程..."
    
    check_kubectl
    check_k8s_connection
    
    deploy_complete
    deploy_hpa
    
    show_status
    show_access_info
}

# 显示帮助信息
show_help() {
    echo "惠农金融项目 Kubernetes 部署脚本"
    echo ""
    echo "用法: $0 [选项]"
    echo ""
    echo "选项:"
    echo "  deploy          完整部署 (默认)"
    echo "  deploy-step     分步部署"
    echo "  deploy-hpa      部署HPA自动扩缩容"
    echo "  cleanup         清理所有资源"
    echo "  status          显示部署状态"
    echo "  restart <服务>  重启指定服务"
    echo "  logs <服务>     查看服务日志"
    echo "  scale <服务> <副本数>  扩缩容服务"
    echo "  help            显示帮助信息"
    echo ""
    echo "服务名称:"
    echo "  huinong-backend, huinong-admin, huinong-users"
    echo ""
    echo "示例:"
    echo "  $0                              # 完整部署"
    echo "  $0 deploy-step                  # 分步部署"
    echo "  $0 restart huinong-backend      # 重启后端服务"
    echo "  $0 logs huinong-backend         # 查看后端日志"
    echo "  $0 scale huinong-backend 5      # 扩容后端到5个副本"
    echo "  $0 cleanup                      # 清理所有资源"
}

# 参数处理
case "${1:-deploy}" in
    deploy)
        main
        ;;
    deploy-step)
        check_kubectl
        check_k8s_connection
        deploy_step_by_step
        show_status
        show_access_info
        ;;
    deploy-hpa)
        deploy_hpa
        ;;
    cleanup)
        cleanup
        ;;
    status)
        show_status
        show_access_info
        ;;
    restart)
        if [ -z "$2" ]; then
            print_error "请指定要重启的服务名称"
            show_help
            exit 1
        fi
        restart_service $2
        ;;
    logs)
        if [ -z "$2" ]; then
            print_error "请指定要查看日志的服务名称"
            show_help
            exit 1
        fi
        show_logs $2
        ;;
    scale)
        if [ -z "$2" ] || [ -z "$3" ]; then
            print_error "请指定服务名称和副本数"
            show_help
            exit 1
        fi
        scale_service $2 $3
        ;;
    help|--help|-h)
        show_help
        ;;
    *)
        print_error "未知选项: $1"
        show_help
        exit 1
        ;;
esac 