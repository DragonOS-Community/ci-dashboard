package services

import (
	"errors"
	"fmt"

	"github.com/dragonos/dragonos-ci-dashboard/internal/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ListProjects 列出所有项目
func ListProjects(c *gin.Context) ([]models.Project, error) {
	var projects []models.Project
	db := getDB(c)
	if err := db.Find(&projects).Error; err != nil {
		return nil, err
	}
	return projects, nil
}

// GetProjectByID 根据ID获取项目
func GetProjectByID(c *gin.Context, id uint64) (*models.Project, error) {
	var project models.Project
	db := getDB(c)
	if err := db.First(&project, id).Error; err != nil {
		return nil, err
	}
	return &project, nil
}

// CreateProject 创建项目
func CreateProject(c *gin.Context, name string, description string) (*models.Project, error) {
	db := getDB(c)

	// 检查项目名称是否已存在
	var existingProject models.Project
	if err := db.Where("name = ?", name).First(&existingProject).Error; err == nil {
		return nil, fmt.Errorf("%w: %s", ErrProjectExists, name)
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("failed to check project existence: %w", err)
	}

	project := &models.Project{
		Name:        name,
		Description: description,
	}

	if err := db.Create(project).Error; err != nil {
		return nil, fmt.Errorf("failed to create project: %w", err)
	}

	return project, nil
}

// UpdateProject 更新项目
func UpdateProject(c *gin.Context, id uint64, name string, description string) (*models.Project, error) {
	db := getDB(c)

	var project models.Project
	if err := db.First(&project, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrProjectNotFound
		}
		return nil, fmt.Errorf("failed to get project: %w", err)
	}

	// 如果名称改变，检查新名称是否已存在
	if name != project.Name {
		var existingProject models.Project
		if err := db.Where("name = ? AND id != ?", name, id).First(&existingProject).Error; err == nil {
			return nil, fmt.Errorf("%w: %s", ErrProjectExists, name)
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("failed to check project existence: %w", err)
		}
	}

	project.Name = name
	project.Description = description

	if err := db.Save(&project).Error; err != nil {
		return nil, fmt.Errorf("failed to update project: %w", err)
	}

	return &project, nil
}

// DeleteProject 删除项目
func DeleteProject(c *gin.Context, id uint64) error {
	db := getDB(c)

	// 检查项目是否存在
	var project models.Project
	if err := db.First(&project, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrProjectNotFound
		}
		return fmt.Errorf("failed to get project: %w", err)
	}

	// 由于外键约束，删除项目会自动删除关联的测试运行和API密钥
	if err := db.Delete(&project).Error; err != nil {
		return fmt.Errorf("failed to delete project: %w", err)
	}

	return nil
}
