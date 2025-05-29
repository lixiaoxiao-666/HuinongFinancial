/*
 数字惠农后端 - 模拟数据插入脚本
 
 使用说明：
 1. 确保数据库已经创建并导入了 app.sql 结构
 2. 执行此脚本前请备份现有数据
 3. 此脚本会插入测试数据，请勿在生产环境使用
 
 执行方法：
 mysql -u username -p database_name < mock_data.sql
 
 Date: 30/05/2025
*/

USE huinong_db;

-- 禁用外键检查，避免插入顺序问题
SET FOREIGN_KEY_CHECKS = 0;

-- 清空所有表数据（可选，谨慎使用）
-- TRUNCATE TABLE api_logs;
-- TRUNCATE TABLE approval_logs;
-- TRUNCATE TABLE articles;
-- TRUNCATE TABLE categories;
-- TRUNCATE TABLE dify_workflow_logs;
-- TRUNCATE TABLE experts;
-- TRUNCATE TABLE file_uploads;
-- TRUNCATE TABLE loan_applications;
-- TRUNCATE TABLE loan_products;
-- TRUNCATE TABLE machines;
-- TRUNCATE TABLE oa_roles;
-- TRUNCATE TABLE oa_users;
-- TRUNCATE TABLE offline_queue;
-- TRUNCATE TABLE rental_orders;
-- TRUNCATE TABLE system_configs;
-- TRUNCATE TABLE user_auths;
-- TRUNCATE TABLE user_sessions;
-- TRUNCATE TABLE user_tags;
-- TRUNCATE TABLE users;

-- ========== 1. 系统配置数据 ==========
INSERT INTO system_configs (config_key, config_value, config_type, config_group, description, is_editable, is_encrypted, created_at, updated_at) VALUES
('app.name', '数字惠农金融服务平台', 'string', 'app', '应用程序名称', 1, 0, NOW(), NOW()),
('app.version', '1.0.0', 'string', 'app', '应用程序版本', 1, 0, NOW(), NOW()),
('app.maintenance_mode', 'false', 'boolean', 'app', '维护模式开关', 1, 0, NOW(), NOW()),
('sms.daily_limit', '100', 'number', 'sms', '每日短信发送限制', 1, 0, NOW(), NOW()),
('loan.max_amount', '1000000', 'number', 'loan', '贷款最大金额(分)', 1, 0, NOW(), NOW()),
('loan.min_amount', '10000', 'number', 'loan', '贷款最小金额(分)', 1, 0, NOW(), NOW()),
('machine.rental_deposit_rate', '0.2', 'number', 'machine', '农机租赁押金比例', 1, 0, NOW(), NOW()),
('security.password_min_length', '8', 'number', 'security', '密码最小长度', 1, 0, NOW(), NOW()),
('dify.default_timeout', '30', 'number', 'dify', 'Dify API默认超时时间', 1, 0, NOW(), NOW()),
('file.upload_max_size', '10485760', 'number', 'file', '文件上传最大大小(字节)', 1, 0, NOW(), NOW());

-- ========== 2. OA角色数据 ==========
INSERT INTO oa_roles (name, display_name, description, permissions, is_super, status, created_at, updated_at) VALUES
('super_admin', '超级管理员', '系统超级管理员，拥有所有权限', '["*"]', 1, 'active', NOW(), NOW()),
('loan_manager', '贷款经理', '负责贷款业务管理', '["loan.*", "user.view", "approval.*"]', 0, 'active', NOW(), NOW()),
('risk_manager', '风控经理', '负责风险评估和控制', '["loan.view", "approval.*", "risk.*"]', 0, 'active', NOW(), NOW()),
('machine_manager', '农机管理员', '负责农机设备管理', '["machine.*", "user.view", "order.*"]', 0, 'active', NOW(), NOW()),
('content_manager', '内容管理员', '负责内容和专家管理', '["article.*", "expert.*", "category.*"]', 0, 'active', NOW(), NOW()),
('customer_service', '客服专员', '负责客户服务支持', '["user.view", "user.update", "message.*"]', 0, 'active', NOW(), NOW());

