#!/bin/bash

# 数字惠农后端 - 模拟数据导入脚本 (Linux/macOS版本)

set -e  # 遇到错误立即退出

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 打印带颜色的消息
print_info() {
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

echo
echo "===================================="
echo "   数字惠农后端 - 模拟数据导入脚本"
echo "===================================="
echo

# 数据库连接信息（请根据实际情况修改）
DB_HOST="10.10.10.10"
DB_PORT="4000"
DB_USER="root"
DB_PASSWORD=""
DB_NAME="app"

# 检查MySQL客户端是否存在
if ! command -v mysql &> /dev/null; then
    print_error "MySQL客户端未安装或未配置到PATH"
    print_info "请先安装MySQL客户端"
    print_info "Ubuntu/Debian: sudo apt-get install mysql-client"
    print_info "CentOS/RHEL: sudo yum install mysql"
    print_info "macOS: brew install mysql-client"
    exit 1
fi

print_info "数据库连接信息:"
print_info "主机: ${DB_HOST}:${DB_PORT}"
print_info "用户: ${DB_USER}"
print_info "数据库: ${DB_NAME}"
echo

# 检查SQL文件是否存在
if [ ! -f "mock_data.sql" ]; then
    print_error "找不到 mock_data.sql 文件"
    print_info "请确保在 backend/configs 目录下运行此脚本"
    exit 1
fi

if [ ! -f "mock_data_part2.sql" ]; then
    print_error "找不到 mock_data_part2.sql 文件"
    print_info "请确保在 backend/configs 目录下运行此脚本"
    exit 1
fi

print_warning "此操作将向数据库插入测试数据"
print_warning "请确保这不是生产环境！"
echo

# 用户确认
read -p "确认要继续吗？(y/N): " -n 1 -r
echo
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    print_info "操作已取消"
    exit 0
fi

echo
print_info "开始导入模拟数据..."
echo

# 构建MySQL连接命令
if [ -z "$DB_PASSWORD" ]; then
    MYSQL_CMD="mysql -h${DB_HOST} -P${DB_PORT} -u${DB_USER} ${DB_NAME}"
else
    MYSQL_CMD="mysql -h${DB_HOST} -P${DB_PORT} -u${DB_USER} -p${DB_PASSWORD} ${DB_NAME}"
fi

# 导入基础数据
print_info "正在导入基础数据 (mock_data.sql)..."
if $MYSQL_CMD < mock_data.sql; then
    print_success "基础数据导入成功"
else
    print_error "基础数据导入失败"
    exit 1
fi

# 导入业务数据
print_info "正在导入业务数据 (mock_data_part2.sql)..."
if $MYSQL_CMD < mock_data_part2.sql; then
    print_success "业务数据导入成功"
else
    print_error "业务数据导入失败"
    exit 1
fi

echo
print_success "所有模拟数据导入完成！"
echo

# 显示数据统计
print_info "正在获取数据统计..."
$MYSQL_CMD -e "
SELECT '=== 数据导入统计 ===' as info;
SELECT 'users' as table_name, COUNT(*) as record_count FROM users
UNION ALL SELECT 'oa_users', COUNT(*) FROM oa_users
UNION ALL SELECT 'loan_products', COUNT(*) FROM loan_products
UNION ALL SELECT 'loan_applications', COUNT(*) FROM loan_applications
UNION ALL SELECT 'machines', COUNT(*) FROM machines
UNION ALL SELECT 'rental_orders', COUNT(*) FROM rental_orders
UNION ALL SELECT 'articles', COUNT(*) FROM articles
UNION ALL SELECT 'experts', COUNT(*) FROM experts
UNION ALL SELECT 'categories', COUNT(*) FROM categories;
SELECT '=== 导入完成 ===' as info;
"

echo
print_info "模拟数据导入完成，您现在可以："
print_info "1. 启动后端服务进行测试"
print_info "2. 使用测试账号登录系统"
print_info "3. 测试各项业务功能"
echo
echo "测试账号信息："
echo "普通用户: 13900000001 (李大牛)"
echo "管理员: admin (admin@huinong.com)"
echo "密码: 请查看数据库中的密码哈希"
echo

print_success "脚本执行完成！" 