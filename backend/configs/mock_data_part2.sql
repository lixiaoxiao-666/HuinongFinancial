/*
 数字惠农后端 - 模拟数据插入脚本 (第二部分)
 
 包含：农机设备、专家、文章、订单、会话、文件等数据
 
 使用说明：
 1. 请先执行 mock_data.sql 基础数据
 2. 然后执行此脚本添加业务数据
 
 执行方法：
 mysql -u username -p database_name < mock_data_part2.sql
 
 Date: 30/05/2025
*/

USE huinong_db;

-- 禁用外键检查
SET FOREIGN_KEY_CHECKS = 0;

-- ========== 11. 文章分类数据 ==========
INSERT INTO categories (name, display_name, description, parent_id, sort_order, icon, status, created_at, updated_at) VALUES
('policy', '政策资讯', '农业政策、法规、补贴等相关资讯', NULL, 1, 'policy-icon', 'active', NOW(), NOW()),
('technology', '技术指导', '农业种植、养殖技术指导文章', NULL, 2, 'tech-icon', 'active', NOW(), NOW()),
('market', '市场信息', '农产品价格、市场行情分析', NULL, 3, 'market-icon', 'active', NOW(), NOW()),
('finance', '金融服务', '农业金融、贷款、保险等服务信息', NULL, 4, 'finance-icon', 'active', NOW(), NOW()),
('machine', '农机设备', '农业机械设备介绍、使用技巧', NULL, 5, 'machine-icon', 'active', NOW(), NOW()),
('news', '行业新闻', '农业行业最新动态和新闻', NULL, 6, 'news-icon', 'active', NOW(), NOW());

-- ========== 12. 专家数据 ==========
INSERT INTO experts (name, title, organization, specialties, phone, email, we_chat, avatar, biography, experience_years, service_areas, rating, rating_count, status, is_verified, verified_at, created_at, updated_at) VALUES
('张农业', '高级农艺师', '山东省农业科学院', '["小麦种植", "玉米栽培", "土壤改良"]', '13700000001', 'zhang.agriculture@shandong.gov.cn', 'zhangny2024', '/uploads/expert_avatar_001.jpg', '从事农业技术推广工作20余年，在粮食作物种植方面具有丰富经验', 22, '["山东省", "河南省", "河北省"]', 4.8, 156, 'active', 1, NOW(), NOW(), NOW()),
('李植保', '植物保护专家', '中国农业大学', '["病虫害防治", "农药使用", "生物防控"]', '13700000002', 'li.protection@cau.edu.cn', 'lizhibo2024', '/uploads/expert_avatar_002.jpg', '植物保护专业博士，专注农作物病虫害绿色防控技术研究', 15, '["全国"]', 4.9, 203, 'active', 1, NOW(), NOW(), NOW()),
('王机械', '农机专家', '雷沃重工股份有限公司', '["农机维修", "设备选型", "操作培训"]', '13700000003', 'wang.machine@lovol.com', 'wangjixie2024', '/uploads/expert_avatar_003.jpg', '农业机械工程师，在大型农机设备研发和维护方面有20年经验', 20, '["华北地区", "华东地区"]', 4.7, 89, 'active', 1, NOW(), NOW(), NOW()),
('陈畜牧', '畜牧养殖专家', '河南牧原食品股份有限公司', '["生猪养殖", "饲料配方", "疫病防控"]', '13700000004', 'chen.livestock@muyuan.com', 'chenxumu2024', '/uploads/expert_avatar_004.jpg', '畜牧兽医专业，专注规模化养殖技术和动物疫病防控', 18, '["河南省", "湖北省", "安徽省"]', 4.6, 124, 'active', 1, NOW(), NOW(), NOW()),
('赵经济', '农业经济师', '农业农村部农村经济研究中心', '["合作社管理", "农业经营", "政策解读"]', '13700000005', 'zhao.economy@moa.gov.cn', 'zhaojingji2024', '/uploads/expert_avatar_005.jpg', '农业经济管理博士，长期从事农村经济政策研究和合作社指导工作', 25, '["全国"]', 4.9, 267, 'active', 1, NOW(), NOW(), NOW()),
('孙园艺', '园艺技术专家', '寿光市农业农村局', '["蔬菜种植", "大棚技术", "无土栽培"]', '13700000006', 'sun.horticulture@shouguang.gov.cn', 'sunyuanyi2024', '/uploads/expert_avatar_006.jpg', '园艺学硕士，在设施蔬菜栽培和现代农业技术推广方面经验丰富', 12, '["山东省", "河北省"]', 4.8, 178, 'active', 1, NOW(), NOW(), NOW());

