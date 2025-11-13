-- 创建项目表
CREATE TABLE IF NOT EXISTS projects (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL COMMENT '项目名称',
    description TEXT COMMENT '项目描述',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_name (name)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='项目表';

-- 创建测试运行记录表
CREATE TABLE IF NOT EXISTS test_runs (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    project_id BIGINT UNSIGNED NOT NULL COMMENT '项目ID',
    branch_name VARCHAR(255) NOT NULL COMMENT '分支名',
    commit_id VARCHAR(40) NOT NULL COMMENT '完整Commit ID',
    commit_short_id VARCHAR(10) NOT NULL COMMENT '短Commit ID',
    test_type VARCHAR(50) NOT NULL DEFAULT 'gvisor' COMMENT '测试类型',
    status ENUM('passed', 'failed', 'running', 'cancelled') NOT NULL DEFAULT 'running' COMMENT '测试状态',
    started_at DATETIME COMMENT '开始时间',
    completed_at DATETIME COMMENT '完成时间',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_project_id (project_id),
    INDEX idx_branch_name (branch_name),
    INDEX idx_commit_id (commit_id),
    INDEX idx_commit_short_id (commit_short_id),
    INDEX idx_status (status),
    INDEX idx_created_at (created_at),
    INDEX idx_test_type (test_type),
    FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='测试运行记录表';

-- 创建测例详情表
CREATE TABLE IF NOT EXISTS test_cases (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    test_run_id BIGINT UNSIGNED NOT NULL COMMENT '测试运行ID',
    name VARCHAR(500) NOT NULL COMMENT '测例名称',
    status ENUM('passed', 'failed', 'skipped') NOT NULL COMMENT '测例状态',
    duration_ms INT UNSIGNED DEFAULT 0 COMMENT '执行时长（毫秒）',
    error_log TEXT COMMENT '错误日志',
    debug_log TEXT COMMENT '调试日志',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_test_run_id (test_run_id),
    INDEX idx_name (name),
    INDEX idx_status (status),
    FOREIGN KEY (test_run_id) REFERENCES test_runs(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='测例详情表';

-- 创建测试输出文件表
CREATE TABLE IF NOT EXISTS test_output_files (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    test_run_id BIGINT UNSIGNED NOT NULL COMMENT '测试运行ID',
    filename VARCHAR(500) NOT NULL COMMENT '文件名',
    file_path VARCHAR(1000) NOT NULL COMMENT '文件路径',
    file_size BIGINT UNSIGNED DEFAULT 0 COMMENT '文件大小（字节）',
    mime_type VARCHAR(100) COMMENT 'MIME类型',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_test_run_id (test_run_id),
    INDEX idx_filename (filename),
    FOREIGN KEY (test_run_id) REFERENCES test_runs(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='测试输出文件表';

-- 创建API密钥表
CREATE TABLE IF NOT EXISTS api_keys (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL COMMENT '密钥名称',
    key_hash VARCHAR(255) NOT NULL COMMENT '密钥哈希值',
    project_id BIGINT UNSIGNED COMMENT '关联项目ID（可选）',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    last_used_at DATETIME COMMENT '最后使用时间',
    expires_at DATETIME COMMENT '过期时间',
    INDEX idx_project_id (project_id),
    INDEX idx_key_hash (key_hash),
    FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='API密钥表';

-- 创建用户表（后台管理）
CREATE TABLE IF NOT EXISTS users (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(100) NOT NULL UNIQUE COMMENT '用户名',
    password_hash VARCHAR(255) NOT NULL COMMENT '密码哈希值',
    role VARCHAR(50) NOT NULL DEFAULT 'admin' COMMENT '角色',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_username (username)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户表';

-- 插入默认项目
INSERT INTO projects (id, name, description) VALUES (1, 'DragonOS', 'DragonOS操作系统项目') ON DUPLICATE KEY UPDATE name=name;

