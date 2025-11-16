-- 移除is_public字段
ALTER TABLE test_runs 
DROP INDEX idx_is_public,
DROP COLUMN is_public;

