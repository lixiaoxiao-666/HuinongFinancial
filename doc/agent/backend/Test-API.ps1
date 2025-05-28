# 数字惠农OA管理系统后端接口测试脚本 (PowerShell版本)
# 测试OA系统的审批管理、用户管理、系统配置等API接口

param(
    [string]$BaseUrl = "http://localhost:8080"
)

$TotalTests = 0
$PassedTests = 0
$FailedTests = 0
$AdminToken = ""
$ReviewerToken = ""

Write-Host "===================================================" -ForegroundColor Cyan
Write-Host "   数字惠农OA管理系统 后端接口测试开始" -ForegroundColor Cyan
Write-Host "===================================================" -ForegroundColor Cyan

# 测试函数
function Test-API {
    param(
        [string]$Method,
        [string]$Endpoint,
        [string]$Description,
        [string]$Data = "",
        [string]$Token = "",
        [int]$ExpectedStatus = 200
    )
    
    $global:TotalTests++
    
    Write-Host "`n[TEST $global:TotalTests] $Description" -ForegroundColor Yellow
    Write-Host "请求: $Method $Endpoint"
    
    try {
        # 构建请求头
        $Headers = @{
            'Content-Type' = 'application/json'
        }
        
        if ($Token -ne "") {
            $Headers['Authorization'] = "Bearer $Token"
        }
        
        # 构建请求参数
        $RequestParams = @{
            Uri = "$BaseUrl$Endpoint"
            Method = $Method
            Headers = $Headers
        }
        
        if ($Data -ne "" -and $Method -ne "GET") {
            $RequestParams['Body'] = $Data
        }
        
        # 发送请求
        $Response = Invoke-WebRequest @RequestParams -ErrorAction Stop
        
        if ($Response.StatusCode -eq $ExpectedStatus) {
            Write-Host "✓ 通过 (状态码: $($Response.StatusCode))" -ForegroundColor Green
            $global:PassedTests++
            
            # 显示响应内容（简化）
            $Content = $Response.Content
            if ($Content.Length -lt 500) {
                Write-Host "响应: $Content"
            } else {
                Write-Host "响应: $($Content.Substring(0, 200))..."
            }
            
            return $Response.Content
        } else {
            Write-Host "✗ 失败 (期望状态码: $ExpectedStatus, 实际状态码: $($Response.StatusCode))" -ForegroundColor Red
            $global:FailedTests++
        }
    } catch {
        $StatusCode = 0
        if ($_.Exception.Response) {
            $StatusCode = [int]$_.Exception.Response.StatusCode
        }
        
        if ($StatusCode -eq $ExpectedStatus) {
            Write-Host "✓ 通过 (状态码: $StatusCode)" -ForegroundColor Green
            $global:PassedTests++
        } else {
            Write-Host "✗ 失败 (期望状态码: $ExpectedStatus, 实际状态码: $StatusCode)" -ForegroundColor Red
            Write-Host "错误: $($_.Exception.Message)" -ForegroundColor Red
            $global:FailedTests++
        }
    }
    
    Write-Host "---------------------------------------------------"
    return ""
}

# 提取Token函数
function Get-TokenFromResponse {
    param([string]$Response)
    
    if ($Response -match '"token":"([^"]*)"') {
        return $Matches[1]
    }
    return ""
}

Write-Host "`n======================== 1. OA用户认证测试 ========================" -ForegroundColor Blue

# 1.1 管理员登录
$LoginData = @{
    username = "admin"
    password = "admin123"
} | ConvertTo-Json

$AdminLoginResponse = Test-API -Method "POST" -Endpoint "/admin/login" -Description "管理员登录" -Data $LoginData

$AdminToken = Get-TokenFromResponse $AdminLoginResponse
if ($AdminToken -ne "") {
    Write-Host "管理员Token: $($AdminToken.Substring(0, [Math]::Min(20, $AdminToken.Length)))..." -ForegroundColor Green
}

# 1.2 审批员登录
$ReviewerLoginData = @{
    username = "reviewer"
    password = "reviewer123"
} | ConvertTo-Json

$ReviewerLoginResponse = Test-API -Method "POST" -Endpoint "/admin/login" -Description "审批员登录" -Data $ReviewerLoginData

$ReviewerToken = Get-TokenFromResponse $ReviewerLoginResponse
if ($ReviewerToken -ne "") {
    Write-Host "审批员Token: $($ReviewerToken.Substring(0, [Math]::Min(20, $ReviewerToken.Length)))..." -ForegroundColor Green
}

Write-Host "`n======================== 2. OA首页/工作台测试 ========================" -ForegroundColor Blue

# 2.1 获取OA首页信息
Test-API -Method "GET" -Endpoint "/admin/dashboard" -Description "获取OA首页/工作台信息" -Token $AdminToken

Write-Host "`n======================== 3. 审批管理测试 ========================" -ForegroundColor Blue

# 3.1 获取待审批申请列表
Test-API -Method "GET" -Endpoint "/admin/loans/applications/pending" -Description "获取待审批申请列表" -Token $ReviewerToken

