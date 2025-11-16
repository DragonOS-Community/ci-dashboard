-- 创建系统配置表
CREATE TABLE IF NOT EXISTS system_configs (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    config_key VARCHAR(255) NOT NULL UNIQUE COMMENT '配置键',
    config_value VARCHAR(1000) NOT NULL COMMENT '配置值',
    description VARCHAR(500) COMMENT '配置描述',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_config_key (config_key)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='系统配置表';

-- 插入默认配置：是否允许上传测试输出文件（默认不允许）
INSERT INTO system_configs (config_key, config_value, description) 
VALUES ('allow_upload_output_files', 'false', '是否允许上传测试输出文件') 
ON DUPLICATE KEY UPDATE config_value=config_value;

