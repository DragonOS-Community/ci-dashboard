package models

import (
	"time"

	"gorm.io/gorm"
)

// SystemConfig 系统配置模型
type SystemConfig struct {
	ID          uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	ConfigKey   string    `gorm:"type:varchar(255);not null;uniqueIndex" json:"config_key"`
	ConfigValue string    `gorm:"type:varchar(1000);not null" json:"config_value"`
	Description string    `gorm:"type:varchar(500)" json:"description"`
	CreatedAt   time.Time `gorm:"type:datetime;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt   time.Time `gorm:"type:datetime;not null;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP" json:"updated_at"`
}

// TableName 指定表名
func (SystemConfig) TableName() string {
	return "system_configs"
}

// BeforeCreate 创建前钩子
func (sc *SystemConfig) BeforeCreate(tx *gorm.DB) error {
	now := time.Now()
	sc.CreatedAt = now
	sc.UpdatedAt = now
	return nil
}

// BeforeUpdate 更新前钩子
func (sc *SystemConfig) BeforeUpdate(tx *gorm.DB) error {
	sc.UpdatedAt = time.Now()
	return nil
}
