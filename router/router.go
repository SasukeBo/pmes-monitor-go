package router

import (
	"fmt"
	"github.com/SasukeBo/configer"
	"github.com/SasukeBo/log"
	"github.com/SasukeBo/pmes-device-monitor/router/handler"
	"github.com/gin-gonic/gin"
)

func Start() {
	r := gin.Default()
	r.Use(gin.Recovery())

	r.POST("/monitor/file/posts", handler.Post()) // 上传文件

	//  API v1
	api1 := r.Group("/monitor/api", handler.HttpRequestLogger(), handler.InjectGinContext())
	{
		api1.POST("v1", handler.API1())
	}

	var port = configer.GetString("service_port")
	log.Info("Start service on [%s] mode", configer.GetEnv("env"))
	log.Info("HTTP service listening on %s", port)
	r.Run(fmt.Sprintf(":%s", port))
}
