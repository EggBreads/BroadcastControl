package main

import (
	"github.com/catenoid-company/wrController/config"
	"github.com/catenoid-company/wrController/controller"
	"github.com/catenoid-company/wrController/docs"
	"github.com/catenoid-company/wrController/logger"
	"github.com/catenoid-company/wrController/middlewares"
	"github.com/catenoid-company/wrController/utils"
	"github.com/catenoid-company/wrController/wrService/monitoring"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"os"
	"runtime"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	utils.RedisClient = setConfigAndRedis()

	defer utils.RedisClient.Close()

	runtime.GOMAXPROCS(runtime.NumCPU())

	if config.WrConfig.JanusHealthCheckIsUse == "true"{
		go monitoring.AgentManageServers()
	}
	// Gin Framework Adjust
	g := setWrControllerServer()

	logger.Info("main", "================ START WrController ================")

	// Start WrController
	err := g.Run(":8888")

	if err != nil{
		logger.Error("[Error] Start Error %s", err.Error())
	}
}

func setConfigAndRedis() *redis.Client {
	//Bind to configuration to environment
	config.WrConfig = config.InitConfig()

	// Logger Set
	logger.Init()
	logger.MonitoringLoggerInit()
	logger.WithStruct(config.WrConfig).Info("main", "WebrtcController Config")

	// Connect to Redis Sentinel by SentinelFailover
	redisClient, err := utils.ConnectSentinel(config.WrConfig)

	if err != nil{
		logger.Error("main","Fail to connect redis sentinel : %s", err)
		os.Exit(1)
		return nil
	}

	utils.MysqlGorm = utils.ConnMysql()

	return redisClient
}

// @title WebRtc Controller Swagger API
// @version 1.0
// @description Webrtc Controller Api

// @contact.name API Support
// @contact.email deuksoo.mun@catenoid.net
func setWrControllerServer() *gin.Engine {
	docs.SwaggerInfo.Schemes = []string{
		"http", "https",
	}
	docs.SwaggerInfo.Host = config.WrConfig.WrcBaseHost
	docs.SwaggerInfo.BasePath = "/" + config.V1
	// Gin Framework Adjust
	g := gin.Default()

	g.Use(middlewares.AccessLogMiddleware())

	// 고유 계정
	account := gin.Accounts{
		config.WrConfig.AuthUser: config.WrConfig.AuthPass,
	}

	// Webrtc BroadCast Router Group
	routerGroup := g.Group(config.V1, gin.BasicAuth(account))
	{
		// BroadCastHandlers properties insert to redisClient and environment configuration
		liveHandlers := controller.LiveApiHandlers{}

		routerGroup.POST(config.CHANNEL,liveHandlers.PrepareBroadCast)
		routerGroup.DELETE(config.CHANNEL,liveHandlers.CancelPrepareBroadCast)

		routerGroup.POST(config.PUBLISH,liveHandlers.StartBroadCasting)

		routerGroup.POST(config.UNPUBLISH,liveHandlers.CloseBroadCasting)
	}

	// Webrtc Monitoring Router Group
	{
		// MonitoringHandlers properties insert to redisClient and environment configuration
		monitoringHandler := controller.MonitoringHandler{}

		routerGroup.GET(config.JANUSINFO,monitoringHandler.GetJanusServersInfo)
		routerGroup.POST(config.JANUSINFO,monitoringHandler.SaveJanusMonitoringHandler)

		routerGroup.GET(config.NGINXINFO,monitoringHandler.GetNginxServersInfo)
		routerGroup.POST(config.NGINXINFO,monitoringHandler.SaveNginxMonitoringHandler)

		routerGroup.POST(config.HEALTH, monitoringHandler.AgentHealthCheck)
	}

	routerGroup.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return g
}
