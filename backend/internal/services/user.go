package services

import (
	"errors"
	"fmt"

	"github.com/dragonos/dragonos-ci-dashboard/internal/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CreateUser 创建用户
func CreateUser(c *gin.Context, username, password string, role models.UserRole) (*models.User, error) {
	db := getDB(c)

	// 检查用户名是否已存在（使用Count避免产生record not found日志）
	var count int64
	if err := db.Model(&models.User{}).Where("username = ?", username).Count(&count).Error; err != nil {
		return nil, fmt.Errorf("failed to check username: %w", err)
	}
	if count > 0 {
		return nil, fmt.Errorf("username already exists")
	}

	// 哈希密码
	passwordHash, err := HashPassword(password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user := &models.User{
		Username:     username,
		PasswordHash: passwordHash,
		Role:         role,
	}

	if err := db.Create(user).Error; err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

// AuthenticateUser 用户认证（用于登录，可能没有 context）
func AuthenticateUser(username, password string) (*models.User, error) {
	var user models.User
	if err := models.DB.Where("username = ?", username).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("invalid username or password")
		}
		return nil, err
	}

	if !CheckPassword(password, user.PasswordHash) {
		return nil, errors.New("invalid username or password")
	}

	return &user, nil
}

// GetUserByID 根据ID获取用户
func GetUserByID(c *gin.Context, id uint64) (*models.User, error) {
	var user models.User
	db := getDB(c)
	if err := db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserByUsername 根据用户名获取用户
func GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	if err := models.DB.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// UpdateUserPassword 更新用户密码
func UpdateUserPassword(username, newPassword string) error {
	user, err := GetUserByUsername(username)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	// 哈希新密码
	passwordHash, err := HashPassword(newPassword)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// 更新密码
	user.PasswordHash = passwordHash
	if err := models.DB.Save(user).Error; err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	return nil
}

// UpdateUserPasswordByID 根据ID更新用户密码
func UpdateUserPasswordByID(c *gin.Context, id uint64, newPassword string) error {
	user, err := GetUserByID(c, id)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	// 哈希新密码
	passwordHash, err := HashPassword(newPassword)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// 更新密码
	user.PasswordHash = passwordHash
	db := getDB(c)
	if err := db.Save(user).Error; err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	return nil
}

// UpdateUserRole 更新用户角色
func UpdateUserRole(username string, role models.UserRole) error {
	user, err := GetUserByUsername(username)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	user.Role = role
	if err := models.DB.Save(user).Error; err != nil {
		return fmt.Errorf("failed to update role: %w", err)
	}

	return nil
}
