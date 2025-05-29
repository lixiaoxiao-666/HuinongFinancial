/*
 Navicat Premium Dump SQL

 Source Server         : TiDB
 Source Server Type    : MySQL
 Source Server Version : 80011 (8.0.11-TiDB-v8.5.1)
 Source Host           : 10.10.10.10:4000
 Source Schema         : app

 Target Server Type    : MySQL
 Target Server Version : 80011 (8.0.11-TiDB-v8.5.1)
 File Encoding         : 65001

 Date: 29/05/2025 13:13:25
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for ai_agent_logs
-- ----------------------------
DROP TABLE IF EXISTS `ai_agent_logs`;
CREATE TABLE `ai_agent_logs`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `log_id` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL,
  `application_id` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL DEFAULT NULL,
  `action_type` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL,
  `agent_type` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL,
  `request_data` json NULL,
  `response_data` json NULL,
  `status` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL,
  `error_message` text CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL,
  `duration` bigint NULL DEFAULT NULL,
  `ip_address` varchar(45) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL DEFAULT NULL,
  `user_agent` text CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL,
  `occurred_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `idx_ai_agent_logs_log_id`(`log_id` ASC) USING BTREE,
  INDEX `idx_ai_agent_logs_application_id`(`application_id` ASC) USING BTREE,
  INDEX `idx_ai_agent_logs_action_type`(`action_type` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 30001 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_bin ROW_FORMAT = Compact;

-- ----------------------------
-- Records of ai_agent_logs
-- ----------------------------
INSERT INTO `ai_agent_logs` VALUES (1, '5ab23a18-e8a6-4b94-a2c0-1e90a73de389', 'test_app_001', 'GET_APPLICATION_INFO', 'DIFY_WORKFLOW', '{\"application_id\": \"test_app_001\"}', '{\"applicant_info\": {\"address\": \"北京市朝阳区测试街道123号\", \"age\": 35, \"id_card_number\": \"110***********1234\", \"is_verified\": true, \"phone\": \"138****8001\", \"real_name\": \"�*�\", \"user_id\": \"user_001\"}, \"application_id\": \"test_app_001\", \"application_info\": {\"amount\": 100000, \"purpose\": \"装修\", \"status\": \"pending_review\", \"submitted_at\": \"2025-05-29T13:10:14.768+08:00\", \"term_months\": 24}, \"external_data\": {\"blacklist_check\": false, \"credit_bureau_score\": 750, \"land_ownership_verified\": true, \"previous_loan_history\": []}, \"financial_info\": {\"account_balance\": 15000, \"annual_income\": 200000, \"credit_score\": 750, \"existing_loans\": 0, \"farming_experience\": \"10年\", \"land_area\": \"10亩\"}, \"product_info\": {\"category\": \"personal_credit\", \"interest_rate_yearly\": \"8.5%\", \"max_amount\": 500000, \"name\": \"个人信用贷\", \"product_id\": \"product_001\"}, \"uploaded_documents\": []}', 'SUCCESS', '', 36, '172.20.0.9', 'python-httpx/0.27.2', '2025-05-29 13:10:28.896', '2025-05-29 13:10:28.896');
INSERT INTO `ai_agent_logs` VALUES (2, '05b68ac8-5c39-4045-a16d-53cb9809534e', '', 'GET_EXTERNAL_DATA', 'DIFY_WORKFLOW', '{\"data_types\": \"credit_report,bank_flow,blacklist_check,government_subsidy\", \"user_id\": \"user_001\"}', '{\"bank_flow\": {\"account_stability\": \"稳定\", \"average_monthly_income\": 6500, \"last_6_months_flow\": [{\"expense\": 4800, \"income\": 7200, \"month\": \"2024-02\"}, {\"expense\": 5100, \"income\": 6800, \"month\": \"2024-01\"}]}, \"blacklist_check\": {\"check_time\": \"2025-05-29T13:10:29+08:00\", \"is_blacklisted\": false}, \"credit_report\": {\"grade\": \"优秀\", \"loan_history\": [], \"overdue_records\": 0, \"report_date\": \"2024-03-01\", \"score\": 750}, \"government_subsidy\": {\"received_subsidies\": [{\"amount\": 1200, \"type\": \"种粮补贴\", \"year\": 2023}, {\"amount\": 3000, \"type\": \"农机购置补贴\", \"year\": 2022}]}, \"user_id\": \"user_001\"}', 'SUCCESS', '', 33, '172.20.0.9', 'python-httpx/0.27.2', '2025-05-29 13:10:29.228', '2025-05-29 13:10:29.228');
INSERT INTO `ai_agent_logs` VALUES (3, 'e4c98d0b-be75-4ed2-a589-1d5d7e0d96e2', '', 'GET_AI_MODEL_CONFIG', 'DIFY_WORKFLOW', '{\"action\": \"get_model_config\"}', '{\"active_models\": [{\"model_id\": \"risk_assessment_v2\", \"model_type\": \"RISK_EVALUATION\", \"status\": \"ACTIVE\", \"thresholds\": {\"high_risk\": 0.9, \"low_risk\": 0.3, \"medium_risk\": 0.7}, \"version\": \"2.1.0\"}, {\"model_id\": \"fraud_detection_v1\", \"model_type\": \"FRAUD_DETECTION\", \"sensitivity\": 0.85, \"status\": \"ACTIVE\", \"version\": \"1.5.2\"}], \"approval_rules\": {\"auto_approval_threshold\": 0.3, \"auto_rejection_threshold\": 0.8, \"max_auto_approval_amount\": 50000, \"required_human_review_conditions\": [\"申请金额超过5万元\", \"信用评分低于700分\", \"存在潜在欺诈风险\"]}, \"business_parameters\": {\"max_debt_to_income_ratio\": 0.5, \"max_loan_amount_by_category\": {\"其他\": 30000, \"种植贷\": 50000, \"设备贷\": 200000}, \"min_credit_score\": 600}}', 'SUCCESS', '', 1, '172.20.0.9', 'python-httpx/0.27.2', '2025-05-29 13:10:29.414', '2025-05-29 13:10:29.414');
INSERT INTO `ai_agent_logs` VALUES (4, '959e6fd4-02b2-4500-bdb5-1a2a25387cb1', 'test_app_001', 'SUBMIT_AI_DECISION_QUERY', 'DIFY_WORKFLOW', '{\"ai_model_version\": \"v1.0.0\", \"analysis_summary\": \"AI风险分析\", \"application_id\": \"test_app_001\", \"approved_amount\": 0, \"approved_term_months\": 12, \"conditions\": [\"无特殊条件\"], \"confidence_score\": 0.5, \"decision\": \"REQUIRE_HUMAN_REVIEW\", \"detailed_analysis\": {\"analysis_summary\": \"AI风险分析\", \"approved_amount\": 0, \"approved_term_months\": 12, \"conditions\": [\"需要审核\"], \"confidence_score\": 0.5, \"decision\": \"REQUIRE_HUMAN_REVIEW\", \"detailed_analysis\": {\"credit_analysis\": \"信用分析\", \"financial_analysis\": \"财务分析\", \"risk_factors\": [\"待评估\"], \"strengths\": [\"待评估\"]}, \"recommendations\": [\"建议审核\"], \"risk_level\": \"MEDIUM\", \"risk_score\": 0.5, \"suggested_interest_rate\": \"5.0%\"}, \"recommendations\": [\"请关注后续通知\"], \"risk_level\": \"MEDIUM\", \"risk_score\": 0.5, \"suggested_interest_rate\": \"5.0%\", \"workflow_id\": \"dify_ai_workflow\"}', '{\"application_id\": \"test_app_001\", \"new_status\": \"MANUAL_REVIEW_REQUIRED\", \"next_step\": \"ASSIGN_TO_REVIEWER\"}', 'SUCCESS', '', 67, '172.20.0.9', 'python-httpx/0.27.2', '2025-05-29 13:12:25.742', '2025-05-29 13:12:25.742');

-- ----------------------------
-- Table structure for ai_analysis_results
-- ----------------------------
DROP TABLE IF EXISTS `ai_analysis_results`;
CREATE TABLE `ai_analysis_results`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `application_id` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL,
  `workflow_execution_id` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL DEFAULT NULL,
  `risk_level` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL,
  `risk_score` decimal(5, 4) NOT NULL,
  `confidence_score` decimal(5, 4) NOT NULL,
  `analysis_summary` text CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL,
  `detailed_analysis` json NULL,
  `risk_factors` json NULL,
  `recommendations` json NULL,
  `a_idecision` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL,
  `approved_amount` decimal(15, 2) NULL DEFAULT NULL,
  `approved_term_months` bigint NULL DEFAULT NULL,
  `suggested_interest_rate` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL DEFAULT NULL,
  `next_action` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL,
  `ai_model_version` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL,
  `processing_time_ms` bigint NOT NULL,
  `processed_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_ai_analysis_results_application_id`(`application_id` ASC) USING BTREE,
  UNIQUE INDEX `idx_ai_analysis_results_workflow_execution_id`(`workflow_execution_id` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 30001 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_bin ROW_FORMAT = Compact;

-- ----------------------------
-- Records of ai_analysis_results
-- ----------------------------
INSERT INTO `ai_analysis_results` VALUES (1, 'test_app_001', 'dify_ai_workflow', 'MEDIUM', 0.5000, 0.5000, 'AI风险分析', '{\"analysis_summary\": \"AI风险分析\", \"approved_amount\": 0, \"approved_term_months\": 12, \"conditions\": [\"需要审核\"], \"confidence_score\": 0.5, \"decision\": \"REQUIRE_HUMAN_REVIEW\", \"detailed_analysis\": {\"credit_analysis\": \"信用分析\", \"financial_analysis\": \"财务分析\", \"risk_factors\": [\"待评估\"], \"strengths\": [\"待评估\"]}, \"recommendations\": [\"建议审核\"], \"risk_level\": \"MEDIUM\", \"risk_score\": 0.5, \"suggested_interest_rate\": \"5.0%\"}', '[{\"description\": \"AI风险分析\", \"factor\": \"risk_score\", \"value\": 0.5}]', '[\"请关注后续通知\"]', 'REQUIRE_HUMAN_REVIEW', 0.00, 12, '5.0%', 'ASSIGN_TO_REVIEWER', 'v1.0.0', 2000, '2025-05-29 13:12:25.684', '2025-05-29 13:12:25.684', '2025-05-29 13:12:25.684');

-- ----------------------------
-- Table structure for ai_model_configs
-- ----------------------------
DROP TABLE IF EXISTS `ai_model_configs`;
CREATE TABLE `ai_model_configs`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `model_id` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL,
  `model_type` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL,
  `version` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL,
  `status` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL DEFAULT 'ACTIVE',
  `configuration` json NULL,
  `thresholds` json NULL,
  `description` text CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL,
  `created_by` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL,
  `activated_at` datetime(3) NULL DEFAULT NULL,
  `deactivated_at` datetime(3) NULL DEFAULT NULL,
  `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `idx_ai_model_configs_model_id`(`model_id` ASC) USING BTREE,
  INDEX `idx_ai_model_configs_model_type`(`model_type` ASC) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_bin ROW_FORMAT = Compact;

-- ----------------------------
-- Records of ai_model_configs
-- ----------------------------

-- ----------------------------
-- Table structure for external_data_queries
-- ----------------------------
DROP TABLE IF EXISTS `external_data_queries`;
CREATE TABLE `external_data_queries`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `query_id` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL,
  `user_id` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL,
  `application_id` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL DEFAULT NULL,
  `data_types` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL,
  `query_result` json NULL,
  `status` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL,
  `error_message` text CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL,
  `query_duration` bigint NULL DEFAULT NULL,
  `queried_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_external_data_queries_application_id`(`application_id` ASC) USING BTREE,
  UNIQUE INDEX `idx_external_data_queries_query_id`(`query_id` ASC) USING BTREE,
  INDEX `idx_external_data_queries_user_id`(`user_id` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 30001 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_bin ROW_FORMAT = Compact;

-- ----------------------------
-- Records of external_data_queries
-- ----------------------------
INSERT INTO `external_data_queries` VALUES (1, '0fdedac5-af10-4d17-a243-f6a135f01131', 'user_001', '', 'credit_report,bank_flow,blacklist_check,government_subsidy', '{\"bank_flow\": {\"account_stability\": \"稳定\", \"average_monthly_income\": 6500, \"last_6_months_flow\": [{\"expense\": 4800, \"income\": 7200, \"month\": \"2024-02\"}, {\"expense\": 5100, \"income\": 6800, \"month\": \"2024-01\"}]}, \"blacklist_check\": {\"check_time\": \"2025-05-29T13:10:29+08:00\", \"is_blacklisted\": false}, \"credit_report\": {\"grade\": \"优秀\", \"loan_history\": [], \"overdue_records\": 0, \"report_date\": \"2024-03-01\", \"score\": 750}, \"government_subsidy\": {\"received_subsidies\": [{\"amount\": 1200, \"type\": \"种粮补贴\", \"year\": 2023}, {\"amount\": 3000, \"type\": \"农机购置补贴\", \"year\": 2022}]}, \"user_id\": \"user_001\"}', 'SUCCESS', '', 0, '2025-05-29 13:10:29.193', '2025-05-29 13:10:29.193', '2025-05-29 13:10:29.193');

-- ----------------------------
-- Table structure for farm_machinery
-- ----------------------------
DROP TABLE IF EXISTS `farm_machinery`;
CREATE TABLE `farm_machinery`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `machinery_id` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL,
  `owner_user_id` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL,
  `type` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL,
  `brand_model` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL DEFAULT NULL,
  `description` text CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL,
  `images` json NULL,
  `daily_rent` decimal(10, 2) NOT NULL,
  `deposit` decimal(10, 2) NULL DEFAULT NULL,
  `location_text` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL DEFAULT NULL,
  `location_geo` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL DEFAULT NULL,
  `status` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL DEFAULT 'AVAILABLE',
  `published_at` datetime(3) NULL DEFAULT NULL,
  `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_farm_machinery_owner_user_id`(`owner_user_id` ASC) USING BTREE,
  INDEX `idx_farm_machinery_type`(`type` ASC) USING BTREE,
  UNIQUE INDEX `idx_farm_machinery_machinery_id`(`machinery_id` ASC) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_bin ROW_FORMAT = Compact;

-- ----------------------------
-- Records of farm_machinery
-- ----------------------------

-- ----------------------------
-- Table structure for loan_application_history
-- ----------------------------
DROP TABLE IF EXISTS `loan_application_history`;
CREATE TABLE `loan_application_history`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `application_id` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL,
  `status_from` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL DEFAULT NULL,
  `status_to` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL,
  `operator_type` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL,
  `operator_id` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL DEFAULT NULL,
  `comments` text CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL,
  `occurred_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_loan_application_history_application_id`(`application_id` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 30001 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_bin ROW_FORMAT = Compact;

-- ----------------------------
-- Records of loan_application_history
-- ----------------------------
INSERT INTO `loan_application_history` VALUES (1, 'test_app_001', 'MANUAL_REVIEW_REQUIRED', 'MANUAL_REVIEW_REQUIRED', 'AI_SYSTEM', 'ai_agent', 'AI决策: REQUIRE_HUMAN_REVIEW, 风险评分: 0.50', '2025-05-29 13:12:25.718');

-- ----------------------------
-- Table structure for loan_applications
-- ----------------------------
DROP TABLE IF EXISTS `loan_applications`;
CREATE TABLE `loan_applications`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `application_id` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL,
  `user_id` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL,
  `product_id` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL,
  `amount_applied` decimal(15, 2) NOT NULL,
  `term_months_applied` bigint NOT NULL,
  `purpose` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL DEFAULT NULL,
  `status` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL,
  `applicant_snapshot` json NULL,
  `submitted_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `ai_risk_score` bigint NULL DEFAULT NULL,
  `ai_suggestion` text CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL,
  `approved_amount` decimal(15, 2) NULL DEFAULT NULL,
  `approved_term_months` bigint NULL DEFAULT NULL,
  `final_decision` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL DEFAULT NULL,
  `decision_reason` text CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL,
  `processed_by` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL DEFAULT NULL,
  `processed_at` datetime(3) NULL DEFAULT NULL,
  `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_loan_applications_product_id`(`product_id` ASC) USING BTREE,
  INDEX `idx_loan_applications_status`(`status` ASC) USING BTREE,
  UNIQUE INDEX `idx_loan_applications_application_id`(`application_id` ASC) USING BTREE,
  INDEX `idx_loan_applications_user_id`(`user_id` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 139120 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_bin ROW_FORMAT = Compact;

-- ----------------------------
-- Records of loan_applications
-- ----------------------------
INSERT INTO `loan_applications` VALUES (1, 'test_app_001', 'user_001', 'product_001', 100000.00, 24, '装修', 'MANUAL_REVIEW_REQUIRED', '{\"address\": \"北京市朝阳区测试街道123号\", \"annual_income\": 200000, \"birth_date\": \"1990-01-01\", \"credit_level\": \"good\", \"education\": \"bachelor\", \"gender\": \"male\", \"has_car\": false, \"has_house\": true, \"id_card_number\": \"110101199001011234\", \"marital_status\": \"single\", \"occupation\": \"软件工程师\", \"phone\": \"13800138001\", \"real_name\": \"张三\", \"work_years\": 5}', '2025-05-29 13:10:14.768', 500, 'AI风险分析', NULL, NULL, '', '', '', NULL, '2025-05-29 13:10:14.768', '2025-05-29 13:12:25.711');

-- ----------------------------
-- Table structure for loan_products
-- ----------------------------
DROP TABLE IF EXISTS `loan_products`;
CREATE TABLE `loan_products`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `product_id` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL,
  `name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL,
  `description` text CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL,
  `category` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL DEFAULT NULL,
  `min_amount` decimal(15, 2) NOT NULL,
  `max_amount` decimal(15, 2) NOT NULL,
  `min_term_months` bigint NOT NULL,
  `max_term_months` bigint NOT NULL,
  `interest_rate_yearly` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL,
  `repayment_methods` json NULL,
  `application_conditions` text CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL,
  `required_documents` json NULL,
  `status` tinyint NOT NULL DEFAULT 0,
  `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `idx_loan_products_product_id`(`product_id` ASC) USING BTREE,
  INDEX `idx_loan_products_category`(`category` ASC) USING BTREE,
  INDEX `idx_loan_products_status`(`status` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 165180 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_bin ROW_FORMAT = Compact;

-- ----------------------------
-- Records of loan_products
-- ----------------------------
INSERT INTO `loan_products` VALUES (1, 'lp_b5d1ac2f-057', '春耕助力贷', '专为春耕生产设计，利率优惠，快速审批', '种植贷', 5000.00, 50000.00, 6, 24, '4.5% - 6.0%', '[\"等额本息\", \"先息后本\"]', '1. 年满18周岁的农户；2. 有稳定的农业收入；3. 信用记录良好', '[{\"desc\": \"申请人身份证\", \"type\": \"ID_CARD\"}, {\"desc\": \"土地承包合同\", \"type\": \"LAND_CONTRACT\"}]', 0, '2025-05-29 13:10:14.289', '2025-05-29 13:10:14.289');
INSERT INTO `loan_products` VALUES (2, 'lp_c274637a-d46', '农机购置贷', '支持农户购买农业机械，助力农业现代化', '设备贷', 10000.00, 200000.00, 12, 60, '5.0% - 7.0%', '[\"等额本息\", \"等额本金\"]', '1. 年满18周岁的农户；2. 有购机需求证明；3. 有还款能力', '[{\"desc\": \"申请人身份证\", \"type\": \"ID_CARD\"}, {\"desc\": \"农机购买合同\", \"type\": \"PURCHASE_CONTRACT\"}]', 0, '2025-05-29 13:10:14.316', '2025-05-29 13:10:14.316');
INSERT INTO `loan_products` VALUES (3, 'lp_70b183c0-af4', '丰收种植贷', '支持大棚种植、果蔬种植等现代农业项目', '种植贷', 8000.00, 80000.00, 12, 36, '4.8% - 6.5%', '[\"等额本息\", \"季末还息到期还本\"]', '1. 有种植经验；2. 有土地使用权证明；3. 无不良信用记录', '[{\"desc\": \"身份证\", \"type\": \"ID_CARD\"}, {\"desc\": \"土地使用权证\", \"type\": \"LAND_USE_CERT\"}]', 0, '2025-05-29 13:10:14.333', '2025-05-29 13:10:14.333');
INSERT INTO `loan_products` VALUES (4, 'lp_38f279e8-fda', '养殖创业贷', '支持家禽、家畜养殖业发展，助力规模化养殖', '养殖贷', 15000.00, 150000.00, 6, 48, '5.2% - 7.5%', '[\"等额本息\", \"按季付息到期还本\"]', '1. 有养殖场地；2. 有相关养殖技术；3. 有销售渠道', '[{\"desc\": \"身份证\", \"type\": \"ID_CARD\"}, {\"desc\": \"养殖场证明\", \"type\": \"FARM_CERT\"}]', 0, '2025-05-29 13:10:14.345', '2025-05-29 13:10:14.345');
INSERT INTO `loan_products` VALUES (5, 'lp_37d6681e-bcf', '智慧农业贷', '支持智能灌溉、物联网设备等现代农业技术应用', '设备贷', 20000.00, 300000.00, 18, 72, '4.8% - 6.8%', '[\"等额本息\", \"等额本金\"]', '1. 有现代农业发展规划；2. 有技术团队支持；3. 有稳定收入来源', '[{\"desc\": \"身份证\", \"type\": \"ID_CARD\"}, {\"desc\": \"项目计划书\", \"type\": \"BUSINESS_PLAN\"}]', 0, '2025-05-29 13:10:14.360', '2025-05-29 13:10:14.360');
INSERT INTO `loan_products` VALUES (6, 'lp_191cc297-a87', '农村电商贷', '支持农产品电商平台、直播带货等新业态发展', '经营贷', 3000.00, 50000.00, 3, 24, '6.0% - 8.5%', '[\"等额本息\", \"按月付息到期还本\"]', '1. 有电商运营经验；2. 有农产品货源；3. 有良好信用记录', '[{\"desc\": \"身份证\", \"type\": \"ID_CARD\"}, {\"desc\": \"营业执照\", \"type\": \"BUSINESS_LICENSE\"}]', 0, '2025-05-29 13:10:14.372', '2025-05-29 13:10:14.372');
INSERT INTO `loan_products` VALUES (7, 'lp_46f44e51-bcf', '合作社发展贷', '支持农民专业合作社扩大经营规模，提升服务能力', '经营贷', 50000.00, 500000.00, 12, 60, '4.2% - 6.2%', '[\"等额本息\", \"按季付息到期还本\"]', '1. 注册满1年的合作社；2. 有稳定的经营收入；3. 无重大违法记录', '[{\"desc\": \"合作社执照\", \"type\": \"COOP_LICENSE\"}, {\"desc\": \"财务报表\", \"type\": \"FINANCIAL_REPORT\"}]', 0, '2025-05-29 13:10:14.389', '2025-05-29 13:10:14.389');
INSERT INTO `loan_products` VALUES (8, 'lp_332b32fc-a28', '绿色农业贷', '支持有机农业、生态农业等绿色环保项目', '种植贷', 12000.00, 120000.00, 12, 48, '4.0% - 5.8%', '[\"等额本息\", \"季末还息到期还本\"]', '1. 有绿色认证或申请中；2. 有环保设施；3. 有市场销路', '[{\"desc\": \"身份证\", \"type\": \"ID_CARD\"}, {\"desc\": \"绿色认证书\", \"type\": \"GREEN_CERT\"}]', 0, '2025-05-29 13:10:14.407', '2025-05-29 13:10:14.407');
INSERT INTO `loan_products` VALUES (9, 'lp_5373d538-a6b', '水产养殖贷', '支持鱼类、虾类、蟹类等水产养殖业发展', '养殖贷', 8000.00, 100000.00, 6, 36, '5.5% - 7.8%', '[\"等额本息\", \"按季付息到期还本\"]', '1. 有水产养殖经验；2. 有养殖场地；3. 有销售渠道', '[{\"desc\": \"身份证\", \"type\": \"ID_CARD\"}, {\"desc\": \"养殖场地证明\", \"type\": \"POND_CERT\"}]', 0, '2025-05-29 13:10:14.427', '2025-05-29 13:10:14.427');
INSERT INTO `loan_products` VALUES (10, 'lp_6127de65-a77', '温室大棚贷', '支持现代化温室大棚建设和设备采购', '设备贷', 30000.00, 500000.00, 24, 84, '4.5% - 6.5%', '[\"等额本息\", \"等额本金\"]', '1. 有大棚建设规划；2. 有土地使用权；3. 有技术支持', '[{\"desc\": \"身份证\", \"type\": \"ID_CARD\"}, {\"desc\": \"建设规划书\", \"type\": \"CONSTRUCTION_PLAN\"}]', 0, '2025-05-29 13:10:14.443', '2025-05-29 13:10:14.443');
INSERT INTO `loan_products` VALUES (11, 'lp_6b037a46-c98', '乡村旅游贷', '支持农家乐、民宿等乡村旅游项目发展', '经营贷', 20000.00, 300000.00, 12, 60, '5.8% - 8.0%', '[\"等额本息\", \"按季付息到期还本\"]', '1. 有旅游资源；2. 有经营许可；3. 有市场调研', '[{\"desc\": \"身份证\", \"type\": \"ID_CARD\"}, {\"desc\": \"旅游经营许可证\", \"type\": \"TOURISM_LICENSE\"}]', 0, '2025-05-29 13:10:14.459', '2025-05-29 13:10:14.459');
INSERT INTO `loan_products` VALUES (12, 'lp_9e130168-3f9', '农产品加工贷', '支持农产品深加工设备采购和技术升级', '设备贷', 25000.00, 400000.00, 18, 72, '4.8% - 6.8%', '[\"等额本息\", \"等额本金\"]', '1. 有加工经验；2. 有原料供应；3. 有销售渠道', '[{\"desc\": \"身份证\", \"type\": \"ID_CARD\"}, {\"desc\": \"食品加工许可证\", \"type\": \"PROCESSING_LICENSE\"}]', 0, '2025-05-29 13:10:14.471', '2025-05-29 13:10:14.471');
INSERT INTO `loan_products` VALUES (13, 'lp_0eac1f4e-46d', '果树种植贷', '支持果园建设、果树种植和果品销售', '种植贷', 10000.00, 150000.00, 12, 60, '4.5% - 6.8%', '[\"等额本息\", \"季末还息到期还本\"]', '1. 有种植经验；2. 有土地承包权；3. 有销售计划', '[{\"desc\": \"身份证\", \"type\": \"ID_CARD\"}, {\"desc\": \"果园种植计划\", \"type\": \"ORCHARD_PLAN\"}]', 0, '2025-05-29 13:10:14.487', '2025-05-29 13:10:14.487');
INSERT INTO `loan_products` VALUES (14, 'lp_480c20f2-66c', '畜牧设备贷', '支持现代化畜牧设备采购和养殖场建设', '设备贷', 15000.00, 250000.00, 12, 60, '5.0% - 7.2%', '[\"等额本息\", \"等额本金\"]', '1. 有畜牧经验；2. 有养殖场地；3. 有环保手续', '[{\"desc\": \"身份证\", \"type\": \"ID_CARD\"}, {\"desc\": \"畜牧养殖许可证\", \"type\": \"LIVESTOCK_PERMIT\"}]', 0, '2025-05-29 13:10:14.505', '2025-05-29 13:10:14.505');
INSERT INTO `loan_products` VALUES (15, 'lp_2341b8c9-794', '农业科技贷', '支持农业科技创新项目和技术研发投入', '经营贷', 30000.00, 600000.00, 18, 72, '4.2% - 6.0%', '[\"等额本息\", \"按季付息到期还本\"]', '1. 有技术团队；2. 有创新项目；3. 有市场前景', '[{\"desc\": \"身份证\", \"type\": \"ID_CARD\"}, {\"desc\": \"技术创新计划书\", \"type\": \"TECH_PLAN\"}]', 0, '2025-05-29 13:10:14.521', '2025-05-29 13:10:14.521');
INSERT INTO `loan_products` VALUES (16, 'product_001', '个人信用贷', '专为个人信用贷款设计，无需抵押，快速审批', 'personal_credit', 10000.00, 500000.00, 6, 36, '8.5%', '[\"等额本息\", \"先息后本\"]', '年收入不低于10万，征信良好', '[{\"desc\": \"身份证\", \"type\": \"ID_CARD\"}, {\"desc\": \"收入证明\", \"type\": \"INCOME_PROOF\"}]', 0, '2025-05-29 13:10:14.752', '2025-05-29 13:10:14.752');

-- ----------------------------
-- Table structure for machinery_leasing_orders
-- ----------------------------
DROP TABLE IF EXISTS `machinery_leasing_orders`;
CREATE TABLE `machinery_leasing_orders`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `order_id` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL,
  `machinery_id` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL,
  `lessee_user_id` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL,
  `lessor_user_id` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL,
  `start_date` date NOT NULL,
  `end_date` date NOT NULL,
  `total_rent` decimal(10, 2) NOT NULL,
  `deposit_amount` decimal(10, 2) NULL DEFAULT NULL,
  `status` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL,
  `lessee_notes` text CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL,
  `lessor_notes` text CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL,
  `confirmed_at` datetime(3) NULL DEFAULT NULL,
  `completed_at` datetime(3) NULL DEFAULT NULL,
  `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_machinery_leasing_orders_lessor_user_id`(`lessor_user_id` ASC) USING BTREE,
  INDEX `idx_machinery_leasing_orders_status`(`status` ASC) USING BTREE,
  UNIQUE INDEX `idx_machinery_leasing_orders_order_id`(`order_id` ASC) USING BTREE,
  INDEX `idx_machinery_leasing_orders_machinery_id`(`machinery_id` ASC) USING BTREE,
  INDEX `idx_machinery_leasing_orders_lessee_user_id`(`lessee_user_id` ASC) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_bin ROW_FORMAT = Compact;

-- ----------------------------
-- Records of machinery_leasing_orders
-- ----------------------------

-- ----------------------------
-- Table structure for oa_users
-- ----------------------------
DROP TABLE IF EXISTS `oa_users`;
CREATE TABLE `oa_users`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `oa_user_id` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL,
  `username` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL,
  `password_hash` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL,
  `role` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL,
  `display_name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL DEFAULT NULL,
  `email` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL DEFAULT NULL,
  `status` tinyint NOT NULL DEFAULT 0,
  `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `idx_oa_users_oa_user_id`(`oa_user_id` ASC) USING BTREE,
  UNIQUE INDEX `idx_oa_users_username`(`username` ASC) USING BTREE,
  INDEX `idx_oa_users_role`(`role` ASC) USING BTREE,
  UNIQUE INDEX `idx_oa_users_email`(`email` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 240545 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_bin ROW_FORMAT = Compact;

-- ----------------------------
-- Records of oa_users
-- ----------------------------
INSERT INTO `oa_users` VALUES (1, 'oa_admin001', 'admin', '$2a$10$5HCRtuvI.2f8.MRQEk7Gze5suCutkgf1uhuhIYLaiak5KdVazs7ga', 'ADMIN', '系统管理员', 'admin@example.com', 0, '2025-05-29 13:10:14.590', '2025-05-29 13:10:14.590');

-- ----------------------------
-- Table structure for system_configurations
-- ----------------------------
DROP TABLE IF EXISTS `system_configurations`;
CREATE TABLE `system_configurations`  (
  `config_key` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL,
  `config_value` text CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL,
  `description` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL DEFAULT NULL,
  `updated_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`config_key`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_bin ROW_FORMAT = Compact;

-- ----------------------------
-- Records of system_configurations
-- ----------------------------
INSERT INTO `system_configurations` VALUES ('ai_approval_enabled', 'true', 'AI审批功能开关', '2025-05-29 13:10:18.949');

-- ----------------------------
-- Table structure for uploaded_files
-- ----------------------------
DROP TABLE IF EXISTS `uploaded_files`;
CREATE TABLE `uploaded_files`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `file_id` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL,
  `user_id` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL,
  `file_name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL,
  `file_type` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL DEFAULT NULL,
  `file_size` bigint NULL DEFAULT NULL,
  `storage_path` varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL,
  `purpose` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL DEFAULT NULL,
  `related_id` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL DEFAULT NULL,
  `uploaded_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_uploaded_files_related_id`(`related_id` ASC) USING BTREE,
  UNIQUE INDEX `idx_uploaded_files_file_id`(`file_id` ASC) USING BTREE,
  INDEX `idx_uploaded_files_user_id`(`user_id` ASC) USING BTREE,
  INDEX `idx_uploaded_files_purpose`(`purpose` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 161472 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_bin ROW_FORMAT = Compact;

-- ----------------------------
-- Records of uploaded_files
-- ----------------------------
INSERT INTO `uploaded_files` VALUES (1, 'file_001', 'user_001', '身份证正面.jpg', 'image/jpeg', 1024000, '/uploads/user_001/id_card_front.jpg', 'id_card', '', '2025-05-29 13:10:14.792', '2025-05-29 13:10:14.792', '2025-05-29 13:10:14.792');
INSERT INTO `uploaded_files` VALUES (2, 'file_002', 'user_001', '收入证明.pdf', 'application/pdf', 2048000, '/uploads/user_001/income_proof.pdf', 'income_proof', '', '2025-05-29 13:10:14.818', '2025-05-29 13:10:14.818', '2025-05-29 13:10:14.818');
INSERT INTO `uploaded_files` VALUES (3, 'file_003', 'user_001', '银行流水.pdf', 'application/pdf', 3072000, '/uploads/user_001/bank_statement.pdf', 'bank_statement', '', '2025-05-29 13:10:14.830', '2025-05-29 13:10:14.830', '2025-05-29 13:10:14.830');

-- ----------------------------
-- Table structure for user_profiles
-- ----------------------------
DROP TABLE IF EXISTS `user_profiles`;
CREATE TABLE `user_profiles`  (
  `user_id` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL,
  `real_name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL DEFAULT NULL,
  `id_card_number` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL DEFAULT NULL,
  `id_card_front_url` varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL DEFAULT NULL,
  `id_card_back_url` varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL DEFAULT NULL,
  `address` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL DEFAULT NULL,
  `gender` tinyint NULL DEFAULT NULL,
  `birth_date` date NULL DEFAULT NULL,
  `occupation` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL DEFAULT NULL,
  `annual_income` decimal(15, 2) NULL DEFAULT NULL,
  `credit_auth_agreed` tinyint(1) NOT NULL DEFAULT 0,
  `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`user_id`) USING BTREE,
  INDEX `idx_user_profiles_id_card_number`(`id_card_number` ASC) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_bin ROW_FORMAT = Compact;

-- ----------------------------
-- Records of user_profiles
-- ----------------------------
INSERT INTO `user_profiles` VALUES ('user_001', '张三', '110101199001011234', '', '', '北京市朝阳区测试街道123号', 1, '1990-01-01', '软件工程师', 200000.00, 1, '2025-05-29 13:10:14.737', '2025-05-29 13:10:14.737');
INSERT INTO `user_profiles` VALUES ('usr_e538f08c-21a', '张三', '31010119900101****', '', '', '上海市浦东新区XX镇XX村', 0, NULL, '', 0.00, 1, '2025-05-29 13:10:14.707', '2025-05-29 13:10:14.707');

-- ----------------------------
-- Table structure for users
-- ----------------------------
DROP TABLE IF EXISTS `users`;
CREATE TABLE `users`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `user_id` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL,
  `phone` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL,
  `password_hash` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL,
  `nickname` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL DEFAULT NULL,
  `avatar_url` varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL DEFAULT NULL,
  `status` tinyint NOT NULL DEFAULT 0,
  `registered_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `last_login_at` datetime(3) NULL DEFAULT NULL,
  `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `idx_users_user_id`(`user_id` ASC) USING BTREE,
  UNIQUE INDEX `idx_users_phone`(`phone` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 120464 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_bin ROW_FORMAT = Compact;

-- ----------------------------
-- Records of users
-- ----------------------------
INSERT INTO `users` VALUES (1, 'usr_e538f08c-21a', '13800138000', '$2a$10$woEwy9.UoU3sgm3d5XucFOGiDSlbPloauLegPsGx/x1Cnsqb/v2oq', '测试农户', '', 0, '2025-05-29 13:10:14.674', NULL, '2025-05-29 13:10:14.674', '2025-05-29 13:10:14.674');
INSERT INTO `users` VALUES (2, 'user_001', '13800138001', '$2a$10$woEwy9.UoU3sgm3d5XucFOGiDSlbPloauLegPsGx/x1Cnsqb/v2oq', 'AI测试用户', '', 0, '2025-05-29 13:10:14.724', NULL, '2025-05-29 13:10:14.724', '2025-05-29 13:10:14.724');

-- ----------------------------
-- Table structure for workflow_executions
-- ----------------------------
DROP TABLE IF EXISTS `workflow_executions`;
CREATE TABLE `workflow_executions`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `execution_id` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL,
  `application_id` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL,
  `workflow_type` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL,
  `status` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL,
  `priority` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL DEFAULT 'NORMAL',
  `started_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `completed_at` datetime(3) NULL DEFAULT NULL,
  `estimated_completion` datetime(3) NULL DEFAULT NULL,
  `current_stage` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL DEFAULT NULL,
  `progress` bigint NOT NULL DEFAULT 0,
  `error_message` text CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL,
  `metadata` json NULL,
  `callback_url` varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL DEFAULT NULL,
  `retry_count` bigint NOT NULL DEFAULT 0,
  `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `idx_workflow_executions_execution_id`(`execution_id` ASC) USING BTREE,
  INDEX `idx_workflow_executions_application_id`(`application_id` ASC) USING BTREE,
  INDEX `idx_workflow_executions_status`(`status` ASC) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_bin ROW_FORMAT = Compact;

-- ----------------------------
-- Records of workflow_executions
-- ----------------------------

SET FOREIGN_KEY_CHECKS = 1;
