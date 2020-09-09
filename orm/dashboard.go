package orm

// 故障日志

import (
	"fmt"
	"github.com/SasukeBo/pmes-device-monitor/orm/types"
	"github.com/jinzhu/gorm"
	"strconv"
)

const DashboardDeviceIDsMapKey = "deviceIDs"

type Dashboard struct {
	gorm.Model
	Name      string    `gorm:"COMMENT:'看板名称';column:name;index;not null"`
	DeviceIDs types.Map `gorm:"COMMENT:'关联设备';column:device_ids;type:JSON"`
}

func (d *Dashboard) Get(id int) error {
	return Model(d).Where("id = ?", id).First(d).Error
}

func (d *Dashboard) GetDeviceIDs() []int {
	var idxs []int
	if v, ok := d.DeviceIDs[DashboardDeviceIDsMapKey]; ok {
		if vs, ok := v.([]interface{}); ok {
			for _, item := range vs {
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
