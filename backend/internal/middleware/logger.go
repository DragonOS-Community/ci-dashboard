package middleware

import (
	"runtime"
	"time"

	"github.com/dragonos/dragonos-ci-dashboard/pkg/logger"
	"github.com/dragonos/dragonos-ci-dashboard/pkg/response"
	"github.com/gin-gonic/gin"
)

// Logger 日志中间件
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 记录请求开始
		logger.LogRequestStart(c)

		// 记录开始时间
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// 处理请求
		c.Next()

		// 计算耗时
		latency := time.Since(start)

		// 记录请求结束
		logger.LogRequestEnd(c, latency)

		// 如果有错误消息，记录错误日志
		if len(c.Errors) > 0 {
			for _, err := range c.Errors {
				logger.LogError(c, logger.ModuleHTTP, err, "request_error path=%s query=%s", path, raw)
			}
		}
	}
}

// Recovery 恢复中间件
func Recovery() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		// 获取堆栈信息
		stack := make([]byte, 4096)
		length := runtime.Stack(stack, false)
		stack = stack[:length]

		// 记录 panic 日志
		logger.LogPanic(c, recovered, stack)

		// 返回错误响应
		response.InternalServerError(c, "Internal server error")
		c.Abort()
	})
}
