#!/bin/bash

# 惠农金融项目一键部署脚本
# 作者: AI Assistant
# 日期: $(date +%Y-%m-%d)

set -e

echo "🚀 惠农金融项目一键部署开始..."

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

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

# 检查必要工具
check_tools() {
    print_message "检查必要工具..."
    
    if ! command -v docker &> /dev/null; then
        print_error "Docker 未安装，请先安装 Docker"
        exit 1
    fi
    
    if ! command -v kubectl &> /dev/null; then
        print_error "kubectl 未安装，请先安装 kubectl"
        exit 1
    fi
    
    print_success "所有必要工具已安装"
}

# 构建镜像
build_images() {
    print_message "开始构建Docker镜像..."
    
    if [ -f "./build-images.sh" ]; then
        ./build-images.sh
    else
        print_error "构建脚本不存在"
        exit 1
    fi
    
    print_success "镜像构建完成"
}

# 部署到Kubernetes
deploy_to_k8s() {
    print_message "开始部署到Kubernetes..."
    
    if [ -f "./deploy-k8s.sh" ]; then
        ./deploy-k8s.sh deploy
    else
        print_error "部署脚本不存在"
        exit 1
    fi
    
    print_success "Kubernetes部署完成"
}

# 本地Docker Compose部署
deploy_local() {
    print_message "开始本地Docker Compose部署..."
    
    if [ -f "./docker-compose.yml" ]; then
        docker-compose down 2>/dev/null || true
        docker-compose up -d --build
        
        print_message "等待服务启动..."
        sleep 10
        
        print_message "检查服务状态..."
        docker-compose ps
        
        print_success "本地部署完成"
        print_message "访问地址:"
        echo "📱 前端用户端: http://localhost:3000"
        echo "🔧 前端管理端: http://localhost:3001"
        echo "🔌 后端API:   http://localhost:8080"
    else
        print_error "docker-compose.yml 不存在"
        exit 1
    fi
}

# 显示帮助信息
show_help() {
    echo "惠农金融项目一键部署脚本"
    echo ""
    echo "用法: $0 [选项]"
    echo ""
    echo "选项:"
    echo "  k8s             构建镜像并部署到Kubernetes (默认)"
    echo "  local           本地Docker Compose部署"
    echo "  build-only      仅构建镜像"
    echo "  deploy-only     仅部署到Kubernetes (需要镜像已存在)"
    echo "  help            显示帮助信息"
    echo ""
    echo "示例:"
    echo "  $0              # 完整K8s部署"
    echo "  $0 local        # 本地部署"
    echo "  $0 build-only   # 仅构建镜像"
}

# 主函数
main() {
    print_message "开始完整部署流程..."
    
    check_tools
    build_images
    deploy_to_k8s
    
    print_success "🎉 一键部署完成！"
}

# 参数处理
case "${1:-k8s}" in
    k8s)
        main
        ;;
    local)
        check_tools
        deploy_local
        ;;
    build-only)
        check_tools
        build_images
        ;;
    deploy-only)
        check_tools
        deploy_to_k8s
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