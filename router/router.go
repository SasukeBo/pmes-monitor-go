package router

import (
	"fmt"
	"github.com/SasukeBo/configer"
	"github.com/SasukeBo/log"
	"github.com/SasukeBo/pmes-device-monitor/router/handler"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"time"
)

func Start() {
	r := gin.Default()
	r.Use(gin.Recovery())
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://localhost:8080"},
		AllowMethods: []string{"POST"},
		AllowHeaders: []string{"Origin", "Content-Type"},
		MaxAge:       12 * time.Hour,
	}))

	//  API v1
	api1 := r.Group("/monitor/api", handler.HttpRequestLogger(), handler.InjectGinContext())
	{
		api1.POST("v1/admin", handler.API1Admin())
		api1.POST("v1", handler.API1())
		api1.POST("v1/import_error_codes", handler.ImportCodes()) // 上传错误代码
	}

	var port = configer.GetString("service_port")
	log.Info("Start service on [%s] mode", configer.GetEnv("env"))
	log.Info("HTTP service listening on %s", port)
	r.Run(fmt.Sprintf(":%s", port))
}
