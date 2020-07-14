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

	log.Info("start service on [%s] mode", configer.GetEnv("env"))
	r.Run(fmt.Sprintf(":%s", configer.GetString("service_port")))
}
