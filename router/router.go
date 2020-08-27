package router

import (
	"fmt"
	"github.com/SasukeBo/configer"
	"github.com/SasukeBo/log"
	"github.com/gin-gonic/gin"
)

func Start() {
	r := gin.Default()
	r.Use(gin.Recovery())

	var port = configer.GetString("service_port")
	log.Info("Start service on [%s] mode", configer.GetEnv("env"))
	log.Info("HTTP service listening on %s", port)
	r.Run(fmt.Sprintf(":%s", port))
}
