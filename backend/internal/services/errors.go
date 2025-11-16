package services

import "errors"

// 统一定义所有服务层错误
var (
	// 项目相关错误
	ErrProjectExists   = errors.New("project with this name already exists")
	ErrProjectNotFound = errors.New("project not found")

	// 测试运行相关错误
	ErrTestRunNotFound = errors.New("test run not found")
)
