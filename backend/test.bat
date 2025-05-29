@echo off
chcp 65001 >nul
echo.
echo ================================
echo   数字惠农后端连接测试脚本
echo ================================
echo.

set APP_NAME=huinong-backend
set BUILD_DIR=bin
set CONFIG_FILE=configs\config.yaml

REM 检查配置文件
if not exist "%CONFIG_FILE%" (
    echo [ERROR] 配置文件不存在: %CONFIG_FILE%
    echo [INFO] 请复制 configs\config.example.yaml 为 %CONFIG_FILE% 并修改配置
    pause
    exit /b 1
)

REM 检查可执行文件是否存在，如果不存在则编译
if not exist "%BUILD_DIR%\%APP_NAME%.exe" (
    echo [INFO] 可执行文件不存在，开始编译...
    call build.bat
    if %errorlevel% neq 0 (
        echo [ERROR] 编译失败
        pause
        exit /b 1
    )
)

echo [INFO] 开始测试数据库和Dify平台连接...
echo.
echo ======== 连接测试日志 ========
echo.

REM 设置测试模式环境变量
set TEST_MODE=true

REM 运行程序进行连接测试
timeout /t 2 /nobreak >nul
%BUILD_DIR%\%APP_NAME%.exe --test-connection

echo.
echo ======== 测试完成 ========
echo.
pause 