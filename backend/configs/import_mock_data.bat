@echo off
chcp 65001 >nul
echo.
echo ====================================
echo   数字惠农后端 - 模拟数据导入脚本
echo ====================================
echo.

REM 设置数据库连接信息（请根据实际情况修改）
set DB_HOST=10.10.10.10
set DB_PORT=4000
set DB_USER=root
set DB_PASSWORD=
set DB_NAME=app

REM 检查MySQL客户端是否存在
mysql --version >nul 2>&1
if %errorlevel% neq 0 (
    echo [ERROR] MySQL客户端未安装或未配置到PATH
    echo [INFO] 请先安装MySQL客户端或配置环境变量
    pause
    exit /b 1
)

echo [INFO] 数据库连接信息:
echo [INFO] 主机: %DB_HOST%:%DB_PORT%
echo [INFO] 用户: %DB_USER%
echo [INFO] 数据库: %DB_NAME%
echo.

REM 检查SQL文件是否存在
if not exist "mock_data.sql" (
    echo [ERROR] 找不到 mock_data.sql 文件
    pause
    exit /b 1
)

if not exist "mock_data_part2.sql" (
    echo [ERROR] 找不到 mock_data_part2.sql 文件
    pause
    exit /b 1
)

echo [WARNING] 此操作将向数据库插入测试数据
echo [WARNING] 请确保这不是生产环境！
echo.
set /p confirm="确认要继续吗？(y/N): "
if /i not "%confirm%"=="y" (
    echo [INFO] 操作已取消
    pause
    exit /b 0
)

echo.
echo [INFO] 开始导入模拟数据...
echo.

REM 导入基础数据
echo [INFO] 正在导入基础数据 (mock_data.sql)...
mysql -h%DB_HOST% -P%DB_PORT% -u%DB_USER% %DB_NAME% < mock_data.sql
if %errorlevel% neq 0 (
    echo [ERROR] 基础数据导入失败
    pause
    exit /b 1
)
echo [SUCCESS] 基础数据导入成功

REM 导入业务数据
echo [INFO] 正在导入业务数据 (mock_data_part2.sql)...
mysql -h%DB_HOST% -P%DB_PORT% -u%DB_USER% %DB_NAME% < mock_data_part2.sql
if %errorlevel% neq 0 (
    echo [ERROR] 业务数据导入失败
    pause
    exit /b 1
)
echo [SUCCESS] 业务数据导入成功

echo.
echo [SUCCESS] 所有模拟数据导入完成！
echo.

REM 显示数据统计
echo [INFO] 正在获取数据统计...
mysql -h%DB_HOST% -P%DB_PORT% -u%DB_USER% %DB_NAME% -e "
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

echo.
echo [INFO] 模拟数据导入完成，您现在可以：
echo [INFO] 1. 启动后端服务进行测试
echo [INFO] 2. 使用测试账号登录系统
echo [INFO] 3. 测试各项业务功能
echo.
echo 测试账号信息：
echo 普通用户: 13900000001 (李大牛)
echo 管理员: admin (admin@huinong.com)
echo 密码: 请查看数据库中的密码哈希
echo.
pause 