-- ========== 3. OA用户数据 ==========
INSERT INTO oa_users (username, email, phone, password_hash, salt, real_name, role_id, department, position, status, last_login_at, last_login_ip, login_count, created_at, updated_at) VALUES
('admin', 'admin@huinong.com', '13800000001', '$2a$10$N.zmdr9k7uOCQb376NoUnuTGJmKWNbYNhYeQ79vNgtV3OIGfftP5C', 'admin_salt_001', '系统管理员', 1, '技术部', '系统管理员', 'active', NOW(), '127.0.0.1', 10, NOW(), NOW()),
('loan001', 'loan.manager@huinong.com', '13800000002', '$2a$10$N.zmdr9k7uOCQb376NoUnuTGJmKWNbYNhYeQ79vNgtV3OIGfftP5C', 'loan_salt_001', '张贷款', 2, '业务部', '贷款经理', 'active', NOW(), '192.168.1.101', 25, NOW(), NOW()),
('risk001', 'risk.manager@huinong.com', '13800000003', '$2a$10$N.zmdr9k7uOCQb376NoUnuTGJmKWNbYNhYeQ79vNgtV3OIGfftP5C', 'risk_salt_001', '李风控', 3, '风控部', '风控经理', 'active', NOW(), '192.168.1.102', 18, NOW(), NOW()),
('machine001', 'machine.manager@huinong.com', '13800000004', '$2a$10$N.zmdr9k7uOCQb376NoUnuTGJmKWNbYNhYeQ79vNgtV3OIGfftP5C', 'machine_salt_001', '王农机', 4, '农机部', '农机管理员', 'active', NOW(), '192.168.1.103', 32, NOW(), NOW()),
('content001', 'content.manager@huinong.com', '13800000005', '$2a$10$N.zmdr9k7uOCQb376NoUnuTGJmKWNbYNhYeQ79vNgtV3OIGfftP5C', 'content_salt_001', '陈内容', 5, '运营部', '内容管理员', 'active', NOW(), '192.168.1.104', 15, NOW(), NOW()),
('service001', 'service@huinong.com', '13800000006', '$2a$10$N.zmdr9k7uOCQb376NoUnuTGJmKWNbYNhYeQ79vNgtV3OIGfftP5C', 'service_salt_001', '赵客服', 6, '客服部', '客服专员', 'active', NOW(), '192.168.1.105', 42, NOW(), NOW());

-- ========== 4. 用户数据 ==========
INSERT INTO users (uuid, username, phone, email, password_hash, salt, user_type, status, real_name, id_card, gender, birthday, province, city, county, address, is_real_name_verified, is_bank_card_verified, is_credit_verified, last_login_time, last_login_ip, login_count, register_ip, register_time, created_at, updated_at) VALUES
('f47ac10b-58cc-4372-a567-0e02b2c3d479', 'farmer001', '13900000001', 'farmer001@example.com', '$2a$10$N.zmdr9k7uOCQb376NoUnuTGJmKWNbYNhYeQ79vNgtV3OIGfftP5C', 'salt_001', 'farmer', 'active', '李大牛', '110101199001011234', 'male', '1990-01-01', '山东省', '济南市', '章丘区', '相公庄街道李家村123号', 1, 1, 1, NOW(), '192.168.1.201', 25, '192.168.1.201', NOW(), NOW(), NOW()),
('f47ac10b-58cc-4372-a567-0e02b2c3d480', 'farmer002', '13900000002', 'farmer002@example.com', '$2a$10$N.zmdr9k7uOCQb376NoUnuTGJmKWNbYNhYeQ79vNgtV3OIGfftP5C', 'salt_002', 'farmer', 'active', '王二麻子', '110101199002021235', 'male', '1990-02-02', '河南省', '郑州市', '中牟县', '韩寺镇王家村456号', 1, 1, 0, NOW(), '192.168.1.202', 18, '192.168.1.202', NOW(), NOW(), NOW()),
('f47ac10b-58cc-4372-a567-0e02b2c3d481', 'coop001', '13900000003', 'coop001@example.com', '$2a$10$N.zmdr9k7uOCQb376NoUnuTGJmKWNbYNhYeQ79vNgtV3OIGfftP5C', 'salt_003', 'cooperative', 'active', '张合作社', '110101199003031236', 'male', '1990-03-03', '江苏省', '徐州市', '沛县', '安国镇张楼村789号', 1, 1, 1, NOW(), '192.168.1.203', 32, '192.168.1.203', NOW(), NOW(), NOW()),
('f47ac10b-58cc-4372-a567-0e02b2c3d482', 'company001', '13900000004', 'company001@example.com', '$2a$10$N.zmdr9k7uOCQb376NoUnuTGJmKWNbYNhYeQ79vNgtV3OIGfftP5C', 'salt_004', 'enterprise', 'active', '刘企业', '110101199004041237', 'male', '1990-04-04', '安徽省', '合肥市', '肥东县', '店埠镇刘集村101号', 1, 1, 1, NOW(), '192.168.1.204', 15, '192.168.1.204', NOW(), NOW(), NOW()),
('f47ac10b-58cc-4372-a567-0e02b2c3d483', 'farmer003', '13900000005', 'farmer003@example.com', '$2a$10$N.zmdr9k7uOCQb376NoUnuTGJmKWNbYNhYeQ79vNgtV3OIGfftP5C', 'salt_005', 'farmer', 'active', '陈小花', '110101199005051238', 'female', '1990-05-05', '湖北省', '襄阳市', '老河口市', '竹林桥镇陈家村202号', 1, 0, 0, NOW(), '192.168.1.205', 8, '192.168.1.205', NOW(), NOW(), NOW());

