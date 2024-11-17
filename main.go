package main

import (
	"flag"
	"fmt"
	"server/core"
	"server/im_chat/chat_models"
	"server/im_group/group_models"
	"server/im_user/user_models"
)

type Options struct {
	DB bool
}

func main() {
	var opt Options
	flag.BoolVar(&opt.DB, "db", false, "db")
	flag.Parse()
	if opt.DB {
		db := core.InitMysql("root:qwer0209@tcp(127.0.0.1:3306)/im_server?charset=utf8mb4&parseTime=True&loc=Local")
		err := db.AutoMigrate(
			&user_models.UserModel{},
			&user_models.FriendModel{},
			&user_models.FriendVerifyModel{},
			&user_models.UserConfModel{},

			&chat_models.ChatModel{},

			&group_models.GroupModel{},
			&group_models.GroupMemberModel{},
			&group_models.GroupMsgModel{},
			&group_models.GroupVerifyModel{},
		)
		if err != nil {
			fmt.Println("表结构生成失败", err)
			return
		}
		fmt.Println("表结构生成成功")
	}
}
