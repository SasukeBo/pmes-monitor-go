package orm

// 设备故障

import (
	"github.com/jinzhu/gorm"
)

type DeviceError struct {
	gorm.Model
	DeviceID uint   `gorm:"COMMENT:'设备ID';not null;index"`
	Index    int    `gorm:"COMMENT:'故障代码中的位置';column:idx;not null"`
	Message  string `gorm:"COMMENT:'错误信息';not null"`
}
