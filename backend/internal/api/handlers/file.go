package handlers

import (
	"net/http"
	"strconv"

	"github.com/dragonos/dragonos-ci-dashboard/internal/config"
	"github.com/dragonos/dragonos-ci-dashboard/internal/services"
	"github.com/dragonos/dragonos-ci-dashboard/pkg/logger"
	"github.com/dragonos/dragonos-ci-dashboard/pkg/response"
	"github.com/gin-gonic/gin"
)

// GetFileByID 下载文件（公开接口）
func GetFileByID(c *gin.Context) {
	testRunIDStr := c.Param("id")
	testRunID, err := strconv.ParseUint(testRunIDStr, 10, 64)
	if err != nil {
		logger.LogWarn(c, logger.ModuleHandler, "get_file_by_id invalid_test_run_id id=%s error=%s", testRunIDStr, err.Error())
		response.BadRequest(c, "Invalid test run ID")
		return
	}

	fileIDStr := c.Param("fileId")
	fileID, err := strconv.ParseUint(fileIDStr, 10, 64)
	if err != nil {
		logger.LogWarn(c, logger.ModuleHandler, "get_file_by_id invalid_file_id id=%s error=%s", fileIDStr, err.Error())
		response.BadRequest(c, "Invalid file ID")
		return
	}

	logger.LogInfo(c, logger.ModuleHandler, "get_file_by_id test_run_id=%d file_id=%d", testRunID, fileID)

	// 检查测试运行是否存在且为公开
	testRun, err := services.GetTestRunByID(c, testRunID)
	if err != nil {
		logger.LogWarn(c, logger.ModuleHandler, "get_file_by_id test_run_not_found test_run_id=%d", testRunID)
		response.NotFound(c, "Test run not found")
		return
	}
	if !testRun.IsPublic {
		logger.LogWarn(c, logger.ModuleHandler, "get_file_by_id test_run_not_public test_run_id=%d", testRunID)
		response.NotFound(c, "Test run not found")
		return
	}

	file, err := services.GetFileByID(c, fileID)
	if err != nil {
		logger.LogWarn(c, logger.ModuleHandler, "get_file_by_id file_not_found file_id=%d", fileID)
		response.NotFound(c, "File not found")
		return
	}

	// 验证文件是否属于该测试运行
	if file.TestRunID != testRunID {
		logger.LogWarn(c, logger.ModuleHandler, "get_file_by_id file_mismatch test_run_id=%d file_test_run_id=%d", testRunID, file.TestRunID)
		response.NotFound(c, "File not found")
		return
	}

	// 打开文件
	fileHandle, err := services.OpenFile(file)
	if err != nil {
		logger.LogError(c, logger.ModuleHandler, err, "get_file_by_id open_file failed file_id=%d filename=%s", fileID, file.Filename)
		response.InternalServerError(c, "Failed to open file")
		return
	}
	defer fileHandle.Close()

	logger.LogInfo(c, logger.ModuleHandler, "get_file_by_id success file_id=%d filename=%s size=%d", fileID, file.Filename, file.FileSize)

	// 设置响应头
	c.Header("Content-Disposition", "attachment; filename="+file.Filename)
	c.Header("Content-Type", file.MimeType)
	c.Header("Content-Length", strconv.FormatUint(file.FileSize, 10))

	// 流式传输文件
	c.DataFromReader(http.StatusOK, int64(file.FileSize), file.MimeType, fileHandle, map[string]string{})
}

