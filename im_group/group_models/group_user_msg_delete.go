package group_models

import "server/common/models"

type GroupUserMsgDeleteModel struct {
	models.Model
	UserID  uint `json:"userID"`
	MsgID   uint `json:"msgID"`
	GroupID uint `json:"groupID"`
}
