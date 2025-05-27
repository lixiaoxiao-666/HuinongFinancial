# 数字惠农APP后端服务接口测试脚本 (PowerShell版本)
# 测试所有已实现的API接口

$BASE_URL = "http://localhost:8080"
$API_BASE = "$BASE_URL/api/v1"

Write-Host "===================================================" -ForegroundColor Yellow
Write-Host "   数字惠农APP后端服务接口测试开始" -ForegroundColor Yellow
Write-Host "===================================================" -ForegroundColor Yellow

# 测试结果统计
$TOTAL_TESTS = 0
$PASSED_TESTS = 0
$FAILED_TESTS = 0

# 存储token的变量
$USER_TOKEN = ""
$ADMIN_TOKEN = ""

# 测试函数
function Test-API {
    param(
        [string]$Method,
        [string]$Endpoint,
        [string]$Description,
        [object]$Body = $null,
        [string]$AuthHeader = "",
        [int]$ExpectedStatus = 200
    )
    
    $Global:TOTAL_TESTS++
    
    Write-Host "`n[TEST $Global:TOTAL_TESTS] $Description" -ForegroundColor Yellow
    Write-Host "请求: $Method $Endpoint"
    
    try {
        # 构建请求头
        $headers = @{
            "Content-Type" = "application/json"
        }
        
        if ($AuthHeader -ne "") {
            $headers["Authorization"] = $AuthHeader
        }
        
        # 发送请求
        if ($Body -ne $null) {
            $jsonBody = $Body | ConvertTo-Json -Depth 10
            $response = Invoke-WebRequest -Uri $Endpoint -Method $Method -Headers $headers -Body $jsonBody -ErrorAction Stop
        } else {
            $response = Invoke-WebRequest -Uri $Endpoint -Method $Method -Headers $headers -ErrorAction Stop
        }
        
        # 检查状态码
        if ($response.StatusCode -eq $ExpectedStatus) {
            Write-Host "✓ 通过" -ForegroundColor Green -NoNewline
            Write-Host " (状态码: $($response.StatusCode))"
            $Global:PASSED_TESTS++
            
            # 显示响应内容
            $responseContent = $response.Content
            if ($responseContent.Length -gt 200) {
                Write-Host "响应: $($responseContent.Substring(0, 200))..."
            } else {
                Write-Host "响应: $responseContent"
            }
            
            return $response.Content | ConvertFrom-Json
        } else {
            Write-Host "✗ 失败" -ForegroundColor Red -NoNewline
            Write-Host " (期望状态码: $ExpectedStatus, 实际状态码: $($response.StatusCode))"
            Write-Host "响应: $($response.Content)"
            $Global:FAILED_TESTS++
            return $null
        }
    }
    catch {
        $statusCode = $_.Exception.Response.StatusCode.value__
        if ($statusCode -eq $ExpectedStatus) {
            Write-Host "✓ 通过" -ForegroundColor Green -NoNewline
            Write-Host " (状态码: $statusCode)"
            $Global:PASSED_TESTS++
        } else {
            Write-Host "✗ 失败" -ForegroundColor Red -NoNewline
            Write-Host " (期望状态码: $ExpectedStatus, 实际状态码: $statusCode)"
            Write-Host "错误: $($_.Exception.Message)"
            $Global:FAILED_TESTS++
        }
        return $null
    }
    finally {
        Write-Host "---------------------------------------------------"
    }
}

