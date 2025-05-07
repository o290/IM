package user_models

import (
	"gorm.io/gorm"
	"server/common/models"
)

// 好友表
type FriendModel struct {
	models.Model
	SendUserID     uint      `json:"sendUserID"` //发起验证方
	SendUserModel  UserModel `gorm:"foreignKey:SendUserID" json:"-""`
	RevUserID      uint      `json:"revUserID"` //接受验证方
	RevUserModel   UserModel `gorm:"foreignKey:RevUserID" json:"-""`
	SendUserNotice string    `gorm:"size:128" json:"sendUserNotice"` //发送方对接收方的备注
	RevUserNotice  string    `gorm:"size:128" json:"revUserNotice"`  //接收方对发送方的备注
}

func (f *FriendModel) IsFriend(db *gorm.DB, A, B uint) bool {
	//判断是否是好友，并把结果即好友关系赋值给f
	err := db.Take(&f, "(send_user_id=? and rev_user_id=?) or (send_user_id=? and rev_user_id=?)", A, B, B, A).Error
	if err == nil {
		return true
	}
	return false
}

// 查询用户的好友
func (f *FriendModel) Friends(db *gorm.DB, userID uint) (list []FriendModel) {
	db.Find(&list, "send_user_id=? or rev_user_id=?", userID, userID)
	return
}
func (f *FriendModel) GetUserNotice(userID uint) string {
	if userID == f.SendUserID {
		//如果我是发起方
		return f.SendUserNotice
	}
	if userID == f.RevUserID {
		return f.RevUserNotice
	}
	return ""
}
