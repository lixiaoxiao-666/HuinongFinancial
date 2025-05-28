#!/bin/bash

# 数字惠农OA管理系统后端接口测试脚本
# 测试OA系统的审批管理、用户管理、系统配置等API接口

BASE_URL="http://localhost:8080"
ADMIN_TOKEN=""
REVIEWER_TOKEN=""

echo "==================================================="
echo "   数字惠农OA管理系统 后端接口测试开始"
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

# 测试函数
test_api() {
    local method=$1
    local endpoint=$2
    local description=$3
    local data="$4"
    local token="$5"
    local expected_status=${6:-200}
    
    TOTAL_TESTS=$((TOTAL_TESTS + 1))
    
    echo -e "\n${YELLOW}[TEST $TOTAL_TESTS]${NC} $description"
    echo "请求: $method $endpoint"
    
    # 构建curl命令
    local curl_cmd="curl -s -w \"\n%{http_code}\" \"$BASE_URL$endpoint\""
    
    if [ "$method" != "GET" ]; then
        curl_cmd="$curl_cmd -X $method"
    fi
    
    if [ -n "$token" ]; then
        curl_cmd="$curl_cmd -H \"Authorization: Bearer $token\""
    fi
    
    if [ -n "$data" ]; then
        curl_cmd="$curl_cmd -H \"Content-Type: application/json\" -d '$data'"
    fi
    
    # 执行请求
    response=$(eval $curl_cmd)
    
    # 提取HTTP状态码
    http_code=$(echo "$response" | tail -n1)
    response_body=$(echo "$response" | head -n -1)
    
    if [ "$http_code" = "$expected_status" ]; then
        echo -e "${GREEN}✓ 通过${NC} (状态码: $http_code)"
        PASSED_TESTS=$((PASSED_TESTS + 1))
        
        # 显示关键响应数据
        if [ ${#response_body} -lt 500 ]; then
            echo "响应: $response_body"
        else
            echo "响应: ${response_body:0:200}..."
        fi
    else
        echo -e "${RED}✗ 失败${NC} (期望状态码: $expected_status, 实际状态码: $http_code)"
        echo "响应: $response_body"
        FAILED_TESTS=$((FAILED_TESTS + 1))
    fi
    
    echo "---------------------------------------------------"
}

# 辅助函数：从响应中提取token
extract_token() {
    local response="$1"
    echo "$response" | grep -o '"token":"[^"]*"' | cut -d'"' -f4
}

echo -e "\n${BLUE}======================== 1. OA用户认证测试 ========================${NC}"

# 1.1 管理员登录
test_api "POST" "/admin/login" "管理员登录" '{
  "username": "admin",
  "password": "admin123"
}' "" 200

# 提取管理员token（简化处理，实际环境中需要解析JSON）
ADMIN_TOKEN=$(curl -s -X POST "$BASE_URL/admin/login" \
  -H "Content-Type: application/json" \
  -d '{"username": "admin", "password": "admin123"}' | \
  grep -o '"token":"[^"]*"' | cut -d'"' -f4)

echo "管理员Token: ${ADMIN_TOKEN:0:20}..."

# 1.2 审批员登录
test_api "POST" "/admin/login" "审批员登录" '{
  "username": "reviewer",
  "password": "reviewer123"
}' "" 200

# 提取审批员token
REVIEWER_TOKEN=$(curl -s -X POST "$BASE_URL/admin/login" \
  -H "Content-Type: application/json" \
  -d '{"username": "reviewer", "password": "reviewer123"}' | \
  grep -o '"token":"[^"]*"' | cut -d'"' -f4)

echo "审批员Token: ${REVIEWER_TOKEN:0:20}..."

echo -e "\n${BLUE}======================== 2. OA首页/工作台测试 ========================${NC}"

# 2.1 获取OA首页信息
test_api "GET" "/admin/dashboard" "获取OA首页/工作台信息" "" "$ADMIN_TOKEN" 200

echo -e "\n${BLUE}======================== 3. 审批管理测试 ========================${NC}"

