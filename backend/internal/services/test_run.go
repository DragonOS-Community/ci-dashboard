package services

import (
	"fmt"
	"time"

	"github.com/dragonos/dragonos-ci-dashboard/internal/models"
)

// TestRunQueryParams 测试运行查询参数
type TestRunQueryParams struct {
	Branch       string
	CommitID     string
	StartTime    *time.Time
	EndTime      *time.Time
	Status       string
	TestCaseName string
	Page         int
	PageSize     int
}

// CreateTestRun 创建测试运行
func CreateTestRun(projectID uint64, branchName, commitID, commitShortID, testType string) (*models.TestRun, error) {
	testRun := &models.TestRun{
		ProjectID:     projectID,
		BranchName:    branchName,
		CommitID:      commitID,
		CommitShortID: commitShortID,
		TestType:      testType,
		Status:        models.TestRunStatusRunning,
	}

	if err := models.DB.Create(testRun).Error; err != nil {
		return nil, fmt.Errorf("failed to create test run: %w", err)
	}

	return testRun, nil
}

// GetTestRunByID 根据ID获取测试运行
func GetTestRunByID(id uint64) (*models.TestRun, error) {
	var testRun models.TestRun
	if err := models.DB.Preload("Project").
		Preload("TestCases").
		Preload("OutputFiles").
		First(&testRun, id).Error; err != nil {
		return nil, err
	}
	return &testRun, nil
}

// QueryTestRuns 查询测试运行列表
func QueryTestRuns(params TestRunQueryParams) ([]models.TestRun, int64, error) {
	var testRuns []models.TestRun
	var total int64

	query := models.DB.Model(&models.TestRun{}).Preload("Project")

	// 分支名过滤（模糊匹配）
	if params.Branch != "" {
		query = query.Where("branch_name LIKE ?", "%"+params.Branch+"%")
	}

	// Commit ID过滤（支持完整或短ID）
	if params.CommitID != "" {
		query = query.Where("commit_id = ? OR commit_short_id = ?", params.CommitID, params.CommitID)
	}

	// 时间范围过滤
	if params.StartTime != nil {
		query = query.Where("created_at >= ?", *params.StartTime)
	}
	if params.EndTime != nil {
		query = query.Where("created_at <= ?", *params.EndTime)
	}

	// 状态过滤
	if params.Status != "" && params.Status != "all" {
		query = query.Where("status = ?", params.Status)
	}

	// 测例名称过滤（通过关联查询）
	if params.TestCaseName != "" {
		query = query.Joins("JOIN test_cases ON test_cases.test_run_id = test_runs.id").
			Where("test_cases.name LIKE ?", "%"+params.TestCaseName+"%")
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页
	page := params.Page
	if page < 1 {
		page = 1
	}
	pageSize := params.PageSize
	if pageSize < 1 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize

	// 查询数据
	if err := query.Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&testRuns).Error; err != nil {
		return nil, 0, err
	}

	return testRuns, total, nil
}

// UpdateTestRunStatus 更新测试运行状态
func UpdateTestRunStatus(id uint64, status models.TestRunStatus) error {
	testRun, err := GetTestRunByID(id)
	if err != nil {
		return err
	}

	testRun.Complete(status)
	return models.DB.Save(testRun).Error
}

// CompleteTestRun 完成测试运行
func CompleteTestRun(id uint64, status models.TestRunStatus) error {
	testRun, err := GetTestRunByID(id)
	if err != nil {
		return err
	}

	testRun.Complete(status)
	return models.DB.Save(testRun).Error
}
