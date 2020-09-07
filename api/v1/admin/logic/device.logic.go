package logic

import (
	"context"
	"fmt"
	"github.com/SasukeBo/configer"
	"github.com/SasukeBo/log"
	"github.com/SasukeBo/pmes-device-monitor/api/v1/admin/model"
	"github.com/SasukeBo/pmes-device-monitor/errormap"
	"github.com/SasukeBo/pmes-device-monitor/orm"
	"github.com/SasukeBo/pmes-device-monitor/orm/types"
	"github.com/jinzhu/copier"
	"github.com/jinzhu/gorm"
	"github.com/tealeg/xlsx/v3"
	"io/ioutil"

	"path/filepath"
)

func AdminSaveErrorCode(ctx context.Context, id int, errors []string) (string, error) {
	var ec orm.DeviceErrorCode
	if err := orm.Model(&ec).Where("id = ?", id).First(&ec).Error; err != nil {
		return "", err
	}

	var errs = make(types.Map)
	errs["messages"] = errors
	ec.Errors = errs
	if err := orm.Save(&ec).Error; err != nil {
		return "", err
	}

	return "ok", nil
}

func AdminDeviceTypeAddErrorCode(ctx context.Context, deviceTypeID int, errors []string) (string, error) {
	var errs = make(types.Map)
	errs["messages"] = errors
	var errorCode = orm.DeviceErrorCode{
		Errors: errs,
	}
	tx := orm.Begin()
	if err := tx.Create(&errorCode).Error; err != nil {
		tx.Rollback()
		return "", err
	}
	var dt orm.DeviceType
	if err := orm.Model(&dt).Where("id = ?", deviceTypeID).First(&dt).Error; err != nil {
		tx.Rollback()
		return "", err
	}

	dt.ErrorCodeID = int(errorCode.ID)
	if err := tx.Save(&dt).Error; err != nil {
		tx.Rollback()
		return "", err
	}

	tx.Commit()
	return "ok", nil
}

func AdminDeviceTypeCreate(ctx context.Context, name string) (string, error) {
	//user := api.CurrentUser(ctx)
	//if user == nil {
	//	return "", errormap.SendGQLError(ctx, errormap.ErrorCodeUnauthenticated, nil)
	//}
	//if !user.IsAdmin {
	//	return "", errormap.SendGQLError(ctx, errormap.ErrorCodePermissionDeny, nil)
	//}

	var deviceType = orm.DeviceType{
		Name: name,
		//UserID: int(user.ID),
	}
	if err := orm.Create(&deviceType).Error; err != nil {
		return "", errormap.SendGQLError(ctx, errormap.ErrorCodeCreateObjectError, err, "device_type")
	}

	return "ok", nil
}

func AdminDeviceTypes(ctx context.Context, search *string, page int, limit int) (*model.DeviceTypeWrap, error) {
	//user := api.CurrentUser(ctx)
	//if user == nil {
	//	return nil, errormap.SendGQLError(ctx, errormap.ErrorCodeUnauthenticated, nil)
	//}
	//if !user.IsAdmin {
	//	return nil, errormap.SendGQLError(ctx, errormap.ErrorCodePermissionDeny, nil)
	//}

	query := orm.Model(&orm.DeviceType{})
	if search != nil {
		query = query.Where("name LIKE ?", fmt.Sprintf("%%%s%%", *search))
	}

	var types []orm.DeviceType
	if err := query.Limit(limit).Offset((page - 1) * limit).Find(&types).Error; err != nil {
		return nil, errormap.SendGQLError(ctx, errormap.ErrorCodeGetObjectFailed, err, "device_type")
	}

	var count int
	if err := query.Count(&count).Error; err != nil {
		return nil, errormap.SendGQLError(ctx, errormap.ErrorCodeGetObjectFailed, err, "device_type_count")
	}

	var outs []*model.DeviceType
	for _, t := range types {
		var out model.DeviceType
		if err := copier.Copy(&out, &t); err != nil {
			continue
		}
		outs = append(outs, &out)
	}

	return &model.DeviceTypeWrap{
		Total: count,
		Types: outs,
	}, nil
}

