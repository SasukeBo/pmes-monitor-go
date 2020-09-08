package logic

import (
	"context"
	"fmt"
	"github.com/SasukeBo/pmes-device-monitor/api/v1/model"
	"github.com/SasukeBo/pmes-device-monitor/orm"
	"strconv"
	"time"
)

func Dashboards(ctx context.Context, search *string, limit int, page int) (*model.DashboardWrap, error) {
	var query = orm.Model(&orm.Dashboard{})

	if search != nil {
		query = query.Where("name LIKE ?", fmt.Sprintf("%%%s%%", *search))
	}

	var dashboards []orm.Dashboard
	if err := query.Limit(limit).Offset((page - 1) * limit).Find(&dashboards).Error; err != nil {
		return nil, err
	}

	var total int
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	var outs []*model.Dashboard
	for _, d := range dashboards {
		var out = model.Dashboard{
			ID:   int(d.ID),
			Name: d.Name,
		}
		outs = append(outs, &out)

		var deviceIDs []int
		if v, ok := d.DeviceIDs[orm.DashboardDeviceIDsMapKey]; ok {
			if items, ok := v.([]interface{}); ok {
				for _, item := range items {
					id, err := strconv.Atoi(fmt.Sprint(item))
					if err != nil {
						continue
					}
					deviceIDs = append(deviceIDs, id)
				}
			}
		}

		out.DeviceTotal = len(deviceIDs)
		if out.DeviceTotal == 0 {
			continue
		}

		sql := orm.Model(&orm.Device{}).Select("count(id), devices.status").Where("id in (?)", deviceIDs)
		rows, err := sql.Group("devices.status").Rows()
		if err != nil {
			continue
		}

		for rows.Next() {
			var count, status int
			if err := rows.Scan(&count, &status); err != nil {
				continue
			}

			if status == orm.DeviceStatusRunning {
				out.RunningTotal = count
			}
			if status == orm.DeviceStatusError {
				out.ErrorTotal = count
			}
		}
	}

	return &model.DashboardWrap{
		Total:      total,
		Dashboards: outs,
	}, nil
}

func DashboardDevices(ctx context.Context, id int) ([]*model.DashboardDevice, error) {
	var board orm.Dashboard
	if err := board.Get(id); err != nil {
		return nil, err
	}

	var deviceIDs []int
	if v, ok := board.DeviceIDs[orm.DashboardDeviceIDsMapKey]; ok {
		if vs, ok := v.([]interface{}); ok {
			for _, item := range vs {
				id, err := strconv.Atoi(fmt.Sprint(item))
				if err != nil {
					continue
				}
				deviceIDs = append(deviceIDs, id)
			}
		}
	}

	return realtimeDeviceAnalyze(deviceIDs)
}

func realtimeDeviceAnalyze(ids []int) ([]*model.DashboardDevice, error) {
	var outs []*model.DashboardDevice
	var now = time.Now()
	var today = time.Date(now.Year(), now.Month(), now.Day(), -8, 0, 0, 0, time.UTC)
	for _, id := range ids {
		var device orm.Device
		if err := device.Get(id); err != nil {
			continue
		}

		var out = model.DashboardDevice{
			ID:     int(device.ID),
			Number: device.Number,
			Status: device.GetStatusString(),
		}

		// 统计产量
		orm.Model(orm.DeviceProduceLog{}).Where(
			"device_id = ? AND created_at > ?", id, today,
		).Select("SUM(total) as total, SUM(ng) as ng").Scan(&out)

		var durations = []int{0, 0, 0, 0}
		rows, err := orm.Model(orm.DeviceStatusLog{}).Where(
			"device_id = ? AND created_at > ?", id, today,
		).Select("SUM(duration), device_status_logs.status").Group("device_status_logs.status").Rows()
		if err == nil {
			var sum, status int
			for rows.Next() {
				if err := rows.Scan(&sum, &status); err != nil {
					continue
				}
				if status < 4 {
					durations[status] = sum
				}
			}
		}
		out.Durations = durations

		var messages []string
		if device.Status == orm.DeviceStatusError {
			var errorCode = device.GetErrorCode()
			msgs := errorCode.GetErrors()

			var errLog orm.DeviceStatusLog
			var idxs []int
			if err := errLog.GetLastError(id); err == nil {
				idxs = errLog.GetErrorIdxs()
			}

			for _, idx := range idxs {
				if idx < len(msgs) {
					messages = append(messages, msgs[idx])
				}
			}
		}
		out.Errors = messages
		outs = append(outs, &out)
	}

	return outs, nil
}
