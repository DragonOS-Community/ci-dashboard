package services

import (
	"github.com/dragonos/dragonos-ci-dashboard/internal/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// getDB 获取数据库实例，使用 gin.Context 创建带 request_id 的 DB 实例
// 这样 GORM 的日志就能包含 request_id
func getDB(c *gin.Context) *gorm.DB {
	if c != nil {
		return models.DBWithContext(c)
	}
	return models.DB
}
