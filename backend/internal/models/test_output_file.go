package models

import (
	"time"

	"gorm.io/gorm"
)

// TestOutputFile 测试输出文件模型
type TestOutputFile struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	TestRunID uint64    `gorm:"type:bigint unsigned;not null;index" json:"test_run_id"`
	Filename  string    `gorm:"type:varchar(500);not null;index" json:"filename"`
	FilePath  string    `gorm:"type:varchar(1000);not null" json:"file_path"`
	FileSize  uint64    `gorm:"type:bigint unsigned;default:0" json:"file_size"`
	MimeType  string    `gorm:"type:varchar(100)" json:"mime_type,omitempty"`
	CreatedAt time.Time `gorm:"type:datetime;not null;default:CURRENT_TIMESTAMP" json:"created_at"`

	// 关联关系
	TestRun TestRun `gorm:"foreignKey:TestRunID" json:"test_run,omitempty"`
}

// TableName 指定表名
func (TestOutputFile) TableName() string {
	return "test_output_files"
}

// BeforeCreate 创建前钩子
func (tof *TestOutputFile) BeforeCreate(tx *gorm.DB) error {
	tof.CreatedAt = time.Now()
	return nil
}
