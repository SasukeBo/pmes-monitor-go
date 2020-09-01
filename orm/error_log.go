package orm

// 故障日志

import (
	"github.com/jinzhu/gorm"
	"time"
)

type DeviceErrorLog struct {
	gorm.Model
	DeviceID  uint      `gorm:"COMMENT:'设备ID';column:device_id;index;not null"`
	ErrorCode string    `gorm:"COMMENT:'故障代码';column:error_code"`
	EndTime   time.Time `gorm:"COMMENT:'故障解决时间';column:end_time"`
}
