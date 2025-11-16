package models

import (
	"fmt"

	"github.com/dragonos/dragonos-ci-dashboard/internal/config"
	"github.com/dragonos/dragonos-ci-dashboard/pkg/logger"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

var DB *gorm.DB

// InitDatabase 初始化数据库连接
func InitDatabase() error {
	dsn := config.AppConfig.Database.DSN()

	// 创建自定义 GORM logger
	gormLog := logger.NewGormLogger()

	// 根据配置的日志级别设置 GORM logger 级别
	var logLevel gormLogger.LogLevel
	switch config.AppConfig.Log.Level {
	case "debug":
		logLevel = gormLogger.Info // GORM 的 Info 级别对应详细的 SQL 日志
	case "info":
		logLevel = gormLogger.Info
	case "warn", "warning":
		logLevel = gormLogger.Warn
	case "error":
		logLevel = gormLogger.Error
	default:
		logLevel = gormLogger.Info
	}

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: gormLog.LogMode(logLevel),
	})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// 自动迁移（仅用于开发，生产环境应使用迁移文件）
	// 注意：在生产环境中，应该使用数据库迁移工具而不是自动迁移
	if err := DB.AutoMigrate(
		&Project{},
		&TestRun{},
		&TestCase{},
		&TestOutputFile{},
		&APIKey{},
		&User{},
		&SystemConfig{},
	); err != nil {
		return fmt.Errorf("failed to auto migrate: %w", err)
	}

	return nil
}

// CloseDatabase 关闭数据库连接
func CloseDatabase() error {
	if DB != nil {
		sqlDB, err := DB.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}
	return nil
}
