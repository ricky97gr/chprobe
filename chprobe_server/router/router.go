package router

import (
	"os"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/ricky97gr/chprobe/chprobe_server/controller"
	"github.com/ricky97gr/chprobe/chprobe_server/middleware"
	"github.com/ricky97gr/chprobe/chprobe_server/utils"
)

func Start() {
	gin.SetMode(utils.GinLogMode)
	utils.StartTime = time.Now().UnixMilli()

	engine := gin.New()
	engine.Use(middleware.Logger(), middleware.Recovery(), middleware.CORS())

	controller.InitPluginManager()

	// API路由组
	api := engine.Group("/api")
	{
		// 不需要认证的接口
		// 健康检查
		api.GET("/health", controller.HealthCheck)
		// 生成安装命令
		api.GET("/install", controller.GenerateInstallCommand)
		// 下载安装包
		api.GET("/download/:filename", controller.DownloadInstaller)
		// 插件下载（无token，agent校验+限速）
		api.GET("/plugin/download/:pluginId", controller.DownloadPlugin)
		// 查询下载任务状态（无token）
		api.GET("/download/status/:taskId", controller.GetDownloadStatus)
		// 下载插件文件（无token）
		api.GET("/download/file", controller.DownloadFile)
		// 登录
		api.POST("/login", controller.Login)

		// 需要认证的接口
		authApi := api.Group("/")
		authApi.Use(middleware.Auth(), middleware.OperationLogger())
		{
			// 客户端管理
			agent := authApi.Group("/agent")
			{
				// 获取客户端列表
				agent.GET("/list", controller.GetAgentList)
				// 获取客户端详情
				agent.GET("/detail/:uuid", controller.GetAgentDetail)
				// 更新客户端状态
				agent.PUT("/status/:uuid", controller.UpdateAgentStatus)
				// 删除客户端
				agent.DELETE("/delete/:uuid", controller.DeleteAgent)
			}
			// 主机管理
			host := authApi.Group("/host")
			{
				// 获取主机列表
				host.GET("/list", controller.GetHostList)
				// 获取主机详情
				host.GET("/detail/:uuid", controller.GetHostDetail)
			}
			// 镜像管理
			image := authApi.Group("/image")
			{
				// 获取镜像列表
				image.GET("/list", controller.GetImageList)
				// 获取镜像详情
				image.GET("/detail/:id", controller.GetImageDetail)
			}
			// 容器管理
			container := authApi.Group("/container")
			{
				// 获取容器列表
				container.GET("/list", controller.GetContainerList)
				// 获取容器详情
				container.GET("/detail/:id", controller.GetContainerDetail)
			}
			// 日志管理
			log := authApi.Group("/log")
			{
				// 获取操作日志列表
				log.GET("/operation/list", controller.GetOperationLogList)
				// 获取访问日志列表
				log.GET("/access/list", controller.GetAccessLogList)
				// 上报系统运行日志
				log.POST("/system/report", controller.ReportSystemLog)
				// 获取系统运行日志列表
				log.GET("/system/list", controller.GetSystemLogList)
				// 获取最新系统运行日志（仪表盘）
				log.GET("/system/latest", controller.GetLatestSystemLog)
				// 获取升级记录列表
				log.GET("/upgrade/list", controller.GetUpgradeRecordList)
			}
			// 用户管理
			user := authApi.Group("/user")
			{
				// 获取用户列表
				user.GET("/list", controller.GetUserList)
				// 新增用户
				user.POST("/create", controller.CreateUser)
				// 更新用户
				user.PUT("/update/:id", controller.UpdateUser)
				// 删除用户
				user.DELETE("/delete/:id", controller.DeleteUser)
				// 重置密码
				user.POST("/reset-password/:id", controller.ResetPassword)
				// 修改密码
				user.POST("/change-password", controller.ChangePassword)
			}

			// 授权管理
			license := authApi.Group("/license")
			{
				// 获取授权信息
				license.GET("/info", controller.GetLicenseInfo)
				// 获取授权详情
				license.GET("/detail/:id", controller.GetLicenseDetail)
				// 上传授权文件/字符串
				license.POST("/upload", controller.UploadLicense)
				// 删除授权
				license.DELETE("/delete/:id", controller.DeleteLicense)
			}

			// 插件管理
			plugin := authApi.Group("/plugin")
			{
				// 获取我的插件列表
				plugin.GET("/my", controller.GetMyPlugins)
				// 安装插件
				plugin.POST("/install", controller.InstallPlugin)
				// 卸载插件
				plugin.DELETE("/uninstall/:pluginId", controller.UninstallPlugin)
				// 切换插件状态
				plugin.PUT("/toggle/:pluginId", controller.TogglePlugin)
				// 更新插件状态
				plugin.POST("/update-status", controller.UpdatePluginStatus)
				// 创建下载任务
				plugin.POST("/download/task", controller.CreateDownloadTask)
				// 查询下载进度（兼容旧路径）
				plugin.GET("/download/status/:taskId", controller.GetDownloadStatus)
			}

			// 插件进程管理
			pluginManager := authApi.Group("/plugin-manager")
			{
				// 启动插件
				pluginManager.POST("/start", controller.StartPlugin)
				// 停止插件
				pluginManager.POST("/stop", controller.StopPlugin)
				// 重启插件
				pluginManager.POST("/restart", controller.RestartPlugin)
				// 获取插件状态
				pluginManager.GET("/status", controller.GetPluginStatus)
				// 获取所有插件
				pluginManager.GET("/list", controller.GetAllPlugins)
				// 获取插件路由
				pluginManager.GET("/routes", controller.GetPluginRoutes)
				// 查询路由
				pluginManager.POST("/route/query", controller.QueryRoute)
				// 转发请求
				pluginManager.POST("/forward", controller.ForwardRequest)
				// 插件健康检查
				pluginManager.GET("/health", controller.GetPluginHealth)
				// 列出已安装的插件
				pluginManager.GET("/installed", controller.ListInstalledPlugins)
				// 初始化插件
				pluginManager.POST("/initialize", controller.InitializePlugin)
				// 关闭所有插件
				pluginManager.POST("/shutdown-all", controller.ShutdownAllPlugins)
			}

			// 系统信息
			system := authApi.Group("/system")
			{
				// 获取系统信息
				system.GET("/info", controller.GetSystemInfo)
				// 获取服务器IP列表
				system.GET("/ips", controller.GetServerIPs)
			}

			// 仪表盘
			dashboard := authApi.Group("/dashboard")
			{
				// 获取仪表盘统计数据
				dashboard.GET("/stats", controller.GetDashboardStats)
			}
		}
	}

	httpPort := os.Getenv("HTTP_PORT")
	if httpPort == "" {
		httpPort = "8080"
	}
	engine.Run(":" + httpPort)
}