# 3.2 获取待审批申请列表（带筛选）
Test-API -Method "GET" -Endpoint "/admin/loans/applications/pending?status_filter=MANUAL_REVIEW_REQUIRED&page=1&limit=5" -Description "获取待审批申请列表（筛选）" -Token $ReviewerToken

Write-Host "`n======================== 4. 系统管理测试 ========================" -ForegroundColor Blue

# 4.1 获取系统统计信息
Test-API -Method "GET" -Endpoint "/admin/system/stats" -Description "获取系统统计信息" -Token $AdminToken

# 4.2 AI审批开关控制
$AIToggleData = @{
    enabled = $true
} | ConvertTo-Json

Test-API -Method "POST" -Endpoint "/admin/system/ai-approval/toggle" -Description "启用AI审批" -Data $AIToggleData -Token $AdminToken

Write-Host "`n======================== 5. 用户管理测试 ========================" -ForegroundColor Blue

# 5.1 获取OA用户列表
Test-API -Method "GET" -Endpoint "/admin/users?page=1&limit=10" -Description "获取OA用户列表" -Token $AdminToken

# 5.2 创建OA用户
$CreateUserData = @{
    username = "test_reviewer_ps"
    password = "password123"
    role = "审批员"
    display_name = "PowerShell测试审批员"
    email = "test_ps@example.com"
} | ConvertTo-Json

Test-API -Method "POST" -Endpoint "/admin/users" -Description "创建OA用户" -Data $CreateUserData -Token $AdminToken

Write-Host "`n======================== 6. 操作日志测试 ========================" -ForegroundColor Blue

# 6.1 获取操作日志
Test-API -Method "GET" -Endpoint "/admin/logs?page=1&limit=10" -Description "获取操作日志" -Token $AdminToken

Write-Host "`n======================== 7. 系统配置测试 ========================" -ForegroundColor Blue

# 7.1 获取系统配置
Test-API -Method "GET" -Endpoint "/admin/configs" -Description "获取系统配置" -Token $AdminToken

# 7.2 更新系统配置
$ConfigData = @{
    config_value = "true"
} | ConvertTo-Json

Test-API -Method "PUT" -Endpoint "/admin/configs/ai_approval_enabled" -Description "更新系统配置" -Data $ConfigData -Token $AdminToken

Write-Host "`n======================== 8. 错误处理测试 ========================" -ForegroundColor Blue

# 8.1 无效token测试
Test-API -Method "GET" -Endpoint "/admin/dashboard" -Description "无效token访问" -Token "invalid_token" -ExpectedStatus 401

# 8.2 参数错误测试
$InvalidLoginData = @{
    username = ""
    password = "short"
} | ConvertTo-Json

Test-API -Method "POST" -Endpoint "/admin/login" -Description "登录参数错误" -Data $InvalidLoginData -ExpectedStatus 400

Write-Host "`n===================================================" -ForegroundColor Cyan
Write-Host "           测试结果统计" -ForegroundColor Cyan
Write-Host "===================================================" -ForegroundColor Cyan
Write-Host "总测试数: $TotalTests"
Write-Host "通过: $PassedTests" -ForegroundColor Green
Write-Host "失败: $FailedTests" -ForegroundColor Red

$SuccessRate = if ($TotalTests -gt 0) { [Math]::Round(($PassedTests * 100 / $TotalTests), 2) } else { 0 }
Write-Host "成功率: $SuccessRate%"

if ($FailedTests -eq 0) {
    Write-Host "`n🎉 所有测试通过！OA管理系统接口运行正常！" -ForegroundColor Green
    exit 0
} else {
    Write-Host "`n❌ 有测试失败，请检查API实现" -ForegroundColor Red
    Write-Host "`n注意事项：" -ForegroundColor Yellow
    Write-Host "- 部分失败可能是因为测试环境中没有相应的数据"
    Write-Host "- 404错误通常表示使用了示例ID，在实际测试中需要使用真实ID"
    Write-Host "- 403错误表示权限控制正常工作"
    Write-Host "- 检查服务是否正常启动并连接到数据库"
    Write-Host "- 确保已创建默认的OA用户账号"
    
    Write-Host "`n默认测试账号：" -ForegroundColor Blue
    Write-Host "管理员：admin / admin123"
    Write-Host "审批员：reviewer / reviewer123"
    
    exit 1
}

Write-Host "`n===================================================" -ForegroundColor Blue
Write-Host "           OA接口功能说明" -ForegroundColor Blue
Write-Host "===================================================" -ForegroundColor Blue
Write-Host "1. 🔐 认证系统：支持OA用户登录和JWT token验证"
Write-Host "2. 📊 工作台：提供系统统计、待办事项、快捷操作"
Write-Host "3. 📋 审批管理：待审批列表、申请详情、审批决策"
Write-Host "4. ⚙️  系统管理：AI审批开关、系统统计、配置管理"
Write-Host "5. 👥 用户管理：OA用户创建、状态管理、权限控制"
Write-Host "6. 📝 操作日志：记录和查询所有操作历史"
Write-Host "7. 🔧 系统配置：灵活的配置项管理"
Write-Host "8. 🛡️  安全控制：权限验证、参数校验、错误处理"