#!/bin/bash

# AI智能体接口测试脚本
# 测试AI智能体相关的所有API接口

BASE_URL="http://localhost:8080"
AI_AGENT_TOKEN="ai_agent_secure_token_123456"
SYSTEM_TOKEN="system_secure_token"

echo "==================================================="
echo "   AI智能体接口测试开始"
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

# 测试用的数据
TEST_APPLICATION_ID="app_20240301_001"
TEST_USER_ID="user_20240301_001"

# 测试函数
test_api() {
    local method=$1
    local endpoint=$2
    local description=$3
    local expected_status=${4:-200}
    local auth_header=$5
    local request_body=$6
    
    TOTAL_TESTS=$((TOTAL_TESTS + 1))
    
    echo -e "\n${YELLOW}[TEST $TOTAL_TESTS]${NC} $description"
    echo "请求: $method $endpoint"
    
    # 构建curl命令
    local curl_cmd="curl -s -w \"\n%{http_code}\" \"$BASE_URL$endpoint\""
    
    if [ ! -z "$auth_header" ]; then
        curl_cmd="$curl_cmd -H \"Authorization: $auth_header\""
    fi
    
    if [ "$method" != "GET" ]; then
        curl_cmd="$curl_cmd -X $method"
    fi
    
    if [ ! -z "$request_body" ]; then
        curl_cmd="$curl_cmd -H \"Content-Type: application/json\" -d '$request_body'"
    fi
    
    # 执行请求
    response=$(eval $curl_cmd)
    
    # 提取HTTP状态码
    http_code=$(echo "$response" | tail -n1)
    response_body=$(echo "$response" | head -n -1)
    
    if [ "$http_code" = "$expected_status" ]; then
        echo -e "${GREEN}✓ 通过${NC} (状态码: $http_code)"
        PASSED_TESTS=$((PASSED_TESTS + 1))
        
        # 美化显示JSON响应
        if [ ! -z "$response_body" ]; then
            echo -e "${BLUE}响应数据:${NC}"
            echo "$response_body" | python3 -m json.tool 2>/dev/null || echo "$response_body"
        fi
    else
        echo -e "${RED}✗ 失败${NC} (期望状态码: $expected_status, 实际状态码: $http_code)"
        echo "响应: $response_body"
        FAILED_TESTS=$((FAILED_TESTS + 1))
    fi
    
    echo "---------------------------------------------------"
}

# 创建测试数据的函数
create_test_data() {
    echo -e "\n${BLUE}准备测试数据...${NC}"
    
    # 这里应该调用真实的API创建测试数据
    # 为了测试，我们假设已经有了测试数据
    echo "使用预设的测试数据: $TEST_APPLICATION_ID, $TEST_USER_ID"
}

echo -e "\n${YELLOW}开始测试AI智能体接口...${NC}"

create_test_data

# 1. 测试获取申请信息接口
test_api "GET" \
    "/api/v1/ai-agent/applications/$TEST_APPLICATION_ID/info" \
    "获取申请信息 - 正常请求" \
    200 \
    "AI-Agent-Token $AI_AGENT_TOKEN"

# 2. 测试获取申请信息接口 - 无认证
test_api "GET" \
    "/api/v1/ai-agent/applications/$TEST_APPLICATION_ID/info" \
    "获取申请信息 - 无认证头" \
    401

# 3. 测试获取申请信息接口 - 错误Token
test_api "GET" \
    "/api/v1/ai-agent/applications/$TEST_APPLICATION_ID/info" \
    "获取申请信息 - 无效Token" \
    401 \
    "AI-Agent-Token invalid_token"

# 4. 测试获取申请信息接口 - 不存在的申请
test_api "GET" \
    "/api/v1/ai-agent/applications/nonexistent_app/info" \
    "获取申请信息 - 不存在的申请" \
    500 \
    "AI-Agent-Token $AI_AGENT_TOKEN"

