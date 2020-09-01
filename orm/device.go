package orm

// 监测设备

import (
	"github.com/jinzhu/gorm"
)

const (
	DeviceStatusStopped int = iota
	DeviceStatusRunning
	DeviceStatusError
	DeviceStatusShutdown
)

type Device struct {
	gorm.Model
	UUID   string `gorm:"column:uuid;unique_index;not null"`
	Name   string `gorm:"COMMENT:'机种';column:name;not null"`
	IP     string `gorm:"column:ip"`
	Number string `gorm:"COMMENT:'编号'"`
	Status int    `gorm:"COMMENT:'设备状态';default:0"`
}
