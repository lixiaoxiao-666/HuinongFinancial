-- 修改user_sessions表中device_name字段长度
-- 从varchar(100)增加到varchar(500)以支持完整的设备信息存储

USE huinong_financial;

-- 修改device_name字段长度
ALTER TABLE user_sessions 
MODIFY COLUMN device_name VARCHAR(500);

-- 添加索引优化查询性能（可选）
CREATE INDEX idx_user_sessions_device_name ON user_sessions(device_name(100)) IF NOT EXISTS; 