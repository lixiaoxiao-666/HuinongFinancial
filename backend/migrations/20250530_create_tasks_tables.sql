-- 创建任务表
CREATE TABLE IF NOT EXISTS `tasks` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT,
    `title` varchar(255) NOT NULL COMMENT '任务标题',
    `description` text COMMENT '任务描述',
    `type` varchar(50) NOT NULL COMMENT '任务类型',
    `priority` varchar(20) DEFAULT 'medium' COMMENT '任务优先级(low,medium,high,urgent)',
    `status` varchar(20) DEFAULT 'pending' COMMENT '任务状态(pending,processing,completed,cancelled)',
    `assigned_to` bigint unsigned NULL COMMENT '分配给的用户ID',
    `created_by` bigint unsigned NOT NULL COMMENT '创建人ID',
    `business_id` bigint unsigned NOT NULL COMMENT '关联业务ID',
    `business_type` varchar(50) NOT NULL COMMENT '业务类型',
    `data` json COMMENT '任务相关数据',
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `completed_at` timestamp NULL COMMENT '完成时间',
    PRIMARY KEY (`id`),
    KEY `idx_tasks_type` (`type`),
    KEY `idx_tasks_status` (`status`),
    KEY `idx_tasks_priority` (`priority`),
    KEY `idx_tasks_assigned_to` (`assigned_to`),
    KEY `idx_tasks_created_by` (`created_by`),
    KEY `idx_tasks_business` (`business_type`, `business_id`),
    KEY `idx_tasks_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='任务表';

-- 创建任务操作记录表
CREATE TABLE IF NOT EXISTS `task_actions` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT,
    `task_id` bigint unsigned NOT NULL COMMENT '任务ID',
    `action` varchar(50) NOT NULL COMMENT '操作类型',
    `comment` text COMMENT '操作备注',
    `operator_id` bigint unsigned NOT NULL COMMENT '操作人ID',
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    KEY `idx_task_actions_task_id` (`task_id`),
    KEY `idx_task_actions_operator_id` (`operator_id`),
    KEY `idx_task_actions_action` (`action`),
    KEY `idx_task_actions_created_at` (`created_at`),
    CONSTRAINT `fk_task_actions_task_id` FOREIGN KEY (`task_id`) REFERENCES `tasks` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='任务操作记录表'; 