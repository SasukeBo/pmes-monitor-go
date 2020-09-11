package handler

import (
	"github.com/SasukeBo/log"
	"github.com/SasukeBo/pmes-device-monitor/websocket"
	"github.com/gin-gonic/gin"
)

func Websocket() gin.HandlerFunc {
	return func(c *gin.Context) {
		conn, err := websocket.NewWsConn(c)
		if err != nil {
			log.Errorln(err)
			return
		}

		defer conn.Close()
		for { // 循环等待客户端发来的订阅、退订消息
			if conn.Receive() {
				return
			}
		}
	}
}
