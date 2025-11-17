package services

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"github.com/dragonos/dragonos-ci-dashboard/internal/config"
	"github.com/dragonos/dragonos-ci-dashboard/internal/models"
	"github.com/dragonos/dragonos-ci-dashboard/pkg/logger"
	"github.com/gin-gonic/gin"
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

// ValidateAPIKey 验证API Key（用于中间件，可能没有 context）
func ValidateAPIKey(apiKey string) (*models.APIKey, error) {
	// 查询所有API Key（中间件调用，使用默认 DB）
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
func CreateAPIKey(c *gin.Context, name string, projectID *uint64, expiresAt *string) (*models.APIKey, string, error) {
	logger.LogInfo(
		c,
		logger.ModuleService,
		"create_api_key started name=%s project_id=%v expires_at=%v",
		name,
		projectID,
		expiresAt,
	)

	// 生成新的API Key
	apiKey, err := GenerateAPIKey()
	if err != nil {
		logger.LogError(c, logger.ModuleService, err, "create_api_key generate_failed name=%s", name)
		return nil, "", fmt.Errorf("failed to generate API key: %w", err)
	}

	// 哈希API Key
	keyHash, err := HashAPIKey(apiKey)
	if err != nil {
		logger.LogError(c, logger.ModuleService, err, "create_api_key hash_failed name=%s", name)
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

	db := getDB(c)
	logger.LogDebug(
		c,
		logger.ModuleService,
		"create_api_key saving_to_db name=%s project_id=%v",
		name,
		projectID,
	)
	if err := db.Create(newKey).Error; err != nil {
		logger.LogError(c, logger.ModuleService, err, "create_api_key db_create_failed name=%s", name)
		return nil, "", fmt.Errorf("failed to create API key: %w", err)
	}

	logger.LogInfo(
		c,
		logger.ModuleService,
		"create_api_key completed api_key_id=%d name=%s project_id=%v",
		newKey.ID,
		newKey.Name,
		newKey.ProjectID,
	)

	return newKey, apiKey, nil
}

// DeleteAPIKey 删除API Key
func DeleteAPIKey(c *gin.Context, id uint64) error {
	db := getDB(c)
	return db.Delete(&models.APIKey{}, id).Error
}

// ListAPIKeys 列出所有API Key
func ListAPIKeys(c *gin.Context) ([]models.APIKey, error) {
	var keys []models.APIKey
	db := getDB(c)
	if err := db.Preload("Project").Find(&keys).Error; err != nil {
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
