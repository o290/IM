package ctype

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

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
	AtMsg        *AtMsg        `json:"atMsg"` //@消息，群聊
}

// Scan取出来的数据
func (m *Msg) Scan(val interface{}) error {
	return json.Unmarshal(val.([]byte), m)
}
func (m Msg) Value() (driver.Value, error) {
	b, err := json.Marshal(m)
	return string(b), err
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
type AtMsg struct {
	UserID  uint   `json:"userID"`
	Content string `json:"content"` //回复的内容
	Msg     *Msg   `json:"msg"`
}