// UploadFile 上传文件（受保护接口）
func UploadFile(c *gin.Context) {
	// 检查是否允许上传测试输出文件
	if !services.IsUploadOutputFilesAllowed() {
		logger.LogWarn(c, logger.ModuleHandler, "upload_file not_allowed")
		response.Forbidden(c, "Uploading test output files is not allowed")
		return
	}

	testRunIDStr := c.Param("id")
	testRunID, err := strconv.ParseUint(testRunIDStr, 10, 64)
	if err != nil {
		logger.LogWarn(c, logger.ModuleHandler, "upload_file invalid_test_run_id id=%s error=%s", testRunIDStr, err.Error())
		response.BadRequest(c, "Invalid test run ID")
		return
	}

	logger.LogInfo(c, logger.ModuleHandler, "upload_file test_run_id=%d", testRunID)

	// 验证测试运行是否存在
	_, err = services.GetTestRunByID(c, testRunID)
	if err != nil {
		logger.LogWarn(c, logger.ModuleHandler, "upload_file test_run_not_found test_run_id=%d", testRunID)
		response.NotFound(c, "Test run not found")
		return
	}

	// 获取上传的文件
	file, err := c.FormFile("file")
	if err != nil {
		logger.LogWarn(c, logger.ModuleHandler, "upload_file no_file_uploaded error=%s", err.Error())
		response.BadRequest(c, "No file uploaded")
		return
	}

	logger.LogInfo(c, logger.ModuleHandler, "upload_file received filename=%s size=%d", file.Filename, file.Size)

	// 检查文件大小
	if file.Size > config.AppConfig.Storage.MaxFileSize {
		logger.LogWarn(c, logger.ModuleHandler, "upload_file file_too_large filename=%s size=%d max_size=%d",
			file.Filename, file.Size, config.AppConfig.Storage.MaxFileSize)
		response.BadRequest(c, "File size exceeds limit")
		return
	}

	// 打开文件
	src, err := file.Open()
	if err != nil {
		logger.LogError(c, logger.ModuleHandler, err, "upload_file open_file failed filename=%s", file.Filename)
		response.InternalServerError(c, "Failed to open file")
		return
	}
	defer src.Close()

	// 保存文件
	outputFile, err := services.SaveFile(c, testRunID, file.Filename, src)
	if err != nil {
		logger.LogError(c, logger.ModuleHandler, err, "upload_file save_file failed test_run_id=%d filename=%s", testRunID, file.Filename)
		response.InternalServerError(c, "Failed to save file")
		return
	}

	logger.LogInfo(c, logger.ModuleHandler, "upload_file success file_id=%d filename=%s", outputFile.ID, outputFile.Filename)
	response.Success(c, outputFile)
}

// GetFilesByTestRunID 获取文件列表（公开接口）
func GetFilesByTestRunID(c *gin.Context) {
	testRunIDStr := c.Param("id")
	testRunID, err := strconv.ParseUint(testRunIDStr, 10, 64)
	if err != nil {
		logger.LogWarn(c, logger.ModuleHandler, "get_files_by_test_run_id invalid_test_run_id id=%s error=%s", testRunIDStr, err.Error())
		response.BadRequest(c, "Invalid test run ID")
		return
	}

	logger.LogInfo(c, logger.ModuleHandler, "get_files_by_test_run_id test_run_id=%d", testRunID)

	// 检查测试运行是否存在且为公开
	testRun, err := services.GetTestRunByID(c, testRunID)
	if err != nil {
		logger.LogWarn(c, logger.ModuleHandler, "get_files_by_test_run_id test_run_not_found test_run_id=%d", testRunID)
		response.NotFound(c, "Test run not found")
		return
	}
	if !testRun.IsPublic {
		logger.LogWarn(c, logger.ModuleHandler, "get_files_by_test_run_id test_run_not_public test_run_id=%d", testRunID)
		response.NotFound(c, "Test run not found")
		return
	}

	files, err := services.GetFilesByTestRunID(c, testRunID)
	if err != nil {
		logger.LogError(c, logger.ModuleHandler, err, "get_files_by_test_run_id failed test_run_id=%d", testRunID)
		response.InternalServerError(c, "Failed to get files")
		return
	}

	logger.LogInfo(c, logger.ModuleHandler, "get_files_by_test_run_id success test_run_id=%d count=%d", testRunID, len(files))
	response.Success(c, files)
}
