-- 数字惠农后端 - 数据库连接测试脚本

USE app;

-- 检查数据库是否存在并可访问
SELECT 
    DATABASE() as current_database,
    NOW() as connection_time,
    '✅ 数据库连接成功！' as status;

-- 显示当前数据库中的所有表
SHOW TABLES;

-- 检查主要表是否存在
SELECT 
    TABLE_NAME as table_name,
    CASE 
        WHEN TABLE_NAME IS NOT NULL THEN '✅ 存在'
        ELSE '❌ 不存在'
    END as status
FROM INFORMATION_SCHEMA.TABLES 
WHERE TABLE_SCHEMA = 'app' 
    AND TABLE_NAME IN (
        'users', 'oa_users', 'loan_products', 'loan_applications',
        'machines', 'rental_orders', 'articles', 'experts', 'categories'
    )
ORDER BY TABLE_NAME; 