package services

import (
	"fmt"
	"io"
	"mime"
	"os"
	"path/filepath"
	"time"

	"github.com/dragonos/dragonos-ci-dashboard/internal/config"
	"github.com/dragonos/dragonos-ci-dashboard/internal/models"
)

// SaveFile 保存文件
func SaveFile(testRunID uint64, filename string, fileContent io.Reader) (*models.TestOutputFile, error) {
	// 创建测试运行的文件目录
	fileDir := filepath.Join(config.AppConfig.Storage.Path, fmt.Sprintf("test_run_%d", testRunID))
	if err := os.MkdirAll(fileDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create file directory: %w", err)
	}

	// 生成唯一文件名（使用时间戳）
	timestamp := time.Now().Unix()
	ext := filepath.Ext(filename)
	baseName := filename[:len(filename)-len(ext)]
	uniqueFilename := fmt.Sprintf("%s_%d%s", baseName, timestamp, ext)
	filePath := filepath.Join(fileDir, uniqueFilename)

	// 保存文件
	dst, err := os.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to create file: %w", err)
	}
	defer dst.Close()

	fileSize, err := io.Copy(dst, fileContent)
	if err != nil {
		return nil, fmt.Errorf("failed to save file: %w", err)
	}

	// 获取MIME类型
	mimeType := mime.TypeByExtension(ext)
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}

	// 创建数据库记录
	outputFile := &models.TestOutputFile{
		TestRunID: testRunID,
		Filename:  filename,
		FilePath:  filePath,
		FileSize:  uint64(fileSize),
		MimeType:  mimeType,
	}

	if err := models.DB.Create(outputFile).Error; err != nil {
		// 如果数据库保存失败，删除已创建的文件
		os.Remove(filePath)
		return nil, fmt.Errorf("failed to create file record: %w", err)
	}

	return outputFile, nil
}

// GetFileByID 根据ID获取文件
func GetFileByID(id uint64) (*models.TestOutputFile, error) {
	var file models.TestOutputFile
	if err := models.DB.First(&file, id).Error; err != nil {
		return nil, err
	}
	return &file, nil
}

// GetFilesByTestRunID 根据测试运行ID获取文件列表
func GetFilesByTestRunID(testRunID uint64) ([]models.TestOutputFile, error) {
	var files []models.TestOutputFile
	if err := models.DB.Where("test_run_id = ?", testRunID).
		Order("created_at DESC").
		Find(&files).Error; err != nil {
		return nil, err
	}
	return files, nil
}

// OpenFile 打开文件
func OpenFile(file *models.TestOutputFile) (*os.File, error) {
	return os.Open(file.FilePath)
}
