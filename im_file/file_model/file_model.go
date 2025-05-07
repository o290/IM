package file_model

import (
	"github.com/google/uuid"
	"server/common/models"
)

type FileModel struct {
	models.Model
	Uid      uuid.UUID `json:"uid"` //文件唯一id /api/file/{uuid}
	UserID   uint      `json:"userID"`
	FileName string    `json:"fileName"`
	Size     int64     `json:"size"` //文件大小
	Path     string    `json:"path"` //文件实际路径
	Hash     string    `json:"hash"` //文件hash
}

func (file *FileModel) WebPath() string {
	return "api/file/" + file.Uid.String()
}
