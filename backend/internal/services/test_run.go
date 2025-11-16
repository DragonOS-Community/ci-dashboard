package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/dragonos/dragonos-ci-dashboard/internal/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// TestRunQueryParams 测试运行查询参数
type TestRunQueryParams struct {
	Branch       string
	CommitID     string
	TestType     string
	StartTime    *time.Time
	EndTime      *time.Time
	Status       string
	TestCaseName string
	Page         int
	PageSize     int
}

// CreateTestRun 创建测试运行
func CreateTestRun(c *gin.Context, projectID uint64, branchName, commitID, commitShortID, testType string) (*models.TestRun, error) {
	testRun := &models.TestRun{
		ProjectID:     projectID,
		BranchName:    branchName,
		CommitID:      commitID,
		CommitShortID: commitShortID,
		TestType:      testType,
		Status:        models.TestRunStatusRunning,
	}

	db := getDB(c)
	if err := db.Create(testRun).Error; err != nil {
		return nil, fmt.Errorf("failed to create test run: %w", err)
	}

	return testRun, nil
}

// GetTestRunByID 根据ID获取测试运行
func GetTestRunByID(c *gin.Context, id uint64) (*models.TestRun, error) {
	var testRun models.TestRun
	db := getDB(c)
	if err := db.Preload("Project").
		Preload("TestCases").
		Preload("OutputFiles").
		First(&testRun, id).Error; err != nil {
		return nil, err
	}
	return &testRun, nil
}

