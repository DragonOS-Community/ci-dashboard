package services

import (
	"fmt"

	"github.com/dragonos/dragonos-ci-dashboard/internal/models"
	"github.com/gin-gonic/gin"
)

// CreateTestCase 创建测例
func CreateTestCase(testRunID uint64, name string, status models.TestCaseStatus, durationMs uint32, errorLog, debugLog string) (*models.TestCase, error) {
	testCase := &models.TestCase{
		TestRunID:  testRunID,
		Name:       name,
		Status:     status,
		DurationMs: durationMs,
		ErrorLog:   errorLog,
		DebugLog:   debugLog,
	}

	if err := models.DB.Create(testCase).Error; err != nil {
		return nil, fmt.Errorf("failed to create test case: %w", err)
	}

	return testCase, nil
}

// BatchCreateTestCases 批量创建测例
func BatchCreateTestCases(c *gin.Context, testRunID uint64, testCases []struct {
	Name       string
	Status     models.TestCaseStatus
	DurationMs uint32
	ErrorLog   string
	DebugLog   string
}) error {
	if len(testCases) == 0 {
		return nil
	}

	cases := make([]models.TestCase, 0, len(testCases))
	for _, tc := range testCases {
		cases = append(cases, models.TestCase{
			TestRunID:  testRunID,
			Name:       tc.Name,
			Status:     tc.Status,
			DurationMs: tc.DurationMs,
			ErrorLog:   tc.ErrorLog,
			DebugLog:   tc.DebugLog,
		})
	}

	db := getDB(c)
	return db.CreateInBatches(cases, 100).Error
}

// GetTestCasesByTestRunID 根据测试运行ID获取测例列表
func GetTestCasesByTestRunID(c *gin.Context, testRunID uint64) ([]models.TestCase, error) {
	var testCases []models.TestCase
	db := getDB(c)
	if err := db.Where("test_run_id = ?", testRunID).
		Order("status DESC, name ASC").
		Find(&testCases).Error; err != nil {
		return nil, err
	}
	return testCases, nil
}

// GetTestCaseByID 根据ID获取测例
func GetTestCaseByID(id uint64) (*models.TestCase, error) {
	var testCase models.TestCase
	if err := models.DB.First(&testCase, id).Error; err != nil {
		return nil, err
	}
	return &testCase, nil
}
