package group_models

import (
	"server/common/models"
	"server/common/models/ctype"
)

// 群验证表
type GroupVerifyModel struct {
	models.Model
	GroupID              uint                        `json:"groupID"` //群id
	GroupModel           GroupModel                  `gorm:"foreignKey:GroupID" json:"-"`
	UserID               uint                        `json:"userID"`                           //需要加群或退群的用户id
	Status               int8                        `json:"status"`                           //0：未操作，1:同意 2:拒绝 3：忽略
	AdditionalMessage    string                      `gorm:"size:32" json:"additionalMessage"` //附加消息
	VerificationQuestion *ctype.VerificationQuestion `json:"verificationQuestion"`
	Type                 int8                        `json:"type"` //1:加群2：退群
}
