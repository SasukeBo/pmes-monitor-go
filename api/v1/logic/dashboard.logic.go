package logic

import (
	"context"
	"fmt"
	"github.com/SasukeBo/pmes-device-monitor/api/v1/model"
	"github.com/SasukeBo/pmes-device-monitor/orm"
	"github.com/jinzhu/copier"
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

func getToday() time.Time {
	var now = time.Now()
	return time.Date(now.Year(), now.Month(), now.Day(), -8, 0, 0, 0, time.UTC)
}

func realtimeDeviceAnalyze(ids []int) ([]*model.DashboardDevice, error) {
	var outs []*model.DashboardDevice
	var today = getToday()
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

		var tx = orm.Begin()
		var lastPLog orm.DeviceProduceLog
		tx.Model(&lastPLog).Where("device_id = ?", id).Last(&lastPLog)
		var lastSLog orm.DeviceStatusLog
		tx.Model(&lastSLog).Where("device_id = ?", id).Last(&lastSLog)
		out.LastProduceLogID = int(lastPLog.ID)
		out.LastStatusLogID = int(lastSLog.ID)
		out.LastStatusTime = lastSLog.CreatedAt

		// 统计产量
		tx.Model(orm.DeviceProduceLog{}).Where(
			"device_id = ? AND created_at > ?", id, today,
		).Select("SUM(total) as total, SUM(ng) as ng").Scan(&out)

		// 统计时间占比
		var durations = []int{0, 0, 0, 0}
		rows, err := tx.Model(orm.DeviceStatusLog{}).Where(
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
		tx.Commit()

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

func DashboardDeviceFresh(ctx context.Context, id int, pid int, sid int) (*model.DashboardDeviceFreshResponse, error) {
	var device orm.Device
	if err := device.Get(id); err != nil {
		return nil, err
	}
	errorCode := device.GetErrorCode()
	msgs := errorCode.GetErrors()

	var pLogs []orm.DeviceProduceLog
	orm.Model(&orm.DeviceProduceLog{}).Where("device_id = ? AND id > ?", id, pid).Find(&pLogs)
	var sLogs []orm.DeviceStatusLog
	orm.Model(&orm.DeviceStatusLog{}).Where("device_id = ? AND id > ?", id, sid).Find(&sLogs)

	var pOuts []*model.DeviceProduceLog
	for _, p := range pLogs {
		var out model.DeviceProduceLog
		if err := copier.Copy(&out, &p); err == nil {
			pOuts = append(pOuts, &out)
		}
	}

	var sOuts []*model.DeviceStatusLog
	for _, s := range sLogs {
		var out model.DeviceStatusLog
		if err := copier.Copy(&out, &s); err == nil {
			out.Status = s.GetStatusString()
			if s.Status == orm.DeviceStatusError {
				idxs := s.GetErrorIdxs()
				var messages []string
				for _, idx := range idxs {
					if idx < len(msgs) {
						messages = append(messages, msgs[idx])
					}
				}
				out.Messages = messages
			}
			sOuts = append(sOuts, &out)
		}
	}

	return &model.DashboardDeviceFreshResponse{
		ProduceLogs: pOuts,
		StatusLogs:  sOuts,
	}, nil
}

func DashboardOverviewAnalyze(ctx context.Context, id int) (*model.DashboardOverviewAnalyzeResponse, error) {
	var ds orm.Dashboard
	if err := ds.Get(id); err != nil {
		return nil, err
	}

	var out model.DashboardOverviewAnalyzeResponse
	deviceIDs := ds.GetDeviceIDs()
	today := getToday()
	orm.Model(orm.DeviceProduceLog{}).Where(
		"device_id in (?) AND created_at > ?", deviceIDs, today,
	).Select("SUM(total) as total, SUM(ng) as ng").Scan(&out)

	var durations = []int{0, 0, 0, 0}
	query := orm.Model(orm.DeviceStatusLog{}).Where("device_id in (?) AND created_at > ?", deviceIDs, today)
	query = query.Select("SUM(duration), device_status_logs.status").Group("device_status_logs.status")
	query = query.Order("device_status_logs.status")
	rows, err := query.Rows()
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
	power := durations[0] + durations[1] + durations[2]
	if power > 0 {
		out.Activation = float64(durations[1]) / float64(power)
	}

	return &out, nil
}

func DashboardDeviceStatus(ctx context.Context, id int) (*model.DashboardDeviceStatusResponse, error) {
	var ds orm.Dashboard
	if err := ds.Get(id); err != nil {
		return nil, err
	}

	var out model.DashboardDeviceStatusResponse
	deviceIDs := ds.GetDeviceIDs()

	query := orm.Model(&orm.Device{}).Where("id in (?)", deviceIDs)
	query = query.Select("COUNT(id), devices.status").Group("devices.status").Order("devices.status")
	rows, err := query.Rows()
	if err == nil {
		var count, status int
		for rows.Next() {
			if err := rows.Scan(&count, &status); err != nil {
				continue
			}
			switch status {
			case orm.DeviceStatusStopped:
				out.Stopped = count
			case orm.DeviceStatusRunning:
				out.Running = count
			case orm.DeviceStatusShutdown:
				out.Offline = count
			case orm.DeviceStatusError:
				out.Error = count
			}
		}
	}

	return &out, nil
}

func DashboardDeviceErrors(ctx context.Context, id int) (*model.DashboardDeviceErrorsResponse, error) {
	var ds orm.Dashboard
	if err := ds.Get(id); err != nil {
		return nil, err
	}

	deviceIDs := ds.GetDeviceIDs()
	query := orm.Model(&orm.DeviceStatusLog{}).Select("COUNT(device_status_logs.id), devices.number")
	query = query.Joins("JOIN devices ON device_status_logs.device_id = devices.id")
	query = query.Where("device_status_logs.device_id in (?) AND device_status_logs.status = ?", deviceIDs, orm.DeviceStatusError)
	query = query.Group("devices.number")
	rows, err := query.Rows()
	if err != nil {
		return nil, err
	}

	var category []string
	var data []int
	for rows.Next() {
		var c string
		var d int
		if err := rows.Scan(&d, &c); err != nil {
			continue
		}
		category = append(category, c)
		data = append(data, d)
	}

	return &model.DashboardDeviceErrorsResponse{
		Category: category,
		Data:     data,
	}, nil
}
