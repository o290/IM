package chat_models

import "server/common/models"

// 置顶用户表
type TopUserModel struct {
	models.Model
	UserID    uint `json:"userID"`
	TopUserID uint `json:"topUserID"`
}
