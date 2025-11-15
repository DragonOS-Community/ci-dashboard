package services

import (
	"errors"
	"fmt"

	"github.com/dragonos/dragonos-ci-dashboard/internal/models"
	"gorm.io/gorm"
)

// 定义错误
var (
	ErrProjectExists   = errors.New("project with this name already exists")
	ErrProjectNotFound = errors.New("project not found")
)

// ListProjects 列出所有项目
func ListProjects() ([]models.Project, error) {
	var projects []models.Project
	if err := models.DB.Find(&projects).Error; err != nil {
		return nil, err
	}
	return projects, nil
}

// GetProjectByID 根据ID获取项目
func GetProjectByID(id uint64) (*models.Project, error) {
	var project models.Project
	if err := models.DB.First(&project, id).Error; err != nil {
		return nil, err
	}
	return &project, nil
}

// CreateProject 创建项目
func CreateProject(name string, description string) (*models.Project, error) {
	// 检查项目名称是否已存在
	var existingProject models.Project
	if err := models.DB.Where("name = ?", name).First(&existingProject).Error; err == nil {
		return nil, fmt.Errorf("%w: %s", ErrProjectExists, name)
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("failed to check project existence: %w", err)
	}

	project := &models.Project{
		Name:        name,
		Description: description,
	}

	if err := models.DB.Create(project).Error; err != nil {
		return nil, fmt.Errorf("failed to create project: %w", err)
	}

	return project, nil
}

// UpdateProject 更新项目
func UpdateProject(id uint64, name string, description string) (*models.Project, error) {
	var project models.Project
	if err := models.DB.First(&project, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrProjectNotFound
		}
		return nil, fmt.Errorf("failed to get project: %w", err)
	}

	// 如果名称改变，检查新名称是否已存在
	if name != project.Name {
		var existingProject models.Project
		if err := models.DB.Where("name = ? AND id != ?", name, id).First(&existingProject).Error; err == nil {
			return nil, fmt.Errorf("%w: %s", ErrProjectExists, name)
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("failed to check project existence: %w", err)
		}
	}

	project.Name = name
	project.Description = description

	if err := models.DB.Save(&project).Error; err != nil {
		return nil, fmt.Errorf("failed to update project: %w", err)
	}

	return &project, nil
}

// DeleteProject 删除项目
func DeleteProject(id uint64) error {
	// 检查项目是否存在
	var project models.Project
	if err := models.DB.First(&project, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrProjectNotFound
		}
		return fmt.Errorf("failed to get project: %w", err)
	}

	// 由于外键约束，删除项目会自动删除关联的测试运行和API密钥
	if err := models.DB.Delete(&project).Error; err != nil {
		return fmt.Errorf("failed to delete project: %w", err)
	}

	return nil
}
