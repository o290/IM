package chat_models

type UserChatDeleteModel struct {
	UserID uint `json:"userID"`
	ChatID uint `json:"chatID"` //聊天记录id
}
