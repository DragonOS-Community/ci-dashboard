package api

import (
	"github.com/dragonos/dragonos-ci-dashboard/internal/api/handlers"
	"github.com/dragonos/dragonos-ci-dashboard/internal/config"
	"github.com/dragonos/dragonos-ci-dashboard/internal/middleware"
	"github.com/gin-gonic/gin"
)

// SetupRouter 设置路由
func SetupRouter() *gin.Engine {
	// 根据环境设置Gin模式
	if config.AppConfig.Log.Level == "debug" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()

	// 全局中间件
	r.Use(middleware.Logger())
	r.Use(middleware.Recovery())
	r.Use(middleware.CORS())

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// API路由组
	apiPrefix := config.AppConfig.Server.APIPrefix
	v1 := r.Group(apiPrefix)

	// 公开接口（无需认证）
	public := v1.Group("")
	{
		public.GET("/test-runs", handlers.GetTestRuns)
		public.GET("/test-runs/:id", handlers.GetTestRunByID)
		public.GET("/test-runs/:id/test-cases", handlers.GetTestCasesByTestRunID)
		public.GET("/test-runs/:id/files", handlers.GetFilesByTestRunID)
		public.GET("/test-runs/:id/output-files/:fileId", handlers.GetFileByID)
		public.GET("/stats/master", handlers.GetMasterBranchStats)
	}

	// 受保护接口（需要API Key）
	protected := v1.Group("")
	protected.Use(middleware.APIKeyAuth())
	{
		protected.POST("/test-runs", handlers.CreateTestRun)
		protected.POST("/test-runs/:id/output-files", handlers.UploadFile)
	}

	// 管理接口（需要JWT认证）
	admin := v1.Group("/admin")
	admin.Use(middleware.JWTAuth())
	{
		admin.GET("/api-keys", handlers.GetAPIKeys)
		admin.POST("/api-keys", handlers.CreateAPIKey)
		admin.DELETE("/api-keys/:id", handlers.DeleteAPIKey)
		// 项目管理接口
		admin.GET("/projects", handlers.GetProjects)
		admin.GET("/projects/:id", handlers.GetProjectByID)
		admin.POST("/projects", handlers.CreateProject)
		admin.PUT("/projects/:id", handlers.UpdateProject)
		admin.DELETE("/projects/:id", handlers.DeleteProject)
		// 个人面板接口
		admin.GET("/profile", handlers.GetProfile)
		admin.PUT("/profile/password", handlers.UpdatePassword)
		// 仪表板接口
		admin.GET("/dashboard/stats", handlers.GetDashboardStats)
		admin.GET("/dashboard/trend", handlers.GetDashboardTrend)
	}

	// 公开的管理登录和注册接口（不需要认证）
	v1.POST("/admin/login", handlers.AdminLogin)
	v1.POST("/admin/register", handlers.AdminRegister)

	return r
}
