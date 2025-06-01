#!/bin/bash

# 惠农金融项目 Docker 镜像构建脚本
# 作者: AI Assistant
# 日期: $(date +%Y-%m-%d)

set -e

echo "🚀 开始构建惠农金融项目 Docker 镜像..."

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 镜像标签
BACKEND_IMAGE="huinong-backend:latest"
ADMIN_IMAGE="huinong-admin:latest"
USERS_IMAGE="huinong-users:latest"

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

# 检查Docker是否安装
check_docker() {
    if ! command -v docker &> /dev/null; then
        print_error "Docker 未安装，请先安装 Docker"
        exit 1
    fi
    print_success "Docker 已安装"
}

# 构建后端镜像
build_backend() {
    print_message "构建后端镜像: $BACKEND_IMAGE"
    cd backend
    docker build -t $BACKEND_IMAGE .
    if [ $? -eq 0 ]; then
        print_success "后端镜像构建成功"
    else
        print_error "后端镜像构建失败"
        exit 1
    fi
    cd ..
}

# 构建前端Admin镜像
build_admin() {
    print_message "构建前端Admin镜像: $ADMIN_IMAGE"
    cd frontend/admin
    docker build -t $ADMIN_IMAGE .
    if [ $? -eq 0 ]; then
        print_success "前端Admin镜像构建成功"
    else
        print_error "前端Admin镜像构建失败"
        exit 1
    fi
    cd ../..
}

# 构建前端Users镜像
build_users() {
    print_message "构建前端Users镜像: $USERS_IMAGE"
    cd frontend/users
    docker build -t $USERS_IMAGE .
    if [ $? -eq 0 ]; then
        print_success "前端Users镜像构建成功"
    else
        print_error "前端Users镜像构建失败"
        exit 1
    fi
    cd ../..
}

# 显示镜像信息
show_images() {
    print_message "显示构建的镜像信息:"
    docker images | grep -E "(huinong-backend|huinong-admin|huinong-users)"
}

# 主函数
main() {
    print_message "开始构建所有镜像..."
    
    check_docker
    
    # 构建所有镜像
    build_backend
    build_admin
    build_users
    
    show_images
    
    print_success "🎉 所有镜像构建完成！"
    print_message "接下来可以运行 ./deploy-k8s.sh 部署到 Kubernetes"
}

# 如果有参数，则只构建指定的镜像
if [ $# -gt 0 ]; then
    case $1 in
        backend)
            check_docker
            build_backend
            ;;
        admin)
            check_docker
            build_admin
            ;;
        users)
            check_docker
            build_users
            ;;
        *)
            print_error "未知参数: $1"
            print_message "用法: $0 [backend|admin|users]"
            exit 1
            ;;
    esac
else
    main
fi 