func ImportErrors(ctx context.Context, deviceID int, fileToken string) (string, error) {
	var file orm.Attachment
	if err := file.GetByToken(fileToken); err != nil {
		if err == gorm.ErrRecordNotFound {
			return "", errormap.SendGQLError(ctx, errormap.ErrorCodeObjectNotFound, err, "attachment")
		}
		return "", errormap.SendGQLError(ctx, errormap.ErrorCodeGetObjectFailed, err, "attachment")
	}

	baseDir := configer.GetString("file_cache_path")

	if filepath.Ext(file.Name) != ".xlsx" {
		return "", errormap.SendGQLError(ctx, errormap.ErrorCodeFileExtensionError, nil, ".xlsx")
	}

	content, err := ioutil.ReadFile(filepath.Join(baseDir, file.Path))
	if err != nil {
		return "", errormap.SendGQLError(ctx, errormap.ErrorCodeGetObjectFailed, err, "attachment")
	}
	xFile, err := xlsx.OpenBinary(content)
	if err != nil {
		return "", errormap.SendGQLError(ctx, errormap.ErrorCodeFileOpenFailedError, err)
	}

	sheet := xFile.Sheets[0]
	tx := orm.Begin()
	var dErr orm.DeviceErrorCode
	var messages []string
	err = sheet.ForEachRow(func(r *xlsx.Row) error {
		if content := r.GetCell(0).Value; content != "" {
			messages = append(messages, content)
			dErr.Errors = types.Map{"messages": messages}
			if err := tx.Create(&dErr).Error; err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		tx.Rollback()
		return "", errormap.SendGQLError(ctx, errormap.ErrorCodeFileOpenFailedError, err)
	}

	tx.Commit()
	return "success", nil
}

func LoadErrorCode(ctx context.Context, id int) *model.ErrorCode {
	var ec orm.DeviceErrorCode
	if err := orm.Model(&orm.DeviceErrorCode{}).Where("id = ?", id).First(&ec).Error; err != nil {
		return nil
	}
	var out model.ErrorCode
	if err := copier.Copy(&out, &ec); err != nil {
		return nil
	}
	errMap := ec.Errors
	if value, ok := errMap["messages"]; ok {
		log.Infoln(value)
		if errs, ok := value.([]interface{}); ok {
			var errStr []string
			for _, v := range errs {
				errStr = append(errStr, fmt.Sprint(v))
			}
			out.Errors = errStr
		}
	}
	return &out
}

func AdminDeviceType(ctx context.Context, id int) (*model.DeviceType, error) {
	var deviceType orm.DeviceType
	if err := orm.Model(&orm.DeviceType{}).Where("id = ?", id).First(&deviceType).Error; err != nil {
		return nil, errormap.SendGQLError(ctx, errormap.ErrorCodeGetObjectFailed, err, "device_type")
	}

	var out model.DeviceType
	if err := copier.Copy(&out, &deviceType); err != nil {
		return nil, errormap.SendGQLError(ctx, errormap.ErrorCodeGetObjectFailed, err, "device_type")
	}

	return &out, nil
}

func LoadDeviceType(ctx context.Context, id int) *model.DeviceType {
	var deviceType orm.DeviceType
	if err := orm.Model(&orm.DeviceType{}).Where("id = ?", id).First(&deviceType).Error; err != nil {
		return nil
	}

	var out model.DeviceType
	if err := copier.Copy(&out, &deviceType); err != nil {
		return nil
	}

	return &out
}

func AdminCreateDevices(ctx context.Context, input model.CreateDeviceInput) (string, error) {
	var deviceType orm.DeviceType
	if err := orm.Model(&orm.DeviceType{}).Where("id = ?", input.DeviceTypeID).First(&deviceType).Error; err != nil {
		return "", nil
	}

	tx := orm.Begin()
	for _, di := range input.DeviceInputs {
		var device = orm.Device{
			Number:       di.Number,
			DeviceTypeID: int(deviceType.ID),
			Mac:          di.Mac,
			Status:       orm.DeviceStatusStopped,
		}
		if di.Address != nil {
			device.Address = *di.Address
		}

		if err := tx.Save(&device).Error; err != nil {
			tx.Rollback()
			return "", err
		}
	}

	tx.Commit()
	return "ok", nil
}

func AdminDevices(ctx context.Context, search *string, page int, limit int) (*model.DeviceWrap, error) {
	var devices []orm.Device
	var query = orm.Model(&orm.Device{})
	if search != nil {
		query = query.Where("number LIKE ?", fmt.Sprintf("%%%s%%", *search))
	}

	var outs []*model.Device
	if err := query.Limit(limit).Offset((page - 1) * limit).Find(&devices).Error; err != nil {
		return nil, err
	}
	for _, d := range devices {
		var out model.Device
		if err := copier.Copy(&out, &d); err != nil {
			continue
		}
		outs = append(outs, &out)
	}

	var count int
	if err := query.Count(&count).Error; err != nil {
		return nil, err
	}

	return &model.DeviceWrap{
		Total:   count,
		Devices: outs,
	}, nil
}

func LoadDevices(ctx context.Context, ids []int) []*model.Device {
	var devices []orm.Device
	if err := orm.Model(&orm.Device{}).Where("id in (?)", ids).Find(&devices).Error; err != nil {
		return nil
	}

	var outs []*model.Device
	for _, d := range devices {
		var out model.Device
		if err := copier.Copy(&out, &d); err != nil {
			continue
		}
		outs = append(outs, &out)
	}

	return outs
}
