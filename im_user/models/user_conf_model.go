package models

import "server/common/models"

type SearchWay int8
type VerifyWay int8

const (
	NotFind SearchWay = iota
	ByID
	ByNickname
)
const (
	NotAdd VerifyWay = iota
	Admmit
	ByVerify
	ByQuestion
	ByQuestionCorrect
)

// 一些用户信息，但是不会经常使用到，所以放在user_conf中
type UserConfModel struct {
	models.Model
	UserID               uint                  `json:"userID"`
	RecallMessage        *string               `json:"recallMessage"` //消息撤回的提示
	FriendOnline         bool                  `json:"friendOnline"`  //好有上线提醒
	Sound                bool                  `json:"sound"`         //提示声音
	SecureLink           bool                  `json:"secureLink"`
	SavePwd              bool                  `json:"savePwd"`
	SearchUser           SearchWay             `json:"searchUser"` //别人查找自己的方式
	FriendVerification   VerifyWay             `json:"friendVerification"`
	VerificationQuestion *VerificationQuestion `json:"verificationQuestion"`
	Online               bool                  `json:"online"` //是否在线
}
type VerificationQuestion struct {
	Problem1 *string `json:"problem1"`
	Problem2 *string `json:"problem2"`
	Problem3 *string `json:"problem3"`
	Answer1  *string `json:"answer1"`
	Answer2  *string `json:"answer2"`
	Answer3  *string `json:"answer3"`
}
