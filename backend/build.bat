@echo off
chcp 65001 >nul
echo.
echo ================================
echo   数字惠农后端编译脚本 (Windows)
echo ================================
echo.

set APP_NAME=huinong-backend
set BUILD_DIR=bin
set BUILD_TIME=%date% %time%
set GO_VERSION=
for /f %%i in ('go version') do set GO_VERSION=%%i

echo [INFO] 编译环境信息
echo [INFO] Go版本: %GO_VERSION%
echo [INFO] 编译时间: %BUILD_TIME%
echo [INFO] 输出目录: %BUILD_DIR%
echo.

REM 检查Go环境
go version >nul 2>&1
if %errorlevel% neq 0 (
    echo [ERROR] Go环境未安装或未配置到PATH
    pause
    exit /b 1
)

REM 创建输出目录
if not exist "%BUILD_DIR%" (
    mkdir "%BUILD_DIR%"
    echo [INFO] 创建输出目录: %BUILD_DIR%
)

REM 检查go.mod
if not exist "go.mod" (
    echo [ERROR] go.mod文件不存在，请确认在正确的项目根目录
    pause
    exit /b 1
)

REM 下载依赖
echo [INFO] 下载Go模块依赖...
go mod tidy
if %errorlevel% neq 0 (
    echo [ERROR] 下载依赖失败
    pause
    exit /b 1
)

REM 检查main.go
if not exist "cmd\server\main.go" (
    echo [ERROR] 找不到主程序文件: cmd\server\main.go
    pause
    exit /b 1
)

REM 设置编译参数
set CGO_ENABLED=0
set GOOS=windows
set GOARCH=amd64

REM 编译程序
echo [INFO] 开始编译...
go build -ldflags "-s -w -X 'main.BuildTime=%BUILD_TIME%'" -o "%BUILD_DIR%\%APP_NAME%.exe" ./cmd/server
if %errorlevel% neq 0 (
    echo [ERROR] 编译失败
    pause
    exit /b 1
)

echo [SUCCESS] 编译成功！
echo [INFO] 可执行文件: %BUILD_DIR%\%APP_NAME%.exe
echo.

REM 检查配置文件
if not exist "configs\config.yaml" (
    echo [WARNING] 配置文件不存在: configs\config.yaml
    echo [INFO] 请确保配置文件存在后再运行程序
    echo.
)

echo [INFO] 编译完成，可以运行以下命令启动服务:
echo [INFO] %BUILD_DIR%\%APP_NAME%.exe
echo.
echo 按任意键退出...
pause >nul 