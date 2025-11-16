package models

import (
	"context"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// DBWithContext 从 gin.Context 创建带 request_id 的 DB 实例
// 这样 GORM 的日志就能包含 request_id
func DBWithContext(c *gin.Context) *gorm.DB {
	if c == nil {
		return DB
	}

	// 从 gin.Context 获取 request_id
	var requestID string
	if id, exists := c.Get("request_id"); exists {
		if rid, ok := id.(string); ok {
			requestID = rid
		}
	}

	// 创建带 request_id 的 context
	ctx := context.WithValue(context.Background(), "request_id", requestID)

	// 返回带 context 的 DB 实例
	return DB.WithContext(ctx)
}
