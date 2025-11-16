-- 添加is_public字段到test_runs表
ALTER TABLE test_runs 
ADD COLUMN is_public BOOLEAN NOT NULL DEFAULT TRUE COMMENT '是否对外公开展示' AFTER status,
ADD INDEX idx_is_public (is_public);

