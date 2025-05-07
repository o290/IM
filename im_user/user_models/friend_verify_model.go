package user_models

import (
	"server/common/models"
	"server/common/models/ctype"
)

// 好友验证表
type FriendVerifyModel struct {
	models.Model
	SendUserID           uint                        `json:"sendUserID"` //发起验证方
	SendUserModel        UserModel                   `gorm:"foreignKey:SendUserID " json:"-"`
	RevUserID            uint                        `json:"revUserID"` //接受验证方
	RevUserModel         UserModel                   `gorm:"foreignKey:RevUserID " json:"-"`
	Status               int8                        `json:"status"`
	SendStatus           int8                        `json:"sendStatus"`                        //发送方状态 4删除
	RevStatus            int8                        `json:"revStatus"`                         //接收方状态 0未操作，1:同意 2:拒绝 3：忽略 4删除
	AdditionalMessage    string                      `gorm:"size:128" json:"additionalMessage"` //附加消息
	VerificationQuestion *ctype.VerificationQuestion `json:"verificationQuestion"`
}
