package handlers

import (
	"fmt"
	"strconv"
	"time"

	"github.com/dragonos/dragonos-ci-dashboard/internal/models"
	"github.com/dragonos/dragonos-ci-dashboard/internal/services"
	"github.com/dragonos/dragonos-ci-dashboard/pkg/logger"
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

	// 记录查询参数
	logger.LogInfo(c, logger.ModuleHandler, "query_test_runs branch=%s commit_id=%s status=%s page=%d page_size=%d",
		params.Branch, params.CommitID, params.Status, params.Page, params.PageSize)

	// 公开接口只返回公开的记录
	testRuns, total, err := services.QueryTestRuns(c, params, false)
	if err != nil {
		logger.LogError(c, logger.ModuleHandler, err, "query_test_runs failed")
		response.InternalServerError(c, "Failed to query test runs")
		return
	}

	logger.LogInfo(c, logger.ModuleHandler, "query_test_runs success total=%d count=%d", total, len(testRuns))

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
		logger.LogWarn(c, logger.ModuleHandler, "invalid_test_run_id id=%s error=%s", idStr, err.Error())
		response.BadRequest(c, "Invalid test run ID")
		return
	}

	logger.LogInfo(c, logger.ModuleHandler, "get_test_run_by_id test_run_id=%d", id)

	testRun, err := services.GetTestRunByID(c, id)
	if err != nil {
		logger.LogWarn(c, logger.ModuleHandler, "test_run_not_found test_run_id=%d", id)
		response.NotFound(c, "Test run not found")
		return
	}

	// 公开接口只返回公开的记录
	if !testRun.IsPublic {
		logger.LogWarn(c, logger.ModuleHandler, "test_run_not_public test_run_id=%d", id)
		response.NotFound(c, "Test run not found")
		return
	}

	logger.LogInfo(c, logger.ModuleHandler, "get_test_run_success test_run_id=%d branch=%s status=%s",
		id, testRun.BranchName, testRun.Status)

	response.Success(c, testRun)
}

// GetTestCasesByTestRunID 获取测例列表（公开接口）
func GetTestCasesByTestRunID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		logger.LogWarn(c, logger.ModuleHandler, "invalid_test_run_id id=%s error=%s", idStr, err.Error())
		response.BadRequest(c, "Invalid test run ID")
		return
	}

	logger.LogInfo(c, logger.ModuleHandler, "get_test_cases_by_test_run_id test_run_id=%d", id)

	// 检查测试运行是否存在且为公开
	testRun, err := services.GetTestRunByID(c, id)
	if err != nil {
		logger.LogWarn(c, logger.ModuleHandler, "test_run_not_found test_run_id=%d", id)
		response.NotFound(c, "Test run not found")
		return
	}
	if !testRun.IsPublic {
		logger.LogWarn(c, logger.ModuleHandler, "test_run_not_public test_run_id=%d", id)
		response.NotFound(c, "Test run not found")
		return
	}

	testCases, err := services.GetTestCasesByTestRunID(c, id)
	if err != nil {
		logger.LogError(c, logger.ModuleHandler, err, "get_test_cases failed test_run_id=%d", id)
		response.InternalServerError(c, "Failed to get test cases")
		return
	}

	logger.LogInfo(c, logger.ModuleHandler, "get_test_cases_success test_run_id=%d count=%d", id, len(testCases))

	response.Success(c, testCases)
}

// GetMasterBranchStats 获取master分支最新测试统计数据（公开接口）
func GetMasterBranchStats(c *gin.Context) {
	logger.LogInfo(c, logger.ModuleHandler, "get_master_branch_stats")

	stats, err := services.GetMasterBranchLatestStats(c)
	if err != nil {
		logger.LogWarn(c, logger.ModuleHandler, "no_master_branch_stats found")
		response.NotFound(c, "No test run found for master branch")
		return
	}

	logger.LogInfo(c, logger.ModuleHandler, "get_master_branch_stats success")
	response.Success(c, stats)
}

// CreateTestRun 创建测试运行（受保护接口）
func CreateTestRun(c *gin.Context) {
	var req struct {
		BranchName string `json:"branch_name" binding:"required"`
		CommitID   string `json:"commit_id" binding:"required"`
		TestType   string `json:"test_type" binding:"required"`
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
		logger.LogWarn(c, logger.ModuleHandler, "create_test_run invalid_request error=%s", err.Error())
		response.BadRequest(c, err.Error())
		return
	}

	logger.LogInfo(c, logger.ModuleHandler, "create_test_run branch=%s commit_id=%s test_type=%s test_cases_count=%d",
		req.BranchName, req.CommitID, req.TestType, len(req.TestCases))

	// 验证 commit_id 最少8位
	if len(req.CommitID) < 8 {
		logger.LogWarn(c, logger.ModuleHandler, "create_test_run invalid_commit_id commit_id=%s", req.CommitID)
		response.BadRequest(c, "commit_id must be at least 8 characters")
		return
	}

	// 验证 test_type 必须为 gvisor
	if req.TestType != string(models.TestTypeGvisor) {
		logger.LogWarn(c, logger.ModuleHandler, "create_test_run invalid_test_type test_type=%s", req.TestType)
		response.BadRequest(c, fmt.Sprintf("test_type must be '%s'", models.TestTypeGvisor))
		return
	}
	testType := req.TestType

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
			logger.LogWarn(c, logger.ModuleHandler, "create_test_run error_log_too_long test_case=%s length=%d",
				req.TestCases[i].Name, len(req.TestCases[i].ErrorLog))
			response.BadRequest(c, "error_log exceeds maximum length of 2048 characters")
			return
		}
		if len(req.TestCases[i].DebugLog) > maxLogLength {
			logger.LogWarn(c, logger.ModuleHandler, "create_test_run debug_log_too_long test_case=%s length=%d",
				req.TestCases[i].Name, len(req.TestCases[i].DebugLog))
			response.BadRequest(c, "debug_log exceeds maximum length of 2048 characters")
			return
		}
	}

	// 创建测试运行
	testRun, err := services.CreateTestRun(
		c,
		defaultProjectID,
		req.BranchName,
		req.CommitID,
		commitShortID,
		testType,
	)
	if err != nil {
		logger.LogError(c, logger.ModuleHandler, err, "create_test_run failed branch=%s commit_id=%s",
			req.BranchName, req.CommitID)
		response.InternalServerError(c, "Failed to create test run")
		return
	}

	logger.LogInfo(c, logger.ModuleHandler, "create_test_run success test_run_id=%d", testRun.ID)

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

		if err := services.BatchCreateTestCases(c, testRun.ID, testCases); err != nil {
			logger.LogError(c, logger.ModuleHandler, err, "batch_create_test_cases failed test_run_id=%d count=%d",
				testRun.ID, len(testCases))
			response.InternalServerError(c, "Failed to create test cases")
			return
		}

		logger.LogInfo(c, logger.ModuleHandler, "batch_create_test_cases success test_run_id=%d count=%d",
			testRun.ID, len(testCases))

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

		logger.LogInfo(c, logger.ModuleHandler, "test_run_completed test_run_id=%d status=%s",
			testRun.ID, finalStatus)
	}

	// 重新加载关联数据
	testRun, _ = services.GetTestRunByID(c, testRun.ID)

	logger.LogInfo(c, logger.ModuleHandler, "create_test_run completed test_run_id=%d", testRun.ID)
	response.Success(c, testRun)
}
