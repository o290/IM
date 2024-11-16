package models

import (
	"server/common/models"
	"time"
)

type ChatMsg struct {
	models.Model
	SendUserID uint       `json:"sendUserID"`
	RevUserID  uint       `json:"revUserID"`
	MsgType    int8       `json:"msgType"`    //消息类型，1：文本2：图片3：视频4：文件5：语音6：语言通话7：视频通话8：撤回消息9：回复消息10：引用消息
	MsgPreview string     `json:"msgPreview"` //消息预览
	Msg        Msg        `json:"msg"`        //消息内容
	SystemMsg  *SystemMsg `json:"systemMsg"`
}

// 系统提示
type SystemMsg struct {
	Type int8 `json:"type"` //违规类型：1：涉黄2：涉恐3：涉政4：不正当言论
}
type Msg struct {
	Type         int8          `json:"type"`
	Content      *string       `json:"content"` //文本消息
	ImgMsg       *ImgMsg       `json:"imgMsg"`
	VideoMsg     *VideoMsg     `json:"videoMsg"`
	FileMsg      *FileMsg      `json:"fileMsg"`
	VoiceMsg     *VoiceMsg     `json:"voiceMsg"`
	VoiceCallMsg *VoiceCallMsg `json:"voiceCallMsg"`
	VideoCall    *VideoMsg     `json:"videoCall"`
	WithdrawMsg  *WithdrawMsg  `json:"withdrawMsg"`
	ReplyMsg     *ReplyMsg     `json:"replyMsg"`
	QuoteMsg     *QuoteMsg     `json:"quoteMsg"`
}
type ImgMsg struct {
	Title string `json:"title"`
	Src   string `json:"src"`
}
type VideoMsg struct {
	Title string `json:"title"`
	Src   string `json:"src"`
	Time  int    `json:"time"` //时长单位秒
}
type FileMsg struct {
	Title string `json:"title"`
	Src   string `json:"src"`
	Size  int64  `json:"size"`
	Type  string `json:"type"`
}
type VoiceMsg struct {
	Title string `json:"title"`
	Src   string `json:"src"`
}
type VoiceCallMsg struct {
	StartTime time.Time `json:"startTime"`
	EndTime   time.Time `json:"endTime"`
	EndReason time.Time `json:"endReason"` //0：发起方挂断1：接收方挂断2：网络原因挂断3：未打通
}
type VideoCall struct {
	StartTime time.Time `json:"startTime"`
	EndTime   time.Time `json:"endTime"`
	EndReason time.Time `json:"endReason"` //0：发起方挂断1：接收方挂断2：网络原因挂断3：未打通
}
type WithdrawMsg struct {
	Content   string `json:"content"`   //提示词
	OriginMsg *Msg   `json:"originMsg"` //原消息
}
type ReplyMsg struct {
	MsgID   uint   `json:"msgID"`   //原消息id
	Content string `json:"content"` //回复的内容
	Msg     *Msg   `json:"msg"`
}
type QuoteMsg struct {
	MsgID   uint   `json:"msgID"`   //原消息id
	Content string `json:"content"` //回复的内容
	Msg     *Msg   `json:"msg"`
}
