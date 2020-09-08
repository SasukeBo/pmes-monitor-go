package orm

// 监测设备

import (
	"fmt"
	"github.com/SasukeBo/pmes-device-monitor/cache"
	"github.com/SasukeBo/pmes-device-monitor/orm/types"
	"github.com/jinzhu/copier"
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

func (d *Device) GetErrorCode() *DeviceErrorCode {
	var dt DeviceType
	if err := Model(&dt).Where("id = ?", d.DeviceTypeID).First(&dt).Error; err != nil {
		return nil
	}

	var dec DeviceErrorCode
	if err := Model(&dec).Where("id = ?", dt.ErrorCodeID).First(&dec).Error; err != nil {
		return nil
	}

	return &dec
}

func (d *Device) GetStatusString() string {
	switch d.Status {
	case DeviceStatusError:
		return "error"
	case DeviceStatusRunning:
		return "running"
	case DeviceStatusStopped:
		return "stopped"
	case DeviceStatusShutdown:
		return "offline"
	default:
		return "offline"
	}
}

func (d *Device) Get(id int) error {
	return Model(d).Where("id = ?", id).First(d).Error
}

func (d *Device) GetByMAC(mac string) error {
	cacheKey := fmt.Sprintf("%s-device", mac)
	v := cache.Get(cacheKey)
	if v != nil {
		device, ok := v.(*Device)
		if ok {
			if err := copier.Copy(d, device); err == nil {
				return nil
			}
		}
	}

	if err := Model(d).Where("mac = ?", mac).First(d).Error; err != nil {
		return err
	}

	cache.Put(cacheKey, d)
	return nil
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

const DeviceErrorCodeErrorsMapKey = "messages"

func (dec *DeviceErrorCode) GetErrors() []string {
	var errs []string
	if v, ok := dec.Errors[DeviceErrorCodeErrorsMapKey]; ok {
		if vs, ok := v.([]interface{}); ok {
			for _, item := range vs {
				errs = append(errs, fmt.Sprint(item))
			}
		}
	}

	return errs
}
