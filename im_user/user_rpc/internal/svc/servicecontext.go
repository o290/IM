package svc

import (
	"gorm.io/gorm"
	"server/core"
	"server/im_user/user_rpc/internal/config"
)

type ServiceContext struct {
	Config config.Config
	DB     *gorm.DB //注入
}

func NewServiceContext(c config.Config) *ServiceContext {
	mysqlDb := core.InitMysql(c.UserMysql.DataSource)
	return &ServiceContext{
		Config: c,
		DB:     mysqlDb,
	}
}
