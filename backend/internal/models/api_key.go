package models

import (
	"time"

	"gorm.io/gorm"
)

// APIKey API密钥模型
type APIKey struct {
	ID         uint64     `gorm:"primaryKey;autoIncrement" json:"id"`
	Name       string     `gorm:"type:varchar(255);not null" json:"name"`
	KeyHash    string     `gorm:"type:varchar(255);not null;index" json:"-"` // 不返回给客户端
	ProjectID  *uint64    `gorm:"type:bigint unsigned;index" json:"project_id,omitempty"`
	CreatedAt  time.Time  `gorm:"type:datetime;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	LastUsedAt *time.Time `gorm:"type:datetime" json:"last_used_at,omitempty"`
	ExpiresAt  *time.Time `gorm:"type:datetime" json:"expires_at,omitempty"`

	// 关联关系
	Project Project `gorm:"foreignKey:ProjectID" json:"project,omitempty"`
}

// TableName 指定表名
func (APIKey) TableName() string {
	return "api_keys"
}

// BeforeCreate 创建前钩子
func (ak *APIKey) BeforeCreate(tx *gorm.DB) error {
	ak.CreatedAt = time.Now()
	return nil
}

// IsExpired 检查密钥是否过期
func (ak *APIKey) IsExpired() bool {
	if ak.ExpiresAt == nil {
		return false
	}
	return time.Now().After(*ak.ExpiresAt)
}

// UpdateLastUsed 更新最后使用时间
func (ak *APIKey) UpdateLastUsed() {
	now := time.Now()
	ak.LastUsedAt = &now
}
