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

 Date: 29/05/2025 12:39:09
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
  INDEX `idx_ai_agent_logs_action_type`(`action_type` ASC) USING BTREE,
  UNIQUE INDEX `idx_ai_agent_logs_log_id`(`log_id` ASC) USING BTREE,
  INDEX `idx_ai_agent_logs_application_id`(`application_id` ASC) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_bin ROW_FORMAT = Compact;

-- ----------------------------
-- Records of ai_agent_logs
-- ----------------------------

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
INSERT INTO `ai_analysis_results` VALUES (1, 'test_app_001', 'dify_ai_workflow_1748493461', 'MEDIUM', 0.5000, 0.5000, 'AI风险分析', '{\"analysis_summary\": \"AI风险分析\", \"approved_amount\": 0, \"approved_term_months\": 12, \"conditions\": [\"需要审核\"], \"confidence_score\": 0.5, \"decision\": \"REQUIRE_HUMAN_REVIEW\", \"detailed_analysis\": {\"credit_analysis\": \"信用分析\", \"financial_analysis\": \"财务分析\", \"risk_factors\": [\"待评估\"], \"strengths\": [\"待评估\"]}, \"recommendations\": [\"建议审核\"], \"risk_level\": \"MEDIUM\", \"risk_score\": 0.5, \"suggested_interest_rate\": \"5.0%\"}', '[\"请关注后续通知\"]', '[\"请关注后续通知\"]', 'REQUIRE_HUMAN_REVIEW', 0.00, 12, '5.0%', 'ASSIGN_TO_REVIEWER', 'v1.0.0', 2000, '2025-05-29 12:37:41.619', '2025-05-29 12:37:37.477', '2025-05-29 12:37:37.477');

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
  INDEX `idx_external_data_queries_user_id`(`user_id` ASC) USING BTREE,
  INDEX `idx_external_data_queries_application_id`(`application_id` ASC) USING BTREE,
  UNIQUE INDEX `idx_external_data_queries_query_id`(`query_id` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 34774 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_bin ROW_FORMAT = Compact;

-- ----------------------------
-- Records of external_data_queries
-- ----------------------------
INSERT INTO `external_data_queries` VALUES (1, 'c223d0aa-67c4-4212-9c03-d5c6c92e4855', 'user_001', '', 'credit_report,bank_flow,blacklist_check,government_subsidy', '{\"bank_flow\": {\"account_stability\": \"稳定\", \"average_monthly_income\": 6500, \"last_6_months_flow\": [{\"expense\": 4800, \"income\": 7200, \"month\": \"2024-02\"}, {\"expense\": 5100, \"income\": 6800, \"month\": \"2024-01\"}]}, \"blacklist_check\": {\"check_time\": \"2025-05-29T12:36:26+08:00\", \"is_blacklisted\": false}, \"credit_report\": {\"grade\": \"优秀\", \"loan_history\": [], \"overdue_records\": 0, \"report_date\": \"2024-03-01\", \"score\": 750}, \"government_subsidy\": {\"received_subsidies\": [{\"amount\": 1200, \"type\": \"种粮补贴\", \"year\": 2023}, {\"amount\": 3000, \"type\": \"农机购置补贴\", \"year\": 2022}]}, \"user_id\": \"user_001\"}', 'SUCCESS', '', 1500, '2025-05-29 12:36:26.637', '2025-05-29 12:36:22.496', '2025-05-29 12:36:22.496');

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
  UNIQUE INDEX `idx_farm_machinery_machinery_id`(`machinery_id` ASC) USING BTREE,
  INDEX `idx_farm_machinery_owner_user_id`(`owner_user_id` ASC) USING BTREE,
  INDEX `idx_farm_machinery_type`(`type` ASC) USING BTREE
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
INSERT INTO `loan_application_history` VALUES (1, 'test_app_001', 'pending_review', 'MANUAL_REVIEW_REQUIRED', 'AI_SYSTEM', 'ai_agent', 'AI决策: REQUIRE_HUMAN_REVIEW, 风险评分: 0.50', '2025-05-29 12:37:41.642');

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
  INDEX `idx_loan_applications_user_id`(`user_id` ASC) USING BTREE,
  INDEX `idx_loan_applications_product_id`(`product_id` ASC) USING BTREE,
  INDEX `idx_loan_applications_status`(`status` ASC) USING BTREE,
  UNIQUE INDEX `idx_loan_applications_application_id`(`application_id` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 136273 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_bin ROW_FORMAT = Compact;

-- ----------------------------
-- Records of loan_applications
-- ----------------------------
INSERT INTO `loan_applications` VALUES (1, 'test_app_001', 'user_001', 'product_001', 100000.00, 24, '装修', 'MANUAL_REVIEW_REQUIRED', '{\"address\": \"北京市朝阳区测试街道123号\", \"annual_income\": 200000, \"birth_date\": \"1990-01-01\", \"credit_level\": \"good\", \"education\": \"bachelor\", \"gender\": \"male\", \"has_car\": false, \"has_house\": true, \"id_card_number\": \"110101199001011234\", \"marital_status\": \"single\", \"occupation\": \"软件工程师\", \"phone\": \"13800138001\", \"real_name\": \"张三\", \"work_years\": 5}', '2025-05-29 12:36:14.573', NULL, 'AI风险分析', NULL, NULL, '', '', '', NULL, '2025-05-29 12:36:14.573', '2025-05-29 12:37:41.634');

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
) ENGINE = InnoDB AUTO_INCREMENT = 161290 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_bin ROW_FORMAT = Compact;

-- ----------------------------
-- Records of loan_products
-- ----------------------------
INSERT INTO `loan_products` VALUES (1, 'lp_c8324f0d-318', '春耕助力贷', '专为春耕生产设计，利率优惠，快速审批', '种植贷', 5000.00, 50000.00, 6, 24, '4.5% - 6.0%', '[\"等额本息\", \"先息后本\"]', '1. 年满18周岁的农户；2. 有稳定的农业收入；3. 信用记录良好', '[{\"desc\": \"申请人身份证\", \"type\": \"ID_CARD\"}, {\"desc\": \"土地承包合同\", \"type\": \"LAND_CONTRACT\"}]', 0, '2025-05-29 12:36:14.085', '2025-05-29 12:36:14.085');
INSERT INTO `loan_products` VALUES (2, 'lp_2f4a314d-dae', '农机购置贷', '支持农户购买农业机械，助力农业现代化', '设备贷', 10000.00, 200000.00, 12, 60, '5.0% - 7.0%', '[\"等额本息\", \"等额本金\"]', '1. 年满18周岁的农户；2. 有购机需求证明；3. 有还款能力', '[{\"desc\": \"申请人身份证\", \"type\": \"ID_CARD\"}, {\"desc\": \"农机购买合同\", \"type\": \"PURCHASE_CONTRACT\"}]', 0, '2025-05-29 12:36:14.113', '2025-05-29 12:36:14.113');
INSERT INTO `loan_products` VALUES (3, 'lp_76f55d87-438', '丰收种植贷', '支持大棚种植、果蔬种植等现代农业项目', '种植贷', 8000.00, 80000.00, 12, 36, '4.8% - 6.5%', '[\"等额本息\", \"季末还息到期还本\"]', '1. 有种植经验；2. 有土地使用权证明；3. 无不良信用记录', '[{\"desc\": \"身份证\", \"type\": \"ID_CARD\"}, {\"desc\": \"土地使用权证\", \"type\": \"LAND_USE_CERT\"}]', 0, '2025-05-29 12:36:14.133', '2025-05-29 12:36:14.133');
INSERT INTO `loan_products` VALUES (4, 'lp_bacc8961-704', '养殖创业贷', '支持家禽、家畜养殖业发展，助力规模化养殖', '养殖贷', 15000.00, 150000.00, 6, 48, '5.2% - 7.5%', '[\"等额本息\", \"按季付息到期还本\"]', '1. 有养殖场地；2. 有相关养殖技术；3. 有销售渠道', '[{\"desc\": \"身份证\", \"type\": \"ID_CARD\"}, {\"desc\": \"养殖场证明\", \"type\": \"FARM_CERT\"}]', 0, '2025-05-29 12:36:14.149', '2025-05-29 12:36:14.149');
INSERT INTO `loan_products` VALUES (5, 'lp_148d136d-ad5', '智慧农业贷', '支持智能灌溉、物联网设备等现代农业技术应用', '设备贷', 20000.00, 300000.00, 18, 72, '4.8% - 6.8%', '[\"等额本息\", \"等额本金\"]', '1. 有现代农业发展规划；2. 有技术团队支持；3. 有稳定收入来源', '[{\"desc\": \"身份证\", \"type\": \"ID_CARD\"}, {\"desc\": \"项目计划书\", \"type\": \"BUSINESS_PLAN\"}]', 0, '2025-05-29 12:36:14.165', '2025-05-29 12:36:14.165');
INSERT INTO `loan_products` VALUES (6, 'lp_47f109b3-d83', '农村电商贷', '支持农产品电商平台、直播带货等新业态发展', '经营贷', 3000.00, 50000.00, 3, 24, '6.0% - 8.5%', '[\"等额本息\", \"按月付息到期还本\"]', '1. 有电商运营经验；2. 有农产品货源；3. 有良好信用记录', '[{\"desc\": \"身份证\", \"type\": \"ID_CARD\"}, {\"desc\": \"营业执照\", \"type\": \"BUSINESS_LICENSE\"}]', 0, '2025-05-29 12:36:14.180', '2025-05-29 12:36:14.180');
INSERT INTO `loan_products` VALUES (7, 'lp_4c036add-995', '合作社发展贷', '支持农民专业合作社扩大经营规模，提升服务能力', '经营贷', 50000.00, 500000.00, 12, 60, '4.2% - 6.2%', '[\"等额本息\", \"按季付息到期还本\"]', '1. 注册满1年的合作社；2. 有稳定的经营收入；3. 无重大违法记录', '[{\"desc\": \"合作社执照\", \"type\": \"COOP_LICENSE\"}, {\"desc\": \"财务报表\", \"type\": \"FINANCIAL_REPORT\"}]', 0, '2025-05-29 12:36:14.195', '2025-05-29 12:36:14.195');
INSERT INTO `loan_products` VALUES (8, 'lp_a5891af0-16b', '绿色农业贷', '支持有机农业、生态农业等绿色环保项目', '种植贷', 12000.00, 120000.00, 12, 48, '4.0% - 5.8%', '[\"等额本息\", \"季末还息到期还本\"]', '1. 有绿色认证或申请中；2. 有环保设施；3. 有市场销路', '[{\"desc\": \"身份证\", \"type\": \"ID_CARD\"}, {\"desc\": \"绿色认证书\", \"type\": \"GREEN_CERT\"}]', 0, '2025-05-29 12:36:14.210', '2025-05-29 12:36:14.210');
INSERT INTO `loan_products` VALUES (9, 'lp_3129d6f9-ef4', '水产养殖贷', '支持鱼类、虾类、蟹类等水产养殖业发展', '养殖贷', 8000.00, 100000.00, 6, 36, '5.5% - 7.8%', '[\"等额本息\", \"按季付息到期还本\"]', '1. 有水产养殖经验；2. 有养殖场地；3. 有销售渠道', '[{\"desc\": \"身份证\", \"type\": \"ID_CARD\"}, {\"desc\": \"养殖场地证明\", \"type\": \"POND_CERT\"}]', 0, '2025-05-29 12:36:14.225', '2025-05-29 12:36:14.225');
INSERT INTO `loan_products` VALUES (10, 'lp_0c50c857-1e7', '温室大棚贷', '支持现代化温室大棚建设和设备采购', '设备贷', 30000.00, 500000.00, 24, 84, '4.5% - 6.5%', '[\"等额本息\", \"等额本金\"]', '1. 有大棚建设规划；2. 有土地使用权；3. 有技术支持', '[{\"desc\": \"身份证\", \"type\": \"ID_CARD\"}, {\"desc\": \"建设规划书\", \"type\": \"CONSTRUCTION_PLAN\"}]', 0, '2025-05-29 12:36:14.240', '2025-05-29 12:36:14.240');
INSERT INTO `loan_products` VALUES (11, 'lp_4f82bb2f-85b', '乡村旅游贷', '支持农家乐、民宿等乡村旅游项目发展', '经营贷', 20000.00, 300000.00, 12, 60, '5.8% - 8.0%', '[\"等额本息\", \"按季付息到期还本\"]', '1. 有旅游资源；2. 有经营许可；3. 有市场调研', '[{\"desc\": \"身份证\", \"type\": \"ID_CARD\"}, {\"desc\": \"旅游经营许可证\", \"type\": \"TOURISM_LICENSE\"}]', 0, '2025-05-29 12:36:14.256', '2025-05-29 12:36:14.256');
INSERT INTO `loan_products` VALUES (12, 'lp_ee1996fd-0d8', '农产品加工贷', '支持农产品深加工设备采购和技术升级', '设备贷', 25000.00, 400000.00, 18, 72, '4.8% - 6.8%', '[\"等额本息\", \"等额本金\"]', '1. 有加工经验；2. 有原料供应；3. 有销售渠道', '[{\"desc\": \"身份证\", \"type\": \"ID_CARD\"}, {\"desc\": \"食品加工许可证\", \"type\": \"PROCESSING_LICENSE\"}]', 0, '2025-05-29 12:36:14.274', '2025-05-29 12:36:14.274');
INSERT INTO `loan_products` VALUES (13, 'lp_5e36142a-a8a', '果树种植贷', '支持果园建设、果树种植和果品销售', '种植贷', 10000.00, 150000.00, 12, 60, '4.5% - 6.8%', '[\"等额本息\", \"季末还息到期还本\"]', '1. 有种植经验；2. 有土地承包权；3. 有销售计划', '[{\"desc\": \"身份证\", \"type\": \"ID_CARD\"}, {\"desc\": \"果园种植计划\", \"type\": \"ORCHARD_PLAN\"}]', 0, '2025-05-29 12:36:14.290', '2025-05-29 12:36:14.290');
INSERT INTO `loan_products` VALUES (14, 'lp_a6e27c07-ff4', '畜牧设备贷', '支持现代化畜牧设备采购和养殖场建设', '设备贷', 15000.00, 250000.00, 12, 60, '5.0% - 7.2%', '[\"等额本息\", \"等额本金\"]', '1. 有畜牧经验；2. 有养殖场地；3. 有环保手续', '[{\"desc\": \"身份证\", \"type\": \"ID_CARD\"}, {\"desc\": \"畜牧养殖许可证\", \"type\": \"LIVESTOCK_PERMIT\"}]', 0, '2025-05-29 12:36:14.306', '2025-05-29 12:36:14.306');
INSERT INTO `loan_products` VALUES (15, 'lp_f171b109-d3a', '农业科技贷', '支持农业科技创新项目和技术研发投入', '经营贷', 30000.00, 600000.00, 18, 72, '4.2% - 6.0%', '[\"等额本息\", \"按季付息到期还本\"]', '1. 有技术团队；2. 有创新项目；3. 有市场前景', '[{\"desc\": \"身份证\", \"type\": \"ID_CARD\"}, {\"desc\": \"技术创新计划书\", \"type\": \"TECH_PLAN\"}]', 0, '2025-05-29 12:36:14.319', '2025-05-29 12:36:14.319');
INSERT INTO `loan_products` VALUES (16, 'product_001', '个人信用贷', '专为个人信用贷款设计，无需抵押，快速审批', 'personal_credit', 10000.00, 500000.00, 6, 36, '8.5%', '[\"等额本息\", \"先息后本\"]', '年收入不低于10万，征信良好', '[{\"desc\": \"身份证\", \"type\": \"ID_CARD\"}, {\"desc\": \"收入证明\", \"type\": \"INCOME_PROOF\"}]', 0, '2025-05-29 12:36:14.556', '2025-05-29 12:36:14.556');

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
  INDEX `idx_machinery_leasing_orders_lessee_user_id`(`lessee_user_id` ASC) USING BTREE,
  INDEX `idx_machinery_leasing_orders_lessor_user_id`(`lessor_user_id` ASC) USING BTREE,
  INDEX `idx_machinery_leasing_orders_status`(`status` ASC) USING BTREE,
  UNIQUE INDEX `idx_machinery_leasing_orders_order_id`(`order_id` ASC) USING BTREE,
  INDEX `idx_machinery_leasing_orders_machinery_id`(`machinery_id` ASC) USING BTREE
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
) ENGINE = InnoDB AUTO_INCREMENT = 239277 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_bin ROW_FORMAT = Compact;

