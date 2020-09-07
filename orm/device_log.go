package orm

// 设备生产日志

import (
	"github.com/jinzhu/gorm"
	"time"
)

// 产量日志
// 五分钟记录一次设备的产量
type DeviceProduceLog struct {
	gorm.Model
	DeviceID    uint `gorm:"COMMENT:'设备ID';column:device_id;index;not null"`
	TotalAmount int  `gorm:"COMMENT:'产量';column:total_amount"`
	OKAmount    int  `gorm:"COMMENT:'良品数';column:ok_amount"`
}

type DeviceErrorLog struct {
	gorm.Model
	DeviceID  uint      `gorm:"COMMENT:'设备ID';column:device_id;index;not null"`
	ErrorCode string    `gorm:"COMMENT:'故障代码';column:error_code"`
	EndTime   time.Time `gorm:"COMMENT:'故障解决时间';column:end_time"`
}

type DeviceStatusLog struct {
	gorm.Model
	Status   int `gorm:"COMMENT:'设备状态';default:0"`
	Duration int `gorm:"COMMENT:'持续时间';column:end_time"`
}