-- ========== 5. 用户认证数据 ==========
INSERT INTO user_auths (user_id, auth_type, auth_status, auth_data, reviewer_id, review_note, reviewed_at, expires_at, created_at, updated_at) VALUES
(1, 'real_name', 'approved', '{"name": "李大牛", "id_card": "110101199001011234", "front_image": "/uploads/id_front_001.jpg", "back_image": "/uploads/id_back_001.jpg"}', 2, '身份信息核实无误', NOW(), DATE_ADD(NOW(), INTERVAL 3 YEAR), NOW(), NOW()),
(1, 'bank_card', 'approved', '{"bank_name": "中国农业银行", "card_number": "6228480123456789012", "cardholder": "李大牛"}', 2, '银行卡信息验证通过', NOW(), DATE_ADD(NOW(), INTERVAL 1 YEAR), NOW(), NOW()),
(1, 'credit', 'approved', '{"credit_score": 750, "debt_ratio": 0.3, "income_verification": true}', 3, '征信记录良好', NOW(), DATE_ADD(NOW(), INTERVAL 6 MONTH), NOW(), NOW()),
(2, 'real_name', 'approved', '{"name": "王二麻子", "id_card": "110101199002021235", "front_image": "/uploads/id_front_002.jpg", "back_image": "/uploads/id_back_002.jpg"}', 2, '身份信息核实无误', NOW(), DATE_ADD(NOW(), INTERVAL 3 YEAR), NOW(), NOW()),
(2, 'bank_card', 'approved', '{"bank_name": "中国建设银行", "card_number": "6217000123456789013", "cardholder": "王二麻子"}', 2, '银行卡信息验证通过', NOW(), DATE_ADD(NOW(), INTERVAL 1 YEAR), NOW(), NOW()),
(3, 'real_name', 'approved', '{"name": "张合作社", "id_card": "110101199003031236", "front_image": "/uploads/id_front_003.jpg", "back_image": "/uploads/id_back_003.jpg"}', 2, '身份信息核实无误', NOW(), DATE_ADD(NOW(), INTERVAL 3 YEAR), NOW(), NOW());

