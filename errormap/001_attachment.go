package errormap

import "net/http"

// ErrorCode E0001

const (
	// 400
	ErrorCodeFileUploadError    = "E0001S0400N0001"
	ErrorCodeFileExtensionError = "E0004S0400N0002"
	// 500
	ErrorCodeFileOpenFailedError = "E0004S0500N0001"
)

func init() {
	register(ErrorCodeFileUploadError, http.StatusBadRequest, langMap{
		ZH_CN: "对不起，文件上传失败，发生了一些错误。",
		EN:    "Sorry, failed to upload the file with some errors.",
	})
	register(ErrorCodeFileExtensionError, http.StatusBadRequest, langMap{
		ZH_CN: "对不起，文件格式不正确，需要{{.field_1}}文件。",
		EN:    "Sorry, the file extension is wrong, {{.field_1}} file in need.",
	})
	register(ErrorCodeFileOpenFailedError, http.StatusBadRequest, langMap{
		ZH_CN: "对不起，打开文件时发生了一些错误。",
		EN:    "Sorry, failed to open the file with some errors.",
	})
}
