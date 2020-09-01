package orm

// 生产日志

import (
	"github.com/jinzhu/gorm"
)

type ProduceLog struct {
	gorm.Model
	DeviceID    uint `gorm:"COMMENT:'设备ID';column:device_id;index;not null"`
	TotalAmount int  `gorm:"COMMENT:'产量';column:total_amount"`
	OKAmount    int  `gorm:"COMMENT:'良品数';column:ok_amount"`
}