-- ========== 6. 用户标签数据 ==========
INSERT INTO user_tags (user_id, tag_type, tag_key, tag_value, source, creator_id, score, expires_at, created_at, updated_at) VALUES
(1, 'behavior', 'active_level', 'high', 'system', NULL, 85, NULL, NOW(), NOW()),
(1, 'financial', 'credit_level', 'excellent', 'system', NULL, 750, NULL, NOW(), NOW()),
(1, 'business', 'crop_type', 'wheat', 'user', 1, 0, NULL, NOW(), NOW()),
(1, 'business', 'farm_size', 'large', 'system', NULL, 0, NULL, NOW(), NOW()),
(2, 'behavior', 'active_level', 'medium', 'system', NULL, 65, NULL, NOW(), NOW()),
(2, 'financial', 'credit_level', 'good', 'system', NULL, 680, NULL, NOW(), NOW()),
(2, 'business', 'crop_type', 'corn', 'user', 2, 0, NULL, NOW(), NOW()),
(3, 'behavior', 'active_level', 'high', 'system', NULL, 90, NULL, NOW(), NOW()),
(3, 'financial', 'credit_level', 'excellent', 'system', NULL, 780, NULL, NOW(), NOW()),
(3, 'business', 'organization_type', 'cooperative', 'system', NULL, 0, NULL, NOW(), NOW());

-- ========== 7. 贷款产品数据 ==========
INSERT INTO loan_products (product_name, product_code, product_type, description, min_amount, max_amount, interest_rate, term_months, required_auth, is_active, sort_order, dify_workflow_id, eligible_user_type, created_at, updated_at) VALUES
('惠农小额贷', 'HN_SMALL_001', 'small_loan', '面向小农户的快速小额贷款产品，审批快速，手续简便', 1000000, 10000000, 0.0580, 12, '["real_name", "bank_card"]', 1, 1, 'wf_loan_small_001', 'farmer', NOW(), NOW()),
('农机设备贷', 'HN_MACHINE_001', 'equipment_loan', '专为购买农机设备设计的专项贷款', 5000000, 50000000, 0.0520, 36, '["real_name", "bank_card", "credit"]', 1, 2, 'wf_loan_machine_001', 'farmer,cooperative', NOW(), NOW()),
('合作社发展贷', 'HN_COOP_001', 'development_loan', '支持农民专业合作社发展壮大的中长期贷款', 10000000, 100000000, 0.0480, 60, '["real_name", "bank_card", "credit"]', 1, 3, 'wf_loan_coop_001', 'cooperative', NOW(), NOW()),
('企业经营贷', 'HN_ENTERPRISE_001', 'business_loan', '面向农业企业的经营周转资金贷款', 20000000, 200000000, 0.0450, 24, '["real_name", "bank_card", "credit"]', 1, 4, 'wf_loan_enterprise_001', 'enterprise', NOW(), NOW()),
('季节性生产贷', 'HN_SEASONAL_001', 'seasonal_loan', '根据农业生产周期设计的季节性贷款', 500000, 20000000, 0.0600, 6, '["real_name", "bank_card"]', 1, 5, 'wf_loan_seasonal_001', 'farmer,cooperative', NOW(), NOW());

