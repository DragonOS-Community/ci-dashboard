package middleware

import (
	"net/http"
	"strings"

	"github.com/dragonos/dragonos-ci-dashboard/internal/models"
	"github.com/dragonos/dragonos-ci-dashboard/internal/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// APIKeyAuth API Key认证中间件
func APIKeyAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "Missing Authorization header",
				"data":    nil,
			})
			c.Abort()
			return
		}

		// 提取Bearer token
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "Invalid Authorization header format",
				"data":    nil,
			})
			c.Abort()
			return
		}

		apiKey := parts[1]

		// 验证API Key
		key, err := services.ValidateAPIKey(apiKey)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusUnauthorized, gin.H{
					"code":    401,
					"message": "Invalid API key",
					"data":    nil,
				})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{
					"code":    500,
					"message": "Failed to validate API key",
					"data":    nil,
				})
			}
			c.Abort()
			return
		}

		// 检查是否过期
		if key.IsExpired() {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "API key expired",
				"data":    nil,
			})
			c.Abort()
			return
		}

		// 更新最后使用时间
		key.UpdateLastUsed()
		models.DB.Save(key)

		// 将API Key信息存储到上下文
		c.Set("api_key", key)
		c.Set("api_key_id", key.ID)
		if key.ProjectID != nil {
			c.Set("project_id", *key.ProjectID)
		}

		c.Next()
	}
}

// JWTAuth JWT认证中间件（用于后台管理）
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "Missing Authorization header",
				"data":    nil,
			})
			c.Abort()
			return
		}

		// 提取Bearer token
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "Invalid Authorization header format",
				"data":    nil,
			})
			c.Abort()
			return
		}

		tokenString := parts[1]

		// 验证JWT token
		claims, err := services.ValidateJWT(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "Invalid or expired token",
				"data":    nil,
			})
			c.Abort()
			return
		}

		// 将用户信息存储到上下文
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)

		c.Next()
	}
}
