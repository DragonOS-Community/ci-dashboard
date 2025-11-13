package handlers

import (
	"net/http"
	"strconv"

	"github.com/dragonos/dragonos-ci-dashboard/internal/config"
	"github.com/dragonos/dragonos-ci-dashboard/internal/services"
	"github.com/dragonos/dragonos-ci-dashboard/pkg/response"
	"github.com/gin-gonic/gin"
)

// GetFileByID 下载文件（公开接口）
func GetFileByID(c *gin.Context) {
	testRunIDStr := c.Param("id")
	testRunID, err := strconv.ParseUint(testRunIDStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid test run ID")
		return
	}

	fileIDStr := c.Param("fileId")
	fileID, err := strconv.ParseUint(fileIDStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid file ID")
		return
	}

	file, err := services.GetFileByID(fileID)
	if err != nil {
		response.NotFound(c, "File not found")
		return
	}

	// 验证文件是否属于该测试运行
	if file.TestRunID != testRunID {
		response.NotFound(c, "File not found")
		return
	}

	// 打开文件
	fileHandle, err := services.OpenFile(file)
	if err != nil {
		response.InternalServerError(c, "Failed to open file")
		return
	}
	defer fileHandle.Close()

	// 设置响应头
	c.Header("Content-Disposition", "attachment; filename="+file.Filename)
	c.Header("Content-Type", file.MimeType)
	c.Header("Content-Length", strconv.FormatUint(file.FileSize, 10))

	// 流式传输文件
	c.DataFromReader(http.StatusOK, int64(file.FileSize), file.MimeType, fileHandle, map[string]string{})
}

// UploadFile 上传文件（受保护接口）
func UploadFile(c *gin.Context) {
	testRunIDStr := c.Param("id")
	testRunID, err := strconv.ParseUint(testRunIDStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid test run ID")
		return
	}

	// 验证测试运行是否存在
	_, err = services.GetTestRunByID(testRunID)
	if err != nil {
		response.NotFound(c, "Test run not found")
		return
	}

	// 获取上传的文件
	file, err := c.FormFile("file")
	if err != nil {
		response.BadRequest(c, "No file uploaded")
		return
	}

	// 检查文件大小
	if file.Size > config.AppConfig.Storage.MaxFileSize {
		response.BadRequest(c, "File size exceeds limit")
		return
	}

	// 打开文件
	src, err := file.Open()
	if err != nil {
		response.InternalServerError(c, "Failed to open file")
		return
	}
	defer src.Close()

	// 保存文件
	outputFile, err := services.SaveFile(testRunID, file.Filename, src)
	if err != nil {
		response.InternalServerError(c, "Failed to save file")
		return
	}

	response.Success(c, outputFile)
}

// GetFilesByTestRunID 获取文件列表（公开接口）
func GetFilesByTestRunID(c *gin.Context) {
	testRunIDStr := c.Param("id")
	testRunID, err := strconv.ParseUint(testRunIDStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid test run ID")
		return
	}

	files, err := services.GetFilesByTestRunID(testRunID)
	if err != nil {
		response.InternalServerError(c, "Failed to get files")
		return
	}

	response.Success(c, files)
}
