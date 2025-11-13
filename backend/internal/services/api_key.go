package services

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"github.com/dragonos/dragonos-ci-dashboard/internal/config"
	"github.com/dragonos/dragonos-ci-dashboard/internal/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// GenerateAPIKey 生成API Key
func GenerateAPIKey() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

// HashAPIKey 哈希API Key
func HashAPIKey(apiKey string) (string, error) {
	// 使用配置的salt和API Key一起哈希
	saltedKey := config.AppConfig.APIKey.HashSalt + apiKey
	hash, err := bcrypt.GenerateFromPassword([]byte(saltedKey), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// ValidateAPIKey 验证API Key
func ValidateAPIKey(apiKey string) (*models.APIKey, error) {
	// 查询所有API Key
	var apiKeys []models.APIKey
	if err := models.DB.Find(&apiKeys).Error; err != nil {
		return nil, err
	}

	// 遍历所有API Key，尝试匹配
	for _, key := range apiKeys {
		saltedKey := config.AppConfig.APIKey.HashSalt + apiKey
		if err := bcrypt.CompareHashAndPassword([]byte(key.KeyHash), []byte(saltedKey)); err == nil {
			return &key, nil
		}
	}

	return nil, gorm.ErrRecordNotFound
}

// CreateAPIKey 创建API Key
func CreateAPIKey(name string, projectID *uint64, expiresAt *string) (*models.APIKey, string, error) {
	// 生成新的API Key
	apiKey, err := GenerateAPIKey()
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate API key: %w", err)
	}

	// 哈希API Key
	keyHash, err := HashAPIKey(apiKey)
	if err != nil {
		return nil, "", fmt.Errorf("failed to hash API key: %w", err)
	}

	// 创建记录
	newKey := &models.APIKey{
		Name:      name,
		KeyHash:   keyHash,
		ProjectID: projectID,
	}

	if expiresAt != nil && *expiresAt != "" {
		// 解析过期时间（这里简化处理，实际应该解析时间字符串）
		// 暂时留空，后续可以完善
	}

	if err := models.DB.Create(newKey).Error; err != nil {
		return nil, "", fmt.Errorf("failed to create API key: %w", err)
	}

	return newKey, apiKey, nil
}

// DeleteAPIKey 删除API Key
func DeleteAPIKey(id uint64) error {
	return models.DB.Delete(&models.APIKey{}, id).Error
}

// ListAPIKeys 列出所有API Key
func ListAPIKeys() ([]models.APIKey, error) {
	var keys []models.APIKey
	if err := models.DB.Preload("Project").Find(&keys).Error; err != nil {
		return nil, err
	}
	return keys, nil
}

// GetAPIKeyByID 根据ID获取API Key
func GetAPIKeyByID(id uint64) (*models.APIKey, error) {
	var key models.APIKey
	if err := models.DB.Preload("Project").First(&key, id).Error; err != nil {
		return nil, err
	}
	return &key, nil
}
