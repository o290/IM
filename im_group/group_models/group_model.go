package group_models

import (
	"server/common/models"
	"server/common/models/ctype"
)

type GroupModel struct {
	models.Model
	Title                string                      `gorm:"size:32" json:"title"`
	Abstract             string                      `gorm:"size:128" json:"abstract"`
	Avatar               string                      `gorm:"size:256" json:"avatar"`
	Creator              uint                        `json:"creator"`
	IsSearch             bool                        `json:"isSearch"` //是否可以被搜索到
	Verification         int8                        `json:"verification"`
	VerificationQuestion *ctype.VerificationQuestion `json:"verificationQuestion"`
	IsInvite             bool                        `json:"isInvite"`
	IsTemporarySession   bool                        `json:"isTemporarySession"`          //是否开启临时会话
	IsProhibition        bool                        `json:"isProhibition"`               //是否开启全员禁言
	Size                 int                         `json:"size"`                        //群规模10 20 100 200 1000 2000
	MemberList           []GroupMemberModel          `gorm:"foreignKey:GroupID" json:"-"` //群成员列表
}

// 问题个数
func (uc GroupModel) ProblemCount() (c int) {
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
