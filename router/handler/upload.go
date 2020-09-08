package handler

import (
	"github.com/gin-gonic/gin"
	//"path/filepath"
)

func Post() gin.HandlerFunc {
	return func(c *gin.Context) {
		//post, err := c.FormFile("file")
		//if err != nil {
		//	errormap.SendHttpError(c, errormap.ErrorCodeFileUploadError, err)
		//	return
		//}
		//
		//dst := configer.GetString("file_cache_path")
		//token, err := uuid.NewRandom()
		//if err != nil {
		//	errormap.SendHttpError(c, errormap.ErrorCodeFileUploadError, err)
		//	return
		//}
		//
		//var relevantPath = filepath.Join(orm.DirUpload, token.String())
		//path := filepath.Join(dst, relevantPath)
		//err = c.SaveUploadedFile(post, path)
		//if err != nil {
		//	errormap.SendHttpError(c, errormap.ErrorCodeFileUploadError, err)
		//	return
		//}
		//
		//fmt.Println(post.Filename, post.Size, post.Header)
		//
		//file := orm.Attachment{
		//	Name:        post.Filename,
		//	Path:        relevantPath,
		//	Token:       token.String(),
		//	ContentType: post.Header["Content-Type"][0],
		//}
		//err = orm.Create(&file).Error
		//if err != nil {
		//	errormap.SendHttpError(c, errormap.ErrorCodeCreateObjectError, err, "file")
		//	return
		//}
		//
		//c.JSON(http.StatusOK, map[string]interface{}{
		//	"token": file.Token,
		//})
		return
	}
}

//func init() {
//	var fileCachePath = configer.GetString("file_cache_path")
//	//p := path.Join(fileCachePath, orm.DirUpload)
//	if err := os.MkdirAll(p, os.ModePerm); err != nil {
//		panic("cannot create templates directory.")
//	}
//}
