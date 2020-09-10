package logic

import (
	"context"
	"fmt"
	"github.com/SasukeBo/pmes-device-monitor/api/v1/model"
	"github.com/SasukeBo/pmes-device-monitor/orm"
)

func HomeDeviceStatusCount(ctx context.Context) (*model.DashboardDeviceStatusResponse, error) {
	query := orm.Model(&orm.Device{}).Select("COUNT(id), status").Group("status")
	rows, err := query.Rows()
	if err != nil {
		return nil, err
	}

	var out model.DashboardDeviceStatusResponse
	for rows.Next() {
		var count, status int
		if err := rows.Scan(&count, &status); err != nil {
			continue
		}

		switch status {
		case orm.DeviceStatusStopped:
			out.Stopped = count
		case orm.DeviceStatusRunning:
			out.Running = count
		case orm.DeviceStatusError:
			out.Error = count
		case orm.DeviceStatusShutdown:
			out.Offline = count
		}
	}

	return &out, nil
}

func HomeRecentDevices(ctx context.Context, ids []int, limit int) ([]*model.DashboardDevice, error) {
	if l := len(ids); limit > l {
		var appendIDs []int
		sql := orm.Model(&orm.Device{})
		if l > 0 {
			sql = sql.Where("id NOT IN (?)", ids)
		}
		sql = sql.Order("updated_at DESC").Limit(limit - l)
		if err := sql.Pluck("id", &appendIDs).Error; err == nil {
			ids = append(ids, appendIDs...)
		}
	}

	return realtimeDeviceAnalyze(ids)
}

func Devices(ctx context.Context, search *string, status *model.DeviceStatus, page int, limit int) (*model.DeviceWrap, error) {
	sql := orm.Model(&orm.Device{})

	if search != nil {
		sql = sql.Where("number LIKE ?", fmt.Sprintf("%%%s%%", *search))
	}

	if status != nil {
		switch *status {
		case model.DeviceStatusError:
			sql = sql.Where("status = ?", orm.DeviceStatusError)
		case model.DeviceStatusRunning:
			sql = sql.Where("status = ?", orm.DeviceStatusRunning)
		case model.DeviceStatusStopped:
			sql = sql.Where("status = ?", orm.DeviceStatusStopped)
		case model.DeviceStatusOffline:
			sql = sql.Where("status = ?", orm.DeviceStatusShutdown)
		}
	}

	var devices []orm.Device
	if err := sql.Limit(limit).Offset((page - 1) * limit).Find(&devices).Error; err != nil {
		return nil, err
	}

	var outs []*model.ListDevice
	for _, d := range devices {
		var out = model.ListDevice{
			ID:      int(d.ID),
			Number:  d.Number,
			Status:  d.GetStatusString(),
			Address: d.Address,
		}

		var dt orm.DeviceType
		if err := orm.Model(&dt).Where("id = ?", d.DeviceTypeID).First(&dt).Error; err == nil {
			out.DeviceType = dt.Name
		}

		sqlA := orm.Model(&orm.DeviceStatusLog{}).Select("SUM(duration), status").Group("status")
		sqlA = sqlA.Where("device_id = ? AND status != ?", d.ID, orm.DeviceStatusShutdown).Order("status ASC")
		if rows, err := sqlA.Rows(); err == nil {
			var running, sum, duration, status int
			for rows.Next() {
				if err := rows.Scan(&duration, &status); err != nil {
					continue
				}

				if status == orm.DeviceStatusRunning {
					running = duration
				}

				sum = sum + duration
			}
			if sum > 0 {
				out.Activation = float64(running) / float64(sum)
			}
		}

		sqlB := orm.Model(&orm.DeviceProduceLog{}).Select("SUM(total), SUM(ng)").Where("device_id = ?", d.ID)
		var result struct {
			Total int
			Ng    int
		}
		if err := sqlB.Scan(&result).Error; err == nil {
			if result.Total > 0 {
				out.Yield = float64(result.Ng) / float64(result.Total)
			}
		}

		outs = append(outs, &out)
	}

	var total int
	if err := sql.Count(&total).Error; err != nil {
		return nil, nil
	}

	return &model.DeviceWrap{
		Devices: outs,
		Total:   total,
	}, nil
}
