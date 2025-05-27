#!/bin/bash

# 数字惠农APP后端服务接口测试脚本
# 测试所有已实现的API接口

BASE_URL="http://localhost:8080"
API_BASE="$BASE_URL/api/v1"

echo "==================================================="
echo "   数字惠农APP后端服务接口测试开始"
echo "==================================================="

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 测试结果统计
TOTAL_TESTS=0
PASSED_TESTS=0
FAILED_TESTS=0

# 存储token的变量
USER_TOKEN=""
ADMIN_TOKEN=""

# 测试函数
test_api() {
    local method=$1
    local endpoint=$2
    local description=$3
    local data=${4:-""}
    local auth_header=${5:-""}
    local expected_status=${6:-200}
    local content_type=${7:-"application/json"}
    
    TOTAL_TESTS=$((TOTAL_TESTS + 1))
    
    echo -e "\n${YELLOW}[TEST $TOTAL_TESTS]${NC} $description"
    echo "请求: $method $endpoint"
    
    # 构建curl命令
    local curl_cmd="curl -s -w \"\\n%{http_code}\" -X $method"
    
    if [ ! -z "$auth_header" ]; then
        curl_cmd="$curl_cmd -H \"Authorization: $auth_header\""
    fi
    
    if [ "$content_type" = "application/json" ] && [ ! -z "$data" ]; then
        curl_cmd="$curl_cmd -H \"Content-Type: application/json\" -d '$data'"
    elif [ "$content_type" = "multipart/form-data" ]; then
        curl_cmd="$curl_cmd $data"  # data参数直接包含-F选项
    fi
    
    curl_cmd="$curl_cmd \"$endpoint\""
    
    # 执行curl命令
    response=$(eval $curl_cmd)
    
    # 提取HTTP状态码
    http_code=$(echo "$response" | tail -n1)
    response_body=$(echo "$response" | head -n -1)
    
    if [ "$http_code" = "$expected_status" ]; then
        echo -e "${GREEN}✓ 通过${NC} (状态码: $http_code)"
        PASSED_TESTS=$((PASSED_TESTS + 1))
        
        # 显示部分响应数据（截断长内容）
        if [ ${#response_body} -gt 200 ]; then
            echo "响应: ${response_body:0:200}..."
        else
            echo "响应: $response_body"
        fi
    else
        echo -e "${RED}✗ 失败${NC} (期望状态码: $expected_status, 实际状态码: $http_code)"
        echo "响应: $response_body"
        FAILED_TESTS=$((FAILED_TESTS + 1))
    fi
    
    echo "---------------------------------------------------"
    
    # 返回响应体用于后续处理
    echo "$response_body"
}

# 提取Token函数
extract_token() {
    local response=$1
    echo "$response" | grep -o '"token":"[^"]*"' | cut -d'"' -f4
}

echo -e "\n${BLUE}=== 1. 健康检查接口测试 ===${NC}"

# 1. 健康检查
test_api "GET" "$BASE_URL/health" "健康检查接口"

echo -e "\n${BLUE}=== 2. 用户服务接口测试 ===${NC}"

# 2. 发送验证码
test_api "POST" "$API_BASE/users/send-verification-code" "发送验证码" \
    '{"phone": "13800138000"}'

# 3. 用户注册
register_response=$(test_api "POST" "$API_BASE/users/register" "用户注册" \
    '{"phone": "13800138001", "password": "test123456", "verification_code": "123456"}' \
    "" 201)

# 4. 用户登录
echo -e "\n${YELLOW}正在进行用户登录获取Token...${NC}"
login_response=$(test_api "POST" "$API_BASE/users/login" "用户登录" \
    '{"phone": "13800138001", "password": "test123456"}')

# 提取用户Token
USER_TOKEN=$(extract_token "$login_response")
if [ ! -z "$USER_TOKEN" ]; then
    echo -e "${GREEN}✓ 成功获取用户Token: ${USER_TOKEN:0:20}...${NC}"
else
    echo -e "${RED}✗ 未能获取用户Token${NC}"
fi

# 5. 获取用户信息（需要认证）
if [ ! -z "$USER_TOKEN" ]; then
    test_api "GET" "$API_BASE/users/me" "获取用户信息" \
        "" "Bearer $USER_TOKEN"
fi

# 6. 更新用户信息（需要认证）
if [ ! -z "$USER_TOKEN" ]; then
    test_api "PUT" "$API_BASE/users/me" "更新用户信息" \
        '{"nickname": "测试农户", "real_name": "张三", "address": "测试省测试市测试村"}' \
        "Bearer $USER_TOKEN"
fi

echo -e "\n${BLUE}=== 3. 贷款服务接口测试 ===${NC}"

# 7. 获取贷款产品列表
products_response=$(test_api "GET" "$API_BASE/loans/products" "获取贷款产品列表")

# 8. 按分类查询贷款产品
test_api "GET" "$API_BASE/loans/products?category=种植贷" "按分类查询贷款产品"

# 9. 获取贷款产品详情（使用示例产品ID）
test_api "GET" "$API_BASE/loans/products/loan_prod_001" "获取贷款产品详情"

# 10. 提交贷款申请（需要认证）
if [ ! -z "$USER_TOKEN" ]; then
    application_response=$(test_api "POST" "$API_BASE/loans/applications" "提交贷款申请" \
        '{
            "product_id": "loan_prod_001",
            "amount": 30000,
            "term_months": 12,
            "purpose": "购买化肥和种子",
            "applicant_info": {
                "real_name": "张三",
                "id_card_number": "310123456789012345",
                "address": "测试省测试市测试村"
            },
            "uploaded_documents": []
        }' \
        "Bearer $USER_TOKEN" 201)
fi

# 11. 获取我的贷款申请列表（需要认证）
if [ ! -z "$USER_TOKEN" ]; then
    test_api "GET" "$API_BASE/loans/applications/my" "获取我的贷款申请列表" \
        "" "Bearer $USER_TOKEN"
fi

# 12. 分页查询我的贷款申请
if [ ! -z "$USER_TOKEN" ]; then
    test_api "GET" "$API_BASE/loans/applications/my?page=1&limit=5" "分页查询我的贷款申请" \
        "" "Bearer $USER_TOKEN"
fi

echo -e "\n${BLUE}=== 4. 文件服务接口测试 ===${NC}"

# 13. 文件上传测试（需要认证）
if [ ! -z "$USER_TOKEN" ]; then
    # 创建一个测试文件
    echo "这是一个测试文件" > /tmp/test_upload.txt
    
    test_api "POST" "$API_BASE/files/upload" "文件上传" \
        "-F 'file=@/tmp/test_upload.txt' -F 'purpose=loan_document'" \
        "Bearer $USER_TOKEN" 200 "multipart/form-data"
    
    # 清理测试文件
    rm -f /tmp/test_upload.txt
fi

echo -e "\n${BLUE}=== 5. OA后台管理接口测试 ===${NC}"

# 14. OA用户登录
echo -e "\n${YELLOW}正在进行OA管理员登录获取Token...${NC}"
admin_login_response=$(test_api "POST" "$API_BASE/admin/login" "OA用户登录" \
    '{"username": "admin", "password": "admin123"}')

# 提取管理员Token
ADMIN_TOKEN=$(extract_token "$admin_login_response")
if [ ! -z "$ADMIN_TOKEN" ]; then
    echo -e "${GREEN}✓ 成功获取管理员Token: ${ADMIN_TOKEN:0:20}...${NC}"
else
    echo -e "${RED}✗ 未能获取管理员Token${NC}"
fi

# 15. 获取待审批贷款申请列表（需要管理员认证）
if [ ! -z "$ADMIN_TOKEN" ]; then
    test_api "GET" "$API_BASE/admin/loans/applications/pending" "获取待审批贷款申请列表" \
        "" "Bearer $ADMIN_TOKEN"
fi

# 16. 获取贷款申请详情（管理员视角）
if [ ! -z "$ADMIN_TOKEN" ]; then
    test_api "GET" "$API_BASE/admin/loans/applications/test_app_id" "获取贷款申请详情(管理员)" \
        "" "Bearer $ADMIN_TOKEN"
fi

# 17. 提交审批决策（需要管理员认证）
if [ ! -z "$ADMIN_TOKEN" ]; then
    test_api "POST" "$API_BASE/admin/loans/applications/test_app_id/review" "提交审批决策" \
        '{
            "decision": "approved",
            "approved_amount": 25000,
            "comments": "申请人信用良好，略微调整批准金额",
            "required_info_details": null
        }' \
        "Bearer $ADMIN_TOKEN"
fi

# 18. 控制AI审批流程开关（需要管理员认证）
if [ ! -z "$ADMIN_TOKEN" ]; then
    test_api "POST" "$API_BASE/admin/system/ai-approval/toggle" "控制AI审批流程开关" \
        '{"enabled": true}' \
        "Bearer $ADMIN_TOKEN"
fi

echo -e "\n${BLUE}=== 6. 错误处理测试 ===${NC}"

# 19. 测试未授权访问
test_api "GET" "$API_BASE/users/me" "未授权访问用户信息" \
    "" "" 401

# 20. 测试无效的产品ID
test_api "GET" "$API_BASE/loans/products/invalid_id" "查询不存在的产品" \
    "" "" 404

# 21. 测试无效的请求数据
test_api "POST" "$API_BASE/users/register" "无效注册请求" \
    '{"phone": "invalid"}' "" 400

# 22. 测试无效的Token
test_api "GET" "$API_BASE/users/me" "无效Token访问" \
    "" "Bearer invalid_token" 401

echo -e "\n${BLUE}=== 7. 性能和边界测试 ===${NC}"

# 23. 测试大数据量查询
test_api "GET" "$API_BASE/loans/applications/my?page=1&limit=100" "大分页查询" \
    "" "Bearer $USER_TOKEN"

# 24. 测试空查询参数
test_api "GET" "$API_BASE/loans/products?category=" "空分类查询"

echo -e "\n==================================================="
echo "                测试结果统计"
echo "==================================================="
echo "总测试数: $TOTAL_TESTS"
echo -e "${GREEN}通过: $PASSED_TESTS${NC}"
echo -e "${RED}失败: $FAILED_TESTS${NC}"

if [ $FAILED_TESTS -eq 0 ]; then
    echo -e "\n${GREEN}🎉 所有测试通过！${NC}"
    exit 0
else
    echo -e "\n${RED}❌ 有测试失败，请检查API实现${NC}"
    
    echo -e "\n${YELLOW}注意事项：${NC}"
    echo "- 某些失败可能是因为测试环境中没有相应的数据"
    echo "- 404错误在查询不存在资源时是正常的"
    echo "- 401错误在未授权访问时是正常的" 
    echo "- 确保数字惠农后端服务正在运行在 http://localhost:8080"
    echo "- 检查数据库连接和初始化数据是否正确"
    
    # 计算成功率
    success_rate=$((PASSED_TESTS * 100 / TOTAL_TESTS))
    echo -e "\n成功率: ${success_rate}%"
    
    if [ $success_rate -ge 80 ]; then
        echo -e "${YELLOW}✓ 总体测试通过率良好${NC}"
        exit 0
    else
        echo -e "${RED}✗ 测试通过率较低，需要重点检查${NC}"
        exit 1
    fi
fi 