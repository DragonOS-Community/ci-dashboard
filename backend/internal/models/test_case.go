package models

import (
	"time"

	"gorm.io/gorm"
)

// TestCaseStatus 测例状态
type TestCaseStatus string

const (
	TestCaseStatusPassed  TestCaseStatus = "passed"
	TestCaseStatusFailed  TestCaseStatus = "failed"
	TestCaseStatusSkipped TestCaseStatus = "skipped"
)

// TestCase 测例详情模型
type TestCase struct {
	ID         uint64         `gorm:"primaryKey;autoIncrement" json:"id"`
	TestRunID  uint64         `gorm:"type:bigint unsigned;not null;index" json:"test_run_id"`
	Name       string         `gorm:"type:varchar(500);not null;index" json:"name"`
	Status     TestCaseStatus `gorm:"type:enum('passed','failed','skipped');not null;index" json:"status"`
	DurationMs uint32         `gorm:"type:int unsigned;default:0" json:"duration_ms"`
	ErrorLog   string         `gorm:"type:text" json:"error_log,omitempty"`
	DebugLog   string         `gorm:"type:text" json:"debug_log,omitempty"`
	CreatedAt  time.Time      `gorm:"type:datetime;not null;default:CURRENT_TIMESTAMP" json:"created_at"`

	// 关联关系
	TestRun TestRun `gorm:"foreignKey:TestRunID" json:"test_run,omitempty"`
}

// TableName 指定表名
func (TestCase) TableName() string {
	return "test_cases"
}

// BeforeCreate 创建前钩子
func (tc *TestCase) BeforeCreate(tx *gorm.DB) error {
	tc.CreatedAt = time.Now()
	return nil
}

// IsFailed 检查测例是否失败
func (tc *TestCase) IsFailed() bool {
	return tc.Status == TestCaseStatusFailed
}
