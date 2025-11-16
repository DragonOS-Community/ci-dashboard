package handlers

import (
	"errors"
	"strconv"

	"github.com/dragonos/dragonos-ci-dashboard/internal/models"
	"github.com/dragonos/dragonos-ci-dashboard/internal/services"
	"github.com/dragonos/dragonos-ci-dashboard/pkg/response"
	"github.com/gin-gonic/gin"
)

// AdminLogin 管理员登录
func AdminLogin(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// 认证用户
	user, err := services.AuthenticateUser(req.Username, req.Password)
	if err != nil {
		response.Unauthorized(c, "Invalid username or password")
		return
	}

	// 生成JWT token
	token, err := services.GenerateJWT(user.ID, user.Username, string(user.Role))
	if err != nil {
		response.InternalServerError(c, "Failed to generate token")
		return
	}

	response.Success(c, gin.H{
		"token":    token,
		"user_id":  user.ID,
		"username": user.Username,
		"role":     user.Role,
	})
}

// AdminRegister 管理员注册
func AdminRegister(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required,min=3,max=100"`
		Password string `json:"password" binding:"required,min=6"`
		Role     string `json:"role" binding:"omitempty,oneof=admin user"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// 默认角色为admin
	role := models.UserRoleAdmin
	if req.Role != "" {
		role = models.UserRole(req.Role)
	}

	// 创建用户
	user, err := services.CreateUser(req.Username, req.Password, role)
	if err != nil {
		if err.Error() == "username already exists" {
			response.BadRequest(c, "Username already exists")
			return
		}
		response.InternalServerError(c, "Failed to create user")
		return
	}

	response.Success(c, gin.H{
		"user_id":  user.ID,
		"username": user.Username,
		"role":     user.Role,
		"message":  "User registered successfully",
	})
}

// GetAPIKeys 获取API密钥列表
func GetAPIKeys(c *gin.Context) {
	keys, err := services.ListAPIKeys()
	if err != nil {
		response.InternalServerError(c, "Failed to get API keys")
		return
	}

	// 不返回密钥哈希值
	for i := range keys {
		keys[i].KeyHash = ""
	}

	response.Success(c, keys)
}

// CreateAPIKey 创建API密钥
func CreateAPIKey(c *gin.Context) {
	var req struct {
		Name      string  `json:"name" binding:"required"`
		ProjectID *uint64 `json:"project_id"`
		ExpiresAt *string `json:"expires_at"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// 创建API密钥（expiresAt 直接使用请求中的字符串）
	apiKey, key, err := services.CreateAPIKey(req.Name, req.ProjectID, req.ExpiresAt)
	if err != nil {
		response.InternalServerError(c, "Failed to create API key")
		return
	}

	response.Success(c, gin.H{
		"id":         apiKey.ID,
		"name":       apiKey.Name,
		"project_id": apiKey.ProjectID,
		"api_key":    key, // 只在创建时返回一次
		"created_at": apiKey.CreatedAt,
		"expires_at": apiKey.ExpiresAt,
	})
}

// DeleteAPIKey 删除API密钥
func DeleteAPIKey(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid API key ID")
		return
	}

	if err := services.DeleteAPIKey(id); err != nil {
		response.InternalServerError(c, "Failed to delete API key")
		return
	}

	response.Success(c, nil)
}

// GetProfile 获取当前用户信息
func GetProfile(c *gin.Context) {
	// 从中间件获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	// 转换为uint64
	id, ok := userID.(uint64)
	if !ok {
		response.InternalServerError(c, "Invalid user ID type")
		return
	}

	// 获取用户信息
	user, err := services.GetUserByID(id)
	if err != nil {
		response.NotFound(c, "User not found")
		return
	}

	// 返回用户信息（不包含密码哈希）
	response.Success(c, gin.H{
		"id":         user.ID,
		"username":   user.Username,
		"role":       user.Role,
		"created_at": user.CreatedAt,
		"updated_at": user.UpdatedAt,
	})
}

