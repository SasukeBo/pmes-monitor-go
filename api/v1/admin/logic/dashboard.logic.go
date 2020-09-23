package logic

import (
	"context"
	"fmt"
	"github.com/SasukeBo/pmes-device-monitor/api/v1/admin/model"
	"github.com/SasukeBo/pmes-device-monitor/orm"
	"github.com/SasukeBo/pmes-device-monitor/orm/types"
	"github.com/jinzhu/copier"
	"github.com/jinzhu/gorm"
	"strconv"
)

func AdminCreateDashboard(ctx context.Context, name string, deviceIDs []int) (string, error) {
	var dashboard = orm.Dashboard{Name: name}
	var ids = make(types.Map)
	ids["deviceIDs"] = deviceIDs
	dashboard.DeviceIDs = ids
	if err := orm.Create(&dashboard).Error; err != nil {
		return "", err
	}

	return "ok", nil
}

func AdminDashboards(ctx context.Context, search *string, page int, limit int) (*model.DashboardWrap, error) {
	var query = orm.Model(&orm.Dashboard{})
	if search != nil {
		query = query.Where("name LIKE ?", fmt.Sprintf("%%%s%%", *search))
	}

	var ds []orm.Dashboard
	if err := query.Limit(limit).Offset((page - 1) * limit).Find(&ds).Error; err != nil {
		return nil, err
	}

	var count int
	if err := query.Count(&count).Error; err != nil {
		return nil, err
	}

	var outs []*model.Dashboard
	for _, d := range ds {
		var out model.Dashboard
		if value, ok := d.DeviceIDs["deviceIDs"]; ok {
			if ids, ok := value.([]interface{}); ok {
				var deviceIDs []int
				for _, v := range ids {
					id, err := strconv.Atoi(fmt.Sprint(v))
					if err != nil {
						continue
					}
					deviceIDs = append(deviceIDs, id)
				}
				out.DeviceIDs = deviceIDs
			}
		}
		if err := copier.Copy(&out, &d); err != nil {
			continue
		}

		outs = append(outs, &out)
	}

	return &model.DashboardWrap{
		Total:      count,
		Dashboards: outs,
	}, nil
}

func AdminDashboardDelete(ctx context.Context, id int) (string, error) {
	var d orm.Dashboard
	if err := orm.Model(&d).Where("id =  ?", id).First(&d).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			return "ERROR", err
		}

		return "OK", nil
	}

	if err := orm.Delete(&d).Error; err != nil {
		return "ERROR", err
	}

	return "OK", nil
}

func AdminDashboard(ctx context.Context, id int) (*model.Dashboard, error) {
	var d orm.Dashboard
	if err := orm.Model(&d).Where("id = ?", id).First(&d).Error; err != nil {
		return nil, err
	}

	var out model.Dashboard
	out.ID = int(d.ID)
	out.CreatedAt = d.CreatedAt
	out.Name = d.Name
	if value, ok := d.DeviceIDs["deviceIDs"]; ok {
		if ids, ok := value.([]interface{}); ok {
			var deviceIDs []int
			for _, v := range ids {
				id, err := strconv.Atoi(fmt.Sprint(v))
				if err != nil {
					continue
				}
				deviceIDs = append(deviceIDs, id)
			}
			out.DeviceIDs = deviceIDs
		}
	}

	return &out, nil
}
