package main

import (
	_ "github.com/SasukeBo/pmes-device-monitor/orm"
	"github.com/SasukeBo/pmes-device-monitor/router"
)

func main() {
	router.Start()
}
