package group_models

import (
	"server/common/models"
	"server/common/models/ctype"
)

// 群聊表
type GroupMsgModel struct {
	models.Model
	GroupID          uint              `json:"groupID"`
	GroupModel       GroupModel        `gorm:"foreignKey:GroupID" json:"-"`
	SendUserID       uint              `json:"sendUserID"`
	GroupMemberID    uint              `json:"groupMemberID"`
	GroupMemberModel *GroupMemberModel `gorm:"foreignKey:GroupMemberID" json:"-"` //对应的群成员
	MsgType          ctype.MsgType     `json:"msgType"`                           //消息类型，1：文本2：图片3：视频4：文件5：语音6：语言通话7：视频通话8：撤回消息9：回复消息10：引用消息11:@消息
	MsgPreview       string            `gorm:"size:64" json:"msgPreview"`         //消息预览
	Msg              ctype.Msg         `json:"msg"`                               //消息内容
	SystemMsg        *ctype.SystemMsg  `json:"systemMsg"`
}

func (chat GroupMsgModel) MsgPreviewMethod() string {
	if chat.SystemMsg != nil {
		switch chat.SystemMsg.Type {
		case 1:
			return "[系统消息]-该消息涉黄，已经被系统拦截"
		case 2:
			return "[系统消息]-该消息涉恐，已经被系统拦截"
		case 3:
			return "[系统消息]-该消息涉政，已经被系统拦截"
		case 4:
			return "[系统消息]-该消息不正当言论，已经被系统拦截"
		}
		return "[系统消息]"
	}
	return chat.Msg.MsgPreview()
}
