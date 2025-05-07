package svc

import (
	"github.com/zeromicro/go-zero/zrpc"
	"gorm.io/gorm"
	"server/core"
	"server/im_file/file_api/internal/config"
	"server/im_user/user_rpc/types/user_rpc"
	"server/im_user/user_rpc/users"
)

type ServiceContext struct {
	Config  config.Config
	UserRpc user_rpc.UsersClient
	DB      *gorm.DB
}

func NewServiceContext(c config.Config) *ServiceContext {
	mysqlDb := core.InitMysql(c.Mysql.DataSource)
	return &ServiceContext{
		Config:  c,
		UserRpc: users.NewUsers(zrpc.MustNewClient(c.UserRpc)),
		DB:      mysqlDb,
	}
}