-- ----------------------------
-- Records of oa_users
-- ----------------------------
INSERT INTO `oa_users` VALUES (1, 'oa_admin001', 'admin', '$2a$10$DP376J0R2bJmgum0qB6cc.XiYj5E314A7C9PY1xbqEtwnjC3BQ14G', 'ADMIN', '系统管理员', 'admin@example.com', 0, '2025-05-29 12:36:14.398', '2025-05-29 12:36:14.398');

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
INSERT INTO `system_configurations` VALUES ('ai_approval_enabled', 'true', 'AI审批功能开关', '2025-05-29 12:36:18.816');

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
) ENGINE = InnoDB AUTO_INCREMENT = 157663 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_bin ROW_FORMAT = Compact;

-- ----------------------------
-- Records of uploaded_files
-- ----------------------------
INSERT INTO `uploaded_files` VALUES (1, 'file_001', 'user_001', '身份证正面.jpg', 'image/jpeg', 1024000, '/uploads/user_001/id_card_front.jpg', 'id_card', '', '2025-05-29 12:36:14.605', '2025-05-29 12:36:14.605', '2025-05-29 12:36:14.605');
INSERT INTO `uploaded_files` VALUES (2, 'file_002', 'user_001', '收入证明.pdf', 'application/pdf', 2048000, '/uploads/user_001/income_proof.pdf', 'income_proof', '', '2025-05-29 12:36:14.634', '2025-05-29 12:36:14.634', '2025-05-29 12:36:14.634');
INSERT INTO `uploaded_files` VALUES (3, 'file_003', 'user_001', '银行流水.pdf', 'application/pdf', 3072000, '/uploads/user_001/bank_statement.pdf', 'bank_statement', '', '2025-05-29 12:36:14.653', '2025-05-29 12:36:14.653', '2025-05-29 12:36:14.653');

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
INSERT INTO `user_profiles` VALUES ('user_001', '张三', '110101199001011234', '', '', '北京市朝阳区测试街道123号', 1, '1990-01-01', '软件工程师', 200000.00, 1, '2025-05-29 12:36:14.539', '2025-05-29 12:36:14.539');
INSERT INTO `user_profiles` VALUES ('usr_f637ccd0-270', '张三', '31010119900101****', '', '', '上海市浦东新区XX镇XX村', 0, NULL, '', 0.00, 1, '2025-05-29 12:36:14.507', '2025-05-29 12:36:14.507');

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
) ENGINE = InnoDB AUTO_INCREMENT = 116163 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_bin ROW_FORMAT = Compact;

-- ----------------------------
-- Records of users
-- ----------------------------
INSERT INTO `users` VALUES (1, 'usr_f637ccd0-270', '13800138000', '$2a$10$Tuos/MuuBQznuJwO6nBzy.Kl5jFkfsw/235cTWk7GmSeLWcPTlQ7u', '测试农户', '', 0, '2025-05-29 12:36:14.484', NULL, '2025-05-29 12:36:14.484', '2025-05-29 12:36:14.484');
INSERT INTO `users` VALUES (2, 'user_001', '13800138001', '$2a$10$Tuos/MuuBQznuJwO6nBzy.Kl5jFkfsw/235cTWk7GmSeLWcPTlQ7u', 'AI测试用户', '', 0, '2025-05-29 12:36:14.524', NULL, '2025-05-29 12:36:14.524', '2025-05-29 12:36:14.524');

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
