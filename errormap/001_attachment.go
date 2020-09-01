package errormap

import "net/http"

// ErrorCode E0001

const (
	// 400
	ErrorCodeFileUploadError = "E0001S0400N0001"
)

func init() {
	register(ErrorCodeFileUploadError, http.StatusBadRequest, langMap{
		ZH_CN: "对不起，文件上传失败，发生了一些错误。",
		EN:    "Sorry, failed to upload the file with some errors.",
	})
}