-- ========== 8. 贷款申请数据 ==========
INSERT INTO loan_applications (application_no, user_id, product_id, apply_amount, loan_amount, apply_term_months, term_months, loan_purpose, expected_use_date, applicant_name, applicant_id_card, applicant_phone, contact_phone, contact_email, monthly_income, yearly_income, income_source, other_debts, farm_area, crop_types, years_of_experience, land_certificate, application_documents, materials_json, status, current_approver, approval_level, auto_approval_passed, submitted_at, approved_amount, approved_term_months, approved_rate, credit_score, risk_level, debt_income_ratio, risk_assessment, dify_conversation_id, ai_recommendation, created_at, updated_at) VALUES
('LA20250530001', 1, 1, 5000000, 5000000, 12, 12, '购买种子化肥', DATE_ADD(NOW(), INTERVAL 7 DAY), '李大牛', '110101199001011234', '13900000001', '13900000001', 'farmer001@example.com', 800000, 9600000, '农业种植收入', 0, 50.5, '["wheat", "corn"]', 15, 'LC20231101001', '["application.pdf", "income_proof.pdf"]', '{"guarantor": "村委会", "collateral": "土地使用权"}', 'approved', 2, 3, 1, NOW(), 5000000, 12, 0.0580, 750, 'low', 0.0000, 'AI评估：申请人信用记录良好，农业经营经验丰富，风险较低', 'conv_20250530_001', '建议批准，可按申请金额和期限放款', NOW(), NOW()),
('LA20250530002', 2, 2, 15000000, 12000000, 36, 36, '购买拖拉机设备', DATE_ADD(NOW(), INTERVAL 15 DAY), '王二麻子', '110101199002021235', '13900000002', '13900000002', 'farmer002@example.com', 600000, 7200000, '农业种植收入', 200000, 80.3, '["corn", "soybean"]', 10, 'LC20231102001', '["application.pdf", "equipment_quote.pdf"]', '{"equipment_brand": "约翰迪尔", "equipment_model": "6B-1404"}', 'pending', 3, 2, 0, NOW(), NULL, NULL, NULL, 680, 'medium', 0.0278, 'AI评估：申请人有一定债务，但在可控范围内', 'conv_20250530_002', '建议适当降低贷款金额至12万元', NOW(), NOW()),
('LA20250530003', 3, 3, 30000000, 0, 60, 0, '扩大合作社规模', DATE_ADD(NOW(), INTERVAL 30 DAY), '张合作社', '110101199003031236', '13900000003', '13900000003', 'coop001@example.com', 1500000, 18000000, '合作社经营收入', 500000, 200.0, '["wheat", "corn", "vegetables"]', 8, 'LC20231103001', '["application.pdf", "business_plan.pdf", "financial_report.pdf"]', '{"members_count": 120, "annual_revenue": 2000000}', 'rejected', 3, 3, 0, NOW(), 0, 0, 0.0000, 620, 'high', 0.0278, 'AI评估：合作社成立时间较短，财务状况不够稳定', 'conv_20250530_003', '建议申请人补充更多财务证明材料后重新申请', NOW(), NOW()),
('LA20250530004', 1, 5, 2000000, 2000000, 6, 6, '春耕生产资金', DATE_ADD(NOW(), INTERVAL 3 DAY), '李大牛', '110101199001011234', '13900000001', '13900000001', 'farmer001@example.com', 800000, 9600000, '农业种植收入', 0, 50.5, '["wheat"]', 15, 'LC20231101001', '["application.pdf", "planting_plan.pdf"]', '{"planting_area": 50.5, "expected_yield": 300}', 'disbursed', 2, 3, 1, NOW(), 2000000, 6, 0.0600, 750, 'low', 0.0000, 'AI评估：季节性贷款，风险可控', 'conv_20250530_004', '建议快速批准，支持春耕生产', NOW(), NOW());

-- ========== 9. 审批日志数据 ==========
INSERT INTO approval_logs (application_id, approver_id, action, step, approval_level, result, status, comment, note, previous_status, new_status, action_time, approved_at, created_at) VALUES
(1, 2, 'approve', 'initial_review', 1, 'approved', 'completed', '初审通过，申请人资质良好', '符合产品准入标准', 'pending', 'under_review', NOW(), NOW(), NOW()),
(1, 3, 'approve', 'risk_assessment', 2, 'approved', 'completed', '风险评估通过，风险等级为低', 'AI评估结果与人工复核一致', 'under_review', 'risk_approved', NOW(), NOW(), NOW()),
(1, 2, 'approve', 'final_approval', 3, 'approved', 'completed', '最终审批通过，准予放款', '可按申请条件执行', 'risk_approved', 'approved', NOW(), NOW(), NOW()),
(2, 2, 'approve', 'initial_review', 1, 'approved', 'completed', '初审通过，转风控评估', '建议关注债务情况', 'pending', 'under_review', NOW(), NOW(), NOW()),
(2, 3, 'pending', 'risk_assessment', 2, 'pending', 'in_progress', '', '等待进一步评估', 'under_review', 'risk_reviewing', NOW(), NULL, NOW()),
(3, 2, 'approve', 'initial_review', 1, 'approved', 'completed', '初审通过，转风控评估', '合作社资质需要重点关注', 'pending', 'under_review', NOW(), NOW(), NOW()),
(3, 3, 'reject', 'risk_assessment', 2, 'rejected', 'completed', '风险评估未通过', '合作社财务状况不够稳定，建议补充材料', 'under_review', 'rejected', NOW(), NOW(), NOW()),
(4, 2, 'approve', 'initial_review', 1, 'approved', 'completed', '初审通过', '季节性贷款，快速审批', 'pending', 'under_review', NOW(), NOW(), NOW()),
(4, 3, 'approve', 'risk_assessment', 2, 'approved', 'completed', '风险评估通过', 'AI评估结果良好', 'under_review', 'risk_approved', NOW(), NOW(), NOW()),
(4, 2, 'approve', 'final_approval', 3, 'approved', 'completed', '最终审批通过，已放款', '支持春耕生产', 'risk_approved', 'disbursed', NOW(), NOW(), NOW());