// QueryTestRuns 查询测试运行列表
// includePrivate 为true时包含私有记录（管理员使用），为false时只返回公开记录（公开接口使用）
func QueryTestRuns(c *gin.Context, params TestRunQueryParams, includePrivate bool) ([]models.TestRun, int64, error) {
	var testRuns []models.TestRun
	var total int64

	db := getDB(c)
	query := db.Model(&models.TestRun{}).Preload("Project")

	// 如果不是管理员查询，只返回公开的记录
	if !includePrivate {
		query = query.Where("is_public = ?", true)
	}

	// 分支名过滤（模糊匹配）
	if params.Branch != "" {
		query = query.Where("branch_name LIKE ?", "%"+params.Branch+"%")
	}

	// Commit ID过滤（前缀匹配）
	if params.CommitID != "" {
		query = query.Where("commit_id LIKE ? OR commit_short_id LIKE ?", params.CommitID+"%", params.CommitID+"%")
	}

	// 测试类型过滤
	if params.TestType != "" {
		query = query.Where("test_type = ?", params.TestType)
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
func UpdateTestRunStatus(c *gin.Context, id uint64, status models.TestRunStatus) error {
	testRun, err := GetTestRunByID(c, id)
	if err != nil {
		return err
	}

	testRun.Complete(status)
	db := getDB(c)
	return db.Save(testRun).Error
}

// CompleteTestRun 完成测试运行
func CompleteTestRun(c *gin.Context, id uint64, status models.TestRunStatus) error {
	testRun, err := GetTestRunByID(c, id)
	if err != nil {
		return err
	}

	testRun.Complete(status)
	db := getDB(c)
	return db.Save(testRun).Error
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
func GetMasterBranchLatestStats(c *gin.Context) (*MasterBranchStats, error) {
	db := getDB(c)

	// 查找master分支最新的已完成测试运行（只返回公开的记录）
	var testRun models.TestRun
	if err := db.Where("branch_name = ? AND status IN (?) AND is_public = ?", "master", []models.TestRunStatus{
		models.TestRunStatusPassed,
		models.TestRunStatusFailed,
	}, true).Order("created_at DESC").First(&testRun).Error; err != nil {
		// 如果没有找到已完成的，尝试找运行中的（只返回公开的记录）
		if err := db.Where("branch_name = ? AND is_public = ?", "master", true).
			Order("created_at DESC").First(&testRun).Error; err != nil {
			return nil, fmt.Errorf("no test run found for master branch")
		}
	}

	// 统计测例数据
	var totalCases, passedCases, failedCases, skippedCases int64
	var duration int64

	if err := db.Model(&models.TestCase{}).
		Where("test_run_id = ?", testRun.ID).
		Count(&totalCases).Error; err != nil {
		return nil, fmt.Errorf("failed to count test cases: %w", err)
	}

	if totalCases > 0 {
		if err := db.Model(&models.TestCase{}).
			Where("test_run_id = ? AND status = ?", testRun.ID, models.TestCaseStatusPassed).
			Count(&passedCases).Error; err != nil {
			return nil, fmt.Errorf("failed to count passed cases: %w", err)
		}

		if err := db.Model(&models.TestCase{}).
			Where("test_run_id = ? AND status = ?", testRun.ID, models.TestCaseStatusFailed).
			Count(&failedCases).Error; err != nil {
			return nil, fmt.Errorf("failed to count failed cases: %w", err)
		}

		if err := db.Model(&models.TestCase{}).
			Where("test_run_id = ? AND status = ?", testRun.ID, models.TestCaseStatusSkipped).
			Count(&skippedCases).Error; err != nil {
			return nil, fmt.Errorf("failed to count skipped cases: %w", err)
		}

		// 计算总耗时
		var result struct {
			TotalDuration int64
		}
		if err := db.Model(&models.TestCase{}).
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

// DashboardStats 仪表板统计数据
type DashboardStats struct {
	// 统计卡片数据
	TotalTests      int64   `json:"total_tests"`       // 总测试数
	TodayRuns       int64   `json:"today_runs"`        // 今日运行
	SuccessRate     float64 `json:"success_rate"`      // 成功率
	AvgDuration     float64 `json:"avg_duration"`      // 平均耗时（秒）
	TotalTestsPrev  int64   `json:"total_tests_prev"`  // 上期总测试数（用于计算趋势）
	TodayRunsPrev   int64   `json:"today_runs_prev"`   // 上期今日运行
	SuccessRatePrev float64 `json:"success_rate_prev"` // 上期成功率
	AvgDurationPrev float64 `json:"avg_duration_prev"` // 上期平均耗时

	// 成功率统计
	SuccessCount int64 `json:"success_count"` // 成功数
	FailedCount  int64 `json:"failed_count"`  // 失败数
	SkippedCount int64 `json:"skipped_count"` // 跳过数
}

// TrendData 趋势数据点
type TrendData struct {
	Date  string `json:"date"`  // 日期
	Count int64  `json:"count"` // 数量
}

// GetDashboardStats 获取仪表板统计数据
func GetDashboardStats(c *gin.Context) (*DashboardStats, error) {
	db := getDB(c)

	now := time.Now()
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	todayEnd := todayStart.Add(24 * time.Hour)

	// 昨天的时间范围（用于计算趋势）
	yesterdayStart := todayStart.AddDate(0, 0, -1)
	yesterdayEnd := todayStart

	// 上周同期（7天前）
	lastWeekStart := todayStart.AddDate(0, 0, -7)

	stats := &DashboardStats{}

	// 1. 总测试数（所有测试运行）
	var totalTests int64
	if err := db.Model(&models.TestRun{}).Count(&totalTests).Error; err != nil {
		return nil, fmt.Errorf("failed to count total tests: %w", err)
	}
	stats.TotalTests = totalTests

	// 上期总测试数（7天前）
	var totalTestsPrev int64
	if err := db.Model(&models.TestRun{}).
		Where("created_at < ?", lastWeekStart).
		Count(&totalTestsPrev).Error; err != nil {
		return nil, fmt.Errorf("failed to count previous total tests: %w", err)
	}
	stats.TotalTestsPrev = totalTestsPrev

	// 2. 今日运行
	var todayRuns int64
	if err := db.Model(&models.TestRun{}).
		Where("created_at >= ? AND created_at < ?", todayStart, todayEnd).
		Count(&todayRuns).Error; err != nil {
		return nil, fmt.Errorf("failed to count today runs: %w", err)
	}
	stats.TodayRuns = todayRuns

	// 上期今日运行（昨天）
	var todayRunsPrev int64
	if err := db.Model(&models.TestRun{}).
		Where("created_at >= ? AND created_at < ?", yesterdayStart, yesterdayEnd).
		Count(&todayRunsPrev).Error; err != nil {
		return nil, fmt.Errorf("failed to count previous today runs: %w", err)
	}
	stats.TodayRunsPrev = todayRunsPrev

	// 3. 成功率统计（所有已完成的测试运行）
	var completedRuns []models.TestRun
	if err := db.Where("status IN (?)", []models.TestRunStatus{
		models.TestRunStatusPassed,
		models.TestRunStatusFailed,
	}).Find(&completedRuns).Error; err != nil {
		return nil, fmt.Errorf("failed to get completed runs: %w", err)
	}

	var successCount, failedCount int64
	for _, run := range completedRuns {
		if run.Status == models.TestRunStatusPassed {
			successCount++
		} else if run.Status == models.TestRunStatusFailed {
			failedCount++
		}
	}

	// 统计测例的成功/失败/跳过数
	var totalSuccessCases, totalFailedCases, totalSkippedCases int64
	if err := db.Model(&models.TestCase{}).
		Where("status = ?", models.TestCaseStatusPassed).
		Count(&totalSuccessCases).Error; err != nil {
		return nil, fmt.Errorf("failed to count success cases: %w", err)
	}
	if err := db.Model(&models.TestCase{}).
		Where("status = ?", models.TestCaseStatusFailed).
		Count(&totalFailedCases).Error; err != nil {
		return nil, fmt.Errorf("failed to count failed cases: %w", err)
	}
	if err := db.Model(&models.TestCase{}).
		Where("status = ?", models.TestCaseStatusSkipped).
		Count(&totalSkippedCases).Error; err != nil {
		return nil, fmt.Errorf("failed to count skipped cases: %w", err)
	}

	stats.SuccessCount = totalSuccessCases
	stats.FailedCount = totalFailedCases
	stats.SkippedCount = totalSkippedCases

	// 计算成功率（基于测例）
	totalCases := totalSuccessCases + totalFailedCases + totalSkippedCases
	if totalCases > 0 {
		stats.SuccessRate = float64(totalSuccessCases) / float64(totalCases) * 100.0
	}

	// 上期成功率（7天前的数据）
	var prevSuccessCases, prevFailedCases, prevSkippedCases int64
	var prevTestRuns []models.TestRun
	if err := db.Where("created_at < ? AND status IN (?)", lastWeekStart, []models.TestRunStatus{
		models.TestRunStatusPassed,
		models.TestRunStatusFailed,
	}).Find(&prevTestRuns).Error; err == nil {
		var prevRunIDs []uint64
		for _, run := range prevTestRuns {
			prevRunIDs = append(prevRunIDs, run.ID)
		}
		if len(prevRunIDs) > 0 {
			db.Model(&models.TestCase{}).
				Where("test_run_id IN (?) AND status = ?", prevRunIDs, models.TestCaseStatusPassed).
				Count(&prevSuccessCases)
			db.Model(&models.TestCase{}).
				Where("test_run_id IN (?) AND status = ?", prevRunIDs, models.TestCaseStatusFailed).
				Count(&prevFailedCases)
			db.Model(&models.TestCase{}).
				Where("test_run_id IN (?) AND status = ?", prevRunIDs, models.TestCaseStatusSkipped).
				Count(&prevSkippedCases)
		}
		prevTotalCases := prevSuccessCases + prevFailedCases + prevSkippedCases
		if prevTotalCases > 0 {
			stats.SuccessRatePrev = float64(prevSuccessCases) / float64(prevTotalCases) * 100.0
		}
	}

	// 4. 平均耗时（所有已完成的测试运行的平均耗时）
	var avgDurationResult struct {
		AvgDuration float64
	}
	if err := db.Model(&models.TestCase{}).
		Select("COALESCE(AVG(duration_ms), 0) / 1000.0 as avg_duration").
		Scan(&avgDurationResult).Error; err != nil {
		return nil, fmt.Errorf("failed to calculate avg duration: %w", err)
	}
	stats.AvgDuration = avgDurationResult.AvgDuration

	// 上期平均耗时
	var avgDurationPrevResult struct {
		AvgDuration float64
	}
	if err := db.Model(&models.TestCase{}).
		Joins("JOIN test_runs ON test_cases.test_run_id = test_runs.id").
		Where("test_runs.created_at < ?", lastWeekStart).
		Select("COALESCE(AVG(test_cases.duration_ms), 0) / 1000.0 as avg_duration").
		Scan(&avgDurationPrevResult).Error; err == nil {
		stats.AvgDurationPrev = avgDurationPrevResult.AvgDuration
	}

	return stats, nil
}

// GetDashboardTrend 获取仪表板趋势数据
func GetDashboardTrend(c *gin.Context, days int) ([]TrendData, error) {
	db := getDB(c)

	now := time.Now()
	startDate := now.AddDate(0, 0, -days)

	var results []struct {
		Date  string
		Count int64
	}

	// 根据数据库类型选择日期格式化函数
	// MySQL使用DATE函数
	if err := db.Model(&models.TestRun{}).
		Select("DATE(created_at) as date, COUNT(*) as count").
		Where("created_at >= ?", startDate).
		Group("DATE(created_at)").
		Order("date ASC").
		Scan(&results).Error; err != nil {
		return nil, fmt.Errorf("failed to get trend data: %w", err)
	}

	trendData := make([]TrendData, 0, len(results))
	for _, r := range results {
		trendData = append(trendData, TrendData{
			Date:  r.Date,
			Count: r.Count,
		})
	}

	return trendData, nil
}

// DeleteTestRun 删除测试运行
func DeleteTestRun(c *gin.Context, id uint64) error {
	db := getDB(c)

	// 检查测试运行是否存在
	var testRun models.TestRun
	if err := db.First(&testRun, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrTestRunNotFound
		}
		return fmt.Errorf("failed to get test run: %w", err)
	}

	// 由于外键约束，删除测试运行会自动删除关联的测例和输出文件
	if err := db.Delete(&testRun).Error; err != nil {
		return fmt.Errorf("failed to delete test run: %w", err)
	}

	return nil
}

// UpdateTestRunVisibility 更新测试运行的可见性
func UpdateTestRunVisibility(c *gin.Context, id uint64, isPublic bool) error {
	db := getDB(c)

	var testRun models.TestRun
	if err := db.First(&testRun, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrTestRunNotFound
		}
		return fmt.Errorf("failed to get test run: %w", err)
	}

	testRun.IsPublic = isPublic
	if err := db.Save(&testRun).Error; err != nil {
		return fmt.Errorf("failed to update test run visibility: %w", err)
	}

	return nil
}
