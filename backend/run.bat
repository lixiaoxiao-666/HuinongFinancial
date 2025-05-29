@echo off
chcp 65001 >nul
setlocal enabledelayedexpansion

echo.
echo ================================
echo   数字惠农后端运行脚本 (Windows)
echo ================================
echo.

set APP_NAME=huinong-backend
set BUILD_DIR=bin
set CONFIG_FILE=configs\config.yaml

REM 检查是否需要编译
set NEED_BUILD=0

REM 检查可执行文件是否存在
if not exist "%BUILD_DIR%\%APP_NAME%.exe" (
    echo [INFO] 可执行文件不存在，需要编译
    set NEED_BUILD=1
) else (
    REM 检查main.go是否比可执行文件新
    for %%i in ("cmd\server\main.go") do set MAIN_TIME=%%~ti
    for %%i in ("%BUILD_DIR%\%APP_NAME%.exe") do set EXE_TIME=%%~ti
    
    REM 简单比较：如果main.go时间戳字符串大于exe时间戳，则需要重新编译
    if "!MAIN_TIME!" gtr "!EXE_TIME!" (
        echo [INFO] 源码已更新，需要重新编译
        set NEED_BUILD=1
    )
)

REM 如果需要编译，先执行编译
if !NEED_BUILD! equ 1 (
    echo [INFO] 开始编译...
    call build.bat
    if !errorlevel! neq 0 (
        echo [ERROR] 编译失败，无法运行程序
        pause
        exit /b 1
    )
    echo.
)

REM 检查配置文件
if not exist "%CONFIG_FILE%" (
    echo [ERROR] 配置文件不存在: %CONFIG_FILE%
    echo [INFO] 请先创建配置文件后再运行程序
    pause
    exit /b 1
)

REM 检查可执行文件
if not exist "%BUILD_DIR%\%APP_NAME%.exe" (
    echo [ERROR] 可执行文件不存在: %BUILD_DIR%\%APP_NAME%.exe
    echo [INFO] 请先执行编译脚本 build.bat
    pause
    exit /b 1
)

echo [INFO] 启动数字惠农后端服务...
echo [INFO] 配置文件: %CONFIG_FILE%
echo [INFO] 可执行文件: %BUILD_DIR%\%APP_NAME%.exe
echo.
echo ======== 服务启动日志 ========
echo.

REM 运行程序
%BUILD_DIR%\%APP_NAME%.exe

REM 程序退出后的处理
echo.
echo ======== 服务已停止 ========
echo.
echo 按任意键退出...
pause >nul 