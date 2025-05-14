// package chat_models
//
// import (
//
//	"server/common/models"
//	"server/common/models/ctype"
//
// )
//
//	type ChatModel struct {
//		models.Model
//		SendUserID uint             `json:"sendUserID"`
//		RevUserID  uint             `json:"revUserID"`
//		MsgType    ctype.MsgType    `json:"msgType"`                   //消息类型，1：文本2：图片3：视频4：文件5：语音6：语言通话7：视频通话8：撤回消息9：回复消息10：引用消息
//		MsgPreview string           `gorm:"size:64" json:"msgPreview"` //消息预览
//		Msg        ctype.Msg        `json:"msg"`                       //消息内容
//		SystemMsg  *ctype.SystemMsg `json:"systemMsg"`
//	}
//
//	func (chat ChatModel) MsgPreviewMethod() string {
//		if chat.SystemMsg != nil {
//			switch chat.SystemMsg.Type {
//			case 1:
//				return "[系统消息]-该消息涉黄，已经被系统拦截"
//			case 2:
//				return "[系统消息]-该消息涉恐，已经被系统拦截"
//			case 3:
//				return "[系统消息]-该消息涉政，已经被系统拦截"
//			case 4:
//				return "[系统消息]-该消息不正当言论，已经被系统拦截"
//			}
//			return "[系统消息]"
//		}
//		return chat.Msg.MsgPreview()
//	}
package chat_models

import (
	"server/common/models"
	"server/common/models/ctype"
)

type ChatModel struct {
	models.Model
	SendUserID uint             `json:"sendUserID"` //发送者id
	RevUserID  uint             `json:"revUserID"`  //接收者id
	MsgType    ctype.MsgType    `json:"msgType"`    //消息类型
	Msg        ctype.Msg        `json:"msg"`        //消息内容
	SystemMsg  *ctype.SystemMsg `json:"systemMsg"`  //系统消息
	MsgPreview string           `json:"msgPreview"` //消息预览
	Status     int              `json:"status"`     //消息状态：0-未发送 1-已发送 2-已接收 3-已读
}

func (chat ChatModel) MsgPreviewMethod() string {
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
