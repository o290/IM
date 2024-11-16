package models

import "server/common/models"

type UserModel struct {
	models.Model
	Pwd      string `json:"pwd"`
	Nickname string `json:"nickname"`
	Abstract string `json:"abstract"`
	Avatar   string `json:"avatar"`
	IP       string `json:"ip"`
	Addr     string `json:"addr"`
}