# 3.1 获取待审批申请列表
test_api "GET" "/admin/loans/applications/pending" "获取待审批申请列表" "" "$REVIEWER_TOKEN" 200

# 3.2 获取待审批申请列表（带筛选）
test_api "GET" "/admin/loans/applications/pending?status_filter=MANUAL_REVIEW_REQUIRED&page=1&limit=5" "获取待审批申请列表（筛选）" "" "$REVIEWER_TOKEN" 200

# 3.3 获取申请详情（需要真实的application_id，这里使用示例ID）
test_api "GET" "/admin/loans/applications/la_test_app_001" "获取申请详情" "" "$REVIEWER_TOKEN" 200

# 3.4 提交审批决策 - 批准
test_api "POST" "/admin/loans/applications/la_test_app_001/review" "提交审批决策（批准）" '{
  "decision": "approved",
  "approved_amount": 25000.00,
  "approved_term_months": 12,
  "comments": "申请人资质良好，批准贷款申请。"
}' "$REVIEWER_TOKEN" 200

# 3.5 提交审批决策 - 拒绝
test_api "POST" "/admin/loans/applications/la_test_app_002/review" "提交审批决策（拒绝）" '{
  "decision": "rejected",
  "comments": "申请人收入证明不足，拒绝贷款申请。"
}' "$REVIEWER_TOKEN" 200

# 3.6 提交审批决策 - 要求补充信息
test_api "POST" "/admin/loans/applications/la_test_app_003/review" "提交审批决策（要求补充信息）" '{
  "decision": "request_more_info",
  "comments": "需要补充收入证明材料",
  "required_info_details": "请提供最近3个月的银行流水和收入证明"
}' "$REVIEWER_TOKEN" 200

echo -e "\n${BLUE}======================== 4. 系统管理测试 ========================${NC}"

# 4.1 获取系统统计信息
test_api "GET" "/admin/system/stats" "获取系统统计信息" "" "$ADMIN_TOKEN" 200

# 4.2 AI审批开关控制 - 启用
test_api "POST" "/admin/system/ai-approval/toggle" "启用AI审批" '{
  "enabled": true
}' "$ADMIN_TOKEN" 200

# 4.3 AI审批开关控制 - 禁用
test_api "POST" "/admin/system/ai-approval/toggle" "禁用AI审批" '{
  "enabled": false
}' "$ADMIN_TOKEN" 200

echo -e "\n${BLUE}======================== 5. 用户管理测试 ========================${NC}"

# 5.1 获取OA用户列表
test_api "GET" "/admin/users?page=1&limit=10" "获取OA用户列表" "" "$ADMIN_TOKEN" 200

# 5.2 获取OA用户列表（角色筛选）
test_api "GET" "/admin/users?role=REVIEWER&page=1&limit=10" "获取OA用户列表（审批员）" "" "$ADMIN_TOKEN" 200

# 5.3 创建OA用户
test_api "POST" "/admin/users" "创建OA用户" '{
  "username": "test_reviewer",
  "password": "password123",
  "role": "审批员",
  "display_name": "测试审批员",
  "email": "test@example.com"
}' "$ADMIN_TOKEN" 200

# 5.4 更新OA用户状态 - 禁用（需要真实的user_id）
test_api "PUT" "/admin/users/oa_test_user_001/status" "更新OA用户状态（禁用）" '{
  "status": 1
}' "$ADMIN_TOKEN" 200

# 5.5 更新OA用户状态 - 启用
test_api "PUT" "/admin/users/oa_test_user_001/status" "更新OA用户状态（启用）" '{
  "status": 0
}' "$ADMIN_TOKEN" 200

echo -e "\n${BLUE}======================== 6. 操作日志测试 ========================${NC}"

# 6.1 获取操作日志
test_api "GET" "/admin/logs?page=1&limit=10" "获取操作日志" "" "$ADMIN_TOKEN" 200

# 6.2 获取操作日志（带筛选）
test_api "GET" "/admin/logs?action=审批申请&start_date=2024-03-01&end_date=2024-03-31&page=1&limit=5" "获取操作日志（筛选）" "" "$ADMIN_TOKEN" 200

