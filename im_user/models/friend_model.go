package models

import "server/common/models"

// 好友表
type FriendModel struct {
	models.Model
	SendUserID uint   `json:"sendUserID"` //发起验证方
	RevUserID  uint   `json:"revUserID"`  //接受验证方
	Notice     string `json:"notice"`     //备注
}
