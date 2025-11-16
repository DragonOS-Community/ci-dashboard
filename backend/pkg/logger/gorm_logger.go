package logger

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm/logger"
)

// GormLogger GORM 自定义日志记录器
type GormLogger struct {
	level logger.LogLevel
}

// NewGormLogger 创建新的 GORM 日志记录器
func NewGormLogger() *GormLogger {
	return &GormLogger{
		level: logger.Info,
	}
}

// LogMode 设置日志级别
func (l *GormLogger) LogMode(level logger.LogLevel) logger.Interface {
	newLogger := *l
	newLogger.level = level
	return &newLogger
}

// Info 记录信息日志
func (l *GormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	// GORM 的 LogLevel：Silent=1, Error=2, Warn=3, Info=4
	if l.level >= 4 {
		logger := WithModule(ModuleGORM)
		// 尝试从 context 获取 request_id
		if requestID := getRequestIDFromContext(ctx); requestID != "" {
			logger = logger.With("request_id", requestID)
		}
		// 将内容合并到 msg 中
		logger.Info(fmt.Sprintf(msg, data...))
	}
}

// Warn 记录警告日志
func (l *GormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	// GORM 的 LogLevel：Silent=1, Error=2, Warn=3, Info=4
	if l.level >= 3 {
		logger := WithModule(ModuleGORM)
		// 尝试从 context 获取 request_id
		if requestID := getRequestIDFromContext(ctx); requestID != "" {
			logger = logger.With("request_id", requestID)
		}
		// 将内容合并到 msg 中
		logger.Warn(fmt.Sprintf(msg, data...))
	}
}

// Error 记录错误日志
func (l *GormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	// GORM 的 LogLevel：Silent=1, Error=2, Warn=3, Info=4
	if l.level >= 2 {
		logger := WithModule(ModuleGORM)
		// 尝试从 context 获取 request_id
		if requestID := getRequestIDFromContext(ctx); requestID != "" {
			logger = logger.With("request_id", requestID)
		}
		// 将内容合并到 msg 中
		logger.Error(fmt.Sprintf(msg, data...))
	}
}

// Trace 记录 SQL 跟踪日志
func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.level <= logger.Silent {
		return
	}

	elapsed := time.Since(begin)
	sql, rows := fc()

	logger := WithModule(ModuleGORM)
	// 尝试从 context 获取 request_id
	if requestID := getRequestIDFromContext(ctx); requestID != "" {
		logger = logger.With("request_id", requestID)
	}

	// 构建日志消息，将 SQL 查询内容合并到 msg 中
	msg := fmt.Sprintf("[%.2fms] [rows:%d] %s", float64(elapsed.Nanoseconds())/1e6, rows, sql)

	// 根据日志级别和错误情况选择日志级别
	// GORM 的 LogLevel 是整数类型：Silent=1, Error=2, Warn=3, Info=4
	switch {
	case err != nil && l.level >= 2: // Error level
		logger.Error(msg,
			"error", err.Error(),
			"duration_ms", elapsed.Milliseconds(),
			"rows", rows,
		)
	case elapsed > 200*time.Millisecond && l.level >= 3: // Warn level
		logger.Warn(msg,
			"duration_ms", elapsed.Milliseconds(),
			"rows", rows,
		)
	case l.level >= 4: // Info level
		logger.Info(msg,
			"duration_ms", elapsed.Milliseconds(),
			"rows", rows,
		)
	}
}

// getRequestIDFromContext 从 context 中获取 request_id
func getRequestIDFromContext(ctx context.Context) string {
	if ctx == nil {
		return ""
	}

	// 尝试从 context 中获取 request_id
	// 支持多种 context 类型
	if requestID, ok := ctx.Value("request_id").(string); ok {
		return requestID
	}

	// 如果 context 是 gin.Context，尝试获取
	if ginCtx, ok := ctx.(interface {
		Value(key interface{}) interface{}
	}); ok {
		if requestID, ok := ginCtx.Value("request_id").(string); ok {
			return requestID
		}
	}

	return ""
}
