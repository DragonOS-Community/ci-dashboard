package models

import (
	"time"

	"gorm.io/gorm"
)

// TestRunStatus 测试运行状态
type TestRunStatus string

const (
	TestRunStatusRunning   TestRunStatus = "running"
	TestRunStatusPassed    TestRunStatus = "passed"
	TestRunStatusFailed    TestRunStatus = "failed"
	TestRunStatusCancelled TestRunStatus = "cancelled"
)

// TestType 测试类型
type TestType string

const (
	// TestTypeGvisor gvisor测试类型
	TestTypeGvisor TestType = "gvisor"
)

// TestRun 测试运行记录模型
type TestRun struct {
	ID            uint64        `gorm:"primaryKey;autoIncrement" json:"id"`
	ProjectID     uint64        `gorm:"type:bigint unsigned;not null;index" json:"project_id"`
	BranchName    string        `gorm:"type:varchar(255);not null;index" json:"branch_name"`
	CommitID      string        `gorm:"type:varchar(40);not null;index" json:"commit_id"`
	CommitShortID string        `gorm:"type:varchar(10);not null;index" json:"commit_short_id"`
	TestType      string        `gorm:"type:varchar(50);not null;default:'gvisor';index" json:"test_type"`
	Status        TestRunStatus `gorm:"type:enum('passed','failed','running','cancelled');not null;default:'running';index" json:"status"`
	IsPublic      bool          `gorm:"type:boolean;not null;default:true;index" json:"is_public"`
	StartedAt     *time.Time    `gorm:"type:datetime" json:"started_at,omitempty"`
	CompletedAt   *time.Time    `gorm:"type:datetime" json:"completed_at,omitempty"`
	CreatedAt     time.Time     `gorm:"type:datetime;not null;default:CURRENT_TIMESTAMP;index" json:"created_at"`

	// 关联关系
	Project     Project          `gorm:"foreignKey:ProjectID" json:"project,omitempty"`
	TestCases   []TestCase       `gorm:"foreignKey:TestRunID" json:"test_cases,omitempty"`
	OutputFiles []TestOutputFile `gorm:"foreignKey:TestRunID" json:"output_files,omitempty"`
}

// TableName 指定表名
func (TestRun) TableName() string {
	return "test_runs"
}

// BeforeCreate 创建前钩子
func (tr *TestRun) BeforeCreate(tx *gorm.DB) error {
	now := time.Now()
	if tr.StartedAt == nil {
		tr.StartedAt = &now
	}
	tr.CreatedAt = now
	return nil
}

// IsCompleted 检查测试是否已完成
func (tr *TestRun) IsCompleted() bool {
	return tr.Status == TestRunStatusPassed || tr.Status == TestRunStatusFailed || tr.Status == TestRunStatusCancelled
}

// Complete 完成测试运行
func (tr *TestRun) Complete(status TestRunStatus) {
	tr.Status = status
	now := time.Now()
	tr.CompletedAt = &now
}