-- ========== 13. 文章数据 ==========
INSERT INTO articles (title, subtitle, content, summary, category, tags, cover_image, author_id, author_name, status, is_top, is_featured, view_count, like_count, share_count, comment_count, seo_title, seo_description, seo_keywords, published_at, created_at, updated_at) VALUES
('2025年农业补贴政策最新解读', '种粮农民直接补贴、农机购置补贴等政策详解', '
<h2>2025年农业补贴政策概览</h2>
<p>根据农业农村部最新发布的文件，2025年将继续实施多项惠农政策...</p>
<h3>种粮农民直接补贴</h3>
<p>补贴标准：每亩补贴120-150元，具体金额由各省区市根据实际情况确定...</p>
<h3>农机购置补贴</h3>
<p>补贴比例：不超过同档产品上年平均销售价格的30%...</p>
<h3>新型农业经营主体支持政策</h3>
<p>对农民专业合作社、家庭农场等新型经营主体给予重点支持...</p>
', '详细解读2025年各项农业补贴政策，包括种粮直补、农机补贴、新型经营主体支持等', 'policy', '["农业补贴", "政策解读", "2025年", "惠农政策"]', '/uploads/article_cover_001.jpg', 5, '陈内容', 'published', 1, 1, 1256, 89, 23, 12, '2025年农业补贴政策最新解读-数字惠农', '了解2025年最新农业补贴政策，种粮直补、农机购置补贴政策详解', '农业补贴,政策解读,种粮补贴,农机补贴', NOW(), NOW(), NOW()),

('小麦春季田间管理技术要点', '抓住关键期，确保小麦丰产丰收', '
<h2>春季小麦管理关键期</h2>
<p>立春后，小麦进入返青期，这是决定当年产量的关键时期...</p>
<h3>返青期管理</h3>
<p>1. 及时浇水：根据土壤墒情和苗情，适时浇好返青水...</p>
<p>2. 追肥管理：结合浇水进行追肥，每亩施尿素15-20公斤...</p>
<h3>病虫害防治</h3>
<p>春季是小麦条纹花叶病、纹枯病等病害的高发期...</p>
<h3>化学除草</h3>
<p>选择晴朗无风天气，在小麦4叶1心后进行化学除草...</p>
', '春季小麦田间管理技术指导，包括返青期管理、病虫害防治、化学除草等要点', 'technology', '["小麦种植", "春季管理", "田间管理", "农业技术"]', '/uploads/article_cover_002.jpg', 1, '张农业', 'published', 0, 1, 892, 156, 45, 28, '小麦春季田间管理技术要点-数字惠农', '小麦春季管理技术指导，返青期管理、病虫害防治等关键技术要点', '小麦种植,春季管理,田间管理,农业技术', NOW(), NOW(), NOW()),

('玉米价格行情分析及后市预测', '2025年玉米市场走势分析', '
<h2>当前玉米价格形势</h2>
<p>据最新市场监测数据，全国玉米平均收购价格为2.85元/公斤...</p>
<h3>价格影响因素分析</h3>
<p>1. 供给侧：2024年玉米产量较上年增加3.2%...</p>
<p>2. 需求侧：饲料需求稳中有升，深加工需求增长...</p>
<p>3. 政策因素：国家收储政策调整...</p>
<h3>后市预测</h3>
<p>综合分析各项因素，预计后期玉米价格将呈现稳中有升态势...</p>
', '分析当前玉米价格形势，解读影响因素，预测后市走势', 'market', '["玉米价格", "市场分析", "价格预测", "农产品"]', '/uploads/article_cover_003.jpg', 5, '陈内容', 'published', 0, 0, 567, 34, 18, 7, '玉米价格行情分析及后市预测-数字惠农', '2025年玉米市场价格走势分析，影响因素及后市预测', '玉米价格,市场行情,价格预测,农产品市场', NOW(), NOW(), NOW()),

('农机购置贷款申请指南', '一文了解农机设备贷款申请流程', '
<h2>农机购置贷款产品介绍</h2>
<p>数字惠农平台为广大农户提供专业的农机设备购置贷款服务...</p>
<h3>贷款条件</h3>
<p>1. 年龄18-65周岁，具有完全民事行为能力...</p>
<p>2. 从事农业生产经营满2年以上...</p>
<p>3. 信用记录良好，无不良征信记录...</p>
<h3>申请流程</h3>
<p>第一步：在线提交申请...</p>
<p>第二步：提供相关材料...</p>
<p>第三步：风险评估...</p>
<p>第四步：审批放款...</p>
<h3>所需材料</h3>
<p>身份证明、收入证明、农机购置合同等...</p>
', '详细介绍农机购置贷款申请条件、流程、所需材料等信息', 'finance', '["农机贷款", "申请指南", "贷款流程", "金融服务"]', '/uploads/article_cover_004.jpg', 5, '陈内容', 'published', 0, 1, 423, 67, 12, 5, '农机购置贷款申请指南-数字惠农', '农机设备贷款申请条件、流程、材料要求详细指南', '农机贷款,贷款申请,申请指南,农业金融', NOW(), NOW(), NOW()),

('现代农业机械化发展趋势', '智能化、绿色化成为主流方向', '
<h2>我国农业机械化现状</h2>
<p>截至2024年，全国农作物耕种收综合机械化率达到73%...</p>
<h3>发展趋势分析</h3>
<p>1. 智能化水平不断提升：GPS导航、自动驾驶技术...</p>
<p>2. 绿色化发展：节能减排、精准施药...</p>
<p>3. 专业化服务：农机合作社、农机服务公司...</p>
<h3>技术创新亮点</h3>
<p>无人机植保、智能收割机、精准播种机等新技术...</p>
', '分析我国农业机械化发展现状，探讨未来发展趋势和技术创新方向', 'machine', '["农业机械化", "智能农机", "发展趋势", "农机技术"]', '/uploads/article_cover_005.jpg', 3, '王农机', 'published', 0, 0, 334, 28, 8, 3, '现代农业机械化发展趋势-数字惠农', '分析我国农业机械化发展现状及未来智能化、绿色化发展趋势', '农业机械化,智能农机,发展趋势,农机技术', NOW(), NOW(), NOW());

-- ========== 14. 农机设备数据 ==========
INSERT INTO machines (machine_code, machine_name, machine_type, brand, model, specifications, description, images, owner_id, owner_type, province, city, county, address, longitude, latitude, hourly_rate, daily_rate, per_acre_rate, deposit_amount, status, available_schedule, rating, rating_count, min_rental_hours, max_advance_days, is_verified, verified_at, created_at, updated_at) VALUES
('HN_TRACTOR_001', '大型拖拉机', '拖拉机', '约翰迪尔', '6B-1404', '{"power": "140马力", "drive": "四驱", "transmission": "动力换挡", "pto": "540/1000rpm"}', '大型四驱拖拉机，适用于大田作业，配备GPS导航系统', '["machine_001_1.jpg", "machine_001_2.jpg", "machine_001_3.jpg"]', 1, 'user', '山东省', '济南市', '章丘区', '相公庄街道农机服务站', 117.123456, 36.654321, 15000, 120000, 8000, 50000, 'available', '{"available_dates": ["2025-06-01", "2025-06-02"], "busy_dates": ["2025-05-25"]}', 4.8, 23, 4, 15, 1, NOW(), NOW(), NOW()),

('HN_HARVESTER_001', '小麦收割机', '收割机', '雷沃重工', 'GM80', '{"cutting_width": "2.8米", "grain_tank": "4500升", "power": "125马力", "fuel_consumption": "18L/h"}', '高效小麦收割机，收割效率高，损失率低，适合中大型地块作业', '["machine_002_1.jpg", "machine_002_2.jpg"]', 2, 'user', '河南省', '郑州市', '中牟县', '韩寺镇农机合作社', 113.987654, 34.765432, 20000, 160000, 12000, 80000, 'available', '{"available_dates": ["2025-06-10", "2025-06-15"], "busy_dates": []}', 4.9, 31, 6, 20, 1, NOW(), NOW(), NOW()),

('HN_PLANTER_001', '玉米播种机', '播种机', '大华宝来', '2BMFJ-6', '{"rows": "6行", "row_spacing": "60cm", "seeding_depth": "3-8cm", "fertilizer_box": "300L"}', '精密玉米播种机，播种精度高，可同时施肥，提高播种效率', '["machine_003_1.jpg", "machine_003_2.jpg"]', 3, 'user', '江苏省', '徐州市', '沛县', '安国镇机械化服务中心', 116.456789, 34.567890, 12000, 96000, 6000, 40000, 'available', '{"available_dates": ["2025-04-15", "2025-04-20"], "busy_dates": []}', 4.7, 18, 4, 10, 1, NOW(), NOW(), NOW()),

('HN_SPRAYER_001', '自走式喷药机', '植保机械', '丰疆智能', '3WPZ-1000', '{"tank_capacity": "1000升", "boom_width": "18米", "spray_pressure": "0.2-0.6MPa", "speed": "8-15km/h"}', '高效自走式喷药机，喷洒均匀，适用于大面积农田植保作业', '["machine_004_1.jpg", "machine_004_2.jpg"]', 1, 'user', '山东省', '济南市', '章丘区', '相公庄街道农机服务站', 117.234567, 36.765432, 18000, 144000, 10000, 60000, 'available', '{"available_dates": ["2025-05-01", "2025-05-05"], "busy_dates": []}', 4.6, 15, 2, 7, 1, NOW(), NOW(), NOW()),

('HN_ROTARY_001', '旋耕机', '耕整机械', '常发佳联', '1GQN-200', '{"working_width": "2.0米", "tillage_depth": "15-25cm", "rotor_type": "反转", "power_requirement": "60-80马力"}', '高效旋耕机，适用于水田和旱田耕整，作业质量好', '["machine_005_1.jpg", "machine_005_2.jpg"]', 4, 'user', '安徽省', '合肥市', '肥东县', '店埠镇农机专业合作社', 117.654321, 31.876543, 8000, 64000, 4000, 20000, 'available', '{"available_dates": ["2025-03-20", "2025-03-25"], "busy_dates": []}', 4.5, 12, 2, 5, 1, NOW(), NOW(), NOW()),

('HN_DRONE_001', '植保无人机', '无人机', '大疆农业', 'T40', '{"load_capacity": "40升", "flight_time": "25分钟", "spray_width": "7米", "positioning": "RTK"}', 'RTK高精度植保无人机，适用于精准施药和植保作业', '["machine_006_1.jpg", "machine_006_2.jpg"]', 5, 'user', '湖北省', '襄阳市', '老河口市', '竹林桥镇农业服务站', 111.567890, 32.123456, 25000, 200000, 15000, 100000, 'available', '{"available_dates": ["2025-04-01", "2025-04-10"], "busy_dates": []}', 4.9, 8, 1, 3, 1, NOW(), NOW(), NOW());

-- ========== 15. 租赁订单数据 ==========
INSERT INTO rental_orders (order_no, machine_id, renter_id, owner_id, start_time, end_time, rental_duration, rental_location, contact_person, contact_phone, billing_method, unit_price, quantity, subtotal_amount, deposit_amount, total_amount, status, payment_method, payment_status, paid_amount, paid_at, owner_confirmed_at, renter_confirmed_at, renter_rating, renter_comment, owner_rating, owner_comment, remarks, created_at, updated_at) VALUES
('RO20250530001', 1, 2, 1, '2025-06-01 08:00:00', '2025-06-01 18:00:00', 10, '河南省郑州市中牟县韩寺镇王家村田地', '王二麻子', '13900000002', 'hourly', 15000, 10, 150000, 50000, 200000, 'completed', 'alipay', 'paid', 200000, NOW(), NOW(), NOW(), 4.8, '设备性能良好，操作方便，服务态度好', 4.9, '客户守时诚信，田地条件好，合作愉快', '春耕作业，效果很好', NOW(), NOW()),

('RO20250530002', 2, 3, 2, '2025-06-10 06:00:00', '2025-06-12 18:00:00', 60, '江苏省徐州市沛县安国镇张楼村小麦田', '张合作社', '13900000003', 'hourly', 20000, 60, 1200000, 80000, 1280000, 'confirmed', 'bank_transfer', 'paid', 1280000, NOW(), NOW(), NULL, NULL, NULL, NULL, NULL, '小麦收割作业，面积约200亩', NOW(), NOW()),

('RO20250530003', 3, 1, 3, '2025-04-15 07:00:00', '2025-04-15 17:00:00', 10, '山东省济南市章丘区相公庄街道李家村玉米地', '李大牛', '13900000001', 'hourly', 12000, 10, 120000, 40000, 160000, 'completed', 'wechat_pay', 'paid', 160000, NOW(), NOW(), NOW(), 4.7, '播种效果不错，深浅一致', 4.8, '合作愉快，期待下次合作', '玉米播种作业', NOW(), NOW()),

('RO20250530004', 4, 4, 1, '2025-05-01 09:00:00', '2025-05-01 15:00:00', 6, '安徽省合肥市肥东县店埠镇刘集村农田', '刘企业', '13900000004', 'hourly', 18000, 6, 108000, 60000, 168000, 'in_progress', 'alipay', 'paid', 168000, NOW(), NOW(), NULL, NULL, NULL, NULL, NULL, '农田植保作业', NOW(), NOW()),

('RO20250530005', 5, 2, 4, '2025-03-20 08:00:00', '2025-03-22 17:00:00', 54, '河南省郑州市中牟县韩寺镇王家村田地', '王二麻子', '13900000002', 'hourly', 8000, 54, 432000, 20000, 452000, 'completed', 'bank_transfer', 'paid', 452000, NOW(), NOW(), NOW(), 4.6, '旋耕效果满意，土地平整', 4.7, '客户配合度高，作业顺利', '春耕土地准备', NOW(), NOW());

-- ========== 16. 用户会话数据 ==========
INSERT INTO user_sessions (user_id, session_id, platform, device_id, device_type, app_version, ip_address, location, access_token, refresh_token, token_expires_at, status, login_time, logout_time, created_at, updated_at) VALUES
(1, 'sess_farmer001_20250530_001', 'mobile', 'device_android_001', 'android', '1.2.0', '192.168.1.201', '山东省济南市章丘区', 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...', 'refresh_token_001', DATE_ADD(NOW(), INTERVAL 7 DAY), 'active', NOW(), NULL, NOW(), NOW()),
(2, 'sess_farmer002_20250530_001', 'mobile', 'device_ios_001', 'ios', '1.2.0', '192.168.1.202', '河南省郑州市中牟县', 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...', 'refresh_token_002', DATE_ADD(NOW(), INTERVAL 7 DAY), 'active', NOW(), NULL, NOW(), NOW()),
(3, 'sess_coop001_20250530_001', 'web', 'browser_chrome_001', 'pc', 'web_1.0', '192.168.1.203', '江苏省徐州市沛县', 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...', 'refresh_token_003', DATE_ADD(NOW(), INTERVAL 7 DAY), 'active', NOW(), NULL, NOW(), NOW()),
(1, 'sess_farmer001_20250529_001', 'mobile', 'device_android_001', 'android', '1.2.0', '192.168.1.201', '山东省济南市章丘区', 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...', 'refresh_token_004', DATE_SUB(NOW(), INTERVAL 1 DAY), 'expired', DATE_SUB(NOW(), INTERVAL 1 DAY), DATE_SUB(NOW(), INTERVAL 2 HOUR), NOW(), NOW()),
(4, 'sess_company001_20250530_001', 'web', 'browser_firefox_001', 'pc', 'web_1.0', '192.168.1.204', '安徽省合肥市肥东县', 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...', 'refresh_token_005', DATE_ADD(NOW(), INTERVAL 7 DAY), 'active', NOW(), NULL, NOW(), NOW());

-- ========== 17. 文件上传数据 ==========
INSERT INTO file_uploads (file_name, original_name, file_path, file_url, file_size, file_type, mime_type, file_hash, uploader_id, uploader_type, business_type, business_id, storage_type, bucket_name, status, is_public, access_count, created_at, updated_at) VALUES
('id_front_001.jpg', '身份证正面.jpg', '/uploads/2025/05/30/id_front_001.jpg', 'https://cdn.huinong.com/uploads/2025/05/30/id_front_001.jpg', 524288, 'image', 'image/jpeg', '5d41402abc4b2a76b9719d911017c592', 1, 'user', 'user_auth', 1, 'local', NULL, 'uploaded', 0, 5, NOW(), NOW()),
('id_back_001.jpg', '身份证背面.jpg', '/uploads/2025/05/30/id_back_001.jpg', 'https://cdn.huinong.com/uploads/2025/05/30/id_back_001.jpg', 487424, 'image', 'image/jpeg', '098f6bcd4621d373cade4e832627b4f6', 1, 'user', 'user_auth', 1, 'local', NULL, 'uploaded', 0, 3, NOW(), NOW()),
('application_001.pdf', '贷款申请书.pdf', '/uploads/2025/05/30/application_001.pdf', 'https://cdn.huinong.com/uploads/2025/05/30/application_001.pdf', 2097152, 'document', 'application/pdf', '5d41402abc4b2a76b9719d911017c593', 1, 'user', 'loan_application', 1, 'local', NULL, 'uploaded', 0, 2, NOW(), NOW()),
('income_proof_001.pdf', '收入证明.pdf', '/uploads/2025/05/30/income_proof_001.pdf', 'https://cdn.huinong.com/uploads/2025/05/30/income_proof_001.pdf', 1572864, 'document', 'application/pdf', '098f6bcd4621d373cade4e832627b4f7', 1, 'user', 'loan_application', 1, 'local', NULL, 'uploaded', 0, 1, NOW(), NOW()),
('machine_001_1.jpg', '拖拉机主图.jpg', '/uploads/2025/05/30/machine_001_1.jpg', 'https://cdn.huinong.com/uploads/2025/05/30/machine_001_1.jpg', 786432, 'image', 'image/jpeg', '5d41402abc4b2a76b9719d911017c594', 1, 'user', 'machine', 1, 'local', NULL, 'uploaded', 1, 45, NOW(), NOW()),
('article_cover_001.jpg', '政策资讯配图.jpg', '/uploads/2025/05/30/article_cover_001.jpg', 'https://cdn.huinong.com/uploads/2025/05/30/article_cover_001.jpg', 612352, 'image', 'image/jpeg', '098f6bcd4621d373cade4e832627b4f8', 5, 'oa_user', 'article', 1, 'local', NULL, 'uploaded', 1, 156, NOW(), NOW()),
('expert_avatar_001.jpg', '专家头像.jpg', '/uploads/2025/05/30/expert_avatar_001.jpg', 'https://cdn.huinong.com/uploads/2025/05/30/expert_avatar_001.jpg', 245760, 'image', 'image/jpeg', '5d41402abc4b2a76b9719d911017c595', 5, 'oa_user', 'expert', 1, 'local', NULL, 'uploaded', 1, 78, NOW(), NOW());

-- ========== 18. 离线队列数据 ==========
INSERT INTO offline_queue (user_id, action_type, request_data, status, retry_count, max_retries, next_retry, error_message, result_data, created_at, updated_at) VALUES
(1, 'sync_user_info', '{"user_id": 1, "update_fields": ["last_login_time", "login_count"]}', 'completed', 0, 3, NULL, NULL, '{"success": true, "updated_at": "2025-05-30T10:30:00Z"}', NOW(), NOW()),
(2, 'send_sms', '{"phone": "13900000002", "template": "loan_approval", "params": {"name": "王二麻子", "amount": "12万元"}}', 'pending', 1, 3, DATE_ADD(NOW(), INTERVAL 5 MINUTE), '短信服务暂时不可用', NULL, NOW(), NOW()),
(3, 'upload_file', '{"file_path": "/uploads/temp/business_plan.pdf", "target_path": "/uploads/2025/05/30/", "business_type": "loan_application"}', 'completed', 0, 3, NULL, NULL, '{"success": true, "file_url": "https://cdn.huinong.com/uploads/2025/05/30/business_plan.pdf"}', NOW(), NOW()),
(1, 'sync_credit_score', '{"user_id": 1, "credit_data": {"score": 750, "level": "excellent"}}', 'failed', 3, 3, NULL, '征信接口连接超时', NULL, NOW(), NOW()),
(4, 'send_email', '{"email": "company001@example.com", "template": "machine_rental_confirm", "params": {"order_no": "RO20250530004"}}', 'completed', 0, 3, NULL, NULL, '{"success": true, "message_id": "email_20250530_001"}', NOW(), NOW());

-- ========== 19. API日志数据 ==========
INSERT INTO api_logs (request_id, method, url, headers, query, body, ip_address, user_agent, status_code, response_body, response_time, user_id, user_type, error_code, error_message, created_at) VALUES
('req_20250530_001', 'POST', '/api/v1/auth/login', '{"Content-Type": "application/json", "User-Agent": "HuinongApp/1.2.0"}', '{}', '{"phone": "13900000001", "password": "password123"}', '192.168.1.201', 'HuinongApp/1.2.0 (Android 12)', 200, '{"code": 200, "message": "登录成功", "data": {"token": "eyJ..."}}', 245, 1, 'user', NULL, NULL, NOW()),
('req_20250530_002', 'GET', '/api/v1/loans/products', '{"Authorization": "Bearer eyJ...", "Content-Type": "application/json"}', '{"page": 1, "limit": 10}', NULL, '192.168.1.201', 'HuinongApp/1.2.0 (Android 12)', 200, '{"code": 200, "data": [...]}', 156, 1, 'user', NULL, NULL, NOW()),
('req_20250530_003', 'POST', '/api/v1/loans/applications', '{"Authorization": "Bearer eyJ...", "Content-Type": "application/json"}', '{}', '{"product_id": 1, "amount": 50000, "purpose": "购买种子化肥"}', '192.168.1.201', 'HuinongApp/1.2.0 (Android 12)', 201, '{"code": 201, "message": "申请提交成功"}', 890, 1, 'user', NULL, NULL, NOW()),
('req_20250530_004', 'GET', '/api/v1/machines', '{"Authorization": "Bearer eyJ...", "Content-Type": "application/json"}', '{"lat": 36.654321, "lng": 117.123456, "radius": 50}', NULL, '192.168.1.202', 'HuinongApp/1.2.0 (iOS 16)', 200, '{"code": 200, "data": [...]}', 234, 2, 'user', NULL, NULL, NOW()),
('req_20250530_005', 'POST', '/api/v1/orders/rental', '{"Authorization": "Bearer eyJ...", "Content-Type": "application/json"}', '{}', '{"machine_id": 1, "start_time": "2025-06-01T08:00:00", "end_time": "2025-06-01T18:00:00"}', '192.168.1.202', 'HuinongApp/1.2.0 (iOS 16)', 201, '{"code": 201, "message": "订单创建成功"}', 567, 2, 'user', NULL, NULL, NOW()),
('req_20250530_006', 'GET', '/api/v1/articles', '{"Content-Type": "application/json"}', '{"category": "policy", "page": 1}', NULL, '192.168.1.203', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) Chrome/120.0.0.0', 200, '{"code": 200, "data": [...]}', 123, NULL, NULL, NULL, NULL, NOW()),
('req_20250530_007', 'POST', '/api/v1/auth/refresh', '{"Authorization": "Bearer refresh_token", "Content-Type": "application/json"}', '{}', '{"refresh_token": "refresh_token_001"}', '192.168.1.201', 'HuinongApp/1.2.0 (Android 12)', 200, '{"code": 200, "data": {"access_token": "new_token"}}', 189, 1, 'user', NULL, NULL, NOW()),
('req_20250530_008', 'GET', '/api/v1/users/profile', '{"Authorization": "Bearer eyJ...", "Content-Type": "application/json"}', '{}', NULL, '192.168.1.204', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) Firefox/120.0', 200, '{"code": 200, "data": {...}}', 167, 4, 'user', NULL, NULL, NOW());

-- 重新启用外键检查
SET FOREIGN_KEY_CHECKS = 1;

-- 显示第二部分插入结果统计
SELECT 'categories' as table_name, COUNT(*) as record_count FROM categories
UNION ALL
SELECT 'experts', COUNT(*) FROM experts
UNION ALL
SELECT 'articles', COUNT(*) FROM articles
UNION ALL
SELECT 'machines', COUNT(*) FROM machines
UNION ALL
SELECT 'rental_orders', COUNT(*) FROM rental_orders
UNION ALL
SELECT 'user_sessions', COUNT(*) FROM user_sessions
UNION ALL
SELECT 'file_uploads', COUNT(*) FROM file_uploads
UNION ALL
SELECT 'offline_queue', COUNT(*) FROM offline_queue
UNION ALL
SELECT 'api_logs', COUNT(*) FROM api_logs;

-- 显示成功消息
SELECT '✅ 第二部分模拟数据插入完成！' as message, 
       '所有业务数据已添加完毕' as note,
       NOW() as completed_at; 