echo -e "\n${BLUE}======================== 7. 系统配置测试 ========================${NC}"

# 7.1 获取系统配置
test_api "GET" "/admin/configs" "获取系统配置" "" "$ADMIN_TOKEN" 200

# 7.2 更新系统配置
test_api "PUT" "/admin/configs/ai_approval_enabled" "更新系统配置" '{
  "config_value": "true"
}' "$ADMIN_TOKEN" 200

# 7.3 更新自定义配置
test_api "PUT" "/admin/configs/max_loan_amount" "更新自定义配置" '{
  "config_value": "500000"
}' "$ADMIN_TOKEN" 200

echo -e "\n${BLUE}======================== 8. 错误处理测试 ========================${NC}"

# 8.1 无效token测试
test_api "GET" "/admin/dashboard" "无效token访问" "" "invalid_token" 401

# 8.2 无权限访问测试（审批员访问管理员功能）
test_api "POST" "/admin/users" "无权限访问（创建用户）" '{
  "username": "unauthorized_test",
  "password": "password123",
  "role": "审批员",
  "display_name": "无权限测试",
  "email": "unauthorized@example.com"
}' "$REVIEWER_TOKEN" 403

# 8.3 参数错误测试
test_api "POST" "/admin/login" "登录参数错误" '{
  "username": "",
  "password": "short"
}' "" 400

# 8.4 资源不存在测试
test_api "GET" "/admin/loans/applications/non_existent_id" "不存在的申请ID" "" "$REVIEWER_TOKEN" 404

echo -e "\n${BLUE}======================== 9. 性能和压力测试示例 ========================${NC}"

# 9.1 并发请求测试（简化版）
echo -e "\n${YELLOW}执行并发请求测试...${NC}"
for i in {1..5}; do
    test_api "GET" "/admin/system/stats" "并发请求 #$i" "" "$ADMIN_TOKEN" 200 &
done
wait

echo -e "\n==================================================="
echo "           测试结果统计"
echo "==================================================="
echo "总测试数: $TOTAL_TESTS"
echo -e "${GREEN}通过: $PASSED_TESTS${NC}"
echo -e "${RED}失败: $FAILED_TESTS${NC}"

success_rate=$(echo "scale=2; $PASSED_TESTS * 100 / $TOTAL_TESTS" | bc)
echo "成功率: ${success_rate}%"

if [ $FAILED_TESTS -eq 0 ]; then
    echo -e "\n${GREEN}🎉 所有测试通过！OA管理系统接口运行正常！${NC}"
    exit 0
else
    echo -e "\n${RED}❌ 有测试失败，请检查API实现${NC}"
    echo -e "\n${YELLOW}注意事项：${NC}"
    echo "- 部分失败可能是因为测试环境中没有相应的数据"
    echo "- 404错误通常表示使用了示例ID，在实际测试中需要使用真实ID"
    echo "- 403错误表示权限控制正常工作"
    echo "- 检查服务是否正常启动并连接到数据库"
    echo "- 确保已创建默认的OA用户账号"
    
    echo -e "\n${BLUE}默认测试账号：${NC}"
    echo "管理员：admin / admin123"
    echo "审批员：reviewer / reviewer123"
    
    exit 1
fi

echo -e "\n${BLUE}==================================================="
echo "           OA接口功能说明"
echo "===================================================${NC}"
echo "1. 🔐 认证系统：支持OA用户登录和JWT token验证"
echo "2. 📊 工作台：提供系统统计、待办事项、快捷操作"
echo "3. 📋 审批管理：待审批列表、申请详情、审批决策"
echo "4. ⚙️  系统管理：AI审批开关、系统统计、配置管理"
echo "5. 👥 用户管理：OA用户创建、状态管理、权限控制"
echo "6. 📝 操作日志：记录和查询所有操作历史"
echo "7. 🔧 系统配置：灵活的配置项管理"
echo "8. 🛡️  安全控制：权限验证、参数校验、错误处理" 