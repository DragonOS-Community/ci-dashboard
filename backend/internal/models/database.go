package models

import (
	"fmt"

	"github.com/dragonos/dragonos-ci-dashboard/internal/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// InitDatabase 初始化数据库连接
func InitDatabase() error {
	dsn := config.AppConfig.Database.DSN()

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
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
