package main

import (
	//_ "github.com/SasukeBo/pmes-device-monitor/orm"
	"github.com/SasukeBo/pmes-device-monitor/mqtt"
	"github.com/SasukeBo/pmes-device-monitor/router"
)

func main() {
	go mqtt.Serve()
	router.Start()
}
