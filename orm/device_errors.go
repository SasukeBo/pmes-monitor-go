package orm

// 设备故障

import (
	"github.com/jinzhu/gorm"
)

type DeviceError struct {
	gorm.Model
	Index   int    `gorm:"COMMENT:'故障代码中的位置';not null"`
	Message string `gorm:"COMMENT:'错误信息';not null"`
}
