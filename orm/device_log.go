package orm

// 设备生产日志

import (
	"fmt"
	"github.com/SasukeBo/pmes-device-monitor/cache"
	"github.com/SasukeBo/pmes-device-monitor/orm/types"
	"github.com/jinzhu/copier"
	"github.com/jinzhu/gorm"
	"strconv"
	"time"
)

// 产量日志
// 一分钟记录一次设备的产量
type DeviceProduceLog struct {
	gorm.Model
	DeviceID     uint `gorm:"COMMENT:'设备ID';column:device_id;index;not null"`
	Total        int  `gorm:"COMMENT:'记录内产量';column:total"`
	NG           int  `gorm:"COMMENT:'记录内总不良';column:ng"`
	CurrentTotal int  `gorm:"COMMENT:'时刻总产量';column:current_total"`
	CurrentNG    int  `gorm:"COMMENT:'时刻总不良';column:current_ng"`
}

func (dpl *DeviceProduceLog) GetLast(mac string) error {
	var key = fmt.Sprintf("%s-last-dpl", mac)
	v := cache.Get(key)
	if v != nil {
		log, ok := v.(*DeviceProduceLog)
		if ok {
			if err := copier.Copy(dpl, log); err == nil {
				return nil
			}
		}
	}

	query := Model(dpl).Joins("JOIN devices ON devices.id = device_produce_logs.device_id")
	if err := query.Where("devices.mac = ?", mac).Order("device_produce_logs.created_at desc").First(dpl).Error; err != nil {
		return err
	}

	cache.Put(key, dpl)
	return nil
}

func (dpl *DeviceProduceLog) Record(device *Device, ct, cn int) error {
	var key = fmt.Sprintf("%s-last-dpl", device.Mac)
	dpl.DeviceID = device.ID
	dpl.CurrentNG = cn
	dpl.CurrentTotal = ct

	var last DeviceProduceLog
	if err := last.GetLast(device.Mac); err != nil {
		dpl.NG = cn
		dpl.Total = ct
	} else {
		now := time.Now()
		if now.Sub(last.CreatedAt) < time.Minute {
			return nil
		}
		if ct >= last.CurrentTotal && cn >= last.CurrentNG { // 数量非骤减，则表示PLC数据未重置
			dpl.Total = ct - last.CurrentTotal
			dpl.NG = cn - last.CurrentNG
		} else { // 否则当前数量统计更新
			dpl.Total = ct
			dpl.NG = cn
		}
	}

	if err := Create(dpl).Error; err != nil {
		return err
	}
	cache.Put(key, dpl)
	return nil
}

const ErrorIdxsKey = "errors"

// 设备状态日志
// 每次发生故障时记录错误序号，状态改变后设置故障结束时间
type DeviceStatusLog struct {
	gorm.Model
	DeviceID  uint      `gorm:"COMMENT:'设备ID';column:device_id;index;not null"`
	ErrorIdxs types.Map `gorm:"COMMENT:'故障代码';column:error_idxs;type:JSON"`
	Status    int       `gorm:"COMMENT:'设备状态';default:0"`
	Duration  int       `gorm:"COMMENT:'持续时间';column:duration;default:0"`
}

func (dsl *DeviceStatusLog) GetErrorIdxs() []int {
	var idxs []int
	if v, ok := dsl.ErrorIdxs[ErrorIdxsKey]; ok {
		if vs, ok := v.([]interface{}); ok {
			for item, _ := range vs {
				idx, err := strconv.Atoi(fmt.Sprint(item))
				if err != nil {
					continue
				}
				idxs = append(idxs, idx)
			}
		}
	}

	return idxs
}

func (dsl *DeviceStatusLog) GetLastError(deviceID int) error {
	return Model(dsl).Where("device_id = ? AND status = ?", deviceID, DeviceStatusError).Order("created_at desc").First(dsl).Error
}

func (dsl *DeviceStatusLog) GetLast(mac string) error {
	var key = fmt.Sprintf("%s-last-dsl", mac)

	if v := cache.Get(key); v != nil {
		if log, ok := v.(*DeviceStatusLog); ok {
			if err := copier.Copy(dsl, log); err == nil {
				return nil
			}
		}
	}

	query := Model(dsl).Joins("JOIN devices ON devices.id = device_status_logs.device_id")
	if err := query.Where("devices.mac = ?", mac).Order("device_status_logs.created_at desc").First(dsl).Error; err != nil {
		return err
	}

	cache.Put(key, dsl)
	return nil
}

func (dsl *DeviceStatusLog) Record(device *Device, status int, errors []int) error {
	dsl.Status = status
	dsl.DeviceID = device.ID
	if status == DeviceStatusError {
		em := make(types.Map)
		em[ErrorIdxsKey] = errors
		dsl.ErrorIdxs = em
	}

	var last DeviceStatusLog
	if err := last.GetLast(device.Mac); err == nil {
		if last.Status == status { // 状态未改变，不进行任何操作
			return nil
		}
		var now = time.Now()
		last.Duration = int(now.Sub(last.CreatedAt) / time.Second)
		Save(&last)
	}

	if err := Create(dsl).Error; err != nil {
		return err
	}
	device.Status = status
	Save(&device)

	var key = fmt.Sprintf("%s-last-dsl", device.Mac)
	cache.Put(key, dsl)
	return nil
}
