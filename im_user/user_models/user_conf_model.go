package user_models

import (
	"server/common/models"
	"server/common/models/ctype"
)

// 一些用户信息，但是不会经常使用到，所以放在user_conf中
type UserConfModel struct {
	models.Model
	UserID               uint                        `json:"userID"`
	UserModel            UserModel                   `gorm:"foreignKey:UserID" json:"_"`
	RecallMessage        *string                     `gorm:"size:32" json:"recallMessage"` //消息撤回的提示
	FriendOnline         bool                        `json:"friendOnline"`                 //好有上线提醒
	Sound                bool                        `json:"sound"`                        //提示声音
	SecureLink           bool                        `json:"secureLink"`
	SavePwd              bool                        `json:"savePwd"`
	SearchUser           int8                        `json:"searchUser"`   //别人查找自己的方式,0不允许别人查找到我，1通过用户号找到我，2可以通过昵称搜索到我
	Verification         int8                        `json:"verification"` //0不允许任何人添加 1 允许任何人添加 2 需要验证消息 3 需要回答问题 4 需要正确回答问题
	VerificationQuestion *ctype.VerificationQuestion `json:"verificationQuestion"`
	Online               bool                        `json:"online"` //是否在线
}

// 问题个数
func (uc UserConfModel) ProblemCount() (c int) {
	if uc.VerificationQuestion != nil {
		if uc.VerificationQuestion.Problem1 != nil {
			c += 1
		}
		if uc.VerificationQuestion.Problem2 != nil {
			c += 1
		}
		if uc.VerificationQuestion.Problem3 != nil {
			c += 1
		}
	}
	return c
}