-- ========== 10. Dify工作流日志数据 ==========
INSERT INTO dify_workflow_logs (application_id, workflow_id, conversation_id, message_id, workflow_type, request_data, response_data, status, result, recommendation, confidence_score, start_time, end_time, duration, token_usage, cost_amount, created_at, updated_at) VALUES
(1, 'wf_loan_small_001', 'conv_20250530_001', 'msg_001', 'loan_approval', '{"user_type": "farmer", "amount": 50000, "income": 96000, "credit_score": 750}', '{"approval_result": "approve", "recommended_amount": 50000, "risk_level": "low", "confidence": 0.92}', 'completed', 'approve', '建议批准，申请人信用记录良好，农业经营经验丰富', 0.92, NOW(), NOW(), 1250, 85, 12, NOW(), NOW()),
(2, 'wf_loan_machine_001', 'conv_20250530_002', 'msg_002', 'loan_approval', '{"user_type": "farmer", "amount": 150000, "income": 72000, "debt_ratio": 0.0278}', '{"approval_result": "approve_with_condition", "recommended_amount": 120000, "risk_level": "medium", "confidence": 0.78}', 'completed', 'approve_with_condition', '建议适当降低贷款金额，控制风险', 0.78, NOW(), NOW(), 1850, 120, 18, NOW(), NOW()),
(3, 'wf_loan_coop_001', 'conv_20250530_003', 'msg_003', 'loan_approval', '{"user_type": "cooperative", "amount": 300000, "establishment_years": 2, "member_count": 120}', '{"approval_result": "reject", "recommended_amount": 0, "risk_level": "high", "confidence": 0.85}', 'completed', 'reject', '合作社成立时间较短，建议补充更多财务证明', 0.85, NOW(), NOW(), 2100, 145, 21, NOW(), NOW()),
(4, 'wf_loan_seasonal_001', 'conv_20250530_004', 'msg_004', 'loan_approval', '{"user_type": "farmer", "amount": 20000, "season": "spring", "crop_type": "wheat"}', '{"approval_result": "approve", "recommended_amount": 20000, "risk_level": "low", "confidence": 0.95}', 'completed', 'approve', '季节性贷款，支持春耕生产，风险可控', 0.95, NOW(), NOW(), 980, 65, 9, NOW(), NOW());

-- 重新启用外键检查
SET FOREIGN_KEY_CHECKS = 1;

-- 显示插入结果统计
SELECT 'system_configs' as table_name, COUNT(*) as record_count FROM system_configs
UNION ALL
SELECT 'oa_roles', COUNT(*) FROM oa_roles
UNION ALL
SELECT 'oa_users', COUNT(*) FROM oa_users
UNION ALL
SELECT 'users', COUNT(*) FROM users
UNION ALL
SELECT 'user_auths', COUNT(*) FROM user_auths
UNION ALL
SELECT 'user_tags', COUNT(*) FROM user_tags
UNION ALL
SELECT 'loan_products', COUNT(*) FROM loan_products
UNION ALL
SELECT 'loan_applications', COUNT(*) FROM loan_applications
UNION ALL
SELECT 'approval_logs', COUNT(*) FROM approval_logs
UNION ALL
SELECT 'dify_workflow_logs', COUNT(*) FROM dify_workflow_logs;

-- 显示成功消息
SELECT '✅ 模拟数据插入完成！' as message, 
       '请检查各表数据是否正确插入' as note,
       NOW() as completed_at; 