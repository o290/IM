package svc

import (
	"gorm.io/gorm"
	"server/core"
	"server/im_chat/chat_rpc/internal/config"
)

type ServiceContext struct {
	Config config.Config
	DB     *gorm.DB //注入
}

func NewServiceContext(c config.Config) *ServiceContext {
	mysqlDb := core.InitMysql(c.ChatMysql.DataSource)
	return &ServiceContext{
		Config: c,
		DB:     mysqlDb,
	}
}
