# 测试数据设置指南

## 概述

为了测试AI智能体接口，需要在数据库中插入一些测试数据。本文档提供了完整的测试数据SQL脚本。

## 测试数据SQL脚本

```sql
-- 1. 插入测试用户
INSERT INTO users (user_id, phone, real_name, id_card, gender, birth_date, address, created_at, updated_at) VALUES
('user_001', '13800138001', '张三', '110101199001011234', 'male', '1990-01-01', '北京市朝阳区测试街道123号', NOW(), NOW());

-- 2. 插入用户画像
INSERT INTO user_profiles (user_id, occupation, annual_income, work_years, education_level, marital_status, has_house, has_car, credit_level, created_at, updated_at) VALUES
('user_001', '软件工程师', 200000.00, 5, 'bachelor', 'single', true, false, 'good', NOW(), NOW());

-- 3. 插入贷款产品
INSERT INTO loan_products (product_id, product_name, product_type, min_amount, max_amount, interest_rate, term_months, requirements, status, created_at, updated_at) VALUES
('product_001', '个人信用贷', 'personal_credit', 10000.00, 500000.00, 8.5, 36, '年收入不低于10万，征信良好', 'active', NOW(), NOW());

-- 4. 插入测试申请
INSERT INTO loan_applications (
    application_id, user_id, product_id, loan_amount, loan_purpose, 
    term_months, annual_income, status, applied_at, created_at, updated_at
) VALUES (
    'test_app_001', 'user_001', 'product_001', 100000.00, '装修', 
    24, 200000.00, 'pending_review', NOW(), NOW(), NOW()
);

-- 5. 插入上传文件记录
INSERT INTO uploaded_files (file_id, user_id, file_name, file_type, file_size, storage_path, purpose, created_at, updated_at) VALUES
('file_001', 'user_001', '身份证正面.jpg', 'image/jpeg', 1024000, '/uploads/user_001/id_card_front.jpg', 'id_card', NOW(), NOW()),
('file_002', 'user_001', '收入证明.pdf', 'application/pdf', 2048000, '/uploads/user_001/income_proof.pdf', 'income_proof', NOW(), NOW()),
('file_003', 'user_001', '银行流水.pdf', 'application/pdf', 3072000, '/uploads/user_001/bank_statement.pdf', 'bank_statement', NOW(), NOW());

-- 6. 插入AI分析结果表结构（如果不存在）
CREATE TABLE IF NOT EXISTS ai_analysis_results (
    id INT AUTO_INCREMENT PRIMARY KEY,
    application_id VARCHAR(50) NOT NULL,
    analysis_type VARCHAR(50) NOT NULL DEFAULT 'risk_assessment',
    risk_score DECIMAL(5,2) NOT NULL,
    risk_level VARCHAR(20) NOT NULL,
    decision VARCHAR(50) NOT NULL,
    confidence_score DECIMAL(5,2),
    reasons TEXT,
    analysis_data JSON,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_application_id (application_id)
);

-- 7. 插入工作流执行记录表结构（如果不存在）
CREATE TABLE IF NOT EXISTS workflow_executions (
    id INT AUTO_INCREMENT PRIMARY KEY,
    execution_id VARCHAR(50) NOT NULL UNIQUE,
    application_id VARCHAR(50) NOT NULL,
    workflow_type VARCHAR(50) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'running',
    priority VARCHAR(20) DEFAULT 'normal',
    callback_url VARCHAR(255),
    started_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    completed_at TIMESTAMP NULL,
    estimated_completion TIMESTAMP NULL,
    execution_data JSON,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_application_id (application_id),
    INDEX idx_execution_id (execution_id)
);
```

## 快速执行测试数据插入

### 方法1：使用MySQL命令行

```bash
# 连接到数据库
mysql -h 10.10.10.10 -P 4000 -u root -p app

# 执行上述SQL脚本
source /path/to/test_data.sql
```

### 方法2：使用DBeaver等工具

1. 连接到数据库：`10.10.10.10:4000`
2. 选择数据库：`app`
3. 复制粘贴上述SQL脚本并执行

### 方法3：创建Go脚本自动插入

```go
// scripts/insert_test_data.go
package main

import (
    "database/sql"
    "fmt"
    "time"
    
    _ "github.com/go-sql-driver/mysql"
)

func main() {
    dsn := "root:@tcp(10.10.10.10:4000)/app?charset=utf8mb4&parseTime=true"
    db, err := sql.Open("mysql", dsn)
    if err != nil {
        panic(err)
    }
    defer db.Close()

    // 插入测试用户
    _, err = db.Exec(`
        INSERT IGNORE INTO users (user_id, phone, real_name, id_card, gender, birth_date, address, created_at, updated_at) VALUES
        ('user_001', '13800138001', '张三', '110101199001011234', 'male', '1990-01-01', '北京市朝阳区测试街道123号', NOW(), NOW())
    `)
    if err != nil {
        fmt.Printf("插入用户失败: %v\n", err)
    } else {
        fmt.Println("✅ 测试用户插入成功")
    }

    // 插入测试申请
    _, err = db.Exec(`
        INSERT IGNORE INTO loan_applications (
            application_id, user_id, product_id, loan_amount, loan_purpose, 
            term_months, annual_income, status, applied_at, created_at, updated_at
        ) VALUES (
            'test_app_001', 'user_001', 'product_001', 100000.00, '装修', 
            24, 200000.00, 'pending_review', NOW(), NOW(), NOW()
        )
    `)
    if err != nil {
        fmt.Printf("插入申请失败: %v\n", err)
    } else {
        fmt.Println("✅ 测试申请插入成功")
    }

    fmt.Println("测试数据插入完成！")
}
```

## 验证测试数据

插入数据后，验证是否成功：

```sql
-- 检查测试数据
SELECT * FROM users WHERE user_id = 'user_001';
SELECT * FROM loan_applications WHERE application_id = 'test_app_001';
SELECT * FROM loan_products WHERE product_id = 'product_001';
```

## 测试接口

数据插入成功后，重新测试接口：

```bash
# 测试获取申请信息
curl -X GET "http://localhost:8080/api/v1/ai-agent/applications/test_app_001/info" \
  -H "Authorization: AI-Agent-Token ai_agent_secure_token_2024_v1" \
  -H "Content-Type: application/json"
```

预期返回完整的申请信息，包括用户信息、产品信息、用户画像等。 