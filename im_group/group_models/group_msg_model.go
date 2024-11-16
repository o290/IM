package models

import (
	"server/common/models"
	"server/common/models/ctype"
)

// 群聊表
type GroupMsgModel struct {
	models.Model
	GroupID    uint             `json:"groupID"`
	GroupModel GroupModel       `gorm:"foreignKey:GroupID" json:"-"`
	SendUserID uint             `json:"sendUserID"`
	MsgType    int8             `json:"msgType"`                   //消息类型，1：文本2：图片3：视频4：文件5：语音6：语言通话7：视频通话8：撤回消息9：回复消息10：引用消息11:@消息
	MsgPreview string           `gorm:"size:64" json:"msgPreview"` //消息预览
	Msg        ctype.Msg        `json:"msg"`                       //消息内容
	SystemMsg  *ctype.SystemMsg `json:"systemMsg"`
}