# 5. 测试提交AI决策结果接口
ai_decision_body='{
  "ai_analysis": {
    "risk_level": "LOW",
    "risk_score": 0.25,
    "confidence_score": 0.92,
    "analysis_summary": "申请人信用状况良好，还款能力强",
    "detailed_analysis": {
      "income_analysis": "年收入8万元，稳定",
      "credit_analysis": "征信良好，无不良记录",
      "asset_analysis": "拥有10亩土地，资产充足"
    },
    "risk_factors": [
      {
        "factor": "credit_score",
        "value": 750,
        "weight": 0.3,
        "risk_contribution": 0.05
      }
    ],
    "recommendations": [
      "建议批准贷款",
      "可给予优惠利率"
    ]
  },
  "ai_decision": {
    "decision": "AUTO_APPROVED",
    "approved_amount": 30000,
    "approved_term_months": 12,
    "suggested_interest_rate": "6.5%",
    "conditions": [],
    "next_action": "AWAIT_FINAL_CONFIRMATION"
  },
  "processing_info": {
    "ai_model_version": "v2.1.0",
    "processing_time_ms": 1500,
    "workflow_id": "workflow_001",
    "processed_at": "2024-03-01T10:35:00Z"
  }
}'

test_api "POST" \
    "/api/v1/ai-agent/applications/$TEST_APPLICATION_ID/ai-decision" \
    "提交AI决策结果 - 自动通过" \
    200 \
    "AI-Agent-Token $AI_AGENT_TOKEN" \
    "$ai_decision_body"

# 6. 测试提交AI决策结果 - 需要人工审核
ai_decision_manual_body='{
  "ai_analysis": {
    "risk_level": "MEDIUM",
    "risk_score": 0.65,
    "confidence_score": 0.78,
    "analysis_summary": "申请金额较大，需要人工审核",
    "detailed_analysis": {},
    "risk_factors": [],
    "recommendations": ["建议人工审核"]
  },
  "ai_decision": {
    "decision": "REQUIRE_HUMAN_REVIEW",
    "approved_amount": 0,
    "approved_term_months": 0,
    "suggested_interest_rate": "",
    "conditions": ["需要提供更多收入证明"],
    "next_action": "AWAIT_HUMAN_REVIEW"
  },
  "processing_info": {
    "ai_model_version": "v2.1.0",
    "processing_time_ms": 2000,
    "workflow_id": "workflow_002",
    "processed_at": "2024-03-01T10:40:00Z"
  }
}'

test_api "POST" \
    "/api/v1/ai-agent/applications/app_20240301_002/ai-decision" \
    "提交AI决策结果 - 需要人工审核" \
    200 \
    "AI-Agent-Token $AI_AGENT_TOKEN" \
    "$ai_decision_manual_body"

# 7. 测试触发AI审批工作流
trigger_workflow_body='{
  "workflow_type": "LOAN_APPROVAL",
  "priority": "NORMAL",
  "callback_url": "https://example.com/api/callback"
}'

test_api "POST" \
    "/api/v1/ai-agent/applications/$TEST_APPLICATION_ID/trigger-workflow" \
    "触发AI审批工作流" \
    200 \
    "System-Token $SYSTEM_TOKEN" \
    "$trigger_workflow_body"

# 8. 测试触发工作流 - 高优先级
trigger_workflow_high_body='{
  "workflow_type": "LOAN_APPROVAL",
  "priority": "HIGH"
}'

test_api "POST" \
    "/api/v1/ai-agent/applications/$TEST_APPLICATION_ID/trigger-workflow" \
    "触发AI审批工作流 - 高优先级" \
    200 \
    "System-Token $SYSTEM_TOKEN" \
    "$trigger_workflow_high_body"

# 9. 测试获取AI模型配置
test_api "GET" \
    "/api/v1/ai-agent/config/models" \
    "获取AI模型配置" \
    200 \
    "AI-Agent-Token $AI_AGENT_TOKEN"

# 10. 测试获取外部数据
test_api "GET" \
    "/api/v1/ai-agent/external-data/$TEST_USER_ID" \
    "获取外部数据 - 默认类型" \
    200 \
    "AI-Agent-Token $AI_AGENT_TOKEN"

# 11. 测试获取外部数据 - 指定类型
test_api "GET" \
    "/api/v1/ai-agent/external-data/$TEST_USER_ID?data_types=credit,bank_flow" \
    "获取外部数据 - 指定类型" \
    200 \
    "AI-Agent-Token $AI_AGENT_TOKEN"

