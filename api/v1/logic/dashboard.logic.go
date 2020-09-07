package logic

import (
	"context"
	"fmt"
	"github.com/SasukeBo/pmes-device-monitor/api/v1/model"
	"github.com/SasukeBo/pmes-device-monitor/orm"
	"strconv"
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
		if v, ok := d.DeviceIDs["deviceIDs"]; ok {
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
