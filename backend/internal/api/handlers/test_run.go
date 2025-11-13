package handlers

import (
	"strconv"
	"time"

	"github.com/dragonos/dragonos-ci-dashboard/internal/models"
	"github.com/dragonos/dragonos-ci-dashboard/internal/services"
	"github.com/dragonos/dragonos-ci-dashboard/pkg/response"
	"github.com/gin-gonic/gin"
)

// GetTestRuns 获取测试运行列表（公开接口）
func GetTestRuns(c *gin.Context) {
	params := services.TestRunQueryParams{
		Page:     1,
		PageSize: 20,
	}

	// 解析查询参数
	if branch := c.Query("branch"); branch != "" {
		params.Branch = branch
	}
	if commitID := c.Query("commit_id"); commitID != "" {
		params.CommitID = commitID
	}
	if startTimeStr := c.Query("start_time"); startTimeStr != "" {
		if startTime, err := time.Parse(time.RFC3339, startTimeStr); err == nil {
			params.StartTime = &startTime
		}
	}
	if endTimeStr := c.Query("end_time"); endTimeStr != "" {
		if endTime, err := time.Parse(time.RFC3339, endTimeStr); err == nil {
			params.EndTime = &endTime
		}
	}
	if status := c.Query("status"); status != "" {
		params.Status = status
	}
	if testCaseName := c.Query("test_case_name"); testCaseName != "" {
		params.TestCaseName = testCaseName
	}
	if page := c.Query("page"); page != "" {
		if p, err := strconv.Atoi(page); err == nil && p > 0 {
			params.Page = p
		}
	}
	if pageSize := c.Query("page_size"); pageSize != "" {
		if ps, err := strconv.Atoi(pageSize); err == nil && ps > 0 {
			params.PageSize = ps
		}
	}

	testRuns, total, err := services.QueryTestRuns(params)
	if err != nil {
		response.InternalServerError(c, "Failed to query test runs")
		return
	}

	response.Success(c, gin.H{
		"test_runs": testRuns,
		"total":     total,
		"page":      params.Page,
		"page_size": params.PageSize,
	})
}

// GetTestRunByID 获取测试运行详情（公开接口）
func GetTestRunByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid test run ID")
		return
	}

	testRun, err := services.GetTestRunByID(id)
	if err != nil {
		response.NotFound(c, "Test run not found")
		return
	}

	response.Success(c, testRun)
}

// GetTestCasesByTestRunID 获取测例列表（公开接口）
func GetTestCasesByTestRunID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid test run ID")
		return
	}

	testCases, err := services.GetTestCasesByTestRunID(id)
	if err != nil {
		response.InternalServerError(c, "Failed to get test cases")
		return
	}

	response.Success(c, testCases)
}

// CreateTestRun 创建测试运行（受保护接口）
func CreateTestRun(c *gin.Context) {
	var req struct {
		BranchName string `json:"branch_name" binding:"required"`
		CommitID   string `json:"commit_id" binding:"required"`
		TestType   string `json:"test_type"`
		TestCases  []struct {
			Name       string `json:"name" binding:"required"`
			Status     string `json:"status" binding:"required"`
			DurationMs uint32 `json:"duration_ms"`
			ErrorLog   string `json:"error_log"`
			DebugLog   string `json:"debug_log"`
		} `json:"test_cases"`
		Status string `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// 验证 commit_id 最少8位
	if len(req.CommitID) < 8 {
		response.BadRequest(c, "commit_id must be at least 8 characters")
		return
	}

	// 验证 test_type 只能为 gvisor
	testType := req.TestType
	if testType == "" {
		testType = "gvisor"
	}
	if testType != "gvisor" {
		response.BadRequest(c, "test_type must be 'gvisor'")
		return
	}

	// 自动截取 commit_short_id（前10位）
	commitShortID := req.CommitID
	if len(commitShortID) > 10 {
		commitShortID = commitShortID[:10]
	}

	// 使用默认项目ID 1
	const defaultProjectID = 1

	// 验证并截断日志长度
	const maxLogLength = 2048
	for i := range req.TestCases {
		if len(req.TestCases[i].ErrorLog) > maxLogLength {
			response.BadRequest(c, "error_log exceeds maximum length of 2048 characters")
			return
		}
		if len(req.TestCases[i].DebugLog) > maxLogLength {
			response.BadRequest(c, "debug_log exceeds maximum length of 2048 characters")
			return
		}
	}

	// 创建测试运行
	testRun, err := services.CreateTestRun(
		defaultProjectID,
		req.BranchName,
		req.CommitID,
		commitShortID,
		testType,
	)
	if err != nil {
		response.InternalServerError(c, "Failed to create test run")
		return
	}

	// 批量创建测例
	if len(req.TestCases) > 0 {
		testCases := make([]struct {
			Name       string
			Status     models.TestCaseStatus
			DurationMs uint32
			ErrorLog   string
			DebugLog   string
		}, 0, len(req.TestCases))

		for _, tc := range req.TestCases {
			testCases = append(testCases, struct {
				Name       string
				Status     models.TestCaseStatus
				DurationMs uint32
				ErrorLog   string
				DebugLog   string
			}{
				Name:       tc.Name,
				Status:     models.TestCaseStatus(tc.Status),
				DurationMs: tc.DurationMs,
				ErrorLog:   tc.ErrorLog,
				DebugLog:   tc.DebugLog,
			})
		}

		if err := services.BatchCreateTestCases(testRun.ID, testCases); err != nil {
			response.InternalServerError(c, "Failed to create test cases")
			return
		}

		// 根据测例状态更新测试运行状态
		allPassed := true
		hasFailed := false
		for _, tc := range req.TestCases {
			if tc.Status == "failed" {
				hasFailed = true
				allPassed = false
				break
			} else if tc.Status != "passed" {
				allPassed = false
			}
		}

		var finalStatus models.TestRunStatus
		if hasFailed {
			finalStatus = models.TestRunStatusFailed
		} else if allPassed {
			finalStatus = models.TestRunStatusPassed
		} else {
			finalStatus = models.TestRunStatusPassed // 默认通过
		}

		// 如果请求中指定了状态，使用请求的状态
		if req.Status != "" {
			finalStatus = models.TestRunStatus(req.Status)
		}

		testRun.Complete(finalStatus)
		models.DB.Save(testRun)
	}

	// 重新加载关联数据
	testRun, _ = services.GetTestRunByID(testRun.ID)

	response.Success(c, testRun)
}