# 文件上传测试函数
function Test-FileUpload {
    param(
        [string]$Endpoint,
        [string]$Description,
        [string]$AuthHeader,
        [int]$ExpectedStatus = 200
    )
    
    $Global:TOTAL_TESTS++
    
    Write-Host "`n[TEST $Global:TOTAL_TESTS] $Description" -ForegroundColor Yellow
    Write-Host "请求: POST $Endpoint"
    
    try {
        # 创建临时测试文件
        $tempFile = [System.IO.Path]::GetTempFileName()
        "这是一个测试文件" | Out-File -FilePath $tempFile -Encoding UTF8
        
        # 构建multipart/form-data请求
        $boundary = [System.Guid]::NewGuid().ToString()
        $headers = @{
            "Authorization" = $AuthHeader
            "Content-Type" = "multipart/form-data; boundary=$boundary"
        }
        
        # 构建multipart内容
        $fileContent = [System.IO.File]::ReadAllBytes($tempFile)
        $fileName = "test_upload.txt"
        
        $bodyLines = @()
        $bodyLines += "--$boundary"
        $bodyLines += "Content-Disposition: form-data; name=`"file`"; filename=`"$fileName`""
        $bodyLines += "Content-Type: text/plain"
        $bodyLines += ""
        $bodyLines += [System.Text.Encoding]::UTF8.GetString($fileContent)
        $bodyLines += "--$boundary"
        $bodyLines += "Content-Disposition: form-data; name=`"purpose`""
        $bodyLines += ""
        $bodyLines += "loan_document"
        $bodyLines += "--$boundary--"
        
        $body = $bodyLines -join "`r`n"
        $bodyBytes = [System.Text.Encoding]::UTF8.GetBytes($body)
        
        # 发送请求
        $response = Invoke-WebRequest -Uri $Endpoint -Method POST -Headers $headers -Body $bodyBytes -ErrorAction Stop
        
        # 清理临时文件
        Remove-Item $tempFile -Force
        
        if ($response.StatusCode -eq $ExpectedStatus) {
            Write-Host "✓ 通过" -ForegroundColor Green -NoNewline
            Write-Host " (状态码: $($response.StatusCode))"
            $Global:PASSED_TESTS++
            Write-Host "响应: $($response.Content)"
            return $response.Content | ConvertFrom-Json
        } else {
            Write-Host "✗ 失败" -ForegroundColor Red -NoNewline
            Write-Host " (期望状态码: $ExpectedStatus, 实际状态码: $($response.StatusCode))"
            $Global:FAILED_TESTS++
            return $null
        }
    }
    catch {
        Write-Host "✗ 失败" -ForegroundColor Red -NoNewline
        Write-Host " - 错误: $($_.Exception.Message)"
        $Global:FAILED_TESTS++
        return $null
    }
    finally {
        Write-Host "---------------------------------------------------"
    }
}

Write-Host "`n=== 1. 健康检查接口测试 ===" -ForegroundColor Blue

# 1. 健康检查
Test-API -Method "GET" -Endpoint "$BASE_URL/health" -Description "健康检查接口"

Write-Host "`n=== 2. 用户服务接口测试 ===" -ForegroundColor Blue

# 2. 发送验证码
$phoneData = @{
    phone = "13800138000"
}
Test-API -Method "POST" -Endpoint "$API_BASE/users/send-verification-code" -Description "发送验证码" -Body $phoneData

# 3. 用户注册
$registerData = @{
    phone = "13800138001"
    password = "test123456"
    verification_code = "123456"
}
$registerResponse = Test-API -Method "POST" -Endpoint "$API_BASE/users/register" -Description "用户注册" -Body $registerData -ExpectedStatus 201

# 4. 用户登录
Write-Host "`n正在进行用户登录获取Token..." -ForegroundColor Yellow
$loginData = @{
    phone = "13800138001"
    password = "test123456"
}
$loginResponse = Test-API -Method "POST" -Endpoint "$API_BASE/users/login" -Description "用户登录" -Body $loginData

# 提取用户Token
if ($loginResponse -and $loginResponse.data.token) {
    $USER_TOKEN = "Bearer " + $loginResponse.data.token
    Write-Host "✓ 成功获取用户Token: $($loginResponse.data.token.Substring(0, 20))..." -ForegroundColor Green
} else {
    Write-Host "✗ 未能获取用户Token" -ForegroundColor Red
}

# 5. 获取用户信息（需要认证）
if ($USER_TOKEN -ne "") {
    Test-API -Method "GET" -Endpoint "$API_BASE/users/me" -Description "获取用户信息" -AuthHeader $USER_TOKEN
}

# 6. 更新用户信息（需要认证）
if ($USER_TOKEN -ne "") {
    $updateData = @{
        nickname = "测试农户"
        real_name = "张三"
        address = "测试省测试市测试村"
    }
    Test-API -Method "PUT" -Endpoint "$API_BASE/users/me" -Description "更新用户信息" -Body $updateData -AuthHeader $USER_TOKEN
}

Write-Host "`n=== 3. 贷款服务接口测试 ===" -ForegroundColor Blue

# 7. 获取贷款产品列表
Test-API -Method "GET" -Endpoint "$API_BASE/loans/products" -Description "获取贷款产品列表"

# 8. 按分类查询贷款产品
Test-API -Method "GET" -Endpoint "$API_BASE/loans/products?category=种植贷" -Description "按分类查询贷款产品"

# 9. 获取贷款产品详情
Test-API -Method "GET" -Endpoint "$API_BASE/loans/products/loan_prod_001" -Description "获取贷款产品详情"

# 10. 提交贷款申请（需要认证）
if ($USER_TOKEN -ne "") {
    $loanApplicationData = @{
        product_id = "loan_prod_001"
        amount = 30000
        term_months = 12
        purpose = "购买化肥和种子"
        applicant_info = @{
            real_name = "张三"
            id_card_number = "310123456789012345"
            address = "测试省测试市测试村"
        }
        uploaded_documents = @()
    }
    Test-API -Method "POST" -Endpoint "$API_BASE/loans/applications" -Description "提交贷款申请" -Body $loanApplicationData -AuthHeader $USER_TOKEN -ExpectedStatus 201
}

# 11. 获取我的贷款申请列表（需要认证）
if ($USER_TOKEN -ne "") {
    Test-API -Method "GET" -Endpoint "$API_BASE/loans/applications/my" -Description "获取我的贷款申请列表" -AuthHeader $USER_TOKEN
}

# 12. 分页查询我的贷款申请
if ($USER_TOKEN -ne "") {
    $pageUrl = "$API_BASE/loans/applications/my?page=1`&limit=5"
    Test-API -Method "GET" -Endpoint $pageUrl -Description "分页查询我的贷款申请" -AuthHeader $USER_TOKEN
}

Write-Host "`n=== 4. 文件服务接口测试 ===" -ForegroundColor Blue

# 13. 文件上传测试（需要认证）
if ($USER_TOKEN -ne "") {
    Test-FileUpload -Endpoint "$API_BASE/files/upload" -Description "文件上传" -AuthHeader $USER_TOKEN
}

Write-Host "`n=== 5. OA后台管理接口测试 ===" -ForegroundColor Blue

# 14. OA用户登录
Write-Host "`n正在进行OA管理员登录获取Token..." -ForegroundColor Yellow
$adminLoginData = @{
    username = "admin"
    password = "admin123"
}
$adminLoginResponse = Test-API -Method "POST" -Endpoint "$API_BASE/admin/login" -Description "OA用户登录" -Body $adminLoginData

# 提取管理员Token
if ($adminLoginResponse -and $adminLoginResponse.data.token) {
    $ADMIN_TOKEN = "Bearer " + $adminLoginResponse.data.token
    Write-Host "✓ 成功获取管理员Token: $($adminLoginResponse.data.token.Substring(0, 20))..." -ForegroundColor Green
} else {
    Write-Host "✗ 未能获取管理员Token" -ForegroundColor Red
}

# 15. 获取待审批贷款申请列表（需要管理员认证）
if ($ADMIN_TOKEN -ne "") {
    Test-API -Method "GET" -Endpoint "$API_BASE/admin/loans/applications/pending" -Description "获取待审批贷款申请列表" -AuthHeader $ADMIN_TOKEN
}

# 16. 获取贷款申请详情（管理员视角）
if ($ADMIN_TOKEN -ne "") {
    Test-API -Method "GET" -Endpoint "$API_BASE/admin/loans/applications/test_app_id" -Description "获取贷款申请详情(管理员)" -AuthHeader $ADMIN_TOKEN
}

# 17. 提交审批决策（需要管理员认证）
if ($ADMIN_TOKEN -ne "") {
    $reviewData = @{
        decision = "approved"
        approved_amount = 25000
        comments = "申请人信用良好，略微调整批准金额"
        required_info_details = $null
    }
    Test-API -Method "POST" -Endpoint "$API_BASE/admin/loans/applications/test_app_id/review" -Description "提交审批决策" -Body $reviewData -AuthHeader $ADMIN_TOKEN
}

# 18. 控制AI审批流程开关（需要管理员认证）
if ($ADMIN_TOKEN -ne "") {
    $toggleData = @{
        enabled = $true
    }
    Test-API -Method "POST" -Endpoint "$API_BASE/admin/system/ai-approval/toggle" -Description "控制AI审批流程开关" -Body $toggleData -AuthHeader $ADMIN_TOKEN
}

Write-Host "`n=== 6. 错误处理测试 ===" -ForegroundColor Blue

# 19. 测试未授权访问
Test-API -Method "GET" -Endpoint "$API_BASE/users/me" -Description "未授权访问用户信息" -ExpectedStatus 401

# 20. 测试无效的产品ID
Test-API -Method "GET" -Endpoint "$API_BASE/loans/products/invalid_id" -Description "查询不存在的产品" -ExpectedStatus 404

# 21. 测试无效的请求数据
$invalidData = @{
    phone = "invalid"
}
Test-API -Method "POST" -Endpoint "$API_BASE/users/register" -Description "无效注册请求" -Body $invalidData -ExpectedStatus 400

# 22. 测试无效的Token
Test-API -Method "GET" -Endpoint "$API_BASE/users/me" -Description "无效Token访问" -AuthHeader "Bearer invalid_token" -ExpectedStatus 401

Write-Host "`n=== 7. 性能和边界测试 ===" -ForegroundColor Blue

# 23. 测试大数据量查询
if ($USER_TOKEN -ne "") {
    $bigPageUrl = "$API_BASE/loans/applications/my?page=1`&limit=100"
    Test-API -Method "GET" -Endpoint $bigPageUrl -Description "大分页查询" -AuthHeader $USER_TOKEN
}

# 24. 测试空查询参数
Test-API -Method "GET" -Endpoint "$API_BASE/loans/products?category=" -Description "空分类查询"

Write-Host "`n===================================================" -ForegroundColor Yellow
Write-Host "                测试结果统计" -ForegroundColor Yellow
Write-Host "===================================================" -ForegroundColor Yellow
Write-Host "总测试数: $TOTAL_TESTS"
Write-Host "通过: $PASSED_TESTS" -ForegroundColor Green
Write-Host "失败: $FAILED_TESTS" -ForegroundColor Red

if ($FAILED_TESTS -eq 0) {
    Write-Host "`n🎉 所有测试通过！" -ForegroundColor Green
} else {
    Write-Host "`n❌ 有测试失败，请检查API实现" -ForegroundColor Red
    
    Write-Host "`n注意事项：" -ForegroundColor Yellow
    Write-Host "- 某些失败可能是因为测试环境中没有相应的数据"
    Write-Host "- 404错误在查询不存在资源时是正常的"
    Write-Host "- 401错误在未授权访问时是正常的" 
    Write-Host "- 确保数字惠农后端服务正在运行在 http://localhost:8080"
    Write-Host "- 检查数据库连接和初始化数据是否正确"
    
    # 计算成功率
    $successRate = [math]::Round(($PASSED_TESTS * 100) / $TOTAL_TESTS, 2)
    Write-Host "`n成功率: $successRate%"
    
    if ($successRate -ge 80) {
        Write-Host "✓ 总体测试通过率良好" -ForegroundColor Yellow
    } else {
        Write-Host "✗ 测试通过率较低，需要重点检查" -ForegroundColor Red
    }
}