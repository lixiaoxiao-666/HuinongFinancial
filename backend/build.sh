#!/bin/bash

# 数字惠农后端编译脚本 (Linux/macOS)
set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 配置变量
APP_NAME="huinong-backend"
BUILD_DIR="bin"
BUILD_TIME=$(date '+%Y-%m-%d %H:%M:%S')
GO_VERSION=$(go version 2>/dev/null || echo "未安装")

echo -e "${BLUE}================================${NC}"
echo -e "${BLUE}  数字惠农后端编译脚本 (Unix)${NC}"
echo -e "${BLUE}================================${NC}"
echo

echo -e "${BLUE}[INFO]${NC} 编译环境信息"
echo -e "${BLUE}[INFO]${NC} Go版本: ${GO_VERSION}"
echo -e "${BLUE}[INFO]${NC} 编译时间: ${BUILD_TIME}"
echo -e "${BLUE}[INFO]${NC} 输出目录: ${BUILD_DIR}"
echo

# 检查Go环境
if ! command -v go &> /dev/null; then
    echo -e "${RED}[ERROR]${NC} Go环境未安装或未配置到PATH"
    exit 1
fi

# 创建输出目录
if [ ! -d "${BUILD_DIR}" ]; then
    mkdir -p "${BUILD_DIR}"
    echo -e "${BLUE}[INFO]${NC} 创建输出目录: ${BUILD_DIR}"
fi

# 检查go.mod
if [ ! -f "go.mod" ]; then
    echo -e "${RED}[ERROR]${NC} go.mod文件不存在，请确认在正确的项目根目录"
    exit 1
fi

# 下载依赖
echo -e "${BLUE}[INFO]${NC} 下载Go模块依赖..."
go mod tidy
if [ $? -ne 0 ]; then
    echo -e "${RED}[ERROR]${NC} 下载依赖失败"
    exit 1
fi

# 检查main.go
if [ ! -f "cmd/server/main.go" ]; then
    echo -e "${RED}[ERROR]${NC} 找不到主程序文件: cmd/server/main.go"
    exit 1
fi

# 设置编译参数
export CGO_ENABLED=0

# 检测操作系统
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case "${ARCH}" in
    x86_64)
        ARCH="amd64"
        ;;
    aarch64|arm64)
        ARCH="arm64"
        ;;
    i386|i686)
        ARCH="386"
        ;;
esac

case "${OS}" in
    linux*)
        GOOS="linux"
        BINARY_NAME="${APP_NAME}"
        ;;
    darwin*)
        GOOS="darwin"
        BINARY_NAME="${APP_NAME}"
        ;;
    *)
        echo -e "${YELLOW}[WARNING]${NC} 未知操作系统: ${OS}，使用默认配置"
        GOOS="linux"
        BINARY_NAME="${APP_NAME}"
        ;;
esac

export GOOS GOARCH=${ARCH}

# 编译程序
echo -e "${BLUE}[INFO]${NC} 开始编译..."
echo -e "${BLUE}[INFO]${NC} 目标平台: ${GOOS}/${ARCH}"

go build -ldflags "-s -w -X 'main.BuildTime=${BUILD_TIME}'" -o "${BUILD_DIR}/${BINARY_NAME}" ./cmd/server

if [ $? -ne 0 ]; then
    echo -e "${RED}[ERROR]${NC} 编译失败"
    exit 1
fi

# 设置可执行权限
chmod +x "${BUILD_DIR}/${BINARY_NAME}"

echo -e "${GREEN}[SUCCESS]${NC} 编译成功！"
echo -e "${BLUE}[INFO]${NC} 可执行文件: ${BUILD_DIR}/${BINARY_NAME}"
echo

# 检查配置文件
if [ ! -f "configs/config.yaml" ]; then
    echo -e "${YELLOW}[WARNING]${NC} 配置文件不存在: configs/config.yaml"
    echo -e "${BLUE}[INFO]${NC} 请确保配置文件存在后再运行程序"
    echo
fi

echo -e "${BLUE}[INFO]${NC} 编译完成，可以运行以下命令启动服务:"
echo -e "${GREEN}./${BUILD_DIR}/${BINARY_NAME}${NC}"
echo 