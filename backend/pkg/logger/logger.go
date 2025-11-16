package logger

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/dragonos/dragonos-ci-dashboard/internal/config"
	"github.com/gin-gonic/gin"
)

// 模块常量定义
const (
	ModuleHTTP     = "http"     // HTTP 中间件
	ModuleHandler  = "handler"  // 请求处理器
	ModuleService  = "service"  // 业务服务层
	ModuleGORM     = "gorm"     // GORM 数据库
	ModuleRecovery = "recovery" // 错误恢复
)

var (
	defaultLogger *slog.Logger
	logFile       *os.File
)

// InitLogger 初始化日志系统
func InitLogger() error {
	var writer io.Writer = os.Stdout

	// 如果配置了日志文件路径，则同时输出到文件
	if config.AppConfig.Log.FilePath != "" {
		// 确保日志文件目录存在
		logDir := filepath.Dir(config.AppConfig.Log.FilePath)
		if err := os.MkdirAll(logDir, 0755); err != nil {
			return err
		}

		// 打开日志文件（追加模式）
		file, err := os.OpenFile(config.AppConfig.Log.FilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			return err
		}

		logFile = file
		// 使用 MultiWriter 同时写入 stdout 和文件
		writer = io.MultiWriter(os.Stdout, file)
	}

	var handler slog.Handler

	// 根据配置选择日志格式
	if config.AppConfig.Log.Format == "json" {
		handler = slog.NewJSONHandler(writer, &slog.HandlerOptions{
			Level: parseLogLevel(config.AppConfig.Log.Level),
		})
	} else {
		handler = slog.NewTextHandler(writer, &slog.HandlerOptions{
			Level: parseLogLevel(config.AppConfig.Log.Level),
		})
	}

	defaultLogger = slog.New(handler)
	return nil
}

// CloseLogger 关闭日志文件（如果打开）
func CloseLogger() error {
	if logFile != nil {
		return logFile.Close()
	}
	return nil
}

// parseLogLevel 解析日志级别
func parseLogLevel(level string) slog.Level {
	switch level {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

// Logger 获取默认日志记录器
func Logger() *slog.Logger {
	if defaultLogger == nil {
		// 如果未初始化，使用默认配置
		handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		})
		defaultLogger = slog.New(handler)
	}
	return defaultLogger
}

// WithModule 创建带模块的日志记录器
func WithModule(module string) *slog.Logger {
	return Logger().With("module", module)
}

// WithRequestID 创建带 request_id 的日志记录器
func WithRequestID(logger *slog.Logger, requestID string) *slog.Logger {
	if requestID != "" {
		return logger.With("request_id", requestID)
	}
	return logger
}

// WithContext 从 gin.Context 创建带 request_id 的日志记录器
func WithContext(logger *slog.Logger, c *gin.Context) *slog.Logger {
	if c != nil {
		if requestID, exists := c.Get("request_id"); exists {
			if id, ok := requestID.(string); ok {
				return logger.With("request_id", id)
			}
		}
	}
	return logger
}

// LogRequestStart 记录 HTTP 请求开始日志
func LogRequestStart(c *gin.Context) {
	logger := WithModule(ModuleHTTP)
	logger = WithContext(logger, c)

	msg := fmt.Sprintf("request_start method=%s path=%s query=%s client_ip=%s",
		c.Request.Method,
		c.Request.URL.Path,
		c.Request.URL.RawQuery,
		c.ClientIP(),
	)
	logger.Info(msg, "user_agent", c.Request.UserAgent())
}

// LogRequestEnd 记录 HTTP 请求结束日志
func LogRequestEnd(c *gin.Context, latency time.Duration) {
	logger := WithModule(ModuleHTTP)
	logger = WithContext(logger, c)

	msg := fmt.Sprintf("request_end method=%s path=%s status=%d latency_ms=%d latency=%s client_ip=%s",
		c.Request.Method,
		c.Request.URL.Path,
		c.Writer.Status(),
		latency.Milliseconds(),
		latency.String(),
		c.ClientIP(),
	)
	logger.Info(msg)
}

// LogRequest 记录 HTTP 请求日志（兼容旧接口）
func LogRequest(c *gin.Context, module string, latency time.Duration) {
	LogRequestEnd(c, latency)
}

// LogError 记录错误日志（支持格式化字符串）
func LogError(c *gin.Context, module string, err error, format string, args ...interface{}) {
	logger := WithModule(module)
	logger = WithContext(logger, c)

	msg := fmt.Sprintf(format, args...)
	if err != nil {
		logger.Error(msg, "error", err.Error())
	} else {
		logger.Error(msg)
	}
}

// LogInfo 记录信息日志（支持格式化字符串）
func LogInfo(c *gin.Context, module string, format string, args ...interface{}) {
	logger := WithModule(module)
	logger = WithContext(logger, c)

	msg := fmt.Sprintf(format, args...)
	logger.Info(msg)
}

// LogDebug 记录调试日志（支持格式化字符串）
func LogDebug(c *gin.Context, module string, format string, args ...interface{}) {
	logger := WithModule(module)
	logger = WithContext(logger, c)

	msg := fmt.Sprintf(format, args...)
	logger.Debug(msg)
}

// LogWarn 记录警告日志（支持格式化字符串）
func LogWarn(c *gin.Context, module string, format string, args ...interface{}) {
	logger := WithModule(module)
	logger = WithContext(logger, c)

	msg := fmt.Sprintf(format, args...)
	logger.Warn(msg)
}

// LogPanic 记录 panic 日志
func LogPanic(c *gin.Context, recovered interface{}, stack []byte) {
	logger := WithModule(ModuleRecovery)
	logger = WithContext(logger, c)

	msg := fmt.Sprintf("panic_recovered panic=%v method=%s path=%s",
		recovered,
		c.Request.Method,
		c.Request.URL.Path,
	)
	logger.Error(msg, "stack", string(stack))
}

// GetCallerInfo 获取调用者信息（用于模块自动识别）
func GetCallerInfo() (string, int) {
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		return "unknown", 0
	}
	return file, line
}

// LogWithFields 记录带自定义字段的日志
func LogWithFields(ctx context.Context, module string, level slog.Level, msg string, fields map[string]interface{}) {
	logger := WithModule(module)

	// 如果有 context，尝试获取 request_id
	if ctx != nil {
		if ginCtx, ok := ctx.(*gin.Context); ok {
			logger = WithContext(logger, ginCtx)
		}
	}

	// 转换字段
	args := make([]any, 0, len(fields)*2)
	for k, v := range fields {
		args = append(args, k, v)
	}

	switch level {
	case slog.LevelDebug:
		logger.Debug(msg, args...)
	case slog.LevelInfo:
		logger.Info(msg, args...)
	case slog.LevelWarn:
		logger.Warn(msg, args...)
	case slog.LevelError:
		logger.Error(msg, args...)
	}
}

// LogJSON 记录 JSON 格式的日志（用于复杂数据结构）
func LogJSON(c *gin.Context, module string, level slog.Level, msg string, data interface{}) {
	logger := WithModule(module)
	logger = WithContext(logger, c)

	jsonData, err := json.Marshal(data)
	if err != nil {
		logger.Error("failed_to_marshal_json", "error", err.Error())
		return
	}

	args := []any{"data", string(jsonData)}

	switch level {
	case slog.LevelDebug:
		logger.Debug(msg, args...)
	case slog.LevelInfo:
		logger.Info(msg, args...)
	case slog.LevelWarn:
		logger.Warn(msg, args...)
	case slog.LevelError:
		logger.Error(msg, args...)
	}
}
