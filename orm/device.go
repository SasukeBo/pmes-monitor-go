package orm

// 监测设备

import (
	"github.com/SasukeBo/pmes-device-monitor/orm/types"
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
	Number       string `gorm:"COMMENT:'设备编号';not null"`
	DeviceTypeID int    `gorm:"COMMENT:'机种';not null"`
	Mac          string `gorm:"COMMENT:'MAC地址';column:mac;not null"`
	Address      string `gorm:"COMMENT:'物理地址';"`
	Status       int    `gorm:"COMMENT:'设备状态';default:0"`
	//UserID       int    `gorm:"COMMENT:'创建人';column:user_id;not null"`
}

type DeviceType struct {
	gorm.Model
	Name string `gorm:"COMMENT:'机种名称';column:name;not null"`
	//UserID      int    `gorm:"COMMENT:'创建人';column:user_id;not null"`
	ErrorCodeID int `gorm:"COMMENT:'故障代码';column:error_code_id"`
}

type DeviceErrorCode struct {
	gorm.Model
	Errors types.Map `gorm:"COMMENT:'故障代码中的位置';type:JSON;not null"`
}