# 12. 测试更新申请状态
update_status_body='{
  "status": "MANUAL_REVIEW_REQUIRED",
  "operator": "ai_system",
  "remarks": "AI分析建议人工审核",
  "metadata": {
    "reason": "高风险申请",
    "review_priority": "HIGH"
  }
}'

test_api "PUT" \
    "/api/v1/ai-agent/applications/$TEST_APPLICATION_ID/status" \
    "更新申请状态" \
    200 \
    "AI-Agent-Token $AI_AGENT_TOKEN" \
    "$update_status_body"

# 13. 测试错误情况 - 缺少必需参数
invalid_decision_body='{
  "ai_analysis": {
    "risk_level": "LOW"
  }
}'

test_api "POST" \
    "/api/v1/ai-agent/applications/$TEST_APPLICATION_ID/ai-decision" \
    "提交AI决策结果 - 缺少必需参数" \
    400 \
    "AI-Agent-Token $AI_AGENT_TOKEN" \
    "$invalid_decision_body"

# 14. 测试错误情况 - 空的申请ID
test_api "GET" \
    "/api/v1/ai-agent/applications//info" \
    "获取申请信息 - 空申请ID" \
    404 \
    "AI-Agent-Token $AI_AGENT_TOKEN"

# 15. 测试工作流触发 - 缺少必需参数
invalid_workflow_body='{
  "priority": "NORMAL"
}'

test_api "POST" \
    "/api/v1/ai-agent/applications/$TEST_APPLICATION_ID/trigger-workflow" \
    "触发工作流 - 缺少必需参数" \
    400 \
    "System-Token $SYSTEM_TOKEN" \
    "$invalid_workflow_body"

echo -e "\n==================================================="
echo "           测试结果统计"
echo "==================================================="
echo "总测试数: $TOTAL_TESTS"
echo -e "${GREEN}通过: $PASSED_TESTS${NC}"
echo -e "${RED}失败: $FAILED_TESTS${NC}"

# 计算成功率
if [ $TOTAL_TESTS -gt 0 ]; then
    success_rate=$(echo "scale=2; $PASSED_TESTS * 100 / $TOTAL_TESTS" | bc -l)
    echo "成功率: ${success_rate}%"
fi

echo -e "\n==================================================="
echo "           详细测试报告"
echo "==================================================="

if [ $FAILED_TESTS -eq 0 ]; then
    echo -e "\n${GREEN}🎉 所有测试通过！${NC}"
    echo -e "${GREEN}AI智能体接口测试全部成功${NC}"
    exit 0
else
    echo -e "\n${RED}❌ 有测试失败，请检查以下问题：${NC}"
    echo -e "\n${YELLOW}可能的原因：${NC}"
    echo "1. 服务器未启动或端口错误"
    echo "2. 数据库中缺少测试数据"
    echo "3. Token配置不正确"
    echo "4. 接口实现存在问题"
    
    echo -e "\n${YELLOW}调试建议：${NC}"
    echo "1. 检查服务器日志: docker logs [container_name]"
    echo "2. 确认数据库连接: 检查数据库连接配置"
    echo "3. 验证Token配置: 检查middleware.go中的Token列表"
    echo "4. 手动创建测试数据"
    
    echo -e "\n${YELLOW}手动测试命令：${NC}"
    echo "# 测试健康检查"
    echo "curl -X GET \"$BASE_URL/health\""
    echo ""
    echo "# 测试获取申请信息"
    echo "curl -X GET \"$BASE_URL/api/v1/ai-agent/applications/$TEST_APPLICATION_ID/info\" \\"
    echo "  -H \"Authorization: AI-Agent-Token $AI_AGENT_TOKEN\""
    
    exit 1
fi

echo -e "\n==================================================="
echo "           性能统计"
echo "==================================================="
echo "测试执行时间: $(date)"
echo "接口覆盖率: 100% (6个核心接口)"
echo "认证方式测试: ✓ AI-Agent-Token, ✓ System-Token"
echo "错误处理测试: ✓ 401, ✓ 400, ✓ 404"
echo "业务场景测试: ✓ 自动通过, ✓ 人工审核, ✓ 高优先级" 