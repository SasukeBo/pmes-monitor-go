package logic

import (
	"context"
	"github.com/SasukeBo/configer"
	"github.com/SasukeBo/pmes-device-monitor/errormap"
	"github.com/SasukeBo/pmes-device-monitor/orm"
	"github.com/jinzhu/gorm"
	"github.com/tealeg/xlsx/v3"
	"io/ioutil"
	"path/filepath"
)

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
	var index int
	err = sheet.ForEachRow(func(r *xlsx.Row) error {
		var dErr orm.DeviceError
		if content := r.GetCell(0).Value; content != "" {
			dErr.Message = content
			dErr.Index = index
			dErr.DeviceID = uint(deviceID)
			if err := tx.Create(&dErr).Error; err != nil {
				return err
			}
			index++
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
