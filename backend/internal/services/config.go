package services

import (
	"errors"
	"fmt"

	"github.com/dragonos/dragonos-ci-dashboard/internal/models"
	"gorm.io/gorm"
)

// GetConfig 根据key获取配置值
func GetConfig(key string) (string, error) {
	var config models.SystemConfig
	if err := models.DB.Where("config_key = ?", key).First(&config).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", fmt.Errorf("config not found: %s", key)
		}
		return "", fmt.Errorf("failed to get config: %w", err)
	}
	return config.ConfigValue, nil
}

// GetConfigWithDefault 根据key获取配置值，如果不存在则返回默认值
func GetConfigWithDefault(key string, defaultValue string) string {
	value, err := GetConfig(key)
	if err != nil {
		return defaultValue
	}
	return value
}

// SetConfig 设置配置值
func SetConfig(key string, value string, description string) error {
	var config models.SystemConfig
	err := models.DB.Where("config_key = ?", key).First(&config).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		// 创建新配置
		config = models.SystemConfig{
			ConfigKey:   key,
			ConfigValue: value,
			Description: description,
		}
		if err := models.DB.Create(&config).Error; err != nil {
			return fmt.Errorf("failed to create config: %w", err)
		}
		return nil
	} else if err != nil {
		return fmt.Errorf("failed to get config: %w", err)
	}

	// 更新现有配置
	config.ConfigValue = value
	if description != "" {
		config.Description = description
	}
	if err := models.DB.Save(&config).Error; err != nil {
		return fmt.Errorf("failed to update config: %w", err)
	}

	return nil
}

// GetAllConfigs 获取所有配置
func GetAllConfigs() ([]models.SystemConfig, error) {
	var configs []models.SystemConfig
	if err := models.DB.Find(&configs).Error; err != nil {
		return nil, fmt.Errorf("failed to get configs: %w", err)
	}
	return configs, nil
}

// GetConfigBool 获取布尔类型配置值
func GetConfigBool(key string) (bool, error) {
	value, err := GetConfig(key)
	if err != nil {
		return false, err
	}
	return value == "true", nil
}

// GetConfigBoolWithDefault 获取布尔类型配置值，如果不存在则返回默认值
func GetConfigBoolWithDefault(key string, defaultValue bool) bool {
	value, err := GetConfigBool(key)
	if err != nil {
		return defaultValue
	}
	return value
}

// SetConfigBool 设置布尔类型配置值
func SetConfigBool(key string, value bool, description string) error {
	strValue := "false"
	if value {
		strValue = "true"
	}
	return SetConfig(key, strValue, description)
}

// IsUploadOutputFilesAllowed 检查是否允许上传测试输出文件
func IsUploadOutputFilesAllowed() bool {
	return GetConfigBoolWithDefault("allow_upload_output_files", false)
}
