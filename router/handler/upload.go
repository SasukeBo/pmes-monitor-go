package handler

import (
	"bytes"
	"github.com/SasukeBo/pmes-device-monitor/errormap"
	"github.com/gin-gonic/gin"
	"github.com/tealeg/xlsx/v3"
	"io"
	"net/http"
)

func ImportCodes() gin.HandlerFunc {
	return func(c *gin.Context) {
		file, _, err := c.Request.FormFile("file")
		if err != nil {
			errormap.SendHttpError(c, errormap.ErrorCodeFileUploadError, err)
			return
		}
		defer file.Close()

		buf := bytes.NewBuffer(nil)
		if _, err := io.Copy(buf, file); err != nil {
			errormap.SendHttpError(c, errormap.ErrorCodeFileUploadError, err)
			return
		}

		xf, err := xlsx.OpenBinary(buf.Bytes())
		if err != nil {
			errormap.SendHttpError(c, errormap.ErrorCodeFileExtensionError, err, ".xlsx")
			return
		}

		sheet := xf.Sheets[0]
		var messages []string
		sheet.ForEachRow(func(r *xlsx.Row) error {
			messages = append(messages, r.GetCell(0).String())
			return nil
		})

		c.JSON(http.StatusOK, messages)
	}
}

//func init() {
//	var fileCachePath = configer.GetString("file_cache_path")
//	//p := path.Join(fileCachePath, orm.DirUpload)
//	if err := os.MkdirAll(p, os.ModePerm); err != nil {
//		panic("cannot create templates directory.")
//	}
//}