// UpdatePassword 更新用户密码
func UpdatePassword(c *gin.Context) {
	// 从中间件获取用户名
	username, exists := c.Get("username")
	if !exists {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	usernameStr, ok := username.(string)
	if !ok {
		response.InternalServerError(c, "Invalid username type")
		return
	}

	var req struct {
		OldPassword string `json:"old_password" binding:"required"`
		NewPassword string `json:"new_password" binding:"required,min=6"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// 验证旧密码
	user, err := services.AuthenticateUser(usernameStr, req.OldPassword)
	if err != nil {
		response.BadRequest(c, "Old password is incorrect")
		return
	}

	// 更新密码
	if err := services.UpdateUserPasswordByID(user.ID, req.NewPassword); err != nil {
		response.InternalServerError(c, "Failed to update password")
		return
	}

	response.Success(c, gin.H{
		"message": "Password updated successfully",
	})
}

// GetDashboardStats 获取仪表板统计数据（管理接口）
func GetDashboardStats(c *gin.Context) {
	stats, err := services.GetDashboardStats()
	if err != nil {
		response.InternalServerError(c, "Failed to get dashboard stats")
		return
	}

	response.Success(c, stats)
}

// GetDashboardTrend 获取仪表板趋势数据（管理接口）
func GetDashboardTrend(c *gin.Context) {
	daysStr := c.DefaultQuery("days", "7")
	days, err := strconv.Atoi(daysStr)
	if err != nil || days <= 0 {
		days = 7
	}

	// 限制天数范围
	if days > 365 {
		days = 365
	}

	trendData, err := services.GetDashboardTrend(days)
	if err != nil {
		response.InternalServerError(c, "Failed to get dashboard trend")
		return
	}

	response.Success(c, trendData)
}

// GetProjects 获取项目列表
func GetProjects(c *gin.Context) {
	projects, err := services.ListProjects()
	if err != nil {
		response.InternalServerError(c, "Failed to get projects")
		return
	}

	response.Success(c, projects)
}

// GetProjectByID 根据ID获取项目
func GetProjectByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid project ID")
		return
	}

	project, err := services.GetProjectByID(id)
	if err != nil {
		response.NotFound(c, "Project not found")
		return
	}

	response.Success(c, project)
}

// CreateProject 创建项目
func CreateProject(c *gin.Context) {
	var req struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	project, err := services.CreateProject(req.Name, req.Description)
	if err != nil {
		if errors.Is(err, services.ErrProjectExists) {
			response.BadRequest(c, err.Error())
			return
		}
		response.InternalServerError(c, "Failed to create project")
		return
	}

	response.Success(c, project)
}

// UpdateProject 更新项目
func UpdateProject(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid project ID")
		return
	}

	var req struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	project, err := services.UpdateProject(id, req.Name, req.Description)
	if err != nil {
		if errors.Is(err, services.ErrProjectExists) {
			response.BadRequest(c, err.Error())
			return
		}
		if errors.Is(err, services.ErrProjectNotFound) {
			response.NotFound(c, "Project not found")
			return
		}
		response.InternalServerError(c, "Failed to update project")
		return
	}

	response.Success(c, project)
}

// DeleteProject 删除项目
func DeleteProject(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid project ID")
		return
	}

	if err := services.DeleteProject(id); err != nil {
		if errors.Is(err, services.ErrProjectNotFound) {
			response.NotFound(c, "Project not found")
			return
		}
		response.InternalServerError(c, "Failed to delete project")
		return
	}

	response.Success(c, nil)
}

// DeleteTestRun 删除测试运行（管理员接口）
func DeleteTestRun(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid test run ID")
		return
	}

	if err := services.DeleteTestRun(id); err != nil {
		if errors.Is(err, services.ErrTestRunNotFound) {
			response.NotFound(c, "Test run not found")
			return
		}
		response.InternalServerError(c, "Failed to delete test run")
		return
	}

	response.Success(c, nil)
}

// UpdateTestRunVisibility 更新测试运行可见性（管理员接口）
func UpdateTestRunVisibility(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid test run ID")
		return
	}

	var req struct {
		IsPublic bool `json:"is_public"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := services.UpdateTestRunVisibility(id, req.IsPublic); err != nil {
		if errors.Is(err, services.ErrTestRunNotFound) {
			response.NotFound(c, "Test run not found")
			return
		}
		response.InternalServerError(c, "Failed to update test run visibility")
		return
	}

	// 返回更新后的测试运行
	testRun, err := services.GetTestRunByID(id)
	if err != nil {
		response.InternalServerError(c, "Failed to get updated test run")
		return
	}

	response.Success(c, testRun)
}

// GetTestRunsAdmin 获取测试运行列表（管理员接口，包含私有记录）
func GetTestRunsAdmin(c *gin.Context) {
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
	if testType := c.Query("test_type"); testType != "" {
		params.TestType = testType
	}
	if status := c.Query("status"); status != "" {
		params.Status = status
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

	// 管理员接口包含私有记录
	testRuns, total, err := services.QueryTestRuns(params, true)
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

// GetSystemConfigs 获取所有系统配置
func GetSystemConfigs(c *gin.Context) {
	configs, err := services.GetAllConfigs()
	if err != nil {
		response.InternalServerError(c, "Failed to get system configs")
		return
	}

	response.Success(c, configs)
}

// GetSystemConfig 根据key获取系统配置
func GetSystemConfig(c *gin.Context) {
	key := c.Param("key")
	if key == "" {
		response.BadRequest(c, "Config key is required")
		return
	}

	value, err := services.GetConfig(key)
	if err != nil {
		response.NotFound(c, "Config not found")
		return
	}

	response.Success(c, gin.H{
		"key":   key,
		"value": value,
	})
}

// UpdateSystemConfig 更新系统配置
func UpdateSystemConfig(c *gin.Context) {
	key := c.Param("key")
	if key == "" {
		response.BadRequest(c, "Config key is required")
		return
	}

	var req struct {
		Value       string `json:"value" binding:"required"`
		Description string `json:"description"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := services.SetConfig(key, req.Value, req.Description); err != nil {
		response.InternalServerError(c, "Failed to update system config")
		return
	}

	response.Success(c, gin.H{
		"key":     key,
		"value":   req.Value,
		"message": "Config updated successfully",
	})
}
