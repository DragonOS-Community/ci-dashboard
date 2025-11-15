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

	// Commit ID过滤（前缀匹配）
	if params.CommitID != "" {
		query = query.Where("commit_id LIKE ? OR commit_short_id LIKE ?", params.CommitID+"%", params.CommitID+"%")
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

// MasterBranchStats master分支统计信息
type MasterBranchStats struct {
	TestRunID     uint64    `json:"test_run_id"`
	BranchName    string    `json:"branch_name"`
	CommitID      string    `json:"commit_id"`
	CommitShortID string    `json:"commit_short_id"`
	TestType      string    `json:"test_type"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
	TotalCases    int64     `json:"total_cases"`
	PassedCases   int64     `json:"passed_cases"`
	FailedCases   int64     `json:"failed_cases"`
	SkippedCases  int64     `json:"skipped_cases"`
	PassRate      float64   `json:"pass_rate"`
	Duration      int64     `json:"duration"` // 总耗时（毫秒）
}

// GetMasterBranchLatestStats 获取master分支最新的测试统计数据
func GetMasterBranchLatestStats() (*MasterBranchStats, error) {
	// 查找master分支最新的已完成测试运行
	var testRun models.TestRun
	if err := models.DB.Where("branch_name = ? AND status IN (?)", "master", []models.TestRunStatus{
		models.TestRunStatusPassed,
		models.TestRunStatusFailed,
	}).Order("created_at DESC").First(&testRun).Error; err != nil {
		// 如果没有找到已完成的，尝试找运行中的
		if err := models.DB.Where("branch_name = ?", "master").
			Order("created_at DESC").First(&testRun).Error; err != nil {
			return nil, fmt.Errorf("no test run found for master branch")
		}
	}

	// 统计测例数据
	var totalCases, passedCases, failedCases, skippedCases int64
	var duration int64

	if err := models.DB.Model(&models.TestCase{}).
		Where("test_run_id = ?", testRun.ID).
		Count(&totalCases).Error; err != nil {
		return nil, fmt.Errorf("failed to count test cases: %w", err)
	}

	if totalCases > 0 {
		if err := models.DB.Model(&models.TestCase{}).
			Where("test_run_id = ? AND status = ?", testRun.ID, models.TestCaseStatusPassed).
			Count(&passedCases).Error; err != nil {
			return nil, fmt.Errorf("failed to count passed cases: %w", err)
		}

		if err := models.DB.Model(&models.TestCase{}).
			Where("test_run_id = ? AND status = ?", testRun.ID, models.TestCaseStatusFailed).
			Count(&failedCases).Error; err != nil {
			return nil, fmt.Errorf("failed to count failed cases: %w", err)
		}

		if err := models.DB.Model(&models.TestCase{}).
			Where("test_run_id = ? AND status = ?", testRun.ID, models.TestCaseStatusSkipped).
			Count(&skippedCases).Error; err != nil {
			return nil, fmt.Errorf("failed to count skipped cases: %w", err)
		}

		// 计算总耗时
		var result struct {
			TotalDuration int64
		}
		if err := models.DB.Model(&models.TestCase{}).
			Select("COALESCE(SUM(duration_ms), 0) as total_duration").
			Where("test_run_id = ?", testRun.ID).
			Scan(&result).Error; err != nil {
			return nil, fmt.Errorf("failed to calculate duration: %w", err)
		}
		duration = result.TotalDuration
	}

	// 计算通过率
	passRate := 0.0
	if totalCases > 0 {
		passRate = float64(passedCases) / float64(totalCases) * 100.0
	}

	stats := &MasterBranchStats{
		TestRunID:     testRun.ID,
		BranchName:    testRun.BranchName,
		CommitID:      testRun.CommitID,
		CommitShortID: testRun.CommitShortID,
		TestType:      testRun.TestType,
		Status:        string(testRun.Status),
		CreatedAt:     testRun.CreatedAt,
		TotalCases:    totalCases,
		PassedCases:   passedCases,
		FailedCases:   failedCases,
		SkippedCases:  skippedCases,
		PassRate:      passRate,
		Duration:      duration,
	}

	return stats, nil
}
