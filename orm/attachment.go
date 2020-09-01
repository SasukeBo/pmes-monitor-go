package orm

// 故障日志

import (
	"github.com/jinzhu/gorm"
)

const (
	DirUpload = "post"
)

type Attachment struct {
	gorm.Model
	Name        string `gorm:"COMMENT:'附件名称';not null"`
	Path        string `gorm:"COMMENT:'相对路径';not null"`
	Token       string `gorm:"COMMENT:'文件Token';not null"`
	ContentType string `gorm:"COMMENT:'文件类型';not null"`
}
