#!/bin/bash

# 数字惠农后端运行脚本 (Linux/macOS)
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
CONFIG_FILE="configs/config.yaml"

echo -e "${BLUE}================================${NC}"
echo -e "${BLUE}  数字惠农后端运行脚本 (Unix)${NC}"
echo -e "${BLUE}================================${NC}"
echo

# 检查是否需要编译
NEED_BUILD=0

# 检查可执行文件是否存在
if [ ! -f "${BUILD_DIR}/${APP_NAME}" ]; then
    echo -e "${BLUE}[INFO]${NC} 可执行文件不存在，需要编译"
    NEED_BUILD=1
else
    # 检查main.go是否比可执行文件新
    if [ "cmd/server/main.go" -nt "${BUILD_DIR}/${APP_NAME}" ]; then
        echo -e "${BLUE}[INFO]${NC} 源码已更新，需要重新编译"
        NEED_BUILD=1
    fi
fi

# 如果需要编译，先执行编译
if [ ${NEED_BUILD} -eq 1 ]; then
    echo -e "${BLUE}[INFO]${NC} 开始编译..."
    ./build.sh
    if [ $? -ne 0 ]; then
        echo -e "${RED}[ERROR]${NC} 编译失败，无法运行程序"
        exit 1
    fi
    echo
fi

# 检查配置文件
if [ ! -f "${CONFIG_FILE}" ]; then
    echo -e "${RED}[ERROR]${NC} 配置文件不存在: ${CONFIG_FILE}"
    echo -e "${BLUE}[INFO]${NC} 请先创建配置文件后再运行程序"
    exit 1
fi

# 检查可执行文件
if [ ! -f "${BUILD_DIR}/${APP_NAME}" ]; then
    echo -e "${RED}[ERROR]${NC} 可执行文件不存在: ${BUILD_DIR}/${APP_NAME}"
    echo -e "${BLUE}[INFO]${NC} 请先执行编译脚本 ./build.sh"
    exit 1
fi

# 检查可执行权限
if [ ! -x "${BUILD_DIR}/${APP_NAME}" ]; then
    echo -e "${YELLOW}[WARNING]${NC} 文件没有可执行权限，正在添加..."
    chmod +x "${BUILD_DIR}/${APP_NAME}"
fi

echo -e "${BLUE}[INFO]${NC} 启动数字惠农后端服务..."
echo -e "${BLUE}[INFO]${NC} 配置文件: ${CONFIG_FILE}"
echo -e "${BLUE}[INFO]${NC} 可执行文件: ${BUILD_DIR}/${APP_NAME}"
echo
echo -e "${GREEN}======== 服务启动日志 ========${NC}"
echo

# 设置信号处理
trap 'echo -e "\n${YELLOW}[INFO]${NC} 接收到停止信号，正在关闭服务..." && exit 0' INT TERM

# 运行程序
./${BUILD_DIR}/${APP_NAME}

# 程序退出后的处理
echo
echo -e "${YELLOW}======== 服务已停止 ========${NC}" 