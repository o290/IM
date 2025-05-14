package chat_models

import (
	"github.com/google/uuid"
	"server/common/models"
	"server/common/models/ctype"
)

type OfflineMsgModel struct {
	models.Model
	SendUserID uint             `json:"sendUserID"`
	RevUserID  uint             `json:"revUserID"`
	MsgUUID    string           `json:"msgUUID"`
	MsgType    ctype.MsgType    `json:"msgType"`
	Msg        ctype.Msg        `json:"msg"`
	SystemMsg  *ctype.SystemMsg `json:"systemMsg"`  //系统消息
	MsgPreview string           `json:"msgPreview"` //消息预览
}

// GenerateMsgUUID 生成唯一的消息 UUID
func (o *OfflineMsgModel) GenerateMsgUUID() {
	o.MsgUUID = uuid.New().String()
}
