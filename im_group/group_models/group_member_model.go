package group_models

import "server/common/models"

// 群成员表
type GroupMemberModel struct {
	models.Model
	GroupID         uint       `json:"groupID"`
	GroupModel      GroupModel `gorm:"foreignKey:GroupID" json:"-"`
	UserID          uint       `json:"userID"`
	MemberNickName  string     `gorm:"size:32" json:"memberNickName"` //群昵称
	Role            int        `json:"role"`                          //1：群主2：管理员3：普通成员
	ProhibitionTime *int       `json:"prohibitionTime"`               //禁言时长，单位分钟
}
